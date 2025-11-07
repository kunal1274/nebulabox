# MVP Testing Plan - NebulaBox Platform

## Overview

This document outlines the testing plan for NebulaBox's first MVP using a simple MERN (MongoDB, Express, React, Node.js) stack application as the test workload.

## Goals

1. **Validate Core Functionality**: Test basic container lifecycle management
2. **Test Networking**: Verify container networking and service discovery
3. **Validate Multi-Container Support**: Test running multiple interconnected containers
4. **Test Resource Management**: Verify CPU, memory, and disk limits
5. **Test Monitoring**: Validate metrics collection and dashboard visibility
6. **Test Shared Runtime**: Validate workspace sharing capabilities
7. **Test Auto-Sleep**: Validate workspace auto-sleep functionality

## Test Application: Simple MERN Stack

### Application Architecture

```
┌─────────────────────────────────────────────────────┐
│                  MERN Stack App                      │
├─────────────────────────────────────────────────────┤
│                                                       │
│  ┌─────────────┐    ┌──────────────┐                │
│  │   React     │───▶│   Express    │                │
│  │  Frontend   │    │   Backend    │                │
│  │  (Port 3000)│    │  (Port 5000) │                │
│  └─────────────┘    └──────┬───────┘                │
│                             │                         │
│                             ▼                         │
│                    ┌──────────────┐                  │
│                    │   MongoDB    │                  │
│                    │  (Port 27017)│                  │
│                    └──────────────┘                  │
│                                                       │
└─────────────────────────────────────────────────────┘
```

### Application Components

1. **React Frontend** (`frontend` container)
   - Simple todo application
   - Connects to Express API
   - Port: 3000

2. **Express Backend** (`backend` container)
   - REST API for todos
   - Connects to MongoDB
   - Port: 5000

3. **MongoDB Database** (`mongodb` container)
   - Persistent data storage
   - Port: 27017

## MVP Testing Phases

### Phase 1: Basic Container Management ✅

**Objective**: Verify basic container lifecycle operations

**Test Cases**:
- [ ] Create React frontend container
- [ ] Create Express backend container
- [ ] Create MongoDB container
- [ ] Start all containers
- [ ] Verify containers are running
- [ ] Stop containers
- [ ] Delete containers

**Success Criteria**:
- All containers can be created from images
- Containers start successfully
- Containers can be stopped and deleted
- Container status updates correctly in dashboard

**Commands**:
```bash
# Pull images
nebula-cli images pull node:18-alpine
nebula-cli images pull mongo:latest

# Create containers
nebula-cli containers create --name frontend --image node:18-alpine
nebula-cli containers create --name backend --image node:18-alpine
nebula-cli containers create --name mongodb --image mongo:latest

# Start containers
nebula-cli containers start frontend
nebula-cli containers start backend
nebula-cli containers start mongodb
```

---

### Phase 2: Container Networking 🌐

**Objective**: Test container networking and communication

**Test Cases**:
- [ ] Create a custom network for MERN stack
- [ ] Connect all containers to the network
- [ ] Verify containers can resolve each other by name
- [ ] Test service discovery registration
- [ ] Verify DNS resolution

**Success Criteria**:
- Containers can communicate on the same network
- Service names resolve to container IPs
- Service discovery shows registered services
- DNS resolution works for service names

**Commands**:
```bash
# Create network
nebula-cli networks create --name mern-network --driver bridge

# Connect containers to network (via dashboard or CLI)
# Register services
nebula-cli services register --name frontend --id frontend-container-id
nebula-cli services register --name backend --id backend-container-id
nebula-cli services register --name mongodb --id mongodb-container-id
```

---

### Phase 3: Port Mapping & Access 🔌

**Objective**: Test port mapping and external access

**Test Cases**:
- [ ] Map frontend port 3000 to host
- [ ] Map backend port 5000 to host
- [ ] Verify ports are allocated correctly
- [ ] Test HTTP access to frontend
- [ ] Test API access to backend
- [ ] Verify port conflict detection

**Success Criteria**:
- Ports are correctly mapped
- Services are accessible from host
- Port allocation avoids conflicts
- Port registry shows allocated ports

**Commands**:
```bash
# Create containers with port mappings
nebula-cli containers create --name frontend --image node:18-alpine \
  --port 3000:3000

nebula-cli containers create --name backend --image node:18-alpine \
  --port 5000:5000

# Check port allocation
nebula-cli ports list
```

---

### Phase 4: Environment Variables & Configuration ⚙️

**Objective**: Test environment variable management

**Test Cases**:
- [ ] Set environment variables for backend (MongoDB connection)
- [ ] Set environment variables for frontend (API URL)
- [ ] Verify variables are applied correctly
- [ ] Test environment variable updates

**Success Criteria**:
- Environment variables are set correctly
- Containers can access environment variables
- Variables persist across container restarts
- Dashboard shows environment variables

