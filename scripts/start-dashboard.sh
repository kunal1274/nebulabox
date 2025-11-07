#!/bin/bash

# Start Dashboard Server Only
# This script starts only the dashboard (not the API)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DASHBOARD_PORT="${DASHBOARD_PORT:-3001}"
DASHBOARD_DIR="$PROJECT_ROOT/web/dashboard"
PID_FILE="/tmp/nebulabox-dashboard.pid"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🎨 Starting NebulaBox Dashboard...${NC}"
echo ""

# Function to check if a port is in use
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

# Function to check if dashboard is already running
check_dashboard_running() {
    # Check PID file
    if [ -f "$PID_FILE" ]; then
        DASHBOARD_PID=$(cat "$PID_FILE" 2>/dev/null || echo "")
        if [ ! -z "$DASHBOARD_PID" ] && ps -p "$DASHBOARD_PID" > /dev/null 2>&1; then
            return 0  # Running
        fi
    fi
    
    # Check by process name
    if pgrep -f "vite.*dashboard" > /dev/null 2>&1 || pgrep -f "npm.*run.*dev" > /dev/null 2>&1; then
        return 0  # Running
    fi
    
    # Check by port
    if check_port $DASHBOARD_PORT; then
        return 0  # Running
    fi
    
    return 1  # Not running
}

# Check if already running
if check_dashboard_running; then
    if [ -f "$PID_FILE" ]; then
        EXISTING_PID=$(cat "$PID_FILE" 2>/dev/null || echo "")
        if [ ! -z "$EXISTING_PID" ] && ps -p "$EXISTING_PID" > /dev/null 2>&1; then
            echo -e "${YELLOW}⚠️  Dashboard already running (PID: $EXISTING_PID)${NC}"
        else
            echo -e "${YELLOW}⚠️  Dashboard appears to be running (checking processes...)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  Dashboard appears to be running${NC}"
    fi
    
    # Show port status
    if check_port $DASHBOARD_PORT; then
        PORT_PID=$(lsof -ti :$DASHBOARD_PORT)
        echo -e "${BLUE}   Port $DASHBOARD_PORT is in use (PID: $PORT_PID)${NC}"
        echo -e "${YELLOW}   Dashboard URL: http://localhost:$DASHBOARD_PORT${NC}"
    fi
    
    echo -e "${YELLOW}   To restart, use: ./scripts/restart-dashboard.sh${NC}"
    exit 0
fi

# Check Dashboard port (may use 3000 if 3001 is taken)
DASHBOARD_ACTUAL_PORT=$DASHBOARD_PORT
if check_port $DASHBOARD_PORT; then
    echo -e "${YELLOW}⚠️  Port $DASHBOARD_PORT in use, trying port 3000...${NC}"
    DASHBOARD_ACTUAL_PORT=3000
    
    if check_port $DASHBOARD_ACTUAL_PORT; then
        echo -e "${RED}❌ Both ports $DASHBOARD_PORT and $DASHBOARD_ACTUAL_PORT are in use${NC}"
        echo -e "${YELLOW}   Please free up a port or stop existing dashboard${NC}"
        exit 1
    fi
fi

# Navigate to dashboard directory
if [ ! -d "$DASHBOARD_DIR" ]; then
    echo -e "${RED}❌ Dashboard directory not found: $DASHBOARD_DIR${NC}"
    exit 1
fi

cd "$DASHBOARD_DIR"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo -e "${BLUE}📦 Installing dashboard dependencies...${NC}"
    npm install
    echo ""
fi

# Start dashboard in background
echo -e "${BLUE}🚀 Starting dashboard on port $DASHBOARD_ACTUAL_PORT...${NC}"
npm run dev > /tmp/nebulabox-dashboard.log 2>&1 &
DASHBOARD_PID=$!
echo $DASHBOARD_PID > "$PID_FILE"

cd "$PROJECT_ROOT"

# Wait for dashboard to be ready
echo -e "${BLUE}⏳ Waiting for dashboard to start...${NC}"
sleep 5

# Detect actual port from logs
if [ -f /tmp/nebulabox-dashboard.log ]; then
    if grep -q "Local:.*http://localhost:" /tmp/nebulabox-dashboard.log 2>/dev/null; then
        DETECTED_PORT=$(grep -o "Local:.*http://localhost:[0-9]*" /tmp/nebulabox-dashboard.log 2>/dev/null | grep -o "[0-9]*" | head -1 || echo "")
        if [ ! -z "$DETECTED_PORT" ]; then
            DASHBOARD_ACTUAL_PORT=$DETECTED_PORT
        fi
    fi
fi

# Wait a bit more and verify
sleep 3

# Check if dashboard is running
if check_port $DASHBOARD_ACTUAL_PORT || check_port $DASHBOARD_PORT; then
    FINAL_PORT=$DASHBOARD_ACTUAL_PORT
    if check_port $DASHBOARD_PORT; then
        FINAL_PORT=$DASHBOARD_PORT
    fi
    
    echo -e "${GREEN}✅ Dashboard started successfully!${NC}"
    echo ""
    echo -e "${GREEN}════════════════════════════════════════${NC}"
    echo -e "${GREEN}   Dashboard is Ready!${NC}"
    echo -e "${GREEN}════════════════════════════════════════${NC}"
    echo ""
    echo -e "${BLUE}🎨 Dashboard:${NC}"
    echo -e "   URL: http://localhost:$FINAL_PORT"
    echo -e "   PID: $DASHBOARD_PID"
    echo -e "   Port: $FINAL_PORT"
    echo ""
    echo -e "${BLUE}📋 Logs:${NC}"
    echo -e "   tail -f /tmp/nebulabox-dashboard.log"
    echo ""
    echo -e "${YELLOW}To stop dashboard:${NC}"
    echo -e "   ./scripts/stop-dashboard.sh"
    echo ""
else
    echo -e "${YELLOW}⚠️  Dashboard may still be starting...${NC}"
    echo -e "${YELLOW}   Check logs: tail -f /tmp/nebulabox-dashboard.log${NC}"
    echo -e "${YELLOW}   PID: $DASHBOARD_PID${NC}"
    exit 1
fi

