# Build and Deploy Guide - Using NebulaBox Build Specification

This guide shows how to build and deploy the MERN MVP using NebulaBox's native Build Specification format.

## Prerequisites

- NebulaBox API server running on `http://localhost:8081`
- NebulaBox Dashboard accessible
- Application code ready in `mvp/app/single-container/app/`

## Step 1: Validate Build Specification

Before building, validate the build spec:

```bash
# Using curl
curl -X POST http://localhost:8081/api/buildspec/validate \
  -H "Content-Type: application/json" \
  -d @mvp/app/single-container/buildspec.json

# Or use dashboard
# Navigate to: Build Spec page
# Paste buildspec.json content
# Click "Validate"
```

Expected response:
```json
{
  "valid": true,
  "message": "Build specification is valid"
}
```

## Step 2: Preview Dockerfile (Optional)

See what Dockerfile will be generated:

```bash
# Using curl
curl -X POST http://localhost:8081/api/buildspec/convert \
  -H "Content-Type: application/json" \
  -d @mvp/app/single-container/buildspec.json

# Or use dashboard
# Click "Preview Dockerfile" button
```

This shows the generated Dockerfile from your build spec.

## Step 3: Build Image from Build Specification

**Build using NebulaBox API (NOT Docker CLI!):**

```bash
# Using NebulaBox API (recommended)
curl -X POST http://localhost:8081/api/buildspec/build \
  -H "Content-Type: application/json" \
  -d @mvp/app/single-container/buildspec.json

# Using dashboard (also recommended)
# 1. Go to Build Spec page
# 2. Upload or paste buildspec.json
# 3. Set build context directory to: mvp/app/single-container/
# 4. Click "Build from Spec"
# 5. Monitor build logs

# Or use the deployment script
./mvp/scripts/deploy-buildspec.sh
```

**Important**: We use NebulaBox to build, NOT Docker CLI commands!

The build process:
1. Validates build spec
2. Converts to Dockerfile
3. Builds image using Docker
4. Tags image as `mern-mvp:latest`
5. Optionally pushes to Nebula Registry

## Step 4: Verify Image

Check that image was created:

```bash
# List images
curl http://localhost:8081/api/images

# Or use dashboard
# Navigate to: Images page
# Look for "mern-mvp:latest"
```

## Step 5: Create Container

Create container from the built image:

```bash
# Using API
curl -X POST http://localhost:8081/api/containers/run \
  -H "Content-Type: application/json" \
  -d '{
    "image": "mern-mvp:latest",
    "name": "mern-mvp",
    "ports": ["3000:3000", "5000:5000", "27017:27017"],
    "detach": true,
    "env": [
      "MONGO_HOST=localhost",
      "MONGO_PORT=27017",
      "MONGO_DB=todos"
    ]
  }'

# Or use dashboard
# Navigate to: Containers page
# Click "Create Container"
# Select image: mern-mvp:latest
# Configure ports, env vars, etc.
# Click "Create"
```

## Step 6: Start Container

Start the container:

```bash
# Using API
curl -X POST http://localhost:8081/api/containers/mern-mvp/start

# Or use dashboard
# Click "Start" button on container
```

## Step 7: Verify Application

Check if application is running:

```bash
# Check container status
curl http://localhost:8081/api/containers

# Check health endpoint
curl http://localhost:5000/health

# Check frontend
curl http://localhost:3000

# Or use dashboard
# View container logs
# Check container status
# Monitor metrics
```

## Step 8: Test Application

Test the MERN application:

```bash
# Create a todo via API
curl -X POST http://localhost:5000/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Todo",
    "completed": false
  }'

# List todos
curl http://localhost:5000/api/todos

# Access frontend
# Open browser: http://localhost:3000
```

## Troubleshooting

### Build Fails

1. **Check build spec validation**
   ```bash
   curl -X POST http://localhost:8081/api/buildspec/validate \
     -H "Content-Type: application/json" \
     -d @mvp/app/single-container/buildspec.json
   ```

2. **Check build logs** (in dashboard)

3. **Verify build context** (all files in correct location)

### Container Won't Start

1. **Check container logs**
   ```bash
   curl http://localhost:8081/api/containers/mern-mvp/logs
   ```

2. **Verify ports aren't in use**
   ```bash
   curl http://localhost:8081/api/ports
   ```

3. **Check health status**
   ```bash
   curl http://localhost:8081/api/containers/mern-mvp
   ```

### Application Not Accessible

1. **Check port mappings** (should be 3000:3000, 5000:5000)

2. **Check firewall** (if running on remote server)

3. **Check container is running**
   ```bash
   curl http://localhost:8081/api/containers
   ```

## Automated Deployment Script

Use the deployment script:

```bash
cd mvp
./scripts/deploy-buildspec.sh
```

This script:
1. Validates build spec
2. Builds image
3. Creates container
4. Starts container
5. Verifies deployment

## Next Steps

After successful deployment:
1. Test basic functionality (Phase 1)
2. Test networking (Phase 2)
3. Test volumes (Phase 5)
4. Test monitoring (Phase 6)
5. Continue with MVP testing plan phases

