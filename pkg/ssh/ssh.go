package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/levanduy/ssh_management/internal/domain"
)

// ConnectToHost executes SSH connection using the system's SSH client
func ConnectToHost(host *domain.Host) error {
	args := buildSSHArgs(host)
	
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// BuildSSHCommand returns the SSH command as a string
func BuildSSHCommand(host *domain.Host) string {
	args := buildSSHArgs(host)
	return "ssh " + strings.Join(args, " ")
}

// buildSSHArgs constructs SSH command arguments
func buildSSHArgs(host *domain.Host) []string {
	var args []string

	// Add port if not default
	if host.Port != 22 {
		args = append(args, "-p", strconv.Itoa(host.Port))
	}

	// Add SSH key if specified
	if host.KeyPath != "" {
		args = append(args, "-i", host.KeyPath)
	}

	// Add the connection string
	connectionString := fmt.Sprintf("%s@%s", host.Username, host.Hostname)
	args = append(args, connectionString)

	return args
}

// ValidateSSHKey checks if the SSH key file exists and has proper permissions
func ValidateSSHKey(keyPath string) error {
	if keyPath == "" {
		return nil // No key specified, which is fine
	}

	info, err := os.Stat(keyPath)
	if err != nil {
		return fmt.Errorf("SSH key file does not exist: %s", keyPath)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("SSH key path is not a regular file: %s", keyPath)
	}

	// Check permissions (should be readable by owner only for private keys)
	mode := info.Mode()
	if mode&0077 != 0 {
		fmt.Printf("Warning: SSH key file %s has overly permissive permissions (%o)\n", keyPath, mode&0777)
		fmt.Printf("Consider running: chmod 600 %s\n", keyPath)
	}

	return nil
}

// TestConnection tests if we can connect to the host without executing commands
func TestConnection(host *domain.Host) error {
	args := buildSSHArgs(host)
	args = append(args, "-o", "ConnectTimeout=5", "-o", "BatchMode=yes", "exit")

	cmd := exec.Command("ssh", args...)
	return cmd.Run()
} 