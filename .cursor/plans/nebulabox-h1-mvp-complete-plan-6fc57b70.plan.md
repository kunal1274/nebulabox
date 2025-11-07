<!-- 6fc57b70-2782-4d8e-b75f-ebb6eebe514d 6673b092-b493-448e-94c0-850f0bcf73f7 -->
# NebulaBox Custom Container Engine - Complete Independence Plan

## Phase 1: Custom Container Runtime Engine (Foundation)

### 1.1 Core Runtime Architecture

**Location**: `internal/engine/` (new module)

Create a completely independent container runtime that uses Linux primitives directly:

- **`engine/runtime.go`**: Main runtime interface and implementation
  - Direct Linux namespace management (PID, Network, Mount, UTS, IPC, User)
  - cgroups v2 integration for resource limits
  - OCI-compliant container spec generation
  - Process isolation without containerd dependency
  - Flexible container grouping support

- **`engine/namespace.go`**: Linux namespace management
  - `CreateNamespace()`: Create isolated namespaces
  - `EnterNamespace()`: Switch process to namespace
  - `CloneWithNamespaces()`: Fork process in new namespaces
  - Support for shared namespaces (for grouped containers)

- **`engine/cgroups.go`**: Resource management
  - CPU limits (cpuset, cpu)
  - Memory limits (memory)
  - I/O limits (blkio)
  - Network bandwidth (net_cls)
  - Group-level resource limits

- **`engine/filesystem.go`**: Container filesystem
  - Root filesystem creation (chroot/overlay)
  - Image layer mounting
  - Volume management
  - Bind mounts
  - Shared volumes between grouped containers

- **`engine/network.go`**: Network isolation
  - Virtual network interfaces (veth pairs)
  - Bridge network creation
  - Port forwarding (iptables/nftables)
  - DNS configuration
  - Internal networking for container groups
  - Service discovery within groups

- **`engine/process.go`**: Process management
  - Container process lifecycle
  - Signal handling
  - Log collection (stdout/stderr)
  - Health checks
  - Process monitoring across groups

### 1.2 Image Management System

**Location**: `internal/engine/images.go`

- **Image format**: Custom NebulaBox image format (compatible with OCI)
  - Layer-based storage (similar to Docker)
  - Image manifest and config
  - Layer deduplication
  - Image registry support

- **Build system**: Direct image building
  - BuildSpec → Image layers
  - Layer caching
  - Multi-stage builds
  - Image tagging and storage
  - Support for building individual service images

### 1.3 Container Lifecycle

**Location**: `internal/engine/container.go`

- Container creation, start, stop, delete
- State management (created, running, stopped, removed)
- Resource tracking
- Metadata storage
- Dependency management (start order for grouped containers)

### 1.4 Flexible Container Grouping System

**Location**: `internal/engine/group.go`

- **Group Definition**: Flexible grouping strategies
  - Single container (all services together)
  - Service-based grouping (frontend separate, backend+DB together)
  - Microservices grouping (each service in own container)
  - Custom grouping (user-defined relationships)

- **Group Management**:
  - `CreateGroup()`: Create container group with strategy
  - `AddContainerToGroup()`: Add container to existing group
  - `RemoveContainerFromGroup()`: Remove container from group
  - `StartGroup()`: Start all containers in group (with dependencies)
  - `StopGroup()`: Stop all containers in group
  - `GetGroupStatus()`: Get status of all containers in group

- **Group Strategies**:
  - **Monolithic**: All services in one container
  - **Frontend-Backend**: Frontend separate, backend+DB together
  - **Three-Tier**: Frontend, Backend, Database separate
  - **Microservices**: Each service in own container
  - **Custom**: User-defined grouping via config

- **Group Networking**:
  - Internal bridge network per group
  - Service discovery within group
  - Port mapping from group to host
  - Inter-group communication

- **Group Configuration**:
  ```json
  {
    "name": "mern-app",
    "strategy": "frontend-backend",
    "containers": [
      {
        "name": "frontend",
        "image": "mern-frontend:latest",
        "ports": ["3000:80"]
      },
      {
        "name": "backend-db",
        "image": "mern-backend-db:latest",
        "ports": ["5000:5000", "27017:27017"]
      }
    ],
    "networking": {
      "internal": true,
      "bridge": "mern-bridge"
    }
  }
  ```


