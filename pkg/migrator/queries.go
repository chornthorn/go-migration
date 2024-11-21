package migrator

import (
	"database/sql"
	"fmt"
)

// getAppliedMigrations returns a map of already applied migrations
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migration rows: %w", err)
	}

	return applied, nil
}

// getAppliedMigrationsWithDetails returns detailed information about applied migrations
func (m *Migrator) getAppliedMigrationsWithDetails() ([]Migration, error) {
	var migrations []Migration

	query := "SELECT version, name, applied_at FROM schema_migrations ORDER BY version"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var migration Migration
		if err := rows.Scan(&migration.Version, &migration.Name, &migration.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration details: %w", err)
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

// GetLastAppliedMigration returns the last applied migration
func (m *Migrator) GetLastAppliedMigration() (*Migration, error) {
	var migration Migration

	query := `
        SELECT version, name, applied_at 
        FROM schema_migrations 
        ORDER BY version DESC 
        LIMIT 1
    `

	err := m.db.QueryRow(query).Scan(
		&migration.Version,
		&migration.Name,
		&migration.AppliedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last applied migration: %w", err)
	}

	return &migration, nil
}
