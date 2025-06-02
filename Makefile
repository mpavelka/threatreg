# Go CLI Application Makefile

.PHONY: build run clean test deps help migrate-apply-sqlite migrate-apply-postgres migrate-status-sqlite migrate-status-postgres migrate-validate-sqlite migrate-validate-postgres migrate-gen-sqlite migrate-gen-postgres

# Treat unknown targets as arguments to the run command
%:
	@:

# Variables
BINARY_NAME=threatreg
MAIN_PATH=.
BUILD_DIR=bin
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)
ATLAS_CMD=$(shell go env GOPATH)/bin/atlas

# Default target
help:
	@echo "🚀 Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make clean          - Clean build files"
	@echo "  make test           - Run tests"
	@echo "  make deps           - Download dependencies"
	@echo "  make install        - Install the binary globally"
	@echo "  make migrate-gen-sqlite    - Generate SQLite migrations from GORM models"
	@echo "  make migrate-gen-postgres  - Generate Postgres migrations from GORM models"
	@echo "  make migrate-apply-sqlite  - Apply pending SQLite migrations"
	@echo "  make migrate-apply-postgres - Apply pending Postgres migrations"
	@echo "  make migrate-status-sqlite - Check SQLite migration status"
	@echo "  make migrate-status-postgres - Check Postgres migration status"
	@echo "  make migrate-validate-sqlite - Validate SQLite migration files"
	@echo "  make migrate-validate-postgres - Validate Postgres migration files"
	@echo "  make setup          - Initial project setup"

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "✅ Build complete!"

# Run the application
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@./$(BINARY_PATH) $(filter-out $@,$(MAKECMDGOALS))

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated!"

# Install binary globally
install: build
	@echo "📦 Installing $(BINARY_NAME) globally..."
	@go install $(MAIN_PATH)
	@echo "✅ $(BINARY_NAME) installed! You can now run '$(BINARY_NAME)' from anywhere."

# Development setup
setup:
	@echo "🚀 Setting up development environment..."
	@chmod +x setup.sh
	@./setup.sh

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

# Lint code (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "🔍 Linting code..."; \
		golangci-lint run; \
	else \
		echo "📥 Install golangci-lint first: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Database migration commands using Atlas

migrate-apply-sqlite:
	@echo "⬆️  Applying SQLite database migrations..."
	@$(ATLAS_CMD) migrate apply --env sqlite --dir file://migrations/sqlite

migrate-apply-postgres:
	@echo "⬆️  Applying Postgres database migrations..."
	@$(ATLAS_CMD) migrate apply --env postgres --dir file://migrations/postgres

migrate-status-sqlite:
	@echo "📊 Checking SQLite migration status..."
	@$(ATLAS_CMD) migrate status --env sqlite --dir file://migrations/sqlite

migrate-status-postgres:
	@echo "📊 Checking Postgres migration status..."
	@$(ATLAS_CMD) migrate status --env postgres --dir file://migrations/postgres

migrate-validate-sqlite:
	@echo "✅ Validating SQLite migration files..."
	@$(ATLAS_CMD) migrate validate --env sqlite --dir file://migrations/sqlite

migrate-validate-postgres:
	@echo "✅ Validating Postgres migration files..."
	@$(ATLAS_CMD) migrate validate --env postgres --dir file://migrations/postgres

migrate-gen-sqlite:
	@echo "📝 Generating SQLite migration from GORM models..."
	@mkdir -p migrations/sqlite
	@read -r DESC?'Enter migration description: '; \
	$(ATLAS_CMD) migrate diff --env sqlite --dir file://migrations/sqlite --description "$$DESC"

migrate-gen-postgres:
	@echo "📝 Generating Postgres migration from GORM models..."
	@mkdir -p migrations/postgres
	@read -r DESC?'Enter migration description: '; \
	$(ATLAS_CMD) migrate diff --env postgres --dir file://migrations/postgres --description "$$DESC"

# Show application status
status: build
	@./$(BINARY_PATH) status
