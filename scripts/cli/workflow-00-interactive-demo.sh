#!/bin/bash

# Workflow 00: Interactive End-to-End Demo
# This is a POC CLI demo script that showcases NebulaBox's unique features
# Different from Docker/Kubernetes - demonstrates unified development-to-deployment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLI_BINARY="$PROJECT_ROOT/bin/nebulabox"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Demo state
CURRENT_STEP=0
TOTAL_STEPS=8
DEMO_IMAGE=""
DEMO_CONTAINER=""
DEMO_GROUP=""

# Check if CLI binary exists
if [ ! -f "$CLI_BINARY" ]; then
    echo -e "${RED}❌ CLI binary not found at $CLI_BINARY${NC}"
    echo "   Run 'make build-cli-test' first"
    exit 1
fi

# Print header
print_header() {
    clear
    echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  ${PURPLE}NebulaBox POC Demo - Interactive End-to-End Workflow${NC}  ${CYAN}║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}What makes NebulaBox different:${NC}"
    echo -e "  • ${GREEN}Unified Development${NC} - Single container for all services (no Docker Compose needed)"
    echo -e "  • ${GREEN}Built-in Collaboration${NC} - Real-time shared workspaces (no VPN/ngrok needed)"
    echo -e "  • ${GREEN}Single Deployment${NC} - Deploy everything together (no split across platforms)"
    echo -e "  • ${GREEN}Flexible Grouping${NC} - Test different architectures easily (monolithic, microservices, etc.)"
    echo ""
    echo -e "${YELLOW}Progress: Step $CURRENT_STEP/$TOTAL_STEPS${NC}"
    echo ""
}

# Print menu
print_menu() {
    echo -e "${CYAN}Available Actions:${NC}"
    echo -e "  ${GREEN}[Enter]${NC} - Continue to next step"
    echo -e "  ${YELLOW}[s]${NC} - Skip this step"
    echo -e "  ${YELLOW}[b]${NC} - Go back to previous step"
    echo -e "  ${YELLOW}[r]${NC} - Restart from beginning"
    echo -e "  ${RED}[q]${NC} - Quit demo"
    echo ""
    echo -ne "${CYAN}Your choice: ${NC}"
}

# Get user input
get_choice() {
    read -r choice
    case "$choice" in
        ""|"c"|"C")
            return 0  # Continue
            ;;
        "s"|"S")
            return 1  # Skip
            ;;
        "b"|"B")
            return 2  # Back
            ;;
        "r"|"R")
            return 3  # Restart
            ;;
        "q"|"Q")
            return 4  # Quit
            ;;
        *)
            echo -e "${RED}Invalid choice. Please try again.${NC}"
            sleep 1
            return 5
            ;;
    esac
}

# Step 1: Introduction
step_1_intro() {
    print_header
    echo -e "${PURPLE}Step 1: Introduction to NebulaBox${NC}"
    echo ""
    echo -e "${BLUE}NebulaBox is a unified container platform that:${NC}"
    echo ""
    echo -e "1. ${GREEN}Simplifies Development${NC}"
    echo "   • Run entire MERN stack in ONE container"
    echo "   • No need for Docker Compose or multiple docker run commands"
    echo "   • All services (frontend, backend, database) work together"
    echo ""
    echo -e "2. ${GREEN}Enables Real-time Collaboration${NC}"
    echo "   • Share workspaces with team members instantly"
    echo "   • No VPN, ngrok, or Tailscale setup required"
    echo "   • Built-in tunneling and file sync"
    echo ""
    echo -e "3. ${GREEN}Unified Deployment${NC}"
    echo "   • Deploy everything to ONE platform (NebulaBox Cloud)"
    echo "   • No splitting frontend to Vercel, backend to Render, DB to MongoDB Atlas"
    echo "   • Single URL, single configuration"
    echo ""
    echo -e "4. ${GREEN}Flexible Architecture Testing${NC}"
    echo "   • Test monolithic, microservices, or custom architectures"
    echo "   • Switch between strategies easily"
    echo "   • Perfect for POC and experimentation"
    echo ""
    print_menu
    get_choice
}

