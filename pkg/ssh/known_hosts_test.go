package ssh

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveFromKnownHostsExactMatch(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	sshDir := filepath.Join(homeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("mkdir .ssh: %v", err)
	}

	knownHostsPath := filepath.Join(sshDir, "known_hosts")
	content := `example.com ssh-ed25519 AAAA
example.com2 ssh-ed25519 BBBB
[example.com]:22 ssh-rsa CCCC
example.com,192.168.1.10 ssh-ed25519 DDDD
`
	if err := os.WriteFile(knownHostsPath, []byte(content), 0600); err != nil {
		t.Fatalf("write known_hosts: %v", err)
	}

	if err := RemoveFromKnownHosts("example.com", 22); err != nil {
		t.Fatalf("RemoveFromKnownHosts: %v", err)
	}

	lines := readNonEmptyLines(t, knownHostsPath)
	if len(lines) != 1 || lines[0] != "example.com2 ssh-ed25519 BBBB" {
		t.Fatalf("unexpected known_hosts lines: %v", lines)
	}
}

func TestRemoveFromKnownHostsPortSpecific(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	sshDir := filepath.Join(homeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("mkdir .ssh: %v", err)
	}

	knownHostsPath := filepath.Join(sshDir, "known_hosts")
	content := `example.com ssh-ed25519 AAAA
[example.com]:2222 ssh-rsa BBBB
`
	if err := os.WriteFile(knownHostsPath, []byte(content), 0600); err != nil {
		t.Fatalf("write known_hosts: %v", err)
	}

	if err := RemoveFromKnownHosts("example.com", 2222); err != nil {
		t.Fatalf("RemoveFromKnownHosts: %v", err)
	}

	lines := readNonEmptyLines(t, knownHostsPath)
	if len(lines) != 1 || lines[0] != "example.com ssh-ed25519 AAAA" {
		t.Fatalf("unexpected known_hosts lines: %v", lines)
	}
}

func readNonEmptyLines(t *testing.T, path string) []string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan file: %v", err)
	}
	return lines
}
