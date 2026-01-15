#!/bin/bash

echo "Building autocommit-cli..."
go build -o bin/autocommit-cli cmd/autocommit-cli/main.go

if [ $? -eq 0 ]; then
    echo "Build successful! Executable is in ./bin/autocommit-cli"
    echo "You can run it using: ./bin/autocommit-cli"
else
    echo "Build failed."
fi