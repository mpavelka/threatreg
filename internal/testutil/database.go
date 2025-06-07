package testutil

import (
	"os"
	"path/filepath"
	"testing"
	"threatreg/internal/config"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/stretchr/testify/require"
)

// SetupTestDatabase creates a fresh test database with all migrations applied.
// It returns a cleanup function that should be called when the test is done.
// This function:
// - Creates a temporary SQLite database file
// - Saves and restores the original configuration
// - Connects to the test database
// - Runs all model migrations
// - Returns a cleanup function that closes the database and removes temp files
func SetupTestDatabase(t *testing.T) func() {
	return SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Application{},
		&models.Threat{},
		&models.Control{},
		&models.ThreatAssignment{},
		&models.ControlAssignment{},
		&models.ThreatControl{},
	)
}

// SetupTestDatabaseWithCustomModels creates a test database with only specific models migrated.
// This is useful when you want to test specific model interactions without the full schema.
func SetupTestDatabaseWithCustomModels(t *testing.T, models ...interface{}) func() {
	// Create a temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Save original config
	originalConfig := config.AppConfig

	// Set test config to use temporary SQLite database
	config.AppConfig = config.Config{
		DatabaseProtocol: "sqlite3",
		DatabaseName:     dbPath,
	}

	// Connect to test database
	err := database.Connect()
	require.NoError(t, err, "Failed to connect to test database")

	// Run migrations for specified models only
	db := database.GetDB()
	require.NotNil(t, db, "Database connection should not be nil")

	if len(models) > 0 {
		err = db.AutoMigrate(models...)
		require.NoError(t, err, "Failed to run migrations")
	}

	// Return cleanup function
	return func() {
		database.Close()
		config.AppConfig = originalConfig
		os.RemoveAll(tempDir)
	}
}
