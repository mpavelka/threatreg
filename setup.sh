#!/bin/bash

# Variables
BUILD_DIR="bin"
BINARY_NAME="threatreg"
BINARY_PATH="$BUILD_DIR/$BINARY_NAME"

echo "🚀 Setting up Threatreg application..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Initialize Go module if go.mod doesn't exist
if [ ! -f "go.mod" ]; then
    echo "📦 Initializing Go module..."
    go mod init threatreg
fi

# Download dependencies and update go.sum
echo "📥 Downloading dependencies..."
go mod download
go mod tidy

# Install migrate CLI if not present
if ! command -v migrate &> /dev/null; then
    echo "📦 Installing migrate CLI..."
    go install -tags 'sqlite3,postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    echo "✅ Migrate CLI installed"
else
    echo "✅ Migrate CLI already installed"
fi

# Verify dependencies
echo "🔍 Verifying dependencies..."
go mod verify

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ Created .env file - please update it with your configuration"
fi

# Build the application
echo "🔨 Building application..."
mkdir -p $BUILD_DIR
go build -o $BINARY_PATH .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    echo "💡 Try running: go clean -modcache && go mod tidy"
    exit 1
fi

# Create migrations directory
echo "🗄️  Creating migrations directory..."
mkdir -p migrations

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Available commands:"
echo "  ./$BINARY_PATH --help                    # Show all commands"
echo "  ./$BINARY_PATH db setup                  # Quick database setup (dev only)"
echo "  ./$BINARY_PATH user create -u name -e email@example.com"
echo "  ./$BINARY_PATH user list                 # List users"
echo "  ./$BINARY_PATH status                    # Show app status"
echo "  ./$BINARY_PATH serve                     # Start server"
echo ""
echo "Database Migrations (use migrate CLI directly):"
echo "  migrate create -ext sql -dir migrations create_users"
echo "  migrate -path migrations -database \$DATABASE_URL up"
echo "  migrate -path migrations -database \$DATABASE_URL down 1"
echo "  migrate -path migrations -database \$DATABASE_URL version"
echo ""
echo "Quick start:"
echo "  1. Create your first migration:"
echo "     migrate create -ext sql -dir migrations create_users_table"
echo "  2. Edit the generated SQL files in ./migrations/"
echo "  3. Run migrations: migrate -path migrations -database \$DATABASE_URL up"
echo "  4. Create a user: ./$BINARY_PATH user create -u admin -e admin@example.com"
echo "  5. Check status: ./$BINARY_PATH status"
