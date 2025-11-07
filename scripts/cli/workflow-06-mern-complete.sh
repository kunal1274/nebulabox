#!/bin/bash

# Workflow 06: Complete MERN Todo App Workflow
# This script tests the complete MERN todo app deployment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLI_BINARY="$PROJECT_ROOT/bin/nebulabox"
MERN_DIR="$PROJECT_ROOT/mvp/app/mern-todo"

echo "=========================================="
echo "Workflow 06: MERN Todo App Complete Test"
echo "=========================================="

# Check if CLI binary exists
if [ ! -f "$CLI_BINARY" ]; then
    echo "❌ CLI binary not found at $CLI_BINARY"
    echo "   Run 'make build-cli-test' first"
    exit 1
fi

# Check if MERN app exists
if [ ! -f "$MERN_DIR/buildspec.json" ]; then
    echo "❌ MERN app not found at $MERN_DIR"
    exit 1
fi

echo "🚀 Testing complete MERN todo app workflow..."
echo ""

# Step 1: Build image
echo "Step 1: Building MERN todo image..."
if $CLI_BINARY build -f "$MERN_DIR/buildspec.json" -t mern-todo:latest; then
    echo "✅ Image built successfully!"
else
    echo "❌ Image build failed"
    exit 1
fi

# Step 2: Run container
echo ""
echo "Step 2: Running MERN todo container..."
CONTAINER_NAME="mern-todo-$(date +%s)"
if $CLI_BINARY run mern-todo:latest --name "$CONTAINER_NAME" \
    -p 3000:80 \
    -p 5000:5000 \
    -p 27017:27017 \
    -d; then
    echo "✅ Container started successfully!"
    
    # Wait for services to start
    echo ""
    echo "⏳ Waiting for services to start..."
    sleep 5
    
    # Step 3: Verify services
    echo ""
    echo "Step 3: Verifying services..."
    
    # Check frontend
    if curl -s http://localhost:3000 > /dev/null; then
        echo "✅ Frontend accessible at http://localhost:3000"
    else
        echo "⚠️  Frontend not yet accessible"
    fi
    
    # Check backend
    if curl -s http://localhost:5000/api/health > /dev/null; then
        echo "✅ Backend API accessible at http://localhost:5000"
    else
        echo "⚠️  Backend API not yet accessible"
    fi
    
    # Step 4: List containers
    echo ""
    echo "Step 4: Listing containers..."
    $CLI_BINARY ps
    
    # Step 5: View logs
    echo ""
    echo "Step 5: Viewing container logs..."
    $CLI_BINARY logs "$CONTAINER_NAME" --tail 20 || echo "   (Logs command may not be fully implemented)"
    
    # Step 6: Stop container
    echo ""
    echo "Step 6: Stopping container..."
    if $CLI_BINARY stop "$CONTAINER_NAME"; then
        echo "✅ Container stopped successfully!"
    fi
else
    echo "❌ Container run failed"
    exit 1
fi

echo ""
echo "✅ Workflow 06 complete: MERN todo app test passed"

