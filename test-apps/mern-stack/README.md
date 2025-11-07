# Simple MERN Stack Application for Testing

This is a simple Todo application built with MongoDB, Express, React, and Node.js for testing NebulaBox platform.

## Application Structure

```
mern-stack/
├── frontend/          # React application
│   ├── src/
│   ├── package.json
│   └── Dockerfile
├── backend/           # Express API
│   ├── src/
│   ├── package.json
│   └── Dockerfile
└── docker-compose.yml # For reference (local testing)
```

## Frontend (React)

**Port**: 3000

**Features**:
- Simple todo list UI
- Add/Edit/Delete todos
- Mark todos as complete
- Connects to backend API

**Environment Variables**:
- `REACT_APP_API_URL`: Backend API URL (default: http://localhost:5000)

## Backend (Express + Node.js)

**Port**: 5000

**Features**:
- REST API for todos
- MongoDB connection
- CORS enabled

**API Endpoints**:
- `GET /api/todos` - List all todos
- `POST /api/todos` - Create new todo
- `PUT /api/todos/:id` - Update todo
- `DELETE /api/todos/:id` - Delete todo
- `GET /health` - Health check endpoint

**Environment Variables**:
- `MONGO_HOST`: MongoDB host (default: mongodb)
- `MONGO_PORT`: MongoDB port (default: 27017)
- `MONGO_DB`: Database name (default: todos)
- `PORT`: Server port (default: 5000)

## MongoDB

**Port**: 27017

**Features**:
- Persistent data storage
- No authentication (for simplicity in testing)

## Quick Start (Local Testing)

```bash
# Using Docker Compose
docker-compose up -d

# Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:5000
# MongoDB: localhost:27017
```

## Deployment to NebulaBox

### Step 1: Build Images

```bash
cd frontend
docker build -t mern-frontend:latest .
docker tag mern-frontend:latest localhost:5000/mern-frontend:latest

cd ../backend
docker build -t mern-backend:latest .
docker tag mern-backend:latest localhost:5000/mern-backend:latest
```

### Step 2: Push to Nebula Registry

```bash
# Login to registry
docker login localhost:5000

# Push images
docker push localhost:5000/mern-frontend:latest
docker push localhost:5000/mern-backend:latest
```

### Step 3: Create Network

```bash
nebula-cli networks create --name mern-network --driver bridge
```

### Step 4: Create Containers

```bash
# MongoDB
nebula-cli containers create \
  --name mongodb \
  --image mongo:latest \
  --network mern-network \
  --port 27017:27017

# Backend
nebula-cli containers create \
  --name backend \
  --image localhost:5000/mern-backend:latest \
  --network mern-network \
  --port 5000:5000 \
  --env MONGO_HOST=mongodb \
  --env MONGO_PORT=27017 \
  --env MONGO_DB=todos

# Frontend
nebula-cli containers create \
  --name frontend \
  --image localhost:5000/mern-frontend:latest \
  --network mern-network \
  --port 3000:3000 \
  --env REACT_APP_API_URL=http://backend:5000
```

### Step 5: Start Containers

```bash
# Start in order: MongoDB -> Backend -> Frontend
nebula-cli containers start mongodb
sleep 5  # Wait for MongoDB to be ready
nebula-cli containers start backend
sleep 2
nebula-cli containers start frontend
```

### Step 6: Access Application

- Frontend: http://localhost:3000 (or mapped host port)
- Backend API: http://localhost:5000/api/todos
- Health Check: http://localhost:5000/health

## Testing Checklist

- [ ] All containers start successfully
- [ ] Frontend can connect to backend
- [ ] Backend can connect to MongoDB
- [ ] API endpoints work correctly
- [ ] Todos can be created, updated, deleted
- [ ] Data persists after container restart
- [ ] Health checks pass
- [ ] Monitoring shows resource usage

