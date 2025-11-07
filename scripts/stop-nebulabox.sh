#!/bin/bash

# Stop NebulaBox API Server and Dashboard

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🛑 Stopping NebulaBox services...${NC}"

# Stop API server
if [ -f /tmp/nebulabox-api.pid ]; then
    API_PID=$(cat /tmp/nebulabox-api.pid)
    if ps -p $API_PID > /dev/null 2>&1; then
        echo -e "${BLUE}Stopping API server (PID: $API_PID)...${NC}"
        kill $API_PID
        rm /tmp/nebulabox-api.pid
        echo -e "${GREEN}✅ API server stopped${NC}"
    else
        rm /tmp/nebulabox-api.pid
    fi
fi

# Also try to kill by port/process name
pkill -f nebulabox-api && echo -e "${GREEN}✅ API server processes stopped${NC}" || true

# Stop Dashboard
if [ -f /tmp/nebulabox-dashboard.pid ]; then
    DASHBOARD_PID=$(cat /tmp/nebulabox-dashboard.pid)
    if ps -p $DASHBOARD_PID > /dev/null 2>&1; then
        echo -e "${BLUE}Stopping Dashboard (PID: $DASHBOARD_PID)...${NC}"
        kill $DASHBOARD_PID
        rm /tmp/nebulabox-dashboard.pid
        echo -e "${GREEN}✅ Dashboard stopped${NC}"
    else
        rm /tmp/nebulabox-dashboard.pid
    fi
fi

# Also try to kill vite processes
pkill -f "vite.*dashboard" && echo -e "${GREEN}✅ Dashboard processes stopped${NC}" || true

echo -e "${GREEN}✅ All NebulaBox services stopped${NC}"

