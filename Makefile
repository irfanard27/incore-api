# Incore API Makefile

.PHONY: help run build clean test backup restore migrate-up migrate-down dev

# Default target
help:
	@echo "Available commands:"
	@echo "  run         - Run the application in development mode"
	@echo "  build       - Build the application binary"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  backup      - Backup PostgreSQL database"
	@echo "  restore     - Restore PostgreSQL database from backup"
	@echo "  migrate-up  - Run database migrations up"
	@echo "  migrate-down- Run database migrations down"
	@echo "  dev         - Run in development with hot reload (if air is installed)"

# Variables
APP_NAME := incore-api
DB_HOST := localhost
DB_PORT := 5432
DB_USER := postgres
DB_PASSWORD := postgres
DB_NAME := incore-db
BACKUP_FILE := database_backup.sql

# Run the application
run:
	@echo "Starting $(APP_NAME)..."
	go run cmd/server/main.go

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) cmd/server/main.go

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	rm -f database_backup.sql

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Backup database
backup:
	@echo "Backing up database $(DB_NAME)..."
	PGPASSWORD="$(DB_PASSWORD)" pg_dump -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) > $(BACKUP_FILE)
	@echo "Backup saved to $(BACKUP_FILE)"

# Restore database
restore:
	@echo "Restoring database $(DB_NAME) from $(BACKUP_FILE)..."
	@if [ ! -f $(BACKUP_FILE) ]; then \
		echo "Backup file $(BACKUP_FILE) not found!"; \
		exit 1; \
	fi
	PGPASSWORD="$(DB_PASSWORD)" psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) < $(BACKUP_FILE)
	@echo "Database restored successfully"

# Run migrations up
migrate-up:
	@echo "Running database migrations up..."
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# Run migrations down
migrate-down:
	@echo "Running database migrations down..."
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Development with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	@if ! command -v air &> /dev/null; then \
		echo "air is not installed. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Skipping..."; \
	fi

# Create database
create-db:
	@echo "Creating database $(DB_NAME)..."
	PGPASSWORD="$(DB_PASSWORD)" createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
	@echo "Database $(DB_NAME) created successfully"

# Drop database (WARNING: This will delete all data)
drop-db:
	@echo "WARNING: This will delete database $(DB_NAME) and all its data!"
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	PGPASSWORD="$(DB_PASSWORD)" dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
	@echo "Database $(DB_NAME) dropped"

# Full setup (create db, migrate, and run)
setup: create-db migrate-up deps run

# Production build and run
prod: build
	@echo "Starting $(APP_NAME) in production mode..."
	./$(APP_NAME)