**Commands**:
```bash
# Set environment variables
nebula-cli containers env set backend \
  MONGO_HOST=mongodb \
  MONGO_PORT=27017 \
  MONGO_DB=todos

nebula-cli containers env set frontend \
  REACT_APP_API_URL=http://backend:5000
```

---

### Phase 5: Volume Mounting & Persistence 💾

**Objective**: Test volume mounting and data persistence

**Test Cases**:
- [ ] Create volume for MongoDB data
- [ ] Mount volume to MongoDB container
- [ ] Create volume for backend logs
- [ ] Mount volume to backend container
- [ ] Verify data persists after container restart
- [ ] Test volume listing and inspection

**Success Criteria**:
- Volumes are created successfully
- Volumes mount correctly to containers
- Data persists across container restarts
- Dashboard shows mounted volumes

**Commands**:
```bash
# Create volumes
nebula-cli volumes create mongo-data
nebula-cli volumes create backend-logs

# Mount volumes (via dashboard or container update)
```

---

### Phase 6: Health Checks & Monitoring 📊

**Objective**: Test health checks and monitoring

**Test Cases**:
- [ ] Configure HTTP health check for frontend
- [ ] Configure HTTP health check for backend
- [ ] Configure TCP health check for MongoDB
- [ ] Verify health check status in dashboard
- [ ] Test metrics collection
- [ ] Verify CPU/RAM usage monitoring
- [ ] Test alert system (if implemented)

**Success Criteria**:
- Health checks execute correctly
- Health status updates in real-time
- Metrics are collected and displayed
- Dashboard shows accurate resource usage

**Commands**:
```bash
# Set health checks (via dashboard or API)
# Frontend: HTTP GET http://localhost:3000/health
# Backend: HTTP GET http://localhost:5000/health
# MongoDB: TCP 27017
```

---

### Phase 7: Container Groups & Hierarchy 📦

**Objective**: Test container grouping functionality

**Test Cases**:
- [ ] Create container group "mern-stack"
- [ ] Add frontend, backend, mongodb to group
- [ ] Verify group hierarchy
- [ ] Start all containers in group
- [ ] Stop all containers in group
- [ ] Test group resource limits

**Success Criteria**:
- Container groups can be created
- Containers can be added to groups
- Group operations affect all containers
- Dashboard shows group hierarchy

**Commands**:
```bash
# Create group
nebula-cli groups create --name mern-stack

# Add containers
nebula-cli groups add --group mern-stack --container frontend
nebula-cli groups add --group mern-stack --container backend
nebula-cli groups add --group mern-stack --container mongodb

# Start group
nebula-cli groups start mern-stack
```

---

### Phase 8: Stack Template Deployment 🚀

**Objective**: Test deploying MERN stack from template

**Test Cases**:
- [ ] Verify MERN template exists (or create it)
- [ ] Deploy MERN stack from template
- [ ] Verify all containers are created
- [ ] Verify networking is configured
- [ ] Verify environment variables are set
- [ ] Test template customization

**Success Criteria**:
- Template deployment works
- All containers are created with correct config
- Networking is automatically configured
- Environment variables are applied

**Commands**:
```bash
# Deploy from template
nebula-cli templates deploy --template mern-stack --name my-mern-app
```

---

### Phase 9: Shared Runtime & Collaboration 👥

**Objective**: Test workspace sharing and collaboration

**Test Cases**:
- [ ] Create shared workspace for MERN stack
- [ ] Invite team member (viewer role)
- [ ] Invite team member (editor role)
- [ ] Test viewer can see but not modify
- [ ] Test editor can modify containers
- [ ] Test session sharing
- [ ] Test audit logging

**Success Criteria**:
- Workspaces can be created
- Invites are sent and accepted
- Role-based permissions work correctly
- Sessions can be shared
- Audit logs are created

**Commands**:
```bash
# Create workspace
nebula-cli shareruntime workspace create \
  --name "MERN Development" \
  --container frontend

# Share workspace
nebula-cli shareruntime share --workspace ws-123
# Get invite link

# Join workspace (as another user)
nebula-cli shareruntime join --token <invite-token>
```

---

### Phase 10: Auto-Sleep & Snapshots 😴

**Objective**: Test workspace auto-sleep and snapshots

**Test Cases**:
- [ ] Configure auto-sleep for workspace
- [ ] Verify workspace goes to sleep after idle
- [ ] Verify snapshot is created on sleep
- [ ] Wake workspace manually
- [ ] Verify workspace restores from snapshot
- [ ] Test auto-wake on access

**Success Criteria**:
- Auto-sleep triggers after idle timeout
- Snapshots are created automatically
- Workspaces can be woken up
- State is restored correctly

**Commands**:
```bash
# Configure auto-sleep
nebula-cli shareruntime autosleep set \
  --workspace ws-123 \
  --enabled true \
  --idle-timeout 30 \
  --create-snapshot true

# Wake workspace
nebula-cli shareruntime wake --workspace ws-123
```

