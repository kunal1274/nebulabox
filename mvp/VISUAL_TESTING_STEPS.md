# Visual Testing Steps - Real-time Dashboard Verification

## ✅ Completed Actions (via API)

1. ✅ Built image `mern-mvp:latest` from `buildspec.json`
2. ✅ Created container `mern-mvp-test` 
3. ✅ Started container
4. ✅ Verified container status

## 📺 Now Check Dashboard Visually

### Step 1: Open Dashboard
```
http://localhost:3001
```

### Step 2: Navigate to Containers Page
1. Click **"Containers"** in sidebar
2. Look for **"mern-mvp-test"** container
3. Verify:
   - ✅ Status: **"running"** (green badge)
   - ✅ Image: **"mern-mvp:latest"**
   - ✅ Ports: **3000:3000, 5000:5000, 27017:27017**

**What You Should See:**
- Container card with name, image, status
- Action buttons: Stop, Restart, Logs, Delete
- Port mappings displayed
- Created timestamp

### Step 3: View Container Details
1. Click on container name or **"View Details"**
2. Check:
   - Environment variables
   - Port mappings
   - Resource usage
   - Health status

### Step 4: Check Container Logs
1. Click **"Logs"** button on container card
2. Should see:
   - MongoDB startup messages
   - Backend API startup
   - Frontend build/serve messages
   - Health check status

### Step 5: Navigate to Images Page
1. Click **"Images"** in sidebar
2. Look for **"mern-mvp"** or **"mern-mvp:latest"**
3. Verify:
   - ✅ Image name and tag
   - ✅ Size information
   - ✅ Created timestamp

### Step 6: Navigate to Monitor Page
1. Click **"Monitor"** in sidebar
2. Check system metrics:
   - CPU usage graph
   - Memory usage graph
   - Container count
   - Network I/O
3. Look for **"mern-mvp-test"** in container list
4. View real-time stats

### Step 7: Test the Running Application

#### Frontend (Todo App)
```
http://localhost:3000
```
**Expected:**
- Todo list interface
- Add/Edit/Delete functionality
- Connection status indicator

#### Backend API Health
```
http://localhost:5000/health
```
**Expected Response:**
```json
{
  "status": "healthy",
  "service": "mern-backend",
  "database": "connected",
  "timestamp": "..."
}
```

#### Backend API (Todos)
```
http://localhost:5000/api/todos
```
**Expected:** List of todos (may be empty initially)

### Step 8: Test Container Lifecycle

1. **Stop Container:**
   - Click **"Stop"** button
   - Status should change to **"stopped"**
   - Application at localhost:3000 should be inaccessible

2. **Start Container:**
   - Click **"Start"** button
   - Status should change to **"running"**
   - Application should be accessible again

3. **Restart Container:**
   - Click **"Restart"** button
   - Container should restart
   - Status should remain **"running"**

### Step 9: Check API Responses

#### Get Container List
```bash
curl http://localhost:8081/api/containers | jq '.[] | select(.name == "mern-mvp-test")'
```

#### Get Container Details
```bash
CONTAINER_ID=$(curl -s http://localhost:8081/api/containers | jq -r '.[] | select(.name == "mern-mvp-test") | .id')
curl http://localhost:8081/api/containers/$CONTAINER_ID | jq '.'
```

#### Get Container Logs
```bash
curl http://localhost:8081/api/containers/$CONTAINER_ID/logs | jq '.logs | .[-10:]'
```

#### Get Container Stats
```bash
curl http://localhost:8081/api/containers/$CONTAINER_ID/stats | jq '.'
```

## 🎯 Comparison: Docker vs NebulaBox

| Docker Command | NebulaBox Equivalent |
|----------------|----------------------|
| `docker build -t mern-mvp:latest .` | Dashboard: Build Spec page → Build |
| `docker run -d -p 3000:3000 mern-mvp:latest` | Dashboard: Containers → Create → Start |
| `docker ps` | Dashboard: Containers page |
| `docker logs <id>` | Dashboard: Containers → Logs button |
| `docker stats <id>` | Dashboard: Monitor page |
| `docker stop <id>` | Dashboard: Containers → Stop button |
| `docker start <id>` | Dashboard: Containers → Start button |
| `docker images` | Dashboard: Images page |

## 📊 What to Verify

### Dashboard Views
- ✅ Container appears in list
- ✅ Container status is correct
- ✅ Port mappings are displayed
- ✅ Actions work (Start/Stop/Restart)
- ✅ Logs are accessible
- ✅ Stats/metrics update

### API Responses
- ✅ Container created successfully
- ✅ Container started successfully
- ✅ Logs are returned
- ✅ Stats are returned
- ✅ Health checks work

### Application
- ✅ Frontend loads at localhost:3000
- ✅ Backend responds at localhost:5000/health
- ✅ Todos can be created/updated/deleted
- ✅ Data persists in MongoDB

## 🐛 Troubleshooting

### Container Not Showing
- Refresh dashboard
- Check browser console for errors
- Verify API is responding: `curl http://localhost:8081/api/containers`

### Application Not Accessible
- Check container status (should be "running")
- Check port mappings in container details
- Check logs for errors
- Verify ports aren't already in use: `lsof -i :3000 -i :5000 -i :27017`

### Logs Not Showing
- Wait a few seconds for logs to populate
- Check API directly: `curl http://localhost:8081/api/containers/{id}/logs`

## 📸 Screenshots to Take

For documentation, capture:
1. Build Spec page with Dockerfile preview
2. Images page showing mern-mvp:latest
3. Containers page with running container
4. Container details view
5. Container logs view
6. Monitor page with metrics
7. Running Todo app at localhost:3000
8. Backend health check response

## 🎉 Success Criteria

✅ Image built from buildspec.json  
✅ Container created and visible in dashboard  
✅ Container started and running  
✅ Application accessible at mapped ports  
✅ Logs available and updating  
✅ Metrics visible in monitor  
✅ Lifecycle operations work (stop/start/restart)  
✅ API responses match expected format  

