# NebulaBox Unified Architecture - Complete Picture

## Overview

NebulaBox unifies the entire development-to-deployment lifecycle into a single, cohesive system. This document explains how all components work together to solve the fragmentation problem.

---

## Problem Statement (What We're Solving)

### Current Fragmented Approach:
```
1. Local Dev: Frontend (port 3000) + Backend (port 5000) + DB (port 27017) - Manual setup
2. Remote Team: Can't access local dev, need VPN/tunnels
3. Deployment: Frontend → Vercel, Backend → Render, DB → MongoDB Atlas - 3 different places
4. Environment Management: Different configs for dev/staging/prod
5. Collaboration: Git conflicts, manual sync, no real-time sharing
```

### NebulaBox Unified Approach:
```
1. Local Dev: Single container with all services (unified)
2. Remote Team: Shared workspace with real-time sync
3. Deployment: Single deployment to NebulaBox Cloud (all services together)
4. Environment Management: Same BuildSpec, different environments
5. Collaboration: Real-time file sync, shared terminals, live preview
```

---

## Architecture Layers

### Layer 1: Development (Local)

#### 1.1 Single Container Architecture

**MERN Todo App Example:**

```
┌─────────────────────────────────────────────────────────┐
│              NebulaBox Container                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Container Namespace (Isolated)                  │  │
│  │                                                   │  │
│  │  ┌──────────────┐  ┌──────────────┐            │  │
│  │  │   MongoDB    │  │   Express    │            │  │
│  │  │   :27017     │  │   :5000      │            │  │
│  │  └──────────────┘  └──────────────┘            │  │
│  │                                                   │  │
│  │  ┌──────────────┐  ┌──────────────┐            │  │
│  │  │   Nginx      │  │   Tailscale  │            │  │
│  │  │   :80        │  │   (tunnel)   │            │  │
│  │  └──────────────┘  └──────────────┘            │  │
│  │                                                   │  │
│  │  All services in ONE container                   │  │
│  │  Port mapping: 3000:80, 5000:5000, 27017:27017 │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

**How It Works:**

1. **BuildSpec Definition** (`buildspec.json`):
   ```json
   {
     "name": "mern-todo",
     "base": { "image": "node:18-alpine" },
     "steps": [
       { "type": "run", "command": "apk add mongodb nginx tailscale" },
       { "type": "copy", "source": "backend", "dest": "/app/backend" },
       { "type": "copy", "source": "frontend", "dest": "/app/frontend" }
     ],
     "cmd": "/app/scripts/start.sh"
   }
   ```

2. **Startup Script** (`start.sh`):
   ```bash
   # Start MongoDB
   mongod --dbpath /app/data --fork
   
   # Start Backend API
   cd /app/backend && node server.js &
   
   # Start Nginx (serves React frontend)
   nginx -g "daemon off;" &
   
   # Start Tailscale (if needed for remote access)
   tailscaled &
   ```

3. **Single Container = All Services**:
   - All processes run in same container namespace
   - Shared filesystem
   - Internal networking (localhost)
   - Port mapping exposes services to host

**Benefits:**
- ✅ One command: `nebulabox run mern-todo:latest`
- ✅ All services start together
- ✅ No manual port management
- ✅ Consistent environment

---

### Layer 2: Grouped/Hierarchical Containers

For complex apps, NebulaBox supports container groups:

```
┌─────────────────────────────────────────────────────────┐
│              Container Group: "mern-app"                │
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │  Container 1 │  │  Container 2 │  │  Container 3 │ │
│  │  Frontend    │  │  Backend    │  │  Database   │ │
│  │  (Nginx)     │  │  (Express)   │  │  (MongoDB)  │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │
│         │                 │                  │          │
│         └─────────────────┼──────────────────┘          │
│                           │                             │
│                  Internal Network                       │
│                  (Isolated Bridge)                      │
└─────────────────────────────────────────────────────────┘
```

**Use Cases:**
- Microservices architecture
- Separate scaling per service
- Independent updates
- Service isolation

**NebulaBox Command:**
```bash
# Create container group
nebulabox group create mern-app

# Add containers to group
nebulabox group add mern-app frontend nginx:latest
nebulabox group add mern-app backend express:latest
nebulabox group add mern-app database mongo:latest

