#!/bin/bash

# NebulaBox API Test Script
# This script tests all API endpoints

API_BASE="http://localhost:8081/api"

echo "🧪 Testing NebulaBox API Server"
echo "================================"
echo ""

# Test health check
echo "1. Testing health check..."
curl -s "$API_BASE/health" | jq '.' 2>/dev/null || curl -s "$API_BASE/health"
echo -e "\n"

# Test system stats
echo "2. Testing system stats..."
curl -s "$API_BASE/system/stats" | jq '.' 2>/dev/null || curl -s "$API_BASE/system/stats"
echo -e "\n"

# Test list containers
echo "3. Testing list containers..."
curl -s "$API_BASE/containers" | jq '.' 2>/dev/null || curl -s "$API_BASE/containers"
echo -e "\n"

# Test list images
echo "4. Testing list images..."
curl -s "$API_BASE/images" | jq '.' 2>/dev/null || curl -s "$API_BASE/images"
echo -e "\n"

# Test create container
echo "5. Testing create container..."
curl -s -X POST "$API_BASE/containers/run" \
  -H "Content-Type: application/json" \
  -d '{"image": "nginx:latest", "name": "test-container", "detach": true}' \
  | jq '.' 2>/dev/null || curl -s -X POST "$API_BASE/containers/run" \
  -H "Content-Type: application/json" \
  -d '{"image": "nginx:latest", "name": "test-container", "detach": true}'
echo -e "\n"

# Test pull image
echo "6. Testing pull image..."
curl -s -X POST "$API_BASE/images/pull" \
  -H "Content-Type: application/json" \
  -d '{"image": "alpine:latest"}' \
  | jq '.' 2>/dev/null || curl -s -X POST "$API_BASE/images/pull" \
  -H "Content-Type: application/json" \
  -d '{"image": "alpine:latest"}'
echo -e "\n"

echo "✅ API testing complete!"
echo ""
echo "💡 Note: Some endpoints may show mock data or placeholders"
echo "   This is expected behavior for the current development phase"