# Step 2: Build Image
step_2_build() {
    print_header
    echo -e "${PURPLE}Step 2: Build Image from BuildSpec${NC}"
    echo ""
    echo -e "${BLUE}Unlike Docker, NebulaBox uses BuildSpec (JSON) instead of Dockerfile:${NC}"
    echo ""
    echo -e "${YELLOW}Key Differences:${NC}"
    echo "  • Docker: Uses Dockerfile (text-based, line-by-line)"
    echo "  • NebulaBox: Uses BuildSpec (structured JSON, easier to generate/modify)"
    echo ""
    
    # Create a demo BuildSpec
    TEST_DIR=$(mktemp -d)
    cat > "$TEST_DIR/demo-buildspec.json" <<'EOF'
{
  "version": "1.0",
  "name": "demo-app",
  "tag": "demo-app:latest",
  "base": {
    "image": "alpine",
    "tag": "latest"
  },
  "workdir": "/app",
  "env": {
    "APP_NAME": "NebulaBox Demo"
  },
  "steps": [
    {
      "type": "run",
      "command": "apk add --no-cache nodejs npm",
      "comment": "Install Node.js"
    },
    {
      "type": "run",
      "command": "echo 'Hello from NebulaBox!' > /app/hello.txt",
      "comment": "Create demo file"
    }
  ]
}
EOF

    echo -e "${CYAN}BuildSpec created:${NC}"
    cat "$TEST_DIR/demo-buildspec.json" | head -15
    echo ""
    echo -e "${GREEN}Building image...${NC}"
    echo ""
    
    if $CLI_BINARY build -f "$TEST_DIR/demo-buildspec.json" -t demo-app:latest 2>&1; then
        DEMO_IMAGE="demo-app:latest"
        echo ""
        echo -e "${GREEN}✅ Image built successfully!${NC}"
        echo -e "${BLUE}Image: $DEMO_IMAGE${NC}"
    else
        echo ""
        echo -e "${YELLOW}⚠️  Build command may not be fully implemented yet${NC}"
        echo -e "${BLUE}   (This is expected in POC phase)${NC}"
        DEMO_IMAGE="demo-app:latest"  # Set anyway for demo flow
    fi
    
    rm -rf "$TEST_DIR"
    echo ""
    print_menu
    get_choice
}

# Step 3: Run Container
step_3_run() {
    print_header
    echo -e "${PURPLE}Step 3: Run Container${NC}"
    echo ""
    echo -e "${BLUE}NebulaBox runs containers with built-in service grouping:${NC}"
    echo ""
    echo -e "${YELLOW}Key Differences:${NC}"
    echo "  • Docker: One container = one service (need multiple containers for MERN)"
    echo "  • NebulaBox: One container can run multiple services together"
    echo "  • NebulaBox: Built-in grouping for flexible architectures"
    echo ""
    
    DEMO_CONTAINER="demo-container-$(date +%s)"
    echo -e "${CYAN}Running container: $DEMO_CONTAINER${NC}"
    echo ""
    
    if $CLI_BINARY run "${DEMO_IMAGE:-alpine:latest}" --name "$DEMO_CONTAINER" -d 2>&1; then
        echo ""
        echo -e "${GREEN}✅ Container started!${NC}"
        echo ""
        echo -e "${CYAN}Listing containers:${NC}"
        $CLI_BINARY ps 2>&1 || echo "   (List command may show mock data)"
    else
        echo ""
        echo -e "${YELLOW}⚠️  Run command may not be fully implemented yet${NC}"
        echo -e "${BLUE}   (This is expected in POC phase)${NC}"
    fi
    
    echo ""
    print_menu
    get_choice
}

