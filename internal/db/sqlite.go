package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"remembrall/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite store
func NewSQLiteStore() (*SQLiteStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dbPath := filepath.Join(homeDir, ".remembrall.db")
	
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteStore{db: db}
	if err := store.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return store, nil
}

// createTables creates the necessary database tables
func (s *SQLiteStore) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app_name TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_app_name ON passwords(app_name);
	`

	_, err := s.db.Exec(query)
	return err
}

// Save stores a new password entry
func (s *SQLiteStore) Save(appName, password string) error {
	query := `
	INSERT INTO passwords (app_name, password, created_at, updated_at)
	VALUES (?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := s.db.Exec(query, appName, password, now, now)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("password for '%s' already exists, use 'update' command to modify it", appName)
		}
		return fmt.Errorf("failed to save password: %w", err)
	}
	
	return nil
}

// Get retrieves a password entry by app name
func (s *SQLiteStore) Get(appName string) (*models.PasswordEntry, error) {
	query := `
	SELECT id, app_name, password, created_at, updated_at
	FROM passwords
	WHERE app_name = ?
	`
	
	row := s.db.QueryRow(query, appName)
	
	var entry models.PasswordEntry
	err := row.Scan(&entry.ID, &entry.AppName, &entry.Password, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no password found for '%s'", appName)
		}
		return nil, fmt.Errorf("failed to retrieve password: %w", err)
	}
	
	return &entry, nil
}

// Update modifies an existing password entry
func (s *SQLiteStore) Update(appName, newPassword string) error {
	query := `
	UPDATE passwords
	SET password = ?, updated_at = ?
	WHERE app_name = ?
	`
	
	result, err := s.db.Exec(query, newPassword, time.Now(), appName)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	
	if affected == 0 {
		return fmt.Errorf("no password found for '%s'", appName)
	}
	
	return nil
}

// List returns all password entries (without decrypted passwords)
func (s *SQLiteStore) List() ([]*models.PasswordEntry, error) {
	query := `
	SELECT id, app_name, password, created_at, updated_at
	FROM passwords
	ORDER BY app_name
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list passwords: %w", err)
	}
	defer rows.Close()
	
	var entries []*models.PasswordEntry
	for rows.Next() {
		var entry models.PasswordEntry
		err := rows.Scan(&entry.ID, &entry.AppName, &entry.Password, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan password entry: %w", err)
		}
		entries = append(entries, &entry)
	}
	
	return entries, nil
}

// Search finds password entries that match the query (fuzzy search)
func (s *SQLiteStore) Search(query string) ([]*models.PasswordEntry, error) {
	sqlQuery := `
	SELECT id, app_name, password, created_at, updated_at
	FROM passwords
	WHERE app_name LIKE ?
	ORDER BY app_name
	`
	
	searchPattern := "%" + query + "%"
	rows, err := s.db.Query(sqlQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search passwords: %w", err)
	}
	defer rows.Close()
	
	var entries []*models.PasswordEntry
	for rows.Next() {
		var entry models.PasswordEntry
		err := rows.Scan(&entry.ID, &entry.AppName, &entry.Password, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan password entry: %w", err)
		}
		entries = append(entries, &entry)
	}
	
	return entries, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}