---

### Phase 11: Performance & Load Testing 📈

**Objective**: Test system performance under load

**Test Cases**:
- [ ] Run all 3 MERN containers simultaneously
- [ ] Generate load on backend API
- [ ] Monitor resource usage
- [ ] Test resource limits enforcement
- [ ] Test concurrent operations
- [ ] Verify system stability

**Success Criteria**:
- System handles concurrent operations
- Resource limits are enforced
- No resource leaks
- Performance metrics are accurate

**Tools**:
```bash
# Use Apache Bench or similar
ab -n 1000 -c 10 http://localhost:5000/api/todos

# Monitor metrics via dashboard
```

---

### Phase 12: Failure Scenarios & Recovery 🔄

**Objective**: Test failure handling and recovery

**Test Cases**:
- [ ] Kill frontend container (simulate failure)
- [ ] Verify health check detects failure
- [ ] Test auto-restart (if configured)
- [ ] Test manual restart
- [ ] Test container recovery from snapshot
- [ ] Test network partition scenarios

**Success Criteria**:
- Failures are detected
- Containers can be recovered
- Data persists through failures
- System remains stable

---

## Test Application Setup

### Step 1: Prepare MERN Application Code

Create a simple todo application structure:

```bash
mkdir mern-test-app
cd mern-test-app

# Frontend structure
mkdir frontend
# React app with todo list UI

# Backend structure
mkdir backend
# Express API with CRUD operations

# Dockerfiles
# Create Dockerfile for frontend and backend
```

### Step 2: Build Application Images

```bash
# Build frontend image
cd frontend
docker build -t mern-frontend:latest .

# Build backend image
cd ../backend
docker build -t mern-backend:latest .
```

### Step 3: Push Images to Registry (or use local)

```bash
# Option 1: Use Nebula Registry
nebula-cli registry push mern-frontend:latest
nebula-cli registry push mern-backend:latest

# Option 2: Use Docker Hub
docker push username/mern-frontend:latest
docker push username/mern-backend:latest
```

## Testing Checklist

### Pre-Testing Setup
- [ ] NebulaBox API server is running
- [ ] NebulaBox dashboard is accessible
- [ ] Containerd is running (or mock mode)
- [ ] Network connectivity verified
- [ ] Test images are available

### Testing Execution
- [ ] Execute Phase 1 tests
- [ ] Execute Phase 2 tests
- [ ] Execute Phase 3 tests
- [ ] Execute Phase 4 tests
- [ ] Execute Phase 5 tests
- [ ] Execute Phase 6 tests
- [ ] Execute Phase 7 tests
- [ ] Execute Phase 8 tests
- [ ] Execute Phase 9 tests
- [ ] Execute Phase 10 tests
- [ ] Execute Phase 11 tests
- [ ] Execute Phase 12 tests

### Post-Testing
- [ ] Collect test results
- [ ] Document issues found
- [ ] Create bug reports
- [ ] Performance analysis
- [ ] Generate test report

## Success Metrics

### Functional Metrics
- ✅ All 12 testing phases pass
- ✅ Zero critical bugs
- ✅ < 5 minor bugs
- ✅ All core features working

### Performance Metrics
- Container creation: < 2 seconds
- Container start: < 1 second
- Network creation: < 500ms
- Dashboard load: < 2 seconds
- API response time: < 100ms

### Reliability Metrics
- System uptime: > 99%
- No memory leaks
- No container crashes
- Successful recovery from failures

## Issues Tracking

### Known Issues
- List any known issues here

### Testing Blockers
- List any blockers here

### Suggestions for Improvement
- List suggestions here

## Test Reports

Test reports will be generated after each phase:
- `test-reports/phase-1-basic-containers.md`
- `test-reports/phase-2-networking.md`
- `test-reports/phase-3-ports.md`
- ... (one per phase)

## Next Steps

1. **Prepare Test Environment**
   - Set up NebulaBox in test mode
   - Prepare MERN application code
   - Create test scripts

2. **Execute Phase 1-3**
   - Basic functionality testing
   - Network testing
   - Port mapping testing

3. **Execute Phase 4-6**
   - Configuration testing
   - Persistence testing
   - Monitoring testing

4. **Execute Phase 7-9**
   - Advanced features testing
   - Collaboration testing

5. **Execute Phase 10-12**
   - Auto-sleep testing
   - Performance testing
   - Failure testing

6. **Generate Final Report**
   - Compile all test results
   - Document findings
   - Create improvement plan

## Resources

### Test Data
- Sample todos for testing
- Test users for collaboration
- Test networks and volumes

### Test Scripts
- Automation scripts for repetitive tests
- Load generation scripts
- Failure injection scripts

### Documentation
- API reference
- Dashboard user guide
- Troubleshooting guide

