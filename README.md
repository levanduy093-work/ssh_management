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
brew tap levanduy093-work/sshm && brew install sshm
```

### Manual Build
```bash
git clone https://github.com/levanduy093-work/ssh_management.git
cd ssh_management && go build -o sshm cmd/sshm/main.go
sudo install -m 755 sshm /usr/local/bin/sshm
```

## 🎮 Usage

```bash
sshm
```

**TUI Controls:**
- `↑/↓` or `j/k` - Navigate hosts
- `Enter` - Connect to selected host
- `/` - Search hosts
- `x` - Delete host
- `r` - Refresh/discover
- `q` - Quit

## 🗑️ Uninstall

```bash
# Homebrew
brew uninstall sshm && brew untap levanduy093-work/sshm

# Manual
sudo rm /usr/local/bin/sshm

# Remove data (optional)
rm -rf ~/.sshm/
```

## 🔧 How It Works

SSH Manager automatically:
1. **Scans** `~/.ssh/known_hosts` for hosts
2. **Detects** usernames from shell history
3. **Resolves** IP addresses
4. **Organizes** everything in a clean TUI

**Data Storage:** `~/.sshm/hosts.db`

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Create Pull Request

## 📄 License

This project is licensed under the **MIT License** - **completely free for personal and commercial use**.

**What this means:**
- ✅ **Free to use** - No cost, no subscription, no hidden fees
- ✅ **Commercial use** - Use in your business or sell products that include this software
- ✅ **Modify** - Change the code to fit your needs
- ✅ **Distribute** - Share with others
- ✅ **Private use** - Use for personal projects

See [LICENSE](LICENSE) for full details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Excellent TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Beautiful styling
- [SQLite](https://sqlite.org/) - Reliable embedded database
- Go community - Amazing ecosystem

## 📞 Contact & Support

- **Email**: levanduy.work@gmail.com
- **GitHub**: [@levanduy093-work](https://github.com/levanduy093-work)
- **Project**: [ssh_management](https://github.com/levanduy093-work/ssh_management)
- **Issues**: [Report Issues](https://github.com/levanduy093-work/ssh_management/issues)

---

⭐ **This software is completely free to use. If this project helps you, please give it a star!** ⭐ 