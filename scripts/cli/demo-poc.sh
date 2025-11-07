#!/bin/bash

# NebulaBox POC Demo Script
# This script demonstrates NebulaBox's unique features for investors/demo

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLI_BINARY="$PROJECT_ROOT/bin/nebulabox"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                                                              ║"
echo "║           NebulaBox - POC Demo for Investors                ║"
echo "║                                                              ║"
echo "║     Unified Container Platform - Different from Docker       ║"
echo "║                                                              ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""

echo -e "${BLUE}What is NebulaBox?${NC}"
echo ""
echo "NebulaBox is a unified container platform that simplifies the entire"
echo "development-to-deployment lifecycle. Unlike Docker or Kubernetes,"
echo "NebulaBox focuses on:"
echo ""
echo -e "  ${GREEN}✓${NC} Unified development (all services in one container)"
echo -e "  ${GREEN}✓${NC} Built-in collaboration (no VPN/ngrok needed)"
echo -e "  ${GREEN}✓${NC} Single deployment platform (no fragmentation)"
echo -e "  ${GREEN}✓${NC} Flexible architecture testing (monolithic to microservices)"
echo ""

echo -e "${YELLOW}Press Enter to start the interactive demo...${NC}"
read -r

# Run interactive demo
"$SCRIPT_DIR/workflow-00-interactive-demo.sh"

