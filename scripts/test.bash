#!/bin/bash
# Exit on first error
set -e

# Navigate to the parent directory where the project is
cd "$(dirname "$0")/.."

# Clear cache and run test
go clean -testcache
go test ./... -v -race -covermode=atomic
