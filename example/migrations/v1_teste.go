package main

import (
	"database/sql"

	"github.com/IsaacDSC/migrations"
	"github.com/IsaacDSC/migrations/migration"
)

func init() {
	migrations.State = append(migrations.State, migration.Migrate{
		Version: 1,
		Up: func(db *sql.Tx) error {
			_, err := db.Exec("CREATE TABLE IF NOT EXISTS teste (id SERIAL PRIMARY KEY, name VARCHAR(255))")
			return err
		},
		Down: func(db *sql.Tx) error {
			_, err := db.Exec("DROP TABLE IF EXISTS teste")
			return err
		},
	})
}