## Phase 2: CLI Workflow (Primary Interface)

### 2.1 Enhanced CLI Commands

**Location**: `internal/cli/` (extend existing)

- **`cli/build.go`**: Build images from BuildSpec
  ```go
  nebulabox build -f buildspec.json -t mern-todo:latest
  ```

- **`cli/run.go`**: Run containers (enhance existing)
  ```go
  nebulabox run mern-todo:latest --name todo-app -p 3000:80 -p 5000:5000
  ```

- **`cli/ps.go`**: List containers
  ```go
  nebulabox ps -a
  ```

- **`cli/exec.go`**: Execute commands in containers
  ```go
  nebulabox exec todo-app /bin/sh
  ```

- **`cli/logs.go`**: View container logs (enhance existing)
  ```go
  nebulabox logs todo-app -f
  ```

- **`cli/images.go`**: Image management (enhance existing)
  ```go
  nebulabox images
  nebulabox rmi mern-todo:latest
  ```


### 2.2 CLI Workflow Guide

**Location**: `docs/CLI_WORKFLOW_GUIDE.md`

Comprehensive guide covering:

- Installation and setup
- Basic commands (build, run, stop, logs)
- BuildSpec format and examples
- Multi-component app deployment (MERN example)
- Remote deployment via CLI
- Troubleshooting

### 2.3 CLI Integration with Engine

**Location**: `internal/cli/engine_client.go`

- Direct integration with custom engine (no API dependency)
- Local-first workflow
- Fallback to API for remote operations

## Phase 3: API Development (Parallel)

### 3.1 Engine API Layer

**Location**: `internal/api/engine/` (new module)

- **`engine/handlers.go`**: HTTP handlers for engine operations
  - `POST /api/v2/containers/run`: Create and run container
  - `GET /api/v2/containers`: List containers
  - `POST /api/v2/containers/:id/start`: Start container
  - `POST /api/v2/containers/:id/stop`: Stop container
  - `GET /api/v2/containers/:id/logs`: Get logs
  - `POST /api/v2/containers/:id/exec`: Execute command

- **`engine/images.go`**: Image API endpoints
  - `POST /api/v2/images/build`: Build from BuildSpec
  - `GET /api/v2/images`: List images
  - `DELETE /api/v2/images/:id`: Delete image

### 3.2 API Server Integration

**Location**: `internal/api/server.go`

- Add engine routes alongside existing routes
- Version routing (`/api/v2/` for new engine, `/api/` for legacy)
- Feature flag for engine selection

## Phase 4: Remote Deployment & Built-in Tunneling (No External Tools)

### 4.1 Built-in Tunneling System (No ngrok/Tailscale Needed)

**Location**: `internal/remote/tunnel/`

- **`tunnel/server.go`**: NebulaBox native tunneling server
  - WebSocket-based secure tunnels
  - Automatic port forwarding
  - Connection multiplexing
  - No external dependencies (ngrok/Tailscale not required)

- **`tunnel/client.go`**: Tunnel client
  - Automatic tunnel establishment
  - Reconnection logic
  - Connection pooling
  - Health monitoring

- **`tunnel/protocol.go`**: Tunnel protocol
  - Secure tunnel establishment
  - Data encryption
  - Authentication
  - Heartbeat mechanism

- **Features**:
  - Zero-config tunneling (automatic)
  - Secure by default (encrypted)
  - Works behind NAT/firewalls
  - No port forwarding needed on router
  - Direct container access from anywhere

### 4.2 WebSocket Server

**Location**: `internal/remote/websocket/`

- **`websocket/server.go`**: WebSocket server implementation
  - Real-time connection management
  - Container lifecycle events
  - Log streaming
  - Exec command streaming
  - Built-in tunneling support

- **`websocket/handlers.go`**: Message handlers
  - Container operations (create, start, stop)
  - Log streaming requests
  - Exec command requests
  - Health check streaming
  - Tunnel establishment

- **`websocket/protocol.go`**: Message protocol
  - Request/response format
  - Event streaming format
  - Error handling
  - Tunnel protocol messages

### 4.3 SSH Agent (Hybrid)

