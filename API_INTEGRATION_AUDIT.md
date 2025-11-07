# API Integration Audit - Comprehensive Analysis

## 🔍 Critical Issues Found

### 1. **Container Storage Not Integrated** ❌
**Problem:** 
- Containers are created via `runContainer` API
- But `listContainers` returns only hardcoded mock data
- Created containers never appear in dashboard

**Root Cause:**
- `containerd.CreateContainer` in mock mode doesn't store containers
- `containerd.ListContainers` always returns same 3 mock containers
- No in-memory storage for containers (unlike images which now have `builtImages`)

**Files Affected:**
- `internal/api/containers.go` - `listContainers` always calls `s.containerd.ListContainers(ctx)`
- `internal/containerd/client.go` - Mock `ListContainers` returns hardcoded data
- `internal/api/server.go` - No container storage map

---

### 2. **Image Storage - Partially Fixed** ✅
**Status:** Fixed for built images from buildspec
- Built images stored in `builtImages` map
- Appear in `/api/images` list
- Works in test/live mode

**Remaining Issue:**
- Images pulled via `PullImage` not stored
- Need to check if pulled images appear in list

---

### 3. **Environment Variables Parsing** ⚠️
**Problem:**
```go
for _, env := range req.Env {
    containerOpts.Environment[env] = env // Simplified for now
}
```
Should parse `KEY=VALUE` format, not use value as key.

**Files:**
- `internal/api/containers.go:213`

---

### 4. **Volume Parsing** ⚠️
**Problem:**
```go
for _, vol := range req.Volume {
    containerOpts.Volumes[vol] = vol // Simplified for now
}
```
Should parse `host:container` or `host:container:ro` format.

**Files:**
- `internal/api/containers.go:218`

---

## 📋 Integration Status by Endpoint

### ✅ Fully Integrated (Test Mode Works)
- `POST /api/buildspec/build` - Stores built images
- `GET /api/images` - Returns built images + mock data in test mode
- `GET /api/mode` - Mode management works
- `PUT /api/mode` - Mode switching works

### ❌ NOT Integrated (Always Returns Mock Data)
- `GET /api/containers` - Always returns 3 hardcoded containers
- `GET /api/containers/:id` - Only finds hardcoded containers
- `POST /api/containers/run` - Creates but doesn't store
- `POST /api/containers/:id/stop` - Doesn't update stored containers
- `POST /api/containers/:id/start` - Doesn't update stored containers

### ⚠️ Partially Integrated
- `GET /api/containers/:id/logs` - Mock logs, but container must exist first

---

## 🔧 Required Fixes

### Priority 1: Container Storage Integration

1. **Add container storage to Server struct:**
   ```go
   builtContainers map[string]*containerd.Container // id -> container
   builtContainersMu sync.Mutex
   ```

2. **Store containers in test/live mode:**
   - In `runContainer`: Store after creation
   - Update status when stopped/started

3. **Merge stored containers in listContainers:**
   - Get containers from containerd (real or mock)
   - Merge with stored containers
   - Respect operating mode

4. **Update container operations:**
   - `stopContainer`: Update stored container status
   - `startContainer`: Update stored container status
   - `deleteContainer`: Remove from storage

### Priority 2: Fix Parsing Issues

1. **Environment variables:** Parse `KEY=VALUE`
2. **Volumes:** Parse `host:container` or `host:container:ro`

### Priority 3: Image Pull Storage

1. **Store pulled images:** When `PullImage` is called, store in `builtImages`
2. **Check image existence:** Before creating container, check if image exists

---

## 🎯 Test Mode Requirements

For test mode to work properly:

1. ✅ Built images stored → `builtImages` map
2. ❌ Created containers stored → Need `builtContainers` map
3. ❌ Container status updates → Need to update stored containers
4. ❌ Container deletion → Need to remove from storage

---

## 📊 Data Flow Analysis

### Current Flow (Broken):
```
POST /api/containers/run
  → containerd.CreateContainer (creates container, returns it)
  → containerd.StartContainer (starts it)
  → Returns container response
  ❌ Container is LOST - not stored anywhere

GET /api/containers
  → containerd.ListContainers (returns hardcoded mock data)
  ❌ Created containers never appear
```

### Fixed Flow (Needed):
```
POST /api/containers/run
  → containerd.CreateContainer (creates container)
  → Store in builtContainers map (test/live mode)
  → containerd.StartContainer (starts it)
  → Update status in storage
  → Returns container response

GET /api/containers
  → Get from containerd (real or mock)
  → Get from builtContainers map (test/live mode)
  → Merge and deduplicate
  → Return combined list
```

---

## 🚀 Implementation Plan

1. **Add container storage** (similar to builtImages)
2. **Integrate container CRUD operations**
3. **Fix environment/volume parsing**
4. **Test end-to-end flow**
5. **Verify dashboard shows created containers**

