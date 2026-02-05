.PHONY: build install clean test homebrew-setup help

# Variables
BINARY_NAME=sshm
BUILD_DIR=bin
INSTALL_PATH=/usr/local/bin
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building SSH Manager..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sshm
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install the application
install: build
	@echo "📦 Installing SSH Manager..."
	sudo install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "🚀 Run 'sshm' to start SSH Manager"

# Uninstall the application
uninstall:
	@echo "🗑️  Uninstalling SSH Manager..."
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Uninstalled from $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "💡 Data directory ~/.sshm/ preserved. Remove manually if needed."

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...
	@echo "✅ Tests complete"

# Build for multiple platforms
build-all:
	@echo "🌍 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/sshm
	
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/sshm
	
	# Linux (Intel)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/sshm
	
	# Linux (ARM)
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/sshm

	# Windows (Intel/AMD64)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/sshm

	# Windows (ARM64)
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd/sshm
	
	@echo "✅ Multi-platform build complete"
	@ls -la $(BUILD_DIR)/

# Setup Homebrew formula
homebrew-setup:
	@echo "🍺 Setting up Homebrew formula..."
	./scripts/setup-homebrew.sh

# Create a release
release: clean test build-all
	@echo "🚀 Creating release..."
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION not set. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	
	# Create git tag
	git tag $(VERSION)
	git push origin $(VERSION)
	
	# Setup Homebrew
	$(MAKE) homebrew-setup
	
	@echo "✅ Release $(VERSION) created!"
	@echo "📋 Next steps:"
	@echo "   1. Create GitHub release with binaries from $(BUILD_DIR)/"
	@echo "   2. Setup Homebrew tap repository"
	@echo "   3. Update package managers"

# Development setup
dev-setup:
	@echo "🛠️  Setting up development environment..."
	go mod tidy
	go mod download
	@echo "✅ Development setup complete"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Show help
help:
	@echo "SSH Manager (sshm) - Available commands:"
	@echo ""
	@echo "🔨 Build commands:"
	@echo "  build         Build the application"
	@echo "  build-all     Build for multiple platforms"
	@echo "  install       Build and install to system"
	@echo "  uninstall     Remove from system"
	@echo ""
	@echo "🧪 Development commands:"
	@echo "  dev-setup     Setup development environment"
	@echo "  test          Run tests"
	@echo "  fmt           Format code"
	@echo "  lint          Lint code"
	@echo "  clean         Clean build artifacts"
	@echo ""
	@echo "🚀 Release commands:"
	@echo "  homebrew-setup Setup Homebrew formula"
	@echo "  release       Create a release (requires VERSION=vX.X.X)"
	@echo ""
	@echo "💡 Examples:"
	@echo "  make build"
	@echo "  make install"
	@echo "  make release VERSION=v1.0.0"
	@echo "  make homebrew-setup" 
