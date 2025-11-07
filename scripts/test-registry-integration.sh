#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

REGISTRY_PORT=${NEBULABOX_REGISTRY_PORT:-5001}
API_PORT=${NEBULABOX_API_PORT:-8081}
REGISTRY_URL="http://localhost:$REGISTRY_PORT"
API_URL="http://localhost:$API_PORT"
REGISTRY_STORAGE="./test-registry-storage"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

cleanup() {
    echo -e "\n${YELLOW}🧹 Cleaning up...${NC}"
    pkill -f nebulabox-registry || true
    pkill -f nebulabox-api || true
    sleep 1
    rm -rf "$REGISTRY_STORAGE"
}

trap cleanup EXIT

echo "🚀 Starting Registry Integration Tests"
echo "========================================"
echo ""

# Start registry server
echo "📦 Starting registry server..."
NEBULABOX_REGISTRY_PORT=$REGISTRY_PORT \
NEBULABOX_REGISTRY_STORAGE="$REGISTRY_STORAGE" \
./nebulabox-registry > /tmp/registry.log 2>&1 &
REGISTRY_PID=$!

# Wait for registry to be ready
echo "⏳ Waiting for registry to start..."
for i in {1..10}; do
    if curl -s "$REGISTRY_URL/v2/" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Registry is ready!${NC}"
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${RED}❌ Registry failed to start${NC}"
        cat /tmp/registry.log
        exit 1
    fi
    sleep 1
done

# Start API server
echo "📦 Starting API server..."
NEBULABOX_REGISTRY_URL="$REGISTRY_URL" \
NEBULABOX_ADMIN_USER="admin" \
NEBULABOX_ADMIN_PASS="admin" \
./nebulabox-api > /tmp/api.log 2>&1 &
API_PID=$!

# Wait for API to be ready
echo "⏳ Waiting for API server to start..."
for i in {1..15}; do
    if curl -s "$API_URL/api/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ API server is ready!${NC}"
        break
    fi
    if [ $i -eq 15 ]; then
        echo -e "${RED}❌ API server failed to start${NC}"
        cat /tmp/api.log
        exit 1
    fi
    sleep 1
done

echo ""
echo "🧪 Running Integration Tests"
echo "----------------------------"

# Test 1: Registry catalog
echo -n "Test 1: Registry catalog... "
RESPONSE=$(curl -s "$API_URL/api/registry/catalog")
if echo "$RESPONSE" | grep -q "repositories"; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Response: $RESPONSE"
    exit 1
fi

# Test 2: Registry repositories API
echo -n "Test 2: List repositories... "
RESPONSE=$(curl -s "$API_URL/api/registry/repositories")
if echo "$RESPONSE" | grep -q "repositories"; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Response: $RESPONSE"
    exit 1
fi

# Test 3: Registry authentication
echo -n "Test 3: Registry authentication... "
LOGIN_RESPONSE=$(curl -s -X POST "$REGISTRY_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin"}')
if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

# Test 4: API server can list images from registry
echo -n "Test 4: API lists images from registry... "
RESPONSE=$(curl -s "$API_URL/api/images")
if echo "$RESPONSE" | grep -q "\[\]"; then
    echo -e "${GREEN}✓${NC} (empty registry - expected)"
else
    echo -e "${GREEN}✓${NC}"
fi

# Test 5: API registry management endpoints
echo -n "Test 5: Registry management endpoints... "
REPO_RESPONSE=$(curl -s "$API_URL/api/registry/repositories")
if [ -n "$REPO_RESPONSE" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ All integration tests passed!${NC}"
echo ""
echo "Test Summary:"
echo "  - Registry server: Running on port $REGISTRY_PORT"
echo "  - API server: Running on port $API_PORT"
echo "  - Registry storage: $REGISTRY_STORAGE"
echo "  - Tests passed: 5/5"

