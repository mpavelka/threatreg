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

	// Run migrations for all models
	db := database.GetDB()
	require.NotNil(t, db, "Database connection should not be nil")

	err = db.AutoMigrate(
		&Product{},
		&Instance{},
		&Threat{},
		&Control{},
		&ThreatAssignment{},
		&ControlAssignment{},
		&ThreatControl{},
	)
	require.NoError(t, err, "Failed to run migrations")

	// Return cleanup function
	return func() {
		database.Close()
		config.AppConfig = originalConfig
		os.RemoveAll(tempDir)
	}
}

func TestThreatAssignmentConstraints_Integration(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	db := database.GetDB()
	require.NotNil(t, db)

	// Create test data
	product := &Product{
		Name:        "Test Product",
		Description: "A test product",
	}
	require.NoError(t, db.Create(product).Error)

	instance := &Instance{
		Name:       "Test Instance",
		InstanceOf: product.ID,
	}
	require.NoError(t, db.Create(instance).Error)

	threat := &Threat{
		Title:       "Test Threat",
		Description: "A test threat",
	}
	require.NoError(t, db.Create(threat).Error)

	t.Run("ValidAssignment_ProductOnly", func(t *testing.T) {
		// Test creating assignment with only ProductID set
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  product.ID,
			InstanceID: uuid.Nil, // Explicitly nil
		}

		err := db.Create(assignment).Error
		assert.NoError(t, err, "Should allow assignment with only ProductID set")

		// Verify assignment was created
		var retrieved ThreatAssignment
		err = db.First(&retrieved, "id = ?", assignment.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, threat.ID, retrieved.ThreatID)
		assert.Equal(t, product.ID, retrieved.ProductID)
		assert.Equal(t, uuid.Nil, retrieved.InstanceID)

		// Cleanup
		db.Delete(assignment)
	})

	t.Run("ValidAssignment_InstanceOnly", func(t *testing.T) {
		// Test creating assignment with only InstanceID set
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  uuid.Nil, // Explicitly nil
			InstanceID: instance.ID,
		}

		err := db.Create(assignment).Error
		assert.NoError(t, err, "Should allow assignment with only InstanceID set")

		// Verify assignment was created
		var retrieved ThreatAssignment
		err = db.First(&retrieved, "id = ?", assignment.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, threat.ID, retrieved.ThreatID)
		assert.Equal(t, uuid.Nil, retrieved.ProductID)
		assert.Equal(t, instance.ID, retrieved.InstanceID)

		// Cleanup
		db.Delete(assignment)
	})

	t.Run("InvalidAssignment_BothSet", func(t *testing.T) {
		// Test creating assignment with both ProductID and InstanceID set
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  product.ID,
			InstanceID: instance.ID, // Both set - should fail
		}

		err := db.Create(assignment).Error
		assert.Error(t, err, "Should reject assignment with both ProductID and InstanceID set")
		assert.Contains(t, err.Error(), "cannot have both ProductID and InstanceID set")
	})

	t.Run("InvalidAssignment_NeitherSet", func(t *testing.T) {
		// Test creating assignment with neither ProductID nor InstanceID set
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  uuid.Nil,
			InstanceID: uuid.Nil, // Both nil - should fail
		}

		err := db.Create(assignment).Error
		assert.Error(t, err, "Should reject assignment with neither ProductID nor InstanceID set")
		assert.Contains(t, err.Error(), "must have either ProductID or InstanceID set")
	})

	t.Run("InvalidUpdate_BothSet", func(t *testing.T) {
		// Create valid assignment first
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  product.ID,
			InstanceID: uuid.Nil,
		}
		require.NoError(t, db.Create(assignment).Error)

		// Try to update to invalid state (both set)
		assignment.InstanceID = instance.ID
		err := db.Save(assignment).Error
		assert.Error(t, err, "Should reject update that sets both ProductID and InstanceID")
		assert.Contains(t, err.Error(), "cannot have both ProductID and InstanceID set")

		// Cleanup
		db.Delete(assignment)
	})

	t.Run("InvalidUpdate_NeitherSet", func(t *testing.T) {
		// Create valid assignment first
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  product.ID,
			InstanceID: uuid.Nil,
		}
		require.NoError(t, db.Create(assignment).Error)

		// Try to update to invalid state (neither set)
		assignment.ProductID = uuid.Nil
		err := db.Save(assignment).Error
		assert.Error(t, err, "Should reject update that clears both ProductID and InstanceID")
		assert.Contains(t, err.Error(), "must have either ProductID or InstanceID set")

		// Cleanup - reset to valid state first
		assignment.ProductID = product.ID
		db.Save(assignment)
		db.Delete(assignment)
	})

	t.Run("ValidUpdate_SwitchFromProductToInstance", func(t *testing.T) {
		// Create assignment with ProductID
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  product.ID,
			InstanceID: uuid.Nil,
		}
		require.NoError(t, db.Create(assignment).Error)

		// Switch to InstanceID (valid transition)
		assignment.ProductID = uuid.Nil
		assignment.InstanceID = instance.ID
		err := db.Save(assignment).Error
		assert.NoError(t, err, "Should allow switching from ProductID to InstanceID")

		// Verify the update
		var retrieved ThreatAssignment
		err = db.First(&retrieved, "id = ?", assignment.ID).Error
		require.NoError(t, err)
		assert.Equal(t, uuid.Nil, retrieved.ProductID)
		assert.Equal(t, instance.ID, retrieved.InstanceID)

		// Cleanup
		db.Delete(assignment)
	})

	t.Run("ValidUpdate_SwitchFromInstanceToProduct", func(t *testing.T) {
		// Create assignment with InstanceID
		assignment := &ThreatAssignment{
			ThreatID:   threat.ID,
			ProductID:  uuid.Nil,
			InstanceID: instance.ID,
		}
		require.NoError(t, db.Create(assignment).Error)

		// Switch to ProductID (valid transition)
		assignment.InstanceID = uuid.Nil
		assignment.ProductID = product.ID
		err := db.Save(assignment).Error
		assert.NoError(t, err, "Should allow switching from InstanceID to ProductID")

		// Verify the update
		var retrieved ThreatAssignment
		err = db.First(&retrieved, "id = ?", assignment.ID).Error
		require.NoError(t, err)
		assert.Equal(t, product.ID, retrieved.ProductID)
		assert.Equal(t, uuid.Nil, retrieved.InstanceID)

		// Cleanup
		db.Delete(assignment)
	})
}

