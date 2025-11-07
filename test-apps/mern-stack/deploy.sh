#!/bin/bash

# MERN Stack Deployment Script for NebulaBox
# This script automates the deployment of the MERN stack to NebulaBox

set -e

API_URL="${NEBULA_API_URL:-http://localhost:8081/api}"
REGISTRY_URL="${NEBULA_REGISTRY_URL:-localhost:5000}"

echo "🚀 Deploying MERN Stack to NebulaBox..."
echo "API URL: $API_URL"
echo "Registry URL: $REGISTRY_URL"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Create Network
echo -e "\n${BLUE}Step 1: Creating network...${NC}"
curl -X POST "$API_URL/networks" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mern-network",
    "driver": "bridge",
    "subnet": "172.20.0.0/16"
  }' || echo "Network may already exist"

# Step 2: Create MongoDB Container
echo -e "\n${BLUE}Step 2: Creating MongoDB container...${NC}"
curl -X POST "$API_URL/containers/run" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "mongo:latest",
    "name": "mern-mongodb",
    "network": "mern-network",
    "ports": ["27017:27017"],
    "detach": true
  }'

sleep 3

# Step 3: Create Backend Container
echo -e "\n${BLUE}Step 3: Creating backend container...${NC}"
curl -X POST "$API_URL/containers/run" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "'"$REGISTRY_URL"'/mern-backend:latest",
    "name": "mern-backend",
    "network": "mern-network",
    "ports": ["5000:5000"],
    "env": [
      "MONGO_HOST=mern-mongodb",
      "MONGO_PORT=27017",
      "MONGO_DB=todos",
      "PORT=5000"
    ],
    "detach": true
  }'

sleep 2

# Step 4: Create Frontend Container
echo -e "\n${BLUE}Step 4: Creating frontend container...${NC}"
curl -X POST "$API_URL/containers/run" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "'"$REGISTRY_URL"'/mern-frontend:latest",
    "name": "mern-frontend",
    "network": "mern-network",
    "ports": ["3000:3000"],
    "env": [
      "REACT_APP_API_URL=http://mern-backend:5000"
    ],
    "detach": true
  }'

# Step 5: Register Services
echo -e "\n${BLUE}Step 5: Registering services...${NC}"
curl -X POST "$API_URL/services/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mern-mongodb",
    "address": "mern-mongodb",
    "port": 27017
  }' || echo "Service registration skipped"

curl -X POST "$API_URL/services/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mern-backend",
    "address": "mern-backend",
    "port": 5000
  }' || echo "Service registration skipped"

curl -X POST "$API_URL/services/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mern-frontend",
    "address": "mern-frontend",
    "port": 3000
  }' || echo "Service registration skipped"

echo -e "\n${GREEN}✅ MERN Stack deployment complete!${NC}"
echo -e "\n${YELLOW}Access points:${NC}"
echo "  Frontend: http://localhost:3000"
echo "  Backend API: http://localhost:5000/api/todos"
echo "  Health Check: http://localhost:5000/health"

echo -e "\n${YELLOW}To check container status:${NC}"
echo "  curl $API_URL/containers"

