# NebulaBox 🚀

A modern DevOps platform that simplifies containerization, orchestration, and deployment workflows.

## Vision

NebulaBox aims to provide Docker-like functionality with a focus on:
- **Simplicity**: No Kubernetes complexity for small and mid companies
- **Intelligence**: AI guidance built into deployment and troubleshooting  
- **Business-first**: Especially optimized for SaaS apps like ERPs

## Development Phases

### Phase 1: Core Workflow Layer (Current)
- CLI tool with containerd backend
- Web dashboard for container management
- Basic image registry
- Simple CI/CD pipelines

### Phase 2: Custom Registry + Build System
- Private image registry
- Custom build specifications
- Security scanning and signing

### Phase 3: Cloud Orchestrator
- Multi-node deployment
- Auto-scaling and health checks
- Service discovery

### Phase 4: Runtime Innovation
- Custom container engine
- Enhanced security and performance

### Phase 5: AI-driven Ops
- Predictive analytics
- Auto-configuration
- ChatOps interface

## Quick Start

```bash
# Install dependencies
brew install go containerd

# Build NebulaBox
go build -o nebulabox ./cmd/nebulabox

# Run your first container
./nebulabox run nginx
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   NebulaBox     │    │   NebulaBox     │    │   NebulaBox     │
│   CLI Tool      │◄──►│   Dashboard     │◄──►│   Registry      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   containerd     │    │   REST API      │    │   Image Store   │
│   Runtime        │    │   Server        │    │   & Auth        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Contributing

This project is in active development. Check the issues for ways to contribute!

## License

MIT License - see LICENSE file for details.
