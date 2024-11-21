# Database Migration Tool

A flexible database migration tool that supports multiple databases.

## Features

- Supports multiple databases (PostgreSQL, MySQL, SQLite)
- Version-controlled migrations
- Up/down migrations
- Migration status tracking
- Transaction support
- Template-based migration generation

## Usage

### Create a new migration

```bash
# Create a basic migration
go run cmd/migrate/main.go -create add_users -template create_table

# Create a column migration
go run cmd/migrate/main.go -create add_email \
    -template add_column \
    -table users \
    -column email \
    -type "VARCHAR(255)"