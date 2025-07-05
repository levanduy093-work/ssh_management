#!/bin/bash

# Setup Homebrew tap for SSH Manager
# This script helps setup the Homebrew formula

set -e

echo "🍺 Setting up Homebrew tap for SSH Manager..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [[ ! -f "cmd/sshm/main.go" ]]; then
    print_error "This script must be run from the project root directory"
    exit 1
fi

# Check if git tag exists
if ! git describe --tags --exact-match HEAD 2>/dev/null; then
    print_warning "No git tag found for current commit"
    print_warning "You may want to create a tag first: git tag v1.0.0"
fi

# Get current version/tag
VERSION=$(git describe --tags --exact-match HEAD 2>/dev/null || echo "v1.0.0")
print_status "Using version: $VERSION"

# Create release archive
print_status "Creating release archive..."
git archive --format=tar.gz --prefix=ssh_management-${VERSION#v}/ HEAD > /tmp/ssh_management-${VERSION}.tar.gz

# Calculate SHA256
print_status "Calculating SHA256..."
if command -v shasum >/dev/null 2>&1; then
    SHA256=$(shasum -a 256 /tmp/ssh_management-${VERSION}.tar.gz | cut -d' ' -f1)
elif command -v sha256sum >/dev/null 2>&1; then
    SHA256=$(sha256sum /tmp/ssh_management-${VERSION}.tar.gz | cut -d' ' -f1)
else
    print_error "Neither shasum nor sha256sum found"
    exit 1
fi

print_status "SHA256: $SHA256"

# Update formula with correct SHA256
print_status "Updating Homebrew formula..."
sed -i.bak "s/PLACEHOLDER_SHA256/$SHA256/" homebrew-formula/sshm.rb
sed -i.bak "s|refs/tags/v1.0.0|refs/tags/$VERSION|" homebrew-formula/sshm.rb

print_status "Formula updated successfully!"

# Instructions for publishing
cat << EOF

🍺 Homebrew Setup Complete!

📋 Next Steps:

1. Create a GitHub release:
   git tag $VERSION
   git push origin $VERSION
   
2. Create a Homebrew tap repository:
   - Create a new GitHub repository named 'homebrew-sshm'
   - Copy the formula: cp homebrew-formula/sshm.rb /path/to/homebrew-sshm/Formula/sshm.rb
   
3. Users can then install with:
   brew tap levanduy/sshm
   brew install sshm

4. Or install directly:
   brew install levanduy/sshm/sshm

📁 Files created/updated:
   - homebrew-formula/sshm.rb (SHA256: $SHA256)
   - /tmp/ssh_management-${VERSION}.tar.gz (test archive)

🔗 Useful links:
   - Homebrew Formula Cookbook: https://docs.brew.sh/Formula-Cookbook
   - How to Create a Tap: https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap

EOF

# Clean up
rm -f /tmp/ssh_management-${VERSION}.tar.gz
rm -f homebrew-formula/sshm.rb.bak

print_status "Setup complete! 🎉" 