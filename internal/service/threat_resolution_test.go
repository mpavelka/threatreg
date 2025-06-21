package service

import (
	"testing"
	"threatreg/internal/models"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestThreatResolutionService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create test data
	product, err := CreateProduct("Test Product", "A test product")
	require.NoError(t, err)

	instance, err := CreateInstance("Test Instance", product.ID)
	require.NoError(t, err)

	threat, err := CreateThreat("Test Threat", "A test threat")
	require.NoError(t, err)

	// Create threat assignments
	productAssignment, err := AssignThreatToProduct(threat.ID, product.ID)
	require.NoError(t, err)

	instanceAssignment, err := AssignThreatToInstance(threat.ID, instance.ID)
	require.NoError(t, err)

	t.Run("CreateThreatResolution_Product", func(t *testing.T) {
		resolution, err := CreateThreatResolution(
			productAssignment.ID,
			nil,         // instanceID
			&product.ID, // productID
			models.ThreatAssignmentResolutionStatusResolved,
			"Product threat resolved",
		)

		require.NoError(t, err)
		assert.NotNil(t, resolution)
		assert.NotEqual(t, uuid.Nil, resolution.ID)
		assert.Equal(t, productAssignment.ID, resolution.ThreatAssignmentID)
		assert.Equal(t, product.ID, resolution.ProductID)
		assert.Equal(t, uuid.Nil, resolution.InstanceID)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, resolution.Status)
		assert.Equal(t, "Product threat resolved", resolution.Description)

		// Cleanup
		DeleteThreatResolution(resolution.ID)
	})

	t.Run("CreateThreatResolution_Instance", func(t *testing.T) {
		resolution, err := CreateThreatResolution(
			instanceAssignment.ID,
			&instance.ID, // instanceID
			nil,          // productID
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Instance threat awaiting resolution",
		)

		require.NoError(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, instanceAssignment.ID, resolution.ThreatAssignmentID)
		assert.Equal(t, uuid.Nil, resolution.ProductID)
		assert.Equal(t, instance.ID, resolution.InstanceID)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusAwaiting, resolution.Status)

		// Cleanup
		DeleteThreatResolution(resolution.ID)
	})

	t.Run("CreateThreatResolution_NeitherProvided", func(t *testing.T) {
		resolution, err := CreateThreatResolution(
			productAssignment.ID,
			nil, // instanceID
			nil, // productID
			models.ThreatAssignmentResolutionStatusResolved,
			"Invalid resolution",
		)

		assert.Error(t, err)
		assert.Nil(t, resolution)
		assert.Contains(t, err.Error(), "either instanceID or productID must be provided")
	})

	t.Run("GetThreatResolution", func(t *testing.T) {
		// Create resolution
		created, err := CreateThreatResolution(
			productAssignment.ID,
			nil,
			&product.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Test description",
		)
		require.NoError(t, err)

		// Get resolution
		retrieved, err := GetThreatResolution(created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.ThreatAssignmentID, retrieved.ThreatAssignmentID)
		assert.Equal(t, created.Status, retrieved.Status)

		// Cleanup
		DeleteThreatResolution(created.ID)
	})

	t.Run("GetThreatResolution_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()
		resolution, err := GetThreatResolution(nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, resolution)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("GetThreatResolutionByThreatAssignmentID", func(t *testing.T) {
		// Create resolution
		created, err := CreateThreatResolution(
			productAssignment.ID,
			nil,
			&product.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Resolved threat",
		)
		require.NoError(t, err)

		// Get by threat assignment ID
		retrieved, err := GetThreatResolutionByThreatAssignmentID(productAssignment.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, productAssignment.ID, retrieved.ThreatAssignmentID)

		// Cleanup
		DeleteThreatResolution(created.ID)
	})

	t.Run("UpdateThreatResolution", func(t *testing.T) {
		// Create resolution
		created, err := CreateThreatResolution(
			productAssignment.ID,
			nil,
			&product.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Original description",
		)
		require.NoError(t, err)

		// Update status and description
		newStatus := models.ThreatAssignmentResolutionStatusResolved
		newDescription := "Updated description"
		updated, err := UpdateThreatResolution(created.ID, &newStatus, &newDescription)

		require.NoError(t, err)
		assert.Equal(t, created.ID, updated.ID)
		assert.Equal(t, newStatus, updated.Status)
		assert.Equal(t, newDescription, updated.Description)

		// Verify update was persisted
		retrieved, err := GetThreatResolution(created.ID)
		require.NoError(t, err)
		assert.Equal(t, newStatus, retrieved.Status)
		assert.Equal(t, newDescription, retrieved.Description)

		// Cleanup
		DeleteThreatResolution(created.ID)
	})

	t.Run("UpdateThreatResolution_PartialUpdate", func(t *testing.T) {
		// Create resolution
		created, err := CreateThreatResolution(
			productAssignment.ID,
			nil,
			&product.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Original description",
		)
		require.NoError(t, err)

		// Update only status
		newStatus := models.ThreatAssignmentResolutionStatusAccepted
		updated, err := UpdateThreatResolution(created.ID, &newStatus, nil)

		require.NoError(t, err)
		assert.Equal(t, newStatus, updated.Status)
		assert.Equal(t, "Original description", updated.Description) // Should remain unchanged

		// Cleanup
		DeleteThreatResolution(created.ID)
	})

	t.Run("UpdateThreatResolution_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()
		newStatus := models.ThreatAssignmentResolutionStatusResolved
		resolution, err := UpdateThreatResolution(nonExistentID, &newStatus, nil)

		assert.Error(t, err)
		assert.Nil(t, resolution)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("ListThreatResolutionsByProductID", func(t *testing.T) {
		// Create multiple resolutions for the product
		resolution1, err := CreateThreatResolution(
			productAssignment.ID,
			nil,
			&product.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Resolution 1",
		)
		require.NoError(t, err)

		// Create another threat assignment for the same product
		threat2, err := CreateThreat("Another Test Threat", "Another test threat")
		require.NoError(t, err)
		assignment2, err := AssignThreatToProduct(threat2.ID, product.ID)
		require.NoError(t, err)

		resolution2, err := CreateThreatResolution(
			assignment2.ID,
			nil,
			&product.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution 2",
		)
		require.NoError(t, err)

		// List resolutions for product
		resolutions, err := ListThreatResolutionsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, resolutions, 2)

		// Verify both resolutions are present
		resolutionIDs := []uuid.UUID{resolutions[0].ID, resolutions[1].ID}
		assert.Contains(t, resolutionIDs, resolution1.ID)
		assert.Contains(t, resolutionIDs, resolution2.ID)

		// Cleanup
		DeleteThreatResolution(resolution1.ID)
		DeleteThreatResolution(resolution2.ID)
	})

	t.Run("ListThreatResolutionsByInstanceID", func(t *testing.T) {
		// Create resolution for instance
		resolution, err := CreateThreatResolution(
			instanceAssignment.ID,
			&instance.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Instance resolution",
		)
		require.NoError(t, err)

		// List resolutions for instance
		resolutions, err := ListThreatResolutionsByInstanceID(instance.ID)
		require.NoError(t, err)
		assert.Len(t, resolutions, 1)
		assert.Equal(t, resolution.ID, resolutions[0].ID)

		// Cleanup
		DeleteThreatResolution(resolution.ID)
	})

	t.Run("DeleteThreatResolution", func(t *testing.T) {
		// Create a new threat assignment for this test
		threat3, err := CreateThreat("Delete Test Threat", "Threat for delete test")
		require.NoError(t, err)
		assignment3, err := AssignThreatToProduct(threat3.ID, product.ID)
		require.NoError(t, err)

		// Create resolution
		created, err := CreateThreatResolution(
			assignment3.ID,
			nil,
			&product.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"To be deleted",
		)
		require.NoError(t, err)

		// Delete resolution
		err = DeleteThreatResolution(created.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = GetThreatResolution(created.ID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteThreatResolution_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := DeleteThreatResolution(nonExistentID)

		// Delete should succeed even if resolution doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})
}
