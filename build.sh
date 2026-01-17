#!/bin/bash
echo "Building autocommit-cli..."
go build -o autocommit-cli cmd/autocommit-cli/main.go
echo "Installing autocommit-cli to /usr/local/bin..."
sudo mv autocommit-cli /usr/local/bin/
echo "Installation complete."
