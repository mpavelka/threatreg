package service

import (
	"testing"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateInheritsThreatsRelationship_Integration(t *testing.T) {

	t.Run("CreateValidRelationship", func(t *testing.T) {
		// Create test components
		childComponent, err := CreateComponent("Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		parentComponent, err := CreateComponent("Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create inherits threats relationship
		relationship, err := CreateInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		assert.NoError(t, err)
		assert.NotNil(t, relationship)

		// Verify the relationship properties
		assert.Equal(t, childComponent.ID, relationship.FromID)
		assert.Equal(t, parentComponent.ID, relationship.ToID)
		assert.Equal(t, string(models.ReservedLabelInheritsThreatsFrom), relationship.Label)
		assert.NotEqual(t, uuid.Nil, relationship.ID)
	})

	t.Run("CreateRelationshipWithNonExistentChild", func(t *testing.T) {
		// Create only parent component
		parentComponent, err := CreateComponent("Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		nonExistentChildID := uuid.New()

		// Try to create relationship with non-existent child
		relationship, err := CreateInheritsThreatsRelationship(nonExistentChildID, parentComponent.ID)
		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "from component does not exist")
	})

	t.Run("CreateRelationshipWithNonExistentParent", func(t *testing.T) {
		// Create only child component
		childComponent, err := CreateComponent("Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		nonExistentParentID := uuid.New()

		// Try to create relationship with non-existent parent
		relationship, err := CreateInheritsThreatsRelationship(childComponent.ID, nonExistentParentID)
		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "to component does not exist")
	})

	t.Run("CreateDuplicateRelationship", func(t *testing.T) {
		// Create test components
		childComponent, err := CreateComponent("Duplicate Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		parentComponent, err := CreateComponent("Duplicate Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create first relationship
		relationship1, err := CreateInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		require.NoError(t, err)
		require.NotNil(t, relationship1)

		// Try to create duplicate relationship
		relationship2, err := CreateInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		assert.Error(t, err)
		assert.Nil(t, relationship2)
		// Should fail due to unique constraint on (from_id, to_id, label)
	})

	t.Run("CreateSelfRelationship", func(t *testing.T) {
		// Create test component
		component, err := CreateComponent("Self Component", "Test component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Try to create relationship from component to itself
		relationship, err := CreateInheritsThreatsRelationship(component.ID, component.ID)
		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "component cannot have a relationship to itself")
	})
}

func TestRemoveInheritsThreatsRelationship_Integration(t *testing.T) {

	t.Run("RemoveExistingRelationship", func(t *testing.T) {
		// Create test components
		childComponent, err := CreateComponent("Remove Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		parentComponent, err := CreateComponent("Remove Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create relationship first
		relationship, err := CreateInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		require.NoError(t, err)
		require.NotNil(t, relationship)

		// Remove the relationship
		err = RemoveInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		assert.NoError(t, err)

		// Verify relationship is removed by trying to get it from component relationship service
		relationshipRepository, err := getComponentRelationshipRepository()
		require.NoError(t, err)
		
		deletedRelationship, err := relationshipRepository.GetByFromAndTo(nil, childComponent.ID, parentComponent.ID)
		assert.Error(t, err) // Should not be found
		assert.Nil(t, deletedRelationship)
	})

	t.Run("RemoveNonExistentRelationship", func(t *testing.T) {
		// Create test components
		childComponent, err := CreateComponent("Non-existent Remove Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		parentComponent, err := CreateComponent("Non-existent Remove Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Try to remove relationship that doesn't exist (should not error - GORM behavior)
		err = RemoveInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		assert.NoError(t, err) // GORM doesn't error when deleting non-existent records
	})

	t.Run("RemoveRelationshipWithNonExistentComponents", func(t *testing.T) {
		nonExistentChildID := uuid.New()
		nonExistentParentID := uuid.New()

		// Try to remove relationship with non-existent components (should not error)
		err := RemoveInheritsThreatsRelationship(nonExistentChildID, nonExistentParentID)
		assert.NoError(t, err) // GORM doesn't error when deleting non-existent records
	})
}

func TestCreateAndRemoveInheritsThreatsRelationship_Integration(t *testing.T) {

	t.Run("CreateThenRemoveRelationship", func(t *testing.T) {
		// Create test components
		childComponent, err := CreateComponent("Full Lifecycle Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		parentComponent, err := CreateComponent("Full Lifecycle Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create relationship
		relationship, err := CreateInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		require.NoError(t, err)
		require.NotNil(t, relationship)

		// Verify relationship exists by getting it through repository
		relationshipRepository, err := getComponentRelationshipRepository()
		require.NoError(t, err)
		
		existingRelationship, err := relationshipRepository.GetByFromAndTo(nil, childComponent.ID, parentComponent.ID)
		require.NoError(t, err)
		require.NotNil(t, existingRelationship)
		assert.Equal(t, string(models.ReservedLabelInheritsThreatsFrom), existingRelationship.Label)

		// Remove the relationship
		err = RemoveInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		require.NoError(t, err)

		// Verify relationship is gone
		deletedRelationship, err := relationshipRepository.GetByFromAndTo(nil, childComponent.ID, parentComponent.ID)
		assert.Error(t, err) // Should not be found
		assert.Nil(t, deletedRelationship)

		// Create the relationship again to ensure it can be recreated
		relationship2, err := CreateInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		assert.NoError(t, err)
		assert.NotNil(t, relationship2)
		assert.Equal(t, string(models.ReservedLabelInheritsThreatsFrom), relationship2.Label)
	})

	t.Run("VerifyReservedLabelIsUsed", func(t *testing.T) {
		// Create test components
		childComponent, err := CreateComponent("Label Verification Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		parentComponent, err := CreateComponent("Label Verification Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create relationship
		relationship, err := CreateInheritsThreatsRelationship(childComponent.ID, parentComponent.ID)
		require.NoError(t, err)
		require.NotNil(t, relationship)

		// Verify the label is exactly the reserved label
		assert.Equal(t, "__inherits_threats_from", relationship.Label)
		assert.Equal(t, string(models.ReservedLabelInheritsThreatsFrom), relationship.Label)

		// Verify the relationship direction (child -> parent)
		assert.Equal(t, childComponent.ID, relationship.FromID)
		assert.Equal(t, parentComponent.ID, relationship.ToID)
	})
}