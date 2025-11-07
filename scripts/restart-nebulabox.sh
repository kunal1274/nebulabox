#!/bin/bash

# NebulaBox Service Restart Script
# Stops and restarts API server and dashboard

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔═══════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     NebulaBox Service Restart Script            ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════╝${NC}"
echo ""

# Step 1: Stop existing services
echo -e "${YELLOW}🛑 Stopping existing services...${NC}"
pkill -f nebulabox-api || true
pkill -f "vite.*dashboard" || true
sleep 2

# Step 2: Build API
echo -e "${BLUE}🔨 Building API server...${NC}"
cd "$PROJECT_ROOT"
make build-api

# Step 3: Start API server
echo -e "${BLUE}🚀 Starting API server...${NC}"
nohup ./nebulabox-api > api.log 2>&1 &
API_PID=$!
echo -e "${GREEN}✅ API server started (PID: $API_PID)${NC}"

# Step 4: Wait for API to be ready
echo -e "${YELLOW}⏳ Waiting for API to be ready...${NC}"
MAX_RETRIES=10
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:8081/api/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ API server is healthy!${NC}"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -e "${YELLOW}   Attempt $RETRY_COUNT/$MAX_RETRIES...${NC}"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}❌ API server failed to start within timeout${NC}"
    exit 1
fi

# Step 5: Show status
echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}   ✅ Services Restarted Successfully!            ${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}📊 API Server:${NC}"
echo "   URL:    http://localhost:8081"
echo "   Health: http://localhost:8081/api/health"
echo "   PID:    $API_PID"
echo "   Logs:   tail -f $PROJECT_ROOT/api.log"
echo ""
echo -e "${BLUE}🎨 Dashboard:${NC}"
echo "   URL:    http://localhost:3001"
echo "   Start:  cd web/dashboard && npm run dev"
echo ""
echo -e "${YELLOW}💡 To stop API server:${NC}"
echo "   kill $API_PID"
echo ""
echo -e "${YELLOW}💡 To view API logs:${NC}"
echo "   tail -f $PROJECT_ROOT/api.log"
echo ""

