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

func TestProductService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("CreateProduct", func(t *testing.T) {
		// Test data
		name := "Test Product"
		description := "A test product description"

		// Create product
		product, err := CreateProduct(name, description)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, product)
		assert.NotEqual(t, uuid.Nil, product.ID)
		assert.Equal(t, name, product.Name)
		assert.Equal(t, description, product.Description)

		// Verify product was actually saved to database
		db := database.GetDB()
		var dbProduct models.Product
		err = db.First(&dbProduct, "id = ?", product.ID).Error
		require.NoError(t, err)
		assert.Equal(t, product.ID, dbProduct.ID)
		assert.Equal(t, name, dbProduct.Name)
		assert.Equal(t, description, dbProduct.Description)
	})

	t.Run("GetProduct", func(t *testing.T) {
		// Create a product first
		name := "Get Test Product"
		description := "Product for get test"
		createdProduct, err := CreateProduct(name, description)
		require.NoError(t, err)

		// Get the product
		retrievedProduct, err := GetProduct(createdProduct.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedProduct)
		assert.Equal(t, createdProduct.ID, retrievedProduct.ID)
		assert.Equal(t, name, retrievedProduct.Name)
		assert.Equal(t, description, retrievedProduct.Description)
	})

	t.Run("GetProduct_NotFound", func(t *testing.T) {
		// Try to get a non-existent product
		nonExistentID := uuid.New()
		product, err := GetProduct(nonExistentID)

		// Should return error and nil product
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		// Create a product first
		originalName := "Original Product"
		originalDescription := "Original description"
		createdProduct, err := CreateProduct(originalName, originalDescription)
		require.NoError(t, err)

		// Update the product
		newName := "Updated Product"
		newDescription := "Updated description"
		updatedProduct, err := UpdateProduct(createdProduct.ID, &newName, &newDescription)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedProduct)
		assert.Equal(t, createdProduct.ID, updatedProduct.ID)
		assert.Equal(t, newName, updatedProduct.Name)
		assert.Equal(t, newDescription, updatedProduct.Description)

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbProduct models.Product
		err = db.First(&dbProduct, "id = ?", createdProduct.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbProduct.Name)
		assert.Equal(t, newDescription, dbProduct.Description)
	})

	t.Run("UpdateProduct_PartialUpdate", func(t *testing.T) {
		// Create a product first
		originalName := "Partial Update Product"
		originalDescription := "Original description"
		createdProduct, err := CreateProduct(originalName, originalDescription)
		require.NoError(t, err)

		// Update only the name
		newName := "New Name Only"
		updatedProduct, err := UpdateProduct(createdProduct.ID, &newName, nil)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedProduct)
		assert.Equal(t, createdProduct.ID, updatedProduct.ID)
		assert.Equal(t, newName, updatedProduct.Name)
		assert.Equal(t, originalDescription, updatedProduct.Description) // Should remain unchanged

		// Verify the partial update was persisted
		db := database.GetDB()
		var dbProduct models.Product
		err = db.First(&dbProduct, "id = ?", createdProduct.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbProduct.Name)
		assert.Equal(t, originalDescription, dbProduct.Description)
	})

	t.Run("UpdateProduct_NotFound", func(t *testing.T) {
		// Try to update a non-existent product
		nonExistentID := uuid.New()
		newName := "New Name"
		product, err := UpdateProduct(nonExistentID, &newName, nil)

		// Should return error and nil product
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		// Create a product first
		name := "Delete Test Product"
		description := "Product to be deleted"
		createdProduct, err := CreateProduct(name, description)
		require.NoError(t, err)

		// Delete the product
		err = DeleteProduct(createdProduct.ID)

		// Assertions
		require.NoError(t, err)

		// Verify the product was actually deleted from database
		db := database.GetDB()
		var dbProduct models.Product
		err = db.First(&dbProduct, "id = ?", createdProduct.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteProduct_NotFound", func(t *testing.T) {
		// Try to delete a non-existent product
		nonExistentID := uuid.New()
		err := DeleteProduct(nonExistentID)

		// Delete should succeed even if product doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})

	t.Run("ListProducts", func(t *testing.T) {
		// Clear any existing products first
		db := database.GetDB()
		db.Exec("DELETE FROM products")

		// Create multiple products
		products := []struct {
			name        string
			description string
		}{
			{"Product 1", "Description 1"},
			{"Product 2", "Description 2"},
			{"Product 3", "Description 3"},
		}

		var createdProducts []*models.Product
		for _, p := range products {
			product, err := CreateProduct(p.name, p.description)
			require.NoError(t, err)
			createdProducts = append(createdProducts, product)
		}

		// List all products
		retrievedProducts, err := ListProducts()

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedProducts, len(products))

		// Verify all created products are in the list
		productMap := make(map[uuid.UUID]models.Product)
		for _, p := range retrievedProducts {
			productMap[p.ID] = p
		}

		for _, created := range createdProducts {
			retrieved, exists := productMap[created.ID]
			assert.True(t, exists, "Created product should exist in list")
			assert.Equal(t, created.Name, retrieved.Name)
			assert.Equal(t, created.Description, retrieved.Description)
		}
	})

	t.Run("ListProducts_Empty", func(t *testing.T) {
		// Clear all products
		db := database.GetDB()
		db.Exec("DELETE FROM products")

		// List products
		products, err := ListProducts()

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, products, 0)
	})

	t.Run("AssignThreatToProduct", func(t *testing.T) {
		// Create a product first
		product, err := CreateProduct("Test Product for Threat", "A test product for threat assignment")
		require.NoError(t, err)

		// Create a threat first
		threat, err := CreateThreat("Test Threat", "A test threat for assignment")
		require.NoError(t, err)

		// Assign threat to product
		assignment, err := AssignThreatToProduct(product.ID, threat.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.NotEqual(t, 0, assignment.ID)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, product.ID, assignment.ProductID)
		assert.Equal(t, uuid.Nil, assignment.InstanceID) // Should be nil for product assignment

		// Verify assignment was saved to database
		db := database.GetDB()
		var dbAssignment models.ThreatAssignment
		err = db.First(&dbAssignment, "id = ?", assignment.ID).Error
		require.NoError(t, err)
		assert.Equal(t, assignment.ThreatID, dbAssignment.ThreatID)
		assert.Equal(t, assignment.ProductID, dbAssignment.ProductID)
		assert.Equal(t, uuid.Nil, dbAssignment.InstanceID)
	})

	t.Run("AssignThreatToProduct_Duplicate", func(t *testing.T) {
		// Create a product and threat
		product, err := CreateProduct("Duplicate Test Product", "A test product for duplicate assignment")
		require.NoError(t, err)

		threat, err := CreateThreat("Duplicate Test Threat", "A test threat for duplicate assignment")
		require.NoError(t, err)

		// Assign threat to product first time
		assignment1, err := AssignThreatToProduct(product.ID, threat.ID)
		require.NoError(t, err)
		require.NotNil(t, assignment1)

		// Try to assign the same threat to the same product again
		assignment2, err := AssignThreatToProduct(product.ID, threat.ID)

		// Should return the existing assignment, not create a new one
		require.NoError(t, err)
		assert.NotNil(t, assignment2)
		assert.Equal(t, assignment1.ID, assignment2.ID)
		assert.Equal(t, assignment1.ThreatID, assignment2.ThreatID)
		assert.Equal(t, assignment1.ProductID, assignment2.ProductID)

		// Verify that only one assignment exists in the database for this threat/product combination
		db := database.GetDB()
		var count int64
		err = db.Model(&models.ThreatAssignment{}).
			Where("threat_id = ? AND product_id = ? AND (instance_id IS NULL OR instance_id = ?)", threat.ID, product.ID, uuid.Nil).
			Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Should only have one assignment for this threat/product combination")
	})

	t.Run("AssignThreatToProduct_InvalidThreatID", func(t *testing.T) {
		// Create a product
		product, err := CreateProduct("Invalid Threat Test Product", "A test product")
		require.NoError(t, err)

		// Try to assign non-existent threat
		nonExistentThreatID := uuid.New()
		assignment, err := AssignThreatToProduct(product.ID, nonExistentThreatID)

		// Should succeed (foreign key constraint allows it, but relationship won't load)
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, nonExistentThreatID, assignment.ThreatID)
		assert.Equal(t, product.ID, assignment.ProductID)
	})

	t.Run("AssignThreatToProduct_InvalidProductID", func(t *testing.T) {
		// Create a threat
		threat, err := CreateThreat("Invalid Product Test Threat", "A test threat")
		require.NoError(t, err)

		// Try to assign to non-existent product
		nonExistentProductID := uuid.New()
		assignment, err := AssignThreatToProduct(nonExistentProductID, threat.ID)

		// Should succeed (foreign key constraint allows it, but relationship won't load)
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, nonExistentProductID, assignment.ProductID)
	})

	t.Run("ListThreatAssignmentsByProductID", func(t *testing.T) {
		// Create a product for testing
		product, err := CreateProduct("Test Product for Assignments", "A test product for threat assignment listing")
		require.NoError(t, err)

		// Create multiple threats for testing
		threat1, err := CreateThreat("Product Threat 1", "First product test threat")
		require.NoError(t, err)
		threat2, err := CreateThreat("Product Threat 2", "Second product test threat")
		require.NoError(t, err)
		threat3, err := CreateThreat("Product Threat 3", "Third product test threat")
		require.NoError(t, err)

		// Assign threats to the product
		assignment1, err := AssignThreatToProduct(product.ID, threat1.ID)
		require.NoError(t, err)
		assignment2, err := AssignThreatToProduct(product.ID, threat2.ID)
		require.NoError(t, err)

		// Assign one threat to a different product to ensure filtering works
		otherProduct, err := CreateProduct("Other Product", "Another product for testing")
		require.NoError(t, err)
		_, err = AssignThreatToProduct(otherProduct.ID, threat3.ID)
		require.NoError(t, err)

		// List threat assignments for our test product
		assignments, err := ListThreatAssignmentsByProductID(product.ID)

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
			assert.Equal(t, product.ID, assignment.ProductID)
			assert.Equal(t, uuid.Nil, assignment.InstanceID) // Should be nil for product assignments
		}
	})

	t.Run("ListThreatAssignmentsByProductID_Empty", func(t *testing.T) {
		// Create a product with no threat assignments
		product, err := CreateProduct("Empty Assignments Product", "A product with no threat assignments")
		require.NoError(t, err)

		// List threat assignments for this product
		assignments, err := ListThreatAssignmentsByProductID(product.ID)

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, assignments, 0)
	})

	t.Run("ListThreatAssignmentsByProductID_InvalidProductID", func(t *testing.T) {
		// Try to list assignments for non-existent product
		nonExistentProductID := uuid.New()
		assignments, err := ListThreatAssignmentsByProductID(nonExistentProductID)

		// Should succeed but return empty slice
		require.NoError(t, err)
		assert.Len(t, assignments, 0)
	})
}

