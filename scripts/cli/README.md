# NebulaBox CLI Workflow Scripts

This directory contains workflow test scripts for the NebulaBox CLI.

## Interactive Demo (POC for Investors)

### Main Demo Script

```bash
./demo-poc.sh
```

This is the main demo script that showcases NebulaBox's unique features. It's designed to be shown to investors or used as a proof-of-concept demonstration.

### Interactive Workflow

```bash
./workflow-00-interactive-demo.sh
```

An interactive end-to-end workflow with navigation:
- **Enter** - Continue to next step
- **s** - Skip this step
- **b** - Go back to previous step
- **r** - Restart from beginning
- **q** - Quit demo

## Individual Workflow Scripts

### Workflow 01: Build Test
```bash
./workflow-01-build.sh
```
Tests basic build functionality from BuildSpec.

### Workflow 02: Run Test
```bash
./workflow-02-run.sh
```
Tests container run, list, and stop operations.

### Workflow 03: Group Test
```bash
./workflow-03-group.sh
```
Tests container grouping functionality (unique feature).

### Workflow 04: Remote Deployment
```bash
./workflow-04-remote.sh
```
Tests remote deployment (Phase 4 feature - coming soon).

### Workflow 05: Cloud Deployment
```bash
./workflow-05-deploy.sh
```
Tests cloud deployment (Phase 6 feature - coming soon).

### Workflow 06: MERN Complete
```bash
./workflow-06-mern-complete.sh
```
Complete MERN todo app workflow from build to deployment.

## Prerequisites

1. Build the CLI binary:
   ```bash
   make build-cli-test
   ```

2. Ensure scripts are executable:
   ```bash
   chmod +x scripts/cli/*.sh
   ```

## Usage

### For Investors/Demo

Run the main POC demo:
```bash
cd scripts/cli
./demo-poc.sh
```

This will:
1. Show NebulaBox's unique features
2. Demonstrate differences from Docker/Kubernetes
3. Walk through interactive workflow

### For Testing

Run individual workflows:
```bash
cd scripts/cli
./workflow-01-build.sh
./workflow-02-run.sh
# etc.
```

## What Makes NebulaBox Different

The demo scripts highlight:

1. **Unified Development** - Single container for entire stack
2. **Built-in Collaboration** - No VPN/ngrok needed
3. **Flexible Architecture** - Test different strategies easily
4. **Unified Deployment** - One platform, not fragmented

## See Also

- [CLI Workflow Guide](../../docs/CLI_WORKFLOW_GUIDE.md)
- [Unified Architecture](../../docs/UNIFIED_ARCHITECTURE.md)

