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
	productComponent, err := CreateComponent("Test Product Component", "A test product component", models.ComponentTypeProduct)
	require.NoError(t, err)

	instanceComponent, err := CreateComponent("Test Instance Component", "A test instance component", models.ComponentTypeInstance)
	require.NoError(t, err)

	threat, err := CreateThreat("Test Threat", "A test threat")
	require.NoError(t, err)

	// Create threat assignments
	productAssignment, err := AssignThreatToComponent(productComponent.ID, threat.ID)
	require.NoError(t, err)

	instanceAssignment, err := AssignThreatToComponent(instanceComponent.ID, threat.ID)
	require.NoError(t, err)

	t.Run("CreateThreatResolution_Product", func(t *testing.T) {
		resolution, err := CreateThreatResolution(
			productAssignment.ID,
			productComponent.ID, // componentID
			models.ThreatAssignmentResolutionStatusResolved,
			"Product threat resolved",
		)

		require.NoError(t, err)
		assert.NotNil(t, resolution)
		assert.NotEqual(t, uuid.Nil, resolution.ID)
		assert.Equal(t, productAssignment.ID, resolution.ThreatAssignmentID)
		assert.Equal(t, productComponent.ID, resolution.ComponentID)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, resolution.Status)
		assert.Equal(t, "Product threat resolved", resolution.Description)

		// Cleanup
		DeleteThreatResolution(resolution.ID)
	})

	t.Run("CreateThreatResolution_Instance", func(t *testing.T) {
		resolution, err := CreateThreatResolution(
			instanceAssignment.ID,
			instanceComponent.ID, // componentID
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Instance threat awaiting resolution",
		)

		require.NoError(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, instanceAssignment.ID, resolution.ThreatAssignmentID)
		assert.Equal(t, instanceComponent.ID, resolution.ComponentID)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusAwaiting, resolution.Status)

		// Cleanup
		DeleteThreatResolution(resolution.ID)
	})

	t.Run("CreateThreatResolution_NilComponentID", func(t *testing.T) {
		resolution, err := CreateThreatResolution(
			productAssignment.ID,
			uuid.Nil, // nil componentID
			models.ThreatAssignmentResolutionStatusResolved,
			"Invalid resolution",
		)

		assert.Error(t, err)
		assert.Nil(t, resolution)
		assert.Contains(t, err.Error(), "ComponentID")
	})

	t.Run("GetThreatResolution", func(t *testing.T) {
		// Create resolution
		created, err := CreateThreatResolution(
			productAssignment.ID,
			productComponent.ID,
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
			productComponent.ID,
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
			productComponent.ID,
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
			productComponent.ID,
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

	t.Run("ListThreatResolutionsByComponentID", func(t *testing.T) {
		// Create multiple resolutions for the product component
		resolution1, err := CreateThreatResolution(
			productAssignment.ID,
			productComponent.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Resolution 1",
		)
		require.NoError(t, err)

		// Create another threat assignment for the same product component
		threat2, err := CreateThreat("Another Test Threat", "Another test threat")
		require.NoError(t, err)
		assignment2, err := AssignThreatToComponent(productComponent.ID, threat2.ID)
		require.NoError(t, err)

		resolution2, err := CreateThreatResolution(
			assignment2.ID,
			productComponent.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution 2",
		)
		require.NoError(t, err)

		// List resolutions for component
		resolutions, err := ListThreatResolutionsByComponentID(productComponent.ID)
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

	t.Run("ListThreatResolutionsByComponentID_Instance", func(t *testing.T) {
		// Create resolution for instance component
		resolution, err := CreateThreatResolution(
			instanceAssignment.ID,
			instanceComponent.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Instance component resolution",
		)
		require.NoError(t, err)

		// List resolutions for instance component
		resolutions, err := ListThreatResolutionsByComponentID(instanceComponent.ID)
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
		assignment3, err := AssignThreatToComponent(productComponent.ID, threat3.ID)
		require.NoError(t, err)

		// Create resolution
		created, err := CreateThreatResolution(
			assignment3.ID,
			productComponent.ID,
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
		component1, _ := CreateComponent("Component 1", "First test component", models.ComponentTypeInstance)
		component2, _ := CreateComponent("Component 2", "Second test component", models.ComponentTypeInstance)
		threat, _ := CreateThreat("Test Threat", "A test threat")
		assignment1, _ := AssignThreatToComponent(component1.ID, threat.ID)
		assignment2, _ := AssignThreatToComponent(component2.ID, threat.ID)

		source, _ := CreateThreatResolution(assignment1.ID, component1.ID, models.ThreatAssignmentResolutionStatusAwaiting, "Source")
		target, _ := CreateThreatResolution(assignment2.ID, component2.ID, models.ThreatAssignmentResolutionStatusResolved, "Target")
		defer func() { DeleteThreatResolution(source.ID); DeleteThreatResolution(target.ID) }()

		err := DelegateResolution(*source, *target)
		require.NoError(t, err)

		updatedSource, _ := GetThreatResolution(source.ID)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, updatedSource.Status)
	})

	t.Run("DelegateResolution_MultipleTimes", func(t *testing.T) {
		component1, _ := CreateComponent("Component 1", "First test component", models.ComponentTypeInstance)
		component2, _ := CreateComponent("Component 2", "Second test component", models.ComponentTypeInstance)
		component3, _ := CreateComponent("Component 3", "Third test component", models.ComponentTypeInstance)
		threat, _ := CreateThreat("Test Threat 2", "A test threat")
		assignment1, _ := AssignThreatToComponent(component1.ID, threat.ID)
		assignment2, _ := AssignThreatToComponent(component2.ID, threat.ID)
		assignment3, _ := AssignThreatToComponent(component3.ID, threat.ID)

		source, err := CreateThreatResolution(assignment1.ID, component1.ID, models.ThreatAssignmentResolutionStatusAwaiting, "Source")
		require.NoError(t, err)
		target1, err := CreateThreatResolution(assignment2.ID, component2.ID, models.ThreatAssignmentResolutionStatusAccepted, "Target1")
		require.NoError(t, err)
		target2, err := CreateThreatResolution(assignment3.ID, component3.ID, models.ThreatAssignmentResolutionStatusResolved, "Target2")
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
		component1, _ := CreateComponent("Component 1", "First test component", models.ComponentTypeInstance)
		component2, _ := CreateComponent("Component 2", "Second test component", models.ComponentTypeInstance)
		threat, _ := CreateThreat("Test Threat 3", "A test threat")
		assignment1, _ := AssignThreatToComponent(component1.ID, threat.ID)
		assignment2, _ := AssignThreatToComponent(component2.ID, threat.ID)

		source, err := CreateThreatResolution(assignment1.ID, component1.ID, models.ThreatAssignmentResolutionStatusAwaiting, "Source")
		require.NoError(t, err)
		target, err := CreateThreatResolution(assignment2.ID, component2.ID, models.ThreatAssignmentResolutionStatusResolved, "Target")
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
		component1, _ := CreateComponent("Component 1", "First test component", models.ComponentTypeInstance)
		component2, _ := CreateComponent("Component 2", "Second test component", models.ComponentTypeInstance)
		component3, _ := CreateComponent("Component 3", "Third test component", models.ComponentTypeInstance)
		threat, _ := CreateThreat("Test Threat 4", "A test threat")
		assignment1, _ := AssignThreatToComponent(component1.ID, threat.ID)
		assignment2, _ := AssignThreatToComponent(component2.ID, threat.ID)
		assignment3, _ := AssignThreatToComponent(component3.ID, threat.ID)

		// Create chain: resA -> resB -> resC
		resA, err := CreateThreatResolution(assignment1.ID, component1.ID, models.ThreatAssignmentResolutionStatusAwaiting, "ResA")
		require.NoError(t, err)
		resB, err := CreateThreatResolution(assignment2.ID, component2.ID, models.ThreatAssignmentResolutionStatusAwaiting, "ResB")
		require.NoError(t, err)
		resC, err := CreateThreatResolution(assignment3.ID, component3.ID, models.ThreatAssignmentResolutionStatusAwaiting, "ResC")
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

func TestGetComponentLevelThreatResolutionWithDelegation(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create test data
	component1, _ := CreateComponent("Component 1", "First test component", models.ComponentTypeInstance)
	component2, _ := CreateComponent("Component 2", "Second test component", models.ComponentTypeInstance)
	threat, _ := CreateThreat("Test Threat", "A test threat")
	assignment1, _ := AssignThreatToComponent(component1.ID, threat.ID)
	assignment2, _ := AssignThreatToComponent(component2.ID, threat.ID)

	// Create two resolutions
	resolution1, err := CreateThreatResolution(assignment1.ID, component1.ID, models.ThreatAssignmentResolutionStatusAwaiting, "Source resolution")
	require.NoError(t, err)
	resolution2, err := CreateThreatResolution(assignment2.ID, component2.ID, models.ThreatAssignmentResolutionStatusResolved, "Target resolution")
	require.NoError(t, err)

	// Create delegation from resolution1 to resolution2
	err = DelegateResolution(*resolution1, *resolution2)
	require.NoError(t, err)

	// Test GetComponentLevelThreatResolutionWithDelegation for resolution1 (has delegation)
	resultWithDelegation, err := GetComponentLevelThreatResolutionWithDelegation(assignment1.ID, component1.ID)
	require.NoError(t, err)
	require.NotNil(t, resultWithDelegation)

	// Verify resolution data
	assert.Equal(t, resolution1.ID, resultWithDelegation.Resolution.ID)
	assert.Equal(t, assignment1.ID, resultWithDelegation.Resolution.ThreatAssignmentID)
	assert.Equal(t, component1.ID, resultWithDelegation.Resolution.ComponentID)
	assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, resultWithDelegation.Resolution.Status)

	// Verify delegation data
	require.NotNil(t, resultWithDelegation.Delegation)
	assert.Equal(t, resolution1.ID, resultWithDelegation.Delegation.DelegatedBy)
	assert.Equal(t, resolution2.ID, resultWithDelegation.Delegation.DelegatedTo)

	// Test GetComponentLevelThreatResolutionWithDelegation for resolution2 (no delegation)
	resultNoDelegation, err := GetComponentLevelThreatResolutionWithDelegation(assignment2.ID, component2.ID)
	require.NoError(t, err)
	require.NotNil(t, resultNoDelegation)

	// Verify resolution data
	assert.Equal(t, resolution2.ID, resultNoDelegation.Resolution.ID)
	assert.Equal(t, assignment2.ID, resultNoDelegation.Resolution.ThreatAssignmentID)
	assert.Equal(t, component2.ID, resultNoDelegation.Resolution.ComponentID)
	assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, resultNoDelegation.Resolution.Status)

	// Verify no delegation
	assert.Nil(t, resultNoDelegation.Delegation)

	// Cleanup
	DeleteThreatResolution(resolution1.ID)
	DeleteThreatResolution(resolution2.ID)
}
