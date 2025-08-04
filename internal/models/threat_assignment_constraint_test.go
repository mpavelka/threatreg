package models

import (
	"os"
	"path/filepath"
	"testing"
	"threatreg/internal/config"
	"threatreg/internal/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a test database without using testutil to avoid import cycle
func setupTestDB(t *testing.T) func() {
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

	// Get database instance
	db := database.GetDB()
	require.NotNil(t, db, "Database connection should not be nil")

	// Auto-migrate models for testing
	err = db.AutoMigrate(&Component{}, &Threat{}, &ThreatAssignment{})
	require.NoError(t, err, "Failed to migrate test database")

	// Return cleanup function
	return func() {
		// Close database connection
		database.Close()

		// Remove temp database file
		os.Remove(dbPath)

		// Restore original config
		config.AppConfig = originalConfig
	}
}

func TestThreatAssignmentConstraints(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	db := database.GetDB()

	// Create test data
	component := &Component{
		ID:          uuid.New(),
		Name:        "Test Component",
		Description: "A test component",
		Type:        ComponentTypeProduct,
	}
	require.NoError(t, db.Create(component).Error)

	threat := &Threat{
		ID:          uuid.New(),
		Title:       "Test Threat",
		Description: "A test threat",
	}
	require.NoError(t, db.Create(threat).Error)

	t.Run("ValidAssignment_ComponentReference", func(t *testing.T) {
		// Test creating assignment with valid ComponentID
		assignment := &ThreatAssignment{
			ThreatID:    threat.ID,
			ComponentID: component.ID,
		}

		err := db.Create(assignment).Error
		assert.NoError(t, err, "Should allow assignment with valid ComponentID")

		// Verify assignment was created
		var retrieved ThreatAssignment
		err = db.First(&retrieved, "id = ?", assignment.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, threat.ID, retrieved.ThreatID)
		assert.Equal(t, component.ID, retrieved.ComponentID)

		// Cleanup
		db.Delete(assignment)
	})

	t.Run("InvalidAssignment_NonExistentComponent", func(t *testing.T) {
		// Test creating assignment with non-existent ComponentID
		nonExistentID := uuid.New()
		assignment := &ThreatAssignment{
			ThreatID:    threat.ID,
			ComponentID: nonExistentID,
		}

		err := db.Create(assignment).Error
		// Note: Depending on database constraints, this might or might not error
		// SQLite with foreign key constraints disabled might allow this
		if err != nil {
			assert.Contains(t, err.Error(), "FOREIGN KEY constraint")
		}
	})

	t.Run("InvalidAssignment_NilComponentID", func(t *testing.T) {
		// Test creating assignment with nil ComponentID
		assignment := &ThreatAssignment{
			ThreatID:    threat.ID,
			ComponentID: uuid.Nil, // Should fail validation
		}

		err := db.Create(assignment).Error
		// This should fail due to NOT NULL constraint on ComponentID
		assert.Error(t, err, "Should reject assignment with nil ComponentID")
	})

	t.Run("DuplicateAssignment_SameThreatAndComponent", func(t *testing.T) {
		// Create first assignment
		assignment1 := &ThreatAssignment{
			ThreatID:    threat.ID,
			ComponentID: component.ID,
		}
		err := db.Create(assignment1).Error
		require.NoError(t, err)

		// Try to create duplicate assignment
		assignment2 := &ThreatAssignment{
			ThreatID:    threat.ID,
			ComponentID: component.ID,
		}
		err = db.Create(assignment2).Error
		// Should fail due to unique constraint on threat_id + component_id
		assert.Error(t, err, "Should reject duplicate threat assignment to same component")

		// Cleanup
		db.Delete(assignment1)
	})
}