**Location**: `internal/remote/ssh/`

- **`ssh/agent.go`**: SSH agent for direct server access
  - Key-based authentication
  - Command execution
  - File transfer (SCP/SFTP)
  - Port forwarding
  - Built-in SSH tunnel (no external setup)

- **`ssh/tunnel.go`**: SSH tunnel management
  - Secure port forwarding
  - Reverse tunnels
  - Connection pooling
  - Automatic tunnel setup

### 4.4 Remote Client

**Location**: `internal/remote/client.go`

- Unified interface for remote operations
- WebSocket for real-time operations
- SSH for direct access
- Built-in tunneling (no ngrok/Tailscale)
- Automatic fallback and retry
- Zero-config remote access

### 4.5 CLI Remote Commands

**Location**: `internal/cli/remote.go`

```bash
# Connect to remote (automatic tunneling, no config needed)
nebulabox remote connect user@remote-server

# Deploy to remote (tunnel established automatically)
nebulabox remote deploy buildspec.json --target remote-server

# Access remote container (via built-in tunnel)
nebulabox remote exec container-id --target remote-server

# List remote containers (via tunnel)
nebulabox remote ps --target remote-server
```

## Phase 5: GitHub Integration & One-Click Deployment

### 5.1 GitHub App Integration

**Location**: `internal/github/`

- **`github/app.go`**: GitHub App OAuth and installation
  - App registration and authentication
  - Repository access permissions
  - Installation webhook handling

- **`github/webhook.go`**: Enhanced webhook handling (extend existing)
  - Push events (auto-deploy on push to main)
  - Pull request events (preview deployments)
  - Release events (production deployments)
  - Workflow events (GitHub Actions integration)

- **`github/actions.go`**: GitHub Actions support
  - NebulaBox Action for CI/CD
  - BuildSpec validation in Actions
  - Automated testing and deployment

### 5.2 One-Click Deployment

**Location**: `internal/deploy/github.go`

- **Repository Detection**: Auto-detect BuildSpec in repo
- **Deploy Button**: GitHub "Deploy to NebulaBox" button
  ```markdown
  [![Deploy to NebulaBox](https://nebulabox.cloud/deploy-button.svg)](https://nebulabox.cloud/deploy?repo=owner/repo)
  ```

- **Deployment Flow**:

  1. User clicks "Deploy to NebulaBox" in GitHub
  2. NebulaBox clones repository
  3. Finds `buildspec.json` or `nebulabox.json`
  4. Builds and deploys automatically
  5. Returns deployment URL

- **Deployment Status**: GitHub commit status integration
  - Show deployment progress in PRs
  - Deployment success/failure status
  - Preview URL in PR comments

### 5.3 BuildSpec as Universal Config

**Location**: `internal/deploy/universal.go`

- **Multi-Platform Support**: Same BuildSpec works for:
  - NebulaBox Cloud
  - Vercel (via adapter)
  - Railway (via adapter)
  - Render (via adapter)
  - Self-hosted NebulaBox

- **Platform Adapters**:
  - `internal/deploy/adapters/vercel.go`: Convert BuildSpec to Vercel config
  - `internal/deploy/adapters/railway.go`: Convert BuildSpec to Railway config
  - `internal/deploy/adapters/render.go`: Convert BuildSpec to Render config

## Phase 6: NebulaBox Cloud Platform

### 6.1 Cloud Infrastructure

**Location**: `internal/cloud/platform/`

- **`platform/server.go`**: Cloud platform API server
  - Multi-tenant architecture
  - Resource quotas and limits
  - Billing integration (future)
  - User management

- **`platform/deployment.go`**: Cloud deployment engine
  - Automatic scaling
  - Load balancing
  - Health monitoring
  - Auto-restart on failure

