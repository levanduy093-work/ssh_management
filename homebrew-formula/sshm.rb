class Sshm < Formula
  desc "Terminal UI for SSH Host Management with auto-discovery and intelligent username detection"
  homepage "https://github.com/levanduy093-work/ssh_management"
  url "https://github.com/levanduy093-work/ssh_management/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "e9a0af78e737e0eaff9b4eeea802eb9bb8266bcb7f86f42bd3f06553f0b35523"
  license "MIT"
  head "https://github.com/levanduy093-work/ssh_management.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the application
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"sshm", "./cmd/sshm"
    
    # Install man page if it exists
    # man1.install "docs/sshm.1" if File.exist?("docs/sshm.1")
    
    # Install shell completions if they exist
    # bash_completion.install "completions/sshm.bash" if File.exist?("completions/sshm.bash")
    # zsh_completion.install "completions/_sshm" if File.exist?("completions/_sshm")
    # fish_completion.install "completions/sshm.fish" if File.exist?("completions/sshm.fish")
  end

  def post_install
    # Create data directory
    (var/"sshm").mkpath
    
    # Backup existing known_hosts if it exists
    if File.exist?(ENV["HOME"] + "/.ssh/known_hosts")
      cp ENV["HOME"] + "/.ssh/known_hosts", ENV["HOME"] + "/.ssh/known_hosts.backup.homebrew"
    end
  end

  test do
    # Test that the binary was installed and runs
    system bin/"sshm", "--help"
  end

  def caveats
    <<~EOS
      SSH Manager (sshm) has been installed!
      
      🚀 Quick Start:
        sshm
      
      📁 Data Storage:
        Database: ~/.sshm/hosts.db
        Backup: ~/.ssh/known_hosts.backup.homebrew
      
      🔍 Auto-Discovery:
        SSH Manager automatically discovers hosts from:
        • ~/.ssh/known_hosts
        • Shell history (~/.zsh_history, ~/.bash_history)
        • SSH config (~/.ssh/config)
      
      ⌨️  TUI Controls:
        ↑/↓ or j/k    Navigate
        Enter         Connect to host
        /             Search hosts
        x             Delete host
        r             Refresh/discover
        q             Quit
      
      For more information, visit:
      https://github.com/levanduy093-work/ssh_management
    EOS
  end
end 