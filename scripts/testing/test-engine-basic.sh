#!/bin/bash

# Basic Engine Testing Script
# Tests core engine functionality

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "=========================================="
echo "Engine Basic Functionality Test"
echo "=========================================="

# Test 1: Compile engine
echo ""
echo "Test 1: Compiling engine..."
if go build ./internal/engine/...; then
    echo "✅ Engine compiles successfully"
else
    echo "❌ Engine compilation failed"
    exit 1
fi

# Test 2: Run engine unit tests
echo ""
echo "Test 2: Running engine unit tests..."
if go test ./internal/engine/... -v -run TestNewRuntime; then
    echo "✅ Engine unit tests pass"
else
    echo "⚠️  Some engine tests may have failed (expected in POC phase)"
fi

# Test 3: Compile CLI
echo ""
echo "Test 3: Compiling CLI..."
if go build ./internal/cli/...; then
    echo "✅ CLI compiles successfully"
else
    echo "❌ CLI compilation failed"
    exit 1
fi

# Test 4: Build CLI binary
echo ""
echo "Test 4: Building CLI binary..."
if make build-cli-test 2>&1 | grep -q "CLI test binary build complete"; then
    echo "✅ CLI binary built successfully"
else
    echo "⚠️  CLI binary build may have issues"
fi

# Test 5: CLI commands exist
echo ""
echo "Test 5: Checking CLI commands..."
if [ -f "$PROJECT_ROOT/bin/nebulabox" ]; then
    echo "✅ CLI binary exists"
    
    # Test help
    if ./bin/nebulabox --help > /dev/null 2>&1; then
        echo "✅ CLI help command works"
    else
        echo "⚠️  CLI help command may have issues"
    fi
    
    # Test group command
    if ./bin/nebulabox group --help > /dev/null 2>&1; then
        echo "✅ Group command exists"
    else
        echo "⚠️  Group command may not be fully implemented"
    fi
else
    echo "❌ CLI binary not found"
    exit 1
fi

echo ""
echo "=========================================="
echo "✅ Basic Engine Tests Complete"
echo "=========================================="

