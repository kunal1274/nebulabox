# Dashboard Visual Walkthrough - Real-time Testing

Follow these steps to visually verify NebulaBox functionality in the dashboard.

## 🎯 Complete Flow Visualization

### Step 1: Build Spec Page - Build the Image

**URL:** http://localhost:3001/buildspec

1. **Open Build Spec page** in dashboard
2. **Click "Load Example"** or **paste buildspec.json** content
3. **Click "Validate"** button
   - ✅ Should show "Valid" in Validation tab
   - ✅ Dockerfile tab should show generated Dockerfile
   - ✅ Should see `FROM node:18-alpine` and all steps
4. **Click "Build"** button
   - ✅ Should show build logs in "Build Logs" tab
   - ✅ Should see: "Building from NebulaBox build spec"
   - ✅ Should see: "Successfully built from NebulaBox build specification"
   - ✅ Message: "Build completed and registered successfully"

**What to Look For:**
- Dockerfile preview showing all build steps
- Validation showing "Valid" status
- Build logs with step-by-step progress
- No error messages

---

### Step 2: Images Page - Verify Image Created

**URL:** http://localhost:3001/images

1. **Navigate to Images page**
2. **Look for `mern-mvp:latest`** in the list
3. **Verify image details:**
   - ✅ Image name: `mern-mvp`
   - ✅ Tag: `latest`
   - ✅ Size displayed
   - ✅ Created timestamp

**What to Look For:**
- Image appears in list
- Image details are correct
- Can see when it was created

---

### Step 3: Containers Page - Create Container

**URL:** http://localhost:3001/containers

1. **Click "Create Container"** button (top right)
2. **Fill in the form:**
   - **Image:** `mern-mvp:latest`
   - **Name:** `mern-mvp-test`
   - **Ports:** 
     - Container Port: `3000`, Host Port: `3000` → Add
     - Container Port: `5000`, Host Port: `5000` → Add
     - Container Port: `27017`, Host Port: `27017` → Add
   - **Environment Variables:**
     - Key: `MONGO_HOST`, Value: `localhost` → Add
     - Key: `MONGO_PORT`, Value: `27017` → Add
     - Key: `MONGO_DB`, Value: `todos` → Add
     - Key: `PORT`, Value: `5000` → Add
3. **Click "Create"** button

**What to Look For:**
- Container appears in container list
- Status shows as "created" or "running"
- Port mappings displayed correctly
- Environment variables listed

---

### Step 4: Container Details - View Configuration

1. **Click on container name** `mern-mvp-test`
2. **View container details:**
   - ✅ Container ID
   - ✅ Image: `mern-mvp:latest`
   - ✅ Status: `running` (green badge)
   - ✅ Port mappings: 3000:3000, 5000:5000, 27017:27017
   - ✅ Environment variables listed
   - ✅ Created timestamp
   - ✅ Resource usage (if available)

**What to Look For:**
- All details are correct
- Status is accurate
- Ports are mapped correctly
- Env vars are shown

---

### Step 5: Container Logs - View Startup

1. **Click "Logs" button** on container card
2. **View container logs:**
   - ✅ MongoDB startup messages
   - ✅ Backend API starting
   - ✅ Frontend building/serving
   - ✅ Health check status
   - ✅ No error messages

**What to Look For:**
- Logs are streaming or available
- Services are starting successfully
- No critical errors
- Health checks passing

---

### Step 6: Monitor Page - View Metrics

**URL:** http://localhost:3001/monitor

1. **Navigate to Monitor page**
2. **Check System Stats:**
   - ✅ CPU Usage percentage
   - ✅ Memory Usage percentage
   - ✅ Disk Usage percentage
   - ✅ Containers Running count
   - ✅ Total Containers count
3. **Check Container Metrics:**
   - ✅ Find `mern-mvp-test` in list
   - ✅ View CPU usage
   - ✅ View Memory usage
   - ✅ View Network I/O
   - ✅ Real-time updates (if available)

**What to Look For:**
- Metrics are updating
- Container appears in metrics list
- Resource usage is reasonable
- Graphs/charts are displaying

---

### Step 7: Test Running Application

#### Frontend (Todo App)
**URL:** http://localhost:3000

1. **Open in browser:** http://localhost:3000
2. **What to see:**
   - ✅ Todo list interface
   - ✅ "Add a new todo..." input field
   - ✅ Connection status (should show "Connected")
   - ✅ Empty list or existing todos
3. **Test functionality:**
   - ✅ Add a todo item
   - ✅ Mark todo as complete
   - ✅ Edit a todo
   - ✅ Delete a todo
   - ✅ Todos persist (refresh page)

**What to Look For:**
- Page loads without errors
- UI is responsive
- Todos can be created/edited/deleted
- Data persists on refresh

#### Backend API Health
**URL:** http://localhost:5000/health

