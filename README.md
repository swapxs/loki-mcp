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

## License

This project is licensed under the MIT License - see the LICENSE file for details.
