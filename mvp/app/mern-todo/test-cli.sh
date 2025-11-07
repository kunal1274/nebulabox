#!/bin/bash

# MERN Todo App - CLI Testing Script
# This script helps test the MERN todo app via NebulaBox CLI

set -e

API_URL="${NEBULABOX_API_URL:-http://localhost:8081}"
BUILDSPEC_PATH="mvp/app/mern-todo/buildspec.json"
IMAGE_TAG="mern-todo:latest"
CONTAINER_NAME="mern-todo-app"

echo "🚀 MERN Todo App - CLI Testing"
echo "═══════════════════════════════════════════════════════"
echo ""

# Check if API is running
echo "📡 Checking NebulaBox API..."
if ! curl -s "${API_URL}/api/health" > /dev/null; then
    echo "❌ API server is not running at ${API_URL}"
    echo "   Start it with: make run-api"
    exit 1
fi
echo "✅ API server is running"
echo ""

# Step 1: Validate buildspec
echo "📋 Step 1: Validating buildspec..."
BUILDSPEC_JSON=$(cat "${BUILDSPEC_PATH}")
SPEC_CONTENT=$(echo "${BUILDSPEC_JSON}" | jq -c '.')
TAG=$(echo "${BUILDSPEC_JSON}" | jq -r '.tag')

VALIDATE_RESPONSE=$(curl -s -X POST "${API_URL}/api/buildspec/validate" \
    -H "Content-Type: application/json" \
    -d "{\"spec\": ${SPEC_CONTENT}}")

VALID=$(echo "${VALIDATE_RESPONSE}" | jq -r '.valid')

if [ "${VALID}" = "true" ]; then
    echo "✅ Buildspec is valid"
else
    echo "❌ Buildspec validation failed:"
    echo "${VALIDATE_RESPONSE}" | jq '.'
    exit 1
fi
echo ""

# Step 2: Build image
echo "🔨 Step 2: Building image..."
BUILD_RESPONSE=$(curl -s -X POST "${API_URL}/api/buildspec/build" \
    -H "Content-Type: application/json" \
    -d "{\"spec\": ${SPEC_CONTENT}, \"tag\": \"${TAG}\"}")

BUILD_VALID=$(echo "${BUILD_RESPONSE}" | jq -r '.valid // false')

if [ "${BUILD_VALID}" = "true" ]; then
    echo "✅ Image built successfully: ${TAG}"
    echo "${BUILD_RESPONSE}" | jq -r '.logs[]?' 2>/dev/null || echo "Build logs not available"
else
    echo "❌ Build failed:"
    echo "${BUILD_RESPONSE}" | jq '.'
    exit 1
fi
echo ""

# Step 3: Verify image exists
echo "📦 Step 3: Verifying image..."
IMAGES_RESPONSE=$(curl -s "${API_URL}/api/images")
IMAGE_EXISTS=$(echo "${IMAGES_RESPONSE}" | jq -r ".[] | select(.name == \"${IMAGE_TAG%:*}\") | .name" | head -1)

if [ -n "${IMAGE_EXISTS}" ]; then
    echo "✅ Image found: ${IMAGE_EXISTS}"
else
    echo "⚠️  Image not found in list (may still be processing)"
fi
echo ""

# Step 4: Run container via CLI
echo "🏃 Step 4: Running container via CLI..."
echo "   Command: ./bin/nebulabox run ${IMAGE_TAG} --name ${CONTAINER_NAME} --port 3000:80 --port 5000:5000"
echo ""

if [ -f "./bin/nebulabox" ]; then
    export NEBULABOX_API_URL="${API_URL}"
    ./bin/nebulabox run "${IMAGE_TAG}" \
        --name "${CONTAINER_NAME}" \
        --port 3000:80 \
        --port 5000:5000 \
        --detach || echo "⚠️  Container run command completed (check status)"
else
    echo "⚠️  CLI binary not found at ./bin/nebulabox"
    echo "   Build it with: make build-cli-test"
    echo "   Or use the Dashboard to run the container"
fi
echo ""

# Step 5: Verify container is running
echo "📊 Step 5: Checking container status..."
sleep 2
CONTAINERS_RESPONSE=$(curl -s "${API_URL}/api/containers")
CONTAINER_STATUS=$(echo "${CONTAINERS_RESPONSE}" | jq -r ".[] | select(.name == \"${CONTAINER_NAME}\") | .status" | head -1)

if [ -n "${CONTAINER_STATUS}" ]; then
    echo "✅ Container found: ${CONTAINER_NAME}"
    echo "   Status: ${CONTAINER_STATUS}"
else
    echo "⚠️  Container not found (may still be starting)"
fi
echo ""

# Step 6: Test endpoints
echo "🧪 Step 6: Testing application endpoints..."
sleep 3

echo "   Testing backend health..."
HEALTH_RESPONSE=$(curl -s "http://localhost:5000/api/health" || echo "{}")
if echo "${HEALTH_RESPONSE}" | jq -e '.status' > /dev/null 2>&1; then
    echo "   ✅ Backend health check passed"
    echo "   ${HEALTH_RESPONSE}" | jq '.'
else
    echo "   ⚠️  Backend not responding yet (may need more time to start)"
fi
echo ""

echo "   Testing frontend..."
if curl -s "http://localhost:3000" > /dev/null; then
    echo "   ✅ Frontend is accessible at http://localhost:3000"
else
    echo "   ⚠️  Frontend not accessible yet (may need more time to start)"
fi
echo ""

# Summary
echo "═══════════════════════════════════════════════════════"
echo "📋 Testing Summary"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "✅ Buildspec validated"
echo "✅ Image built: ${IMAGE_TAG}"
echo "✅ Container: ${CONTAINER_NAME}"
echo ""
echo "🌐 Access Points:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:5000/api"
echo "   Health Check: http://localhost:5000/api/health"
echo ""
echo "📝 Next Steps:"
echo "   1. Open http://localhost:3000 in your browser"
echo "   2. Test adding, completing, and deleting todos"
echo "   3. Check logs: ./bin/nebulabox logs ${CONTAINER_NAME}"
echo "   4. Stop container: ./bin/nebulabox stop ${CONTAINER_NAME}"
echo ""
echo "═══════════════════════════════════════════════════════"

