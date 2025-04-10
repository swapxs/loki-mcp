# Loki MCP Server

A Go-based server implementation for the Model Context Protocol (MCP).

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

By default, the server runs on port 8080. You can change this by setting the PORT environment variable:

```bash
PORT=3000 go run ./cmd/server
```

## Project Structure

```
.
├── cmd/
│   ├── server/       # MCP server implementation
│   └── client/       # Client for testing the MCP server
├── internal/
│   ├── handlers/     # HTTP handlers
│   └── models/       # Data models
├── pkg/
│   └── utils/        # Utility functions and shared code
└── go.mod            # Go module definition
```

## MCP Server

The Loki MCP Server implements the Model Context Protocol (MCP) and provides a `calculate` tool that can perform basic arithmetic operations:

- add: Addition (x + y)
- subtract: Subtraction (x - y)
- multiply: Multiplication (x * y)
- divide: Division (x / y)

### Testing the MCP Server

You can test the MCP server using the provided client:

```bash
# Run the client with default values (add 5 3)
make run-client

# Or run the client with custom values
go run ./cmd/client <operation> <x> <y>

# Examples:
go run ./cmd/client add 10 5
go run ./cmd/client subtract 10 5
go run ./cmd/client multiply 10 5
go run ./cmd/client divide 10 5
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

## Using with Claude Desktop

You can use this MCP server with Claude Desktop to add a calculator tool to Claude. Follow these steps:

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
    "calculator": {
      "command": "path/to/loki-mcp-server",
      "args": [],
      "env": {},
      "disabled": false,
      "autoApprove": ["calculate"]
    }
  }
}
```

Make sure to replace `path/to/loki-mcp-server` with the absolute path to the built binary or source code.

4. Restart Claude Desktop.

5. You can now use the calculator tool in Claude by asking it to perform calculations:
   - "Calculate 5 + 3"
   - "What is 10 * 7?"
   - "Divide 20 by 4"

## License

This project is licensed under the MIT License - see the LICENSE file for details.
