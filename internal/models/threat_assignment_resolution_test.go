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
	"gorm.io/gorm"
)

// setupTestDBWithResolution creates a test database with resolution model
func setupTestDBWithResolution(t *testing.T) func() {
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
		&ThreatAssignmentResolution{},
		&ThreatAssignmentResolutionDelegation{},
	)
	require.NoError(t, err, "Failed to run migrations")

	// Return cleanup function
	return func() {
		database.Close()
		config.AppConfig = originalConfig
		os.RemoveAll(tempDir)
	}
}

func TestThreatAssignmentResolutionConstraints_Integration(t *testing.T) {
	cleanup := setupTestDBWithResolution(t)
	defer cleanup()

	// Create test data
	product := createTestProduct(t)
	instance := createTestInstance(t, product.ID)
	threat := createTestThreat(t)
	threatAssignment := createTestThreatAssignment(t, threat.ID, product.ID, uuid.Nil)

	t.Run("ValidResolution_ProductOnly", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          product.ID,
			InstanceID:         uuid.Nil,
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.NoError(t, err, "Should allow resolution with only ProductID set")

		// Verify resolution was created
		retrieved, err := repo.GetByID(nil, resolution.ID)
		assert.NoError(t, err)
		assert.Equal(t, threatAssignment.ID, retrieved.ThreatAssignmentID)
		assert.Equal(t, product.ID, retrieved.ProductID)
		assert.Equal(t, uuid.Nil, retrieved.InstanceID)
		assert.Equal(t, ThreatAssignmentResolutionStatusResolved, retrieved.Status)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})

	t.Run("ValidResolution_InstanceOnly", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          uuid.Nil,
			InstanceID:         instance.ID,
			Status:             ThreatAssignmentResolutionStatusAwaiting,
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.NoError(t, err, "Should allow resolution with only InstanceID set")

		// Verify resolution was created
		retrieved, err := repo.GetByID(nil, resolution.ID)
		assert.NoError(t, err)
		assert.Equal(t, threatAssignment.ID, retrieved.ThreatAssignmentID)
		assert.Equal(t, uuid.Nil, retrieved.ProductID)
		assert.Equal(t, instance.ID, retrieved.InstanceID)
		assert.Equal(t, ThreatAssignmentResolutionStatusAwaiting, retrieved.Status)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})

	t.Run("InvalidResolution_BothSet", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          product.ID,
			InstanceID:         instance.ID, // Both set - should fail
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.Error(t, err, "Should reject resolution with both ProductID and InstanceID set")
		assert.Contains(t, err.Error(), "cannot have both ProductID and InstanceID set")
	})

	t.Run("InvalidResolution_NeitherSet", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          uuid.Nil,
			InstanceID:         uuid.Nil, // Both nil - should fail
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.Error(t, err, "Should reject resolution with neither ProductID nor InstanceID set")
		assert.Contains(t, err.Error(), "must have either ProductID or InstanceID set")
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          product.ID,
			InstanceID:         uuid.Nil,
			Status:             "invalid_status",
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.Error(t, err, "Should reject resolution with invalid status")
		assert.Contains(t, err.Error(), "invalid status")
	})

	t.Run("ValidStatuses", func(t *testing.T) {
		validStatuses := []ThreatAssignmentResolutionStatus{
			ThreatAssignmentResolutionStatusResolved,
			ThreatAssignmentResolutionStatusAwaiting,
			ThreatAssignmentResolutionStatusAccepted,
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))

		for i, status := range validStatuses {
			resolution := &ThreatAssignmentResolution{
				ThreatAssignmentID: threatAssignment.ID,
				ProductID:          product.ID,
				InstanceID:         uuid.Nil,
				Status:             status,
				Description:        "Test resolution",
			}

			err := repo.Create(nil, resolution)
			assert.NoError(t, err, "Should allow valid status: %s", status)

			// Cleanup
			repo.Delete(nil, resolution.ID)

			// Avoid unique constraint violations by using different threats
			if i < len(validStatuses)-1 {
				newThreat := createTestThreat(t)
				threatAssignment = createTestThreatAssignment(t, newThreat.ID, product.ID, uuid.Nil)
			}
		}
	})
}

