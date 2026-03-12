package migration

import "database/sql"

type Migrate struct {
	Version int
	Up      func(db *sql.Tx) error
	Down    func(db *sql.Tx) error
}
