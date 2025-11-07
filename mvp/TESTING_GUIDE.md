# MVP Testing Guide - NebulaBox as Docker Alternative

This guide walks you through testing NebulaBox's Docker-like functionality with the MERN stack MVP.

## Prerequisites

### 1. Start Services

**Terminal 1 - API Server:**
```bash
cd /home/serverratxen/Documents/cursor-projects/nebulabox
./nebulabox-api
# Or: make api
```

**Terminal 2 - Dashboard:**
```bash
cd /home/serverratxen/Documents/cursor-projects/nebulabox/web/dashboard
npm run dev
```

### 2. Verify Services

```bash
# Check API
curl http://localhost:8081/api/health

# Check Dashboard
curl http://localhost:3001
```

## Step-by-Step Testing Flow

### Step 1: Build Image from Build Spec

#### Via Dashboard:
1. Open http://localhost:3001
2. Navigate to **"Build Spec"** page
3. Click **"Load Example"** or paste your `buildspec.json`
4. Click **"Validate"** - should show Dockerfile preview
5. Click **"Build"** - should show build logs
6. Image `mern-mvp:latest` should be created

#### Via API (for verification):
```bash
# Validate buildspec
curl -X POST http://localhost:8081/api/buildspec/validate \
  -H "Content-Type: application/json" \
  -d '{
    "spec": {
      "version": "1.0",
      "name": "mern-mvp",
      "tag": "mern-mvp:latest",
      "base": {"image": "node", "tag": "18-alpine"},
      "steps": [...]
    }
  }'

# Build image
curl -X POST http://localhost:8081/api/buildspec/build \
  -H "Content-Type: application/json" \
  -d '{
    "spec": {...},
    "tag": "mern-mvp:latest"
  }'
```

### Step 2: Verify Image in Dashboard

1. Navigate to **"Images"** page
2. Should see `mern-mvp:latest` in the list
3. Check image details (size, created date, etc.)

#### Via API:
```bash
curl http://localhost:8081/api/images
```

### Step 3: Create Container

#### Via Dashboard:
1. Navigate to **"Containers"** page
2. Click **"Create Container"** button
3. Fill in form:
   - **Image**: `mern-mvp:latest`
   - **Name**: `mern-mvp-test`
   - **Ports**: 
     - `3000:3000` (Frontend)
     - `5000:5000` (Backend API)
     - `27017:27017` (MongoDB)
   - **Environment Variables** (optional):
     - `MONGO_HOST=localhost`
     - `MONGO_PORT=27017`
     - `MONGO_DB=todos`
     - `PORT=5000`
4. Click **"Create"**
5. Container should appear in the list

#### Via API:
```bash
curl -X POST http://localhost:8081/api/containers/run \
  -H "Content-Type: application/json" \
  -d '{
    "image": "mern-mvp:latest",
    "name": "mern-mvp-test",
    "ports": ["3000:3000", "5000:5000", "27017:27017"],
    "env": [
      "MONGO_HOST=localhost",
      "MONGO_PORT=27017",
      "MONGO_DB=todos",
      "PORT=5000"
    ],
    "detach": true
  }'
```

### Step 4: Start Container

1. In **"Containers"** page, find your container
2. Click **"Start"** button
3. Status should change to "running"
4. Check container logs

#### Via API:
```bash
# Start container
curl -X POST http://localhost:8081/api/containers/{container-id}/start

# Get container logs
curl http://localhost:8081/api/containers/{container-id}/logs
```

### Step 5: Verify Container is Running

#### Check Status:
1. Container status should be **"running"**
2. Port mappings should be active
3. Health check should pass (if configured)

#### Access Services:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:5000/health
- **MongoDB**: localhost:27017

#### Via API:
```bash
# Get container details
curl http://localhost:8081/api/containers/{container-id}

# Get container stats
curl http://localhost:8081/api/containers/{container-id}/stats
```

### Step 6: Test MERN Application

1. Open http://localhost:3000 in browser
2. Should see the Todo App
3. Add a todo item
4. Verify it persists (backend + MongoDB working)
5. Check backend health: http://localhost:5000/health

### Step 7: Monitor Container

1. Navigate to **"Monitor"** page in dashboard
2. View container metrics:
   - CPU usage
   - Memory usage
   - Network I/O
   - Disk I/O

#### Via API:
```bash
# Get system stats
curl http://localhost:8081/api/system/stats

# Get container stats
curl http://localhost:8081/api/containers/{container-id}/stats
```

### Step 8: Test Container Lifecycle

1. **Stop** container - Click "Stop" button
2. **Start** container - Click "Start" button
3. **Restart** container - Click "Restart" button
4. **View Logs** - Check logs in dashboard
5. **Delete** container - Click "Delete" button (optional)

#### Via API:
```bash
# Stop
curl -X POST http://localhost:8081/api/containers/{container-id}/stop

# Start
curl -X POST http://localhost:8081/api/containers/{container-id}/start

# Restart
curl -X POST http://localhost:8081/api/containers/{container-id}/restart

# Delete
curl -X DELETE http://localhost:8081/api/containers/{container-id}
```

## API Response Examples

### Image Build Response
```json
{
  "valid": true,
  "dockerfile": "FROM node:18-alpine\n...",
  "tag": "mern-mvp:latest",
  "logs": [
    "[+] Building from NebulaBox build spec",
    " => exporting to image",
    "Successfully built from NebulaBox build specification"
  ],
  "message": "Build completed and registered successfully"
}
```

### Container Create Response
```json
{
  "id": "container-123",
  "name": "mern-mvp-test",
  "image": "mern-mvp:latest",
  "status": "created",
  "created": "2024-10-31T12:00:00Z",
  "ports": [
    {"container": 3000, "host": 3000},
    {"container": 5000, "host": 5000},
    {"container": 27017, "host": 27017}
  ]
}
```

### Container Stats Response
```json
{
  "cpu": {
    "usage": 45.2,
    "percent": 45.2
  },
  "memory": {
    "usage": 524288000,
    "limit": 1073741824,
    "percent": 48.8
  },
  "network": {
    "rx_bytes": 1024000,
    "tx_bytes": 2048000
  }
}
```

## Comparison: Docker vs NebulaBox

| Feature | Docker | NebulaBox |
|---------|--------|-----------|
| Build from file | `docker build` | Build Spec JSON |
| Create container | `docker run` | Dashboard/API |
| List containers | `docker ps` | Dashboard/API |
| View logs | `docker logs` | Dashboard/API |
| Monitor | `docker stats` | Dashboard/API |
| Images | `docker images` | Dashboard/API |

## Troubleshooting

### Container won't start
- Check logs in dashboard
- Verify ports aren't already in use
- Check environment variables

### Image not found
- Rebuild the image
- Check Images page for existing images

### Services not accessible
- Verify port mappings
- Check container status
- Review container logs

## Next Steps

1. Test volume mounting
2. Test environment variables
3. Test health checks
4. Test networking
5. Test monitoring and alerts

