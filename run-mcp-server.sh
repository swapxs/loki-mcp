#!/bin/bash

# Set environment variables
export HOME="/Users/scottlepper"
export GOPATH="/Users/scottlepper/go"
export GOMODCACHE="/Users/scottlepper/go/pkg/mod"
export GOCACHE="/Users/scottlepper/Library/Caches/go-build"

# Change to the project directory
cd /Users/scottlepper/code/loki-mcp-server

# Run the server
go run ./cmd/server/main.go
