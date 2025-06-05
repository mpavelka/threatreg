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


## Generating Migrations

This project uses [Atlas](https://atlasgo.io/) for database migrations with GORM model integration.

### Atlas Installation

Atlas CLI needs to be installed manually:

```bash
# Install Atlas CLI
go install ariga.io/atlas/cmd/atlas@latest

# Install Atlas GORM provider
go get -u ariga.io/atlas-provider-gorm

# Make sure atlas is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
atlas version
```

To generate migrations for each database, use the following commands:

### SQLite
- Generate: `make migrate-gen-sqlite`
- Apply: `make migrate-apply-sqlite`
- Status: `make migrate-status-sqlite`
- Validate: `make migrate-validate-sqlite`

### Postgres
- Generate: `make migrate-gen-postgres`
- Apply: `make migrate-apply-postgres`
- Status: `make migrate-status-postgres`
- Validate: `make migrate-validate-postgres`

### Migration Workflow

1. **Update your GORM models:**
```go
// internal/models/product.go
type Product struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key"`
    Name        string    `gorm:"not null"`
    Description string
    gorm.Model
}
```

2. **Generate migration from models:**
```bash
$(go env GOPATH)/bin/atlas migrate diff $MIGRATION_SUFFIX --env sqlite --dir file://migrations/sqlite
```
This creates a new migration file in the `migrations/` directory based on the difference between your current database schema and your GORM models.

3. **Create a down migration file:**
For each generated migration (e.g., `20250602124850.sql`), create a corresponding down migration file named `20250602124850_down.sql` in the same directory. This file should contain SQL to revert the changes. For example:
```sql
-- Drop index and table for products (rollback)
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP TABLE IF EXISTS products;
```

4. **Review the generated migration:**
Atlas will show you the SQL changes and create a migration file like `migrations/20240101120000.sql`.

5. **Apply the migration:**
```bash
atlas migrate apply --env sqlite
```

6. **Verify it worked:**
```bash
atlas migrate status --env sqlite
```

## Troubleshooting

To apply only one migration (go one migration up):

```sh
atlas migrate apply --env sqlite --dir file://migrations/sqlite --limit 1
```

To roll back (revert) the last migration (go one migration down):

```sh
atlas migrate down --env sqlite --dir file://migrations/sqlite --steps 1
```

Replace `sqlite` and the directory as needed for Postgres.

### Environment Variables

Set your database URL:

```bash
# SQLite (default)
export DATABASE_URL="sqlite3://app.db"

# PostgreSQL
export DATABASE_URL="postgresql://username:password@localhost:5432/database_name"
```

## 🏭 Production Deployment

### Build for Production
```bash
# Build optimized binary
go build -ldflags="-w -s" -o myapp .

# Or with Make
make build
```
