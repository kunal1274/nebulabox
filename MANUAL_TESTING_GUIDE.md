# NebulaBox Manual Testing Guide

## Step-by-Step Testing Instructions

This guide will walk you through testing all NebulaBox components manually.

---

## Prerequisites

1. **Ensure you're in the project directory:**
   ```bash
   cd /home/serverratxen/Documents/cursor-projects/nebulabox
   ```

2. **Check Go is installed:**
   ```bash
   go version
   ```
   Expected: Should show Go version (1.22+)

---

## Step 1: Build the CLI Binary

### 1.1 Build the CLI
```bash
make build-cli-test
```

**Expected Output:**
```
🔨 Building NebulaBox CLI for testing...
go build -o bin/nebulabox ./cmd/nebulabox
✅ CLI test binary build complete!
```

### 1.2 Verify Binary Exists
```bash
ls -lh bin/nebulabox
```

**Expected Output:**
```
-rwxr-xr-x 1 user user [size] bin/nebulabox
```

### 1.3 Test Binary Works
```bash
./bin/nebulabox --help
```

**Expected Output:**
```
NebulaBox is a modern DevOps platform that simplifies containerization, 
orchestration, and deployment workflows. It provides Docker-like functionality 
with enhanced simplicity and intelligence.

Usage:
  nebulabox [command]

Available Commands:
  api         Start the NebulaBox API server
  build       Build an image from a Dockerfile
  completion  Generate the autocompletion script for the specified shell
  group       Manage container groups
  help        Help about any command
  list        List containers
  logs        Fetch the logs of a container
  pull        Pull an image from a registry
  push        Push an image to a registry
  ps          List containers
  run         Run a container from an image
  stop        Stop a running container
  version     Show NebulaBox version information
```

✅ **Checkpoint 1:** If you see the help output, CLI binary is working!

---

## Step 2: Test Version Command

```bash
./bin/nebulabox version
```

**Expected Output:**
```
🚀 NebulaBox
Version:     0.1.0-alpha
Build:      Phase 1 Development
Go Version: 1.22+
Platform:   [your platform]

🎯 Current Phase: Core Workflow Layer
📋 Features: CLI, Basic Commands, Placeholder Integration
🔮 Next: containerd Runtime Integration
```

✅ **Checkpoint 2:** Version command works!

---

## Step 3: Test Container Listing

### 3.1 List Containers (ps command)
```bash
./bin/nebulabox ps
```

**Expected Output:**
```
📋 NebulaBox Containers:
CONTAINER ID    IMAGE           STATUS          NAMES
----------------------------------------------------
mock-001        nginx:latest    running         web-server
mock-002        postgres:13     running         db-server
```

**Note:** This shows mock data (expected in POC phase).

### 3.2 List Containers (list command - alias)
```bash
./bin/nebulabox list
```

**Expected Output:** Same as above

### 3.3 List All Containers (including stopped)
```bash
./bin/nebulabox list --all
```

**Expected Output:** Similar to above, may show more containers

✅ **Checkpoint 3:** Container listing works!

---

## Step 4: Test Group Commands

### 4.1 Show Group Help
```bash
./bin/nebulabox group --help
```

**Expected Output:**
```
Create and manage container groups for flexible architecture testing

Usage:
  nebulabox group [command]

Available Commands:
  create      Create a container group
  list        List container groups
  start       Start a container group
  status      Show group status
  stop        Stop a container group
```

### 4.2 List Groups (should be empty)
```bash
./bin/nebulabox group list
```

**Expected Output:**
```
No groups found
```

✅ **Checkpoint 4:** Group commands are available!

---

## Step 5: Test Engine Compilation

### 5.1 Compile Engine Module
```bash
go build ./internal/engine/...
```

**Expected Output:** No errors (silent success)

### 5.2 Run Engine Unit Tests
```bash
go test ./internal/engine/... -v
```

**Expected Output:**
```
=== RUN   TestNewRuntime
--- PASS: TestNewRuntime (0.00s)
=== RUN   TestRuntimeStoragePaths
--- PASS: TestRuntimeStoragePaths (0.00s)
=== RUN   TestListContainersEmpty
--- PASS: TestListContainersEmpty (0.00s)
=== RUN   TestListImagesEmpty
--- PASS: TestListImagesEmpty (0.00s)
=== RUN   TestListGroupsEmpty
--- PASS: TestListGroupsEmpty (0.00s)
PASS
ok  	github.com/nebulabox/nebulabox/internal/engine	0.012s
```

✅ **Checkpoint 5:** Engine compiles and tests pass!

---

## Step 6: Test CLI Integration Tests

### 6.1 Run CLI Tests
```bash
go test ./internal/cli/tests/... -v
```

