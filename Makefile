# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=loki-mcp-server
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/server

.PHONY: all build clean test run deps tidy help

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

test:
	$(GOTEST) -v ./...

run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

run-client:
	go run ./cmd/client add 5 3

deps:
	$(GOGET) -v -t ./...

tidy:
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

help:
	@echo "Make commands:"
	@echo "  build       - Build the binary"
	@echo "  clean       - Remove binary and cache files"
	@echo "  test        - Run tests"
	@echo "  run         - Build and run the binary"
	@echo "  deps        - Get dependencies"
	@echo "  tidy        - Tidy go.mod file"
	@echo "  build-linux - Cross-compile for Linux"
	@echo "  help        - Display this help message"
