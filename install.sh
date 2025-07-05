#!/bin/bash

# SSH Manager Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/levanduy/ssh_management/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="sshm"
INSTALL_DIR="/usr/local/bin"
REPO_URL="https://github.com/levanduy/ssh_management"
API_URL="https://api.github.com/repos/levanduy/ssh_management/releases/latest"

# Functions
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

# Detect OS and Architecture
detect_platform() {
    local os
    local arch
    
    case "$(uname -s)" in
        Darwin*)
            os="darwin"
            ;;
        Linux*)
            os="linux"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            ;;
    esac
    
    case "$(uname -m)" in
        x86_64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release info
get_latest_release() {
    if command -v curl >/dev/null 2>&1; then
        curl -s "$API_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -q -O - "$API_URL"
    else
        print_error "Neither curl nor wget is available. Please install one of them."
    fi
}

# Download and install
install_sshm() {
    local platform
    local download_url
    local temp_dir
    local binary_path
    
    print_info "Detecting platform..."
    platform=$(detect_platform)
    print_info "Platform detected: $platform"
    
    print_info "Fetching latest release information..."
    local release_info
    release_info=$(get_latest_release)
    
    if [ -z "$release_info" ]; then
        print_error "Failed to fetch release information"
    fi
    
    # Extract download URL for the platform
    download_url=$(echo "$release_info" | grep "browser_download_url.*${BINARY_NAME}-${platform}.tar.gz" | cut -d '"' -f 4)
    
    if [ -z "$download_url" ]; then
        print_error "No release found for platform: $platform"
    fi
    
    print_info "Downloading SSH Manager..."
    temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    if command -v curl >/dev/null 2>&1; then
        curl -sSL "$download_url" -o "${BINARY_NAME}.tar.gz"
    else
        wget -q "$download_url" -O "${BINARY_NAME}.tar.gz"
    fi
    
    print_info "Extracting archive..."
    tar -xzf "${BINARY_NAME}.tar.gz"
    
    # Find the binary
    binary_path=$(find . -name "${BINARY_NAME}-${platform}" -type f)
    if [ -z "$binary_path" ]; then
        print_error "Binary not found in archive"
    fi
    
    print_info "Installing to $INSTALL_DIR..."
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        print_info "Installing requires sudo access..."
        sudo mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup
    cd - >/dev/null
    rm -rf "$temp_dir"
    
    print_success "SSH Manager installed successfully!"
    print_info "You can now run: $BINARY_NAME"
}

# Check if already installed
check_existing() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_version
        current_version=$($BINARY_NAME version 2>/dev/null | head -n1 | awk '{print $NF}' || echo "unknown")
        print_warning "SSH Manager is already installed (version: $current_version)"
        print_info "This will update to the latest version."
        read -p "Continue? [y/N]: " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled."
            exit 0
        fi
    fi
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$($BINARY_NAME version 2>/dev/null | head -n1 | awk '{print $NF}' || echo "unknown")
        print_success "Installation verified! Version: $version"
        print_info ""
        print_info "🚀 Quick start:"
        print_info "   $BINARY_NAME add myserver    # Add your first SSH host"
        print_info "   $BINARY_NAME                 # Launch interactive mode"
        print_info "   $BINARY_NAME --help          # Show all commands"
        print_info ""
        print_info "📖 Documentation: $REPO_URL"
    else
        print_error "Installation verification failed"
    fi
}

# Main function
main() {
    echo -e "${BLUE}"
    echo "╔═══════════════════════════════════════╗"
    echo "║        SSH Manager Installer         ║"
    echo "║    Easy SSH host management tool     ║"
    echo "╚═══════════════════════════════════════╝"
    echo -e "${NC}"
    
    check_existing
    install_sshm
    verify_installation
}

# Run main function
main "$@" 