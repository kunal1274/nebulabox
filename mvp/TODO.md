# MVP Testing TODO List

## Phase 1: Application Development (Tasks 1-4)

### ✅ Preparation Complete
- [x] MVP folder structure created
- [x] Build specification created (`buildspec.json`)
- [x] NebulaBox services started (API + Dashboard)
- [x] Testing approach defined

### 📝 Application Code (In Progress)

- [ ] **mvp-1**: Create MERN application structure
  - [ ] Frontend folder (`app/frontend/`)
  - [ ] Backend folder (`app/backend/`)
  - [ ] Scripts folder (`app/scripts/`)
  - [ ] Package.json files

- [ ] **mvp-2**: Create React frontend
  - [ ] Todo list UI component
  - [ ] API integration (connect to backend)
  - [ ] Environment variable for API URL
  - [ ] Health check endpoint

- [ ] **mvp-3**: Create Express backend
  - [ ] REST API endpoints (GET, POST, PUT, DELETE /api/todos)
  - [ ] MongoDB connection setup
  - [ ] CORS configuration
  - [ ] Health check endpoint (`/health`)

- [ ] **mvp-4**: Create start script
  - [ ] MongoDB startup script
  - [ ] Backend startup (wait for MongoDB)
  - [ ] Frontend startup (serve built files)
  - [ ] Process management

---

## Phase 2: Build & Deploy (Tasks 5-8)

### 🔨 Build Specification

- [ ] **mvp-5**: Test buildspec.json validation
  - [ ] Validate in NebulaBox Dashboard
  - [ ] Verify no errors
  - [ ] Preview generated Dockerfile

- [ ] **mvp-6**: Build image from buildspec.json
  - [ ] Build via NebulaBox API or Dashboard
  - [ ] Verify image created
  - [ ] Check image size and layers

- [ ] **mvp-7**: Create container from built image
  - [ ] Configure ports (3000, 5000, 27017)
  - [ ] Set environment variables
  - [ ] Create container via NebulaBox

- [ ] **mvp-8**: Start container and verify
  - [ ] Start container
  - [ ] Verify MongoDB running
  - [ ] Verify Backend API responding
  - [ ] Verify Frontend accessible
  - [ ] Test full MERN stack functionality

---

## Phase 3: Basic Features Testing (Tasks 9-15)

### 🌐 Networking & Configuration

- [ ] **mvp-9**: Test container networking
  - [ ] Verify port mappings (3000, 5000, 27017)
  - [ ] Test external access to services
  - [ ] Verify port allocation

- [ ] **mvp-10**: Test environment variables
  - [ ] Set env vars via NebulaBox dashboard
  - [ ] Verify env vars in container
  - [ ] Test env var updates

- [ ] **mvp-11**: Test volume mounting
  - [ ] Create volume for MongoDB data
  - [ ] Mount volume to container
  - [ ] Verify data persistence after restart

- [ ] **mvp-12**: Configure health checks
  - [ ] Set HTTP health check for backend
  - [ ] Verify health status in dashboard
  - [ ] Test health check intervals

### 📊 Monitoring

- [ ] **mvp-13**: Test monitoring and metrics
  - [ ] View CPU usage in dashboard
  - [ ] View memory usage
  - [ ] View container metrics
  - [ ] Test real-time updates

- [ ] **mvp-14**: Test container lifecycle
  - [ ] Stop container
  - [ ] Start container
  - [ ] Restart container
  - [ ] Verify state persistence

- [ ] **mvp-15**: Test logs viewing
  - [ ] View container logs in dashboard
  - [ ] Test log search/filter
  - [ ] Verify log streaming

---

## Phase 4: Advanced Features (Tasks 16-20)

### 🔗 Networking & Groups

- [ ] **mvp-16**: Create custom network
  - [ ] Create network via NebulaBox
  - [ ] Connect container to network
  - [ ] Test network isolation

- [ ] **mvp-17**: Test service discovery
  - [ ] Register services
  - [ ] Resolve service names
  - [ ] Test DNS resolution

- [ ] **mvp-18**: Create container group
  - [ ] Create group for MERN stack
  - [ ] Add container to group
  - [ ] Configure shared resources

- [ ] **mvp-19**: Test container group operations
  - [ ] Start entire group
  - [ ] Stop entire group
  - [ ] View group hierarchy

- [ ] **mvp-20**: Test stack templates
  - [ ] Create MERN stack template
  - [ ] Deploy from template
  - [ ] Verify template deployment

---

## Phase 5: Shared Runtime Testing (Tasks 21-25)

### 👥 Collaboration Features

- [ ] **mvp-21**: Create Shared Runtime workspace
  - [ ] Create workspace for MERN container
  - [ ] Verify workspace creation
  - [ ] Configure workspace settings

- [ ] **mvp-22**: Test team collaboration
  - [ ] Invite team member
  - [ ] Test viewer role (read-only)
  - [ ] Test editor role (can modify)
  - [ ] Test permissions enforcement

- [ ] **mvp-23**: Test auto-sleep
  - [ ] Configure auto-sleep settings
  - [ ] Verify workspace goes to sleep on idle
  - [ ] Verify snapshot created
  - [ ] Test wake from sleep

- [ ] **mvp-24**: Create snapshot
  - [ ] Create snapshot of running container
  - [ ] Verify snapshot created
  - [ ] View snapshot details

- [ ] **mvp-25**: Test snapshot restore
  - [ ] Restore from snapshot
  - [ ] Verify container state restored
  - [ ] Verify data persistence

---

## Phase 6: Stress Testing & Validation (Tasks 26-30)

### ⚡ Performance & Reliability

- [ ] **mvp-26**: Load testing
  - [ ] Generate traffic to MERN API
  - [ ] Monitor resource usage
  - [ ] Test concurrent requests
  - [ ] Analyze performance metrics

- [ ] **mvp-27**: Test failure scenarios
  - [ ] Simulate container crash
  - [ ] Test auto-restart (if configured)
  - [ ] Test manual recovery
  - [ ] Verify data integrity

- [ ] **mvp-28**: Document test results
  - [ ] Record findings for each phase
  - [ ] Document issues found
  - [ ] Note performance metrics
  - [ ] Capture screenshots

- [ ] **mvp-29**: Create MVP test report
  - [ ] Compile all test results
  - [ ] Include screenshots
  - [ ] Performance benchmarks
  - [ ] Recommendations for improvements

- [ ] **mvp-30**: Final review
  - [ ] Verify all features working
  - [ ] End-to-end workflow test
  - [ ] Validate NebulaBox can replace Docker
  - [ ] Prepare summary report

---

## Progress Tracking

### Overall Progress: 0/30 tasks completed

**Current Phase**: Phase 1 - Application Development

**Next Task**: mvp-1 - Create MERN application structure

---

## Notes

- All testing done via NebulaBox (not Docker CLI)
- Build using NebulaBox Build Specification
- Deploy via NebulaBox Dashboard/API
- Monitor via NebulaBox Dashboard

---

## Quick Commands

```bash
# Start services
./scripts/start-nebulabox.sh

# Access dashboard
http://localhost:3001

# Test API
curl http://localhost:8081/api/containers
```

