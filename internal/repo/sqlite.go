package repo

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/levanduy/ssh_management/internal/domain"
	_ "modernc.org/sqlite"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(dbPath string) (*SQLiteRepo, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := ensureDir(dir); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &SQLiteRepo{db: db}
	if err := repo.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

func (r *SQLiteRepo) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS hosts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		hostname TEXT NOT NULL,
		port INTEGER DEFAULT 22,
		username TEXT NOT NULL,
		key_path TEXT DEFAULT '',
		description TEXT DEFAULT '',
		tags TEXT DEFAULT '',
		last_used DATETIME DEFAULT CURRENT_TIMESTAMP,
		use_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepo) Create(host *domain.Host) error {
	now := time.Now()
	host.CreatedAt = now
	host.UpdatedAt = now

	query := `
	INSERT INTO hosts (name, hostname, port, username, key_path, description, tags, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		host.Name, host.Hostname, host.Port, host.Username,
		host.KeyPath, host.Description, host.Tags,
		host.CreatedAt, host.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create host: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	host.ID = int(id)
	return nil
}

func (r *SQLiteRepo) GetAll() ([]*domain.Host, error) {
	query := `
	SELECT id, name, hostname, port, username, key_path, description, tags,
		   last_used, use_count, created_at, updated_at
	FROM hosts ORDER BY last_used DESC, name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query hosts: %w", err)
	}
	defer rows.Close()

	var hosts []*domain.Host
	for rows.Next() {
		host := &domain.Host{}
		err := rows.Scan(
			&host.ID, &host.Name, &host.Hostname, &host.Port,
			&host.Username, &host.KeyPath, &host.Description, &host.Tags,
			&host.LastUsed, &host.UseCount, &host.CreatedAt, &host.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan host: %w", err)
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (r *SQLiteRepo) GetByID(id int) (*domain.Host, error) {
	query := `
	SELECT id, name, hostname, port, username, key_path, description, tags,
		   last_used, use_count, created_at, updated_at
	FROM hosts WHERE id = ?
	`
	host := &domain.Host{}
	err := r.db.QueryRow(query, id).Scan(
		&host.ID, &host.Name, &host.Hostname, &host.Port,
		&host.Username, &host.KeyPath, &host.Description, &host.Tags,
		&host.LastUsed, &host.UseCount, &host.CreatedAt, &host.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("host with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get host by id: %w", err)
	}

	return host, nil
}

func (r *SQLiteRepo) GetByName(name string) (*domain.Host, error) {
	query := `
	SELECT id, name, hostname, port, username, key_path, description, tags,
		   last_used, use_count, created_at, updated_at
	FROM hosts WHERE name = ?
	`
	host := &domain.Host{}
	err := r.db.QueryRow(query, name).Scan(
		&host.ID, &host.Name, &host.Hostname, &host.Port,
		&host.Username, &host.KeyPath, &host.Description, &host.Tags,
		&host.LastUsed, &host.UseCount, &host.CreatedAt, &host.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("host with name '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get host by name: %w", err)
	}

	return host, nil
}

func (r *SQLiteRepo) Update(host *domain.Host) error {
	host.UpdatedAt = time.Now()

	query := `
	UPDATE hosts SET 
		name = ?, hostname = ?, port = ?, username = ?, 
		key_path = ?, description = ?, tags = ?, updated_at = ?
	WHERE id = ?
	`
	_, err := r.db.Exec(query,
		host.Name, host.Hostname, host.Port, host.Username,
		host.KeyPath, host.Description, host.Tags, host.UpdatedAt,
		host.ID)

	if err != nil {
		return fmt.Errorf("failed to update host: %w", err)
	}

	return nil
}

func (r *SQLiteRepo) Delete(id int) error {
	query := `DELETE FROM hosts WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete host: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("host with id %d not found", id)
	}

	return nil
}

func (r *SQLiteRepo) Search(query string) ([]*domain.Host, error) {
	searchQuery := `
	SELECT id, name, hostname, port, username, key_path, description, tags,
		   last_used, use_count, created_at, updated_at
	FROM hosts 
	WHERE name LIKE ? OR hostname LIKE ? OR description LIKE ? OR tags LIKE ?
	ORDER BY last_used DESC, name ASC
	`
	pattern := "%" + strings.ToLower(query) + "%"
	rows, err := r.db.Query(searchQuery, pattern, pattern, pattern, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search hosts: %w", err)
	}
	defer rows.Close()

	var hosts []*domain.Host
	for rows.Next() {
		host := &domain.Host{}
		err := rows.Scan(
			&host.ID, &host.Name, &host.Hostname, &host.Port,
			&host.Username, &host.KeyPath, &host.Description, &host.Tags,
			&host.LastUsed, &host.UseCount, &host.CreatedAt, &host.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan host: %w", err)
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (r *SQLiteRepo) IncrementUseCount(id int) error {
	query := `
	UPDATE hosts SET 
		use_count = use_count + 1,
		last_used = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to increment use count: %w", err)
	}

	return nil
}

func (r *SQLiteRepo) Close() error {
	return r.db.Close()
}

// Helper function to ensure directory exists
func ensureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
} 