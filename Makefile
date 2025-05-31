# Go CLI Application Makefile

.PHONY: build run clean test deps help migrate-up migrate-down migrate-create

# Variables
BINARY_NAME=myapp
MAIN_PATH=.
BUILD_DIR=bin

# Default target
help:
	@echo "🚀 Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make clean          - Clean build files"
	@echo "  make test           - Run tests"
	@echo "  make deps           - Download dependencies"
	@echo "  make install        - Install the binary globally"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-create NAME=<name> - Create new migration"
	@echo "  make setup          - Initial project setup"

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete!"

# Run the application
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

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
	@echo "✅ $(BINARY_NAME) installed! You can now run 'myapp' from anywhere."

# Database migration commands
migrate-up: build
	@echo "⬆️  Running database migrations..."
	@./$(BINARY_NAME) db up

migrate-down: build
	@echo "⬇️  Rolling back last migration..."
	@./$(BINARY_NAME) db down

migrate-create: build
ifndef NAME
	@echo "❌ Please provide a migration name: make migrate-create NAME=add_users_table"
else
	@echo "📝 Creating migration: $(NAME)"
	@./$(BINARY_NAME) db create "$(NAME)"
endif

# Development setup
setup:
	@echo "🚀 Setting up development environment..."
	@chmod +x setup.sh
	@./setup.sh

# Development server with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@if command -v air > /dev/null; then \
		echo "🔄 Starting development server with auto-reload..."; \
		air; \
	else \
		echo "📥 Installing air for auto-reload..."; \
		go install github.com/cosmtrek/air@latest; \
		echo "🔄 Starting development server..."; \
		air; \
	fi

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

# Show application status
status: build
	@./$(BINARY_NAME) status

# Create a user quickly
user: build
	@./$(BINARY_NAME) user create --username admin --email admin@example.com

# Start the server
serve: build
	@./$(BINARY_NAME) serve