# Step 4: Container Grouping
step_4_group() {
    print_header
    echo -e "${PURPLE}Step 4: Container Grouping (Unique Feature)${NC}"
    echo ""
    echo -e "${BLUE}NebulaBox allows flexible container grouping:${NC}"
    echo ""
    echo -e "${GREEN}Grouping Strategies:${NC}"
    echo "  1. Monolithic - All services in one container"
    echo "  2. Frontend-Backend - Frontend separate, backend+DB together"
    echo "  3. Three-Tier - Frontend, Backend, Database separate"
    echo "  4. Microservices - Each service in own container"
    echo "  5. Custom - Your own architecture"
    echo ""
    echo -e "${YELLOW}Why this is different:${NC}"
    echo "  • Docker: Manual container management, need Docker Compose for groups"
    echo "  • Kubernetes: Complex YAML, overkill for small apps"
    echo "  • NebulaBox: Simple JSON config, test different architectures easily"
    echo ""
    
    DEMO_GROUP="demo-group-$(date +%s)"
    echo -e "${CYAN}Example Group Configuration:${NC}"
    cat <<EOF
{
  "name": "$DEMO_GROUP",
  "strategy": "frontend-backend",
  "containers": [
    {
      "name": "frontend",
      "image": "nginx:alpine",
      "ports": ["3000:80"]
    },
    {
      "name": "backend-db",
      "image": "node:alpine",
      "ports": ["5000:5000", "27017:27017"]
    }
  ]
}
EOF
    echo ""
    echo -e "${YELLOW}⚠️  Group creation via CLI will be available in next phase${NC}"
    echo ""
    print_menu
    get_choice
}

# Step 5: Shared Workspace
step_5_workspace() {
    print_header
    echo -e "${PURPLE}Step 5: Shared Workspace (Unique Feature)${NC}"
    echo ""
    echo -e "${BLUE}NebulaBox enables real-time collaboration:${NC}"
    echo ""
    echo -e "${GREEN}Features:${NC}"
    echo "  • Real-time file synchronization (CRDT-based)"
    echo "  • Shared terminal access"
    echo "  • Live preview for all team members"
    echo "  • Built-in tunneling (no VPN/ngrok needed)"
    echo ""
    echo -e "${YELLOW}Why this is different:${NC}"
    echo "  • Docker: No built-in collaboration, need external tools"
    echo "  • Kubernetes: No collaboration features"
    echo "  • NebulaBox: Collaboration built-in from the start"
    echo ""
    echo -e "${CYAN}Example Commands (coming soon):${NC}"
    echo "  nebulabox workspace create my-app --share"
    echo "  nebulabox workspace invite my-app team@example.com"
    echo "  nebulabox workspace join my-app"
    echo ""
    echo -e "${YELLOW}⚠️  Workspace commands will be available in next phase${NC}"
    echo ""
    print_menu
    get_choice
}

# Step 6: Remote Deployment
step_6_remote() {
    print_header
    echo -e "${PURPLE}Step 6: Remote Deployment (No External Tools)${NC}"
    echo ""
    echo -e "${BLUE}NebulaBox has built-in remote deployment:${NC}"
    echo ""
    echo -e "${GREEN}Features:${NC}"
    echo "  • Built-in tunneling (no ngrok/Tailscale needed)"
    echo "  • WebSocket-based real-time connection"
    echo "  • SSH fallback for direct access"
    echo "  • Zero-config remote access"
    echo ""
    echo -e "${YELLOW}Why this is different:${NC}"
    echo "  • Docker: Need to setup port forwarding, VPN, or ngrok"
    echo "  • Kubernetes: Complex networking setup"
    echo "  • NebulaBox: Just connect - tunneling is automatic"
    echo ""
    echo -e "${CYAN}Example Commands (coming soon):${NC}"
    echo "  nebulabox remote connect user@remote-server"
    echo "  nebulabox remote deploy buildspec.json --target remote"
    echo "  nebulabox remote ps --target remote"
    echo ""
    echo -e "${YELLOW}⚠️  Remote commands will be available in Phase 4${NC}"
    echo ""
    print_menu
    get_choice
}

# Step 7: Unified Deployment
step_7_deploy() {
    print_header
    echo -e "${PURPLE}Step 7: Unified Cloud Deployment${NC}"
    echo ""
    echo -e "${BLUE}NebulaBox deploys everything to ONE platform:${NC}"
    echo ""
    echo -e "${GREEN}Current Fragmented Approach:${NC}"
    echo "  • Frontend → Vercel (different config)"
    echo "  • Backend → Render (different config)"
    echo "  • Database → MongoDB Atlas (separate service)"
    echo "  • Result: 3 different places, 3 different configs"
    echo ""
    echo -e "${GREEN}NebulaBox Unified Approach:${NC}"
    echo "  • Everything → NebulaBox Cloud (single platform)"
    echo "  • Same BuildSpec for all environments"
    echo "  • Single URL, single configuration"
    echo "  • Result: One place, one config, unified access"
    echo ""
    echo -e "${CYAN}Example Commands (coming soon):${NC}"
    echo "  nebulabox cloud login"
    echo "  nebulabox cloud deploy"
    echo "  nebulabox cloud deployments"
    echo ""
    echo -e "${YELLOW}⚠️  Cloud commands will be available in Phase 6${NC}"
    echo ""
    print_menu
    get_choice
}

