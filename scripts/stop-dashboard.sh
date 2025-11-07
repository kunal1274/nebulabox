#!/bin/bash

# Stop Dashboard Server and Verify Status
# This script stops the dashboard and verifies if it's running

set -e

DASHBOARD_PORT="${DASHBOARD_PORT:-3001}"
PID_FILE="/tmp/nebulabox-dashboard.pid"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Function to check if a port is in use
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

# Function to check if dashboard is running
check_dashboard_running() {
    local is_running=false
    local status_msg=""
    local found_pid=""
    
    # Check by port FIRST - this is the most reliable check
    if check_port $DASHBOARD_PORT; then
        PORT_PID=$(lsof -ti :$DASHBOARD_PORT 2>/dev/null | head -1 || echo "")
        if [ ! -z "$PORT_PID" ]; then
            # Verify it's a node process listening on the port
            PORT_CMD=$(ps -p "$PORT_PID" -o comm= 2>/dev/null || echo "")
            # Check if it's a node process (vite runs as node)
            if [ "$PORT_CMD" = "node" ]; then
                PORT_ARGS=$(ps -p "$PORT_PID" -o args= 2>/dev/null || echo "")
                # Check if it's vite (not a build process)
                if echo "$PORT_ARGS" | grep -q "vite" && ! echo "$PORT_ARGS" | grep -q "vite build"; then
                    is_running=true
                    found_pid="$PORT_PID"
                    status_msg="Dashboard is running on port $DASHBOARD_PORT (PID: $PORT_PID)"
                fi
            fi
        fi
    fi
    
    # Check PID file (if port check didn't find it)
    if [ "$is_running" = "false" ] && [ -f "$PID_FILE" ]; then
        DASHBOARD_PID=$(cat "$PID_FILE" 2>/dev/null || echo "")
        if [ ! -z "$DASHBOARD_PID" ] && ps -p "$DASHBOARD_PID" > /dev/null 2>&1; then
            is_running=true
            found_pid="$DASHBOARD_PID"
            status_msg="Dashboard is running (PID from file: $DASHBOARD_PID)"
        fi
    fi
    
    # If port check didn't find it, check by process name (but exclude build processes)
    if [ "$is_running" = "false" ]; then
        # Check all vite processes and see if any are from dashboard (but not build)
        for vite_pid in $(pgrep -f "vite" 2>/dev/null || true); do
            if [ ! -z "$vite_pid" ]; then
                VITE_CMD=$(ps -p "$vite_pid" -o args= 2>/dev/null || echo "")
                # Must be from dashboard directory and NOT a build process
                if echo "$VITE_CMD" | grep -q "dashboard" && ! echo "$VITE_CMD" | grep -q "vite build"; then
                    is_running=true
                    found_pid="$vite_pid"
                    status_msg="Dashboard is running (Vite PID: $vite_pid)"
                    break
                fi
            fi
        done
        
        # Also check for npm run dev
        if [ "$is_running" = "false" ]; then
            for npm_pid in $(pgrep -f "npm.*run.*dev" 2>/dev/null || true); do
                if [ ! -z "$npm_pid" ]; then
                    NPM_CWD=$(ps -p "$npm_pid" -o cwd= 2>/dev/null || echo "")
                    if echo "$NPM_CWD" | grep -q "dashboard"; then
                        is_running=true
                        found_pid="$npm_pid"
                        status_msg="Dashboard is running (npm PID: $npm_pid)"
                        break
                    fi
                fi
            done
        fi
    fi
    
    echo "$is_running|$status_msg|$found_pid"
}

