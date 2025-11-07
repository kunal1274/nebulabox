# Quick Start - MVP Visual Testing

## 🎯 Goal
Test NebulaBox as a Docker alternative by building and running the MERN stack MVP.

## ✅ Step 1: Build Image (Dashboard)

1. **Open Dashboard:** http://localhost:3001
2. **Go to:** Build Spec page (left sidebar)
3. **Action:**
   - Click **"Load Example"** OR paste content from `mvp/app/single-container/buildspec.json`
   - Click **"Validate"** button
   - ✅ Should see Dockerfile preview
   - Click **"Build"** button
   - ✅ Should see build logs: "Successfully built from NebulaBox build specification"

**Verify API:**
```bash
curl http://localhost:8081/api/images | jq '.[] | select(.name | contains("mern"))'
```

---

## ✅ Step 2: Create Container (Dashboard)

1. **Go to:** Containers page
2. **Click:** "Create Container" button (top right)
3. **Fill Form:**
   - **Image:** `mern-mvp:latest`
   - **Name:** `mern-mvp-test`
   - **Ports:** 
     - Add: `3000:3000`
     - Add: `5000:5000`
     - Add: `27017:27017`
   - **Environment Variables:**
     - `MONGO_HOST=localhost`
     - `MONGO_PORT=27017`
     - `MONGO_DB=todos`
     - `PORT=5000`
4. **Click:** "Create" button
5. **Verify:** Container appears in list with status

**Verify API:**
```bash
curl http://localhost:8081/api/containers | jq '.[] | {name, image, status}'
```

---

## ✅ Step 3: Start Container (Dashboard)

1. **Find container:** `mern-mvp-test` in list
2. **Click:** "Start" button
3. **Verify:** Status changes to "running" (green badge)
4. **Click:** "Logs" button
5. **Check logs:** Should see MongoDB, Backend, Frontend starting

**Verify API:**
```bash
CONTAINER_ID=$(curl -s http://localhost:8081/api/containers | jq -r '.[] | select(.name == "mern-mvp-test") | .id')
curl http://localhost:8081/api/containers/$CONTAINER_ID | jq '{name, status, ports}'
```

---

## ✅ Step 4: Test Application

### Frontend (Todo App)
**Open:** http://localhost:3000

**Test:**
- ✅ Page loads
- ✅ Add a todo item
- ✅ Mark as complete
- ✅ Edit/Delete todos

### Backend Health
**Open:** http://localhost:5000/health

**Expected:**
```json
{
  "status": "healthy",
  "service": "mern-backend",
  "database": "connected"
}
```

### Backend API
**Test:**
```bash
# Get todos
curl http://localhost:5000/api/todos

# Create todo
curl -X POST http://localhost:5000/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Test from API", "completed": false}'
```

---

## ✅ Step 5: Monitor (Dashboard)

1. **Go to:** Monitor page
2. **Check:**
   - System stats (CPU, Memory, Disk)
   - Container count
   - Find `mern-mvp-test` in container metrics

---

## ✅ Step 6: Verify Docker-like Operations

### Stop Container
- Click "Stop" → Status = "stopped" → App inaccessible

### Start Container  
- Click "Start" → Status = "running" → App accessible

### Restart Container
- Click "Restart" → Container restarts → Still running

### View Logs
- Click "Logs" → See all container output

### View Stats
- Monitor page → See resource usage

---

## 📊 What You're Testing

✅ **Build:** Like `docker build` - but from JSON spec  
✅ **Run:** Like `docker run` - but via dashboard form  
✅ **List:** Like `docker ps` - visual container list  
✅ **Logs:** Like `docker logs` - in dashboard  
✅ **Stats:** Like `docker stats` - real-time metrics  
✅ **Lifecycle:** Stop/Start/Restart - all in dashboard  

---

## 🎥 Watch the Flow

1. Build Spec page → Build image
2. Images page → See built image
3. Containers page → Create container
4. Containers page → Start container
5. Containers page → View logs
6. Monitor page → View metrics
7. Browser → Test application

**Everything visual, no CLI needed!**

---

## 🐛 If Something Doesn't Work

### Image Not Showing
- Refresh Images page
- Rebuild from Build Spec page

### Container Not Creating
- Check image exists
- Check form filled correctly
- Check browser console (F12)

### App Not Accessible
- Check container status = "running"
- Check ports mapped correctly
- Check container logs for errors

---

## 📸 What to Screenshot

1. Build Spec with Dockerfile preview
2. Build logs showing success
3. Images page with mern-mvp:latest
4. Containers page with running container
5. Container logs view
6. Monitor page with metrics
7. Running Todo app
8. Backend health check

---

## 🎉 Success = You Can Use Dashboard Instead of Docker CLI!