# Step 8: Summary
step_8_summary() {
    print_header
    echo -e "${PURPLE}Step 8: Summary - What Makes NebulaBox Unique${NC}"
    echo ""
    echo -e "${GREEN}✅ Key Differentiators:${NC}"
    echo ""
    echo -e "1. ${CYAN}Unified Development${NC}"
    echo "   • Single container for entire stack (MERN, LAMP, etc.)"
    echo "   • No Docker Compose complexity"
    echo "   • Faster local development"
    echo ""
    echo -e "2. ${CYAN}Built-in Collaboration${NC}"
    echo "   • Real-time shared workspaces"
    echo "   • No VPN/ngrok setup needed"
    echo "   • Team members can code together instantly"
    echo ""
    echo -e "3. ${CYAN}Flexible Architecture Testing${NC}"
    echo "   • Test monolithic, microservices, or custom architectures"
    echo "   • Switch strategies easily"
    echo "   • Perfect for POC and experimentation"
    echo ""
    echo -e "4. ${CYAN}Unified Deployment${NC}"
    echo "   • Deploy everything to one platform"
    echo "   • Same BuildSpec for all environments"
    echo "   • No platform fragmentation"
    echo ""
    echo -e "5. ${CYAN}Simpler Than Kubernetes${NC}"
    echo "   • No complex YAML configurations"
    echo "   • No need for k8s expertise"
    echo "   • Simple JSON-based BuildSpec"
    echo ""
    echo -e "${BLUE}NebulaBox is not trying to replace Docker or Kubernetes.${NC}"
    echo -e "${BLUE}It's a different approach - unified, simple, and collaboration-first.${NC}"
    echo ""
    echo -e "${GREEN}🎉 Demo Complete!${NC}"
    echo ""
    print_menu
    get_choice
}

# Main demo loop
main() {
    local step=1
    local choice_result=0
    
    while true; do
        case $step in
            1)
                step_1_intro
                choice_result=$?
                ;;
            2)
                step_2_build
                choice_result=$?
                ;;
            3)
                step_3_run
                choice_result=$?
                ;;
            4)
                step_4_group
                choice_result=$?
                ;;
            5)
                step_5_workspace
                choice_result=$?
                ;;
            6)
                step_6_remote
                choice_result=$?
                ;;
            7)
                step_7_deploy
                choice_result=$?
                ;;
            8)
                step_8_summary
                choice_result=$?
                ;;
            *)
                echo -e "${GREEN}Demo complete!${NC}"
                exit 0
                ;;
        esac
        
        case $choice_result in
            0)  # Continue
                CURRENT_STEP=$step
                step=$((step + 1))
                if [ $step -gt $TOTAL_STEPS ]; then
                    step=8  # Go to summary
                fi
                ;;
            1)  # Skip
                CURRENT_STEP=$step
                step=$((step + 1))
                if [ $step -gt $TOTAL_STEPS ]; then
                    step=8
                fi
                ;;
            2)  # Back
                if [ $step -gt 1 ]; then
                    step=$((step - 1))
                    CURRENT_STEP=$step
                else
                    echo -e "${YELLOW}Already at first step${NC}"
                    sleep 1
                fi
                ;;
            3)  # Restart
                step=1
                CURRENT_STEP=0
                DEMO_IMAGE=""
                DEMO_CONTAINER=""
                DEMO_GROUP=""
                ;;
            4)  # Quit
                echo -e "${YELLOW}Exiting demo...${NC}"
                exit 0
                ;;
            5)  # Invalid
                continue
                ;;
        esac
    done
}

# Run demo
main

