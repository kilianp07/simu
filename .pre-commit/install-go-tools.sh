#!/bin/bash

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

if ! command_exists goimports; then
  echo "Installing goimports..."
  go get -u golang.org/x/tools/cmd/goimports
fi

if ! command_exists gocyclo; then
  echo "Installing gocyclo..."
  go get -u github.com/fzipp/gocyclo
fi

if ! command_exists golangci-lint; then
  echo "Installing golangci-lint..."
  go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
fi

if ! command_exists gocritic; then
  echo "Installing gocritic..."
  go get -u github.com/go-critic/go-critic/cmd/gocritic
fi

goimports --version
gocyclo --version
golangci-lint --version
gocritic --version
