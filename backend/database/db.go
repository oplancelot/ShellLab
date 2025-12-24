// Package database provides SQLite database operations for ShellLab
package database

import (
	"database/sql"
	"fmt"
	"sync"

	"shelllab/backend/database/schema"

	_ "modernc.org/sqlite"
)

// SQLiteDB wraps the SQLite database connection
type SQLiteDB struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Enable WAL mode for better performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &SQLiteDB{db: db}, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// InitSchema creates the database schema if it doesn't exist
func (s *SQLiteDB) InitSchema() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create core tables
	if _, err := s.db.Exec(schema.CoreSchema()); err != nil {
		return fmt.Errorf("failed to create core schema: %w", err)
	}

	// Create AtlasLoot tables
	if _, err := s.db.Exec(schema.AtlasLootSchema()); err != nil {
		return fmt.Errorf("failed to create atlasloot schema: %w", err)
	}

	// Create locale tables
	if _, err := s.db.Exec(schema.LocaleSchema()); err != nil {
		return fmt.Errorf("failed to create locale schema: %w", err)
	}

	return nil
}

// DB returns the underlying sql.DB for direct queries
func (s *SQLiteDB) DB() *sql.DB {
	return s.db
}
