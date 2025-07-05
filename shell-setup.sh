#!/bin/bash

# Shell Setup Script for SSH Manager
# Adds convenient aliases and PATH setup

SHELL_CONFIG=""

# Detect shell and config file
if [ -n "$ZSH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
else
    echo "❌ Unsupported shell. Please manually add to your shell config."
    exit 1
fi

echo "🔧 Setting up shell integration for SSH Manager..."

# Create backup
if [ -f "$SHELL_CONFIG" ]; then
    cp "$SHELL_CONFIG" "${SHELL_CONFIG}.backup.$(date +%Y%m%d_%H%M%S)"
    echo "📁 Backup created: ${SHELL_CONFIG}.backup.$(date +%Y%m%d_%H%M%S)"
fi

# Add SSH Manager section
cat >> "$SHELL_CONFIG" << 'EOF'

# SSH Manager aliases and functions
alias ssh-list="sshm list"
alias ssh-add="sshm add"
alias ssh-search="sshm search"
alias ssm="sshm"  # Short alias

# Quick connect function
ssh-connect() {
    if [ $# -eq 0 ]; then
        sshm
    else
        sshm connect "$1"
    fi
}
EOF

echo "✅ Shell configuration updated!"
echo "🔄 Please run: source $SHELL_CONFIG"
echo ""
echo "🎯 New aliases available:"
echo "   ssm                  # Short alias for sshm"
echo "   ssh-list            # List all hosts"
echo "   ssh-add             # Add new host"
echo "   ssh-search <query>  # Search hosts"
echo "   ssh-connect [host]  # Connect to host or open TUI" 