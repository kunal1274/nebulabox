#!/bin/bash

# Workflow 01: Test Build Workflow
# This script tests the basic build functionality of NebulaBox CLI

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLI_BINARY="$PROJECT_ROOT/bin/nebulabox"

echo "=========================================="
echo "Workflow 01: Build Test"
echo "=========================================="

# Check if CLI binary exists
if [ ! -f "$CLI_BINARY" ]; then
    echo "❌ CLI binary not found at $CLI_BINARY"
    echo "   Run 'make build-cli-test' first"
    exit 1
fi

# Create a test BuildSpec
TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

cat > "$TEST_DIR/buildspec.json" <<EOF
{
  "version": "1.0",
  "name": "test-app",
  "tag": "test-app:latest",
  "base": {
    "image": "alpine",
    "tag": "latest"
  },
  "workdir": "/app",
  "steps": [
    {
      "type": "run",
      "command": "echo 'Hello from NebulaBox'",
      "comment": "Test command"
    }
  ]
}
EOF

echo "📦 Building image from BuildSpec..."
echo "   BuildSpec: $TEST_DIR/buildspec.json"

# Build the image
if $CLI_BINARY build -f "$TEST_DIR/buildspec.json" -t test-app:latest; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed"
    exit 1
fi

echo ""
echo "✅ Workflow 01 complete: Build test passed"

