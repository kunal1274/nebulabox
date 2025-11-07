# Step-by-Step Manual Testing Guide

## 🚀 Quick Start

Follow these steps in order. Copy and paste each command.

---

## ✅ Step 1: Build the CLI (30 seconds)

```bash
cd /home/serverratxen/Documents/cursor-projects/nebulabox
make build-cli-test
```

**What to check:**
- Should see: `✅ CLI test binary build complete!`
- Binary should exist: `ls -lh bin/nebulabox`

**✅ If successful:** Move to Step 2

---

## ✅ Step 2: Test Help Command (10 seconds)

```bash
./bin/nebulabox --help
```

**What to check:**
- Should see list of available commands
- Should see: `group`, `ps`, `list`, `version`, etc.

**✅ If successful:** Move to Step 3

---

## ✅ Step 3: Test Version Command (10 seconds)

```bash
./bin/nebulabox version
```

**What to check:**
- Should see version information
- Should show: `Version: 0.1.0-alpha`

**✅ If successful:** Move to Step 4

---

## ✅ Step 4: Test Container Listing (10 seconds)

```bash
./bin/nebulabox ps
```

**What to check:**
- Should see container list (may show mock data)
- Should see: `CONTAINER ID`, `IMAGE`, `STATUS`, `NAMES`

**Expected output:**
```
📋 NebulaBox Containers:
CONTAINER ID    IMAGE           STATUS          NAMES
----------------------------------------------------
mock-001        nginx:latest    running         web-server
mock-002        postgres:13     running         db-server
```

**✅ If successful:** Move to Step 5

---

## ✅ Step 5: Test Group Commands (20 seconds)

### 5.1 Show Group Help
```bash
./bin/nebulabox group --help
```

**What to check:**
- Should see group subcommands: `create`, `list`, `start`, `stop`, `status`

### 5.2 List Groups
```bash
./bin/nebulabox group list
```

**What to check:**
- Should see: `No groups found` (this is correct - no groups created yet)

**✅ If successful:** Move to Step 6

---

## ✅ Step 6: Test Engine Compilation (30 seconds)

```bash
go build ./internal/engine/...
```

**What to check:**
- Should complete without errors (silent success)

**✅ If successful:** Move to Step 7

---

## ✅ Step 7: Run Engine Unit Tests (30 seconds)

```bash
go test ./internal/engine/... -v
```

**What to check:**
- Should see 5 tests pass:
  - `TestNewRuntime` ✅
  - `TestRuntimeStoragePaths` ✅
  - `TestListContainersEmpty` ✅
  - `TestListImagesEmpty` ✅
  - `TestListGroupsEmpty` ✅

**Expected output:**
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
```

**✅ If successful:** Move to Step 8

---

## ✅ Step 8: Run CLI Integration Tests (1 minute)

```bash
go test ./internal/cli/tests/... -v
```

**What to check:**
- Should see multiple tests pass
- Some may be skipped (expected): `TestImageList`, `TestImageDelete`

**Expected output:**
```
=== RUN   TestBuildCommand
--- PASS: TestBuildCommand (0.01s)
=== RUN   TestContainerList
--- PASS: TestContainerList (0.01s)
=== RUN   TestContainerRun
--- PASS: TestContainerRun (3.51s)
...
PASS
```

**✅ If successful:** Move to Step 9

---

## ✅ Step 9: Run Automated Test Scripts (1 minute)

### 9.1 Test Engine Basic
```bash
./scripts/testing/test-engine-basic.sh
```

**What to check:**
- Should see: `✅ Basic Engine Tests Complete`
- All tests should pass

### 9.2 Test Containers
```bash
./scripts/testing/test-containers.sh
```

**What to check:**
- Should see: `✅ Container Operations Test Complete`
- Most tests should pass

**✅ If successful:** Move to Step 10

---

## ✅ Step 10: Test Interactive Demo (5 minutes)

```bash
./scripts/cli/demo-poc.sh
```

**What to check:**
- Should see welcome screen
- Press Enter to start
- Navigate through 8 steps using:
  - **Enter** = Continue
  - **s** = Skip
  - **b** = Go back
  - **r** = Restart
  - **q** = Quit

**What you'll see:**
- Step 1: Introduction to NebulaBox
- Step 2: Build Image from BuildSpec
- Step 3: Run Container
- Step 4: Container Grouping
- Step 5: Shared Workspace
- Step 6: Remote Deployment
- Step 7: Unified Cloud Deployment
- Step 8: Summary

**✅ If successful:** All testing complete!

---

## 📋 Testing Checklist

Use this to track your progress:

- [ ] Step 1: CLI builds successfully
- [ ] Step 2: Help command works
- [ ] Step 3: Version command works
- [ ] Step 4: Container listing works
- [ ] Step 5: Group commands work
- [ ] Step 6: Engine compiles
- [ ] Step 7: Engine tests pass (5/5)
- [ ] Step 8: CLI tests pass (7/9, 2 skipped)
- [ ] Step 9: Automated scripts work
- [ ] Step 10: Interactive demo works

---

## 🎯 Expected Results Summary

### ✅ Should Work:
- CLI binary builds
- Help, version, ps, list commands
- Group commands (list, help)
- Engine compilation
- Engine unit tests (5/5)
- CLI integration tests (7/9)
- Automated test scripts
- Interactive demo

### ⚠️ Expected Issues (POC Phase):
- Container create/start/stop may need root
- Build command may need file path format
- Images command not yet implemented
- Some operations show mock data

---

## 🚨 Troubleshooting

### Issue: "Command not found"
**Solution:** Make sure you're in the project directory:
```bash
cd /home/serverratxen/Documents/cursor-projects/nebulabox
```

### Issue: "Permission denied"
**Solution:** Make scripts executable:
```bash
chmod +x scripts/cli/*.sh
chmod +x scripts/testing/*.sh
```

### Issue: "CLI binary not found"
**Solution:** Rebuild:
```bash
make build-cli-test
```

### Issue: Tests fail
**Solution:** Check if it's an expected POC limitation (see Expected Issues above)

---

## 📊 Quick Test Summary

Run this to see everything at once:

```bash
echo "=== Quick Test Summary ===" && \
echo "1. CLI Binary:" && ls -lh bin/nebulabox && \
echo "" && echo "2. Version:" && ./bin/nebulabox version | head -3 && \
echo "" && echo "3. Containers:" && ./bin/nebulabox ps | head -5 && \
echo "" && echo "4. Groups:" && ./bin/nebulabox group list
```

---

## ✅ Success Criteria

**You've successfully tested if:**
1. ✅ CLI builds without errors
2. ✅ Help and version commands work
3. ✅ Container listing shows output
4. ✅ Group commands are available
5. ✅ Engine tests pass (5/5)
6. ✅ CLI tests pass (7/9)
7. ✅ Automated scripts complete
8. ✅ Interactive demo runs

**If all above work → Testing Complete! 🎉**

---

## 📚 Additional Resources

- **Detailed Guide:** `MANUAL_TESTING_GUIDE.md`
- **Quick Steps:** `TEST_STEPS.md`
- **Testing Summary:** `TESTING_SUMMARY.md`

---

## Next Steps After Testing

1. ✅ Document any unexpected issues
2. ✅ Note which features work perfectly
3. ✅ Prepare demo using interactive script
4. ✅ Continue with Phase 3 (API Development)

---

**Happy Testing! 🚀**

