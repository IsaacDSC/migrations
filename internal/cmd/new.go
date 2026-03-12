package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/IsaacDSC/migrations/internal/cfg"
	"github.com/IsaacDSC/migrations/migration"
)

func New(filename string, state []migration.Migrate) {
	pathMigration := cfg.GetMigrationsPath()

	if filename == "" {
		log.Fatal("filename is required, use: go run main.go new <filename>")
	}

	version := len(state) + 1
	path := fmt.Sprintf("%s/v%d_%s.go", pathMigration, version, filename)
	content := `
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
					return nil
				},
				Down: func(db *sql.Tx) error {
					return nil
				},
			})
		}

		`
	os.WriteFile(path, []byte(content), 0644)

	// execute go fmt ./migrations/...
	exec.Command("go", "fmt", fmt.Sprintf("%s/...", pathMigration)).Run()

	log.Printf("Migration %d created at %s", version, path)
}
