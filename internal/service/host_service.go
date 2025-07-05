package service

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/levanduy/ssh_management/internal/domain"
	"github.com/levanduy/ssh_management/pkg/ssh"
)

type HostService struct {
	repo domain.HostRepository
}

func NewHostService(repo domain.HostRepository) *HostService {
	return &HostService{repo: repo}
}

func (s *HostService) CreateHost(name, hostname, username string, port int, keyPath, description, tags string) (*domain.Host, error) {
	if name == "" || hostname == "" || username == "" {
		return nil, fmt.Errorf("name, hostname and username are required")
	}

	if port <= 0 || port > 65535 {
		port = 22 // Default SSH port
	}

	// Validate key path if provided
	if keyPath != "" {
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("SSH key file does not exist: %s", keyPath)
		}
	}

	// Resolve IP address
	ipAddress := s.resolveIPAddress(hostname)

	host := &domain.Host{
		Name:        name,
		Hostname:    hostname,
		IPAddress:   ipAddress,
		Port:        port,
		Username:    username,
		KeyPath:     keyPath,
		Description: description,
		Tags:        tags,
	}

	if err := s.repo.Create(host); err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	return host, nil
}

func (s *HostService) GetAllHosts() ([]*domain.Host, error) {
	return s.repo.GetAll()
}

func (s *HostService) GetHostByID(id int) (*domain.Host, error) {
	return s.repo.GetByID(id)
}

func (s *HostService) GetHostByName(name string) (*domain.Host, error) {
	return s.repo.GetByName(name)
}

func (s *HostService) UpdateHost(host *domain.Host) error {
	if host.Name == "" || host.Hostname == "" || host.Username == "" {
		return fmt.Errorf("name, hostname and username are required")
	}

	if host.Port <= 0 || host.Port > 65535 {
		host.Port = 22
	}

	// Validate key path if provided
	if host.KeyPath != "" {
		if _, err := os.Stat(host.KeyPath); os.IsNotExist(err) {
			return fmt.Errorf("SSH key file does not exist: %s", host.KeyPath)
		}
	}

	return s.repo.Update(host)
}

func (s *HostService) DeleteHost(id int) error {
	return s.repo.Delete(id)
}

// DeleteHostFromBoth deletes host from both database and known_hosts file
func (s *HostService) DeleteHostFromBoth(id int) error {
	// First get the host details
	host, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get host: %w", err)
	}

	// Delete from database
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete from database: %w", err)
	}

	// Delete from known_hosts
	if err := ssh.RemoveFromKnownHosts(host.Hostname, host.Port); err != nil {
		return fmt.Errorf("failed to remove from known_hosts: %w", err)
	}

	return nil
}

func (s *HostService) SearchHosts(query string) ([]*domain.Host, error) {
	if query == "" {
		return s.repo.GetAll()
	}
	return s.repo.Search(query)
}

func (s *HostService) ConnectToHost(id int) error {
	// Increment use count
	return s.repo.IncrementUseCount(id)
}

// GetDefaultConfigPath returns the default configuration directory
func GetDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".sshm"
	}
	return filepath.Join(homeDir, ".sshm")
}

// GetDefaultDatabasePath returns the default database file path
func GetDefaultDatabasePath() string {
	return filepath.Join(GetDefaultConfigPath(), "hosts.db")
}

// ParseTags splits comma-separated tags into a slice
func ParseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}

	var result []string
	for _, tag := range strings.Split(tags, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// JoinTags joins a slice of tags into a comma-separated string
func JoinTags(tags []string) string {
	var cleanTags []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			cleanTags = append(cleanTags, tag)
		}
	}
	return strings.Join(cleanTags, ", ")
}

