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

func TestDelegationFunctionality(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("DelegateResolution_UpdatesStatus", func(t *testing.T) {
		product, _ := CreateProduct("Test Product", "A test product")
		instance1, _ := CreateInstance("Instance 1", product.ID)
		instance2, _ := CreateInstance("Instance 2", product.ID)
		threat, _ := CreateThreat("Test Threat", "A test threat")
		assignment1, _ := AssignThreatToInstance(threat.ID, instance1.ID)
		assignment2, _ := AssignThreatToInstance(threat.ID, instance2.ID)

		source, _ := CreateThreatResolution(assignment1.ID, &instance1.ID, nil, models.ThreatAssignmentResolutionStatusAwaiting, "Source")
		target, _ := CreateThreatResolution(assignment2.ID, &instance2.ID, nil, models.ThreatAssignmentResolutionStatusResolved, "Target")
		defer func() { DeleteThreatResolution(source.ID); DeleteThreatResolution(target.ID) }()

		err := DelegateResolution(*source, *target)
		require.NoError(t, err)

		updatedSource, _ := GetThreatResolution(source.ID)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, updatedSource.Status)
	})

	t.Run("DelegateResolution_MultipleTimes", func(t *testing.T) {
		product, _ := CreateProduct("Test Product 2", "A test product")
		instance1, _ := CreateInstance("Instance 1", product.ID)
		instance2, _ := CreateInstance("Instance 2", product.ID)
		instance3, _ := CreateInstance("Instance 3", product.ID)
		threat, _ := CreateThreat("Test Threat 2", "A test threat")
		assignment1, _ := AssignThreatToInstance(threat.ID, instance1.ID)
		assignment2, _ := AssignThreatToInstance(threat.ID, instance2.ID)
		assignment3, _ := AssignThreatToInstance(threat.ID, instance3.ID)

		source, err := CreateThreatResolution(assignment1.ID, &instance1.ID, nil, models.ThreatAssignmentResolutionStatusAwaiting, "Source")
		require.NoError(t, err)
		target1, err := CreateThreatResolution(assignment2.ID, &instance2.ID, nil, models.ThreatAssignmentResolutionStatusAccepted, "Target1")
		require.NoError(t, err)
		target2, err := CreateThreatResolution(assignment3.ID, &instance3.ID, nil, models.ThreatAssignmentResolutionStatusResolved, "Target2")
		require.NoError(t, err)
		defer func() {
			DeleteThreatResolution(source.ID)
			DeleteThreatResolution(target1.ID)
			DeleteThreatResolution(target2.ID)
		}()

		err = DelegateResolution(*source, *target1)
		require.NoError(t, err)

		delegated1, err := GetDelegatedToResolutionByDelegatedByID(source.ID)
		require.NoError(t, err)
		require.NotNil(t, delegated1)
		assert.Equal(t, target1.ID, delegated1.ID)

		err = DelegateResolution(*source, *target2)
		require.NoError(t, err)

		delegated2, err := GetDelegatedToResolutionByDelegatedByID(source.ID)
		require.NoError(t, err)
		require.NotNil(t, delegated2)
		assert.Equal(t, target2.ID, delegated2.ID)

		updatedSource, err := GetThreatResolution(source.ID)
		require.NoError(t, err)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, updatedSource.Status)
	})

	t.Run("UpdateThreatResolution_RemovesDelegation", func(t *testing.T) {
		product, _ := CreateProduct("Test Product 3", "A test product")
		instance1, _ := CreateInstance("Instance 1", product.ID)
		instance2, _ := CreateInstance("Instance 2", product.ID)
		threat, _ := CreateThreat("Test Threat 3", "A test threat")
		assignment1, _ := AssignThreatToInstance(threat.ID, instance1.ID)
		assignment2, _ := AssignThreatToInstance(threat.ID, instance2.ID)

		source, err := CreateThreatResolution(assignment1.ID, &instance1.ID, nil, models.ThreatAssignmentResolutionStatusAwaiting, "Source")
		require.NoError(t, err)
		target, err := CreateThreatResolution(assignment2.ID, &instance2.ID, nil, models.ThreatAssignmentResolutionStatusResolved, "Target")
		require.NoError(t, err)
		defer func() { DeleteThreatResolution(source.ID); DeleteThreatResolution(target.ID) }()

		err = DelegateResolution(*source, *target)
		require.NoError(t, err)

		delegated, err := GetDelegatedToResolutionByDelegatedByID(source.ID)
		require.NoError(t, err)
		require.NotNil(t, delegated)

		newStatus := models.ThreatAssignmentResolutionStatusAccepted
		newDesc := "Updated"
		_, err = UpdateThreatResolution(source.ID, &newStatus, &newDesc)
		require.NoError(t, err)

		delegatedAfter, err := GetDelegatedToResolutionByDelegatedByID(source.ID)
		require.NoError(t, err)
		assert.Nil(t, delegatedAfter)

		updated, err := GetThreatResolution(source.ID)
		require.NoError(t, err)
		assert.Equal(t, newStatus, updated.Status)
		assert.Equal(t, newDesc, updated.Description)
	})

	t.Run("DelegationChain_StatusPropagation", func(t *testing.T) {
		product, _ := CreateProduct("Test Product 4", "A test product")
		instance1, _ := CreateInstance("Instance 1", product.ID)
		instance2, _ := CreateInstance("Instance 2", product.ID)
		instance3, _ := CreateInstance("Instance 3", product.ID)
		threat, _ := CreateThreat("Test Threat 4", "A test threat")
		assignment1, _ := AssignThreatToInstance(threat.ID, instance1.ID)
		assignment2, _ := AssignThreatToInstance(threat.ID, instance2.ID)
		assignment3, _ := AssignThreatToInstance(threat.ID, instance3.ID)

		// Create chain: resA -> resB -> resC
		resA, err := CreateThreatResolution(assignment1.ID, &instance1.ID, nil, models.ThreatAssignmentResolutionStatusAwaiting, "ResA")
		require.NoError(t, err)
		resB, err := CreateThreatResolution(assignment2.ID, &instance2.ID, nil, models.ThreatAssignmentResolutionStatusAwaiting, "ResB")
		require.NoError(t, err)
		resC, err := CreateThreatResolution(assignment3.ID, &instance3.ID, nil, models.ThreatAssignmentResolutionStatusAwaiting, "ResC")
		require.NoError(t, err)
		defer func() {
			DeleteThreatResolution(resA.ID)
			DeleteThreatResolution(resB.ID)
			DeleteThreatResolution(resC.ID)
		}()

		// Create delegation chain: A -> B -> C
		err = DelegateResolution(*resA, *resB)
		require.NoError(t, err)
		err = DelegateResolution(*resB, *resC)
		require.NoError(t, err)

		// Update final resolution (C) status
		newStatus := models.ThreatAssignmentResolutionStatusResolved
		_, err = UpdateThreatResolution(resC.ID, &newStatus, nil)
		require.NoError(t, err)

		// Verify all resolutions in chain have updated status
		updatedA, err := GetThreatResolution(resA.ID)
		require.NoError(t, err)
		updatedB, err := GetThreatResolution(resB.ID)
		require.NoError(t, err)
		updatedC, err := GetThreatResolution(resC.ID)
		require.NoError(t, err)

		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, updatedA.Status, "ResA should inherit status from chain end")
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, updatedB.Status, "ResB should inherit status from chain end")
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, updatedC.Status, "ResC should have updated status")
	})
}