# Start entire group
nebulabox group start mern-app
```

---

### Layer 3: Remote Team Collaboration

#### 3.1 Shared Workspace Architecture

```
┌─────────────────────────────────────────────────────────┐
│              NebulaBox Shared Workspace                │
│                                                         │
│  Developer A (Local)          Developer B (Remote)      │
│  ┌──────────────┐            ┌──────────────┐          │
│  │  Container   │            │  Container   │          │
│  │  (Running)   │◄───────────┤  (Shared)    │          │
│  └──────────────┘            └──────────────┘          │
│         │                          │                    │
│         └──────────┬───────────────┘                    │
│                    │                                     │
│         ┌──────────▼──────────┐                         │
│         │  NebulaBox Sync     │                         │
│         │  (CRDT + WebSocket) │                         │
│         └─────────────────────┘                         │
│                    │                                     │
│         ┌──────────▼──────────┐                         │
│         │  File Sync Manager  │                         │
│         │  - Real-time sync    │                         │
│         │  - Conflict resolution│                       │
│         └─────────────────────┘                         │
└─────────────────────────────────────────────────────────┘
```

**How It Works:**

1. **Workspace Creation**:
   ```bash
   # Developer A creates shared workspace
   nebulabox workspace create mern-todo-dev --share
   ```

2. **Team Invitation**:
   ```bash
   # Developer A invites Developer B
   nebulabox workspace invite mern-todo-dev user@example.com
   ```

3. **Real-Time Collaboration**:
   - **File Sync**: Changes sync in real-time via CRDT
   - **Terminal Sharing**: Multiple users can access same container
   - **Live Preview**: All see same running app
   - **Conflict Resolution**: Automatic merge for code conflicts

4. **Access Methods**:
   - **WebSocket**: Real-time connection
   - **SSH Tunnel**: Direct access via Tailscale/VPN
   - **Web Dashboard**: Browser-based access

**Benefits:**
- ✅ Multiple developers work simultaneously
- ✅ Real-time code sync
- ✅ Shared terminal access
- ✅ Live preview for all team members
- ✅ No VPN setup needed

---

### Layer 4: Unified Deployment

#### 4.1 Single Deployment Target

**Current Problem:**
```
Frontend → Vercel (different config)
Backend → Render (different config)
Database → MongoDB Atlas (separate service)
Result: 3 different places, 3 different configs, fragmented
```

**NebulaBox Solution:**
```
All Services → NebulaBox Cloud (single deployment)
Result: One place, one config, unified access
```

#### 4.2 Deployment Architecture

```
┌─────────────────────────────────────────────────────────┐
│              NebulaBox Cloud Platform                  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │  Deployment: mern-todo-prod                      │  │
│  │                                                   │  │
│  │  ┌──────────────┐  ┌──────────────┐            │  │
│  │  │  Container   │  │  Container   │            │  │
│  │  │  (All in 1)  │  │  (Grouped)   │            │  │
│  │  │  - MongoDB   │  │  - Frontend   │            │  │
│  │  │  - Backend  │  │  - Backend    │            │  │
│  │  │  - Frontend │  │  - Database   │            │  │
│  │  │  - Nginx    │  │  - Services   │            │  │
│  │  └──────────────┘  └──────────────┘            │  │
│  │                                                   │  │
│  │  URL: https://mern-todo.nebulabox.cloud          │  │
│  │  All services accessible via single domain       │  │
│  └─────────────────────────────────────────────────┘  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │  Deployment: mern-todo-staging                  │  │
│  │  (Same BuildSpec, different environment)         │  │
│  └─────────────────────────────────────────────────┘  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │  Deployment: mern-todo-dev                      │  │
│  │  (Development environment)                       │  │
│  └─────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

**Deployment Flow:**

1. **Single BuildSpec** (same for all environments):
   ```json
   {
     "name": "mern-todo",
     "env": {
       "NODE_ENV": "${ENV}",  // Injected per environment
       "MONGODB_URI": "${MONGODB_URI}"
     }
   }
   ```

2. **Deploy Command**:
   ```bash
   # Deploy to production
   nebulabox deploy --env prod

   # Deploy to staging
   nebulabox deploy --env staging

   # Deploy to dev
   nebulabox deploy --env dev
   ```

3. **Environment Variables** (managed per environment):
   - Production: `MONGODB_URI=prod-db-url`
   - Staging: `MONGODB_URI=staging-db-url`
   - Dev: `MONGODB_URI=dev-db-url`

**Benefits:**
- ✅ Single deployment target
- ✅ Same BuildSpec for all environments
- ✅ Environment-specific configs
- ✅ Unified access (one URL)
- ✅ Easy rollback

---

### Layer 5: Remote Access & Multi-Environment

#### 5.1 Remote Team Access

```
┌─────────────────────────────────────────────────────────┐
│              Remote Team Access                         │
│                                                         │
│  Developer (US)    Developer (EU)    Developer (Asia)   │
│      │                  │                  │            │
│      └──────────────────┼──────────────────┘            │
│                         │                                │
│              ┌──────────▼──────────┐                     │
│              │  NebulaBox Cloud    │                     │
│              │  (Single Platform)  │                     │
│              └──────────┬──────────┘                     │
│                         │                                │
│              ┌──────────▼──────────┐                     │
│              │  Deployment         │                     │
│              │  (All Services)     │                     │
│              └─────────────────────┘                     │
│                                                         │
│  All developers access same deployment                  │
│  Real-time collaboration enabled                        │
└─────────────────────────────────────────────────────────┘
```

