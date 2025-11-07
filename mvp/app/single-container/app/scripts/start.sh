#!/bin/sh

# Start script for MERN stack in single container
# Starts: MongoDB → Backend → Frontend

set -e

echo "🚀 Starting MERN Stack Application..."
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
MONGO_DATA_DIR=${MONGO_DATA_DIR:-/app/data}
MONGO_LOG_DIR=${MONGO_LOG_DIR:-/app/logs}
BACKEND_DIR=/app/backend
FRONTEND_DIR=/app/frontend

# Create directories
mkdir -p "$MONGO_DATA_DIR" "$MONGO_LOG_DIR"

# Start MongoDB
echo -e "${BLUE}Starting MongoDB...${NC}"
mongod --dbpath "$MONGO_DATA_DIR" \
       --logpath "$MONGO_LOG_DIR/mongodb.log" \
       --fork \
       --bind_ip 0.0.0.0 \
       --port ${MONGO_PORT:-27017}

# Wait for MongoDB to be ready
echo -e "${BLUE}Waiting for MongoDB to be ready...${NC}"
sleep 3
for i in 1 2 3 4 5; do
  if mongosh --eval "db.adminCommand('ping')" --quiet > /dev/null 2>&1; then
    echo -e "${GREEN}✅ MongoDB is ready!${NC}"
    break
  fi
  if [ $i -eq 5 ]; then
    echo -e "${YELLOW}⚠️  MongoDB may still be starting...${NC}"
  fi
  sleep 1
done

# Start Backend
echo -e "${BLUE}Starting Backend API...${NC}"
cd "$BACKEND_DIR"
MONGO_HOST=${MONGO_HOST:-localhost} \
MONGO_PORT=${MONGO_PORT:-27017} \
MONGO_DB=${MONGO_DB:-todos} \
PORT=${PORT:-5000} \
node src/server.js &
BACKEND_PID=$!

# Wait for backend to be ready
echo -e "${BLUE}Waiting for Backend to be ready...${NC}"
sleep 3
for i in 1 2 3 4 5; do
  if curl -s http://localhost:5000/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend API is ready!${NC}"
    break
  fi
  if [ $i -eq 5 ]; then
    echo -e "${YELLOW}⚠️  Backend may still be starting...${NC}"
  fi
  sleep 1
done

# Start Frontend (serve built files)
echo -e "${BLUE}Starting Frontend...${NC}"
cd "$FRONTEND_DIR"

# Check if frontend is built
if [ ! -d "dist" ]; then
  echo -e "${YELLOW}Frontend not built, building now...${NC}"
  npm install --silent
  npm run build --silent
fi

# Serve frontend
npm run serve &
FRONTEND_PID=$!

echo ""
echo -e "${GREEN}═══════════════════════════════════════${NC}"
echo -e "${GREEN}   MERN Stack Started Successfully!${NC}"
echo -e "${GREEN}═══════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}📊 Services:${NC}"
echo "   MongoDB:  localhost:${MONGO_PORT:-27017}"
echo "   Backend:   http://localhost:5000"
echo "   Frontend:  http://localhost:3000"
echo ""
echo -e "${BLUE}📋 Health Checks:${NC}"
echo "   Backend:  http://localhost:5000/health"
echo "   API:      http://localhost:5000/api/todos"
echo ""
echo -e "${GREEN}✅ All services running!${NC}"
echo ""

# Keep script running and handle shutdown
trap "echo 'Shutting down...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; mongod --shutdown --dbpath $MONGO_DATA_DIR 2>/dev/null; exit" SIGTERM SIGINT

# Wait for processes
wait

