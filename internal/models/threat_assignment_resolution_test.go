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

	// Get database instance
	db := database.GetDB()
	require.NotNil(t, db, "Database connection should not be nil")

	// Run migrations for all models
	err = db.AutoMigrate(
		&Component{},
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
	component := createTestComponent(t, ComponentTypeProduct)
	instanceComponent := createTestComponent(t, ComponentTypeInstance)
	threat := createTestThreat(t)
	threatAssignment := createTestThreatAssignment(t, threat.ID, component.ID)

	t.Run("ValidResolution_ComponentOnly", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ComponentID:        component.ID,
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.NoError(t, err, "Should allow resolution with ComponentID set")

		// Verify resolution was created
		retrieved, err := repo.GetByID(nil, resolution.ID)
		assert.NoError(t, err)
		assert.Equal(t, threatAssignment.ID, retrieved.ThreatAssignmentID)
		assert.Equal(t, component.ID, retrieved.ComponentID)
		assert.Equal(t, ThreatAssignmentResolutionStatusResolved, retrieved.Status)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})

	t.Run("ValidResolution_InstanceComponent", func(t *testing.T) {
		// Create separate threat assignment for instance component
		instanceThreatAssignment := createTestThreatAssignment(t, threat.ID, instanceComponent.ID)
		
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: instanceThreatAssignment.ID,
			ComponentID:        instanceComponent.ID,
			Status:             ThreatAssignmentResolutionStatusAwaiting,
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.NoError(t, err, "Should allow resolution with ComponentID set")

		// Verify resolution was created
		retrieved, err := repo.GetByID(nil, resolution.ID)
		assert.NoError(t, err)
		assert.Equal(t, instanceThreatAssignment.ID, retrieved.ThreatAssignmentID)
		assert.Equal(t, instanceComponent.ID, retrieved.ComponentID)
		assert.Equal(t, ThreatAssignmentResolutionStatusAwaiting, retrieved.Status)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})

	t.Run("InvalidResolution_NilComponentID", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ComponentID:        uuid.Nil, // Nil ComponentID - should fail
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Test resolution",
		}

		repo := NewThreatAssignmentResolutionRepository(getTestDB(t))
		err := repo.Create(nil, resolution)
		assert.Error(t, err, "Should reject resolution with nil ComponentID")
	})


	t.Run("InvalidStatus", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ComponentID:        component.ID,
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
				ComponentID:        component.ID,
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
				threatAssignment = createTestThreatAssignment(t, newThreat.ID, component.ID)
			}
		}
	})
}

func TestThreatAssignmentResolutionRepository_Integration(t *testing.T) {
	cleanup := setupTestDBWithResolution(t)
	defer cleanup()

	// Create test data
	component := createTestComponent(t, ComponentTypeProduct)
	instanceComponent := createTestComponent(t, ComponentTypeInstance)
	threat := createTestThreat(t)
	threatAssignment := createTestThreatAssignment(t, threat.ID, component.ID)

	repo := NewThreatAssignmentResolutionRepository(getTestDB(t))

	t.Run("CreateAndGet", func(t *testing.T) {
		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ComponentID:        component.ID,
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
		assert.Equal(t, resolution.ComponentID, retrieved.ComponentID)
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
			ComponentID:        component.ID,
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

	t.Run("ListByComponentID", func(t *testing.T) {
		// Create resolutions for the component
		resolution1 := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment.ID,
			ComponentID:        component.ID,
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Resolution 1",
		}
		err := repo.Create(nil, resolution1)
		require.NoError(t, err)

		// Create another threat assignment for second resolution
		threat2 := createTestThreat(t)
		threatAssignment2 := createTestThreatAssignment(t, threat2.ID, component.ID)
		resolution2 := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment2.ID,
			ComponentID:        component.ID,
			Status:             ThreatAssignmentResolutionStatusAwaiting,
			Description:        "Resolution 2",
		}
		err = repo.Create(nil, resolution2)
		require.NoError(t, err)

		// List by ComponentID
		resolutions, err := repo.ListByComponentID(nil, component.ID)
		require.NoError(t, err)
		assert.Len(t, resolutions, 2)

		// Cleanup
		repo.Delete(nil, resolution1.ID)
		repo.Delete(nil, resolution2.ID)
	})

	t.Run("ListByInstanceComponentID", func(t *testing.T) {
		// Create threat assignment for instance component
		threatAssignmentForInstance := createTestThreatAssignment(t, threat.ID, instanceComponent.ID)

		resolution := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignmentForInstance.ID,
			ComponentID:        instanceComponent.ID,
			Status:             ThreatAssignmentResolutionStatusResolved,
			Description:        "Instance resolution",
		}
		err := repo.Create(nil, resolution)
		require.NoError(t, err)

		// List by ComponentID
		resolutions, err := repo.ListByComponentID(nil, instanceComponent.ID)
		require.NoError(t, err)
		assert.Len(t, resolutions, 1)
		assert.Equal(t, resolution.ID, resolutions[0].ID)

		// Cleanup
		repo.Delete(nil, resolution.ID)
	})
}

// Helper functions for creating test data
func createTestComponent(t *testing.T, componentType ComponentType) *Component {
	component := &Component{
		Name:        "Test Component",
		Description: "A test component",
		Type:        componentType,
	}
	err := getTestDB(t).Create(component).Error
	require.NoError(t, err)
	return component
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

func createTestThreatAssignment(t *testing.T, threatID, componentID uuid.UUID) *ThreatAssignment {
	assignment := &ThreatAssignment{
		ThreatID:    threatID,
		ComponentID: componentID,
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
	component := createTestComponent(t, ComponentTypeProduct)
	threat := createTestThreat(t)
	threatAssignment := createTestThreatAssignment(t, threat.ID, component.ID)

	// Create two threat resolutions for delegation
	resolution1 := &ThreatAssignmentResolution{
		ThreatAssignmentID: threatAssignment.ID,
		ComponentID:        component.ID,
		Status:             ThreatAssignmentResolutionStatusAwaiting,
		Description:        "Source resolution",
	}

	threat2 := createTestThreat(t)
	threatAssignment2 := createTestThreatAssignment(t, threat2.ID, component.ID)
	resolution2 := &ThreatAssignmentResolution{
		ThreatAssignmentID: threatAssignment2.ID,
		ComponentID:        component.ID,
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
		threatAssignment4 := createTestThreatAssignment(t, threat4.ID, component.ID)
		resolution4 := &ThreatAssignmentResolution{
			ThreatAssignmentID: threatAssignment4.ID,
			ComponentID:        component.ID,
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
