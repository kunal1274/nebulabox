# NebulaBox Services Status

## Quick Access

### 🌐 Dashboard UI
**URL**: http://localhost:3001 (or http://localhost:3000)

The dashboard should be accessible at one of these URLs. Open in your browser to see the NebulaBox UI.

### 📡 API Server
**URL**: http://localhost:8081/api

Test API:
```bash
curl http://localhost:8081/api/containers
```

## What You'll See in Dashboard

Once the dashboard loads, you'll see:

1. **Sidebar Navigation** with sections:
   - Dashboard (overview)
   - Containers
   - Images
   - Registry
   - Build Spec ⭐ (for our MVP!)
   - Security
   - Orchestrator
   - Runtime
   - AI Ops
   - Groups
   - Composition
   - Templates
   - Shared Runtime
   - Snapshots
   - Ephemeral Runtime
   - Monitor
   - Networks
   - Services
   - Teams
   - Tenants
   - Settings

2. **Dashboard Page** showing:
   - Running containers count
   - CPU usage
   - Memory usage
   - Disk usage
   - Quick links

## For MVP Testing

When we build and deploy the MERN app, you'll see:

1. **Build Spec Page**: Validate and build from `buildspec.json`
2. **Images Page**: See the built `mern-mvp:latest` image
3. **Containers Page**: See the running `mern-mvp` container
4. **Monitor Page**: See resource usage, logs, metrics

## Check Status

```bash
# Check if API is running
curl http://localhost:8081/api/containers

# Check dashboard
curl http://localhost:3001 | head -5

# View logs
tail -20 /tmp/nebulabox-api.log
tail -20 /tmp/nebulabox-dashboard.log
```

## Next Steps

1. ✅ Open Dashboard: http://localhost:3001
2. ✅ Verify you can navigate between pages
3. ✅ Check that Build Spec page works
4. ✅ Ready to create MERN app!

