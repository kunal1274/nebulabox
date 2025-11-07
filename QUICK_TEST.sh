#!/bin/bash

# Quick Test Script - Run all basic tests in sequence
# This is a convenience script for quick verification

set -e

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           NebulaBox Quick Test Suite                         ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

# Test function
test_step() {
    local name="$1"
    local command="$2"
    
    echo -n "Testing: $name... "
    if eval "$command" > /tmp/test_output.log 2>&1; then
        echo -e "${GREEN}✅ PASS${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ FAIL${NC}"
        ((FAILED++))
        return 1
    fi
}

# Test 1: Build CLI
echo "Step 1: Building CLI..."
if make build-cli-test > /dev/null 2>&1; then
    echo -e "${GREEN}✅ CLI built successfully${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ CLI build failed${NC}"
    ((FAILED++))
    exit 1
fi
echo ""

# Test 2: Version command
test_step "Version command" "./bin/nebulabox version"
echo ""

# Test 3: Help command
test_step "Help command" "./bin/nebulabox --help"
echo ""

# Test 4: List containers
test_step "List containers (ps)" "./bin/nebulabox ps"
echo ""

# Test 5: Group list
test_step "Group list" "./bin/nebulabox group list"
echo ""

# Test 6: Engine compilation
test_step "Engine compilation" "go build ./internal/engine/..."
echo ""

# Test 7: Engine unit tests
echo "Step 7: Running engine unit tests..."
if go test ./internal/engine/... -v > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Engine tests pass${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  Some engine tests may have issues${NC}"
    ((FAILED++))
fi
echo ""

# Test 8: CLI compilation
test_step "CLI compilation" "go build ./internal/cli/..."
echo ""

# Summary
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                      Test Summary                             ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Run interactive demo: ./scripts/cli/demo-poc.sh"
    echo "  2. Run detailed tests: ./scripts/testing/test-engine-basic.sh"
    echo "  3. See manual guide: cat MANUAL_TESTING_GUIDE.md"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some tests failed (may be expected in POC phase)${NC}"
    echo ""
    echo "Check MANUAL_TESTING_GUIDE.md for detailed testing instructions"
    exit 1
fi

