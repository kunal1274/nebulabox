# Test Execution Worksheet - NebulaBox CLI

## Test Session Information

**Date**: _______________  
**Tester Name**: _______________  
**Environment**: 
- OS: _______________
- Go Version: _______________
- NebulaBox Version: _______________
- Test Duration: _______________

---

## Test Execution Log

### Test Case: TC-001 - CLI Version Command

**Test ID**: TC-001  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox version` | Command executes without error | | | |
| 2 | Check output contains "NebulaBox" | Output shows "🚀 NebulaBox" | | | |
| 3 | Check version number | Output shows "Version: 0.1.0-alpha" | | | |
| 4 | Check build info | Output shows "Build: Phase 1 Development" | | | |
| 5 | Check Go version | Output shows "Go Version: 1.22+" | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-002 - CLI Help Command

**Test ID**: TC-002  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox --help` | Command executes without error | | | |
| 2 | Check for "Available Commands" section | Section is displayed | | | |
| 3 | Verify "build" command listed | "build" appears in list | | | |
| 4 | Verify "images" command listed | "images" appears in list | | | |
| 5 | Verify "hierarchy" command listed | "hierarchy" appears in list | | | |
| 6 | Verify "group" command listed | "group" appears in list | | | |
| 7 | Verify "list" command listed | "list" appears in list | | | |
| 8 | Verify "ps" command listed | "ps" appears in list | | | |
| 9 | Verify "run" command listed | "run" appears in list | | | |
| 10 | Verify "stop" command listed | "stop" appears in list | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-003 - List Containers (Empty State)

**Test ID**: TC-003  
**Priority**: Medium  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox ps` | Command executes without error | | | |
| 2 | Check output | Shows "No containers found" or empty table | | | |
| 3 | Verify no error messages | No error output | | | |
| 4 | Run `./bin/nebulabox list` | Command executes without error | | | |
| 5 | Verify same output as ps | Same output as step 2 | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-004 - Pull Image

**Test ID**: TC-004  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox pull nginx:latest` | Command executes | | | |
| 2 | Check for "Pulling image" message | Output shows "⬇️ Pulling image: nginx:latest" | | | |
| 3 | Wait for completion | Command completes successfully | | | |
| 4 | Check for success message | Output shows "✅ Image pulled successfully" | | | |
| 5 | Verify no error messages | No error output | | | |
| 6 | Run `./bin/nebulabox images` | Command executes | | | |
| 7 | Verify nginx:latest in list | nginx:latest appears in image list | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-005 - List Images