// AutoDiscoverFromKnownHosts discovers new SSH hosts from ~/.ssh/known_hosts
func (s *HostService) AutoDiscoverFromKnownHosts() (int, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, fmt.Errorf("cannot access home directory: %v", err)
	}

	// Only scan known_hosts
	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	if _, err := os.Stat(knownHostsPath); err != nil {
		return 0, nil // File doesn't exist, return 0 new hosts
	}

	hosts := s.parseKnownHosts(knownHostsPath)
	if len(hosts) == 0 {
		return 0, nil // No hosts found
	}

	// Remove duplicates and merge
	mergedHosts := s.mergeKnownHosts(hosts)

	newHostsCount := 0
	// Import hosts
	for _, host := range mergedHosts {
		// Check if host already exists
		if existingHost, err := s.GetHostByName(host.Name); err == nil && existingHost != nil {
			// Host exists, check if we have better information
			shouldUpdate := false

			// Update username if current one is just system username and we found a better one
			if existingHost.Username == s.getCurrentUsername() && host.Username != s.getCurrentUsername() {
				existingHost.Username = host.Username
				shouldUpdate = true
			}

			// Update IP if we don't have one or found a better one
			if existingHost.IPAddress == "" {
				ipAddress := s.resolveIPAddress(existingHost.Hostname)
				if ipAddress != "" {
					existingHost.IPAddress = ipAddress
					shouldUpdate = true
				}
			}

			if shouldUpdate {
				s.UpdateHost(existingHost)
			}

			continue // Skip existing hosts
		}

		// Add the host
		description := fmt.Sprintf("Auto-detected from %s (%s)", host.Source, host.KeyType)

		_, err := s.CreateHost(
			host.Name,
			host.Hostname,
			host.Username,
			host.Port,
			"", // No specific key path from detection
			description,
			"ssh-detected",
		)
		if err == nil {
			newHostsCount++
		}
	}

	return newHostsCount, nil
}

type KnownHost struct {
	Name     string
	Hostname string
	Username string
	Port     int
	Source   string // "known_hosts"
	KeyType  string // ssh-ed25519, ssh-rsa, etc.
}

func (s *HostService) parseKnownHosts(knownHostsPath string) []KnownHost {
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
				if p, err := s.parsePort(matches[2]); err == nil {
					port = p
				}
			}
		}

		// Generate name from hostname
		name := s.generateKnownHostName(actualHostname)

		host := KnownHost{
			Name:     name,
			Hostname: actualHostname,
			Username: s.parseSSHConfig(actualHostname), // Try to get username from SSH config
			Port:     port,
			Source:   "known_hosts",
			KeyType:  keyType,
		}

		hosts = append(hosts, host)
	}

	return hosts
}

func (s *HostService) generateKnownHostName(hostname string) string {
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

func (s *HostService) getCurrentUsername() string {
	if username := os.Getenv("USER"); username != "" {
		return username
	}
	if username := os.Getenv("USERNAME"); username != "" {
		return username
	}
	return "user" // Fallback
}

func (s *HostService) parsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 22, err
	}
	return port, nil
}

func (s *HostService) mergeKnownHosts(hosts []KnownHost) []KnownHost {
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

// parseShellHistory tries to find SSH commands from shell history to get username
func (s *HostService) parseShellHistory(hostname string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return s.getCurrentUsername()
	}

	// Try different shell history files
	historyFiles := []string{
		filepath.Join(homeDir, ".zsh_history"),
		filepath.Join(homeDir, ".bash_history"),
		filepath.Join(homeDir, ".history"),
	}

	for _, historyFile := range historyFiles {
		if username := s.parseHistoryFile(historyFile, hostname); username != "" {
			return username
		}
	}

	return s.getCurrentUsername() // Fallback
}

func (s *HostService) parseHistoryFile(historyFile, hostname string) string {
	file, err := os.Open(historyFile)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Look for SSH commands
		if username := s.extractUsernameFromSSHCommand(line, hostname); username != "" {
			return username
		}
	}

	return ""
}

