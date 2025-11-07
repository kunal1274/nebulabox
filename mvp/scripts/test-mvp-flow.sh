#!/bin/bash

# MVP End-to-End Testing Script
# Tests the complete flow: Build → Create → Start → Verify

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
MVP_DIR="$PROJECT_ROOT/mvp/app/single-container"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

API_URL="http://localhost:8081/api"
BUILDSPEC_PATH="$MVP_DIR/buildspec.json"

echo -e "${BLUE}╔═══════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   NebulaBox MVP End-to-End Testing Flow         ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════╝${NC}"
echo ""

# Step 1: Check services
echo -e "${BLUE}Step 1: Checking services...${NC}"
API_HEALTH=$(curl -s "$API_URL/health" | jq -r '.status // "unknown"')
if [ "$API_HEALTH" = "healthy" ]; then
    echo -e "${GREEN}✅ API server is healthy${NC}"
else
    echo -e "${RED}❌ API server is not healthy${NC}"
    exit 1
fi

# Step 2: Build image from buildspec
echo -e "${BLUE}Step 2: Building image from buildspec.json...${NC}"
BUILD_REQUEST=$(jq -c '{spec: ., tag: "mern-mvp:latest"}' "$BUILDSPEC_PATH")
BUILD_RESPONSE=$(curl -s -X POST "$API_URL/buildspec/build" \
    -H "Content-Type: application/json" \
    -d "$BUILD_REQUEST")

BUILD_VALID=$(echo "$BUILD_RESPONSE" | jq -r '.valid // false')
if [ "$BUILD_VALID" = "true" ]; then
    echo -e "${GREEN}✅ Image built successfully: mern-mvp:latest${NC}"
    echo "$BUILD_RESPONSE" | jq '{tag, message, logs_count: (.logs | length)}'
else
    echo -e "${RED}❌ Build failed${NC}"
    echo "$BUILD_RESPONSE" | jq '.'
    exit 1
fi

# Step 3: Create container
echo ""
echo -e "${BLUE}Step 3: Creating container...${NC}"
CONTAINER_DATA='{
    "image": "mern-mvp:latest",
    "name": "mern-mvp-test",
    "ports": ["3000:3000", "5000:5000", "27017:27017"],
    "env": [
        "MONGO_HOST=localhost",
        "MONGO_PORT=27017",
        "MONGO_DB=todos",
        "PORT=5000"
    ],
    "detach": true
}'

CONTAINER_RESPONSE=$(curl -s -X POST "$API_URL/containers/run" \
    -H "Content-Type: application/json" \
    -d "$CONTAINER_DATA")

CONTAINER_ID=$(echo "$CONTAINER_RESPONSE" | jq -r '.id // .container.id // .[0].id // empty')
if [ -z "$CONTAINER_ID" ] || [ "$CONTAINER_ID" = "null" ]; then
    # Try to get from list
    CONTAINER_ID=$(curl -s "$API_URL/containers" | jq -r '.[] | select(.name == "mern-mvp-test") | .id' | head -1)
fi

if [ -n "$CONTAINER_ID" ] && [ "$CONTAINER_ID" != "null" ]; then
    echo -e "${GREEN}✅ Container created: $CONTAINER_ID${NC}"
    echo "$CONTAINER_RESPONSE" | jq '{id, name, image, status}' 2>/dev/null || echo "$CONTAINER_RESPONSE"
    echo "$CONTAINER_ID" > /tmp/mvp-container-id.txt
else
    echo -e "${YELLOW}⚠️  Container ID not in response, checking container list...${NC}"
    curl -s "$API_URL/containers" | jq '.[] | {id, name, image}' | head -10
    CONTAINER_ID=$(curl -s "$API_URL/containers" | jq -r '.[] | select(.name == "mern-mvp-test") | .id' | head -1)
    if [ -n "$CONTAINER_ID" ]; then
        echo -e "${GREEN}✅ Found container: $CONTAINER_ID${NC}"
        echo "$CONTAINER_ID" > /tmp/mvp-container-id.txt
    else
        echo -e "${RED}❌ Container not found${NC}"
        exit 1
    fi
fi

# Step 4: Start container
echo ""
echo -e "${BLUE}Step 4: Starting container...${NC}"
START_RESPONSE=$(curl -s -X POST "$API_URL/containers/$CONTAINER_ID/start")
echo "$START_RESPONSE" | jq '.' 2>/dev/null || echo "$START_RESPONSE"
echo -e "${GREEN}✅ Container start command sent${NC}"

# Wait for container to initialize
echo -e "${YELLOW}⏳ Waiting 5 seconds for services to initialize...${NC}"
sleep 5

# Step 5: Get container status
echo ""
echo -e "${BLUE}Step 5: Container Status...${NC}"
CONTAINER_STATUS=$(curl -s "$API_URL/containers/$CONTAINER_ID")
echo "$CONTAINER_STATUS" | jq '{id, name, image, status, ports, created}' 2>/dev/null || echo "$CONTAINER_STATUS"

# Step 6: Get container logs
echo ""
echo -e "${BLUE}Step 6: Container Logs (last 15 lines)...${NC}"
LOGS_RESPONSE=$(curl -s "$API_URL/containers/$CONTAINER_ID/logs?tail=15")
echo "$LOGS_RESPONSE" | jq -r '.logs[]? // .logs // .[]? // .' 2>/dev/null | tail -15 || echo "$LOGS_RESPONSE" | tail -15

# Step 7: Test application endpoints
echo ""
echo -e "${BLUE}Step 7: Testing Application Endpoints...${NC}"
echo -e "${YELLOW}Backend Health:${NC}"
curl -s http://localhost:5000/health 2>/dev/null | jq '.' || echo -e "${RED}Backend not responding yet${NC}"

echo -e "${YELLOW}Frontend:${NC}"
curl -s -I http://localhost:3000 2>/dev/null | head -3 || echo -e "${RED}Frontend not responding yet${NC}"

# Step 8: System stats
echo ""
echo -e "${BLUE}Step 8: System Statistics...${NC}"
curl -s "$API_URL/system/stats" | jq '{cpuUsage, memoryUsage, containersRunning, containersTotal}'

# Summary
echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}   ✅ End-to-End Flow Completed!                   ${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}📺 Check Dashboard:${NC}"
echo "   http://localhost:3001/containers"
echo ""
echo -e "${BLUE}🧪 Test Application:${NC}"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:5000/health"
echo ""
echo -e "${BLUE}📋 Container ID:${NC}"
echo "   $CONTAINER_ID"
echo ""
echo -e "${YELLOW}💡 View logs:${NC}"
echo "   curl $API_URL/containers/$CONTAINER_ID/logs | jq -r '.logs[]'"
echo ""

