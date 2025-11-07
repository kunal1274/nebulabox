# CLI Testing Framework Plan

## Overview
Set up a comprehensive CLI testing framework for NebulaBox CLI commands.

## Structure

```
internal/cli/
├── cmd/           # CLI commands
├── tests/         # Test utilities and framework
│   ├── mocks/     # Mock implementations
│   ├── fixtures/  # Test fixtures
│   └── helpers/   # Test helper functions
└── testutils/     # Shared test utilities
```

## Testing Strategy

### 1. Unit Tests (90%+ coverage)
- Test each command in isolation
- Mock API calls and containerd operations
- Test validation, error handling, output formatting

### 2. Integration Tests
- Test CLI → API server communication
- Test with real containerd (optional)
- Test with mock mode

### 3. E2E CLI Workflow Tests
- Complete user journeys via CLI
- Test command chaining
- Test state persistence

## Test Utilities Needed

1. **Mock API Client** - Mock HTTP responses
2. **Mock Containerd Client** - Mock container operations
3. **Test Fixtures** - Sample containers, images, workspaces
4. **Command Runner** - Helper to execute commands and capture output
5. **Assertion Helpers** - Custom assertions for CLI output

## Commands to Test

- `nebulabox container list`
- `nebulabox container run`
- `nebulabox container start/stop/delete`
- `nebulabox image list/pull/build`
- `nebulabox workspace create/share`
- `nebulabox build`
- `nebulabox deploy`

