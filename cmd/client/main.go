package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Request represents a JSON-RPC request
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// Response represents a JSON-RPC response
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ErrorObject    `json:"error,omitempty"`
}

// ErrorObject represents a JSON-RPC error object
type ErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToolResult represents the result of a call_tool request
type ToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError"`
}

// ContentItem represents a content item in a tool result
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: client <operation> <x> <y>")
		fmt.Println("Operations: add, subtract, multiply, divide")
		fmt.Println("Example: client add 5 3")
		os.Exit(1)
	}

	operation := os.Args[1]
	x, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatalf("Invalid number for x: %v", err)
	}

	y, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		log.Fatalf("Invalid number for y: %v", err)
	}

	// Start the MCP server in a separate process
	cmd := exec.Command("go", "run", "./cmd/server/main.go")

	// Connect stdin and stdout to the MCP server
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Set up stderr to be displayed
	cmd.Stderr = os.Stderr

	// Start the server
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Create a reader for the stdout
	reader := bufio.NewReader(stdout)

	// Call the calculate tool
	fmt.Printf("Calling calculate tool with operation=%s, x=%.2f, y=%.2f\n", operation, x, y)

	// Create the request
	req := Request{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "calculate",
			"arguments": map[string]interface{}{
				"operation": operation,
				"x":         x,
				"y":         y,
			},
		},
	}

	// Marshal the request to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	fmt.Printf("Sending request: %s\n", string(reqJSON))

	// Send the request to the server
	_, err = stdin.Write(append(reqJSON, '\n'))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	// Read the response from the server
	respJSON, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Received response: %s\n", string(respJSON))

	// Unmarshal the response
	var resp Response
	if err := json.Unmarshal(respJSON, &resp); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check for errors
	if resp.Error != nil {
		log.Fatalf("Error from server: %s", resp.Error.Message)
	}

	// Unmarshal the result
	var result ToolResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		log.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Process the result
	if result.IsError {
		fmt.Printf("Error: %s\n", result.Content[0].Text)
	} else {
		// Extract the result from the text response
		resultText := result.Content[0].Text
		resultText = strings.TrimSpace(resultText)

		fmt.Printf("Result: %s %s %s = %s\n", os.Args[2], getOperationSymbol(operation), os.Args[3], resultText)
	}

	// Terminate the server
	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Failed to kill server process: %v", err)
	}
}

func getOperationSymbol(operation string) string {
	switch operation {
	case "add":
		return "+"
	case "subtract":
		return "-"
	case "multiply":
		return "*"
	case "divide":
		return "/"
	default:
		return operation
	}
}
