# 🚀 SSH Manager v1.0.0

Terminal UI for SSH Host Management with auto-discovery and intelligent username detection.

## ✨ Features

- 🔍 **Smart Auto-Discovery**: Automatically discovers SSH hosts from `~/.ssh/known_hosts`
- 🧠 **Username Detection**: Intelligently detects usernames from shell history
- 🌐 **IP Resolution**: Resolves and displays IP addresses for all hosts
- 🖥️ **Beautiful TUI**: Clean terminal interface with intuitive navigation
- ⚡ **Instant Connect**: Connect to any host with just Enter
- 🗑️ **Safe Deletion**: Remove hosts from both database and known_hosts
- 📊 **Usage Stats**: Track connection frequency and usage patterns

---

## 🍺 Installation via Homebrew (Recommended)

### ⚡ Quick Install (Copy and paste these commands)

```bash
# Step 1: Add the SSH Manager tap
brew tap levanduy093-work/sshm
```

```bash
# Step 2: Install SSH Manager
brew install sshm
```

```bash
# Step 3: Launch SSH Manager
sshm
```

### 🔄 Alternative: One-line install

```bash
brew install levanduy093-work/sshm/sshm && sshm
```

### 📋 What happens during installation:

1. **Downloads source code** from this GitHub repository
2. **Compiles Go binary** automatically using your system's Go
3. **Installs to** `/usr/local/bin/sshm` (or `/opt/homebrew/bin/sshm` on Apple Silicon)
4. **Creates data directory** at `~/.sshm/`
5. **Backs up** your `~/.ssh/known_hosts` to `~/.ssh/known_hosts.backup.homebrew`

### ✅ Verify installation:

```bash
# Check if sshm is installed
which sshm
```

```bash
# Check version
sshm --version 2>/dev/null || echo "Version info available in TUI"
```

```bash
# Launch SSH Manager
sshm
```

---

## 📦 Alternative Installation Methods

### 🔧 Manual Installation (macOS/Linux)

```bash
# Download and extract
curl -sSL https://github.com/levanduy093-work/ssh_management/archive/refs/tags/v1.0.0.tar.gz | tar -xz
```

```bash
# Enter directory
cd ssh_management-1.0.0
```

```bash
# Build and install
make install
```

### 🏗️ Build from Source

```bash
# Clone repository
git clone https://github.com/levanduy093-work/ssh_management.git
```

```bash
# Enter directory
cd ssh_management
```

```bash
# Build and install
go build -o bin/sshm cmd/sshm/main.go
sudo install -m 755 bin/sshm /usr/local/bin/sshm
```

---

## 🎮 Usage Guide

### 🚀 Launch SSH Manager

```bash
sshm
```

**What happens automatically:**
1. 🔍 Scans `~/.ssh/known_hosts` for SSH hosts
2. 🧠 Detects usernames from shell history (`~/.zsh_history`, `~/.bash_history`)
3. 🌐 Resolves IP addresses for all discovered hosts
4. 🖥️ Displays everything in a beautiful terminal interface

### ⌨️ TUI Controls

```
Navigation:
  ↑/k        Move up
  ↓/j        Move down
  Enter      Connect to selected host
  
Search & Filter:
  /          Search hosts
  Esc        Cancel search
  
Host Management:
  x          Delete host (removes from database + known_hosts)
  r          Refresh (re-scan for new hosts)
  
General:
  q          Quit SSH Manager
  ?          Show help
```

### 📊 Example Interface

```
┌─ SSH Manager ─────────────────────────────────────────────┐
│ Total hosts: 5 • Selected: 1                             │
│                                                           │
│ webserver (admin@webserver.example.com:22) [203.0.113.10]│
│ Auto-detected from known_hosts (ssh-rsa) • Used 3 times  │
│                                                           │
│ database (dbuser@database.prod.com:22) [198.51.100.25]   │
│ Auto-detected from known_hosts (ssh-rsa) • ssh-detected  │
│                                                           │
│ api-server (deploy@api.staging.net:22) [192.0.2.50]      │
│ Auto-detected from known_hosts (ssh-ed25519) • Used 1 time│
└───────────────────────────────────────────────────────────┘

⌨️  Controls: ↑/k up • ↓/j down • / search • enter connect • x delete • r refresh • q quit
```

