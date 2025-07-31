package models

import "time"

// PasswordEntry represents a stored password entry
type PasswordEntry struct {
	ID          int       `db:"id"`
	AppName     string    `db:"app_name"`
	Password    string    `db:"password"` // This will be encrypted
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// PasswordStore defines the interface for password storage operations
type PasswordStore interface {
	Save(appName, password string) error
	Get(appName string) (*PasswordEntry, error)
	Update(appName, newPassword string) error
	List() ([]*PasswordEntry, error)
	Search(query string) ([]*PasswordEntry, error)
	Close() error
}