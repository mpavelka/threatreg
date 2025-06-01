# Go CLI Application Makefile

.PHONY: build run clean test deps help migrate-up migrate-down migrate-create

# Treat unknown targets as arguments to the run command
%:
	@:

# Variables
BINARY_NAME=threatreg
MAIN_PATH=.
BUILD_DIR=bin
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)

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

# Show application status
status: build
	@./$(BINARY_PATH) status