func TestListThreatAssignmentsByProductIDWithResolutionByInstanceID(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("Happy Flow", func(t *testing.T) {
		// Create product1 and product2
		product1, err := CreateProduct("Product 1", "First test product")
		require.NoError(t, err)
		product2, err := CreateProduct("Product 2", "Second test product")
		require.NoError(t, err)

		// Create instance1 (for product1) and instance2 (for product2)
		instance1, err := CreateInstance("Instance 1", product1.ID)
		require.NoError(t, err)
		instance2, err := CreateInstance("Instance 2", product2.ID)
		require.NoError(t, err)

		// Create threat1 and threat2
		threat1, err := CreateThreat("Threat 1", "First test threat")
		require.NoError(t, err)
		threat2, err := CreateThreat("Threat 2", "Second test threat")
		require.NoError(t, err)

		// Assign threat1 to product1
		assignment1, err := AssignThreatToProduct(product1.ID, threat1.ID)
		require.NoError(t, err)

		// Assign threat2 to instance2
		assignment2, err := AssignThreatToInstance(threat2.ID, instance2.ID)
		require.NoError(t, err)

		// Create threat resolutions for both instance1 and instance2
		resolution1, err := CreateThreatResolution(
			assignment1.ID,
			&instance1.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution for instance1",
		)
		require.NoError(t, err)

		resolution2, err := CreateThreatResolution(
			assignment2.ID,
			&instance2.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution for instance2",
		)
		require.NoError(t, err)

		// Delegate resolution1 to resolution2
		err = DelegateResolution(*resolution1, *resolution2)
		require.NoError(t, err)

		// Mark resolution2 as resolved
		resolvedStatus := models.ThreatAssignmentResolutionStatusResolved
		_, err = UpdateThreatResolution(resolution2.ID, &resolvedStatus, nil)
		require.NoError(t, err)

		// Test ListThreatAssignmentsByProductIDWithResolutionByInstanceID for product1 filtered by instance1
		results, err := ListThreatAssignmentsByProductIDWithResolutionByInstanceID(product1.ID, instance1.ID)
		require.NoError(t, err)

		// Assertions
		assert.Len(t, results, 1, "Should return one threat assignment for product1")

		result := results[0]
		assert.Equal(t, assignment1.ID, result.ID)
		assert.Equal(t, threat1.ID, result.ThreatID)
		assert.Equal(t, product1.ID, result.ProductID)
		assert.Equal(t, uuid.Nil, result.InstanceID) // Product assignment should have nil instance ID

		// Verify threat relationship is loaded
		assert.Equal(t, threat1.ID, result.Threat.ID)
		assert.Equal(t, "Threat 1", result.Threat.Title)

		// Verify product relationship is loaded
		assert.Equal(t, product1.ID, result.Product.ID)
		assert.Equal(t, "Product 1", result.Product.Name)

		// Verify resolution status - should show resolved because resolution1 was delegated to resolution2 which is resolved
		assert.NotNil(t, result.ResolutionStatus)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, *result.ResolutionStatus)

		// Verify delegation status
		assert.True(t, result.IsDelegated, "Should show as delegated")

		// Test with different instance filter - should return same assignment but no resolution info
		resultsOtherInstance, err := ListThreatAssignmentsByProductIDWithResolutionByInstanceID(product1.ID, instance2.ID)
		require.NoError(t, err)

		assert.Len(t, resultsOtherInstance, 1, "Should still return the threat assignment for product1")
		otherResult := resultsOtherInstance[0]
		assert.Equal(t, assignment1.ID, otherResult.ID)

		// But resolution info should be nil since we filtered by instance2 but resolution1 is for instance1
		assert.Nil(t, otherResult.ResolutionStatus, "Should not have resolution status for different instance")
		assert.False(t, otherResult.IsDelegated, "Should not show as delegated for different instance")

		// Test with product2 - should return no assignments since threat2 was assigned to instance2, not product2
		resultsProduct2, err := ListThreatAssignmentsByProductIDWithResolutionByInstanceID(product2.ID, instance2.ID)
		require.NoError(t, err)
		assert.Len(t, resultsProduct2, 0, "Should return no assignments for product2 since threat2 was assigned to instance2")
	})
}
