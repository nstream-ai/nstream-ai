#!/bin/bash

# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Create proto output directory
mkdir -p proto/auth proto/cluster

# Generate Go code from proto files
protoc --go_out=. --go_opt=module=github.com/nstream-ai/nstream-ai-mothership \
    --go-grpc_out=. --go-grpc_opt=module=github.com/nstream-ai/nstream-ai-mothership \
    proto/auth.proto

protoc --go_out=. --go_opt=module=github.com/nstream-ai/nstream-ai-mothership \
    --go-grpc_out=. --go-grpc_opt=module=github.com/nstream-ai/nstream-ai-mothership \
    proto/cluster.proto

# Move generated files to the correct location
# mv proto/gen/github.com/nstream-ai/nstream-ai-mothership/proto/* proto/
# rm -rf proto/gen 