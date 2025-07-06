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
		&models.Product{},
		&models.Instance{},
		&models.Relationship{},
	)
	defer cleanup()

	t.Run("AddRelationshipConsumesApiOf_InstanceToInstance", func(t *testing.T) {
		// Create test products and instances
		product1, err := CreateProduct("Test Product 1", "First test product")
		require.NoError(t, err)
		product2, err := CreateProduct("Test Product 2", "Second test product")
		require.NoError(t, err)

		instance1, err := CreateInstance("Test Instance 1", product1.ID)
		require.NoError(t, err)
		instance2, err := CreateInstance("Test Instance 2", product2.ID)
		require.NoError(t, err)

		// Add relationship between instances
		err = AddRelationshipConsumesApiOf(&instance1.ID, &instance2.ID, nil, nil)
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
				assert.Equal(t, instance1.ID, *rel.FromInstanceID)
				assert.Equal(t, instance2.ID, *rel.ToInstanceID)
				assert.Nil(t, rel.FromProductID)
				assert.Nil(t, rel.ToProductID)
			}
			if rel.Type == "API_CONSUMED_BY" {
				consumedByFound = true
				assert.Equal(t, instance2.ID, *rel.FromInstanceID)
				assert.Equal(t, instance1.ID, *rel.ToInstanceID)
				assert.Nil(t, rel.FromProductID)
				assert.Nil(t, rel.ToProductID)
			}
		}
		assert.True(t, consumesFound, "CONSUMES_API_OF relationship should be created")
		assert.True(t, consumedByFound, "API_CONSUMED_BY relationship should be created")
	})

	t.Run("AddRelationshipConsumesApiOf_ProductToProduct", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test products
		product1, err := CreateProduct("API Product 1", "First API product")
		require.NoError(t, err)
		product2, err := CreateProduct("API Product 2", "Second API product")
		require.NoError(t, err)

		// Add relationship between products
		err = AddRelationshipConsumesApiOf(nil, nil, &product1.ID, &product2.ID)
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
				assert.Equal(t, product1.ID, *rel.FromProductID)
				assert.Equal(t, product2.ID, *rel.ToProductID)
				assert.Nil(t, rel.FromInstanceID)
				assert.Nil(t, rel.ToInstanceID)
			}
			if rel.Type == "API_CONSUMED_BY" {
				consumedByFound = true
				assert.Equal(t, product2.ID, *rel.FromProductID)
				assert.Equal(t, product1.ID, *rel.ToProductID)
				assert.Nil(t, rel.FromInstanceID)
				assert.Nil(t, rel.ToInstanceID)
			}
		}
		assert.True(t, consumesFound, "CONSUMES_API_OF relationship should be created")
		assert.True(t, consumedByFound, "API_CONSUMED_BY relationship should be created")
	})

	t.Run("AddRelationshipConsumesApiOf_MixedTypesError", func(t *testing.T) {
		// Create test data
		product1, err := CreateProduct("Mixed Product", "Mixed test product")
		require.NoError(t, err)
		instance1, err := CreateInstance("Mixed Instance", product1.ID)
		require.NoError(t, err)

		// Try to create relationship between instance and product (should fail)
		err = AddRelationshipConsumesApiOf(&instance1.ID, nil, nil, &product1.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid relationship: must be between two instances or two products")
	})

	t.Run("AddRelationshipConsumesApiOf_InvalidRelationshipError", func(t *testing.T) {
		// Try to create relationship with invalid parameters
		err := AddRelationshipConsumesApiOf(nil, nil, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid relationship: must be between two instances or two products")
	})

	t.Run("AddRelationship_ValidInstanceRelationship", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		product1, err := CreateProduct("Rel Product 1", "Relationship test product 1")
		require.NoError(t, err)
		product2, err := CreateProduct("Rel Product 2", "Relationship test product 2")
		require.NoError(t, err)

		instance1, err := CreateInstance("Rel Instance 1", product1.ID)
		require.NoError(t, err)
		instance2, err := CreateInstance("Rel Instance 2", product2.ID)
		require.NoError(t, err)

		// Add relationship with vice versa
		err = AddRelationship(&instance1.ID, nil, &instance2.ID, nil, "DEPENDS_ON", "DEPENDENCY_OF")
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
				assert.Equal(t, instance1.ID, *rel.FromInstanceID)
				assert.Equal(t, instance2.ID, *rel.ToInstanceID)
			}
			if rel.Type == "DEPENDENCY_OF" {
				viceVersaFound = true
				assert.Equal(t, instance2.ID, *rel.FromInstanceID)
				assert.Equal(t, instance1.ID, *rel.ToInstanceID)
			}
		}
		assert.True(t, primaryFound, "Primary relationship should be created")
		assert.True(t, viceVersaFound, "Vice versa relationship should be created")
	})

	t.Run("AddRelationship_ValidProductRelationship", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		product1, err := CreateProduct("Rel Product A", "Relationship test product A")
		require.NoError(t, err)
		product2, err := CreateProduct("Rel Product B", "Relationship test product B")
		require.NoError(t, err)

		// Add relationship without vice versa
		err = AddRelationship(nil, &product1.ID, nil, &product2.ID, "RELATED_TO", "")
		require.NoError(t, err)

		// Verify only one relationship was created
		relationships, err := ListRelationships()
		require.NoError(t, err)
		assert.Len(t, relationships, 1)

		rel := relationships[0]
		assert.Equal(t, "RELATED_TO", rel.Type)
		assert.Equal(t, product1.ID, *rel.FromProductID)
		assert.Equal(t, product2.ID, *rel.ToProductID)
		assert.Nil(t, rel.FromInstanceID)
		assert.Nil(t, rel.ToInstanceID)
	})

	t.Run("AddRelationship_ValidationErrors", func(t *testing.T) {
		// Create test data
		product1, err := CreateProduct("Validation Product", "Validation test product")
		require.NoError(t, err)
		instance1, err := CreateInstance("Validation Instance", product1.ID)
		require.NoError(t, err)

		// Test multiple "from" attributes
		err = AddRelationship(&instance1.ID, &product1.ID, nil, &product1.ID, "TEST", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one 'from' attribute must be provided")

		// Test no "from" attributes
		err = AddRelationship(nil, nil, nil, &product1.ID, "TEST", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one 'from' attribute must be provided")

		// Test multiple "to" attributes
		err = AddRelationship(&instance1.ID, nil, &instance1.ID, &product1.ID, "TEST", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one 'to' attribute must be provided")

		// Test no "to" attributes
		err = AddRelationship(&instance1.ID, nil, nil, nil, "TEST", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one 'to' attribute must be provided")
	})

	t.Run("DeleteRelationshipById", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		product1, err := CreateProduct("Delete Product", "Delete test product")
		require.NoError(t, err)
		instance1, err := CreateInstance("Delete Instance 1", product1.ID)
		require.NoError(t, err)
		instance2, err := CreateInstance("Delete Instance 2", product1.ID)
		require.NoError(t, err)

		// Create relationship
		err = AddRelationship(&instance1.ID, nil, &instance2.ID, nil, "TEST_RELATIONSHIP", "")
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
		product1, err := CreateProduct("Get Product", "Get test product")
		require.NoError(t, err)
		instance1, err := CreateInstance("Get Instance 1", product1.ID)
		require.NoError(t, err)
		instance2, err := CreateInstance("Get Instance 2", product1.ID)
		require.NoError(t, err)

		// Create relationship
		err = AddRelationship(&instance1.ID, nil, &instance2.ID, nil, "GET_TEST", "")
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
		assert.Equal(t, instance1.ID, *relationship.FromInstanceID)
		assert.Equal(t, instance2.ID, *relationship.ToInstanceID)
	})

	t.Run("GetRelationship_NotFound", func(t *testing.T) {
		// Try to get a non-existent relationship
		nonExistentID := uuid.New()
		relationship, err := GetRelationship(nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("ListRelationshipsByFromInstanceID", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		product1, err := CreateProduct("List Product", "List test product")
		require.NoError(t, err)
		instance1, err := CreateInstance("List Instance 1", product1.ID)
		require.NoError(t, err)
		instance2, err := CreateInstance("List Instance 2", product1.ID)
		require.NoError(t, err)
		instance3, err := CreateInstance("List Instance 3", product1.ID)
		require.NoError(t, err)

		// Create relationships from instance1
		err = AddRelationship(&instance1.ID, nil, &instance2.ID, nil, "REL_TYPE_1", "")
		require.NoError(t, err)
		err = AddRelationship(&instance1.ID, nil, &instance3.ID, nil, "REL_TYPE_2", "")
		require.NoError(t, err)

		// Create relationship from instance2 (should not appear in our list)
		err = AddRelationship(&instance2.ID, nil, &instance3.ID, nil, "REL_TYPE_3", "")
		require.NoError(t, err)

		// List relationships from instance1
		relationships, err := ListRelationshipsByFromInstanceID(instance1.ID)
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Verify all relationships are from instance1
		for _, rel := range relationships {
			assert.Equal(t, instance1.ID, *rel.FromInstanceID)
		}
	})

	t.Run("ListRelationshipsByFromProductID", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		product1, err := CreateProduct("List Product 1", "List test product 1")
		require.NoError(t, err)
		product2, err := CreateProduct("List Product 2", "List test product 2")
		require.NoError(t, err)
		product3, err := CreateProduct("List Product 3", "List test product 3")
		require.NoError(t, err)

		// Create relationships from product1
		err = AddRelationship(nil, &product1.ID, nil, &product2.ID, "PROD_REL_1", "")
		require.NoError(t, err)
		err = AddRelationship(nil, &product1.ID, nil, &product3.ID, "PROD_REL_2", "")
		require.NoError(t, err)

		// Create relationship from product2 (should not appear in our list)
		err = AddRelationship(nil, &product2.ID, nil, &product3.ID, "PROD_REL_3", "")
		require.NoError(t, err)

		// List relationships from product1
		relationships, err := ListRelationshipsByFromProductID(product1.ID)
		require.NoError(t, err)
		assert.Len(t, relationships, 2)

		// Verify all relationships are from product1
		for _, rel := range relationships {
			assert.Equal(t, product1.ID, *rel.FromProductID)
		}
	})

	t.Run("ListRelationshipsByType", func(t *testing.T) {
		// Clear existing relationships
		db := database.GetDB()
		db.Exec("DELETE FROM relationships")

		// Create test data
		product1, err := CreateProduct("Type Product 1", "Type test product 1")
		require.NoError(t, err)
		product2, err := CreateProduct("Type Product 2", "Type test product 2")
		require.NoError(t, err)
		product3, err := CreateProduct("Type Product 3", "Type test product 3")
		require.NoError(t, err)

		// Create relationships of different types
		err = AddRelationship(nil, &product1.ID, nil, &product2.ID, "TYPE_A", "")
		require.NoError(t, err)
		err = AddRelationship(nil, &product1.ID, nil, &product3.ID, "TYPE_A", "")
		require.NoError(t, err)
		err = AddRelationship(nil, &product2.ID, nil, &product3.ID, "TYPE_B", "")
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