1. **Open in browser or use curl:**
   ```bash
   curl http://localhost:5000/health
   ```
2. **Expected response:**
   ```json
   {
     "status": "healthy",
     "service": "mern-backend",
     "database": "connected",
     "timestamp": "..."
   }
   ```

**What to Look For:**
- Returns JSON response
- Status is "healthy"
- Database shows "connected"
- Response is quick (< 1 second)

#### Backend API (Todos)
**URL:** http://localhost:5000/api/todos

1. **Test API endpoints:**
   ```bash
   # Get todos
   curl http://localhost:5000/api/todos
   
   # Create todo
   curl -X POST http://localhost:5000/api/todos \
     -H "Content-Type: application/json" \
     -d '{"title": "Test Todo", "completed": false}'
   ```

**What to Look For:**
- API responds correctly
- Todos can be created via API
- Todos persist in database

---

### Step 8: Container Lifecycle Operations

#### Stop Container
1. **On Containers page**, find `mern-mvp-test`
2. **Click "Stop" button**
3. **Verify:**
   - ✅ Status changes to "stopped"
   - ✅ Badge turns red/gray
   - ✅ Application at localhost:3000 becomes inaccessible

#### Start Container
1. **Click "Start" button**
2. **Verify:**
   - ✅ Status changes to "running"
   - ✅ Badge turns green
   - ✅ Application at localhost:3000 becomes accessible again
   - ✅ Logs show startup messages

#### Restart Container
1. **Click "Restart" button**
2. **Verify:**
   - ✅ Container restarts
   - ✅ Status remains "running"
   - ✅ Logs show restart messages
   - ✅ Application continues working

**What to Look For:**
- All lifecycle actions work
- Status updates correctly
- Application responds to status changes
- Logs reflect the changes

---

### Step 9: API Response Verification

Open browser DevTools (F12) → Network tab to see API calls:

#### Create Container Request
- **Endpoint:** `POST /api/containers/run`
- **Request Body:** JSON with image, name, ports, env
- **Response:** Container object with id, name, status

#### List Containers Request
- **Endpoint:** `GET /api/containers`
- **Response:** Array of container objects

#### Container Details Request
- **Endpoint:** `GET /api/containers/{id}`
- **Response:** Full container details

#### Container Logs Request
- **Endpoint:** `GET /api/containers/{id}/logs`
- **Response:** Logs array or object

**What to Look For:**
- API calls are successful (200 status)
- Responses contain expected data
- No CORS errors
- Responses are well-formed JSON

---

## 📊 Side-by-Side: Docker vs NebulaBox

### Docker Commands
```bash
# Build
docker build -t mern-mvp:latest -f Dockerfile .

# Run
docker run -d \
  --name mern-mvp-test \
  -p 3000:3000 \
  -p 5000:5000 \
  -p 27017:27017 \
  -e MONGO_HOST=localhost \
  -e MONGO_PORT=27017 \
  mern-mvp:latest

# List
docker ps

# Logs
docker logs mern-mvp-test

# Stats
docker stats mern-mvp-test
```

### NebulaBox Dashboard
- **Build:** Build Spec page → Build button
- **Run:** Containers page → Create → Fill form → Create
- **List:** Containers page → See all containers
- **Logs:** Containers page → Click Logs button
- **Stats:** Monitor page → View metrics

**Same functionality, visual interface!**

---

## ✅ Success Checklist

- [ ] Image built successfully from buildspec
- [ ] Image appears in Images page
- [ ] Container created successfully
- [ ] Container visible in Containers page
- [ ] Container status is "running"
- [ ] Port mappings are correct
- [ ] Environment variables are set
- [ ] Container logs are accessible
- [ ] Frontend loads at localhost:3000
- [ ] Backend responds at localhost:5000/health
- [ ] Todo app functionality works
- [ ] Stop/Start/Restart operations work
- [ ] Metrics visible in Monitor page
- [ ] API responses are correct

---

## 🐛 Common Issues & Solutions

### Container Not Appearing
- **Check:** Refresh dashboard (F5)
- **Check:** Browser console for errors
- **Check:** API is responding: `curl http://localhost:8081/api/containers`

### Application Not Accessible
- **Check:** Container status is "running"
- **Check:** Ports are correctly mapped
- **Check:** No port conflicts: `lsof -i :3000 -i :5000`
- **Check:** Container logs for errors

### Logs Not Showing
- **Check:** Wait a few seconds for logs to populate
- **Check:** Container is running
- **Check:** API directly: `curl http://localhost:8081/api/containers/{id}/logs`

---

## 📸 Screenshots to Capture

For documentation:
1. Build Spec page with Dockerfile preview
2. Build logs showing successful build
3. Images page with mern-mvp:latest
4. Containers page with mern-mvp-test running
5. Container details view
6. Container logs view
7. Monitor page with metrics
8. Running Todo app
9. Browser DevTools showing API calls

