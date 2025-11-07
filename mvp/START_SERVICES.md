# Starting NebulaBox Services for MVP Testing

## Quick Start

```bash
# From project root
./scripts/start-nebulabox.sh
```

This will:
1. ✅ Start API server on port 8081
2. ✅ Start Dashboard on port 3000
3. ✅ Show access URLs
4. ✅ Display log locations

## Manual Start

### Start API Server

```bash
# Build if needed
make build-api

# Start API server
./nebulabox-api

# Or in background
./nebulabox-api > /tmp/nebulabox-api.log 2>&1 &
```

### Start Dashboard

```bash
cd web/dashboard

# Install dependencies (first time only)
npm install

# Start dev server
npm run dev

# Or in background
npm run dev > /tmp/nebulabox-dashboard.log 2>&1 &
```

## Access URLs

After starting:

- **Dashboard UI**: http://localhost:3000
- **API Server**: http://localhost:8081
- **API Endpoint**: http://localhost:8081/api

## Verify Services

```bash
# Check API
curl http://localhost:8081/api/containers

# Check Dashboard
curl http://localhost:3000

# Check ports
lsof -i :8081
lsof -i :3000
```

## View Logs

```bash
# API logs
tail -f /tmp/nebulabox-api.log

# Dashboard logs
tail -f /tmp/nebulabox-dashboard.log
```

## Stop Services

```bash
# Use stop script
./scripts/stop-nebulabox.sh

# Or manually
pkill -f nebulabox-api
pkill -f "vite.*dashboard"
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :8081
lsof -i :3000

# Kill process
kill -9 <PID>
```

### API Not Responding

```bash
# Check API logs
tail -20 /tmp/nebulabox-api.log

# Check if API is running
ps aux | grep nebulabox-api
```

### Dashboard Not Loading

```bash
# Check dashboard logs
tail -20 /tmp/nebulabox-dashboard.log

# Verify API is running (dashboard needs API)
curl http://localhost:8081/api/containers
```

## Next Steps

Once both services are running:
1. Open Dashboard: http://localhost:3000
2. Verify you can see containers, images, etc.
3. Proceed with MVP testing!

