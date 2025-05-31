# Go CLI Application with Database Migrations

A complete Go application template with database migrations, CLI interface, and user management. Built with modern Go practices and robust tooling.

## 🚀 Quick Start

```bash
# Setup everything
chmod +x setup.sh && ./setup.sh

# Or use Make
make setup
```

## ✨ Features

- **CLI Interface** with Cobra - Professional command-line interface
- **Database Migrations** with golang-migrate - Version-controlled schema changes
- **Multi-database Support** - SQLite, PostgreSQL (MySQL ready)
- **Configuration Management** - Environment variables with Viper
- **User Management** - Complete CRUD operations
- **Production Ready** - Proper error handling, logging, and project structure

## 📋 Requirements

- Go 1.24+
- Make (optional, for convenience commands)

## 🏗️ Project Structure

```
myapp/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and setup
│   ├── db.go              # Database management commands
│   ├── user.go            # User management commands
│   └── status.go          # Status and serve commands
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database models and operations
│   └── migrations/        # Migration management
├── migrations/            # SQL migration files
├── .env.example          # Environment configuration template
├── go.mod                # Go module dependencies
├── Makefile              # Development commands
└── README.md             # This file
```

## 🗄️ Database Migrations

This project uses the official [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool for database migrations.

### Installation

The `migrate` CLI is automatically installed during setup, or install it manually:

```bash
# Install migrate CLI
go install -tags 'sqlite3,postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verify installation
migrate -version
```

### Migration Commands

```bash
# Create a new migration
migrate create -ext sql -dir migrations create_users_table

# Run all pending migrations
migrate -path migrations -database $DATABASE_URL up

# Rollback one migration
migrate -path migrations -database $DATABASE_URL down 1

# Check migration status
migrate -path migrations -database $DATABASE_URL version

# Force migration version (recovery)
migrate -path migrations -database $DATABASE_URL force 1

# Drop all tables (DANGER!)
migrate -path migrations -database $DATABASE_URL drop
```

### Make Shortcuts

For convenience, you can use Make commands:

```bash
# Create migration
make migrate-create NAME=add_posts_table

# Run migrations
make migrate-up

# Rollback last migration  
make migrate-down

# Check status
make migrate-status

# Force version (for recovery)
make migrate-force VERSION=1
```

### Migration Workflow

1. **Create a migration:**
```bash
migrate create -ext sql -dir migrations add_users_table
```

This creates two files:
- `migrations/000001_add_users_table.up.sql` - Forward migration
- `migrations/000001_add_users_table.down.sql` - Rollback migration

2. **Edit the migration files:**

**`000001_add_users_table.up.sql`:**
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

**`000001_add_users_table.down.sql`:**
```sql
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
```

3. **Run the migration:**
```bash
migrate -path migrations -database $DATABASE_URL up
```

4. **Verify it worked:**
```bash
./threatreg status
```

### Environment Variables

Set your database URL:

```bash
# SQLite (default)
export DATABASE_URL="sqlite3://app.db"

# PostgreSQL
export DATABASE_URL="postgresql://username:password@localhost:5432/database_name"
```

### Production Deployment

For production, generate SQL scripts for review:

```bash
# Generate SQL for all pending migrations
migrate -path migrations -database $DATABASE_URL up -dry-run > migration_script.sql

# Review the SQL
cat migration_script.sql

# Apply manually or through your deployment process
```

### Troubleshooting

**Migration fails:**
```bash
# Check current version
migrate -path migrations -database $DATABASE_URL version

# Force to specific version if needed
migrate -path migrations -database $DATABASE_URL force 1
```

**Dirty database state:**
```bash
# Check logs and fix manually, then force to clean version
migrate -path migrations -database $DATABASE_URL force 1
```

## 📚 Usage

### Database Commands

```bash
# Initialize migration system
./myapp db init

# Create a new migration
./myapp db create "add_users_table"

# Run all pending migrations
./myapp db up

# Rollback last migration
./myapp db down

# Check migration status
./myapp db status

# Reset database (DANGER!)
./myapp db reset

# Create tables directly (development only)
./myapp db setup
```

### User Management

```bash
# Create a new user
./myapp user create --username john --email john@example.com
./myapp user create -u jane -e jane@example.com

# List all users
./myapp user list

# Show specific user
./myapp user show 1

# Delete a user
./myapp user delete 1
```

### Application Commands

```bash
# Show application status
./myapp status

# Start the server (implement your HTTP server)
./myapp serve --port 8080

# Show help
./myapp --help
./myapp db --help
./myapp user --help
```

## 🗄️ Database Configuration

### SQLite (Default)
```bash
DATABASE_URL=sqlite3://app.db
```

### PostgreSQL
```bash
DATABASE_URL=postgresql://username:password@localhost:5432/database_name
```

### Environment Variables
Create a `.env` file:
```bash
# Copy from example
cp .env.example .env

# Edit your settings
DATABASE_URL=sqlite3://app.db
APP_ENV=development
SECRET_KEY=your-secret-key-here
PORT=8080
```

## 🔄 Development Workflow

### 1. Create a Migration
```bash
# Create migration file
./myapp db create "add_posts_table"

# Edit the generated files in migrations/
# - migrations/TIMESTAMP_add_posts_table.up.sql
# - migrations/TIMESTAMP_add_posts_table.down.sql
```

### 2. Run Migration
```bash
# Apply migration
./myapp db up

# Check status
./myapp db status
```

### 3. Update Models
Edit `internal/database/database.go` to add your new models and operations.

### 4. Test Your Changes
```bash
# Create test data
./myapp user create -u testuser -e test@example.com

# Check application status
./myapp status
```

## 🛠️ Make Commands

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Build and run
make test           # Run tests
make clean          # Clean build files
make deps           # Update dependencies
make install        # Install globally
make fmt            # Format code
make lint           # Lint code (requires golangci-lint)

# Database shortcuts
make migrate-up     # Run migrations
make migrate-down   # Rollback migration
make migrate-create NAME="migration_name"  # Create migration

# Development
make dev            # Auto-reload server (requires air)
make status         # Show app status
make user           # Create admin user
make serve          # Start server
```

## 📝 Creating Migrations

### Example Migration Files

**up migration** (`migrations/20240530120000_create_users.up.sql`):
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

**down migration** (`migrations/20240530120000_create_users.down.sql`):
```sql
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
```

## 🏭 Production Deployment

### Build for Production
```bash
# Build optimized binary
go build -ldflags="-w -s" -o myapp .

# Or with Make
make build
```

### Environment Setup
```bash
# Production environment variables
export DATABASE_URL="postgresql://user:pass@localhost:5432/prod_db"
export APP_ENV="production"
export SECRET_KEY="your-production-secret"

# Run migrations
./myapp db up

# Start server
./myapp serve --port 8080
```

### Docker (Optional)
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-w -s" -o myapp .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
CMD ["./myapp", "serve"]
```

## 🔍 Troubleshooting

### Common Issues

**Migration fails:**
```bash
# Check current status
./myapp db status

# Force to specific version if needed
./myapp db force 1
```

**Database connection issues:**
```bash
# Check your DATABASE_URL in .env
# Verify database server is running
./myapp status
```

**Build errors:**
```bash
# Update dependencies
go mod tidy

# Clean and rebuild
make clean && make build
```

## 📦 Dependencies

- **[cobra](https://github.com/spf13/cobra)** - CLI framework
- **[viper](https://github.com/spf13/viper)** - Configuration management
- **[sqlx](https://github.com/jmoiron/sqlx)** - SQL toolkit
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migrations
- **[godotenv](https://github.com/joho/godotenv)** - Environment variables
- **Database drivers:** SQLite, PostgreSQL

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Happy coding! 🚀** 

This template gives you a solid foundation for building CLI applications in Go with proper database management, migrations, and modern Go practices.
