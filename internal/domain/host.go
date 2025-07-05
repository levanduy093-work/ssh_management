package domain

import (
	"time"
)

// Host represents an SSH host configuration
type Host struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Hostname    string    `json:"hostname" db:"hostname"`
	Port        int       `json:"port" db:"port"`
	Username    string    `json:"username" db:"username"`
	KeyPath     string    `json:"key_path" db:"key_path"`
	Description string    `json:"description" db:"description"`
	Tags        string    `json:"tags" db:"tags"`
	LastUsed    time.Time `json:"last_used" db:"last_used"`
	UseCount    int       `json:"use_count" db:"use_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Repository interface for host operations
type HostRepository interface {
	Create(host *Host) error
	GetAll() ([]*Host, error)
	GetByID(id int) (*Host, error)
	GetByName(name string) (*Host, error)
	Update(host *Host) error
	Delete(id int) error
	Search(query string) ([]*Host, error)
	IncrementUseCount(id int) error
}

// Config represents application configuration
type Config struct {
	DatabasePath string `json:"database_path"`
	DefaultPort  int    `json:"default_port"`
	DefaultUser  string `json:"default_user"`
} 