**Expected Output:**
```
=== RUN   TestBuildCommand
--- PASS: TestBuildCommand (0.01s)
=== RUN   TestContainerList
--- PASS: TestContainerList (0.01s)
=== RUN   TestContainerRun
--- PASS: TestContainerRun (3.51s)
=== RUN   TestContainerStop
--- PASS: TestContainerStop (0.31s)
=== RUN   TestImageList
--- SKIP: TestImageList (0.00s)
=== RUN   TestImagePull
--- PASS: TestImagePull (2.01s)
=== RUN   TestImageDelete
--- SKIP: TestImageDelete (0.00s)
=== RUN   TestContainerLifecycleIntegration
--- PASS: TestContainerLifecycleIntegration (3.83s)
=== RUN   TestImageWorkflowIntegration
--- PASS: TestImageWorkflowIntegration (2.01s)
PASS
```

✅ **Checkpoint 6:** CLI integration tests pass!

---

## Step 7: Run Automated Test Scripts

### 7.1 Test Engine Basic Functionality
```bash
./scripts/testing/test-engine-basic.sh
```

**Expected Output:**
```
==========================================
Engine Basic Functionality Test
==========================================

Test 1: Compiling engine...
✅ Engine compiles successfully

Test 2: Running engine unit tests...
✅ Engine unit tests pass

Test 3: Compiling CLI...
✅ CLI compiles successfully

Test 4: Building CLI binary...
✅ CLI binary built successfully

Test 5: Checking CLI commands...
✅ CLI binary exists
✅ CLI help command works
✅ Group command exists

==========================================
✅ Basic Engine Tests Complete
==========================================
```

### 7.2 Test Container Operations
```bash
./scripts/testing/test-containers.sh
```

**Expected Output:**
```
==========================================
Container Operations Test
==========================================

Test 1: Listing containers...
✅ List containers command works

Test 2: Listing images...
⚠️  List images may show error (expected in POC)

Test 3: Listing groups...
✅ List groups command works

Test 4: Getting version...
✅ Version command works

==========================================
✅ Container Operations Test Complete
==========================================
```

✅ **Checkpoint 7:** Automated tests pass!

---

## Step 8: Test Workflow Scripts

### 8.1 Test Build Workflow
```bash
./scripts/cli/workflow-01-build.sh
```

**Expected Output:**
```
==========================================
Workflow 01: Build Test
==========================================
📦 Building image from BuildSpec...
   BuildSpec: /tmp/tmp.XXXXXX/buildspec.json
```

**Note:** May show error about build command needing file path (expected in POC)

### 8.2 Test Run Workflow
```bash
./scripts/cli/workflow-02-run.sh
```

**Expected Output:**
```
==========================================
Workflow 02: Run Test
==========================================
🚀 Running container...
```

**Note:** May show errors if container operations need root (expected)

### 8.3 Test Group Workflow
```bash
./scripts/cli/workflow-03-group.sh
```

**Expected Output:**
```
==========================================
Workflow 03: Container Group Test
==========================================
📦 Creating container group...
   Note: Group creation via CLI will be available in next phase
```

✅ **Checkpoint 8:** Workflow scripts are ready!

---

## Step 9: Run Interactive Demo (POC for Investors)

### 9.1 Run Main POC Demo
```bash
./scripts/cli/demo-poc.sh
```

**Expected Output:**
```
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║           NebulaBox - POC Demo for Investors                ║
║                                                              ║
║     Unified Container Platform - Different from Docker       ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

What is NebulaBox?
...
```

**Instructions:**
- Press Enter to start the interactive demo
- Use navigation: Enter (continue), s (skip), b (back), r (restart), q (quit)
- Follow through all 8 steps

### 9.2 Run Interactive Workflow Directly
```bash
./scripts/cli/workflow-00-interactive-demo.sh
```

**Expected Output:** Interactive menu with 8 steps

✅ **Checkpoint 9:** Interactive demo works!

---

## Step 10: Test Group Creation (Manual)

### 10.1 Create a Test Group Spec
```bash
cat > /tmp/test-group.json <<'EOF'
{
  "name": "test-group",
  "strategy": "frontend-backend",
  "containers": [
    {
      "name": "frontend",
      "image": "nginx:alpine",
      "ports": ["3000:80"]
    },
    {
      "name": "backend",
      "image": "alpine:latest",
      "ports": ["5000:5000"]
    }
  ],
  "networking": {
    "internal": true,
    "bridge": "test-bridge"
  }
}
EOF
```

### 10.2 Try to Create Group
```bash
./bin/nebulabox group create --file /tmp/test-group.json
```

**Expected Output:**
- May show error about engine not being fully connected (expected in POC)
- OR may show success if engine client works

