#!/bin/bash

# SSH Manager Local Setup Script
# Usage: ./setup.sh

set -e

echo "🔨 Building SSH Manager..."
go build -o sshm ./cmd/sshm

echo "📦 Installing SSH Manager globally..."

# Check if /usr/local/bin is writable
if [ -w /usr/local/bin ]; then
    cp sshm /usr/local/bin/
    chmod +x /usr/local/bin/sshm
    echo "✅ SSH Manager installed to /usr/local/bin/sshm"
else
    echo "🔑 Installing requires sudo access..."
    sudo cp sshm /usr/local/bin/
    sudo chmod +x /usr/local/bin/sshm
    echo "✅ SSH Manager installed to /usr/local/bin/sshm"
fi

echo "🧹 Cleaning up local binary..."
rm -f sshm

echo ""
echo "🎉 Setup complete!"
echo "You can now run: sshm"
echo ""
echo "🚀 Quick start:"
echo "   sshm add myserver    # Add your first SSH host"
echo "   sshm                 # Launch interactive mode"
echo "   sshm --help          # Show all commands" 