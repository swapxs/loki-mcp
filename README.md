# Loki MCP Server

[![CI](https://github.com/scottlepp/loki-mcp/workflows/CI/badge.svg)](https://github.com/scottlepp/loki-mcp/actions/workflows/ci.yml)

A Go-based server implementation for the Model Context Protocol (MCP) with Grafana Loki integration.

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Building and Running

Build and run the server:

```bash
# Build the server
go build -o loki-mcp-server ./cmd/server

# Run the server
./loki-mcp-server
```

Or run directly with Go:

```bash
go run ./cmd/server
```

The server communicates using stdin/stdout following the Model Context Protocol (MCP). This makes it suitable for use with Claude Desktop and other MCP-compatible clients. It does not run as an HTTP server on a port.

## Project Structure

```
.
├── cmd/
│   ├── server/       # MCP server implementation
│   └── client/       # Client for testing the MCP server
├── internal/
│   ├── handlers/     # Tool handlers
│   └── models/       # Data models
├── pkg/
│   └── utils/        # Utility functions and shared code
└── go.mod            # Go module definition
```

## MCP Server

The Loki MCP Server implements the Model Context Protocol (MCP) and provides the following tools:

### Loki Query Tool

The `loki_query` tool allows you to query Grafana Loki log data:

- Required parameters:
  - `query`: LogQL query string

- Optional parameters:
  - `url`: The Loki server URL (default: from LOKI_URL environment variable or http://localhost:3100)
  - `start`: Start time for the query (default: 1h ago)
  - `end`: End time for the query (default: now)
  - `limit`: Maximum number of entries to return (default: 100)

#### Environment Variables

The Loki query tool supports the following environment variables:

- `LOKI_URL`: Default Loki server URL to use if not specified in the request

### Testing the MCP Server

You can test the MCP server using the provided client:

```bash
# Build the client
go build -o loki-mcp-client ./cmd/client

# Loki query examples:
./loki-mcp-client loki_query "{job=\"varlogs\"}"
./loki-mcp-client loki_query "{job=\"varlogs\"}" "-1h" "now" 100

# Using environment variable:
export LOKI_URL="http://localhost:3100"
./loki-mcp-client loki_query "{job=\"varlogs\"}"
```

## Docker Support

You can build and run the MCP server using Docker:

```bash
# Build the Docker image
docker build -t loki-mcp-server .

# Run the server
docker run --rm -i loki-mcp-server
```

Alternatively, you can use Docker Compose:

```bash
# Build and run with Docker Compose
docker-compose up --build
```

### Local Testing with Loki

The project includes a complete Docker Compose setup to test Loki queries locally:

1. Start the Docker Compose environment:
   ```bash
   docker-compose up -d
   ```

   This will start:
   - A Loki server on port 3100
   - A Grafana instance on port 3000 (pre-configured with Loki as a data source)
   - A log generator container that sends sample logs to Loki
   - The Loki MCP server

2. Use the provided test script to query logs:
   ```bash
   # Run with default parameters (queries last 15 minutes of logs)
   ./test-loki-query.sh
   
   # Query for error logs
   ./test-loki-query.sh '{job="varlogs"} |= "ERROR"'
   
   # Specify a custom time range and limit
   ./test-loki-query.sh '{job="varlogs"}' '-1h' 'now' 50
   ```

3. Insert dummy logs for testing:
   ```bash
   # Insert 10 dummy logs with default settings
   ./insert-loki-logs.sh
   
   # Insert 20 logs with custom job and app name
   ./insert-loki-logs.sh --num 20 --job "custom-job" --app "my-app"
   
   # Insert logs with custom environment and interval
   ./insert-loki-logs.sh --env "production" --interval 0.5
   
   # Show help message
   ./insert-loki-logs.sh --help
   ```

4. Access the Grafana UI at http://localhost:3000 to explore logs visually.

## Server-Sent Events (SSE) Support

The server now supports two modes of communication:
1. Standard input/output (stdin/stdout) following the Model Context Protocol (MCP)
2. HTTP Server with Server-Sent Events (SSE) endpoint for integration with tools like n8n

The default port for the HTTP server is 8080, but can be configured using the `SSE_PORT` environment variable.

### Server Endpoints

When running in HTTP mode, the server exposes the following endpoints:

- SSE Endpoint: `http://localhost:8080/sse` - For real-time event streaming
- MCP Endpoint: `http://localhost:8080/mcp` - For MCP protocol messaging

### Using Docker with SSE

When running the server with Docker, make sure to expose port 8080:

```bash
# Build the Docker image
docker build -t loki-mcp-server .

# Run the server with port mapping
docker run -p 8080:8080 --rm -i loki-mcp-server
```

### n8n Integration

You can integrate the Loki MCP Server with n8n workflows:

1. Install the MCP Client Tools node in n8n

2. Configure the node with these parameters:
   - **SSE Endpoint**: `http://your-server-address:8080/sse` (replace with your actual server address)
   - **Authentication**: Choose appropriate authentication if needed
   - **Tools to Include**: Choose which Loki tools to expose to the AI Agent

3. Connect the MCP Client Tool node to an AI Agent node that will use the Loki querying capabilities

Example workflow:
Trigger → MCP Client Tool (Loki server) → AI Agent (Claude)

## Architecture

The Loki MCP Server uses a modular architecture:

- **Server**: The main MCP server implementation in `cmd/server/main.go`
- **Client**: A test client in `cmd/client/main.go` for interacting with the MCP server
- **Handlers**: Individual tool handlers in `internal/handlers/`
  - `loki.go`: Grafana Loki query functionality

## Using with Claude Desktop

You can use this MCP server with Claude Desktop to add Loki query tools. Follow these steps:

### Option 1: Using the Compiled Binary

1. Build the server:
```bash
go build -o loki-mcp-server ./cmd/server
```

2. Add the configuration to your Claude Desktop configuration file using `claude_desktop_config_binary.json`.

### Option 2: Using Go Run with a Shell Script

1. Make the script executable:
```bash
chmod +x run-mcp-server.sh
```

2. Add the configuration to your Claude Desktop configuration file using `claude_desktop_config_script.json`.

### Option 3: Using Docker (Recommended)

1. Build the Docker image:
```bash
docker build -t loki-mcp-server .
```

2. Add the configuration to your Claude Desktop configuration file using `claude_desktop_config_docker.json`.

### Configuration Details

The Claude Desktop configuration file is located at:
- On macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- On Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- On Linux: `~/.config/Claude/claude_desktop_config.json`

You can use one of the example configurations provided in this repository:
- `claude_desktop_config.json`: Generic template
- `claude_desktop_config_example.json`: Example using `go run` with the current path
- `claude_desktop_config_binary.json`: Example using the compiled binary
- `claude_desktop_config_script.json`: Example using a shell script (recommended for `go run`)
- `claude_desktop_config_docker.json`: Example using Docker (most reliable)

**Notes**:
- When using `go run` with Claude Desktop, you may need to set several environment variables in both the script and the configuration file:
  - `HOME`: The user's home directory
  - `GOPATH`: The Go workspace directory
  - `GOMODCACHE`: The Go module cache directory
  - `GOCACHE`: The Go build cache directory
  
  These are required to ensure Go can find its modules and build cache when run from Claude Desktop.

- Using Docker is the most reliable approach as it packages all dependencies and environment variables in a container.

Or create your own configuration:

```json
{
  "mcpServers": {
    "lokiserver": {
      "command": "path/to/loki-mcp-server",
      "args": [],
      "env": {
        "LOKI_URL": "http://localhost:3100"
      },
      "disabled": false,
      "autoApprove": ["loki_query"]
    }
  }
}
```

Make sure to replace `path/to/loki-mcp-server` with the absolute path to the built binary or source code.

4. Restart Claude Desktop.

5. You can now use the tools in Claude:
   - Loki query examples:
     - "Query Loki for logs with the query {job=\"varlogs\"}"
     - "Find error logs from the last hour in Loki using query {job=\"varlogs\"} |= \"ERROR\""
     - "Show me the most recent 50 logs from Loki with job=varlogs"

## Using with Cursor

You can also integrate the Loki MCP server with the Cursor editor. To do this, add the following configuration to your Cursor settings:

Docker configuration:

```json
{
  "mcpServers": {
    "loki-mcp-server": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "-e", "LOKI_URL=http://host.docker.internal:3100", "loki-mcp-server:latest"]
    }
  }
}
```

After adding this configuration, restart Cursor, and you'll be able to use the Loki query tool directly within the editor.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Running Tests

The project includes comprehensive unit tests and CI/CD workflows to ensure reliability:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run tests with race detection  
go test -race ./...
```

### Timestamp Bug Fix (Issue #3)

This project previously had a critical bug where timestamps were displayed as year 2262 instead of correct dates. This has been fixed and regression tests are in place:

- **Root Cause**: Loki returns timestamps in nanoseconds, but the code was incorrectly treating them as seconds and multiplying by 1,000,000,000
- **Fix**: Correctly handle nanosecond timestamps from Loki  
- **Testing**: Comprehensive tests ensure timestamps display correctly (e.g., 2024, 2023) instead of 2262
- **CI Protection**: Automated tests prevent regression of this critical bug

The tests specifically verify:
- ✅ Timestamps show correct years instead of 2262
- ✅ Multiple timestamp formats work correctly
- ✅ Invalid timestamps have proper fallback behavior
- ✅ Integration with real Loki instances
