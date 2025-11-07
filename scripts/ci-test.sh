#!/bin/bash

# NebulaBox CI Test Script
# Comprehensive test runner for CI environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "🚀 NebulaBox CI Test Suite"
echo "=========================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track test results
TESTS_PASSED=0
TESTS_FAILED=0

run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -n "Running $test_name... "
    if eval "$test_command" > /tmp/test-output.log 2>&1; then
        echo -e "${GREEN}✓${NC}"
        ((TESTS_PASSED++))
        return 0
    else
        echo -e "${RED}✗${NC}"
        echo "Error output:"
        cat /tmp/test-output.log | tail -20
        ((TESTS_FAILED++))
        return 1
    fi
}

# Step 1: Verify Go environment
echo "📦 Verifying Go environment..."
go version
go env GOROOT GOPATH
echo ""

# Step 2: Download dependencies
echo "📥 Downloading dependencies..."
go mod download
go mod verify
echo ""

# Step 3: Format check
echo "📝 Checking code formatting..."
if gofmt -l . | grep -q .; then
    echo -e "${YELLOW}⚠️  Some files are not formatted:${NC}"
    gofmt -l . | head -10
    echo "Run 'go fmt ./...' to fix"
else
    echo -e "${GREEN}✓ Code is properly formatted${NC}"
fi
echo ""

# Step 4: Run unit tests
echo "🧪 Running unit tests..."
run_test "Unit Tests" "make test-unit"
echo ""

# Step 5: Run integration tests
echo "🔗 Running integration tests..."
if run_test "Integration Tests" "make test-integration"; then
    echo ""
else
    echo -e "${YELLOW}⚠️  Integration tests failed (may require external services)${NC}"
    echo ""
fi

# Step 6: Run benchmarks (non-blocking)
echo "📊 Running benchmarks..."
if run_test "Benchmarks" "make benchmark-api"; then
    echo ""
else
    echo -e "${YELLOW}⚠️  Benchmarks failed (non-blocking)${NC}"
    echo ""
fi

# Step 7: Build check
echo "🔨 Building application..."
if run_test "Build API" "make build-api"; then
    echo ""
else
    echo -e "${RED}✗ Build failed!${NC}"
    exit 1
fi

if run_test "Build Test Runner" "make build-test"; then
    echo ""
else
    echo -e "${RED}✗ Test runner build failed!${NC}"
    exit 1
fi

# Step 8: Generate coverage
echo "📈 Generating test coverage..."
if run_test "Coverage Report" "make test-unit-coverage"; then
    if [ -f coverage.html ]; then
        echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
    fi
    echo ""
fi

# Summary
echo "=========================="
echo "📊 Test Summary"
echo "=========================="
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}Failed: $TESTS_FAILED${NC}"
    echo ""
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
fi

