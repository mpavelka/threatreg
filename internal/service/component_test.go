package service

import (
	"testing"
	"threatreg/internal/database"
	"threatreg/internal/models"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestComponentService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("CreateComponent_Product", func(t *testing.T) {
		// Test data
		name := "Test Product Component"
		description := "A test product component description"
		componentType := models.ComponentTypeProduct

		// Create component
		component, err := CreateComponent(name, description, componentType)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, component)
		assert.NotEqual(t, uuid.Nil, component.ID)
		assert.Equal(t, name, component.Name)
		assert.Equal(t, description, component.Description)
		assert.Equal(t, componentType, component.Type)

		// Verify component was actually saved to database
		db := database.GetDB()
		var dbComponent models.Component
		err = db.First(&dbComponent, "id = ?", component.ID).Error
		require.NoError(t, err)
		assert.Equal(t, component.ID, dbComponent.ID)
		assert.Equal(t, name, dbComponent.Name)
		assert.Equal(t, description, dbComponent.Description)
		assert.Equal(t, componentType, dbComponent.Type)
	})

	t.Run("CreateComponent_Instance", func(t *testing.T) {
		// Test data
		name := "Test Instance Component"
		description := "A test instance component description"
		componentType := models.ComponentTypeInstance

		// Create component
		component, err := CreateComponent(name, description, componentType)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, component)
		assert.NotEqual(t, uuid.Nil, component.ID)
		assert.Equal(t, name, component.Name)
		assert.Equal(t, description, component.Description)
		assert.Equal(t, componentType, component.Type)

		// Verify component was actually saved to database
		db := database.GetDB()
		var dbComponent models.Component
		err = db.First(&dbComponent, "id = ?", component.ID).Error
		require.NoError(t, err)
		assert.Equal(t, component.ID, dbComponent.ID)
		assert.Equal(t, name, dbComponent.Name)
		assert.Equal(t, description, dbComponent.Description)
		assert.Equal(t, componentType, dbComponent.Type)
	})

	t.Run("GetComponent", func(t *testing.T) {
		// Create a component first
		name := "Get Test Component"
		description := "Component for get test"
		componentType := models.ComponentTypeProduct
		createdComponent, err := CreateComponent(name, description, componentType)
		require.NoError(t, err)

		// Get the component
		retrievedComponent, err := GetComponent(createdComponent.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedComponent)
		assert.Equal(t, createdComponent.ID, retrievedComponent.ID)
		assert.Equal(t, name, retrievedComponent.Name)
		assert.Equal(t, description, retrievedComponent.Description)
		assert.Equal(t, componentType, retrievedComponent.Type)
	})

	t.Run("GetComponent_NotFound", func(t *testing.T) {
		// Try to get a non-existent component
		nonExistentID := uuid.New()
		component, err := GetComponent(nonExistentID)

		// Should return error and nil component
		assert.Error(t, err)
		assert.Nil(t, component)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateComponent", func(t *testing.T) {
		// Create a component first
		originalName := "Original Component"
		originalDescription := "Original description"
		componentType := models.ComponentTypeInstance
		createdComponent, err := CreateComponent(originalName, originalDescription, componentType)
		require.NoError(t, err)

		// Update the component
		newName := "Updated Component"
		newDescription := "Updated description"
		updatedComponent, err := UpdateComponent(createdComponent.ID, &newName, &newDescription)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedComponent)
		assert.Equal(t, createdComponent.ID, updatedComponent.ID)
		assert.Equal(t, newName, updatedComponent.Name)
		assert.Equal(t, newDescription, updatedComponent.Description)
		assert.Equal(t, componentType, updatedComponent.Type) // Type should remain unchanged

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbComponent models.Component
		err = db.First(&dbComponent, "id = ?", createdComponent.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbComponent.Name)
		assert.Equal(t, newDescription, dbComponent.Description)
		assert.Equal(t, componentType, dbComponent.Type)
	})

	t.Run("UpdateComponent_PartialUpdate", func(t *testing.T) {
		// Create a component first
		originalName := "Partial Update Component"
		originalDescription := "Original description"
		componentType := models.ComponentTypeProduct
		createdComponent, err := CreateComponent(originalName, originalDescription, componentType)
		require.NoError(t, err)

		// Update only the name
		newName := "New Name Only"
		updatedComponent, err := UpdateComponent(createdComponent.ID, &newName, nil)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedComponent)
		assert.Equal(t, createdComponent.ID, updatedComponent.ID)
		assert.Equal(t, newName, updatedComponent.Name)
		assert.Equal(t, originalDescription, updatedComponent.Description) // Should remain unchanged
		assert.Equal(t, componentType, updatedComponent.Type)

		// Verify the partial update was persisted
		db := database.GetDB()
		var dbComponent models.Component
		err = db.First(&dbComponent, "id = ?", createdComponent.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbComponent.Name)
		assert.Equal(t, originalDescription, dbComponent.Description)
		assert.Equal(t, componentType, dbComponent.Type)
	})

	t.Run("UpdateComponent_NotFound", func(t *testing.T) {
		// Try to update a non-existent component
		nonExistentID := uuid.New()
		newName := "New Name"
		component, err := UpdateComponent(nonExistentID, &newName, nil)

		// Should return error and nil component
		assert.Error(t, err)
		assert.Nil(t, component)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteComponent", func(t *testing.T) {
		// Create a component first
		name := "Delete Test Component"
		description := "Component to be deleted"
		componentType := models.ComponentTypeInstance
		createdComponent, err := CreateComponent(name, description, componentType)
		require.NoError(t, err)

		// Delete the component
		err = DeleteComponent(createdComponent.ID)

		// Assertions
		require.NoError(t, err)

		// Verify the component was actually deleted from database
		db := database.GetDB()
		var dbComponent models.Component
		err = db.First(&dbComponent, "id = ?", createdComponent.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteComponent_NotFound", func(t *testing.T) {
		// Try to delete a non-existent component
		nonExistentID := uuid.New()
		err := DeleteComponent(nonExistentID)

		// Delete should succeed even if component doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})

	t.Run("ListComponents", func(t *testing.T) {
		// Clear any existing components first
		db := database.GetDB()
		db.Exec("DELETE FROM components")

		// Create multiple components
		components := []struct {
			name        string
			description string
			componentType models.ComponentType
		}{
			{"Component 1", "Description 1", models.ComponentTypeProduct},
			{"Component 2", "Description 2", models.ComponentTypeInstance},
			{"Component 3", "Description 3", models.ComponentTypeProduct},
		}

		var createdComponents []*models.Component
		for _, c := range components {
			component, err := CreateComponent(c.name, c.description, c.componentType)
			require.NoError(t, err)
			createdComponents = append(createdComponents, component)
		}

		// List all components
		retrievedComponents, err := ListComponents()

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedComponents, len(components))

		// Verify all created components are in the list
		componentMap := make(map[uuid.UUID]models.Component)
		for _, c := range retrievedComponents {
			componentMap[c.ID] = c
		}

		for _, created := range createdComponents {
			retrieved, exists := componentMap[created.ID]
			assert.True(t, exists, "Created component should exist in list")
			assert.Equal(t, created.Name, retrieved.Name)
			assert.Equal(t, created.Description, retrieved.Description)
			assert.Equal(t, created.Type, retrieved.Type)
		}
	})

	t.Run("ListComponents_Empty", func(t *testing.T) {
		// Clear all components
		db := database.GetDB()
		db.Exec("DELETE FROM components")

		// List components
		components, err := ListComponents()

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, components, 0)
	})

	t.Run("ListComponentsByType", func(t *testing.T) {
		// Clear any existing components first
		db := database.GetDB()
		db.Exec("DELETE FROM components")

		// Create components of different types
		productComponent1, err := CreateComponent("Product 1", "Description 1", models.ComponentTypeProduct)
		require.NoError(t, err)
		productComponent2, err := CreateComponent("Product 2", "Description 2", models.ComponentTypeProduct)
		require.NoError(t, err)
		instanceComponent1, err := CreateComponent("Instance 1", "Description 3", models.ComponentTypeInstance)
		require.NoError(t, err)

		// List product components
		productComponents, err := ListComponentsByType(models.ComponentTypeProduct)
		require.NoError(t, err)
		assert.Len(t, productComponents, 2)

		productIDs := []uuid.UUID{productComponents[0].ID, productComponents[1].ID}
		assert.Contains(t, productIDs, productComponent1.ID)
		assert.Contains(t, productIDs, productComponent2.ID)

		// List instance components
		instanceComponents, err := ListComponentsByType(models.ComponentTypeInstance)
		require.NoError(t, err)
		assert.Len(t, instanceComponents, 1)
		assert.Equal(t, instanceComponent1.ID, instanceComponents[0].ID)
	})

	t.Run("FilterComponents", func(t *testing.T) {
		// Clear any existing components first
		db := database.GetDB()
		db.Exec("DELETE FROM components")

		// Create components with different names
		component1, err := CreateComponent("Web Application", "A web app", models.ComponentTypeProduct)
		require.NoError(t, err)
		component2, err := CreateComponent("Mobile App", "A mobile app", models.ComponentTypeInstance)
		require.NoError(t, err)
		component3, err := CreateComponent("Database Server", "A database", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Filter by "App" - should match components with "App" in name
		filteredComponents, err := FilterComponents("App")
		require.NoError(t, err)
		assert.Len(t, filteredComponents, 2)

		componentIDs := []uuid.UUID{filteredComponents[0].ID, filteredComponents[1].ID}
		assert.Contains(t, componentIDs, component1.ID)
		assert.Contains(t, componentIDs, component2.ID)

		// Filter by "Database" - should match one component
		databaseComponents, err := FilterComponents("Database")
		require.NoError(t, err)
		assert.Len(t, databaseComponents, 1)
		assert.Equal(t, component3.ID, databaseComponents[0].ID)

		// Filter by non-existent term
		noMatchComponents, err := FilterComponents("NonExistent")
		require.NoError(t, err)
		assert.Len(t, noMatchComponents, 0)
	})

	t.Run("AssignThreatToComponent", func(t *testing.T) {
		// Create a component first
		component, err := CreateComponent("Test Component for Threat", "A test component for threat assignment", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create a threat first
		threat, err := CreateThreat("Test Threat", "A test threat for assignment")
		require.NoError(t, err)

		// Assign threat to component
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.NotEqual(t, 0, assignment.ID)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, component.ID, assignment.ComponentID)

		// Verify assignment was saved to database
		db := database.GetDB()
		var dbAssignment models.ThreatAssignment
		err = db.First(&dbAssignment, "id = ?", assignment.ID).Error
		require.NoError(t, err)
		assert.Equal(t, assignment.ThreatID, dbAssignment.ThreatID)
		assert.Equal(t, assignment.ComponentID, dbAssignment.ComponentID)
	})

	t.Run("AssignThreatToComponent_Duplicate", func(t *testing.T) {
		// Create a component and threat
		component, err := CreateComponent("Duplicate Test Component", "A test component for duplicate assignment", models.ComponentTypeInstance)
		require.NoError(t, err)

		threat, err := CreateThreat("Duplicate Test Threat", "A test threat for duplicate assignment")
		require.NoError(t, err)

		// Assign threat to component first time
		assignment1, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)
		require.NotNil(t, assignment1)

		// Try to assign the same threat to the same component again
		assignment2, err := AssignThreatToComponent(component.ID, threat.ID)

		// Should return the existing assignment, not create a new one
		require.NoError(t, err)
		assert.NotNil(t, assignment2)
		assert.Equal(t, assignment1.ID, assignment2.ID)
		assert.Equal(t, assignment1.ThreatID, assignment2.ThreatID)
		assert.Equal(t, assignment1.ComponentID, assignment2.ComponentID)

		// Verify that only one assignment exists in the database for this threat/component combination
		db := database.GetDB()
		var count int64
		err = db.Model(&models.ThreatAssignment{}).
			Where("threat_id = ? AND component_id = ?", threat.ID, component.ID).
			Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Should only have one assignment for this threat/component combination")
	})

	t.Run("AssignThreatToComponent_InvalidThreatID", func(t *testing.T) {
		// Create a component
		component, err := CreateComponent("Invalid Threat Test Component", "A test component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Try to assign non-existent threat
		nonExistentThreatID := uuid.New()
		assignment, err := AssignThreatToComponent(component.ID, nonExistentThreatID)

		// Should succeed (foreign key constraint allows it, but relationship won't load)
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, nonExistentThreatID, assignment.ThreatID)
		assert.Equal(t, component.ID, assignment.ComponentID)
	})

	t.Run("AssignThreatToComponent_InvalidComponentID", func(t *testing.T) {
		// Create a threat
		threat, err := CreateThreat("Invalid Component Test Threat", "A test threat")
		require.NoError(t, err)

		// Try to assign to non-existent component
		nonExistentComponentID := uuid.New()
		assignment, err := AssignThreatToComponent(nonExistentComponentID, threat.ID)

		// Should succeed (foreign key constraint allows it, but relationship won't load)
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, nonExistentComponentID, assignment.ComponentID)
	})

	t.Run("ListThreatAssignmentsByComponentID", func(t *testing.T) {
		// Create a component for testing
		component, err := CreateComponent("Test Component for Assignments", "A test component for threat assignment listing", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create multiple threats for testing
		threat1, err := CreateThreat("Component Threat 1", "First component test threat")
		require.NoError(t, err)
		threat2, err := CreateThreat("Component Threat 2", "Second component test threat")
		require.NoError(t, err)
		threat3, err := CreateThreat("Component Threat 3", "Third component test threat")
		require.NoError(t, err)

		// Assign threats to the component
		assignment1, err := AssignThreatToComponent(component.ID, threat1.ID)
		require.NoError(t, err)
		assignment2, err := AssignThreatToComponent(component.ID, threat2.ID)
		require.NoError(t, err)

		// Assign one threat to a different component to ensure filtering works
		otherComponent, err := CreateComponent("Other Component", "Another component for testing", models.ComponentTypeInstance)
		require.NoError(t, err)
		_, err = AssignThreatToComponent(otherComponent.ID, threat3.ID)
		require.NoError(t, err)

		// List threat assignments for our test component
		assignments, err := ListThreatAssignmentsByComponentID(component.ID)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, assignments, 2)

		// Check that we got the correct assignments
		assignmentIDs := []int{assignments[0].ID, assignments[1].ID}
		assert.Contains(t, assignmentIDs, assignment1.ID)
		assert.Contains(t, assignmentIDs, assignment2.ID)

		// Verify threat relationships are loaded
		for _, assignment := range assignments {
			assert.NotEmpty(t, assignment.Threat.Title)
			assert.Equal(t, component.ID, assignment.ComponentID)
		}
	})

	t.Run("ListThreatAssignmentsByComponentID_Empty", func(t *testing.T) {
		// Create a component with no threat assignments
		component, err := CreateComponent("Empty Assignments Component", "A component with no threat assignments", models.ComponentTypeInstance)
		require.NoError(t, err)

		// List threat assignments for this component
		assignments, err := ListThreatAssignmentsByComponentID(component.ID)

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, assignments, 0)
	})

	t.Run("ListThreatAssignmentsByComponentID_InvalidComponentID", func(t *testing.T) {
		// Try to list assignments for non-existent component
		nonExistentComponentID := uuid.New()
		assignments, err := ListThreatAssignmentsByComponentID(nonExistentComponentID)

		// Should succeed but return empty slice
		require.NoError(t, err)
		assert.Len(t, assignments, 0)
	})
}

