package database

import (
	"database/sql"
	"fmt"
)

// Plugin add plugin gen_random_uuid
func Plugin(db *sql.DB) error {
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err != nil {
		return fmt.Errorf("failed to create extension uuid-ossp: %w", err)
	}

	return nil
}

// CreateTable creates the migrations table if it doesn't exist
func CreateTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), version INT NOT NULL, rollback BOOLEAN NOT NULL DEFAULT FALSE, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to create table migrations: %w", err)
	}

	return nil
}
