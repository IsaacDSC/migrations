package cfg

import "os"

func GetMigrationsPath() string {
	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		return "./migrations"
	}

	return path
}