# Main function
main() {
    echo -e "${BLUE}🔍 Checking Dashboard Status...${NC}"
    echo ""
    
    # Check if running
    result=$(check_dashboard_running)
    is_running=$(echo "$result" | cut -d'|' -f1)
    status_msg=$(echo "$result" | cut -d'|' -f2)
    found_pid=$(echo "$result" | cut -d'|' -f3)
    
    if [ "$is_running" = "true" ]; then
        echo -e "${YELLOW}📊 Status: $status_msg${NC}"
        echo ""
        echo -e "${BLUE}🛑 Stopping Dashboard...${NC}"
        
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
        
        # Stop by process name - more aggressive matching
        # Kill vite processes in dashboard directory
        pgrep -f "vite" | while read pid; do
            if ps -p "$pid" -o args= 2>/dev/null | grep -q "dashboard"; then
                kill "$pid" 2>/dev/null || true
            fi
        done
        
        # Kill npm run dev processes in dashboard directory  
        pgrep -f "npm.*run.*dev" | while read pid; do
            if ps -p "$pid" -o cwd= 2>/dev/null | grep -q "dashboard"; then
                kill "$pid" 2>/dev/null || true
            fi
        done
        
        sleep 2
        
        # Force kill any remaining vite processes in dashboard
        pgrep -f "vite" | while read pid; do
            if ps -p "$pid" -o args= 2>/dev/null | grep -q "dashboard"; then
                kill -9 "$pid" 2>/dev/null || true
            fi
        done
        
        # Force kill any remaining npm processes
        pgrep -f "npm.*run.*dev" | while read pid; do
            if ps -p "$pid" -o cwd= 2>/dev/null | grep -q "dashboard"; then
                kill -9 "$pid" 2>/dev/null || true
            fi
        done
        
        # Stop by port - kill all processes using the port
        if check_port $DASHBOARD_PORT; then
            PORT_PIDS=$(lsof -ti :$DASHBOARD_PORT 2>/dev/null || echo "")
            if [ ! -z "$PORT_PIDS" ]; then
                echo -e "${BLUE}   Stopping processes on port $DASHBOARD_PORT...${NC}"
                echo "$PORT_PIDS" | while read port_pid; do
                    if [ ! -z "$port_pid" ]; then
                        PORT_CMD=$(ps -p "$port_pid" -o comm= 2>/dev/null || echo "")
                        if [ "$PORT_CMD" = "node" ] || [ "$PORT_CMD" = "vite" ]; then
                            echo -e "${BLUE}     Killing PID: $port_pid ($PORT_CMD)...${NC}"
                            kill "$port_pid" 2>/dev/null || true
                        fi
                    fi
                done
                sleep 2
                # Force kill if still running
                if check_port $DASHBOARD_PORT; then
                    lsof -ti :$DASHBOARD_PORT 2>/dev/null | while read port_pid; do
                        PORT_CMD=$(ps -p "$port_pid" -o comm= 2>/dev/null || echo "")
                        if [ "$PORT_CMD" = "node" ] || [ "$PORT_CMD" = "vite" ]; then
                            kill -9 "$port_pid" 2>/dev/null || true
                        fi
                    done
                fi
            fi
        fi
        
        # Verify stopped
        sleep 2
        result=$(check_dashboard_running)
        is_running_after=$(echo "$result" | cut -d'|' -f1)
        
        if [ "$is_running_after" = "false" ]; then
            echo -e "${GREEN}✅ Dashboard stopped successfully${NC}"
        else
            echo -e "${RED}⚠️  Dashboard may still be running. Trying force stop...${NC}"
            # More aggressive force kill
            pgrep -f "vite" | while read pid; do
                if ps -p "$pid" -o args= 2>/dev/null | grep -q "dashboard"; then
                    kill -9 "$pid" 2>/dev/null || true
                fi
            done
            pgrep -f "npm.*run.*dev" | while read pid; do
                if ps -p "$pid" -o cwd= 2>/dev/null | grep -q "dashboard"; then
                    kill -9 "$pid" 2>/dev/null || true
                fi
            done
            # Kill by port as last resort
            if check_port $DASHBOARD_PORT; then
                lsof -ti :$DASHBOARD_PORT 2>/dev/null | xargs kill -9 2>/dev/null || true
            fi
            sleep 2
        fi
    else
        echo -e "${GREEN}✅ Dashboard is not running${NC}"
    fi
    
    echo ""
    echo -e "${BLUE}📊 Final Status Check:${NC}"
    
    # Final verification
    result=$(check_dashboard_running)
    is_running_final=$(echo "$result" | cut -d'|' -f1)
    status_msg_final=$(echo "$result" | cut -d'|' -f2)
    
    if [ "$is_running_final" = "true" ]; then
        echo -e "${RED}❌ Dashboard is still running: $status_msg_final${NC}"
        echo -e "${YELLOW}   Try manually killing processes:${NC}"
        echo -e "${YELLOW}   lsof -ti :$DASHBOARD_PORT | xargs kill -9${NC}"
        echo -e "${YELLOW}   Or: pkill -9 -f vite${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ Dashboard is not running${NC}"
        if check_port $DASHBOARD_PORT; then
            echo -e "${YELLOW}⚠️  Port $DASHBOARD_PORT is still in use (may be another process)${NC}"
        else
            echo -e "${BLUE}   Port $DASHBOARD_PORT is free${NC}"
        fi
        exit 0
    fi
}

# Handle command line arguments
case "${1:-}" in
    "check"|"status")
        result=$(check_dashboard_running)
        is_running=$(echo "$result" | cut -d'|' -f1)
        status_msg=$(echo "$result" | cut -d'|' -f2)
        
        if [ "$is_running" = "true" ]; then
            echo -e "${YELLOW}📊 $status_msg${NC}"
            exit 0
        else
            echo -e "${GREEN}✅ Dashboard is not running${NC}"
            exit 1
        fi
        ;;
    "stop")
        main
        ;;
    *)
        main
        ;;
esac