---

## 🔧 How Auto-Discovery Works

### 📂 Data Sources

1. **`~/.ssh/known_hosts`** - Extracts hostnames and IP addresses
2. **Shell History** - Finds SSH commands to detect usernames:
   - `~/.zsh_history`
   - `~/.bash_history`
   - `~/.history`
3. **SSH Config** - Reads `~/.ssh/config` for user settings
4. **DNS Resolution** - Resolves IP addresses for hostnames

### 🧠 Username Detection Examples

```bash
# Your shell history contains:
ssh admin@webserver.example.com
ssh dbuser@database.prod.com
ssh deploy@api.staging.net

# SSH Manager automatically detects:
# webserver.example.com → username: admin
# database.prod.com → username: dbuser  
# api.staging.net → username: deploy
```

---

## 📁 Data Storage

### 🗂️ File Locations

```bash
# Main database
~/.sshm/hosts.db

# Backup files (created automatically)
~/.ssh/known_hosts.backup
~/.ssh/known_hosts.backup.homebrew
```

### 🔍 What's stored

- Host information (name, hostname, port, username)
- IP addresses (resolved automatically)
- Usage statistics (connection count, last used)
- Discovery metadata (source, key type)

---

## 🗑️ Uninstallation

### 🍺 Homebrew Uninstall

```bash
# Remove SSH Manager
brew uninstall sshm
```

```bash
# Remove the tap (optional)
brew untap levanduy093-work/sshm
```

```bash
# Remove data directory (optional)
rm -rf ~/.sshm/
```

```bash
# Restore original known_hosts (if needed)
cp ~/.ssh/known_hosts.backup.homebrew ~/.ssh/known_hosts
```

### 🔧 Manual Uninstall

```bash
# Remove binary
sudo rm /usr/local/bin/sshm
```

```bash
# Remove data (optional)
rm -rf ~/.sshm/
```

---

## 🐛 Troubleshooting

### ❌ Common Issues

**1. `brew tap` fails:**
```bash
# Make sure you have Homebrew installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

**2. `brew install sshm` fails:**
```bash
# Update Homebrew
brew update

# Try again
brew install sshm
```

**3. `sshm` command not found:**
```bash
# Check installation path
brew --prefix

# Add to PATH (add to ~/.zshrc or ~/.bashrc)
export PATH="/opt/homebrew/bin:$PATH"  # Apple Silicon
export PATH="/usr/local/bin:$PATH"     # Intel Mac
```

**4. No hosts detected:**
```bash
# Check if known_hosts exists
ls -la ~/.ssh/known_hosts

# Check shell history
ls -la ~/.zsh_history ~/.bash_history

# Manual refresh in SSH Manager
# Press 'r' key in the TUI
```

### 🔧 System Requirements

- **OS**: macOS 10.15+ or Linux
- **Dependencies**: Go 1.24+ (installed automatically by Homebrew)
- **Files**: `~/.ssh/known_hosts` (for auto-discovery)

---

## 🎯 Quick Start Examples

### 🚀 First Time User

```bash
# Install
brew tap levanduy093-work/sshm && brew install sshm

# Launch (auto-discovers your SSH hosts)
sshm

# Use arrow keys to navigate, Enter to connect
```

### 🔄 Daily Usage

```bash
# Quick launch
sshm

# Search for specific host
# Press '/' then type hostname

# Connect to host
# Navigate to host, press Enter
```

### 🧹 Clean Uninstall

```bash
# Remove everything
brew uninstall sshm
brew untap levanduy093-work/sshm
rm -rf ~/.sshm/
```

---

## 📞 Support & Links

- 📖 **Documentation**: [GitHub Repository](https://github.com/levanduy093-work/ssh_management)
- 🐛 **Report Issues**: [GitHub Issues](https://github.com/levanduy093-work/ssh_management/issues)
- 💡 **Feature Requests**: [GitHub Issues](https://github.com/levanduy093-work/ssh_management/issues)
- 🍺 **Homebrew Tap**: [levanduy093-work/homebrew-sshm](https://github.com/levanduy093-work/homebrew-sshm)

---

**🎉 Enjoy using SSH Manager! 🎉**

*If this tool helps you, please give it a ⭐ on GitHub!* 