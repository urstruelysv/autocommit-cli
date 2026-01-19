#!/bin/bash
echo "Building autocommit-cli..."
# Ensure the bin directory exists
mkdir -p bin
# Build a statically linked binary
go build -ldflags "-s -w -extldflags '-static'" -o bin/autocommit-cli cmd/autocommit-cli/main.go
echo "Build complete. Binary located at bin/autocommit-cli"
