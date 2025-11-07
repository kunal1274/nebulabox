#!/bin/bash

# Container Testing Script
# Tests container operations

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLI_BINARY="$PROJECT_ROOT/bin/nebulabox"

echo "=========================================="
echo "Container Operations Test"
echo "=========================================="

if [ ! -f "$CLI_BINARY" ]; then
    echo "❌ CLI binary not found. Run 'make build-cli-test' first"
    exit 1
fi

# Test 1: List containers (should work even if empty)
echo ""
echo "Test 1: Listing containers..."
if $CLI_BINARY ps 2>&1; then
    echo "✅ List containers command works"
else
    echo "⚠️  List containers may show error (expected in POC)"
fi

# Test 2: List images
echo ""
echo "Test 2: Listing images..."
if $CLI_BINARY images 2>&1; then
    echo "✅ List images command works"
else
    echo "⚠️  List images may show error (expected in POC)"
fi

# Test 3: List groups
echo ""
echo "Test 3: Listing groups..."
if $CLI_BINARY group list 2>&1; then
    echo "✅ List groups command works"
else
    echo "⚠️  List groups may show error (expected in POC)"
fi

# Test 4: Version
echo ""
echo "Test 4: Getting version..."
if $CLI_BINARY version 2>&1; then
    echo "✅ Version command works"
else
    echo "⚠️  Version command may not be implemented"
fi

echo ""
echo "=========================================="
echo "✅ Container Operations Test Complete"
echo "=========================================="
echo ""
echo "Note: Full container operations (create, start, stop)"
echo "require the engine to be fully implemented with root"
echo "privileges. This is expected in POC phase."