### 10.3 List Groups Again
```bash
./bin/nebulabox group list
```

**Expected Output:**
- If group was created: Shows group info
- If not: "No groups found"

✅ **Checkpoint 10:** Group creation structure works!

---

## Step 11: Verify All Components

### 11.1 Check Engine Files
```bash
ls -la internal/engine/
```

**Expected Files:**
```
runtime.go
namespace.go
cgroups.go
filesystem.go
network.go
process.go
images.go
container.go
group.go
runtime_test.go
```

### 11.2 Check CLI Files
```bash
ls -la internal/cli/*.go | grep -v test
```

**Expected Files:**
```
containers.go
engine_client.go
group.go
images.go
logs.go
root.go
run.go
version.go
```

### 11.3 Check Workflow Scripts
```bash
ls -la scripts/cli/*.sh
```

**Expected Files:**
```
demo-poc.sh
workflow-00-interactive-demo.sh
workflow-01-build.sh
workflow-02-run.sh
workflow-03-group.sh
workflow-04-remote.sh
workflow-05-deploy.sh
workflow-06-mern-complete.sh
```

✅ **Checkpoint 11:** All files are in place!

---

## Step 12: Final Verification

### 12.1 Run Complete Test Suite
```bash
# Test engine
go test ./internal/engine/... -v

# Test CLI
go test ./internal/cli/tests/... -v

# Test compilation
go build ./internal/engine/...
go build ./internal/cli/...
```

### 12.2 Check Documentation
```bash
ls -la docs/ | grep -i cli
```

**Expected:**
```
CLI_WORKFLOW_GUIDE.md
```

### 12.3 View Testing Summary
```bash
cat TESTING_SUMMARY.md
```

✅ **Checkpoint 12:** Everything verified!

---

## Testing Checklist

Use this checklist to track your testing:

- [ ] Step 1: CLI binary builds successfully
- [ ] Step 2: Version command works
- [ ] Step 3: Container listing (ps/list) works
- [ ] Step 4: Group commands available
- [ ] Step 5: Engine compiles and tests pass
- [ ] Step 6: CLI integration tests pass
- [ ] Step 7: Automated test scripts work
- [ ] Step 8: Workflow scripts are ready
- [ ] Step 9: Interactive demo works
- [ ] Step 10: Group creation structure works
- [ ] Step 11: All files are in place
- [ ] Step 12: Final verification complete

---

## Expected Issues (POC Phase)

These are **expected** and **normal** in the POC phase:

1. **Container create/start/stop errors**
   - Reason: Requires root privileges for namespaces/cgroups
   - Status: Structure ready, needs root for full testing

2. **Build command errors**
   - Reason: May need BuildSpec file path format
   - Status: Command structure ready

3. **Images command not found**
   - Reason: Not yet implemented
   - Status: Expected, will be added in next phase

4. **Network/filesystem operations**
   - Reason: Require root privileges
   - Status: Structure complete, needs privileged testing

---

## Quick Test Commands Summary

```bash
# Build
make build-cli-test

# Basic tests
./bin/nebulabox --help
./bin/nebulabox version
./bin/nebulabox ps
./bin/nebulabox group list

# Unit tests
go test ./internal/engine/... -v
go test ./internal/cli/tests/... -v

# Automated scripts
./scripts/testing/test-engine-basic.sh
./scripts/testing/test-containers.sh

# Interactive demo
./scripts/cli/demo-poc.sh
```

---

## Success Criteria

✅ **All tests pass if:**
1. CLI binary builds without errors
2. Help and version commands work
3. Container listing shows output (even if mock)
4. Group commands are available
5. Engine unit tests pass (5/5)
6. CLI integration tests pass (7/9, 2 skipped)
7. Automated test scripts complete
8. Interactive demo runs

---

## Next Steps After Testing

1. ✅ Document any issues found
2. ✅ Note which features work vs need root
3. ✅ Prepare demo for investors using interactive script
4. ✅ Continue with Phase 3 (API Development)

---

## Troubleshooting

### Issue: "CLI binary not found"
**Solution:**
```bash
make build-cli-test
```

### Issue: "Permission denied" on scripts
**Solution:**
```bash
chmod +x scripts/cli/*.sh
chmod +x scripts/testing/*.sh
```

### Issue: "No test files" error
**Solution:** This is normal for some packages. Focus on packages with `*_test.go` files.

### Issue: Container operations fail
**Solution:** This is expected in POC. Full operations need root privileges.

---

## Questions?

If you encounter any unexpected errors or issues, check:
1. Go version (should be 1.22+)
2. You're in the correct directory
3. All files are present (use `ls` commands above)
4. Scripts are executable (`chmod +x`)

Happy Testing! 🚀

