#!/bin/bash

# Restart Dashboard Server
# This script stops the dashboard (if running) and starts it again

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🔄 Restarting NebulaBox Dashboard...${NC}"
echo ""

# Step 1: Stop dashboard if running
echo -e "${YELLOW}🛑 Stopping dashboard (if running)...${NC}"

DASHBOARD_PORT="${DASHBOARD_PORT:-3001}"
PID_FILE="/tmp/nebulabox-dashboard.pid"

# Stop using PID file
if [ -f "$PID_FILE" ]; then
    DASHBOARD_PID=$(cat "$PID_FILE" 2>/dev/null || echo "")
    if [ ! -z "$DASHBOARD_PID" ] && ps -p "$DASHBOARD_PID" > /dev/null 2>&1; then
        echo -e "${BLUE}   Stopping process (PID: $DASHBOARD_PID)...${NC}"
        kill "$DASHBOARD_PID" 2>/dev/null || true
        sleep 1
        # Force kill if still running
        if ps -p "$DASHBOARD_PID" > /dev/null 2>&1; then
            kill -9 "$DASHBOARD_PID" 2>/dev/null || true
        fi
        rm -f "$PID_FILE"
    else
        rm -f "$PID_FILE"
    fi
fi

# Stop by process name
pkill -f "vite.*dashboard" 2>/dev/null && sleep 1 || true
pkill -f "npm.*run.*dev" 2>/dev/null && sleep 1 || true

# Stop by port
if lsof -Pi :$DASHBOARD_PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    PORT_PID=$(lsof -ti :$DASHBOARD_PORT)
    if [ ! -z "$PORT_PID" ]; then
        echo -e "${BLUE}   Stopping process on port $DASHBOARD_PORT (PID: $PORT_PID)...${NC}"
        kill "$PORT_PID" 2>/dev/null || true
        sleep 1
        if lsof -Pi :$DASHBOARD_PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
            kill -9 "$PORT_PID" 2>/dev/null || true
        fi
    fi
fi

# Wait a moment for processes to fully stop
sleep 2

# Verify stopped
if [ -f "$PID_FILE" ] && ps -p "$(cat $PID_FILE 2>/dev/null)" > /dev/null 2>&1; then
    echo -e "${RED}⚠️  Dashboard process still running, forcing stop...${NC}"
    pkill -9 -f "vite.*dashboard" 2>/dev/null || true
    pkill -9 -f "npm.*run.*dev" 2>/dev/null || true
    rm -f "$PID_FILE"
    sleep 1
fi

if lsof -Pi :$DASHBOARD_PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo -e "${RED}⚠️  Port $DASHBOARD_PORT still in use, forcing stop...${NC}"
    lsof -ti :$DASHBOARD_PORT | xargs kill -9 2>/dev/null || true
    sleep 1
fi

echo -e "${GREEN}✅ Dashboard stopped${NC}"
echo ""

# Step 2: Start dashboard
echo -e "${BLUE}🚀 Starting dashboard...${NC}"
"$SCRIPT_DIR/start-dashboard.sh"

# Exit with the exit code from start script
exit $?

