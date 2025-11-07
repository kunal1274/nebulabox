#!/bin/bash

# Workflow 03: Test Container Grouping
# This script tests container grouping functionality

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLI_BINARY="$PROJECT_ROOT/bin/nebulabox"

echo "=========================================="
echo "Workflow 03: Container Group Test"
echo "=========================================="

# Check if CLI binary exists
if [ ! -f "$CLI_BINARY" ]; then
    echo "❌ CLI binary not found at $CLI_BINARY"
    echo "   Run 'make build-cli-test' first"
    exit 1
fi

GROUP_NAME="test-group-$(date +%s)"

echo "📦 Creating container group..."
echo "   Group: $GROUP_NAME"

# Create group configuration
TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

cat > "$TEST_DIR/group.json" <<EOF
{
  "name": "$GROUP_NAME",
  "strategy": "frontend-backend",
  "containers": [
    {
      "name": "frontend",
      "image": "nginx:alpine",
      "ports": ["3000:80"]
    },
    {
      "name": "backend",
      "image": "alpine:latest",
      "ports": ["5000:5000"]
    }
  ],
  "networking": {
    "internal": true,
    "bridge": "${GROUP_NAME}-bridge"
  }
}
EOF

# Create group (when CLI command is available)
echo "   Note: Group creation via CLI will be available in next phase"
echo "   Group config saved to: $TEST_DIR/group.json"

echo ""
echo "✅ Workflow 03 complete: Group test prepared"

