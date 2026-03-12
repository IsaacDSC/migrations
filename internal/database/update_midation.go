package database

import (
	"database/sql"
	"fmt"
)

func UpdateMigration(tx *sql.Tx, version int) error {
	_, err := tx.Exec("UPDATE migrations SET rollback = true WHERE version = $1", version)
	if err != nil {
		return fmt.Errorf("failed to update migration: %w", err)
	}

	return nil
}
