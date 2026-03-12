package cmd

import (
	"database/sql"
	"log"
	"sort"

	"github.com/IsaacDSC/migrations/migration"
)

type (
	UpdateMigration func(db *sql.Tx, version int) error
)

func Down(db *sql.DB, updateFn UpdateMigration, dbVersion int, state []migration.Migrate) {
	// down somente da ultima migration aplicada
	if dbVersion == 0 {
		log.Println("No migrations to revert")
		return
	}

	if dbVersion > len(state) {
		log.Fatalf("db version %d exceeds known migrations (%d)", dbVersion, len(state))
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	// order by version 1,2,3...
	sort.Slice(state[dbVersion:], func(i, j int) bool {
		return state[dbVersion:][i].Version < state[dbVersion:][j].Version
	})

	migration := state[dbVersion-1]

	if migration.Version == 0 {
		log.Fatalf("migration %d not found", dbVersion)
	}

	if err := migration.Down(tx); err != nil {
		panic(err)
	}

	if err := updateFn(tx, migration.Version); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Migration %d reverted", migration.Version)
}
