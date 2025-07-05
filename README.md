# SSH Manager (sshm)

<div align="center">

![SSH Manager](https://img.shields.io/badge/SSH-Manager-blue?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**🚀 Giao diện TUI đơn giản để quản lý SSH hosts**

*Tự động phát hiện từ `~/.ssh/known_hosts` và cung cấp giao diện terminal tương tác*

</div>

## ✨ Tính năng

- 🔍 **Auto-discovery**: Tự động phát hiện SSH hosts từ `~/.ssh/known_hosts`
- 🖥️ **TUI Interface**: Giao diện terminal tương tác đẹp mắt
- ⚡ **Kết nối nhanh**: Chọn và kết nối SSH chỉ với vài phím
- 📊 **Thống kê sử dụng**: Theo dõi tần suất sử dụng các hosts
- 🏷️ **Thông tin chi tiết**: Hiển thị user, port, mô tả cho mỗi host
- 💾 **Nhẹ nhàng**: Single binary, không cần cấu hình phức tạp

## 🎯 Triết lý

SSH Manager được thiết kế theo nguyên tắc **đơn giản và hiệu quả**:

- ✅ **Zero config**: Không cần setup, chỉ cần chạy `sshm`
- ✅ **Auto-discovery**: Tự động phát hiện hosts đã kết nối
- ✅ **TUI-only**: Chỉ giao diện terminal, không có CLI commands rườm rà
- ✅ **Lightweight**: Tập trung vào task chính: browse và connect

## 🚀 Cài đặt

### Cài đặt nhanh
```bash
# Download và cài đặt từ GitHub Releases
curl -sSL https://github.com/levanduy/ssh_management/releases/latest/download/install.sh | bash
```

### Build từ source
```bash
git clone https://github.com/levanduy/ssh_management.git
cd ssh_management
make install
```

### Manual build
```bash
go build -o sshm ./cmd/sshm
sudo cp sshm /usr/local/bin/
```

## 🎮 Sử dụng

### Khởi động SSH Manager
```bash
sshm
```

Chỉ cần vậy thôi! SSH Manager sẽ:
1. 🔍 Tự động scan `~/.ssh/known_hosts` để tìm hosts
2. 🖥️ Hiển thị giao diện TUI với danh sách hosts
3. ⚡ Cho phép bạn browse và kết nối ngay lập tức

### Điều khiển TUI

```
┌─ SSH Manager ─────────────────────────┐
│ ID   NAME            HOST             │
│ ──────────────────────────────────── │
│ 1    server-prod     192.168.1.100    │
│ 2    github          github.com       │
│ 3    vps-dev         dev.example.com  │
└──────────────────────────────────────┘

⌨️  Phím tắt:
  ↑/↓    Điều hướng
  Enter  Kết nối SSH
  /      Tìm kiếm
  d      Xóa host
  r      Refresh (scan lại)
  q      Thoát
```

### Auto-Discovery

SSH Manager tự động phát hiện hosts mỗi khi chạy:

```bash
# Khi bạn đã SSH đến hosts mới:
ssh user@newserver.com

# Sau đó chạy sshm, nó sẽ tự động thêm:
sshm
# → "🔍 Auto-discovered 1 new SSH host(s)"
```

### Tắt auto-discovery (nếu cần)
```bash
sshm --auto-discovery=false
```

## 🏗️ Kiến trúc

```
sshm (single binary)
├── TUI Interface (Bubble Tea)
├── Auto-discovery (known_hosts)
├── SQLite Database (~/.sshm/)
└── SSH Integration (system ssh)
```

**Tech Stack:**
- **Language**: Go 1.24+
- **TUI**: Bubble Tea + Bubbles + Lipgloss  
- **Database**: SQLite (modernc.org/sqlite)
- **SSH**: System SSH client

## ⚙️ Cấu hình

### Database Location
```bash
# Mặc định
~/.sshm/hosts.db

# Custom database path
sshm --db /custom/path/hosts.db
```

### SSH Key Management
SSH Manager sử dụng SSH client của hệ thống:
- ✅ SSH keys (`~/.ssh/`)
- ✅ SSH agent
- ✅ SSH config (`~/.ssh/config`)
- ❌ Không lưu trữ passwords

## 🎯 Use Cases

### Developer Workflow
```bash
# 1. Kết nối đến servers trong ngày
ssh user@prod-web-01
ssh deploy@staging-api  
ssh admin@monitoring

# 2. Sau đó dùng SSH Manager để browse nhanh
sshm
# → Tất cả servers xuất hiện trong TUI
# → Chọn và kết nối chỉ với Enter
```

### Khác biệt với tools khác

| Tool | Approach | SSH Manager |
|------|----------|-------------|
| `ssh` | Manual typing | 🔍 Auto-discovery |
| `ssh-config` | Manual config | ⚡ Zero config |
| Complex tools | Many commands | 🎯 TUI-only |

## 🗑️ Gỡ cài đặt

```bash
# Xóa binary
sudo rm /usr/local/bin/sshm

# Xóa tất cả dữ liệu (tùy chọn)
rm -rf ~/.sshm/
```

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
- [SQLite](https://sqlite.org/) - Database nhẹ và tin cậy
- Go community - Ecosystem tuyệt vời

## 📞 Liên hệ

- GitHub: [@levanduy](https://github.com/levanduy)
- Email: your.email@example.com

---

⭐ **Nếu project hữu ích, hãy cho một star để ủng hộ!** ⭐ 