# Operating Modes - Mock, Test, Live

NebulaBox supports three operating modes to support different use cases: development, testing, and production.

## Modes Overview

### Mock Mode (`mock`)
- **Purpose:** Development and demos
- **Data:** Returns only predefined mock data
- **Persistence:** No persistence - all data is static
- **Use Case:** UI development, demos, testing frontend without backend

**Characteristics:**
- Always returns mock containers, images, networks
- Built images are not stored
- Containers cannot be created
- All API responses are static

### Test Mode (`test`) - Default
- **Purpose:** UAT sandbox, integration testing
- **Data:** Stores built images and containers in memory
- **Persistence:** In-memory (cleared on server restart)
- **Use Case:** User acceptance testing, integration testing, MVP validation

**Characteristics:**
- Built images from `buildspec.json` are stored and appear in Images list
- Containers can be created and managed
- All data persists during server session
- Perfect for testing workflows end-to-end

### Live Mode (`live`)
- **Purpose:** Production deployment
- **Data:** Full persistence via registry and storage
- **Persistence:** Persistent storage, registry integration
- **Use Case:** Production environments, actual deployments

**Characteristics:**
- Full registry integration
- Persistent storage across restarts
- Real containerd integration (when implemented)
- Production-grade features

## Setting the Mode

### Environment Variable
```bash
export NEBULABOX_MODE=test  # or "mock" or "live"
```

### API Endpoint
```bash
# Get current mode
curl http://localhost:8081/api/mode

# Set mode
curl -X PUT http://localhost:8081/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode": "test"}'
```

### Frontend API
```typescript
// Get current mode
const mode = await api.getMode()

// Set mode
await api.setMode('test')
```

## Mode Behavior

### Images List (`GET /api/images`)

| Mode | Registry Images | Built Images | Mock Images |
|------|----------------|--------------|-------------|
| `mock` | ❌ | ❌ | ✅ |
| `test` | ✅ (if available) | ✅ | ✅ (if no images) |
| `live` | ✅ | ✅ | ❌ |

**Example:**
- **Mock mode:** Only shows nginx, postgres, node (mock data)
- **Test mode:** Shows mock images + any images built from buildspec
- **Live mode:** Shows registry images + built images (no mock data)

### Built Images Storage

When building an image via `POST /api/buildspec/build`:

1. **Mock mode:** Image is NOT stored, only build logs returned
2. **Test mode:** Image is stored in memory, appears in Images list
3. **Live mode:** Image is stored in registry + memory, persistent

### Containers

- **Mock mode:** Cannot create containers (only mock containers shown)
- **Test mode:** Can create containers, stored in memory
- **Live mode:** Can create containers, stored persistently

## Default Mode

The default mode is **`test`** to support UAT sandbox workflows.

This means:
- ✅ Built images are automatically stored
- ✅ Images appear in the dashboard immediately
- ✅ Containers can be created and tested
- ✅ Full end-to-end testing is possible

## Switching Modes

When switching from `test`/`live` to `mock`:
- Built images are cleared
- All test data is removed
- Only mock data is shown

When switching to `test`/`live`:
- Existing built images are preserved
- New images can be added
- Full functionality is available

## Use Cases

### Development (Mock Mode)
```bash
export NEBULABOX_MODE=mock
# UI development, no backend needed
# Static mock data for all endpoints
```

### Testing (Test Mode)
```bash
export NEBULABOX_MODE=test
# Build images from buildspec
# Create containers
# Test full workflows
# All data in memory
```

### Production (Live Mode)
```bash
export NEBULABOX_MODE=live
# Full registry integration
# Persistent storage
# Production deployments
```

## API Response Examples

### Get Mode
```json
{
  "mode": "test",
  "description": "Test mode: UAT sandbox - Built images and containers are stored in memory",
  "builtImages": 2,
  "totalImages": 5
}
```

### Set Mode
```json
{
  "message": "Mode changed successfully",
  "oldMode": "mock",
  "newMode": "test",
  "mode": "test",
  "builtImagesCleared": false
}
```

## Testing Workflow

1. **Start server in test mode:**
   ```bash
   NEBULABOX_MODE=test ./nebulabox-api
   ```

2. **Build image from buildspec:**
   ```bash
   curl -X POST http://localhost:8081/api/buildspec/build \
     -H "Content-Type: application/json" \
     -d @buildspec.json
   ```

3. **Verify image appears:**
   ```bash
   curl http://localhost:8081/api/images | jq '.[] | select(.name == "mern-mvp")'
   ```

4. **Create container:**
   ```bash
   curl -X POST http://localhost:8081/api/containers/run \
     -H "Content-Type: application/json" \
     -d '{"image": "mern-mvp:latest", "name": "test-container"}'
   ```

5. **All operations work end-to-end!**

## Dashboard Indicator

The dashboard should display the current mode to users:
- **Mock Mode:** Shows badge "MOCK MODE"
- **Test Mode:** Shows badge "TEST MODE" (or "UAT SANDBOX")
- **Live Mode:** Shows badge "LIVE MODE"

This helps users understand what data they're working with.

