package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetVersion returns the last version of the migrations table
func GetVersion(db *sql.DB) (int, error) {
	var dbVersion int
	err := db.QueryRow("SELECT count(version) FROM migrations WHERE rollback = FALSE;").Scan(&dbVersion)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to get version: %w", err)
	}

	return dbVersion, nil
}
