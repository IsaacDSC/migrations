package database

import (
	"database/sql"
	"fmt"
)

func InsertMigration(tx *sql.Tx, version int) error {
	_, err := tx.Exec("INSERT INTO migrations (version, created_at) VALUES ($1, NOW())", version)
	if err != nil {
		return fmt.Errorf("failed to insert migration: %w", err)
	}

	return nil
}
