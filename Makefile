# NebulaBox Makefile

.PHONY: build clean test run install dev deps

# Build the NebulaBox CLI
build:
	@echo "🔨 Building NebulaBox..."
	go build -o nebulabox ./cmd/nebulabox
	@echo "✅ Build complete!"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f nebulabox
	go clean
	@echo "✅ Clean complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed!"

# Development mode - build and run
dev: build
	@echo "🚀 Running NebulaBox in dev mode..."
	./nebulabox --help

# Install NebulaBox to system
install: build
	@echo "📥 Installing NebulaBox..."
	sudo cp nebulabox /usr/local/bin/
	@echo "✅ NebulaBox installed to /usr/local/bin/"

# Run with example command
demo: build
	@echo "🎬 Running demo..."
	./nebulabox run nginx
	@echo "📋 Listing containers..."
	./nebulabox list
	@echo "ℹ️  Showing version..."
	./nebulabox version

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run
	@echo "✅ Linting complete!"

# Cross-compile for multiple platforms
cross-compile:
	@echo "🌍 Cross-compiling..."
	GOOS=linux GOARCH=amd64 go build -o nebulabox-linux-amd64 ./cmd/nebulabox
	GOOS=windows GOARCH=amd64 go build -o nebulabox-windows-amd64.exe ./cmd/nebulabox
	GOOS=darwin GOARCH=amd64 go build -o nebulabox-darwin-amd64 ./cmd/nebulabox
	@echo "✅ Cross-compilation complete!"

# Setup development environment
setup: deps
	@echo "🛠️  Setting up development environment..."
	@echo "📋 Prerequisites:"
	@echo "  - Go 1.22+ installed"
	@echo "  - containerd installed (brew install containerd)"
	@echo "  - Docker (optional, for testing)"
	@echo ""
	@echo "🎯 Next steps:"
	@echo "  1. Run 'make demo' to test the CLI"
	@echo "  2. Check out the README.md for development phases"
	@echo "  3. Start implementing containerd integration"
	@echo "✅ Setup complete!"
