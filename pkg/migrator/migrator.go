package migrator

import (
	"build-migration/pkg/migrator/templates"
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

type Migrator struct {
	db     *sql.DB
	config *Config
}

// New creates a new Migrator instance
//func New(config *Config) (*Migrator, error) {
//	db, err := sql.Open(config.Driver, config.DSN)
//	if err != nil {
//		return nil, fmt.Errorf("failed to connect to database: %w", err)
//	}
//
//	if err := db.Ping(); err != nil {
//		return nil, fmt.Errorf("failed to ping database: %w", err)
//	}
//
//	return &Migrator{
//		db:     db,
//		config: config,
//	}, nil
//}

// Close closes the database connection
func (m *Migrator) Close() error {
	return m.db.Close()
}

func (m *Migrator) previewSQL(content string) string {
	// Basic SQL formatting
	lines := strings.Split(content, "\n")
	var formatted []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			formatted = append(formatted, "    "+line)
		}
	}
	return strings.Join(formatted, "\n")
}

func (m *Migrator) CreateMigration(name, templateType string, data MigrationData) error {
	timestamp := time.Now().Format("20060102150405")
	data.Name = name
	data.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	upTmpl, downTmpl := templates.GetTemplatesByDriver(m.config.Driver, templateType)

	// Create up migration
	upFile := filepath.Join(m.config.GetMigrationsDir(), fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	if err := m.createFileFromTemplate(upFile, upTmpl, data); err != nil {
		return fmt.Errorf("failed to create up migration: %w", err)
	}

	// Create down migration
	downFile := filepath.Join(m.config.GetMigrationsDir(), fmt.Sprintf("%s_%s.down.sql", timestamp, name))
	if err := m.createFileFromTemplate(downFile, downTmpl, data); err != nil {
		return fmt.Errorf("failed to create down migration: %w", err)
	}

	// Print detailed information about the created migration
	fmt.Printf("\nCreated new migration for table '%s':\n", data.TableName)
	fmt.Printf("Operation: %s\n", templateType)
	switch templateType {
	case "add_column", "add_column_with_default", "add_column_with_fk":
		fmt.Printf("Column: %s %s\n", data.ColumnName, data.ColumnType)
		if data.DefaultValue != "" {
			fmt.Printf("Default: %s\n", data.DefaultValue)
		}
		if data.ReferenceTable != "" {
			fmt.Printf("References: %s(%s)\n", data.ReferenceTable, data.ReferenceColumn)
		}
	case "create_table":
		fmt.Printf("New table will be created with default columns (id, created_at, updated_at)\n")
	}

	// Preview SQL
	//fmt.Printf("\nSQL Preview:\n")
	//fmt.Printf("Up Migration:\n%s\n", m.previewSQL(upContent))
	//fmt.Printf("\nDown Migration:\n%s\n", m.previewSQL(downContent))

	fmt.Printf("\nFiles created:\n")
	fmt.Printf("  Up:   %s\n", upFile)
	fmt.Printf("  Down: %s\n", downFile)
	fmt.Printf("\nTo run this migration:\n")
	fmt.Printf("  go run cmd/migrate/main.go -up\n")
	return nil
}

// execInTransaction executes a function within a transaction
func (m *Migrator) execInTransaction(fn func(*sql.Tx) error) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	if err := m.InitSchema(); err != nil {
		return err
	}

	files, err := m.getMigrationFiles()
	if err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	for _, file := range files {
		base := filepath.Base(file)
		version := base[:14]
		name := strings.TrimSuffix(base[15:], ".up.sql")

		if !applied[version] {
			fmt.Printf("Applying migration %s: %s\n", version, name)

			content, err := m.getMigrationContent(file)
			if err != nil {
				return err
			}

			err = m.execInTransaction(func(tx *sql.Tx) error {
				// Execute the up migration
				if _, err := tx.Exec(content); err != nil {
					return fmt.Errorf("failed to execute up migration: %w", err)
				}

				// Insert into schema_migrations
				query := fmt.Sprintf(
					"INSERT INTO schema_migrations (version, name) VALUES (%s, %s)",
					m.config.Dialect.PlaceholderSQL(1),
					m.config.Dialect.PlaceholderSQL(2),
				)
				if _, err := tx.Exec(query, version, name); err != nil {
					return fmt.Errorf("failed to record migration: %w", err)
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", version, err)
			}

			fmt.Printf("Successfully applied migration %s\n", version)
		}
	}

	return nil
}

// Down Then you can use it like this:
func (m *Migrator) Down() error {
	lastMigration, err := m.GetLastAppliedMigration()
	if err != nil {
		return err
	}
	if lastMigration == nil {
		return fmt.Errorf("no migrations to roll back")
	}

	downFile := filepath.Join(
		m.config.GetMigrationsDir(),
		fmt.Sprintf("%s_%s.down.sql", lastMigration.Version, lastMigration.Name),
	)

	content, err := m.getMigrationContent(downFile)
	if err != nil {
		return err
	}

	err = m.execInTransaction(func(tx *sql.Tx) error {
		// Execute the down migration
		if _, err := tx.Exec(content); err != nil {
			return fmt.Errorf("failed to execute down migration: %w", err)
		}

		// Delete from schema_migrations
		query := fmt.Sprintf(
			"DELETE FROM schema_migrations WHERE version = %s",
			m.config.Dialect.PlaceholderSQL(1),
		)
		if _, err := tx.Exec(query, lastMigration.Version); err != nil {
			return fmt.Errorf("failed to delete migration record: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to roll back migration %s: %w", lastMigration.Version, err)
	}

	fmt.Printf("Successfully rolled back migration %s: %s\n", lastMigration.Version, lastMigration.Name)
	return nil
}

// Reset rolls back all migrations
func (m *Migrator) Reset() error {
	for {
		if err := m.Down(); err != nil {
			if err.Error() == "no migrations to roll back" {
				break
			}
			return err
		}
	}
	return nil
}

// Refresh rolls back all migrations and runs them again
func (m *Migrator) Refresh() error {
	if err := m.Reset(); err != nil {
		return err
	}
	return m.Up()
}

// InitSchema initializes the schema_migrations table
func (m *Migrator) InitSchema() error {
	return m.execInTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(m.config.Dialect.CreateMigrationsTableSQL())
		if err != nil {
			return fmt.Errorf("failed to create migrations table: %w", err)
		}
		return nil
	})
}

func (m *Migrator) createFileFromTemplate(filename, tmplContent string, data MigrationData) error {
	tmpl, err := template.New(filepath.Base(filename)).Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return m.createMigrationFile(filename, buf.Bytes())
}
