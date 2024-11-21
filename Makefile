.PHONY: build test clean help

BINARY_NAME=migrate
BUILD_DIR=bin

help: ## Show this help
    @echo "Available commands:"
    @grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the migration tool
    @mkdir -p $(BUILD_DIR)
    @go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/migrate/main.go

install: build ## Install the migration tool
    @cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

test: ## Run tests
    @go test ./... -v

clean: ## Clean build directory
    @rm -rf $(BUILD_DIR)

# Migration commands
create: ## Create a new migration (usage: make create name=migration_name [template=template_type])
    @./$(BUILD_DIR)/$(BINARY_NAME) -create $(name) $(if $(template),-template $(template))

up: ## Run all pending migrations
    @./$(BUILD_DIR)/$(BINARY_NAME) -up

down: ## Rollback the last migration
    @./$(BUILD_DIR)/$(BINARY_NAME) -down

status: ## Show migration status
    @./$(BUILD_DIR)/$(BINARY_NAME) -status

reset: ## Reset all migrations
    @./$(BUILD_DIR)/$(BINARY_NAME) -reset

refresh: ## Refresh all migrations (reset + up)
    @./$(BUILD_DIR)/$(BINARY_NAME) -refresh