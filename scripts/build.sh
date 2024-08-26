#!/bin/bash

set -e

# Print commands
set -x

# Go to the project root directory
cd "$(dirname "$0")/.."

# Ensure dependencies are up to date
go mod tidy

# Build the application
go build -o beatradar ./cmd/beatradar

echo "Build completed successfully. Binary: beatradar"

# Optionally, build Docker image
if [ "$1" == "--docker" ]; then
    docker build -t beatradar:latest .
    echo "Docker image built: beatradar:latest"
fi