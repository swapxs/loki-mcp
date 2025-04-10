#!/bin/bash

# Default Loki URL
export LOKI_URL="http://localhost:3100"

# Build the MCP server and client
echo "Building MCP server and client..."
go build -o loki-mcp-server ./cmd/server
go build -o loki-mcp-client ./cmd/client

# Check if Docker Compose is running
if ! docker ps | grep -q loki; then
    echo "Starting Docker Compose environment..."
    docker-compose up -d
    
    # Wait for Loki to start
    echo "Waiting for Loki to start..."
    sleep 10
fi

# Default query parameters
QUERY="{job=\"varlogs\"}"
START="-15m"
END="now"
LIMIT=100

# Parse command line arguments
if [ $# -ge 1 ]; then
    QUERY="$1"
fi

if [ $# -ge 2 ]; then
    START="$2"
fi

if [ $# -ge 3 ]; then
    END="$3"
fi

if [ $# -ge 4 ]; then
    LIMIT="$4"
fi

echo "Querying Loki with:"
echo "  URL:    $LOKI_URL (from environment variable)"
echo "  Query:  $QUERY"
echo "  Start:  $START"
echo "  End:    $END"
echo "  Limit:  $LIMIT"
echo ""

# Run the client without explicitly specifying the URL (will use environment)
./loki-mcp-client loki_query "$QUERY" "$START" "$END" "$LIMIT" 