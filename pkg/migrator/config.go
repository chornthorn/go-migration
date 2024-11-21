package migrator

import (
	"database/sql"
	"fmt"
)

type Config struct {
	Driver  string // e.g., "postgres", "mysql", "sqlite"
	DSN     string
	Dialect Dialect
}

func validateDriver(driver string) error {
	drivers := sql.Drivers()
	for _, d := range drivers {
		if d == driver {
			return nil
		}
	}
	return fmt.Errorf("unsupported driver: %s. Available drivers: %v", driver, drivers)
}

func New(config *Config) (*Migrator, error) {
	// Validate driver
	if err := validateDriver(config.Driver); err != nil {
		return nil, err
	}

	db, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Migrator{
		db:     db,
		config: config,
	}, nil
}

func NewConfig(driver, dsn string) *Config {
	// Normalize driver name
	switch driver {
	case "postgres", "postgresql":
		driver = "postgres" // Important: use "postgres" not "postgresql"
	case "mysql":
		driver = "mysql"
	case "sqlite", "sqlite3":
		driver = "sqlite3"
	default:
		panic(fmt.Sprintf("unsupported driver: %s", driver))
	}

	config := &Config{
		Driver: driver,
		DSN:    dsn,
	}

	switch driver {
	case "postgres":
		config.Dialect = &PostgresDialect{}
	case "mysql":
		config.Dialect = &MySQLDialect{}
	case "sqlite3":
		config.Dialect = &SQLiteDialect{}
	}

	return config
}

// GetMigrationsDir returns the directory for migrations based on the driver
func (c *Config) GetMigrationsDir() string {
	switch c.Driver {
	case "postgres", "postgresql":
		return "migrations/postgresql"
	case "mysql":
		return "migrations/mysql"
	case "sqlite", "sqlite3":
		return "migrations/sqlite"
	default:
		return "migrations/" + c.Driver
	}
}
