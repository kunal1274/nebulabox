#!/bin/bash

# Check Dashboard Server Status
# Quick script to verify if dashboard is running

DASHBOARD_PORT="${DASHBOARD_PORT:-3001}"
PID_FILE="/tmp/nebulabox-dashboard.pid"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 Checking Dashboard Status...${NC}"
echo ""

# Check PID file
if [ -f "$PID_FILE" ]; then
    DASHBOARD_PID=$(cat "$PID_FILE" 2>/dev/null || echo "")
    if [ ! -z "$DASHBOARD_PID" ] && ps -p "$DASHBOARD_PID" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Dashboard is running${NC}"
        echo -e "   PID: $DASHBOARD_PID"
        echo -e "   PID File: $PID_FILE"
    else
        echo -e "${YELLOW}⚠️  PID file exists but process not found${NC}"
    fi
fi

# Check by process name
VITE_PID=$(pgrep -f "vite.*dashboard" | head -1 || echo "")
NPM_PID=$(pgrep -f "npm.*run.*dev" | head -1 || echo "")

if [ ! -z "$VITE_PID" ] || [ ! -z "$NPM_PID" ]; then
    if [ ! -z "$VITE_PID" ]; then
        echo -e "${GREEN}✅ Dashboard process found${NC}"
        echo -e "   Vite PID: $VITE_PID"
    fi
    if [ ! -z "$NPM_PID" ]; then
        echo -e "${GREEN}✅ npm process found${NC}"
        echo -e "   npm PID: $NPM_PID"
    fi
fi

# Check by port
if lsof -Pi :$DASHBOARD_PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    PORT_PID=$(lsof -ti :$DASHBOARD_PORT)
    PORT_INFO=$(lsof -i :$DASHBOARD_PORT | tail -n +2)
    echo -e "${GREEN}✅ Port $DASHBOARD_PORT is in use${NC}"
    echo -e "   PID: $PORT_PID"
    echo -e "   Details:"
    echo "$PORT_INFO" | sed 's/^/   /'
else
    echo -e "${BLUE}ℹ️  Port $DASHBOARD_PORT is free${NC}"
fi

# Summary
echo ""
if [ -f "$PID_FILE" ] && ps -p "$(cat $PID_FILE 2>/dev/null)" > /dev/null 2>&1 || \
   [ ! -z "$VITE_PID" ] || \
   [ ! -z "$NPM_PID" ] || \
   lsof -Pi :$DASHBOARD_PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${GREEN}═══════════════════════════════════════${NC}"
    echo -e "${GREEN}   Dashboard Status: RUNNING${NC}"
    echo -e "${GREEN}═══════════════════════════════════════${NC}"
    echo -e "${BLUE}   URL: http://localhost:$DASHBOARD_PORT${NC}"
    exit 0
else
    echo -e "${BLUE}═══════════════════════════════════════${NC}"
    echo -e "${BLUE}   Dashboard Status: NOT RUNNING${NC}"
    echo -e "${BLUE}═══════════════════════════════════════${NC}"
    exit 1
fi

