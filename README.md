# Migrations


## Commands

Create a new migration. This command generates a file at `./migrations/{version}_{name}.go`.
```sh
go run ./migrations/ new <filename>

```
#### Example output file
```go
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

```


### Apply migrations to the database

```sh
go run ./migrations/ up
```

### Roll back the last applied migration
```sh
go run ./migrations/ down
```

### Get migration version
This command returns the database version and the local migrations version.

```sh
go run ./migrations/ version
```
##### Example output:
```
[*] Database version: 1
[*] Migrations version: 1

```
### Show command usage
This command prints standard information on how to use the migrations CLI.

```sh
go run ./migrations/ help
```

## How to implement migrations in your project

### Download the library
```
go get -u github.com/IsaacDSC/migrations
```

### Create folder and migration file

#### Use postgres


```sh
mkdir -p "./migrations" && echo '
package main

import (
	"database/sql"

	"github.com/IsaacDSC/migrations"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://idsc:admin@localhost:5432/example?sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	migrations.Start(db)
}' >> ./migrations/main.go
```

#### Use mysql

```sh
mkdir -p "./migrations" && echo '
package main

import (
	"database/sql"

	"github.com/IsaacDSC/migrations"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "mysql://idsc:admin@localhost:5432/example?sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	migrations.Start(db)
}' >> ./migrations/main.go
```