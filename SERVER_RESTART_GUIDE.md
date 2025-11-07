# NebulaBox Server Restart Guide

This guide helps you restart NebulaBox services (API and Dashboard) when needed.

## Quick Start

### Automated Restart (Recommended)

```bash
# Restart all services
./scripts/restart-nebulabox.sh

# Or with make
make restart
```

### Manual Restart

#### 1. Stop Services

```bash
# Find running processes
ps aux | grep nebulabox
ps aux | grep vite

# Kill processes (replace PID with actual process IDs)
kill <API_PID>
kill <VITE_PID>

# Or use pkill
pkill -f nebulabox-api
pkill -f "vite.*dashboard"
```

#### 2. Start Services

**Terminal 1 - API Server:**
```bash
cd /home/serverratxen/Documents/cursor-projects/nebulabox
./nebulabox-api
# Or
make api
```

**Terminal 2 - Dashboard:**
```bash
cd /home/serverratxen/Documents/cursor-projects/nebulabox/web/dashboard
npm run dev
```

## Detailed Steps

### Method 1: Using Makefile

```bash
# Build and start API
make api

# Start dashboard (in another terminal)
cd web/dashboard && npm run dev
```

### Method 2: Direct Execution

```bash
# Build API first
cd /home/serverratxen/Documents/cursor-projects/nebulabox
go build -o nebulabox-api ./cmd/api

# Run API
./nebulabox-api

# Dashboard (in another terminal)
cd web/dashboard
npm run dev
```

### Method 3: Background Processes

```bash
# Start API in background
nohup ./nebulabox-api > api.log 2>&1 &

# Start dashboard in background
cd web/dashboard
nohup npm run dev > dashboard.log 2>&1 &
```

## Verify Services

### Check API Health

```bash
curl http://localhost:8081/api/health
```

Expected response:
```json
{
  "service": "nebulabox-api",
  "status": "healthy",
  "timestamp": 1234567890,
  "version": "0.1.0-alpha"
}
```

### Check Dashboard

Open in browser:
```
http://localhost:3001
```

### Check Ports

```bash
# Check if ports are in use
lsof -i :8081  # API
lsof -i :3001  # Dashboard
```

## Common Issues

### Port Already in Use

```bash
# Find what's using the port
lsof -i :8081
lsof -i :3001

# Kill the process
kill -9 <PID>
```

### CORS Issues

After restarting API, make sure CORS is configured for `localhost:3001`:
- Check `internal/api/server.go`
- Should include `http://localhost:3001` in `AllowOrigins`

### Build Errors

If API won't build:
```bash
# Clean and rebuild
make clean
make build-api
```

## Auto-Restart Script

Create `scripts/restart-nebulabox.sh`:

```bash
#!/bin/bash

echo "🛑 Stopping NebulaBox services..."

# Kill existing processes
pkill -f nebulabox-api
pkill -f "vite.*dashboard"

sleep 2

echo "🔨 Building API..."
cd "$(dirname "$0")/.."
make build-api

echo "🚀 Starting API server..."
./nebulabox-api &
API_PID=$!

echo "⏳ Waiting for API to be ready..."
sleep 5

# Check API health
if curl -s http://localhost:8081/api/health > /dev/null; then
    echo "✅ API server is healthy (PID: $API_PID)"
else
    echo "❌ API server failed to start"
    exit 1
fi

echo ""
echo "📊 API Server: http://localhost:8081"
echo "🎨 Dashboard: Start manually with 'cd web/dashboard && npm run dev'"
echo ""
echo "💡 To stop: kill $API_PID"
```

Make it executable:
```bash
chmod +x scripts/restart-nebulabox.sh
```

## Integration with E2E Tests

E2E tests automatically start/restart services. See `e2e/playwright.config.ts` for configuration.

## Troubleshooting

### Services Won't Start

1. Check if ports are free:
   ```bash
   lsof -i :8081 -i :3001
   ```

2. Check logs:
   ```bash
   tail -f api.log
   tail -f dashboard.log
   ```

3. Rebuild:
   ```bash
   make clean && make build
   ```

### Services Start But Don't Respond

1. Check firewall:
   ```bash
   sudo ufw status
   ```

2. Check CORS configuration in `internal/api/server.go`

3. Verify environment variables:
   ```bash
   env | grep NEBULABOX
   ```

## Quick Reference

| Service | Port | Command | Health Check |
|---------|------|---------|--------------|
| API | 8081 | `./nebulabox-api` | `curl http://localhost:8081/api/health` |
| Dashboard | 3001 | `cd web/dashboard && npm run dev` | `http://localhost:3001` |
| Registry | 5001 | `./nebulabox-registry` | `curl http://localhost:5001/v2/` |

