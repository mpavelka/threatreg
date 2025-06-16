package database

import (
	"fmt"
	"strings"
	"threatreg/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Connect establishes a database connection using GORM
func Connect() error {
	dbURL := config.BuildDatabaseURL()

	var gormDB *gorm.DB
	var err error

	if strings.HasPrefix(dbURL, "sqlite3://") {
		dataSourceName := strings.TrimPrefix(dbURL, "sqlite3://")
		gormDB, err = gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else if strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://") {
		gormDB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		return fmt.Errorf("unsupported database URL format: %s", dbURL)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	db = gormDB

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// GetDB returns the GORM database instance
func GetDB() *gorm.DB {
	return db
}

func GetDBOrError() (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return db, nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
