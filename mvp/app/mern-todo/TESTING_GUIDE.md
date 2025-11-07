# MERN Todo App - Testing Guide

This guide shows how to build, run, and test the MERN Todo App using NebulaBox via CLI and Dashboard.

## 📋 App Overview

**MERN Todo App** is a full-stack application:
- **MongoDB**: Database for storing todos
- **Express**: Backend API server (port 5000)
- **React**: Frontend UI (served via Nginx on port 80)
- **Node.js**: Runtime environment

## 🚀 Method 1: Build & Run via Dashboard (Recommended)

### Step 1: Build via Dashboard BuildSpec Page

1. **Start NebulaBox API Server** (if not running):
   ```bash
   cd /home/serverratxen/Documents/cursor-projects/nebulabox
   make run-api
   # Or: go run cmd/api/main.go
   ```

2. **Start Dashboard** (if not running):
   ```bash
   cd web/dashboard
   npm run dev
   # Dashboard runs on http://localhost:3001
   ```

3. **Navigate to BuildSpec Page**:
   - Open http://localhost:3001/buildspec
   - Click "Load Example" or paste the buildspec content

4. **Load BuildSpec**:
   ```bash
   # Copy the buildspec.json content
   cat mvp/app/mern-todo/buildspec.json
   ```
   - Paste into the BuildSpec editor
   - Click "Validate" to verify
   - Click "Build" to build the image

5. **Verify Build**:
   - Check the "Build Logs" tab for build progress
   - Image should be tagged as `mern-todo:latest`

### Step 2: Run Container via Dashboard

1. **Navigate to Containers Page**:
   - Go to http://localhost:3001/containers
   - Click "Create New Container"

2. **Configure Container**:
   - **Image**: `mern-todo:latest`
   - **Name**: `mern-todo-app`
   - **Ports**: Map port 3000 (host) → 80 (container) for frontend
   - **Ports**: Map port 5000 (host) → 5000 (container) for backend API
   - Click "Create"

3. **Verify Container**:
   - Container should appear in the list
   - Status should be "running"
   - Check logs if needed

### Step 3: Test the App

1. **Access Frontend**:
   - Open http://localhost:3000 in your browser
   - You should see the MERN Todo App UI

2. **Test Functionality**:
   - ✅ Add a new todo
   - ✅ Mark todo as complete
   - ✅ Delete a todo
   - ✅ Refresh the list

3. **Check Backend API**:
   - Test health: http://localhost:5000/api/health
   - Get todos: http://localhost:5000/api/todos
   - Should return JSON array of todos

## 🔧 Method 2: Build & Run via CLI

### Step 1: Build via CLI (Using BuildSpec API)

The CLI doesn't directly support buildspec files yet, but you can use the API:

```bash
# Set API URL
export NEBULABOX_API_URL=http://localhost:8081

# Build using buildspec via API
curl -X POST http://localhost:8081/api/buildspec/build \
  -H "Content-Type: application/json" \
  -d @mvp/app/mern-todo/buildspec.json
```

Or use the dashboard's buildspec page (Method 1) which is easier.

### Step 2: Run Container via CLI

```bash
# Run the container
./bin/nebulabox run mern-todo:latest \
  --name mern-todo-app \
  --port 3000:80 \
  --port 5000:5000

# Or using the API URL environment variable
NEBULABOX_API_URL=http://localhost:8081 \
  ./bin/nebulabox run mern-todo:latest \
  --name mern-todo-app \
  --port 3000:80 \
  --port 5000:5000
```

### Step 3: Verify Container

```bash
# List containers
./bin/nebulabox list

# Check logs
./bin/nebulabox logs mern-todo-app
```

### Step 4: Test the App

Same as Method 1, Step 3 above.

## 🧪 Manual Testing Checklist

### Frontend Tests
- [ ] Page loads without errors
- [ ] Can add a new todo
- [ ] Can mark todo as complete
- [ ] Can delete a todo
- [ ] Todos persist after refresh
- [ ] Error handling works (disconnect backend)

### Backend API Tests
- [ ] `GET /api/health` returns 200
- [ ] `GET /api/todos` returns todos array
- [ ] `POST /api/todos` creates new todo
- [ ] `PUT /api/todos/:id` updates todo
- [ ] `DELETE /api/todos/:id` deletes todo

### Integration Tests
- [ ] Frontend can communicate with backend
- [ ] MongoDB stores data correctly
- [ ] Data persists across container restarts (if using volumes)
- [ ] Health checks work
- [ ] Logs are accessible

## 📊 Monitoring

### View Container Logs
```bash
# Via CLI
./bin/nebulabox logs mern-todo-app

# Via Dashboard
# Navigate to Containers → Click on container → View Logs
```

### Check Container Status
```bash
# Via CLI
./bin/nebulabox list

# Via Dashboard
# Navigate to Containers page
```

### Monitor Resources
- Check CPU/Memory usage in Dashboard
- View metrics in Monitor page
- Check health status

## 🐛 Troubleshooting

### Container won't start
- Check logs: `./bin/nebulabox logs mern-todo-app`
- Verify image exists: Check Images page
- Check port conflicts: Ensure ports 3000 and 5000 are free

### Backend API not responding
- Check MongoDB is running inside container
- Verify backend logs
- Test health endpoint: `curl http://localhost:5000/api/health`

### Frontend can't connect to backend
- Verify REACT_APP_API_URL environment variable
- Check nginx proxy configuration
- Ensure backend is running on port 5000

### MongoDB connection issues
- Check MongoDB logs in container
- Verify MONGODB_URI environment variable
- Ensure MongoDB data directory has proper permissions

## 📝 Next Steps

1. **E2E Testing**: Create Playwright tests for the todo app
2. **Performance Testing**: Load test the API endpoints
3. **CI/CD Integration**: Set up automated builds
4. **Multi-container**: Split into separate containers (backend, frontend, db)

## 🔗 Useful Commands

```bash
# Stop container
./bin/nebulabox stop mern-todo-app

# Remove container
# (via dashboard or API)

# Rebuild image
# Use dashboard BuildSpec page or API

# View all containers
./bin/nebulabox list

# View all images
# (via dashboard Images page)
```

