# NebulaBox CLI Workflow Guide

## Overview

NebulaBox CLI provides a unified interface for container management, from development to deployment. This guide covers all CLI workflows and demonstrates how NebulaBox differs from Docker and Kubernetes.

## Installation

```bash
# Build from source
make build-cli-test

# The binary will be at: bin/nebulabox
```

## Quick Start

### Interactive Demo

Run the interactive POC demo to see NebulaBox's unique features:

```bash
./scripts/cli/demo-poc.sh
```

Or run the interactive workflow:

```bash
./scripts/cli/workflow-00-interactive-demo.sh
```

## Basic Commands

### Build Image

Build an image from a BuildSpec (JSON format, not Dockerfile):

```bash
nebulabox build -f buildspec.json -t my-app:latest
```

**Key Difference from Docker:**
- Docker: Uses Dockerfile (text-based)
- NebulaBox: Uses BuildSpec (structured JSON, easier to generate/modify)

### Run Container

Run a container from an image:

```bash
nebulabox run my-app:latest --name my-container -p 3000:80 -d
```

**Key Difference from Docker:**
- Docker: One container = one service
- NebulaBox: One container can run multiple services (MERN stack together)

### List Containers

```bash
nebulabox ps
nebulabox ps -a  # All containers
```

### Stop Container

```bash
nebulabox stop my-container
```

### View Logs

```bash
nebulabox logs my-container
nebulabox logs my-container -f  # Follow logs
```

## Container Groups (Unique Feature)

NebulaBox allows flexible container grouping for testing different architectures.

### Create Group

Create a container group from a group specification:

```bash
nebulabox group create --file group.json
```

Example `group.json`:

```json
{
  "name": "mern-app",
  "strategy": "frontend-backend",
  "containers": [
    {
      "name": "frontend",
      "image": "nginx:alpine",
      "ports": ["3000:80"]
    },
    {
      "name": "backend-db",
      "image": "node:alpine",
      "ports": ["5000:5000", "27017:27017"]
    }
  ],
  "networking": {
    "internal": true,
    "bridge": "mern-bridge"
  }
}
```

### Group Strategies

- **monolithic**: All services in one container
- **frontend-backend**: Frontend separate, backend+DB together
- **three-tier**: Frontend, Backend, Database separate
- **microservices**: Each service in own container
- **custom**: Your own architecture

### Start/Stop Group

```bash
nebulabox group start group-id
nebulabox group stop group-id
```

### List Groups

```bash
nebulabox group list
```

### Group Status

```bash
nebulabox group status group-id
```

**Key Difference from Docker/Kubernetes:**
- Docker: Need Docker Compose for groups (separate tool)
- Kubernetes: Complex YAML, overkill for small apps
- NebulaBox: Simple JSON config, test different architectures easily

## Workflow Scripts

### Workflow 01: Build Test

```bash
./scripts/cli/workflow-01-build.sh
```

Tests basic build functionality.

### Workflow 02: Run Test

```bash
./scripts/cli/workflow-02-run.sh
```

Tests container run, list, and stop operations.

### Workflow 03: Group Test

```bash
./scripts/cli/workflow-03-group.sh
```

Tests container grouping functionality.

### Workflow 04: Remote Deployment

```bash
./scripts/cli/workflow-04-remote.sh
```

Tests remote deployment (Phase 4 feature).

### Workflow 05: Cloud Deployment

```bash
./scripts/cli/workflow-05-deploy.sh
```

Tests cloud deployment (Phase 6 feature).

### Workflow 06: MERN Complete

```bash
./scripts/cli/workflow-06-mern-complete.sh
```

Complete MERN todo app workflow from build to deployment.

## What Makes NebulaBox Different

### 1. Unified Development

**Docker Approach:**
```bash
# Need multiple containers
docker run -d --name mongodb mongo
docker run -d --name backend --link mongodb node
docker run -d --name frontend nginx
# Or use Docker Compose (separate tool)
```

**NebulaBox Approach:**
```bash
# Single container with all services
nebulabox run mern-todo:latest --name todo-app
# All services (MongoDB, Express, React, Nginx) run together
```

### 2. Built-in Collaboration

**Traditional Approach:**
- Setup VPN (complex)
- Use ngrok/Tailscale (external tools)
- Manual port forwarding

**NebulaBox Approach:**
- Built-in tunneling (automatic)
- Real-time file sync
- Shared workspaces (coming in Phase 4)

### 3. Flexible Architecture Testing

**Docker/Kubernetes:**
- Need to rewrite configs for different architectures
- Complex YAML/Compose files

**NebulaBox:**
- Switch strategies easily
- Test monolithic → microservices → custom
- Perfect for POC and experimentation

### 4. Unified Deployment

**Current Fragmented Approach:**
- Frontend → Vercel
- Backend → Render
- Database → MongoDB Atlas
- Result: 3 different places, 3 different configs

**NebulaBox Approach:**
- Everything → NebulaBox Cloud
- Same BuildSpec for all environments
- Single URL, single configuration

## Next Steps

1. **Try the Interactive Demo**: `./scripts/cli/demo-poc.sh`
2. **Run Workflow Scripts**: Test each workflow individually
3. **Build Your First App**: Create a BuildSpec and build an image
4. **Test Container Groups**: Experiment with different architectures

## Troubleshooting

### CLI Binary Not Found

```bash
make build-cli-test
```

### Build Fails

Check that:
- BuildSpec JSON is valid
- Base image exists
- You have write permissions

### Container Won't Start

Check that:
- Image exists
- Ports aren't already in use
- You have necessary permissions

## See Also

- [BuildSpec Documentation](../docs/BUILD_SPEC.md)
- [Unified Architecture](../docs/UNIFIED_ARCHITECTURE.md)
- [MERN Todo App Guide](../mvp/app/mern-todo/README.md)

