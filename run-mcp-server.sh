#!/bin/bash

# Set environment variables
export HOME="/Users/<username>"
export GOPATH="/Users/<username>/go"
export GOMODCACHE="/Users/<username>/go/pkg/mod"
export GOCACHE="/Users/<username>/Library/Caches/go-build"

# Change to the project directory
cd /Users/<username>/code/loki-mcp-server

# Run the server
go run ./cmd/server/main.go
