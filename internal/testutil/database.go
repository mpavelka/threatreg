package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"threatreg/internal/config"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// SetupTestDatabase creates a fresh test database with all migrations applied.
// It returns a cleanup function that should be called when the test is done.
// This function:
// - Creates a temporary PostgreSQL database (or falls back to SQLite for basic tests)
// - Saves and restores the original configuration
// - Connects to the test database
// - Runs all model migrations
// - Returns a cleanup function that closes the database and removes temp files
func SetupTestDatabase(t *testing.T) func() {
	return SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentRelationship{},
		&models.ComponentAttribute{},
		&models.Domain{},
		&models.Tag{},
		&models.Threat{},
		&models.Control{},
		&models.ThreatAssignment{},
		&models.ControlAssignment{},
		&models.ThreatControl{},
		&models.ThreatAssignmentResolution{},
		&models.ThreatAssignmentResolutionDelegation{},
		&models.Relationship{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
}

// SetupTestDatabaseWithCustomModels creates a test database with only specific models migrated.
// This is useful when you want to test specific model interactions without the full schema.
// Now defaults to PostgreSQL for advanced features like recursive CTEs.
func SetupTestDatabaseWithCustomModels(t *testing.T, models ...interface{}) func() {
	// Save original config
	originalConfig := config.AppConfig

	// Try PostgreSQL first for advanced features, fallback to SQLite for basic tests
	pgURL := os.Getenv("TEST_DATABASE_URL")
	if pgURL == "" {
		pgURL = "postgresql://threatreg_test:threatreg_test_password@localhost:5433/threatreg_test?sslmode=disable"
	}

	// Set test config to use PostgreSQL
	config.AppConfig = config.Config{
		DatabaseURL: pgURL,
	}

	// Connect to test database
	err := database.Connect()
	if err != nil {
		// Fallback to SQLite for basic tests if PostgreSQL is not available
		t.Logf("PostgreSQL not available (%v), falling back to SQLite (advanced tree queries will be skipped)", err)

		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "test.db")

		config.AppConfig = config.Config{
			DatabaseProtocol: "sqlite3",
			DatabaseName:     dbPath,
		}

		err = database.Connect()
		require.NoError(t, err, "Failed to connect to SQLite test database")
	}

	// Clean up any existing test data BEFORE migrations to handle schema changes
	db := database.GetDB()
	require.NotNil(t, db, "Database connection should not be nil")
	
	if IsPostgreSQL() {
		cleanupTestData(t, db)
	}

	// Run migrations for specified models only
	if len(models) > 0 {
		err = db.AutoMigrate(models...)
		require.NoError(t, err, "Failed to run migrations")
	}

	// Return cleanup function
	return func() {
		database.Close()
		config.AppConfig = originalConfig
	}
}

// IsPostgreSQL returns true if the current test database is PostgreSQL
func IsPostgreSQL() bool {
	db := database.GetDB()
	if db == nil {
		return false
	}

	// Get the database name from the GORM dialector
	dbName := db.Dialector.Name()
	return strings.Contains(dbName, "postgres")
}

// RequirePostgreSQL skips the test if not running on PostgreSQL
func RequirePostgreSQL(t *testing.T) {
	if !IsPostgreSQL() {
		t.Skip("This test requires PostgreSQL for advanced recursive CTE features")
	}
}

// cleanupTestData drops all tables to ensure fresh schema for UUID migration
// This is only used for PostgreSQL since SQLite uses temporary databases
func cleanupTestData(t *testing.T, db *gorm.DB) {
	// List of tables that need to be dropped in dependency order
	// (child tables first, then parent tables)
	tablesToDrop := []string{
		"threat_assignment_resolution_delegations",
		"control_assignments", 
		"threat_assignment_resolutions",
		"threat_assignments",
		"component_relationships",
		"threat_controls",
		"pattern_conditions",    // threat pattern condition table
		"threat_patterns",       // threat pattern table
		"component_attributes",  // component attributes table
		"component_tags",        // many2many join table
		"domain_components",     // many2many join table
		"controls",
		"threats", 
		"components",
		"domains",
		"tags",
		"relationships",
		// Add any other tables as needed
	}

	for _, tableName := range tablesToDrop {
		// Drop tables completely to handle schema changes (like int -> UUID migration)
		result := db.Exec("DROP TABLE IF EXISTS " + tableName + " CASCADE")
		// Ignore errors for tables that might not exist in the specific test
		if result.Error != nil {
			t.Logf("Warning: Failed to drop table %s: %v", tableName, result.Error)
		}
	}
}
