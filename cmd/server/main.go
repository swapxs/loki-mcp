package main

import (
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/server"

	"github.com/scottlepper/loki-mcp-server/internal/handlers"
)

const (
	version = "0.1.0"
	// Default timeout values
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Loki MCP Server",
		version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Add Loki query tool
	lokiQueryTool := handlers.NewLokiQueryTool()
	s.AddTool(lokiQueryTool, handlers.HandleLokiQuery)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