func TestGetInstanceLevelThreatResolutionWithDelegation(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create test data
	product, _ := CreateProduct("Test Product", "A test product")
	instance1, _ := CreateInstance("Instance 1", product.ID)
	instance2, _ := CreateInstance("Instance 2", product.ID)
	threat, _ := CreateThreat("Test Threat", "A test threat")
	assignment1, _ := AssignThreatToInstance(threat.ID, instance1.ID)
	assignment2, _ := AssignThreatToInstance(threat.ID, instance2.ID)

	// Create two resolutions
	resolution1, err := CreateThreatResolution(assignment1.ID, &instance1.ID, nil, models.ThreatAssignmentResolutionStatusAwaiting, "Source resolution")
	require.NoError(t, err)
	resolution2, err := CreateThreatResolution(assignment2.ID, &instance2.ID, nil, models.ThreatAssignmentResolutionStatusResolved, "Target resolution")
	require.NoError(t, err)

	// Create delegation from resolution1 to resolution2
	err = DelegateResolution(*resolution1, *resolution2)
	require.NoError(t, err)

	// Test GetInstanceLevelThreatResolutionWithDelegation for resolution1 (has delegation)
	resultWithDelegation, err := GetInstanceLevelThreatResolutionWithDelegation(assignment1.ID, instance1.ID)
	require.NoError(t, err)
	require.NotNil(t, resultWithDelegation)

	// Verify resolution data
	assert.Equal(t, resolution1.ID, resultWithDelegation.Resolution.ID)
	assert.Equal(t, assignment1.ID, resultWithDelegation.Resolution.ThreatAssignmentID)
	assert.Equal(t, instance1.ID, resultWithDelegation.Resolution.InstanceID)
	assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, resultWithDelegation.Resolution.Status)

	// Verify delegation data
	require.NotNil(t, resultWithDelegation.Delegation)
	assert.Equal(t, resolution1.ID, resultWithDelegation.Delegation.DelegatedBy)
	assert.Equal(t, resolution2.ID, resultWithDelegation.Delegation.DelegatedTo)

	// Test GetInstanceLevelThreatResolutionWithDelegation for resolution2 (no delegation)
	resultNoDelegation, err := GetInstanceLevelThreatResolutionWithDelegation(assignment2.ID, instance2.ID)
	require.NoError(t, err)
	require.NotNil(t, resultNoDelegation)

	// Verify resolution data
	assert.Equal(t, resolution2.ID, resultNoDelegation.Resolution.ID)
	assert.Equal(t, assignment2.ID, resultNoDelegation.Resolution.ThreatAssignmentID)
	assert.Equal(t, instance2.ID, resultNoDelegation.Resolution.InstanceID)
	assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, resultNoDelegation.Resolution.Status)

	// Verify no delegation
	assert.Nil(t, resultNoDelegation.Delegation)

	// Cleanup
	DeleteThreatResolution(resolution1.ID)
	DeleteThreatResolution(resolution2.ID)
}
