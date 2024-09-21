#!/bin/bash

# Test Build
echo "Running Tests..."
go test ./...
echo "Tests complete."

# Linux Build
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/server-linux-amd64
echo "Linux build complete."

# Windows Build
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o build/server-windows-amd64.exe
echo "Windows build complete."

# Mac Build (Intel)
echo "Building for Mac (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o build/server-darwin-amd64
echo "Mac (Intel) build complete."

# Mac Build (Apple Silicon)
echo "Building for Mac (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o build/server-darwin-arm64
echo "Mac (Apple Silicon) build complete."