func TestThreatAssignmentRepository_Integration(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	db := database.GetDB()
	require.NotNil(t, db)

	// Create test data
	product := &Product{
		Name:        "Test Product",
		Description: "A test product",
	}
	require.NoError(t, db.Create(product).Error)

	instance := &Instance{
		Name:       "Test Instance",
		InstanceOf: product.ID,
	}
	require.NoError(t, db.Create(instance).Error)

	threat := &Threat{
		Title:       "Test Threat",
		Description: "A test threat",
	}
	require.NoError(t, db.Create(threat).Error)

	repo := NewThreatAssignmentRepository(db)

	t.Run("AssignThreatToProduct_Success", func(t *testing.T) {
		assignment, err := repo.AssignThreatToProduct(nil, threat.ID, product.ID)

		assert.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, product.ID, assignment.ProductID)
		assert.Equal(t, uuid.Nil, assignment.InstanceID)

		// Cleanup
		repo.Delete(nil, assignment.ID)
	})

	t.Run("AssignThreatToProduct_Duplicate", func(t *testing.T) {
		// Create first assignment
		assignment1, err := repo.AssignThreatToProduct(nil, threat.ID, product.ID)
		require.NoError(t, err)

		// Try to create duplicate assignment
		assignment2, err := repo.AssignThreatToProduct(nil, threat.ID, product.ID)

		// Should return existing assignment without error
		assert.NoError(t, err)
		assert.NotNil(t, assignment2)
		assert.Equal(t, assignment1.ID, assignment2.ID)

		// Cleanup
		repo.Delete(nil, assignment1.ID)
	})

	t.Run("AssignThreatToInstance_Success", func(t *testing.T) {
		assignment, err := repo.AssignThreatToInstance(nil, threat.ID, instance.ID)

		assert.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, uuid.Nil, assignment.ProductID)
		assert.Equal(t, instance.ID, assignment.InstanceID)

		// Cleanup
		repo.Delete(nil, assignment.ID)
	})

	t.Run("AssignThreatToInstance_Duplicate", func(t *testing.T) {
		// Create first assignment
		assignment1, err := repo.AssignThreatToInstance(nil, threat.ID, instance.ID)
		require.NoError(t, err)

		// Try to create duplicate assignment
		assignment2, err := repo.AssignThreatToInstance(nil, threat.ID, instance.ID)

		// Should return existing assignment without error
		assert.NoError(t, err)
		assert.NotNil(t, assignment2)
		assert.Equal(t, assignment1.ID, assignment2.ID)

		// Cleanup
		repo.Delete(nil, assignment1.ID)
	})
}
