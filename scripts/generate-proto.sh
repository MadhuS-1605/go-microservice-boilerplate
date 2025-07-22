#!/bin/bash

set -e

echo "Generating protobuf files..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "protoc is required but not installed. Please install Protocol Buffers compiler."
    echo "Visit: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Check if Go protobuf plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing Go protobuf plugins..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate common proto files
echo "Generating common protobuf files..."
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  internal/proto/common/common.proto


# Generate user proto files
echo "Generating user protobuf files..."
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  -I. \
  internal/proto/user/user.proto

# Generate product proto files
echo "Generating product protobuf files..."
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  -I. \
  internal/proto/product/product.proto

echo "Protobuf generation complete!"
