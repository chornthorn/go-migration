package migrator

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Status returns the current status of all migrations
func (m *Migrator) Status() error {
	files, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, file := range files {
		base := filepath.Base(file)
		version := base[:14]
		name := strings.TrimSuffix(base[15:], ".up.sql")

		status := "Pending"
		if applied[version] {
			status = "Applied"
		}

		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("Status: %s\n", status)
		fmt.Println("----------------")
	}

	return nil
}

// HasPendingMigrations checks if there are any pending migrations
func (m *Migrator) HasPendingMigrations() (bool, error) {
	files, err := m.getMigrationFiles()
	if err != nil {
		return false, err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return false, err
	}

	for _, file := range files {
		base := filepath.Base(file)
		version := base[:14]
		if !applied[version] {
			return true, nil
		}
	}

	return false, nil
}
