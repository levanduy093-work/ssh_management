package ssh

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RemoveFromKnownHosts removes a host from ~/.ssh/known_hosts file
func RemoveFromKnownHosts(hostname string, port int) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot access home directory: %v", err)
	}

	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	if _, err := os.Stat(knownHostsPath); err != nil {
		return nil // File doesn't exist, nothing to remove
	}

	// Read the file
	file, err := os.Open(knownHostsPath)
	if err != nil {
		return fmt.Errorf("cannot open known_hosts: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Create patterns to match the hostname
	patterns := []string{
		hostname,                               // hostname
		fmt.Sprintf("[%s]:%d", hostname, port), // [hostname]:port
	}

	for scanner.Scan() {
		line := scanner.Text()
		shouldRemove := false

		// Check if this line contains our hostname
		for _, pattern := range patterns {
			if strings.Contains(line, pattern) {
				shouldRemove = true
				break
			}
		}

		// Keep the line if it doesn't match our hostname
		if !shouldRemove {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading known_hosts: %v", err)
	}

	// Write back to file
	file.Close() // Close the read handle

	writeFile, err := os.Create(knownHostsPath)
	if err != nil {
		return fmt.Errorf("cannot write to known_hosts: %v", err)
	}
	defer writeFile.Close()

	for _, line := range lines {
		if _, err := writeFile.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to known_hosts: %v", err)
		}
	}

	return nil
}
