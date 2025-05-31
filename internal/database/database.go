package database

import (
	"fmt"
	"strings"
	"threatreg/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type DB struct {
	*sqlx.DB
}

var db *DB

// Connect establishes a database connection
func Connect() error {
	dbURL := config.BuildDatabaseURL()

	// Parse the database URL to determine the driver
	var driverName string
	var dataSourceName string

	if strings.HasPrefix(dbURL, "sqlite3://") {
		driverName = "sqlite3"
		dataSourceName = strings.TrimPrefix(dbURL, "sqlite3://")
	} else if strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://") {
		driverName = "postgres"
		dataSourceName = dbURL
	} else {
		return fmt.Errorf("unsupported database URL format: %s", dbURL)
	}

	// Connect to database
	sqlxDB, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlxDB.SetMaxOpenConns(25)
	sqlxDB.SetMaxIdleConns(25)
	sqlxDB.SetConnMaxLifetime(5 * time.Minute)

	db = &DB{sqlxDB}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Database connected successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