func TestThreatAssignmentResolutionRepository_Integration(t *testing.T) {
	cleanup := setupTestDBWithResolution(t)
	defer cleanup()

	// Create test data
	product := createTestProduct(t)
	instance := createTestInstance(t, product.ID)
	threat := createTestThreat(t)
	threatAssignment := createTestThreatAssignment(t, threat.ID, product.ID, uuid.Nil)

	repo := NewThreatAssignmentResolutionRepository(getTestDB(t))

	t.Run("CreateAndGet", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          product.ID,
			InstanceID:         uuid.Nil,
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Test resolution description",
		}

		// Create
		err := repo.Create(nil, resolution)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, resolution.ID)

		// Get by ID
		retrieved, err := repo.GetByID(nil, resolution.ID)
		require.NoError(t, err)
		assert.Equal(t, resolution.ThreatAssignmentID, retrieved.ThreatAssignmentID)
		assert.Equal(t, resolution.ProductID, retrieved.ProductID)
		assert.Equal(t, resolution.Status, retrieved.Status)
		assert.Equal(t, resolution.Description, retrieved.Description)

		// Get by ThreatAssignmentID
		byThreatAssignment, err := repo.GetOneByThreatAssignmentID(nil, threatAssignment.ID)
		require.NoError(t, err)
		assert.Equal(t, resolution.ID, byThreatAssignment.ID)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})

	t.Run("Update", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          product.ID,
			InstanceID:         uuid.Nil,
			Status:             ThreatAssignmentResolutionStatusAwaiting,
			Description:        "Original description",
		}

		err := repo.Create(nil, resolution)
		require.NoError(t, err)

		// Update
		resolution.Status = ThreatAssignmentResolutionStatusAccepted
		resolution.Description = "Updated description"
		err = repo.Update(nil, resolution)
		require.NoError(t, err)

		// Verify update
		retrieved, err := repo.GetByID(nil, resolution.ID)
		require.NoError(t, err)
		assert.Equal(t, ThreatAssignmentResolutionStatusAccepted, retrieved.Status)
		assert.Equal(t, "Updated description", retrieved.Description)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})

	t.Run("ListByProductID", func(t *testing.T) {
		// Create resolutions for the product
		resolution1 := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ProductID:          product.ID,
			InstanceID:         uuid.Nil,
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Resolution 1",
		}
		err := repo.Create(nil, resolution1)
		require.NoError(t, err)

		// Create another threat assignment for second resolution
		threat2 := createTestThreat(t)
		threatAssignment2 := createTestThreatAssignment(t, threat2.ID, product.ID, uuid.Nil)
		resolution2 := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment2.ID,
			ProductID:          product.ID,
			InstanceID:         uuid.Nil,
			Status:             ThreatAssignmentResolutionStatusAwaiting,
			Description:        "Resolution 2",
		}
		err = repo.Create(nil, resolution2)
		require.NoError(t, err)

		// List by ProductID
		resolutions, err := repo.ListByProductID(nil, product.ID)
		require.NoError(t, err)
		assert.Len(t, resolutions, 2)

		// Cleanup
		repo.Delete(nil, resolution1.ID)
		repo.Delete(nil, resolution2.ID)
	})

	t.Run("ListByInstanceID", func(t *testing.T) {
		// Create threat assignment for instance
		threatAssignmentForInstance := createTestThreatAssignment(t, threat.ID, uuid.Nil, instance.ID)

		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignmentForInstance.ID,
			ProductID:          uuid.Nil,
			InstanceID:         instance.ID,
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Instance resolution",
		}
		err := repo.Create(nil, resolution)
		require.NoError(t, err)

		// List by InstanceID
		resolutions, err := repo.ListByInstanceID(nil, instance.ID)
		require.NoError(t, err)
		assert.Len(t, resolutions, 1)
		assert.Equal(t, resolution.ID, resolutions[0].ID)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})
}

// Helper functions for creating test data
func createTestProduct(t *testing.T) *Product {
	product := &Product{
		Name:        "Test Product",
		Description: "A test product",
	}
	err := getTestDB(t).Create(product).Error
	require.NoError(t, err)
	return product
}

