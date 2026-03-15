package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/IsaacDSC/migrations/migration"
)

type (
	InsertMigration func(db *sql.Tx, version int) error
)

func Up(db *sql.DB, insertFn InsertMigration, dbVersion int, state []migration.Migrate) error {
	if len(state[dbVersion:]) == 0 {
		log.Println("No migrations to apply")
		return nil
	}

	// order by version 1,2,3...
	sort.Slice(state[dbVersion:], func(i, j int) bool {
		return state[dbVersion:][i].Version < state[dbVersion:][j].Version
	})

	// if the versions are repeated, return an error
	for i := 0; i < len(state[dbVersion:])-1; i++ {
		if state[dbVersion:][i].Version == state[dbVersion:][i+1].Version {
			return fmt.Errorf("migration %d is repeated", state[dbVersion:][i].Version)
		}
	}

	for _, migration := range state[dbVersion:] {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		defer tx.Rollback()

		if err := migration.Up(tx); err != nil {
			panic(err)
		}

		if err := insertFn(tx, migration.Version); err != nil {
			panic(err)
		}

		if err := tx.Commit(); err != nil {
			panic(err)
		}

		log.Printf("Migration %d applied", migration.Version)
	}

	return nil
}
