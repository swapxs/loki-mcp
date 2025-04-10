package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// LokiResult represents the structure of Loki query results
type LokiResult struct {
	Status string   `json:"status"`
	Data   LokiData `json:"data"`
	Error  string   `json:"error,omitempty"`
}

// LokiData represents the data portion of Loki results
type LokiData struct {
	ResultType string      `json:"resultType"`
	Result     []LokiEntry `json:"result"`
}

// LokiEntry represents a single log entry from Loki
type LokiEntry struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"` // [timestamp, log line]
}

// Environment variable name for Loki URL
const EnvLokiURL = "LOKI_URL"

// Default Loki URL when environment variable is not set
const DefaultLokiURL = "http://localhost:3100"

// NewLokiQueryTool creates and returns a tool for querying Grafana Loki
func NewLokiQueryTool() mcp.Tool {
	// Get Loki URL from environment variable or use default
	lokiURL := os.Getenv(EnvLokiURL)
	if lokiURL == "" {
		lokiURL = DefaultLokiURL
	}

	return mcp.NewTool("loki_query",
		mcp.WithDescription("Run a query against Grafana Loki"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("LogQL query string"),
		),
		mcp.WithString("url",
			mcp.Description(fmt.Sprintf("Loki server URL (default: %s from %s env var)", lokiURL, EnvLokiURL)),
			mcp.DefaultString(lokiURL),
		),
		mcp.WithString("start",
			mcp.Description("Start time for the query (default: 1h ago)"),
		),
		mcp.WithString("end",
			mcp.Description("End time for the query (default: now)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of entries to return (default: 100)"),
		),
	)
}

// HandleLokiQuery handles Loki query tool requests
func HandleLokiQuery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	queryString := request.Params.Arguments["query"].(string)

	// Get Loki URL from request arguments, if not present check environment
	var lokiURL string
	if urlArg, ok := request.Params.Arguments["url"].(string); ok && urlArg != "" {
		lokiURL = urlArg
	} else {
		// Fallback to environment variable
		lokiURL = os.Getenv(EnvLokiURL)
		if lokiURL == "" {
			lokiURL = DefaultLokiURL
		}
	}

	// Set defaults for optional parameters
	start := time.Now().Add(-1 * time.Hour).Unix()
	end := time.Now().Unix()
	limit := 100

	// Override defaults if parameters are provided
	if startStr, ok := request.Params.Arguments["start"].(string); ok && startStr != "" {
		startTime, err := parseTime(startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %v", err)
		}
		start = startTime.Unix()
	}

	if endStr, ok := request.Params.Arguments["end"].(string); ok && endStr != "" {
		endTime, err := parseTime(endStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %v", err)
		}
		end = endTime.Unix()
	}

	if limitVal, ok := request.Params.Arguments["limit"].(float64); ok {
		limit = int(limitVal)
	}

	// Build query URL
	queryURL, err := buildLokiQueryURL(lokiURL, queryString, start, end, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to build query URL: %v", err)
	}

	// Execute query
	result, err := executeLokiQuery(ctx, queryURL)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}

	// Format results
	formattedResult, err := formatLokiResults(result)
	if err != nil {
		return nil, fmt.Errorf("failed to format results: %v", err)
	}

	return mcp.NewToolResultText(formattedResult), nil
}

// parseTime parses a time string in various formats
func parseTime(timeStr string) (time.Time, error) {
	// Handle "now" keyword
	if timeStr == "now" {
		return time.Now(), nil
	}

	// Handle relative time strings like "-1h", "-30m"
	if len(timeStr) > 0 && timeStr[0] == '-' {
		duration, err := time.ParseDuration(timeStr)
		if err == nil {
			return time.Now().Add(duration), nil
		}
	}

	// Try parsing as RFC3339
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return t, nil
	}

	// Try other common formats
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported time format: %s", timeStr)
}

// buildLokiQueryURL constructs the Loki query URL
func buildLokiQueryURL(baseURL, query string, start, end int64, limit int) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Add path for Loki query API only if not already included
	if !strings.Contains(u.Path, "loki/api/v1") {
		if u.Path == "" || u.Path == "/" {
			u.Path = "/loki/api/v1/query_range"
		} else {
			u.Path = fmt.Sprintf("%s/loki/api/v1/query_range", u.Path)
		}
	} else {
		// If path already contains loki/api/v1, just append query_range if not present
		if !strings.HasSuffix(u.Path, "query_range") {
			u.Path = fmt.Sprintf("%s/query_range", u.Path)
		}
	}

	// Add query parameters
	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start))
	q.Set("end", fmt.Sprintf("%d", end))
	q.Set("limit", fmt.Sprintf("%d", limit))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// executeLokiQuery sends the HTTP request to Loki
func executeLokiQuery(ctx context.Context, queryURL string) (*LokiResult, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	// Execute request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result LokiResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Check for Loki errors
	if result.Status == "error" {
		return nil, fmt.Errorf("Loki error: %s", result.Error)
	}

	return &result, nil
}

// formatLokiResults formats the Loki query results into a readable string
func formatLokiResults(result *LokiResult) (string, error) {
	if len(result.Data.Result) == 0 {
		return "No logs found matching the query", nil
	}

	var output string
	output = fmt.Sprintf("Found %d streams:\n\n", len(result.Data.Result))

	for i, entry := range result.Data.Result {
		// Format stream labels
		streamInfo := "Stream "
		if len(entry.Stream) > 0 {
			streamInfo += "("
			first := true
			for k, v := range entry.Stream {
				if !first {
					streamInfo += ", "
				}
				streamInfo += fmt.Sprintf("%s=%s", k, v)
				first = false
			}
			streamInfo += ")"
		}

		output += fmt.Sprintf("%s %d:\n", streamInfo, i+1)

		// Format log entries
		for _, val := range entry.Values {
			if len(val) >= 2 {
				// Parse timestamp
				ts, err := strconv.ParseFloat(val[0], 64)
				if err == nil {
					// Convert to time
					timestamp := time.Unix(0, int64(ts*1000000000))
					output += fmt.Sprintf("[%s] %s\n", timestamp.Format(time.RFC3339), val[1])
				} else {
					output += fmt.Sprintf("[%s] %s\n", val[0], val[1])
				}
			}
		}
		output += "\n"
	}

	return output, nil
}