func (s *HostService) extractUsernameFromSSHCommand(command, hostname string) string {
	// Remove timestamp prefix from zsh history (: 1234567890:0;ssh ...)
	if strings.Contains(command, ";") {
		parts := strings.Split(command, ";")
		if len(parts) > 1 {
			command = strings.Join(parts[1:], ";")
		}
	}

	command = strings.TrimSpace(command)

	// Look for SSH commands
	if !strings.HasPrefix(command, "ssh ") {
		return ""
	}

	// Parse SSH command arguments
	args := strings.Fields(command)
	if len(args) < 2 {
		return ""
	}

	// Look for user@host pattern
	for _, arg := range args[1:] {
		// Skip flags
		if strings.HasPrefix(arg, "-") {
			continue
		}

		// Check if this argument contains @ (user@host format)
		if strings.Contains(arg, "@") {
			parts := strings.Split(arg, "@")
			if len(parts) >= 2 {
				username := parts[0]
				hostPart := strings.Join(parts[1:], "@")

				// More flexible hostname matching
				if hostPart == hostname ||
					strings.Contains(hostPart, hostname) ||
					strings.Contains(hostname, hostPart) {
					return username
				}
			}
		} else {
			// Also check direct hostname matches (ssh hostname)
			if arg == hostname ||
				strings.Contains(arg, hostname) ||
				strings.Contains(hostname, arg) {
				// No username specified, check previous args for -l flag
				for i := 1; i < len(args)-1; i++ {
					if args[i] == "-l" && i+1 < len(args) {
						return args[i+1]
					}
				}
			}
		}
	}

	return ""
}

// Enhanced parseSSHConfig that also tries shell history
func (s *HostService) parseSSHConfig(hostname string) string {
	// First try SSH config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return s.parseShellHistory(hostname)
	}

	configPath := filepath.Join(homeDir, ".ssh", "config")
	file, err := os.Open(configPath)
	if err != nil {
		// SSH config doesn't exist, try shell history
		return s.parseShellHistory(hostname)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentHost string
	var currentUser string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse SSH config directives
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		directive := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch directive {
		case "host":
			currentHost = value
			currentUser = "" // Reset user for new host
		case "hostname":
			if currentHost != "" && (value == hostname || currentHost == hostname) {
				// We found a matching host section
			}
		case "user":
			if currentHost != "" && currentUser == "" {
				currentUser = value
			}
		}

		// Check if we have a match
		if currentHost != "" && currentUser != "" {
			// Check if this host matches our hostname
			if currentHost == hostname || strings.Contains(currentHost, hostname) {
				return currentUser
			}
		}
	}

	// SSH config didn't have the info, try shell history
	return s.parseShellHistory(hostname)
}

// resolveIPAddress tries to resolve hostname to IP address
func (s *HostService) resolveIPAddress(hostname string) string {
	// If hostname is already an IP address, return it
	if net.ParseIP(hostname) != nil {
		return hostname
	}

	// First try to find IP from known_hosts (might have direct IP entries)
	if ip := s.getIPFromKnownHosts(hostname); ip != "" {
		return ip
	}

	// Try to resolve hostname to IP
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "" // Return empty if resolution fails
	}

	// Return the first IPv4 address found
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}

	// Return the first IPv6 address if no IPv4 found
	if len(ips) > 0 {
		return ips[0].String()
	}

	return "" // Return empty if no IP found
}

// getIPFromKnownHosts checks if known_hosts has direct IP entries for this hostname
func (s *HostService) getIPFromKnownHosts(hostname string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	file, err := os.Open(knownHostsPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		knownHost := parts[0]

		// Skip hashed hostnames
		if strings.HasPrefix(knownHost, "|1|") {
			continue
		}

		// Check if this entry is for our hostname and contains an IP
		if strings.Contains(knownHost, hostname) {
			// Extract IP if present (could be hostname,ip format)
			if strings.Contains(knownHost, ",") {
				hostParts := strings.Split(knownHost, ",")
				for _, part := range hostParts {
					if net.ParseIP(part) != nil {
						return part
					}
				}
			}
			// Check if the hostname itself is an IP
			if net.ParseIP(knownHost) != nil {
				return knownHost
			}
		}
	}

	return ""
}
