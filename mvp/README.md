# MVP Testing - NebulaBox

This folder contains all MVP testing materials, including the test application and testing scripts.

## ⚠️ IMPORTANT: Testing Concept

**We are testing NebulaBox by using NebulaBox itself!**

- ✅ Use **NebulaBox Build Specification** (not Dockerfile)
- ✅ Use **NebulaBox API** (not docker CLI)
- ✅ Use **NebulaBox Dashboard** (not Docker Desktop)
- ✅ Use **NebulaBox Runtime** (not direct Docker daemon)

**We build and run our MVP using NebulaBox to prove NebulaBox works!**

See `CONCEPT.md` for full details.

## Structure

```
mvp/
├── README.md                 # This file
├── CONCEPT.md                # ⭐ Testing concept (read this first!)
├── APPROACH.md               # Detailed approach explanation
├── BUILD_AND_DEPLOY.md       # Build and deployment guide
├── app/                      # Test application code
│   └── single-container/    # Single container MERN stack
│       ├── buildspec.json   # ⭐ NebulaBox Build Specification
│       └── app/
│           ├── frontend/     # React application
│           ├── backend/      # Express API
│           └── scripts/
│               └── start.sh  # Start script
├── scripts/                  # Testing and deployment scripts
│   └── deploy-buildspec.sh   # Automated deployment (NebulaBox API)
└── docs/                     # MVP-specific documentation
    └── test-results/          # Test results will go here
```

## Quick Start

### 1. Understand the Concept
Read `CONCEPT.md` to understand we're using NebulaBox, not Docker.

### 2. Prepare Application
```bash
# Application code goes here (we'll create it)
mvp/app/single-container/app/
```

### 3. Build Using NebulaBox
```bash
# Validate buildspec
curl -X POST http://localhost:8081/api/buildspec/validate \
  -H "Content-Type: application/json" \
  -d @mvp/app/single-container/buildspec.json

# Build via NebulaBox Dashboard
# Go to: Build Spec page → Upload buildspec.json → Build
```

### 4. Deploy Using NebulaBox
```bash
# Create container via API
curl -X POST http://localhost:8081/api/containers/run \
  -H "Content-Type: application/json" \
  -d '{
    "image": "mern-mvp:latest",
    "name": "mern-mvp",
    "ports": ["3000:3000", "5000:5000"]
  }'

# Or use dashboard: Containers → Create Container
```

### 5. Test Everything
Monitor, manage, and test entirely through NebulaBox Dashboard and API.

## Testing Phases

See `../docs/MVP_TESTING_PLAN.md` for complete 12-phase testing plan.

## Key Files

- `buildspec.json` - NebulaBox Build Specification (replaces Dockerfile)
- `CONCEPT.md` - Core testing concept (read this!)
- `APPROACH.md` - Detailed approach explanation
- `BUILD_AND_DEPLOY.md` - Step-by-step deployment guide

## Notes

- We use **NebulaBox Build Specification**, not Dockerfile
- All operations via **NebulaBox API** or Dashboard
- No Docker CLI commands in workflow
- We're proving NebulaBox can replace Docker!
