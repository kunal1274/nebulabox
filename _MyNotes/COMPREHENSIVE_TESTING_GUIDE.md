# Comprehensive Testing Guide - NebulaBox CLI

## Table of Contents
1. [User Manual](#user-manual)
2. [Test Cases Overview](#test-cases-overview)
3. [User Journeys](#user-journeys)
4. [Detailed Test Cases](#detailed-test-cases)
5. [Data Flow Diagrams](#data-flow-diagrams)
6. [Test Execution Template](#test-execution-template)

---

## User Manual

### Installation & Setup

#### Prerequisites
- Linux system (Ubuntu 20.04+ recommended)
- Go 1.22+ installed
- Root/sudo access (for container operations)
- 2GB+ free disk space

#### Build NebulaBox
```bash
cd /home/serverratxen/Documents/cursor-projects/nebulabox
make build-cli-test
```

#### Verify Installation
```bash
# Option 1: Use nbx (short form) - recommended
nbx version
nbx --help

# Option 2: Use full path
./bin/nbx version
./bin/nbx --help

# Option 3: Use full command name (if installed)
nebulabox version
nebulabox --help
```

### Basic Commands

#### Container Management
```bash
# List containers
nbx ps
nbx list

# Run a container
nbx run nginx:latest --name web-server

# Stop a container
nbx stop web-server

# View logs
nbx logs web-server
```

#### Image Management
```bash
# List images
nbx images

# Pull an image
nbx pull nginx:latest

# Remove an image
nbx rmi nginx:latest
```

#### Build Images
```bash
# Build from BuildSpec
nbx build . -f buildspec.json -t myapp:latest
```

#### Hierarchical Containers
```bash
# Create nested container
nbx hierarchy create parent-id --file child-spec.json

# View hierarchy tree
nbx hierarchy tree container-id

# List all hierarchies
nbx hierarchy list
```

#### Container Groups
```bash
# Create a group
nbx group create --file group-spec.json

# List groups
nbx group list

# Start a group
nbx group start group-id
```

---

## Test Cases Overview

### Test Categories
1. **TC-001 to TC-010**: Basic CLI Commands
2. **TC-011 to TC-020**: Container Lifecycle
3. **TC-021 to TC-030**: Image Management
4. **TC-031 to TC-040**: Build Operations
5. **TC-041 to TC-050**: Hierarchical Containers
6. **TC-051 to TC-060**: Container Groups
7. **TC-061 to TC-070**: Integration Tests
8. **TC-071 to TC-080**: Error Handling

---

## User Journeys

### Journey 1: First-Time User - Basic Container Operations
**User Story**: As a new user, I want to run my first container so I can understand how NebulaBox works.

**Flow**:
1. Check version → 2. List containers (empty) → 3. Pull image → 4. Run container → 5. List containers (with data) → 6. View logs → 7. Stop container

### Journey 2: Developer - Build and Deploy Application
**User Story**: As a developer, I want to build an image from my code and run it as a container.

**Flow**:
1. Create BuildSpec → 2. Build image → 3. List images → 4. Run container from built image → 5. Verify running → 6. Stop container

### Journey 3: DevOps Engineer - Hierarchical Container Management
**User Story**: As a DevOps engineer, I want to create nested containers and groups for complex applications.

**Flow**:
1. Create root container → 2. Create nested container → 3. Create group → 4. Add containers to group → 5. View hierarchy → 6. Start group → 7. Verify all running

### Journey 4: System Administrator - Image Management
**User Story**: As a system admin, I want to manage images efficiently (pull, list, delete).

**Flow**:
1. Pull multiple images → 2. List all images → 3. Identify unused → 4. Delete unused images → 5. Verify deletion

---

## Detailed Test Cases

### TC-001: CLI Version Command
**Test ID**: TC-001  
**Test Description**: Verify version command displays correct information  
**User Story**: As a user, I want to check NebulaBox version to ensure I'm using the correct version.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox version` | Command executes without error | | | |
| 2 | Check output contains "NebulaBox" | Output shows "🚀 NebulaBox" | | | |
| 3 | Check version number | Output shows "Version: 0.1.0-alpha" | | | |
| 4 | Check build info | Output shows "Build: Phase 1 Development" | | | |
| 5 | Check Go version | Output shows "Go Version: 1.22+" | | | |

**Data Flow**:
```
User Input → CLI Parser → Version Command → Display Version Info → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-002: CLI Help Command
**Test ID**: TC-002  
**Test Description**: Verify help command shows all available commands  
**User Story**: As a user, I want to see all available commands when I need help.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
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

**Data Flow**:
```
User Input → CLI Parser → Help Command → Load Command Registry → Display Commands → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-003: List Containers (Empty State)
**Test ID**: TC-003  
**Test Description**: Verify list command shows empty state correctly  
**User Story**: As a user, I want to see when no containers are running.  
**Priority**: Medium  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox ps` | Command executes without error | | | |
| 2 | Check output | Shows "No containers found" or empty table | | | |
| 3 | Verify no error messages | No error output | | | |
| 4 | Run `./bin/nebulabox list` | Command executes without error | | | |
| 5 | Verify same output as ps | Same output as step 2 | | | |

**Data Flow**:
```
User Input → CLI Parser → List Command → Engine Client → Runtime.ListContainers() → 
Empty Result → Format Output → "No containers found" → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-004: Pull Image
**Test ID**: TC-004  
**Test Description**: Verify pull command downloads images successfully  
**User Story**: As a user, I want to pull images from registries to use in containers.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox pull nginx:latest` | Command executes | | | |
| 2 | Check for "Pulling image" message | Output shows "⬇️ Pulling image: nginx:latest" | | | |
| 3 | Wait for completion | Command completes successfully | | | |
| 4 | Check for success message | Output shows "✅ Image pulled successfully" | | | |
| 5 | Verify no error messages | No error output | | | |
| 6 | Run `./bin/nebulabox images` | Command executes | | | |
| 7 | Verify nginx:latest in list | nginx:latest appears in image list | | | |

**Data Flow**:
```
User Input → CLI Parser → Pull Command → Engine Client → Runtime.PullImage() → 
Image Manager → Download Image → Store Image → Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-005: List Images
**Test ID**: TC-005  
**Test Description**: Verify images command lists all available images  
**User Story**: As a user, I want to see all images I have available.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox images` | Command executes without error | | | |
| 2 | Check for table header | Shows "IMAGE ID", "REPOSITORY", "TAG", "SIZE", "CREATED" | | | |
| 3 | If images exist, verify format | Each image shows: ID (12 chars), repo, tag, size, created date | | | |
| 4 | If no images, verify message | Shows "No images found" | | | |
| 5 | Verify no error messages | No error output | | | |

**Data Flow**:
```
User Input → CLI Parser → Images Command → Engine Client → Runtime.ListImages() → 
Image Manager → Get All Images → Format Table → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-006: Run Container
**Test ID**: TC-006  
**Test Description**: Verify run command creates and starts containers  
**User Story**: As a user, I want to run containers from images.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
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

**Data Flow**:
```
User Input → CLI Parser → Run Command → Engine Client → 
Runtime.CreateContainer() → Container Manager → Create Namespaces → 
Create CGroup → Create Filesystem → Create Network → 
Runtime.StartContainer() → Process Manager → Start Process → 
Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-007: Stop Container
**Test ID**: TC-007  
**Test Description**: Verify stop command stops running containers  
**User Story**: As a user, I want to stop running containers.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
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

**Data Flow**:
```
User Input → CLI Parser → Stop Command → Engine Client → 
Runtime.StopContainer() → Process Manager → Send SIGTERM → 
Wait for Exit → Update Container State → Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-008: View Container Logs
**Test ID**: TC-008  
**Test Description**: Verify logs command displays container logs  
**User Story**: As a user, I want to view container logs for debugging.  
**Priority**: Medium  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure container running (from TC-006) | Container "test-web" exists and running | | | |
| 2 | Run `./bin/nebulabox logs test-web` | Command executes | | | |
| 3 | Check for log header | Output shows "📜 Logs for container: test-web" | | | |
| 4 | Verify logs displayed | Logs shown (may be empty if no output yet) | | | |
| 5 | If no logs, verify message | Shows "No logs available" or similar | | | |
| 6 | Verify no error messages | No error output | | | |

**Data Flow**:
```
User Input → CLI Parser → Logs Command → Engine Client → 
Runtime.GetContainerLogs() → Process Manager → CollectLogs() → 
Format Logs → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-009: Build Image from BuildSpec
**Test ID**: TC-009  
**Test Description**: Verify build command creates images from BuildSpec  
**User Story**: As a developer, I want to build images from BuildSpec files.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
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

**Data Flow**:
```
User Input → CLI Parser → Build Command → Read BuildSpec File → 
Parse JSON → Engine Client → Runtime.BuildImage() → 
Image Manager → Process Build Steps → Create Image Layers → 
Store Image → Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-010: Delete Image
**Test ID**: TC-010  
**Test Description**: Verify rmi command deletes images  
**User Story**: As a user, I want to remove unused images to free space.  
**Priority**: Medium  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure image exists (from TC-004 or TC-009) | Image available | | | |
| 2 | Run `./bin/nebulabox rmi test-app:latest` | Command executes | | | |
| 3 | Check for "Removing image" message | Output shows "🗑️ Removing image: test-app:latest" | | | |
| 4 | Check for success message | Output shows "✅ Image removed: test-app:latest" | | | |
| 5 | Verify no error messages | No error output | | | |
| 6 | Run `./bin/nebulabox images` | Command executes | | | |
| 7 | Verify image not in list | "test-app:latest" not in image list | | | |

**Data Flow**:
```
User Input → CLI Parser → Rmi Command → Engine Client → 
Runtime.DeleteImage() → Image Manager → Check Dependencies → 
Delete Image Files → Remove from Registry → Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-011: Container Lifecycle - Complete Flow
**Test ID**: TC-011  
**Test Description**: Verify complete container lifecycle (create, start, stop, delete)  
**User Story**: As a user, I want to manage containers through their complete lifecycle.  
**Priority**: High  
**Test Type**: Integration

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Pull image | Image pulled successfully (TC-004) | | | |
| 2 | Run container | Container created and started (TC-006) | | | |
| 3 | Verify running | Container appears in ps with "running" status | | | |
| 4 | View logs | Logs accessible (TC-008) | | | |
| 5 | Stop container | Container stopped successfully (TC-007) | | | |
| 6 | Verify stopped | Container status shows "stopped" | | | |
| 7 | Start container again | Container can be started again | | | |
| 8 | Verify running again | Container status shows "running" | | | |
| 9 | Stop and verify cleanup | Container stops and resources released | | | |

**Data Flow**:
```
[Complete Lifecycle Flow]
Pull Image → Create Container → Start Container → Running State → 
Stop Container → Stopped State → (Optional: Start Again) → 
Delete Container → Cleanup Resources
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-041: Create Hierarchical Container
**Test ID**: TC-041  
**Test Description**: Verify hierarchical container creation (container within container)  
**User Story**: As a DevOps engineer, I want to create nested containers for complex architectures.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
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

**Data Flow**:
```
User Input → CLI Parser → Hierarchy Create Command → Read Spec File → 
Engine Client → Hierarchy.CreateNestedContainer() → 
Verify Parent Exists → Create Child Container → 
Link to Parent → Update Hierarchy → Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-042: View Hierarchy Tree
**Test ID**: TC-042  
**Test Description**: Verify hierarchy tree command displays container relationships  
**User Story**: As a user, I want to visualize container hierarchies.  
**Priority**: Medium  
**Test Type**: Functional

#### Test Steps
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

**Data Flow**:
```
User Input → CLI Parser → Hierarchy Tree Command → Engine Client → 
Hierarchy.GetHierarchy() → Build Tree Structure → 
Format Tree Output → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-051: Create Container Group
**Test ID**: TC-051  
**Test Description**: Verify group creation command works  
**User Story**: As a DevOps engineer, I want to group containers for orchestration.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Create group spec JSON | File "group-spec.json" created | | | |
| 2 | Run `./bin/nebulabox group create --file group-spec.json` | Command executes | | | |
| 3 | Check for success message | Output shows "✅ Group created successfully" | | | |
| 4 | Verify group ID shown | Group ID displayed | | | |
| 5 | Verify group name shown | Group name displayed | | | |
| 6 | Run `./bin/nebulabox group list` | Command executes | | | |
| 7 | Verify group in list | Group appears in list | | | |

**Data Flow**:
```
User Input → CLI Parser → Group Create Command → Read Spec File → 
Engine Client → Runtime.CreateGroup() → Group Manager → 
Create Group → Store Group → Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-052: Add Container to Group
**Test ID**: TC-052  
**Test Description**: Verify adding container to group (cross-hierarchy)  
**User Story**: As a user, I want to add containers from any hierarchy level to groups.  
**Priority**: High  
**Test Type**: Functional

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Ensure container exists (from TC-006) | Container "test-web" exists | | | |
| 2 | Ensure group exists (from TC-051) | Group "web-group" exists | | | |
| 3 | Run `./bin/nebulabox hierarchy add-group test-web web-group` | Command executes | | | |
| 4 | Check for success message | Output shows "✅ Container test-web added to group web-group" | | | |
| 5 | Verify no error messages | No error output | | | |
| 6 | Run `./bin/nebulabox group status web-group` | Command executes | | | |
| 7 | Verify container in group | "test-web" appears in group container list | | | |

**Data Flow**:
```
User Input → CLI Parser → Hierarchy Add-Group Command → Engine Client → 
Hierarchy.AddContainerToGroup() → Verify Container Exists → 
Verify Group Exists → Add Container to Group → Update Group → 
Success Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-061: Error Handling - Invalid Command
**Test ID**: TC-061  
**Test Description**: Verify error handling for invalid commands  
**User Story**: As a user, I want clear error messages when I make mistakes.  
**Priority**: Medium  
**Test Type**: Negative

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox invalid-command` | Command executes | | | |
| 2 | Check for error message | Output shows "Error: unknown command" | | | |
| 3 | Verify suggestion shown | Suggests similar commands or shows help | | | |
| 4 | Verify exit code | Exit code is non-zero | | | |

**Data Flow**:
```
User Input → CLI Parser → Command Not Found → Error Handler → 
Generate Error Message → Suggest Alternatives → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-062: Error Handling - Missing Required Arguments
**Test ID**: TC-062  
**Test Description**: Verify error handling for missing arguments  
**User Story**: As a user, I want clear errors when required arguments are missing.  
**Priority**: Medium  
**Test Type**: Negative

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox run` | Command executes | | | |
| 2 | Check for error message | Output shows "Error: requires at least 1 arg(s)" | | | |
| 3 | Verify usage shown | Usage information displayed | | | |
| 4 | Run `./bin/nebulabox stop` | Command executes | | | |
| 5 | Check for error message | Output shows similar error for missing container ID | | | |

**Data Flow**:
```
User Input → CLI Parser → Validate Arguments → Missing Required Arg → 
Error Handler → Generate Error Message → Show Usage → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

### TC-063: Error Handling - Non-existent Container
**Test ID**: TC-063  
**Test Description**: Verify error handling for operations on non-existent containers  
**User Story**: As a user, I want clear errors when I reference containers that don't exist.  
**Priority**: Medium  
**Test Type**: Negative

#### Test Steps
| Step | Action | Expected Outcome | Actual Outcome | Pass/Fail | Notes |
|------|--------|------------------|----------------|-----------|-------|
| 1 | Run `./bin/nebulabox stop nonexistent-container` | Command executes | | | |
| 2 | Check for error message | Output shows "container not found" or similar | | | |
| 3 | Verify error is clear | Error message is descriptive | | | |
| 4 | Run `./bin/nebulabox logs nonexistent-container` | Command executes | | | |
| 5 | Check for error message | Similar error message shown | | | |

**Data Flow**:
```
User Input → CLI Parser → Stop/Logs Command → Engine Client → 
Runtime.GetContainer() → Container Not Found → Error Handler → 
Generate Error Message → User Output
```

**Solution/Bug/Fix**: [To be filled during testing]

---

## Data Flow Diagrams

### Container Creation Flow
```
┌─────────┐
│  User   │
└────┬────┘
     │ run nginx:latest --name web
     ▼
┌─────────────────┐
│  CLI Parser     │
└────┬────────────┘
     │ Parse arguments
     ▼
┌─────────────────┐
│ Engine Client   │
└────┬────────────┘
     │ CreateContainer(spec)
     ▼
┌─────────────────┐
│ NebulaRuntime   │
└────┬────────────┘
     │
     ├─► NamespaceManager.CreateNamespaceSet()
     ├─► CGroupManager.CreateCGroup()
     ├─► FilesystemManager.CreateRootfs()
     ├─► NetworkManager.SetupContainerNetwork()
     └─► ProcessManager.StartContainer()
     │
     ▼
┌─────────────────┐
│ Container State │
│   (Running)     │
└─────────────────┘
```

### Image Pull Flow
```
┌─────────┐
│  User   │
└────┬────┘
     │ pull nginx:latest
     ▼
┌─────────────────┐
│  CLI Parser     │
└────┬────────────┘
     │ Parse image name
     ▼
┌─────────────────┐
│ Engine Client   │
└────┬────────────┘
     │ PullImage(image)
     ▼
┌─────────────────┐
│ NebulaRuntime   │
└────┬────────────┘
     │
     ├─► ImageManager.PullImage()
     │   ├─► Contact Registry
     │   ├─► Download Manifest
     │   ├─► Download Layers
     │   └─► Store Image
     │
     ▼
┌─────────────────┐
│ Image Stored    │
│  (Available)    │
└─────────────────┘
```

### Hierarchical Container Flow
```
┌─────────┐
│  User   │
└────┬────┘
     │ hierarchy create parent --file child-spec.json
     ▼
┌─────────────────┐
│  CLI Parser     │
└────┬────────────┘
     │ Parse arguments
     ▼
┌─────────────────┐
│ Engine Client   │
└────┬────────────┘
     │ CreateNestedContainer(parentID, spec)
     ▼
┌─────────────────┐
│ HierarchyManager│
└────┬────────────┘
     │
     ├─► Verify Parent Exists
     ├─► Create Child Container
     ├─► Link Child to Parent
     ├─► Update Hierarchy Tree
     └─► Set Depth Level
     │
     ▼
┌─────────────────┐
│ Nested Container │
│  (In Hierarchy)  │
└─────────────────┘
```

---

## Test Execution Template

### Test Execution Log

**Test Session**: [Date/Time]  
**Tester**: [Name]  
**Environment**: [OS, Go Version, etc.]

| Test ID | Test Description | Status | Actual Outcome | Issues Found | Solution/Fix |
|---------|------------------|--------|----------------|--------------|--------------|
| TC-001 | CLI Version Command | | | | |
| TC-002 | CLI Help Command | | | | |
| TC-003 | List Containers (Empty) | | | | |
| TC-004 | Pull Image | | | | |
| TC-005 | List Images | | | | |
| TC-006 | Run Container | | | | |
| TC-007 | Stop Container | | | | |
| TC-008 | View Container Logs | | | | |
| TC-009 | Build Image from BuildSpec | | | | |
| TC-010 | Delete Image | | | | |
| TC-011 | Container Lifecycle | | | | |
| TC-041 | Create Hierarchical Container | | | | |
| TC-042 | View Hierarchy Tree | | | | |
| TC-051 | Create Container Group | | | | |
| TC-052 | Add Container to Group | | | | |
| TC-061 | Error Handling - Invalid Command | | | | |
| TC-062 | Error Handling - Missing Args | | | | |
| TC-063 | Error Handling - Non-existent Container | | | | |

### Test Summary

**Total Tests**: [Count]  
**Passed**: [Count]  
**Failed**: [Count]  
**Blocked**: [Count]  
**Pass Rate**: [Percentage]%

### Issues Log

| Issue ID | Test ID | Severity | Description | Status | Solution |
|----------|---------|----------|-------------|--------|----------|
| | | | | | |

---

## Sample Test Data

### BuildSpec Example (buildspec.json)
```json
{
  "version": "1.0",
  "name": "test-app",
  "tag": "latest",
  "base": {
    "image": "alpine",
    "tag": "latest"
  },
  "workdir": "/app",
  "env": {
    "NODE_ENV": "production"
  },
  "expose": [3000],
  "steps": [
    {
      "type": "run",
      "command": "apk add --no-cache nodejs npm",
      "comment": "Install Node.js"
    },
    {
      "type": "copy",
      "source": "./",
      "dest": "/app",
      "comment": "Copy application files"
    },
    {
      "type": "run",
      "command": "npm install",
      "comment": "Install dependencies"
    }
  ],
  "health": {
    "type": "http",
    "path": "/health",
    "port": 3000,
    "interval": 30,
    "timeout": 5,
    "retries": 3
  },
  "labels": {
    "maintainer": "test@example.com",
    "version": "1.0.0"
  }
}
```

### Container Spec Example (child-spec.json)
```json
{
  "id": "child-container",
  "name": "child",
  "image": "alpine:latest",
  "command": ["/bin/sh", "-c", "sleep 3600"],
  "env": {
    "PARENT_ID": "parent-container"
  },
  "ports": {},
  "volumes": {},
  "labels": {
    "type": "child",
    "parent": "parent-container"
  }
}
```

### Group Spec Example (group-spec.json)
```json
{
  "name": "web-group",
  "strategy": "frontend-backend",
  "containers": [
    {
      "name": "frontend",
      "image": "nginx:alpine",
      "ports": ["3000:80"]
    },
    {
      "name": "backend",
      "image": "node:18",
      "ports": ["5000:5000"]
    }
  ],
  "networking": {
    "internal": true,
    "bridge": "web-bridge"
  }
}
```

---

## Quick Reference

### Command Cheat Sheet
```bash
# Basic Info
nbx version                     # or: nebulabox version
nbx --help                      # or: nebulabox --help

# Containers
nbx ps                          # List containers
nbx run IMAGE [OPTIONS]         # Run container
nbx stop CONTAINER              # Stop container
nbx logs CONTAINER              # View logs

# Images
nbx images                      # List images
nbx pull IMAGE                  # Pull image
nbx rmi IMAGE                   # Remove image
nbx build PATH [OPTIONS]        # Build image

# Hierarchy
nbx hierarchy create PARENT --file SPEC
nbx hierarchy list [ROOT]
nbx hierarchy tree [CONTAINER]
nbx hierarchy add-group CONTAINER GROUP

# Groups
nbx group create --file SPEC
nbx group list
nbx group start GROUP
nbx group stop GROUP
nbx group status GROUP
```

---

**End of Comprehensive Testing Guide**