func createTestInstance(t *testing.T, productID uuid.UUID) *Instance {
	instance := &Instance{
		Name:       "Test Instance",
		InstanceOf: productID,
	}
	err := getTestDB(t).Create(instance).Error
	require.NoError(t, err)
	return instance
}

func createTestThreat(t *testing.T) *Threat {
	threat := &Threat{
		Title:       "Test Threat " + uuid.New().String()[:8],
		Description: "A test threat",
	}
	err := getTestDB(t).Create(threat).Error
	require.NoError(t, err)
	return threat
}

func createTestThreatAssignment(t *testing.T, threatID, productID, instanceID uuid.UUID) *ThreatAssignment {
	assignment := &ThreatAssignment{
		ThreatID:   threatID,
		ProductID:  productID,
		InstanceID: instanceID,
	}
	err := getTestDB(t).Create(assignment).Error
	require.NoError(t, err)
	return assignment
}

func getTestDB(t *testing.T) *gorm.DB {
	return database.GetDB()
}

func TestThreatAssignmentResolutionDelegation_Integration(t *testing.T) {
	cleanup := setupTestDBWithResolution(t)
	defer cleanup()

	// Create test data
	product := createTestProduct(t)
	threat := createTestThreat(t)
	threatAssignment := createTestThreatAssignment(t, threat.ID, product.ID, uuid.Nil)

	// Create two threat resolutions for delegation
	resolution1 := &ThreatAssignmentResolution{
		ThreatAssignmentID: threatAssignment.ID,
		ProductID:          product.ID,
		InstanceID:         uuid.Nil,
		Status:             ThreatAssignmentResolutionStatusAwaiting,
		Description:        "Source resolution",
	}

	threat2 := createTestThreat(t)
	threatAssignment2 := createTestThreatAssignment(t, threat2.ID, product.ID, uuid.Nil)
	resolution2 := &ThreatAssignmentResolution{
		ThreatAssignmentID: threatAssignment2.ID,
		ProductID:          product.ID,
		InstanceID:         uuid.Nil,
		Status:             ThreatAssignmentResolutionStatusAwaiting,
		Description:        "Target resolution",
	}

	resolutionRepo := NewThreatAssignmentResolutionRepository(getTestDB(t))
	delegationRepo := NewThreatAssignmentResolutionDelegationRepository(getTestDB(t))

	// Create resolutions
	err := resolutionRepo.Create(nil, resolution1)
	require.NoError(t, err)
	err = resolutionRepo.Create(nil, resolution2)
	require.NoError(t, err)

	t.Run("TestUniqueDelegationConstraint", func(t *testing.T) {
		// Create new resolutions for this test
		threat4 := createTestThreat(t)
		threatAssignment4 := createTestThreatAssignment(t, threat4.ID, product.ID, uuid.Nil)
		resolution4 := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment4.ID,
			ProductID:          product.ID,
			InstanceID:         uuid.Nil,
			Status:             ThreatAssignmentResolutionStatusAwaiting,
			Description:        "Resolution for unique test",
		}

		err := resolutionRepo.Create(nil, resolution4)
		require.NoError(t, err)

		// Create first delegation
		delegation1 := &ThreatAssignmentResolutionDelegation{
			DelegatedBy: resolution4.ID,
			DelegatedTo: resolution2.ID,
		}

		err = delegationRepo.CreateThreatAssignmentResolutionDelegation(nil, delegation1)
		require.NoError(t, err)

		// Try to create duplicate delegation - should fail
		delegation2 := &ThreatAssignmentResolutionDelegation{
			DelegatedBy: resolution4.ID,
			DelegatedTo: resolution2.ID,
		}

		err = delegationRepo.CreateThreatAssignmentResolutionDelegation(nil, delegation2)
		assert.Error(t, err, "Should not allow duplicate delegations")

		// Cleanup
		delegationRepo.DeleteThreatAssignmentResolutionDelegation(nil, delegation1.ID)
		resolutionRepo.Delete(nil, resolution4.ID)
	})

	// Cleanup
	resolutionRepo.Delete(nil, resolution1.ID)
	resolutionRepo.Delete(nil, resolution2.ID)
}
