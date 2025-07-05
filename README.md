# SSH Manager (sshm)

<div align="center">

![SSH Manager](https://img.shields.io/badge/SSH-Manager-blue?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**🚀 Terminal UI for SSH Host Management**

*Zero-config SSH host manager with auto-discovery and intelligent username detection*

</div>

## ✨ Features

- 🔍 **Smart Auto-Discovery**: Automatically discovers SSH hosts from `~/.ssh/known_hosts`
- 🧠 **Username Detection**: Intelligently detects usernames from shell history
- 🌐 **IP Resolution**: Resolves and displays IP addresses for all hosts
- 🖥️ **Beautiful TUI**: Clean terminal interface with intuitive navigation
- ⚡ **Instant Connect**: Connect to any host with just Enter
- 🗑️ **Safe Deletion**: Remove hosts from both database and known_hosts
- 📊 **Usage Stats**: Track connection frequency and usage patterns
- 💾 **Lightweight**: Single binary, no complex configuration needed

## 🎯 Philosophy

SSH Manager follows the **"Just Works"** principle:

- ✅ **Zero Configuration**: No setup required - just run `sshm`
- ✅ **Auto-Discovery**: Automatically finds all your SSH hosts
- ✅ **Smart Detection**: Detects usernames from your actual SSH usage
- ✅ **TUI-Only**: Clean interface focused on the task at hand
- ✅ **Safe Operations**: Backup and restore capabilities

## 🚀 Installation

### Homebrew (Recommended)
```bash
# Add tap
brew tap levanduy/sshm

# Install SSH Manager
brew install sshm

# Or install directly
brew install levanduy/sshm/sshm
```

### Quick Install from Source
```bash
# Clone and install
git clone https://github.com/levanduy/ssh_management.git
cd ssh_management
make install
```

### Manual Build
```bash
# Clone repository
git clone https://github.com/levanduy/ssh_management.git
cd ssh_management

# Build and install
go build -o bin/sshm cmd/sshm/main.go
sudo install -m 755 bin/sshm /usr/local/bin/sshm
```

### Build Requirements
- Go 1.24+
- Linux/macOS (Unix-like systems)

## 🎮 Usage

### Launch SSH Manager
```bash
sshm
```

That's it! SSH Manager will automatically:
1. 🔍 Scan `~/.ssh/known_hosts` for hosts
2. 🧠 Detect usernames from shell history (`~/.zsh_history`, `~/.bash_history`)
3. 🌐 Resolve IP addresses for all hosts
4. 🖥️ Display everything in a beautiful TUI

### TUI Controls

```
┌─ SSH Manager ─────────────────────────────────────────────┐
│ Total hosts: 7 • Selected: 1                             │
│                                                           │
│ levanduy (duy@levanduy.local:22) [192.168.1.100]         │
│ Auto-detected from known_hosts (ssh-ed25519) • Used 5 times │
│                                                           │
│ webserver (admin@webserver.example.com:22) [203.0.113.10] │
│ Auto-detected from known_hosts (ssh-rsa) • ssh-detected   │
│                                                           │
│ database (dbuser@database.prod.com:22) [198.51.100.25]   │
│ Auto-detected from known_hosts (ssh-rsa) • ssh-detected   │
└───────────────────────────────────────────────────────────┘

⌨️  Controls:
  ↑/k up • ↓/j down • / search • enter connect • x delete • r refresh • q quit
```

### Key Features in Action

#### Smart Username Detection
```bash
# Your SSH history:
ssh admin@webserver.example.com
ssh dbuser@database.prod.com
ssh deploy@api.staging.net

# SSH Manager automatically detects:
# webserver.example.com → username: admin
# database.prod.com → username: dbuser
# api.staging.net → username: deploy
```

#### IP Address Resolution
```bash
# Automatically resolves and displays:
# webserver.example.com → [203.0.113.10]
# database.prod.com → [198.51.100.25]
# api.staging.net → [192.0.2.50]
```

#### Safe Host Deletion
```bash
# Press 'x' on any host to see:
┌─ Delete Host Confirmation ─────────────────────┐
│                                                │
│ Host: webserver                                │
│ Connection: admin@webserver.example.com:22     │
│ IP: 203.0.113.10                              │
│                                                │
│ This will remove the host from:                │
│ • SSH Manager database                         │
│ • ~/.ssh/known_hosts file                     │
│                                                │
│ Continue? (y/N)                                │
│                                                │
│ Press 'y' to confirm • 'n' or 'Esc' to cancel │
└────────────────────────────────────────────────┘
```

## 🏗️ Architecture

```
sshm (single binary)
├── TUI Interface (Bubble Tea + Lipgloss)
├── Auto-Discovery Engine
│   ├── known_hosts parser
│   ├── Shell history analyzer
│   └── IP address resolver
├── SQLite Database (~/.sshm/)
└── SSH Integration (system ssh)
```

**Tech Stack:**
- **Language**: Go 1.24+
- **TUI**: Bubble Tea + Bubbles + Lipgloss
- **Database**: SQLite (modernc.org/sqlite)
- **SSH**: System SSH client integration

## 🔧 How It Works

### Auto-Discovery Process
1. **Parse known_hosts**: Extracts hostnames and IP addresses
2. **Analyze shell history**: Finds SSH commands with usernames
3. **Match and merge**: Combines information from both sources
4. **Resolve IPs**: Uses system DNS to resolve missing IP addresses
5. **Update database**: Stores or updates host information

### Username Detection Sources
1. **SSH Config** (`~/.ssh/config`) - Highest priority
2. **Shell History** (`~/.zsh_history`, `~/.bash_history`) - Smart parsing
3. **System Username** - Fallback option

### Data Storage
```bash
# Database location
~/.sshm/hosts.db

# Backup your known_hosts (automatically created)
~/.ssh/known_hosts.backup
```

## 🎯 Use Cases

### Daily Developer Workflow
```bash
# Throughout the day, you SSH to various servers:
ssh admin@prod-web-01
ssh deploy@staging-api
ssh monitor@monitoring.cloud.io

# Later, just run SSH Manager:
sshm
# → All servers appear with correct usernames
# → IP addresses resolved and displayed
# → Connect to any server with just Enter
```

### Server Administration
```bash
# Manage multiple environments:
ssh root@prod-db-01
ssh admin@staging-web-02
ssh backup@backup.internal

# SSH Manager organizes everything:
# → Tracks usage frequency
# → Shows connection details
# → Enables quick switching between servers
```

## 🛠️ Configuration

### Environment Variables
```bash
# Custom database path
export SSHM_DB_PATH="/custom/path/hosts.db"

# Disable auto-discovery on startup
export SSHM_AUTO_DISCOVERY=false
```

### SSH Integration
SSH Manager leverages your existing SSH setup:
- ✅ SSH keys (`~/.ssh/`)
- ✅ SSH agent
- ✅ SSH config (`~/.ssh/config`)
- ✅ Known hosts (`~/.ssh/known_hosts`)
- ❌ No password storage

## 🗑️ Uninstall

```bash
# Remove binary
sudo rm /usr/local/bin/sshm

# Remove all data (optional)
rm -rf ~/.sshm/

# Restore original known_hosts if needed
cp ~/.ssh/known_hosts.backup ~/.ssh/known_hosts
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Create Pull Request

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Excellent TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Beautiful styling
- [SQLite](https://sqlite.org/) - Reliable embedded database
- Go community - Amazing ecosystem

## 📞 Contact

- GitHub: [@levanduy](https://github.com/levanduy)
- Project: [ssh_management](https://github.com/levanduy/ssh_management)

---

⭐ **If this project helps you, please give it a star!** ⭐ 