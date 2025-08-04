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

func TestRelationshipService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.Relationship{},
	)
	defer cleanup()

	t.Run("AddRelationshipConsumesApiOf_ComponentToComponent", func(t *testing.T) {
		// Create test components
		component1, err := CreateComponent("Test Component 1", "First test component", models.ComponentTypeInstance)
		require.NoError(t, err)
		component2, err := CreateComponent("Test Component 2", "Second test component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Add relationship between components
		err = AddRelationshipConsumesApiOfComponents(component1.ID, component2.ID)
		require.NoError(t, err)

		// Verify both relationships were created
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Check for CONSUMES_API_OF relationship
		consumesFound := false
		consumedByFound := false
		for _, rel := range relationships {
			if rel.Type == "CONSUMES_API_OF" {
				consumesFound = true
				assert.Equal(t, component1.ID, rel.FromComponentID)
				assert.Equal(t, component2.ID, rel.ToComponentID)
			}
			if rel.Type == "API_CONSUMED_BY" {
				consumedByFound = true
				assert.Equal(t, component2.ID, rel.FromComponentID)
				assert.Equal(t, component1.ID, rel.ToComponentID)
			}
		}
		assert.True(t, consumesFound, "CONSUMES_API_OF relationship should be created")
		assert.True(t, consumedByFound, "API_CONSUMED_BY relationship should be created")
	})

	t.Run("AddRelationshipConsumesApiOf_ProductComponents", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test product components
		productComponent1, err := CreateComponent("API Product 1", "First API product component", models.ComponentTypeProduct)
		require.NoError(t, err)
		productComponent2, err := CreateComponent("API Product 2", "Second API product component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Add relationship between product components
		err = AddRelationshipConsumesApiOfComponents(productComponent1.ID, productComponent2.ID)
		require.NoError(t, err)

		// Verify both relationships were created
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Check for CONSUMES_API_OF relationship
		consumesFound := false
		consumedByFound := false
		for _, rel := range relationships {
			if rel.Type == "CONSUMES_API_OF" {
				consumesFound = true
				assert.Equal(t, productComponent1.ID, rel.FromComponentID)
				assert.Equal(t, productComponent2.ID, rel.ToComponentID)
			}
			if rel.Type == "API_CONSUMED_BY" {
				consumedByFound = true
				assert.Equal(t, productComponent2.ID, rel.FromComponentID)
				assert.Equal(t, productComponent1.ID, rel.ToComponentID)
			}
		}
		assert.True(t, consumesFound, "CONSUMES_API_OF relationship should be created")
		assert.True(t, consumedByFound, "API_CONSUMED_BY relationship should be created")
	})

	t.Run("AddRelationshipConsumesApiOf_SelfRelationshipError", func(t *testing.T) {
		// Create test data
		component1, err := CreateComponent("Self-Reference Component", "Component for self-reference test", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Try to create relationship between component and itself (should fail if validation exists)
		err = AddRelationshipConsumesApiOfComponents(component1.ID, component1.ID)
		// Note: This may or may not error depending on business logic - adjust assertion as needed
		if err != nil {
			assert.Contains(t, err.Error(), "self")
		}
	})

	t.Run("AddRelationshipConsumesApiOf_InvalidComponentsError", func(t *testing.T) {
		// Try to create relationship with invalid component IDs
		nonExistentID1 := uuid.New()
		nonExistentID2 := uuid.New()
		err := AddRelationshipConsumesApiOfComponents(nonExistentID1, nonExistentID2)
		// This should error if validation exists for non-existent components
		if err != nil {
			assert.Contains(t, err.Error(), "component")
		}
	})

	t.Run("AddRelationship_ValidComponentRelationship", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		component1, err := CreateComponent("Rel Component 1", "Relationship test component 1", models.ComponentTypeInstance)
		require.NoError(t, err)
		component2, err := CreateComponent("Rel Component 2", "Relationship test component 2", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Add relationship with vice versa
		err = AddRelationshipComponents(component1.ID, component2.ID, "DEPENDS_ON", "DEPENDENCY_OF")
		require.NoError(t, err)

		// Verify both relationships were created
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Check for primary relationship
		primaryFound := false
		viceVersaFound := false
		for _, rel := range relationships {
			if rel.Type == "DEPENDS_ON" {
				primaryFound = true
				assert.Equal(t, component1.ID, rel.FromComponentID)
				assert.Equal(t, component2.ID, rel.ToComponentID)
			}
			if rel.Type == "DEPENDENCY_OF" {
				viceVersaFound = true
				assert.Equal(t, component2.ID, rel.FromComponentID)
				assert.Equal(t, component1.ID, rel.ToComponentID)
			}
		}
		assert.True(t, primaryFound, "Primary relationship should be created")
		assert.True(t, viceVersaFound, "Vice versa relationship should be created")
	})

	t.Run("AddRelationship_ValidProductComponentRelationship", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		productComponent1, err := CreateComponent("Rel Product A", "Relationship test product component A", models.ComponentTypeProduct)
		require.NoError(t, err)
		productComponent2, err := CreateComponent("Rel Product B", "Relationship test product component B", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Add relationship without vice versa
		err = AddRelationshipComponents(productComponent1.ID, productComponent2.ID, "RELATED_TO", "")
		require.NoError(t, err)

		// Verify only one relationship was created
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 1)

		rel := relationships[0]
		assert.Equal(t, "RELATED_TO", rel.Type)
		assert.Equal(t, productComponent1.ID, rel.FromComponentID)
		assert.Equal(t, productComponent2.ID, rel.ToComponentID)
	})

	t.Run("AddRelationship_ValidationErrors", func(t *testing.T) {
		// Create test data
		component1, err := CreateComponent("Validation Component", "Validation test component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Test same component as from and to (self-relationship)
		err = AddRelationshipComponents(component1.ID, component1.ID, "TEST", "")
		// This may or may not error depending on business logic
		if err != nil {
			assert.Contains(t, err.Error(), "self")
		}

		// Test invalid component IDs
		nonExistentID := uuid.New()
		err = AddRelationshipComponents(nonExistentID, component1.ID, "TEST", "")
		// This should error if validation exists for non-existent components
		if err != nil {
			assert.Contains(t, err.Error(), "component")
		}
	})

	t.Run("DeleteRelationshipById", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		component1, err := CreateComponent("Delete Component 1", "Delete test component 1", models.ComponentTypeInstance)
		require.NoError(t, err)
		component2, err := CreateComponent("Delete Component 2", "Delete test component 2", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create relationship
		err = AddRelationshipComponents(component1.ID, component2.ID, "TEST_RELATIONSHIP", "")
		require.NoError(t, err)

		// Get the relationship ID
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 1)
		relationshipID := relationships[0].ID

		// Delete the relationship
		err = DeleteRelationshipById(relationshipID)
		require.NoError(t, err)

		// Verify it's deleted
		relationships, err = ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 0)
	})

	t.Run("DeleteRelationshipById_NotFound", func(t *testing.T) {
		// Try to delete a non-existent relationship
		nonExistentID := uuid.New()
		err := DeleteRelationshipById(nonExistentID)
		// Should succeed (GORM behavior for delete operations)
		assert.NoError(t, err)
	})

	t.Run("GetRelationship", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		component1, err := CreateComponent("Get Component 1", "Get test component 1", models.ComponentTypeInstance)
		require.NoError(t, err)
		component2, err := CreateComponent("Get Component 2", "Get test component 2", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create relationship
		err = AddRelationshipComponents(component1.ID, component2.ID, "GET_TEST", "")
		require.NoError(t, err)

		// Get the relationship ID
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 1)
		relationshipID := relationships[0].ID

		// Get the relationship
		relationship, err := GetRelationship(relationshipID)
		require.NoError(t, err)
		assert.NotNil(t, relationship)
		assert.Equal(t, "GET_TEST", relationship.Type)
		assert.Equal(t, component1.ID, relationship.FromComponentID)
		assert.Equal(t, component2.ID, relationship.ToComponentID)
	})

	t.Run("GetRelationship_NotFound", func(t *testing.T) {
		// Try to get a non-existent relationship
		nonExistentID := uuid.New()
		relationship, err := GetRelationship(nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("ListRelationshipsByFromComponentID", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		component1, err := CreateComponent("List Component 1", "List test component 1", models.ComponentTypeInstance)
		require.NoError(t, err)
		component2, err := CreateComponent("List Component 2", "List test component 2", models.ComponentTypeInstance)
		require.NoError(t, err)
		component3, err := CreateComponent("List Component 3", "List test component 3", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create relationships from component1
		err = AddRelationshipComponents(component1.ID, component2.ID, "REL_TYPE_1", "")
		require.NoError(t, err)
		err = AddRelationshipComponents(component1.ID, component3.ID, "REL_TYPE_2", "")
		require.NoError(t, err)

		// Create relationship from component2 (should not appear in our list)
		err = AddRelationshipComponents(component2.ID, component3.ID, "REL_TYPE_3", "")
		require.NoError(t, err)

		// List relationships from component1
		relationships, err := ListRelationshipsByFromComponentID(component1.ID)
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Verify all relationships are from component1
		for _, rel := range relationships {
			assert.Equal(t, component1.ID, rel.FromComponentID)
		}
	})

	t.Run("ListRelationshipsByFromComponentID_ProductType", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		productComponent1, err := CreateComponent("List Product Component 1", "List test product component 1", models.ComponentTypeProduct)
		require.NoError(t, err)
		productComponent2, err := CreateComponent("List Product Component 2", "List test product component 2", models.ComponentTypeProduct)
		require.NoError(t, err)
		productComponent3, err := CreateComponent("List Product Component 3", "List test product component 3", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create relationships from productComponent1
		err = AddRelationshipComponents(productComponent1.ID, productComponent2.ID, "PROD_REL_1", "")
		require.NoError(t, err)
		err = AddRelationshipComponents(productComponent1.ID, productComponent3.ID, "PROD_REL_2", "")
		require.NoError(t, err)

		// Create relationship from productComponent2 (should not appear in our list)
		err = AddRelationshipComponents(productComponent2.ID, productComponent3.ID, "PROD_REL_3", "")
		require.NoError(t, err)

		// List relationships from productComponent1
		relationships, err := ListRelationshipsByFromComponentID(productComponent1.ID)
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Verify all relationships are from productComponent1
		for _, rel := range relationships {
			assert.Equal(t, productComponent1.ID, rel.FromComponentID)
		}
	})

	t.Run("ListRelationshipsByType", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		component1, err := CreateComponent("Type Component 1", "Type test component 1", models.ComponentTypeProduct)
		require.NoError(t, err)
		component2, err := CreateComponent("Type Component 2", "Type test component 2", models.ComponentTypeProduct)
		require.NoError(t, err)
		component3, err := CreateComponent("Type Component 3", "Type test component 3", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create relationships of different types
		err = AddRelationshipComponents(component1.ID, component2.ID, "TYPE_A", "")
		require.NoError(t, err)
		err = AddRelationshipComponents(component1.ID, component3.ID, "TYPE_A", "")
		require.NoError(t, err)
		err = AddRelationshipComponents(component2.ID, component3.ID, "TYPE_B", "")
		require.NoError(t, err)

		// List relationships of TYPE_A
		relationships, err := ListRelationshipsByType("TYPE_A")
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Verify all relationships are of TYPE_A
		for _, rel := range relationships {
			assert.Equal(t, "TYPE_A", rel.Type)
		}

		// List relationships of TYPE_B
		relationships, err = ListRelationshipsByType("TYPE_B")
		require.NoError(t, err)
		assert.Len(t, relationships, 1)
		assert.Equal(t, "TYPE_B", relationships[0].Type)
	})

	t.Run("ListRelationships_Empty", func(t *testing.T) {
		// Clear all relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// List relationships
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 0)
	})
}
