.PHONY: build clean test lint install dev release uninstall uninstall-all

# Variables
BINARY_NAME=sshm
BUILD_DIR=dist
MAIN_PATH=./cmd/sshm
VERSION?=v1.0.0

# Default target
all: clean build

# Build for development
build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "-X github.com/levanduy/ssh_management/internal/cli.version=$(VERSION)" -o $(BINARY_NAME) $(MAIN_PATH)

# Build and auto-install globally
build-install: build
	@echo "Installing $(BINARY_NAME) globally..."
	@if [ -w /usr/local/bin ]; then \
		cp $(BINARY_NAME) /usr/local/bin/; \
		chmod +x /usr/local/bin/$(BINARY_NAME); \
		echo "✅ $(BINARY_NAME) installed to /usr/local/bin"; \
	else \
		echo "🔑 Installing requires sudo access..."; \
		sudo cp $(BINARY_NAME) /usr/local/bin/; \
		sudo chmod +x /usr/local/bin/$(BINARY_NAME); \
		echo "✅ $(BINARY_NAME) installed to /usr/local/bin"; \
	fi
	@echo "🚀 You can now run: $(BINARY_NAME)"

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	mkdir -p $(BUILD_DIR)
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X github.com/levanduy/ssh_management/internal/cli.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X github.com/levanduy/ssh_management/internal/cli.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X github.com/levanduy/ssh_management/internal/cli.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X github.com/levanduy/ssh_management/internal/cli.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

# Install to local bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Development mode - build and run
dev: build
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go mod download
	go mod tidy

# Create release archives
release: build-all
	@echo "Creating release archives..."
	cd $(BUILD_DIR) && \
	tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
	tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	tar -czf $(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64

# Uninstall from system
uninstall:
	@echo "Uninstalling SSH Manager..."
	@if [ -f /usr/local/bin/$(BINARY_NAME) ]; then \
		rm -f /usr/local/bin/$(BINARY_NAME) || sudo rm -f /usr/local/bin/$(BINARY_NAME); \
		echo "✅ SSH Manager binary removed."; \
	else \
		echo "⚠️  SSH Manager not found in /usr/local/bin/"; \
	fi
	@echo "💡 Data preserved at ~/.sshm/ - use 'make uninstall-all' to remove all data"

# Complete uninstall (binary + data)
uninstall-all:
	@echo "Removing SSH Manager completely..."
	@if [ -f /usr/local/bin/$(BINARY_NAME) ]; then \
		rm -f /usr/local/bin/$(BINARY_NAME) || sudo rm -f /usr/local/bin/$(BINARY_NAME); \
		echo "✅ SSH Manager binary removed."; \
	else \
		echo "⚠️  SSH Manager not found in /usr/local/bin/"; \
	fi
	@if [ -d ~/.sshm ]; then \
		rm -rf ~/.sshm; \
		echo "✅ SSH Manager data removed."; \
	else \
		echo "⚠️  No data found at ~/.sshm/"; \
	fi

# Show help
help:
	@echo "Available targets:"
	@echo "  build-install  - 🚀 Build and install globally (RECOMMENDED)"
	@echo "  build          - Build binary for current platform"
	@echo "  build-all      - Build for all supported platforms"
	@echo "  install        - Install existing binary to /usr/local/bin"
	@echo "  uninstall      - Remove binary only (preserve data)"
	@echo "  uninstall-all  - Remove binary and all data"
	@echo "  dev            - Build and run for development"
	@echo "  test           - Run tests"
	@echo "  lint           - Run linter"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Update dependencies"
	@echo "  release        - Create release archives"
	@echo "  help           - Show this help"
	@echo ""
	@echo "Quick setup: make build-install" 