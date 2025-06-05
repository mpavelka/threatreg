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
}
