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

## 🍺 Installation

### Homebrew (Recommended)

```bash
# Add tap and install
brew tap levanduy093-work/sshm
brew install sshm
```

### One-line install
```bash
brew install levanduy093-work/sshm/sshm
```

### Manual Installation

```bash
# Download and extract
curl -sSL https://github.com/levanduy093-work/ssh_management/archive/refs/tags/v1.0.0.tar.gz | tar -xz
cd ssh_management-1.0.0

# Build and install
go build -o sshm cmd/sshm/main.go
sudo install -m 755 sshm /usr/local/bin/sshm
```

---

## 🎮 Usage

```bash
# Launch SSH Manager
sshm
```

**TUI Controls:**
- `↑/↓` or `j/k` - Navigate
- `Enter` - Connect to host
- `/` - Search hosts
- `x` - Delete host
- `r` - Refresh/discover
- `q` - Quit

---

## 🗑️ Uninstallation

### Homebrew
```bash
# Remove application
brew uninstall sshm
brew untap levanduy093-work/sshm

# Remove data (optional)
rm -rf ~/.sshm/
```

### Manual
```bash
# Remove binary
sudo rm /usr/local/bin/sshm

# Remove data (optional)
rm -rf ~/.sshm/
```

---

## 📄 License

This project is licensed under the **MIT License** - free for personal and commercial use.

## 📞 Contact & Support

- **Email**: levanduy.work@gmail.com
- **Repository**: [GitHub](https://github.com/levanduy093-work/ssh_management)
- **Issues**: [Report Issues](https://github.com/levanduy093-work/ssh_management/issues)

---

**🎉 Enjoy using SSH Manager!** 

*This software is completely free to use. If this tool helps you, please give it a ⭐ on GitHub!* 