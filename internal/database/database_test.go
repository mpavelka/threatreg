package database

import (
	"strings"
	"testing"
	"threatreg/internal/config"

	"github.com/spf13/viper"
)

func TestConnect(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig

	tests := []struct {
		name      string
		config    config.Config
		expectErr bool
		errMsg    string
	}{
		{
			name: "SQLite connection",
			config: config.Config{
				DatabaseURL:      "sqlite3://:memory:",
				DatabaseProtocol: "sqlite3",
				DatabaseName:     ":memory:",
			},
			expectErr: false,
		},
		{
			name: "SQLite with granular config",
			config: config.Config{
				DatabaseURL:      "",
				DatabaseProtocol: "sqlite3",
				DatabaseName:     ":memory:",
			},
			expectErr: false,
		},
		{
			name: "Invalid database URL format",
			config: config.Config{
				DatabaseURL:      "invalid://format",
				DatabaseProtocol: "invalid",
			},
			expectErr: true,
			errMsg:    "unsupported database URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test config
			config.AppConfig = tt.config
			
			// Reset any existing connection
			if db != nil {
				db.Close()
				db = nil
			}

			err := Connect()

			if tt.expectErr {
				if err == nil {
					t.Errorf("Connect() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Connect() error = %v, expected to contain %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Connect() error = %v, expected no error", err)
				}
				
				// Test that we can get the database instance
				dbInstance := GetDB()
				if dbInstance == nil {
					t.Error("GetDB() returned nil after successful connection")
				}
				
				// Clean up
				Close()
			}
		})
	}

	// Restore original config
	config.AppConfig = originalConfig
}

func TestDatabaseConnectionWithSQLite(t *testing.T) {
	// Set up in-memory SQLite for testing
	originalConfig := config.AppConfig
	config.AppConfig = config.Config{
		DatabaseProtocol: "sqlite3",
		DatabaseName:     ":memory:",
	}

	// Connect to database
	err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer func() {
		Close()
		config.AppConfig = originalConfig
	}()

	t.Run("Database connection is established", func(t *testing.T) {
		dbInstance := GetDB()
		if dbInstance == nil {
			t.Error("GetDB() returned nil after successful connection")
		}

		// Test that we can ping the database
		if err := dbInstance.Ping(); err != nil {
			t.Errorf("Database ping failed: %v", err)
		}
	})
}

func TestBuildDatabaseURLIntegration(t *testing.T) {
	// Test that Connect() uses BuildDatabaseURL() correctly
	originalConfig := config.AppConfig
	
	// Reset viper to ensure clean state
	viper.Reset()
	
	// Test with granular config (no DatabaseURL)
	config.AppConfig = config.Config{
		DatabaseURL:      "",
		DatabaseProtocol: "sqlite3",
		DatabaseName:     ":memory:",
	}

	// Reset any existing connection
	if db != nil {
		db.Close()
		db = nil
	}

	err := Connect()
	if err != nil {
		t.Fatalf("Connect() with granular config error = %v", err)
	}
	
	dbInstance := GetDB()
	if dbInstance == nil {
		t.Error("GetDB() returned nil after successful connection with granular config")
	}
	
	Close()

	// Test with DatabaseURL taking precedence
	config.AppConfig = config.Config{
		DatabaseURL:      "sqlite3://:memory:",
		DatabaseProtocol: "postgres", // This should be ignored
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
	}

	// Reset connection
	db = nil

	err = Connect()
	if err != nil {
		t.Fatalf("Connect() with DatabaseURL precedence error = %v", err)
	}
	
	dbInstance = GetDB()
	if dbInstance == nil {
		t.Error("GetDB() returned nil after successful connection with DatabaseURL precedence")
	}
	
	Close()

	// Restore original config
	config.AppConfig = originalConfig
}