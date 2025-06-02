package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestBuildDatabaseURL(t *testing.T) {
	// Save original config
	originalConfig := AppConfig

	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "DatabaseURL takes precedence",
			config: Config{
				DatabaseURL:      "sqlite3://custom.db",
				DatabaseProtocol: "postgres",
				DatabaseHost:     "localhost",
				DatabasePort:     "5432",
				DatabaseName:     "testdb",
			},
			expected: "sqlite3://custom.db",
		},
		{
			name: "SQLite with default values",
			config: Config{
				DatabaseURL:      "",
				DatabaseProtocol: "sqlite3",
			},
			expected: "sqlite3://app.db",
		},
		{
			name: "SQLite with custom database name",
			config: Config{
				DatabaseURL:      "",
				DatabaseProtocol: "sqlite3",
				DatabaseName:     "custom.db",
			},
			expected: "sqlite3://custom.db",
		},
		{
			name: "PostgreSQL without authentication",
			config: Config{
				DatabaseURL:      "",
				DatabaseProtocol: "postgres",
				DatabaseHost:     "localhost",
				DatabasePort:     "5432",
				DatabaseName:     "testdb",
			},
			expected: "postgresql://localhost:5432/testdb",
		},
		{
			name: "PostgreSQL with authentication",
			config: Config{
				DatabaseURL:      "",
				DatabaseProtocol: "postgres",
				DatabaseUsername: "testuser",
				DatabasePassword: "testpass",
				DatabaseHost:     "localhost",
				DatabasePort:     "5432",
				DatabaseName:     "testdb",
			},
			expected: "postgresql://testuser:testpass@localhost:5432/testdb",
		},
		{
			name: "PostgreSQL with username only",
			config: Config{
				DatabaseURL:      "",
				DatabaseProtocol: "postgres",
				DatabaseUsername: "testuser",
				DatabaseHost:     "localhost",
				DatabasePort:     "5432",
				DatabaseName:     "testdb",
			},
			expected: "postgresql://testuser@localhost:5432/testdb",
		},
		{
			name: "PostgreSQL with defaults",
			config: Config{
				DatabaseURL:      "",
				DatabaseProtocol: "postgres",
			},
			expected: "postgresql://localhost:5432/postgres",
		},
		{
			name: "Unknown protocol",
			config: Config{
				DatabaseURL:      "",
				DatabaseProtocol: "unknown",
			},
			expected: "unknown://app.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test config
			AppConfig = tt.config

			result := BuildDatabaseURL()
			if result != tt.expected {
				t.Errorf("BuildDatabaseURL() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Restore original config
	AppConfig = originalConfig
}

func TestConfigGetters(t *testing.T) {
	// Save original config
	originalConfig := AppConfig

	testConfig := Config{
		DatabaseURL:      "test://db",
		DatabaseProtocol: "postgres",
		DatabaseUsername: "testuser",
		DatabasePassword: "testpass",
		DatabaseHost:     "testhost",
		DatabasePort:     "5432",
		DatabaseName:     "testdb",
		Environment:      "test",
		SecretKey:        "testsecret",
	}

	AppConfig = testConfig

	tests := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{"GetDatabaseURL", GetDatabaseURL, "test://db"},
		{"GetDatabaseProtocol", GetDatabaseProtocol, "postgres"},
		{"GetDatabaseUsername", GetDatabaseUsername, "testuser"},
		{"GetDatabasePassword", GetDatabasePassword, "testpass"},
		{"GetDatabaseHost", GetDatabaseHost, "testhost"},
		{"GetDatabasePort", GetDatabasePort, "5432"},
		{"GetDatabaseName", GetDatabaseName, "testdb"},
		{"GetEnvironment", GetEnvironment, "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			if result != tt.expected {
				t.Errorf("%s() = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}

	// Restore original config
	AppConfig = originalConfig
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalEnv := map[string]string{
		"APP_DATABASE_URL":      os.Getenv("APP_DATABASE_URL"),
		"APP_DATABASE_PROTOCOL": os.Getenv("APP_DATABASE_PROTOCOL"),
		"APP_DATABASE_USERNAME": os.Getenv("APP_DATABASE_USERNAME"),
		"APP_DATABASE_PASSWORD": os.Getenv("APP_DATABASE_PASSWORD"),
		"APP_DATABASE_HOST":     os.Getenv("APP_DATABASE_HOST"),
		"APP_DATABASE_PORT":     os.Getenv("APP_DATABASE_PORT"),
		"APP_DATABASE_NAME":     os.Getenv("APP_DATABASE_NAME"),
		"APP_ENVIRONMENT":       os.Getenv("APP_ENVIRONMENT"),
		"APP_SECRET_KEY":        os.Getenv("APP_SECRET_KEY"),
	}

	// Clean environment
	for key := range originalEnv {
		os.Unsetenv(key)
	}

	// Set test environment variables
	testEnvs := map[string]string{
		"APP_DATABASE_URL":      "postgres://test.db",
		"APP_DATABASE_PROTOCOL": "postgres",
		"APP_DATABASE_USERNAME": "envuser",
		"APP_DATABASE_PASSWORD": "envpass",
		"APP_DATABASE_HOST":     "envhost",
		"APP_DATABASE_PORT":     "9999",
		"APP_DATABASE_NAME":     "envdb",
		"APP_ENVIRONMENT":       "testing",
		"APP_SECRET_KEY":        "envsecret",
	}

	for key, value := range testEnvs {
		os.Setenv(key, value)
	}

	// Reset viper
	viper.Reset()

	// Load config
	err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Test that environment variables were loaded
	tests := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{"DatabaseURL", GetDatabaseURL, "postgres://test.db"},
		{"DatabaseProtocol", GetDatabaseProtocol, "postgres"},
		{"DatabaseUsername", GetDatabaseUsername, "envuser"},
		{"DatabasePassword", GetDatabasePassword, "envpass"},
		{"DatabaseHost", GetDatabaseHost, "envhost"},
		{"DatabasePort", GetDatabasePort, "9999"},
		{"DatabaseName", GetDatabaseName, "envdb"},
		{"Environment", GetEnvironment, "testing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			if result != tt.expected {
				t.Errorf("%s() = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}

	// Restore original environment
	for key := range testEnvs {
		os.Unsetenv(key)
	}
	for key, value := range originalEnv {
		if value != "" {
			os.Setenv(key, value)
		}
	}

	// Reset viper and reload original config
	viper.Reset()
	Load()
}