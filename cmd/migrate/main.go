package main

import (
	"build-migration/pkg/migrator"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	var (
		// Basic operations
		create = flag.String("create", "", "Create a new migration")
		up     = flag.Bool("up", false, "Run all pending migrations")
		down   = flag.Bool("down", false, "Rollback the last migration")
		status = flag.Bool("status", false, "Show migration status")
		reset  = flag.Bool("reset", false, "Rollback all migrations")

		// Migration creation options
		template = flag.String("template", "", "Template type (create_table, add_column, etc.)")
		table    = flag.String("table", "", "Table name")
		column   = flag.String("column", "", "Column name")
		colType  = flag.String("type", "", "Column type")

		// Database configuration
		driver = flag.String("driver", "postgresql", "Database driver (postgresql, mysql, sqlite)")
		dsn    = flag.String("dsn", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Phnom_Penh", "Database connection string")
	)
	flag.Parse()

	// Create configuration
	config := migrator.NewConfig(*driver, *dsn)

	// Initialize migrator
	m, err := migrator.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	// Handle commands
	switch {
	case *create != "":
		data := migrator.MigrationData{
			TableName:  *table,
			ColumnName: *column,
			ColumnType: *colType,
		}
		if err := m.CreateMigration(*create, *template, data); err != nil {
			log.Fatal(err)
		}

	case *up:
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully ran migrations")

	case *down:
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully rolled back migration")

	case *status:
		if err := m.Status(); err != nil {
			log.Fatal(err)
		}

	case *reset:
		if err := m.Reset(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully reset all migrations")

	default:
		flag.Usage()
		os.Exit(1)
	}
}
