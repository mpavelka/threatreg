# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Build and Run
- `make build` - Build the application to `bin/threatreg`
- `make run` - Build and run the application
- `./bin/threatreg` - Run the built binary directly
- `make status` - Show application status

### Database Operations
- `make migrate-up` - Run all pending migrations
- `make migrate-down` - Rollback last migration
- `make migrate-create NAME="migration_name"` - Create new migration files
- `./bin/threatreg db status` - Check migration status
- `./bin/threatreg db setup` - Create tables directly (development only)

### Code Quality
- `make fmt` - Format Go code
- `make lint` - Lint code (requires golangci-lint)
- `make test` - Run tests

### Dependencies
- `make deps` - Download and tidy Go module dependencies

## Architecture Overview

This is a Go CLI application for threat registry management built with:

### Core Components
- **Cobra CLI Framework**: All commands are in `cmd/` with `cmd/root.go` as the entry point
- **Database Layer**: `internal/database/database.go` provides SQLX-based data access with connection pooling
- **Configuration**: `internal/config/config.go` uses Viper for environment-based config with .env support
- **Migrations**: Uses golang-migrate for version-controlled schema changes

### Database Support
- **Multi-database**: Supports SQLite (default) and PostgreSQL
- **Connection String Format**: `sqlite3://app.db` or `postgresql://user:pass@host:port/db`
- **Models**: User and Post structs with repository pattern methods

### Project Structure Pattern
- `cmd/` - CLI command definitions (one file per command group)
- `internal/config/` - Configuration management
- `internal/database/` - Database models and operations
- `pkg/migrations/` - Migration management utilities
- `main.go` - Application entry point

### Environment Configuration
Uses Viper with these key variables:
- `DATABASE_URL` - Database connection string
- `APP_ENV` - Environment (development/production)
- `SECRET_KEY` - Application secret

### CLI Command Structure
- Root command: `threatreg`
- Database commands: `threatreg db [init|create|up|down|status|setup|reset]`
- Status command: `threatreg status`
- Serve command: `threatreg serve`

The application follows Go project layout conventions with clear separation between CLI interface, business logic, and data persistence layers.