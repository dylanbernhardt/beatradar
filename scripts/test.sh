#!/bin/bash

set -e

# Print commands
set -x

# Go to the project root directory
cd "$(dirname "$0")/.."

# Run tests
go test -v ./...

# Run tests with race condition checking
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "Tests completed. Coverage report generated: coverage.html"