#!/bin/bash

# Deploy MERN MVP using NebulaBox Build Specification
# This script automates building from buildspec.json and deploying

set -e

API_URL="${NEBULA_API_URL:-http://localhost:8081/api}"
BUILD_SPEC_PATH="${1:-../app/single-container/buildspec.json}"
BUILD_CONTEXT="${2:-../app/single-container}"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🚀 Deploying MERN MVP using NebulaBox Build Specification${NC}"
echo "API URL: $API_URL"
echo "Build Spec: $BUILD_SPEC_PATH"
echo "Build Context: $BUILD_CONTEXT"
echo ""

# Step 1: Validate Build Specification
echo -e "${BLUE}Step 1: Validating build specification...${NC}"
VALIDATION=$(curl -s -X POST "$API_URL/buildspec/validate" \
  -H "Content-Type: application/json" \
  -d @"$BUILD_SPEC_PATH")

if echo "$VALIDATION" | grep -q '"valid":true'; then
  echo -e "${GREEN}✅ Build specification is valid${NC}"
else
  echo -e "${RED}❌ Build specification validation failed:${NC}"
  echo "$VALIDATION"
  exit 1
fi

# Step 2: Build Image from Build Specification
echo -e "\n${BLUE}Step 2: Building image from build specification...${NC}"
echo -e "${YELLOW}Note: This requires buildspec.json content and build context${NC}"
echo -e "${YELLOW}For now, please use the dashboard Build Spec page to build${NC}"
echo -e "${YELLOW}Or use Docker directly after converting buildspec${NC}"

# Convert to Dockerfile first
echo -e "\n${BLUE}Step 2a: Converting buildspec to Dockerfile...${NC}"
DOCKERFILE=$(curl -s -X POST "$API_URL/buildspec/convert" \
  -H "Content-Type: application/json" \
  -d @"$BUILD_SPEC_PATH")

if [ -z "$DOCKERFILE" ]; then
  echo -e "${RED}❌ Failed to convert buildspec${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Buildspec converted to Dockerfile${NC}"
echo -e "${YELLOW}Preview (first 10 lines):${NC}"
echo "$DOCKERFILE" | head -10

# Save Dockerfile for manual build if needed
DOCKERFILE_PATH="/tmp/mern-mvp.Dockerfile"
echo "$DOCKERFILE" > "$DOCKERFILE_PATH"
echo -e "${BLUE}Full Dockerfile saved to: $DOCKERFILE_PATH${NC}"

# Step 3: Build using NebulaBox API (NebulaBox native, not Docker CLI)
echo -e "\n${BLUE}Step 3: Building image using NebulaBox API...${NC}"
echo -e "${YELLOW}Note: Using NebulaBox Build API, not Docker CLI${NC}"

# Read buildspec.json content
BUILD_SPEC_CONTENT=$(cat "$BUILD_SPEC_PATH")

# Build via NebulaBox API
echo -e "${BLUE}Building via: POST $API_URL/buildspec/build${NC}"
BUILD_RESULT=$(curl -s -X POST "$API_URL/buildspec/build" \
  -H "Content-Type: application/json" \
  -d "$BUILD_SPEC_CONTENT")

if echo "$BUILD_RESULT" | grep -q '"image"'; then
  echo -e "${GREEN}✅ Image built successfully via NebulaBox${NC}"
  echo "$BUILD_RESULT"
elif echo "$BUILD_RESULT" | grep -q "error\|failed"; then
  echo -e "${RED}❌ Build failed via NebulaBox API:${NC}"
  echo "$BUILD_RESULT"
  echo -e "${YELLOW}Fallback: Please build using NebulaBox Dashboard Build Spec page${NC}"
else
  echo -e "${YELLOW}⚠️  Build API may not be fully implemented. Using dashboard build:${NC}"
  echo "  1. Go to NebulaBox Dashboard"
  echo "  2. Navigate to 'Build Spec' page"
  echo "  3. Upload or paste buildspec.json"
  echo "  4. Click 'Build from Spec'"
  echo "  5. Wait for build to complete"
fi

# Step 4: Create Container
echo -e "\n${BLUE}Step 4: Creating container...${NC}"
CONTAINER_CREATE=$(curl -s -X POST "$API_URL/containers/run" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "mern-mvp:latest",
    "name": "mern-mvp",
    "ports": ["3000:3000", "5000:5000", "27017:27017"],
    "detach": true,
    "env": [
      "MONGO_HOST=localhost",
      "MONGO_PORT=27017",
      "MONGO_DB=todos",
      "PORT=5000"
    ]
  }')

if echo "$CONTAINER_CREATE" | grep -q '"id"'; then
  echo -e "${GREEN}✅ Container created${NC}"
  CONTAINER_ID=$(echo "$CONTAINER_CREATE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
  echo "Container ID: $CONTAINER_ID"
else
  echo -e "${YELLOW}⚠️  Container may already exist or error occurred${NC}"
  echo "$CONTAINER_CREATE"
fi

# Step 5: Start Container
echo -e "\n${BLUE}Step 5: Starting container...${NC}"
sleep 2
START_RESULT=$(curl -s -X POST "$API_URL/containers/mern-mvp/start" || echo "{}")

if echo "$START_RESULT" | grep -q '"status":"running"'; then
  echo -e "${GREEN}✅ Container started${NC}"
else
  echo -e "${YELLOW}⚠️  Starting container...${NC}"
fi

# Step 6: Verify Deployment
echo -e "\n${BLUE}Step 6: Verifying deployment...${NC}"
sleep 3

# Check container status
CONTAINER_STATUS=$(curl -s "$API_URL/containers" | grep -o '"name":"mern-mvp"[^}]*' || echo "")
if [ -n "$CONTAINER_STATUS" ]; then
  echo -e "${GREEN}✅ Container found${NC}"
  echo "$CONTAINER_STATUS"
else
  echo -e "${YELLOW}⚠️  Container status check failed${NC}"
fi

# Check health
echo -e "\n${BLUE}Checking application health...${NC}"
sleep 5
if curl -s -f http://localhost:5000/health > /dev/null; then
  echo -e "${GREEN}✅ Backend health check passed${NC}"
else
  echo -e "${YELLOW}⚠️  Backend not ready yet (may need more time)${NC}"
fi

echo -e "\n${GREEN}✅ Deployment complete!${NC}"
echo -e "\n${YELLOW}Access points:${NC}"
echo "  Frontend: http://localhost:3000"
echo "  Backend API: http://localhost:5000/api/todos"
echo "  Health Check: http://localhost:5000/health"
echo ""
echo -e "${BLUE}To check container status:${NC}"
echo "  curl $API_URL/containers"
echo ""
echo -e "${BLUE}To view logs:${NC}"
echo "  curl $API_URL/containers/mern-mvp/logs"