#### 5.2 Environment Management

```
┌─────────────────────────────────────────────────────────┐
│              Environment Management                     │
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Production │  │   Staging    │  │   Development│  │
│  │              │  │              │  │              │  │
│  │  URL:        │  │  URL:        │  │  URL:        │  │
│  │  prod.app    │  │  staging.app │  │  dev.app     │  │
│  │              │  │              │  │              │  │
│  │  DB: prod    │  │  DB: staging │  │  DB: dev     │  │
│  │  Config: A   │  │  Config: B   │  │  Config: C   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                         │
│  Same BuildSpec, Different Environments                 │
└─────────────────────────────────────────────────────────┘
```

**Environment Configuration:**
```yaml
# nebulabox-env.yml
environments:
  production:
    domain: mern-todo.com
    mongodb_uri: mongodb://prod-cluster
    node_env: production
    replicas: 3
    
  staging:
    domain: staging.mern-todo.com
    mongodb_uri: mongodb://staging-cluster
    node_env: staging
    replicas: 1
    
  development:
    domain: dev.mern-todo.com
    mongodb_uri: mongodb://dev-cluster
    node_env: development
    replicas: 1
```

---

## Complete Workflow: Development → Deployment

### Step 1: Local Development
```bash
# Clone repo
git clone https://github.com/user/mern-todo.git
cd mern-todo

# Build and run locally (all services in one container)
nebulabox build -f buildspec.json -t mern-todo:latest
nebulabox run mern-todo:latest --name mern-dev -p 3000:80 -p 5000:5000

# Access: http://localhost:3000
# All services running in single container
```

### Step 2: Share with Team
```bash
# Create shared workspace
nebulabox workspace create mern-todo-dev --share

# Invite team members
nebulabox workspace invite mern-todo-dev team@example.com

# Team members join
nebulabox workspace join mern-todo-dev

# Real-time collaboration enabled
# - Code changes sync automatically
# - Shared terminal access
# - Live preview for all
```

### Step 3: Deploy to Cloud
```bash
# Deploy to development environment
nebulabox deploy --env dev

# Deploy to staging
nebulabox deploy --env staging

# Deploy to production
nebulabox deploy --env prod
```

### Step 4: Remote Team Access
```bash
# All team members access same deployment
# - Production: https://mern-todo.nebulabox.cloud
# - Staging: https://staging.mern-todo.nebulabox.cloud
# - Dev: https://dev.mern-todo.nebulabox.cloud

# Real-time collaboration continues in cloud
# - Code changes sync
# - Shared debugging
# - Live monitoring
```

---

## Key Architectural Principles

### 1. **Unified Containerization**
- All services in one container (or grouped containers)
- Single BuildSpec defines everything
- No manual service orchestration

### 2. **Real-Time Collaboration**
- CRDT-based file sync
- WebSocket for real-time updates
- Shared terminal and preview

### 3. **Single Deployment Target**
- One platform (NebulaBox Cloud)
- All services deployed together
- Unified access and management

### 4. **Environment Abstraction**
- Same BuildSpec for all environments
- Environment-specific configs
- Easy environment switching

### 5. **Remote-First Design**
- Built for distributed teams
- No VPN required
- Secure access via WebSocket/SSH

---

## Comparison: Before vs After

### Before (Fragmented):
```
Local Dev: Manual setup (3 services, 3 ports)
Remote Access: VPN/Tailscale setup
Deployment: Vercel + Render + MongoDB Atlas (3 places)
Collaboration: Git only, no real-time
Environment: Different configs per platform
```

### After (Unified):
```
Local Dev: Single container (1 command)
Remote Access: Built-in workspace sharing
Deployment: NebulaBox Cloud (1 place)
Collaboration: Real-time sync + shared preview
Environment: Same BuildSpec, different configs
```

---

## Technical Stack

### Container Runtime
- Custom NebulaBox engine (Linux namespaces + cgroups)
- No Docker/containerd dependency
- Direct kernel integration

### Collaboration
- CRDT for conflict-free file sync
- WebSocket for real-time updates
- SSH for direct access

### Deployment
- NebulaBox Cloud (Kubernetes-like orchestration)
- Automatic scaling
- Load balancing
- Health monitoring

### Networking
- Internal container networking
- Port mapping to host
- Custom domains per environment
- SSL/TLS certificates

---

## Next Steps

1. **Review Architecture**: Confirm this matches your vision
2. **Update Plan**: Incorporate all 7 requirements
3. **Implementation**: Start with core engine, then add layers

This architecture unifies everything into a single, cohesive system that solves the fragmentation problem completely.

