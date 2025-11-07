# CLI Testing Framework

## Overview
Comprehensive testing framework for NebulaBox CLI commands.

## Structure

```
internal/cli/tests/
├── mocks/          # Mock implementations
│   ├── api_mock.go     # Mock API client
│   └── containerd_mock.go  # Mock containerd client
├── fixtures/       # Test fixtures
│   └── testdata/       # Sample data files
├── helpers/        # Test helper functions
│   ├── command_runner.go   # Execute commands and capture output
│   └── assertions.go       # Custom assertions
└── testutils/      # Shared test utilities
    └── test_utils.go
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

## Commands to Test

- `nebulabox run` - Start containers
- `nebulabox list` - List containers
- `nebulabox stop` - Stop containers
- `nebulabox logs` - View logs
- `nebulabox build` - Build images
- `nebulabox push/pull` - Registry operations
- `nebulabox version` - Version info

