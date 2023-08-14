#!/bin/bash

# Build script for Go server

# Setup
APP="chatproxy"

# Helper functions
info() {
  echo -e "\033[1;34m$@\033[0m"
}

error() {
  echo -e "\033[1;31m$@\033[0m" >&2
}

# Validate deps
info "Validating dependencies..."
go mod verify || error "Dependency validation failed!"

# Build binary
info "Building binary..."
env GOOS=linux GOARCH=amd64 go build -o ${APP} server.go || error "Build failed!"

# Run tests
info "Running tests..."
go test ./... || error "Tests failed!"

# Build Docker image
#info "Building Docker image..."
#docker build -t ${APP}:latest . || error "Docker build failed!"

info "Build completed!"

info "Copying to server"
scp ./${APP} root@whereyogi.com:/mnt/volume_nyc1_01/new_chatproxy
info "remove local build"
rm ./${APP}