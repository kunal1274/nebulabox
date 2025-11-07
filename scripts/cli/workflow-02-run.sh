#!/bin/bash

# Workflow 02: Test Run Workflow
# This script tests running containers with NebulaBox CLI

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLI_BINARY="$PROJECT_ROOT/bin/nebulabox"

echo "=========================================="
echo "Workflow 02: Run Test"
echo "=========================================="

# Check if CLI binary exists
if [ ! -f "$CLI_BINARY" ]; then
    echo "❌ CLI binary not found at $CLI_BINARY"
    echo "   Run 'make build-cli-test' first"
    exit 1
fi

CONTAINER_NAME="test-container-$(date +%s)"
IMAGE_NAME="alpine:latest"

echo "🚀 Running container..."
echo "   Image: $IMAGE_NAME"
echo "   Name: $CONTAINER_NAME"

# Pull image first (if needed)
echo "📥 Pulling image..."
$CLI_BINARY pull "$IMAGE_NAME" || echo "   (Image pull skipped - may already exist)"

# Run container
if $CLI_BINARY run "$IMAGE_NAME" --name "$CONTAINER_NAME" -d; then
    echo "✅ Container started successfully!"
    
    # Wait a bit
    sleep 2
    
    # List containers
    echo ""
    echo "📋 Listing containers..."
    $CLI_BINARY ps
    
    # Stop container
    echo ""
    echo "🛑 Stopping container..."
    if $CLI_BINARY stop "$CONTAINER_NAME"; then
        echo "✅ Container stopped successfully!"
    else
        echo "⚠️  Container stop failed (may have already exited)"
    fi
else
    echo "❌ Container run failed"
    exit 1
fi

echo ""
echo "✅ Workflow 02 complete: Run test passed"

