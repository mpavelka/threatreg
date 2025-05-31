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

func TestDatabaseOperationsWithSQLite(t *testing.T) {
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

	// Create tables for testing
	err = CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	t.Run("Create and retrieve user", func(t *testing.T) {
		// Create a user
		user, err := CreateUser("testuser", "test@example.com")
		if err != nil {
			t.Fatalf("CreateUser() error = %v", err)
		}

		if user.Username != "testuser" {
			t.Errorf("CreateUser() username = %v, want testuser", user.Username)
		}
		if user.Email != "test@example.com" {
			t.Errorf("CreateUser() email = %v, want test@example.com", user.Email)
		}

		// Retrieve the user
		retrievedUser, err := GetUserByID(user.ID)
		if err != nil {
			t.Fatalf("GetUserByID() error = %v", err)
		}

		if retrievedUser.Username != user.Username {
			t.Errorf("GetUserByID() username = %v, want %v", retrievedUser.Username, user.Username)
		}
	})

	t.Run("Get all users", func(t *testing.T) {
		users, err := GetAllUsers()
		if err != nil {
			t.Fatalf("GetAllUsers() error = %v", err)
		}

		if len(users) == 0 {
			t.Error("GetAllUsers() returned empty slice, expected at least one user")
		}
	})

	t.Run("Get user count", func(t *testing.T) {
		count, err := GetUserCount()
		if err != nil {
			t.Fatalf("GetUserCount() error = %v", err)
		}

		if count == 0 {
			t.Error("GetUserCount() returned 0, expected at least 1")
		}
	})

	t.Run("Get post count", func(t *testing.T) {
		count, err := GetPostCount()
		if err != nil {
			t.Fatalf("GetPostCount() error = %v", err)
		}

		// Post count should be 0 initially
		if count != 0 {
			t.Errorf("GetPostCount() = %v, want 0", count)
		}
	})

	t.Run("Delete user", func(t *testing.T) {
		// Create a user to delete
		user, err := CreateUser("deleteuser", "delete@example.com")
		if err != nil {
			t.Fatalf("CreateUser() error = %v", err)
		}

		// Delete the user
		err = DeleteUser(user.ID)
		if err != nil {
			t.Fatalf("DeleteUser() error = %v", err)
		}

		// Try to retrieve deleted user
		_, err = GetUserByID(user.ID)
		if err == nil {
			t.Error("GetUserByID() expected error for deleted user, got none")
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		_, err := GetUserByID(99999)
		if err == nil {
			t.Error("GetUserByID() expected error for non-existent user, got none")
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		err := DeleteUser(99999)
		if err == nil {
			t.Error("DeleteUser() expected error for non-existent user, got none")
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