**Test ID**: TC-005  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox images` | Command executes without error | | | |
| 2 | Check for table header | Shows "IMAGE ID", "REPOSITORY", "TAG", "SIZE", "CREATED" | | | |
| 3 | If images exist, verify format | Each image shows: ID (12 chars), repo, tag, size, created date | | | |
| 4 | If no images, verify message | Shows "No images found" | | | |
| 5 | Verify no error messages | No error output | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-006 - Run Container

**Test ID**: TC-006  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure image exists (pull if needed) | Image nginx:latest available | | | |
| 2 | Run `./bin/nebulabox run nginx:latest --name test-web` | Command executes | | | |
| 3 | Check for "Starting container" message | Output shows "🚀 NebulaBox: Starting container..." | | | |
| 4 | Check for "Pulling image" (if needed) | May show pull message if image not cached | | | |
| 5 | Check for "Creating container" | Output shows "📦 Creating container..." | | | |
| 6 | Check for "Starting container" | Output shows "🔄 Starting container: <id>" | | | |
| 7 | Check for success message | Output shows "✅ Container started successfully!" | | | |
| 8 | Verify container ID displayed | Container ID shown in output | | | |
| 9 | Verify container name displayed | Container name "test-web" shown | | | |
| 10 | Run `./bin/nebulabox ps` | Command executes | | | |
| 11 | Verify container in list | "test-web" appears in container list | | | |
| 12 | Verify status is "running" | Container status shows "running" | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-007 - Stop Container

**Test ID**: TC-007  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure container running (from TC-006) | Container "test-web" exists and running | | | |
| 2 | Run `./bin/nebulabox stop test-web` | Command executes | | | |
| 3 | Check for "Stopping container" message | Output shows "🛑 Stopping container: test-web" | | | |
| 4 | Check for success message | Output shows "✅ Container test-web stopped" | | | |
| 5 | Verify no error messages | No error output | | | |
| 6 | Run `./bin/nebulabox ps` | Command executes | | | |
| 7 | Verify container not in running list | Container not shown (or status "stopped") | | | |
| 8 | Run `./bin/nebulabox ps --all` | Command executes | | | |
| 9 | Verify container in list with stopped status | Container shown with "stopped" status | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-008 - View Container Logs

**Test ID**: TC-008  
**Priority**: Medium  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure container running (from TC-006) | Container "test-web" exists and running | | | |
| 2 | Run `./bin/nebulabox logs test-web` | Command executes | | | |
| 3 | Check for log header | Output shows "📜 Logs for container: test-web" | | | |
| 4 | Verify logs displayed | Logs shown (may be empty if no output yet) | | | |
| 5 | If no logs, verify message | Shows "No logs available" or similar | | | |
| 6 | Verify no error messages | No error output | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-009 - Build Image from BuildSpec

**Test ID**: TC-009  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Create test directory | Directory created | | | |
| 2 | Create buildspec.json file | File created with valid BuildSpec | | | |
| 3 | Run `./bin/nebulabox build . -f buildspec.json -t test-app:latest` | Command executes | | | |
| 4 | Check for "Building image" message | Output shows "🔨 Building image from BuildSpec..." | | | |
| 5 | Check for "Building image" progress | Output shows "📦 Building image: test-app:latest" | | | |
| 6 | Wait for completion | Command completes successfully | | | |
| 7 | Check for success message | Output shows "✅ Image built successfully!" | | | |
| 8 | Verify image ID displayed | Image ID shown in output | | | |
| 9 | Verify image name/tag displayed | "test-app:latest" shown | | | |
| 10 | Run `./bin/nebulabox images` | Command executes | | | |
| 11 | Verify test-app:latest in list | "test-app:latest" appears in image list | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-010 - Delete Image

**Test ID**: TC-010  
**Priority**: Medium  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure image exists (from TC-004 or TC-009) | Image available | | | |
| 2 | Run `./bin/nebulabox rmi test-app:latest` | Command executes | | | |
| 3 | Check for "Removing image" message | Output shows "🗑️ Removing image: test-app:latest" | | | |
| 4 | Check for success message | Output shows "✅ Image removed: test-app:latest" | | | |
| 5 | Verify no error messages | No error output | | | |
| 6 | Run `./bin/nebulabox images` | Command executes | | | |
| 7 | Verify image not in list | "test-app:latest" not in image list | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-041 - Create Hierarchical Container

**Test ID**: TC-041  
**Priority**: High  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Create root container (from TC-006) | Root container "parent" exists | | | |
| 2 | Create child container spec JSON | File "child-spec.json" created | | | |
| 3 | Run `./bin/nebulabox hierarchy create parent --file child-spec.json` | Command executes | | | |
| 4 | Check for success message | Output shows "✅ Nested container created" | | | |
| 5 | Verify parent ID shown | Parent container ID displayed | | | |
| 6 | Verify child ID shown | Child container ID displayed | | | |
| 7 | Verify depth shown | Depth level displayed (e.g., "Depth: 1") | | | |
| 8 | Run `./bin/nebulabox hierarchy list parent` | Command executes | | | |
| 9 | Verify child in hierarchy | Child container listed under parent | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

### Test Case: TC-042 - View Hierarchy Tree

**Test ID**: TC-042  
**Priority**: Medium  
**Status**: ⬜ Not Started | ⬜ In Progress | ✅ Pass | ❌ Fail | ⬜ Blocked

#### Execution Steps

| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure hierarchical containers exist (from TC-041) | Parent and child containers exist | | | |
| 2 | Run `./bin/nebulabox hierarchy tree parent-id` | Command executes | | | |
| 3 | Check for tree header | Output shows "🌳 Hierarchy Tree for: parent-id" | | | |
| 4 | Verify parent shown | Parent container displayed | | | |
| 5 | Verify child shown | Child container displayed with indentation | | | |
| 6 | Verify depth indicators | Depth levels shown correctly | | | |
| 7 | Run `./bin/nebulabox hierarchy tree` (no ID) | Command executes | | | |
| 8 | Verify all trees shown | All hierarchy trees displayed | | | |

**Overall Result**: ⬜ Pass | ❌ Fail  
**Issues Found**: 
_________________________________________________________________
_________________________________________________________________

**Solution/Bug/Fix**: 
_________________________________________________________________
_________________________________________________________________

---

## Test Summary

### Overall Statistics

| Category | Count |
|----------|-------|
| Total Tests Executed | _____ |
| Tests Passed | _____ |
| Tests Failed | _____ |
| Tests Blocked | _____ |
| Pass Rate | _____% |

### Critical Issues

| Issue ID | Test ID | Severity | Description | Status |
|----------|---------|----------|-------------|--------|
| | | | | |
| | | | | |

### Test Coverage

- [ ] Basic CLI Commands (TC-001 to TC-010)
- [ ] Container Lifecycle (TC-011)
- [ ] Image Management (TC-004, TC-005, TC-010)
- [ ] Build Operations (TC-009)
- [ ] Hierarchical Containers (TC-041, TC-042)
- [ ] Container Groups (TC-051, TC-052)
- [ ] Error Handling (TC-061 to TC-063)

---

## Notes and Observations

_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________

---

**End of Test Execution Worksheet**

