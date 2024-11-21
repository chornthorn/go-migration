package migrator

import "fmt"

type Dialect interface {
	CreateMigrationsTableSQL() string
	PlaceholderSQL(n int) string
	QuoteIdentifier(name string) string
}

// PostgresDialect implements Dialect for PostgreSQL
type PostgresDialect struct{}

func (d *PostgresDialect) CreateMigrationsTableSQL() string {
	return `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(14) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    )`
}

func (d *PostgresDialect) PlaceholderSQL(n int) string {
	return fmt.Sprintf("$%d", n)
}

func (d *PostgresDialect) QuoteIdentifier(name string) string {
	return `"` + name + `"`
}

// MySQLDialect implements Dialect for MySQL
type MySQLDialect struct{}

func (d *MySQLDialect) CreateMigrationsTableSQL() string {
	return `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(14) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
}

func (d *MySQLDialect) PlaceholderSQL(n int) string {
	return "?"
}

func (d *MySQLDialect) QuoteIdentifier(name string) string {
	return "`" + name + "`"
}

// SQLiteDialect implements Dialect for SQLite
type SQLiteDialect struct{}

func (d *SQLiteDialect) CreateMigrationsTableSQL() string {
	return `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(14) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
}

func (d *SQLiteDialect) PlaceholderSQL(n int) string {
	return "?"
}

func (d *SQLiteDialect) QuoteIdentifier(name string) string {
	return `"` + name + `"`
}
