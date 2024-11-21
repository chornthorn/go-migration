package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// getMigrationFiles returns all migration files sorted by version
func (m *Migrator) getMigrationFiles() ([]string, error) {
	migrationsDir := m.config.GetMigrationsDir()
	pattern := filepath.Join(migrationsDir, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration files: %w", err)
	}

	sort.Strings(files)
	return files, nil
}

// createMigrationFile creates a new migration file
func (m *Migrator) createMigrationFile(filename string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	if err := os.WriteFile(filename, content, 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	return nil
}

// getMigrationContent reads the content of a migration file
func (m *Migrator) getMigrationContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read migration file: %w", err)
	}
	return string(content), nil
}
