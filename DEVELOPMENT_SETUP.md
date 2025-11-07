# 🛠️ NebulaBox Development Setup

## Prerequisites

### **1. Install Go**
```bash
# Option 1: Snap (Recommended)
sudo snap install go --classic

# Option 2: APT
sudo apt update
sudo apt install golang-go

# Verify installation
go version
```

### **2. Install Node.js (for Dashboard)**
```bash
# Install Node.js 18+
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify installation
node --version
npm --version
```

### **3. Install containerd (Optional for Testing)**
```bash
# Install containerd
sudo apt update
sudo apt install containerd

# Start containerd service
sudo systemctl start containerd
sudo systemctl enable containerd
```

## 🚀 Quick Start

### **1. Clone and Setup**
```bash
git clone <repository-url>
cd nebulabox

# Install Go dependencies
go mod tidy
```

### **2. Build Everything**
```bash
# Build CLI tool
make build

# Build API server
make build-api

# Verify builds
ls -la nebulabox nebulabox-api
```

### **3. Start Development Environment**

**Option A: API Server Only**
```bash
# Start API server
make api
# API will be available at http://localhost:8080
```

**Option B: Full Stack (API + Dashboard)**
```bash
# Terminal 1: Start API server
make api

# Terminal 2: Start Dashboard
cd web/dashboard
npm install
npm run dev
# Dashboard will be available at http://localhost:3000
```

## 🧪 Testing the Setup

### **1. Test CLI Tool**
```bash
# Show help
./nebulabox --help

# Run demo
make demo

# Test individual commands
./nebulabox version
./nebulabox list
./nebulabox run nginx
```

### **2. Test API Server**
```bash
# Health check
curl http://localhost:8080/api/health

# List containers
curl http://localhost:8080/api/containers

# System stats
curl http://localhost:8080/api/system/stats
```

### **3. Test Dashboard**
1. Open http://localhost:3000 in browser
2. Navigate through different pages
3. Try creating a container
4. Check system monitoring

## 📁 Project Structure

```
nebulabox/
├── cmd/
│   ├── nebulabox/     # CLI tool entry point
│   └── api/           # API server entry point
├── internal/
│   ├── cli/           # CLI command implementations
│   ├── containerd/    # Container runtime client
│   └── api/           # REST API server
├── web/
│   └── dashboard/     # React dashboard
├── Makefile           # Build and run commands
└── README.md          # Project documentation
```

## 🔧 Development Commands

### **Build Commands**
```bash
make build          # Build CLI tool
make build-api      # Build API server
make clean          # Clean build artifacts
```

### **Run Commands**
```bash
make demo           # Run CLI demo
make api            # Start API server
make dev-full       # Start API + Dashboard
```

### **Development Commands**
```bash
make fmt            # Format Go code
make lint           # Lint Go code
make test           # Run tests
make deps           # Install dependencies
```

## 🐛 Troubleshooting

### **Go Not Found**
```bash
# Check if Go is installed
which go

# If not found, install via snap
sudo snap install go --classic

# Add to PATH if needed
export PATH=$PATH:/snap/bin
```

### **Node.js Issues**
```bash
# Check Node version
node --version

# If version is too old, install newer version
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### **Port Already in Use**
```bash
# Check what's using port 8080
sudo lsof -i :8080

# Kill process if needed
sudo kill -9 <PID>

# Or change port in server.go
```

### **CORS Issues**
- API server is configured for `http://localhost:3000`
- If using different port, update CORS config in `internal/api/server.go`

## 📚 Next Steps

1. **Test the complete setup** using the commands above
2. **Start developing** new features
3. **Check the roadmap** in `README.md`
4. **Join the development** by implementing Phase 2 features

## 🎯 Development Phases

- ✅ **Phase 1**: CLI tool with containerd backend
- ✅ **Phase 2**: Real containerd integration
- ✅ **Phase 3**: React dashboard with UI/UX
- ✅ **Phase 4**: REST API server
- 🔄 **Phase 5**: Real containerd integration (replace mock)
- 📋 **Phase 6+**: Advanced features (see roadmap)

## 💡 Tips

- Use `make dev-full` for full-stack development
- Check logs in terminal for debugging
- API server logs all requests and errors
- Dashboard uses mock data when API is unavailable
- All endpoints return JSON with consistent error format

---

**Happy coding!** 🚀
