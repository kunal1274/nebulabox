#!/bin/bash

# Start NebulaBox API Server and Dashboard
# This script starts both services for MVP testing

set -e

API_PORT="${NEBULA_API_PORT:-8081}"
DASHBOARD_PORT="${DASHBOARD_PORT:-3001}"
API_URL="http://localhost:${API_PORT}"
DASHBOARD_URL="http://localhost:${DASHBOARD_PORT}"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🚀 Starting NebulaBox for MVP Testing${NC}"
echo ""

# Check if already running
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

# Check API port
if check_port $API_PORT; then
    echo -e "${YELLOW}⚠️  API server already running on port $API_PORT${NC}"
else
    echo -e "${BLUE}Starting API server on port $API_PORT...${NC}"
    
    # Build API if needed
    if [ ! -f "./nebulabox-api" ]; then
        echo -e "${BLUE}Building API server...${NC}"
        make build-api
    fi
    
    # Start API in background
    echo -e "${BLUE}Starting API server...${NC}"
    ./nebulabox-api > /tmp/nebulabox-api.log 2>&1 &
    API_PID=$!
    echo $API_PID > /tmp/nebulabox-api.pid
    
    # Wait for API to be ready
    echo -e "${BLUE}Waiting for API server to start...${NC}"
    for i in {1..30}; do
        if curl -s "$API_URL/api/health" > /dev/null 2>&1 || curl -s "$API_URL/api/containers" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ API server started!${NC}"
            break
        fi
        sleep 1
        if [ $i -eq 30 ]; then
            echo -e "${RED}❌ API server failed to start${NC}"
            exit 1
        fi
    done
fi

# Check Dashboard port (may use 3000 if 3001 is taken)
DASHBOARD_ACTUAL_PORT=$DASHBOARD_PORT
if check_port $DASHBOARD_PORT; then
    echo -e "${YELLOW}⚠️  Port $DASHBOARD_PORT in use, dashboard will use next available port${NC}"
    DASHBOARD_ACTUAL_PORT=3000
fi

if check_port $DASHBOARD_ACTUAL_PORT; then
    echo -e "${YELLOW}⚠️  Dashboard already running on port $DASHBOARD_ACTUAL_PORT${NC}"
else
    echo -e "${BLUE}Starting Dashboard (will use available port)...${NC}"
    cd web/dashboard
    
    # Install dependencies if needed
    if [ ! -d "node_modules" ]; then
        echo -e "${BLUE}Installing dashboard dependencies...${NC}"
        npm install
    fi
    
    # Start dashboard in background
    echo -e "${BLUE}Starting dashboard...${NC}"
    npm run dev > /tmp/nebulabox-dashboard.log 2>&1 &
    DASHBOARD_PID=$!
    echo $DASHBOARD_PID > /tmp/nebulabox-dashboard.pid
    cd ../..
    
    # Wait for dashboard to be ready
    echo -e "${BLUE}Waiting for dashboard to start...${NC}"
    sleep 8
    
    # Detect actual port from logs
    if grep -q "Local:.*http://localhost:" /tmp/nebulabox-dashboard.log 2>/dev/null; then
        DASHBOARD_ACTUAL_PORT=$(grep -o "Local:.*http://localhost:[0-9]*" /tmp/nebulabox-dashboard.log 2>/dev/null | grep -o "[0-9]*" | head -1 || echo "$DASHBOARD_PORT")
    fi
    
    # Check if dashboard is running
    if check_port $DASHBOARD_ACTUAL_PORT || check_port $DASHBOARD_PORT; then
        echo -e "${GREEN}✅ Dashboard started on port $DASHBOARD_ACTUAL_PORT!${NC}"
        DASHBOARD_URL="http://localhost:${DASHBOARD_ACTUAL_PORT}"
    else
        echo -e "${YELLOW}⚠️  Dashboard may still be starting...${NC}"
        echo -e "${YELLOW}Check logs: tail -f /tmp/nebulabox-dashboard.log${NC}"
    fi
fi

echo ""
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}   NebulaBox is Ready!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}📡 API Server:${NC}"
echo -e "   URL: $API_URL"
echo -e "   Health: $API_URL/api/health"
echo -e "   Containers: $API_URL/api/containers"
echo ""
echo -e "${BLUE}🎨 Dashboard:${NC}"
if [ -n "$DASHBOARD_ACTUAL_PORT" ]; then
    echo -e "   URL: http://localhost:${DASHBOARD_ACTUAL_PORT}"
    echo -e "   Open in browser: http://localhost:${DASHBOARD_ACTUAL_PORT}"
else
    echo -e "   URL: $DASHBOARD_URL (or check logs for actual port)"
    echo -e "   Open in browser: $DASHBOARD_URL"
fi
echo ""
echo -e "${BLUE}📋 Logs:${NC}"
echo -e "   API: tail -f /tmp/nebulabox-api.log"
echo -e "   Dashboard: tail -f /tmp/nebulabox-dashboard.log"
echo ""
echo -e "${YELLOW}To stop services:${NC}"
echo -e "   ./scripts/stop-nebulabox.sh"
echo ""
echo -e "${GREEN}Ready for MVP testing! 🚀${NC}"

