package cmd

import (
	"database/sql"
	"log"
	"sort"

	"github.com/IsaacDSC/migrations/migration"
)

type (
	InsertMigration func(db *sql.Tx, version int) error
)

func Up(db *sql.DB, insertFn InsertMigration, dbVersion int, state []migration.Migrate) {
	if len(state[dbVersion:]) == 0 {
		log.Println("No migrations to apply")
		return
	}

	// order by version 1,2,3...
	sort.Slice(state[dbVersion:], func(i, j int) bool {
		return state[dbVersion:][i].Version < state[dbVersion:][j].Version
	})

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

}
