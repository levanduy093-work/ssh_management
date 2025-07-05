package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type KnownHost struct {
	Name     string
	Hostname string
	Username string
	Port     int
	Source   string // "known_hosts"
	KeyType  string // ssh-ed25519, ssh-rsa, etc.
}

// detectFromSSHFilesSilent runs detection in silent mode for auto-discovery
func detectFromSSHFilesSilent() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return // Silently fail
	}

	// Only scan known_hosts
	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	if _, err := os.Stat(knownHostsPath); err != nil {
		return // File doesn't exist, silently return
	}

	hosts := parseKnownHosts(knownHostsPath)
	if len(hosts) == 0 {
		return // No hosts found
	}

	// Remove duplicates and merge
	mergedHosts := mergeKnownHosts(hosts)

	// Import hosts silently
	for _, host := range mergedHosts {
		// Check if host already exists
		if existingHost, err := hostService.GetHostByName(host.Name); err == nil && existingHost != nil {
			continue // Skip existing hosts silently
		}

		// Add the host
		description := fmt.Sprintf("Auto-detected from %s (%s)", host.Source, host.KeyType)

		_, _ = hostService.CreateHost(
			host.Name,
			host.Hostname,
			host.Username,
			host.Port,
			"", // No specific key path from detection
			description,
			"ssh-detected",
		)
		// Ignore errors in silent mode
	}
}

func parseKnownHosts(knownHostsPath string) []KnownHost {
	file, err := os.Open(knownHostsPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var hosts []KnownHost
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse known_hosts line: hostname keytype publickey
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		hostname := parts[0]
		keyType := parts[1]

		// Skip hashed hostnames (starting with |1|)
		if strings.HasPrefix(hostname, "|1|") {
			continue
		}

		// Handle [hostname]:port format
		actualHostname := hostname
		port := 22

		if strings.HasPrefix(hostname, "[") && strings.Contains(hostname, "]:") {
			// Format: [hostname]:port
			re := regexp.MustCompile(`\[([^\]]+)\]:(\d+)`)
			matches := re.FindStringSubmatch(hostname)
			if len(matches) == 3 {
				actualHostname = matches[1]
				if p, err := parsePort(matches[2]); err == nil {
					port = p
				}
			}
		}

		// Generate name from hostname
		name := generateKnownHostName(actualHostname)

		host := KnownHost{
			Name:     name,
			Hostname: actualHostname,
			Username: getCurrentUsername(), // Default to current user
			Port:     port,
			Source:   "known_hosts",
			KeyType:  keyType,
		}

		hosts = append(hosts, host)
	}

	return hosts
}

func generateKnownHostName(hostname string) string {
	// Extract meaningful name from hostname
	parts := strings.Split(hostname, ".")
	if len(parts) > 0 {
		name := parts[0]
		// Remove any non-alphanumeric characters except hyphens
		name = regexp.MustCompile(`[^a-zA-Z0-9\-]`).ReplaceAllString(name, "")
		if name != "" {
			return name
		}
	}

	// Fallback to hostname with dots replaced by hyphens
	return strings.ReplaceAll(hostname, ".", "-")
}

func getCurrentUsername() string {
	if username := os.Getenv("USER"); username != "" {
		return username
	}
	if username := os.Getenv("USERNAME"); username != "" {
		return username
	}
	return "user" // Fallback
}

func parsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 22, err
	}
	return port, nil
}

func mergeKnownHosts(hosts []KnownHost) []KnownHost {
	hostMap := make(map[string]KnownHost)

	for _, host := range hosts {
		key := fmt.Sprintf("%s@%s:%d", host.Username, host.Hostname, host.Port)

		// If we already have this host, prefer the one with better key type
		if existing, exists := hostMap[key]; exists {
			// Prefer ed25519 over rsa
			if host.KeyType == "ssh-ed25519" && existing.KeyType != "ssh-ed25519" {
				hostMap[key] = host
			}
		} else {
			hostMap[key] = host
		}
	}

	// Convert back to slice
	var merged []KnownHost
	for _, host := range hostMap {
		merged = append(merged, host)
	}

	return merged
}
