# Quick Testing Steps - Copy & Paste

## Fastest Way to Test Everything

Just copy and paste these commands one by one:

### 1. Build CLI
```bash
make build-cli-test
```

### 2. Test Basic Commands
```bash
./bin/nebulabox --help
./bin/nebulabox version
./bin/nebulabox ps
./bin/nebulabox group list
```

### 3. Test Engine
```bash
go build ./internal/engine/...
go test ./internal/engine/... -v
```

### 4. Test CLI Integration
```bash
go test ./internal/cli/tests/... -v
```

### 5. Run Automated Tests
```bash
./scripts/testing/test-engine-basic.sh
./scripts/testing/test-containers.sh
```

### 6. Run Interactive Demo
```bash
./scripts/cli/demo-poc.sh
```

---

## Expected Results

✅ **All should work:**
- CLI builds
- Commands respond
- Tests pass
- Demo runs

⚠️ **Expected issues (POC phase):**
- Container create/start needs root
- Some commands show mock data
- Build command may need file path

---

## Full Manual Guide

For detailed step-by-step instructions, see:
```bash
cat MANUAL_TESTING_GUIDE.md
```

