# Go CLI Application with Database Migrations

A complete Go application template with database migrations, CLI interface.

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
- **Production Ready** - Proper error handling, logging, and project structure

## 📋 Requirements

- Go 1.24+
- Make (optional, for convenience commands)


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

## 🗄️ Database Configuration

### SQLite (Default)
```bash
DATABASE_URL=sqlite3://app.db
```

### PostgreSQL
```bash
DATABASE_URL=postgresql://username:password@localhost:5432/database_name
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
