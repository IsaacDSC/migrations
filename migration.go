package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/IsaacDSC/migrations/internal/cmd"
	"github.com/IsaacDSC/migrations/internal/database"
	"github.com/IsaacDSC/migrations/migration"
)

var State []migration.Migrate

func Start(db *sql.DB) {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./migrations/ up|down|new [filename]")
	}

	if err := database.Plugin(db); err != nil {
		panic(err)
	}

	if err := database.CreateTable(db); err != nil {
		panic(err)
	}

	dbVersion, err := database.GetVersion(db)
	if err != nil {
		panic(err)
	}

	inputCmd := os.Args[1]

	switch inputCmd {
	case "new":
		if len(os.Args) < 3 {
			log.Fatalln("filename is required, use: go run main.go new <filename>")
		}

		cmd.New(os.Args[2], State)
	case "up":
		if err := cmd.Up(db, database.InsertMigration, dbVersion, State); err != nil {
			log.Fatal(err)
		}
	case "down":
		cmd.Down(db, database.UpdateMigration, dbVersion, State)
	case "help":
		cmd.Help()
	case "version":
		fmt.Printf("[*] Database version: %d\n", dbVersion)
		fmt.Printf("[*] Migrations version: %d\n", len(State))
	}
}
