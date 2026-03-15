package db

import (
	"database/sql"
	"fmt"
	"os"
)

// RunMigrations reads and executes SQL files in order.
func RunMigrations(db *sql.DB, files ...string) error {
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("exec %s: %w", f, err)
		}
		fmt.Printf("✓ Ran migration: %s\n", f)
	}
	return nil
}
