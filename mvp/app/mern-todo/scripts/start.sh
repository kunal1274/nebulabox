#!/bin/sh

# Start MongoDB in background
echo "🔄 Starting MongoDB..."
mongod --dbpath /app/data --logpath /app/logs/mongodb.log --fork --bind_ip_all

# Wait for MongoDB to be ready
echo "⏳ Waiting for MongoDB to be ready..."
sleep 5

# Check if MongoDB is running
if ! pgrep mongod > /dev/null; then
  echo "❌ MongoDB failed to start"
  exit 1
fi

echo "✅ MongoDB started"

# Start Backend API
echo "🔄 Starting Backend API..."
cd /app/backend
node server.js &
BACKEND_PID=$!

# Wait for backend to be ready
echo "⏳ Waiting for Backend API to be ready..."
sleep 3

# Check if backend is running
if ! kill -0 $BACKEND_PID 2>/dev/null; then
  echo "❌ Backend API failed to start"
  exit 1
fi

echo "✅ Backend API started (PID: $BACKEND_PID)"

# Start Nginx
echo "🔄 Starting Nginx..."
nginx -g "daemon off;" &
NGINX_PID=$!

# Wait for nginx to be ready
sleep 2

echo "✅ Nginx started (PID: $NGINX_PID)"
echo ""
echo "🎉 MERN Todo App is running!"
echo "📡 Frontend: http://localhost:3000"
echo "🔌 Backend API: http://localhost:5000/api"
echo "💾 MongoDB: mongodb://localhost:27017/todos"
echo ""

# Keep script running and handle shutdown
trap "echo '🛑 Shutting down...'; kill $BACKEND_PID $NGINX_PID; mongod --shutdown; exit" SIGTERM SIGINT

# Wait for processes
wait

