package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL      string `mapstructure:"database_url"`
	DatabaseProtocol string `mapstructure:"database_protocol"`
	DatabaseUsername string `mapstructure:"database_username"`
	DatabasePassword string `mapstructure:"database_password"`
	DatabaseHost     string `mapstructure:"database_host"`
	DatabasePort     string `mapstructure:"database_port"`
	DatabaseName     string `mapstructure:"database_name"`
	Environment      string `mapstructure:"environment"`
	SecretKey        string `mapstructure:"secret_key"`
	APIHost          string `mapstructure:"api_host"`
	APIPort          string `mapstructure:"api_port"`
}

var AppConfig Config

func Load() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't fail if it doesn't exist
		fmt.Println("No .env file found, using environment variables")
	}

	// Set defaults
	viper.SetDefault("database_url", "")
	viper.SetDefault("database_protocol", "sqlite3")
	viper.SetDefault("database_username", "")
	viper.SetDefault("database_password", "")
	viper.SetDefault("database_host", "")
	viper.SetDefault("database_port", "")
	viper.SetDefault("database_name", "app.db")
	viper.SetDefault("environment", "development")
	viper.SetDefault("secret_key", "your-secret-key-change-this")
	viper.SetDefault("api_host", "localhost")
	viper.SetDefault("api_port", "8080")

	// Configure environment variable handling
	viper.SetEnvPrefix("APP") // Look for APP_* variables
	viper.AutomaticEnv()      // Automatically bind environment variables

	// Replace underscores with dots for nested config (optional)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal into config struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}

	return nil
}

func GetDatabaseURL() string {
	return AppConfig.DatabaseURL
}

func GetDatabaseProtocol() string {
	return AppConfig.DatabaseProtocol
}

func GetDatabaseUsername() string {
	return AppConfig.DatabaseUsername
}

func GetDatabasePassword() string {
	return AppConfig.DatabasePassword
}

func GetDatabaseHost() string {
	return AppConfig.DatabaseHost
}

func GetDatabasePort() string {
	return AppConfig.DatabasePort
}

func GetDatabaseName() string {
	return AppConfig.DatabaseName
}

func GetEnvironment() string {
	return AppConfig.Environment
}

// GetAPIHost returns the API server host configuration.
func GetAPIHost() string {
	return AppConfig.APIHost
}

// GetAPIPort returns the API server port configuration.
func GetAPIPort() string {
	return AppConfig.APIPort
}

// BuildDatabaseURL constructs a database URL from individual components
// If DatabaseURL is set, it takes precedence
func BuildDatabaseURL() string {
	if AppConfig.DatabaseURL != "" {
		return AppConfig.DatabaseURL
	}

	protocol := AppConfig.DatabaseProtocol
	if protocol == "" {
		protocol = "sqlite3"
	}

	switch protocol {
	case "sqlite3":
		dbName := AppConfig.DatabaseName
		if dbName == "" {
			dbName = "app.db"
		}
		return fmt.Sprintf("sqlite3://%s", dbName)
	case "postgres", "postgresql":
		host := AppConfig.DatabaseHost
		if host == "" {
			host = "localhost"
		}
		port := AppConfig.DatabasePort
		if port == "" {
			port = "5432"
		}
		dbName := AppConfig.DatabaseName
		if dbName == "" {
			dbName = "postgres"
		}

		var connStr strings.Builder
		connStr.WriteString("postgresql://")

		if AppConfig.DatabaseUsername != "" {
			connStr.WriteString(AppConfig.DatabaseUsername)
			if AppConfig.DatabasePassword != "" {
				connStr.WriteString(":")
				connStr.WriteString(AppConfig.DatabasePassword)
			}
			connStr.WriteString("@")
		}

		connStr.WriteString(host)
		connStr.WriteString(":")
		connStr.WriteString(port)
		connStr.WriteString("/")
		connStr.WriteString(dbName)

		return connStr.String()
	default:
		return fmt.Sprintf("%s://app.db", protocol)
	}
}