func TestListThreatAssignmentsByComponentIDWithResolutionByComponentID(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("Happy Flow", func(t *testing.T) {
		// Create component1 and component2
		component1, err := CreateComponent("Component 1", "First test component", models.ComponentTypeProduct)
		require.NoError(t, err)
		component2, err := CreateComponent("Component 2", "Second test component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat1 and threat2
		threat1, err := CreateThreat("Threat 1", "First test threat")
		require.NoError(t, err)
		threat2, err := CreateThreat("Threat 2", "Second test threat")
		require.NoError(t, err)

		// Assign threat1 to component1
		assignment1, err := AssignThreatToComponent(component1.ID, threat1.ID)
		require.NoError(t, err)

		// Assign threat2 to component2
		assignment2, err := AssignThreatToComponent(component2.ID, threat2.ID)
		require.NoError(t, err)

		// Create threat resolutions for both components
		resolution1, err := CreateThreatResolution(
			assignment1.ID,
			component1.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution for component1",
		)
		require.NoError(t, err)

		resolution2, err := CreateThreatResolution(
			assignment2.ID,
			component2.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution for component2",
		)
		require.NoError(t, err)

		// Delegate resolution1 to resolution2
		err = DelegateResolution(*resolution1, *resolution2)
		require.NoError(t, err)

		// Mark resolution2 as resolved
		resolvedStatus := models.ThreatAssignmentResolutionStatusResolved
		_, err = UpdateThreatResolution(resolution2.ID, &resolvedStatus, nil)
		require.NoError(t, err)

		// Test ListThreatAssignmentsByComponentIDWithResolutionByComponentID for component1 filtered by component1
		results, err := ListThreatAssignmentsByComponentIDWithResolutionByComponentID(component1.ID, component1.ID)
		require.NoError(t, err)

		// Assertions
		assert.Len(t, results, 1, "Should return one threat assignment for component1")

		result := results[0]
		assert.Equal(t, assignment1.ID, result.ID)
		assert.Equal(t, threat1.ID, result.ThreatID)
		assert.Equal(t, component1.ID, result.ComponentID)

		// Verify threat relationship is loaded
		assert.Equal(t, threat1.ID, result.Threat.ID)
		assert.Equal(t, "Threat 1", result.Threat.Title)

		// Verify component relationship is loaded
		assert.Equal(t, component1.ID, result.Component.ID)
		assert.Equal(t, "Component 1", result.Component.Name)

		// Verify resolution status - should show resolved because resolution1 was delegated to resolution2 which is resolved
		assert.NotNil(t, result.ResolutionStatus)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, *result.ResolutionStatus)

		// Verify delegation status
		assert.True(t, result.IsDelegated, "Should show as delegated")

		// Test with different component filter - should return same assignment but no resolution info
		resultsOtherComponent, err := ListThreatAssignmentsByComponentIDWithResolutionByComponentID(component1.ID, component2.ID)
		require.NoError(t, err)

		assert.Len(t, resultsOtherComponent, 1, "Should still return the threat assignment for component1")
		otherResult := resultsOtherComponent[0]
		assert.Equal(t, assignment1.ID, otherResult.ID)

		// But resolution info should be nil since we filtered by component2 but resolution1 is for component1
		assert.Nil(t, otherResult.ResolutionStatus, "Should not have resolution status for different component")
		assert.False(t, otherResult.IsDelegated, "Should not show as delegated for different component")

		// Test with component2 - should return assignments for component2
		resultsComponent2, err := ListThreatAssignmentsByComponentIDWithResolutionByComponentID(component2.ID, component2.ID)
		require.NoError(t, err)
		assert.Len(t, resultsComponent2, 1, "Should return assignments for component2")
		
		result2 := resultsComponent2[0]
		assert.Equal(t, assignment2.ID, result2.ID)
		assert.Equal(t, threat2.ID, result2.ThreatID)
		assert.Equal(t, component2.ID, result2.ComponentID)
		
		// Should show resolved status
		assert.NotNil(t, result2.ResolutionStatus)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, *result2.ResolutionStatus)
		assert.False(t, result2.IsDelegated, "Component2 resolution should not show as delegated")
	})
}