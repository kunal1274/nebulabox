# NebulaBox Makefile

.PHONY: build build-api build-test build-registry registry-test clean run api registry test help

# Build targets
build: build-api build-test build-registry build-cli
	@echo "✅ All builds complete!"

build-api:
	@echo "🔨 Building NebulaBox API server..."
	go build -o nebulabox-api ./cmd/api
	@echo "✅ API server build complete!"

build-test:
	@echo "🔨 Building NebulaBox test runner..."
	go build -o nebulabox-test ./cmd/test
	@echo "✅ Test runner build complete!"

build-registry:
	@echo "🔨 Building Nebula Registry server..."
	go build -o nebulabox-registry ./cmd/registry
	@echo "✅ Registry server build complete!"

build-cli:
	@echo "🔨 Building NebulaBox CLI..."
	go build -o nebulabox-cli ./cmd/cli
	@echo "✅ CLI build complete!"

build-cli-test:
	@echo "🔨 Building NebulaBox CLI for testing..."
	@mkdir -p bin
	go build -o bin/nebulabox ./cmd/nebulabox
	@ln -sf nebulabox bin/nbx 2>/dev/null || true
	@chmod +x bin/nbx 2>/dev/null || true
	@echo "✅ CLI test binary build complete!"
	@echo "✅ Shortcut 'nbx' created!"

# Clean targets
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f nebulabox-api nebulabox-test nebulabox-registry nebulabox-cli
	rm -rf bin/
	rm -rf test-results/
	rm -rf registry-storage/
	@echo "✅ Clean complete!"

# Run targets
run: build-api
	@echo "🚀 Starting NebulaBox API server..."
	./nebulabox-api

api: build-api
	@echo "🚀 Starting NebulaBox API server..."
	./nebulabox-api

registry: build-registry
	@echo "🚀 Starting Nebula Registry server..."
	./nebulabox-registry

registry-test: build-registry
	@echo "🧪 Testing Nebula Registry..."
	@NEBULABOX_REGISTRY_PORT=5001 NEBULABOX_REGISTRY_STORAGE=./test-registry-storage ./nebulabox-registry & \
	 sleep 2 && bash scripts/test-registry.sh || true; \
	 pkill -f nebulabox-registry || true; \
	 rm -rf ./test-registry-storage

registry-test-unit:
	@echo "🧪 Running registry unit tests..."
	go test -v ./internal/registry/...

registry-test-integration: build-registry build-api
	@echo "🔗 Running registry integration tests..."
	@bash scripts/test-registry-integration.sh

registry-health:
	@echo "🏥 Checking registry health..."
	@curl -s http://localhost:$${NEBULABOX_REGISTRY_PORT:-5001}/v2/ | head -1 || echo "❌ Registry not responding"

test: build-test
	@echo "🧪 Running NebulaBox tests..."
	./nebulabox-test run

test-list: build-test
	@echo "📋 Listing NebulaBox tests..."
	./nebulabox-test list

test-category: build-test
	@echo "📂 Running tests by category..."
	./nebulabox-test category $(CATEGORY)

test-unit:
	@echo "🧪 Running Go unit tests..."
	go test -v ./internal/api/... ./internal/registry/... ./internal/containerd/...

test-cli:
	@echo "🧪 Running CLI tests..."
	@$(MAKE) build-cli-test
	go test -v ./internal/cli/tests/...

test-registry:
	@echo "🧪 Running registry-specific tests..."
	go test -v -run TestRegistry ./internal/api/... ./internal/registry/...

test-unit-coverage:
	@echo "📊 Running unit tests with coverage..."
	go test -v -coverprofile=coverage.out ./internal/api/... ./internal/registry/... ./internal/containerd/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-integration:
	@echo "🔗 Running integration tests..."
	go test -v -tags=integration ./internal/api/... -run TestIntegration

test-e2e:
	@echo "🌐 Running end-to-end tests..."
	@bash scripts/run-e2e-tests.sh

test-e2e-with-server:
	@echo "🌐 Running end-to-end tests (with server startup)..."
	@START_SERVER=true bash scripts/run-e2e-tests.sh

benchmark:
	@echo "📊 Running benchmarks..."
	go test -bench=. -benchmem -benchtime=3s ./internal/api/... ./internal/containerd/...

benchmark-api:
	@echo "📊 Running API benchmarks..."
	go test -bench=BenchmarkAPI -benchmem -benchtime=3s ./internal/api/...

benchmark-containerd:
	@echo "📊 Running containerd benchmarks..."
	go test -bench=. -benchmem -benchtime=3s ./internal/containerd/...

load-test:
	@echo "⚡ Running load tests..."
	go test -v -run TestLoad ./internal/api/...

benchmark-all:
	@echo "📊 Running all benchmarks..."
	go test -bench=. -benchmem -benchtime=3s ./internal/...

ci-test:
	@echo "🚀 Running CI test suite..."
	@bash scripts/ci-test.sh

ci-full:
	@echo "🚀 Running full CI pipeline..."
	@bash scripts/ci-test.sh
	@echo "🔨 Building all components..."
	@$(MAKE) build
	@echo "✅ CI pipeline completed!"

# Development targets
dev-full: build
	@echo "🚀 Starting full development environment..."
	@echo "API Server: http://localhost:8081"
	@echo "Dashboard: http://localhost:3000"
	@echo "Press Ctrl+C to stop all services"

# Help target
help:
	@echo "NebulaBox Makefile"
	@echo "=================="
	@echo ""
	@echo "Build targets:"
	@echo "  build              Build all components (API, test, registry)"
	@echo "  build-api          Build API server only"
	@echo "  build-test         Build test runner only"
	@echo "  build-registry     Build registry server only"
	@echo ""
	@echo "Run targets:"
	@echo "  run                Build and run API server"
	@echo "  api                Build and run API server"
	@echo "  registry           Build and run registry server"
	@echo "  test               Build and run all tests"
	@echo "  test-list          List all available tests"
	@echo "  test-category      Run tests by category"
	@echo ""
	@echo "Test targets:"
	@echo "  test-unit          Run Go unit tests"
	@echo "  test-integration   Run integration tests"
	@echo "  test-e2e           Run end-to-end tests"
	@echo "  test-registry      Run registry-specific tests"
	@echo "  registry-test      Run registry functionality tests"
	@echo "  registry-test-unit Run registry unit tests"
	@echo "  registry-test-integration  Run registry integration tests"
	@echo "  registry-health    Check registry health"
	@echo ""
	@echo "Benchmark targets:"
	@echo "  benchmark          Run all benchmarks"
	@echo "  benchmark-api      Run API benchmarks"
	@echo "  benchmark-containerd  Run containerd benchmarks"
	@echo "  load-test          Run load tests"
	@echo ""
	@echo "CI/CD targets:"
	@echo "  ci-test            Run CI test suite"
	@echo "  ci-full            Run full CI pipeline"
	@echo ""
	@echo "Clean targets:"
	@echo "  clean              Remove all build artifacts"
	@echo ""
	@echo "Development targets:"
	@echo "  dev-full           Start full development environment"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run"
	@echo "  make registry"
	@echo "  make registry-test"
	@echo "  make test-registry"
	@echo "  make clean"