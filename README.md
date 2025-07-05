# SSH Manager (sshm)

🚀 **Quản lý SSH dễ dàng và trực quan trên Terminal**

SSH Manager là một công cụ terminal giúp bạn quản lý các kết nối SSH một cách dễ dàng, nhanh chóng và trực quan. Được viết bằng Go với giao diện TUI (Terminal User Interface) đẹp mắt.

![Terminal Interface](https://img.shields.io/badge/Interface-Terminal%20TUI-blue)
![Go Version](https://img.shields.io/badge/Go-1.24+-green)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)

## ✨ Tính năng chính

### 🔗 Quản lý Host SSH
- ➕ **Thêm, sửa, xóa** SSH hosts dễ dàng
- 📊 **Theo dõi thống kê** sử dụng (số lần kết nối, lần cuối)
- 🏷️ **Gắn tags** để phân loại hosts
- 📝 **Mô tả chi tiết** cho từng host

### 🎯 Giao diện thân thiện
- 🖥️ **TUI tương tác** với Bubble Tea - đẹp mắt và mượt mà
- 🔍 **Tìm kiếm nhanh** theo tên, hostname, mô tả, tags
- ⌨️ **Phím tắt trực quan** - không cần nhớ nhiều lệnh
- 📱 **Responsive** - tự động điều chỉnh theo kích thước terminal

### ⚡ Hiệu suất cao
- 🗄️ **SQLite local** - nhanh và nhẹ
- 🔒 **An toàn** - không lưu trữ mật khẩu, chỉ sử dụng SSH keys
- 📦 **Single binary** - không phụ thuộc external dependencies

## 📦 Cài đặt

### Via Homebrew (Khuyến nghị)
```bash
# Sẽ có sẵn sau khi phát hành
brew tap levanduy/ssh_management
brew install sshm
```

### Quick Setup (Recommended)
```bash
# Clone repository
git clone https://github.com/levanduy/ssh_management.git
cd ssh_management

# One-command setup (build + install globally)
./setup.sh
```

### Manual Installation
```bash
# Clone repository
git clone https://github.com/levanduy/ssh_management.git
cd ssh_management

# Option 1: Using Makefile
make build-install

# Option 2: Manual steps
go build -o sshm ./cmd/sshm
sudo mv sshm /usr/local/bin/
```

## 🚀 Sử dụng

### 1. Interactive TUI Mode (Khuyến nghị)
```bash
# Khởi chạy giao diện tương tác
sshm
# hoặc
sshm tui
```

**Phím tắt trong TUI:**
- `↑/↓` - Di chuyển trong danh sách
- `Enter` - Kết nối SSH
- `/` - Tìm kiếm
- `d` - Xóa host
- `r` - Refresh danh sách
- `q` - Thoát

### 2. CLI Commands

#### Thêm host mới
```bash
sshm add myserver
# Hệ thống sẽ hỏi thông tin: hostname, username, port, SSH key...

# Hoặc thêm trực tiếp
sshm add prod-server --hostname 192.168.1.100 --username ubuntu --port 22
```

#### Liệt kê tất cả hosts
```bash
sshm list
# hoặc
sshm ls
```

#### Kết nối SSH
```bash
sshm connect myserver
# hoặc theo ID
sshm connect 1
```

#### Tìm kiếm hosts
```bash
sshm search production
sshm search nginx
```

#### Xem thông tin chi tiết
```bash
sshm show myserver
```

#### Xóa host
```bash
sshm remove myserver
# Hệ thống sẽ xác nhận trước khi xóa
```

## 📋 Ví dụ sử dụng

### Thêm một server production
```bash
sshm add prod-web
# Host name: prod-web
# Hostname/IP: 192.168.1.100
# Username: ubuntu
# Port [22]: 22
# SSH key path: ~/.ssh/prod_key
# Description: Production web server
# Tags: production, web, nginx
```

### Kết nối nhanh
```bash
# Mở TUI và chọn server
sshm

# Hoặc kết nối trực tiếp
sshm connect prod-web
```

## ⚙️ Cấu hình

### Database Location
Mặc định sshm lưu dữ liệu tại: `~/.sshm/hosts.db`

Có thể thay đổi bằng flag `--db`:
```bash
sshm --db /custom/path/hosts.db list
```

### SSH Key Management
sshm không lưu trữ mật khẩu, chỉ sử dụng:
- SSH keys (khuyến nghị)
- SSH agent
- System SSH client

## 🏗️ Kiến trúc

```
ssh_management/
├── cmd/sshm/              # Entry point
├── internal/
│   ├── domain/            # Business models
│   ├── service/           # Business logic  
│   ├── repo/              # Data persistence (SQLite)
│   ├── cli/               # CLI commands (Cobra)
│   └── ui/                # TUI interface (Bubble Tea)
├── pkg/ssh/               # SSH utilities
└── README.md
```

**Tech Stack:**
- **Language**: Go 1.24+
- **CLI Framework**: Cobra
- **TUI Framework**: Bubble Tea + Bubbles + Lipgloss
- **Database**: SQLite (modernc.org/sqlite)
- **Architecture**: Clean Architecture

## 🤝 Đóng góp

1. Fork repository
2. Tạo feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Tạo Pull Request

## 📄 License

Dự án được phát hành dưới MIT License. Xem [LICENSE](LICENSE) để biết thêm chi tiết.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Framework TUI tuyệt vời
- [Cobra](https://github.com/spf13/cobra) - CLI framework mạnh mẽ
- [SQLite](https://sqlite.org/) - Database nhẹ và tin cậy

## 📞 Liên hệ

- GitHub: [@levanduy](https://github.com/levanduy)
- Email: your.email@example.com

---

⭐ **Nếu project hữu ích, hãy cho một star để ủng hộ!** ⭐ 