- **`platform/networking.go`**: Cloud networking
  - Custom domains
  - SSL/TLS certificates (Let's Encrypt)
  - CDN integration
  - DDoS protection

### 6.2 Cloud Dashboard

**Location**: `web/cloud-dashboard/` (new React app)

- **Project Management**: Create, list, delete projects
- **Deployment History**: View past deployments
- **Logs & Metrics**: Real-time logs and monitoring
- **Settings**: Domain, environment variables, scaling
- **Billing**: Usage and billing (future)

### 6.3 Cloud CLI

**Location**: `internal/cli/cloud.go`

```bash
# Login to NebulaBox Cloud
nebulabox cloud login

# Deploy from current directory
nebulabox cloud deploy

# Deploy from GitHub
nebulabox cloud deploy --github owner/repo

# List deployments
nebulabox cloud deployments

# View logs
nebulabox cloud logs <deployment-id>
```

### 6.4 Multi-Platform Deployment

**Location**: `internal/deploy/multi_platform.go`

- **Unified Deployment Interface**:
  ```bash
  # Deploy to NebulaBox Cloud
  nebulabox deploy --platform nebulabox
  
  # Deploy to Vercel
  nebulabox deploy --platform vercel
  
  # Deploy to Railway
  nebulabox deploy --platform railway
  
  # Deploy to all platforms
  nebulabox deploy --platform all
  ```

- **Platform Detection**: Auto-detect platform from environment
- **Config Conversion**: Convert BuildSpec to platform-specific configs

## Phase 7: MERN Todo App - GitHub Deployment

### 7.1 GitHub Repository Setup

**Location**: `mvp/app/mern-todo/`

- **`.github/workflows/nebulabox.yml`**: GitHub Actions workflow
  - Auto-deploy on push to main
  - Preview deployments for PRs
  - Build and test before deploy

- **`nebulabox.json`**: Deployment configuration
  - Platform selection (nebulabox, vercel, railway)
  - Environment variables
  - Domain configuration
  - Scaling settings

- **`README.md`**: Updated with deployment instructions
  - One-click deployment button
  - Manual deployment steps
  - Platform options

### 7.2 Deployment Templates

**Location**: `templates/deployment/`

- **GitHub Template**: Repository template with NebulaBox setup
- **Vercel Adapter Template**: Vercel-specific config
- **Railway Adapter Template**: Railway-specific config
- **Multi-Platform Template**: Works with all platforms

### 7.3 Testing

- Local deployment test
- GitHub one-click deployment test
- NebulaBox Cloud deployment test
- Multi-platform deployment test (Vercel, Railway)
- Remote deployment test (WebSocket)
- Remote deployment test (SSH)
- Multi-component verification

## Phase 6: Documentation

### 6.1 Architecture Documentation

**Location**: `docs/ARCHITECTURE.md`

- Custom engine architecture
- Namespace and cgroup usage
- Network isolation
- Filesystem management
- Image format specification

### 6.2 Developer Guide

**Location**: `docs/DEVELOPER_GUIDE.md`

- Building the engine
- Adding new features
- Testing procedures
- Debugging guide

### 6.3 Migration Guide

**Location**: `docs/MIGRATION_GUIDE.md`

- Moving from containerd to custom engine
- API migration (v1 to v2)
- Feature parity checklist

## Implementation Order

1. **Week 1-2**: Core runtime engine (namespaces, cgroups, filesystem)
2. **Week 2-3**: CLI workflow (build, run, ps, logs) + CLI guide
3. **Week 3-4**: API development (parallel with CLI)
4. **Week 4-5**: Image management and build system
5. **Week 5-6**: Remote deployment (WebSocket + SSH)
6. **Week 6-7**: MERN todo app integration and testing
7. **Week 7**: Documentation and polish

## Security (Post-MVP)

Security implementation deferred until after MVP:

- Authentication/authorization
- TLS/SSL encryption
- Network policies
- Container sandboxing
- Audit logging

## Project Organization Structure

### Module Organization

```
nebulabox/
├── internal/
│   ├── engine/              # Custom container engine (NEW - completely independent)
│   │   ├── runtime.go
│   │   ├── namespace.go
│   │   ├── cgroups.go
│   │   ├── filesystem.go
│   │   ├── network.go
│   │   ├── process.go
│   │   ├── images.go
│   │   ├── container.go
│   │   └── group.go          # Container grouping system
│   │
│   ├── cli/                  # CLI commands (enhance existing)
│   │   ├── build.go
│   │   ├── run.go
│   │   ├── ps.go
│   │   ├── exec.go
│   │   ├── logs.go
│   │   ├── images.go
│   │   ├── engine_client.go  # Direct engine integration
│   │   ├── remote.go         # Remote deployment commands
│   │   └── cloud.go          # Cloud deployment commands
│   │
│   ├── api/
│   │   ├── engine/           # Engine API layer (NEW)
│   │   │   ├── handlers.go
│   │   │   └── images.go
│   │   └── server.go         # Add v2 routes
│   │
│   ├── remote/               # Remote deployment (NEW)
│   │   ├── websocket/
│   │   │   ├── server.go
│   │   │   ├── handlers.go
│   │   │   └── protocol.go
│   │   ├── ssh/
│   │   │   ├── agent.go
│   │   │   └── tunnel.go
│   │   ├── tunnel/           # Built-in tunneling (NEW)
│   │   │   ├── server.go
│   │   │   ├── client.go
│   │   │   └── protocol.go
│   │   └── client.go
│   │
│   ├── deploy/               # Deployment system (NEW)
│   │   ├── github.go
│   │   ├── universal.go
│   │   ├── multi_platform.go
│   │   └── adapters/
│   │       ├── vercel.go
│   │       ├── railway.go
│   │       └── render.go
│   │
│   ├── github/               # GitHub integration (NEW)
│   │   ├── app.go
│   │   ├── webhook.go
│   │   └── actions.go
│   │
│   ├── cloud/                # Cloud platform (NEW)
│   │   └── platform/
│   │       ├── server.go
│   │       ├── deployment.go
│   │       └── networking.go
│   │
│   ├── orchestrator/          # Cluster orchestration (enhance existing)
│   │   ├── node.go
│   │   ├── deployment.go
│   │   ├── health.go
│   │   └── cluster.go        # k3s/k8s/swarm support (NEW)
│   │
│   └── shareruntime/         # Shared workspaces (existing - enhance)
│       ├── workspace.go
│       ├── sync.go
│       └── tunnel.go
│
├── docs/                      # Documentation (separate folder)
│   ├── ARCHITECTURE.md
│   ├── CLI_WORKFLOW_GUIDE.md
│   ├── MERN_CLI_WORKFLOW.md
│   ├── UNIFIED_ARCHITECTURE.md
│   ├── DEPLOYMENT_GUIDE.md
│   ├── DEVELOPER_GUIDE.md
│   └── MIGRATION_GUIDE.md
│
├── scripts/                   # Scripts (identifiable, organized)
│   ├── engine/               # Engine-related scripts
│   │   ├── test-engine.sh
│   │   └── build-engine.sh
│   ├── cli/                  # CLI workflow test scripts
│   │   ├── workflow-01-build.sh
│   │   ├── workflow-02-run.sh
│   │   ├── workflow-03-group.sh
│   │   ├── workflow-04-remote.sh
│   │   └── workflow-05-deploy.sh
│   ├── testing/              # Testing scripts
│   │   ├── test-local.sh
│   │   ├── test-remote.sh
│   │   └── test-cluster.sh
│   ├── deployment/           # Deployment scripts
│   │   ├── deploy-local.sh
│   │   ├── deploy-github.sh
│   │   └── deploy-cloud.sh
│   └── setup/                # Setup scripts
│       ├── setup-engine.sh
│       └── setup-cli.sh
│
├── mvp/app/mern-todo/        # MERN todo app
│   ├── buildspec.json
│   ├── nebulabox.json        # Deployment config
│   └── .github/workflows/
│       └── nebulabox.yml
│
└── templates/                 # Templates
    └── deployment/
        ├── github-template/
        └── multi-platform-template/
```

## Key Files to Create/Modify

### New Engine Module Files (`internal/engine/`):

- `internal/engine/runtime.go` - Main runtime interface
- `internal/engine/namespace.go` - Linux namespace management
- `internal/engine/cgroups.go` - Resource management
- `internal/engine/filesystem.go` - Container filesystem
- `internal/engine/network.go` - Network isolation
- `internal/engine/process.go` - Process management
- `internal/engine/images.go` - Image management
- `internal/engine/container.go` - Container lifecycle
- `internal/engine/group.go` - Container grouping system

### New Remote Module Files (`internal/remote/`):

- `internal/remote/websocket/server.go`
- `internal/remote/websocket/handlers.go`
- `internal/remote/websocket/protocol.go`
- `internal/remote/ssh/agent.go`
- `internal/remote/ssh/tunnel.go`
- `internal/remote/tunnel/server.go` - Built-in tunneling
- `internal/remote/tunnel/client.go`
- `internal/remote/tunnel/protocol.go`
- `internal/remote/client.go`

### New Deployment Module Files (`internal/deploy/`):

- `internal/deploy/github.go`
- `internal/deploy/universal.go`
- `internal/deploy/multi_platform.go`
- `internal/deploy/adapters/vercel.go`
- `internal/deploy/adapters/railway.go`
- `internal/deploy/adapters/render.go`

### New GitHub Module Files (`internal/github/`):

- `internal/github/app.go`
- `internal/github/webhook.go` (enhance existing)
- `internal/github/actions.go`

### New Cloud Module Files (`internal/cloud/platform/`):

- `internal/cloud/platform/server.go`
- `internal/cloud/platform/deployment.go`
- `internal/cloud/platform/networking.go`

### New Orchestrator Files (`internal/orchestrator/`):

- `internal/orchestrator/cluster.go` - k3s/k8s/swarm support

### Enhanced CLI Files (`internal/cli/`):

- `internal/cli/build.go` (enhance)
- `internal/cli/run.go` (enhance)
- `internal/cli/ps.go` (new or enhance)
- `internal/cli/exec.go` (new or enhance)
- `internal/cli/logs.go` (enhance)
- `internal/cli/images.go` (enhance)
- `internal/cli/engine_client.go` (new)
- `internal/cli/remote.go` (new)
- `internal/cli/cloud.go` (new)
- `internal/cli/group.go` (new) - Container group commands

### New API Files (`internal/api/engine/`):

- `internal/api/engine/handlers.go`
- `internal/api/engine/images.go`

### Documentation Files (`docs/`):

- `docs/ARCHITECTURE.md` - Engine architecture
- `docs/CLI_WORKFLOW_GUIDE.md` - Complete CLI guide
- `docs/MERN_CLI_WORKFLOW.md` - MERN app workflow
- `docs/UNIFIED_ARCHITECTURE.md` - Unified architecture (already created)
- `docs/DEPLOYMENT_GUIDE.md` - Deployment guide
- `docs/DEVELOPER_GUIDE.md` - Developer guide
- `docs/MIGRATION_GUIDE.md` - Migration guide
- `docs/CLUSTER_GUIDE.md` - k3s/k8s/swarm deployment guide

### Workflow Test Scripts (`scripts/cli/`):

- `scripts/cli/workflow-01-build.sh` - Test build workflow
- `scripts/cli/workflow-02-run.sh` - Test run workflow
- `scripts/cli/workflow-03-group.sh` - Test container grouping
- `scripts/cli/workflow-04-remote.sh` - Test remote deployment
- `scripts/cli/workflow-05-deploy.sh` - Test cloud deployment
- `scripts/cli/workflow-06-mern-complete.sh` - Complete MERN workflow

### Testing Scripts (`scripts/testing/`):

- `scripts/testing/test-local.sh` - Local development testing
- `scripts/testing/test-remote.sh` - Remote collaboration testing
- `scripts/testing/test-cluster.sh` - Cluster deployment testing
- `scripts/testing/test-github.sh` - GitHub integration testing

### Deployment Scripts (`scripts/deployment/`):

- `scripts/deployment/deploy-local.sh` - Local deployment
- `scripts/deployment/deploy-github.sh` - GitHub deployment
- `scripts/deployment/deploy-cloud.sh` - Cloud deployment
- `scripts/deployment/deploy-cluster.sh` - Cluster deployment (k3s/k8s/swarm)

### Engine Scripts (`scripts/engine/`):

- `scripts/engine/test-engine.sh` - Test engine functionality
- `scripts/engine/build-engine.sh` - Build engine
- `scripts/engine/benchmark-engine.sh` - Performance benchmarks

### Modified Files:

- `internal/api/server.go` (add v2 routes, engine routes)
- `internal/cli/root.go` (add new commands)
- `mvp/app/mern-todo/buildspec.json` (verify compatibility)
- `mvp/app/mern-todo/nebulabox.json` (new - deployment config)
- `internal/orchestrator/deployment.go` (enhance for cluster support)