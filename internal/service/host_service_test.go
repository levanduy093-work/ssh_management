package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseKnownHosts(t *testing.T) {
	t.Setenv("USER", "testuser")
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	sshDir := filepath.Join(homeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("mkdir .ssh: %v", err)
	}

	knownHostsPath := filepath.Join(sshDir, "known_hosts")
	content := `# comment
example.com ssh-ed25519 AAAA
[host.local]:2222 ssh-rsa BBBB
host.local,192.168.1.10 ssh-ed25519 CCCC
|1|hash ssh-rsa DDDD
`
	if err := os.WriteFile(knownHostsPath, []byte(content), 0600); err != nil {
		t.Fatalf("write known_hosts: %v", err)
	}

	svc := &HostService{}
	hosts := svc.parseKnownHosts(knownHostsPath)

	if len(hosts) != 3 {
		t.Fatalf("expected 3 hosts, got %d", len(hosts))
	}

	assertHost := func(hostname string, port int, name string) {
		t.Helper()
		for _, host := range hosts {
			if host.Hostname == hostname && host.Port == port {
				if host.Name != name {
					t.Fatalf("expected name %q for %s:%d, got %q", name, hostname, port, host.Name)
				}
				if host.Username != "testuser" {
					t.Fatalf("expected username %q for %s:%d, got %q", "testuser", hostname, port, host.Username)
				}
				return
			}
		}
		t.Fatalf("host not found: %s:%d", hostname, port)
	}

	assertHost("example.com", 22, "example")
	assertHost("host.local", 2222, "host")
	assertHost("host.local", 22, "host")
}

func TestParseSSHConfig(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USER", "fallback")

	sshDir := filepath.Join(homeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("mkdir .ssh: %v", err)
	}

	configPath := filepath.Join(sshDir, "config")
	config := `Host *.example.com !bad.example.com
  User alice
Host bad.example.com
  User bob
Host alias
  HostName real.host
  User dan
Host test
  User carol
`
	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	svc := &HostService{}

	if got := svc.parseSSHConfig("good.example.com"); got != "alice" {
		t.Fatalf("expected user alice, got %q", got)
	}
	if got := svc.parseSSHConfig("bad.example.com"); got != "bob" {
		t.Fatalf("expected user bob, got %q", got)
	}
	if got := svc.parseSSHConfig("real.host"); got != "dan" {
		t.Fatalf("expected user dan, got %q", got)
	}
	if got := svc.parseSSHConfig("test"); got != "carol" {
		t.Fatalf("expected user carol, got %q", got)
	}
}
