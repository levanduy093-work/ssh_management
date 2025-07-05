package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/levanduy/ssh_management/internal/domain"
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

	host := &domain.Host{
		Name:        name,
		Hostname:    hostname,
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