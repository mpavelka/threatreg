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

### Code Quality
- `make fmt` - Format Go code
- `make lint` - Lint code (requires golangci-lint)
- `make test` - Run tests

### Dependencies
- `make deps` - Download and tidy Go module dependencies

## Project Overview

**Threatreg** is a CLI-based threat registry management system designed for cybersecurity threat modeling and risk management. It enables organizations to manage products, track threats, define security controls, and map their relationships for comprehensive risk assessment.

### Core Domain Concepts

The application models a threat management domain with these key entities:

#### **Primary Entities**
- **Product**: Software products/systems (e.g., "E-commerce Platform", "Mobile App")
- **Application**: Specific instances/deployments of products (e.g., "Production E-commerce", "Staging E-commerce")
- **Threat**: Security threats and vulnerabilities (e.g., "SQL Injection", "Cross-Site Scripting")
- **Control**: Security controls and countermeasures (e.g., "Input Validation", "WAF Implementation")

#### **Relationship Models**
- **ThreatAssignment**: Links threats to specific products/applications
- **ControlAssignment**: Associates controls with threat assignments for mitigation
- **ThreatControl**: General mapping between threats and applicable controls

### Architecture Overview

This is a Go CLI application built with:

#### **Core Components**
- **Cobra CLI Framework**: All commands in `cmd/` with `cmd/root.go` as entry point
- **GORM Database Layer**: `internal/database/database.go` provides ORM-based data access
- **Configuration**: `internal/config/config.go` uses Viper for environment-based config
- **Service Layer**: `internal/service/` contains business logic and repository patterns
- **Migrations**: Uses Atlas for schema versioning from GORM models

#### **Database Support**
- **Multi-database**: Supports SQLite (default) and PostgreSQL
- **Connection String Format**: `sqlite3://app.db` or `postgresql://user:pass@host:port/db`
- **UUID Primary Keys**: All entities use UUIDs for identification
- **GORM Models**: Rich domain models with relationships and lifecycle hooks

#### **Project Structure**
- `cmd/` - CLI command definitions (one file per command group)
- `internal/config/` - Configuration management with Viper
- `internal/database/` - GORM models and database operations
- `internal/service/` - Business logic and repository patterns
- `internal/models/` - Domain models (Product, Threat, Control, etc.)
- `migrations/` - Atlas migration files organized by database type
- `main.go` - Application entry point

#### **Environment Configuration**
Uses Viper with these key variables:
- `APP_DATABASE_URL` - Database connection string
- `APP_ENVIRONMENT` - Environment (development/production)
- Additional config loaded from .env files

### Current CLI Commands

#### **Available Commands**
- **`threatreg status`** - Show application and database connectivity status
- **`threatreg product create --name "Name" --description "Desc"`** - Create new products
- **`threatreg product get --id <uuid>`** - Retrieve specific product
- **`threatreg product update --id <uuid> --name "Name"`** - Update product info
- **`threatreg product delete --id <uuid>`** - Remove products
- **`threatreg product list`** - List all products

#### **Implementation Status**
- ✅ **Product management**: Fully implemented with complete CRUD operations
- 🚧 **Threat management**: Models defined, CLI commands pending
- 🚧 **Control management**: Models defined, CLI commands pending
- 🚧 **Assignment management**: Models defined, CLI commands pending

### Development Patterns

#### **Service Layer Pattern**
- Repository pattern with GORM operations
- Transaction support for safe updates
- Proper error handling and validation
- UUID-based entity identification

#### **Database Architecture**
- GORM ORM with multiple database driver support
- Atlas-powered migrations from model definitions
- Connection pooling and proper cleanup
- Multi-environment database configurations

The application follows modern Go practices with hexagonal architecture, domain-driven design principles, and comprehensive error handling for production use.