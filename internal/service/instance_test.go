package service

import (
	"fmt"
	"testing"
	"threatreg/internal/database"
	"threatreg/internal/models"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestInstanceService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create a test product first for foreign key reference
	testProduct, err := CreateProduct("Test Product", "A test product for instances")
	require.NoError(t, err)

	t.Run("CreateInstance", func(t *testing.T) {
		// Test data
		name := "Test Instance"
		instanceOf := testProduct.ID

		// Create instance
		instance, err := CreateInstance(name, instanceOf)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, instance)
		assert.NotEqual(t, uuid.Nil, instance.ID)
		assert.Equal(t, name, instance.Name)
		assert.Equal(t, instanceOf, instance.InstanceOf)

		// Verify instance was actually saved to database
		db := database.GetDB()
		var dbInstance models.Instance
		err = db.First(&dbInstance, "id = ?", instance.ID).Error
		require.NoError(t, err)
		assert.Equal(t, instance.ID, dbInstance.ID)
		assert.Equal(t, name, dbInstance.Name)
		assert.Equal(t, instanceOf, dbInstance.InstanceOf)
	})

	t.Run("GetInstance", func(t *testing.T) {
		// Create an instance first
		name := "Get Test Instance"
		instanceOf := testProduct.ID
		createdInstance, err := CreateInstance(name, instanceOf)
		require.NoError(t, err)

		// Get the instance
		retrievedInstance, err := GetInstance(createdInstance.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedInstance)
		assert.Equal(t, createdInstance.ID, retrievedInstance.ID)
		assert.Equal(t, name, retrievedInstance.Name)
		assert.Equal(t, instanceOf, retrievedInstance.InstanceOf)
		assert.Equal(t, testProduct.ID, retrievedInstance.Product.ID)
		assert.Equal(t, testProduct.Name, retrievedInstance.Product.Name)
	})

	t.Run("GetInstance_NotFound", func(t *testing.T) {
		// Try to get a non-existent instance
		nonExistentID := uuid.New()
		instance, err := GetInstance(nonExistentID)

		// Should return error and nil instance
		assert.Error(t, err)
		assert.Nil(t, instance)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateInstance", func(t *testing.T) {
		// Create another product for update test
		anotherProduct, err := CreateProduct("Another Product", "Another test product")
		require.NoError(t, err)

		// Create an instance first
		originalName := "Original Instance"
		originalInstanceOf := testProduct.ID
		createdInstance, err := CreateInstance(originalName, originalInstanceOf)
		require.NoError(t, err)

		// Update the instance
		newName := "Updated Instance"
		newInstanceOf := anotherProduct.ID
		updatedInstance, err := UpdateInstance(createdInstance.ID, &newName, &newInstanceOf)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedInstance)
		assert.Equal(t, createdInstance.ID, updatedInstance.ID)
		assert.Equal(t, newName, updatedInstance.Name)
		assert.Equal(t, newInstanceOf, updatedInstance.InstanceOf)

		// Also verify the Product relationship is loaded correctly
		assert.Equal(t, anotherProduct.ID, updatedInstance.Product.ID)
		assert.Equal(t, anotherProduct.Name, updatedInstance.Product.Name)

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbInstance models.Instance
		err = db.First(&dbInstance, "id = ?", createdInstance.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbInstance.Name)
		assert.Equal(t, newInstanceOf, dbInstance.InstanceOf)
	})

	t.Run("UpdateInstance_PartialUpdate", func(t *testing.T) {
		// Create an instance first
		originalName := "Partial Update Instance"
		originalInstanceOf := testProduct.ID
		createdInstance, err := CreateInstance(originalName, originalInstanceOf)
		require.NoError(t, err)

		// Update only the name
		newName := "New Name Only"
		updatedInstance, err := UpdateInstance(createdInstance.ID, &newName, nil)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedInstance)
		assert.Equal(t, createdInstance.ID, updatedInstance.ID)
		assert.Equal(t, newName, updatedInstance.Name)
		assert.Equal(t, originalInstanceOf, updatedInstance.InstanceOf) // Should remain unchanged

		// Verify the partial update was persisted
		db := database.GetDB()
		var dbInstance models.Instance
		err = db.First(&dbInstance, "id = ?", createdInstance.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbInstance.Name)
		assert.Equal(t, originalInstanceOf, dbInstance.InstanceOf)
	})

	t.Run("UpdateInstance_NotFound", func(t *testing.T) {
		// Try to update a non-existent instance
		nonExistentID := uuid.New()
		newName := "New Name"
		newInstanceOf := testProduct.ID
		instance, err := UpdateInstance(nonExistentID, &newName, &newInstanceOf)

		// Should return error and nil instance
		assert.Error(t, err)
		assert.Nil(t, instance)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteInstance", func(t *testing.T) {
		// Create an instance first
		name := "Delete Test Instance"
		instanceOf := testProduct.ID
		createdInstance, err := CreateInstance(name, instanceOf)
		require.NoError(t, err)

		// Delete the instance
		err = DeleteInstance(createdInstance.ID)

		// Assertions
		require.NoError(t, err)

		// Verify the instance was actually deleted from database
		db := database.GetDB()
		var dbInstance models.Instance
		err = db.First(&dbInstance, "id = ?", createdInstance.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteInstance_NotFound", func(t *testing.T) {
		// Try to delete a non-existent instance
		nonExistentID := uuid.New()
		err := DeleteInstance(nonExistentID)

		// Delete should succeed even if instance doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})

	t.Run("ListInstances", func(t *testing.T) {
		// Clear any existing instances first
		db := database.GetDB()
		db.Exec("DELETE FROM instances")

		// Create multiple instances
		var createdInstances []*models.Instance
		for i := 0; i < 3; i++ {
			name := fmt.Sprintf("Instance %d", i+1)
			instance, err := CreateInstance(name, testProduct.ID)
			require.NoError(t, err)
			createdInstances = append(createdInstances, instance)
		}

		// List all instances
		retrievedInstances, err := ListInstances()

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedInstances, 3)

		// Verify all created instances are in the list
		instanceMap := make(map[uuid.UUID]models.Instance)
		for _, a := range retrievedInstances {
			instanceMap[a.ID] = a
		}

		for _, created := range createdInstances {
			retrieved, exists := instanceMap[created.ID]
			assert.True(t, exists, "Created instance should exist in list")
			assert.Equal(t, created.Name, retrieved.Name)
			assert.Equal(t, created.InstanceOf, retrieved.InstanceOf)
			assert.Equal(t, testProduct.ID, retrieved.Product.ID)
		}
	})

	t.Run("ListInstances_Empty", func(t *testing.T) {
		// Clear all instances
		db := database.GetDB()
		db.Exec("DELETE FROM instances")

		// List instances
		instances, err := ListInstances()

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, instances, 0)
	})

	t.Run("ListInstancesByProductID", func(t *testing.T) {
		// Clear all instances first
		db := database.GetDB()
		db.Exec("DELETE FROM instances")

		// Create another product
		anotherProduct, err := CreateProduct("Another Product", "Another test product")
		require.NoError(t, err)

		// Create instances for both products
		var testProductApps []*models.Instance
		var anotherProductApps []*models.Instance

		for i := 0; i < 2; i++ {
			app1Name := fmt.Sprintf("Test Product App %d", i+1)
			app1, err := CreateInstance(app1Name, testProduct.ID)
			require.NoError(t, err)
			testProductApps = append(testProductApps, app1)

			app2Name := fmt.Sprintf("Another Product App %d", i+1)
			app2, err := CreateInstance(app2Name, anotherProduct.ID)
			require.NoError(t, err)
			anotherProductApps = append(anotherProductApps, app2)
		}

		// List instances by test product ID
		retrievedApps, err := ListInstancesByProductID(testProduct.ID)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedApps, 2)

		for _, app := range retrievedApps {
			assert.Equal(t, testProduct.ID, app.InstanceOf)
			assert.Equal(t, testProduct.ID, app.Product.ID)
		}

		// List instances by another product ID
		retrievedApps2, err := ListInstancesByProductID(anotherProduct.ID)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedApps2, 2)

		for _, app := range retrievedApps2 {
			assert.Equal(t, anotherProduct.ID, app.InstanceOf)
			assert.Equal(t, anotherProduct.ID, app.Product.ID)
		}
	})

	t.Run("ListInstancesByProductID_Empty", func(t *testing.T) {
		// Clear all instances
		db := database.GetDB()
		db.Exec("DELETE FROM instances")

		// List instances for a product with no instances
		instances, err := ListInstancesByProductID(testProduct.ID)

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, instances, 0)
	})

	t.Run("FilterInstances", func(t *testing.T) {
		// Clear all instances first
		db := database.GetDB()
		db.Exec("DELETE FROM instances")

		// Create multiple products for testing
		webProduct, err := CreateProduct("Web Platform", "Web-based platform")
		require.NoError(t, err)
		mobileProduct, err := CreateProduct("Mobile App", "Mobile instance")
		require.NoError(t, err)

		// Create instances with different names
		testApps := []struct {
			name      string
			productID uuid.UUID
		}{
			{"Production Web", webProduct.ID},
			{"Staging Web", webProduct.ID},
			{"Development Web", webProduct.ID},
			{"Production Mobile", mobileProduct.ID},
			{"Test Mobile", mobileProduct.ID},
			{"API Gateway", webProduct.ID},
		}

		var createdApps []*models.Instance
		for _, ta := range testApps {
			app, err := CreateInstance(ta.name, ta.productID)
			require.NoError(t, err)
			createdApps = append(createdApps, app)
		}

		t.Run("FilterByInstanceName", func(t *testing.T) {
			// Filter by instance name containing "Production"
			apps, err := FilterInstances("Production", "")
			require.NoError(t, err)
			assert.Len(t, apps, 2)

			for _, app := range apps {
				assert.Contains(t, app.Name, "Production")
			}
		})

		t.Run("FilterByProductName", func(t *testing.T) {
			// Filter by product name containing "Web"
			apps, err := FilterInstances("", "Web")
			require.NoError(t, err)
			assert.Len(t, apps, 4) // All Web Platform instances

			for _, app := range apps {
				assert.Contains(t, app.Product.Name, "Web")
			}
		})

		t.Run("FilterByBothNames", func(t *testing.T) {
			// Filter by both instance and product name
			apps, err := FilterInstances("Production", "Mobile")
			require.NoError(t, err)
			assert.Len(t, apps, 1)

			app := apps[0]
			assert.Contains(t, app.Name, "Production")
			assert.Contains(t, app.Product.Name, "Mobile")
		})

		t.Run("FilterCaseInsensitive", func(t *testing.T) {
			// Test case insensitive filtering
			apps, err := FilterInstances("production", "web")
			require.NoError(t, err)
			assert.Len(t, apps, 1)

			app := apps[0]
			assert.Equal(t, "Production Web", app.Name)
			assert.Contains(t, app.Product.Name, "Web")
		})

		t.Run("FilterNoMatch", func(t *testing.T) {
			// Filter with no matches
			apps, err := FilterInstances("NonExistent", "")
			require.NoError(t, err)
			assert.Len(t, apps, 0)
		})

		t.Run("FilterEmptyStrings", func(t *testing.T) {
			// Filter with empty strings should return all
			apps, err := FilterInstances("", "")
			require.NoError(t, err)
			assert.Len(t, apps, 6) // All created instances
		})

		t.Run("FilterPartialMatch", func(t *testing.T) {
			// Filter with partial name match
			apps, err := FilterInstances("Web", "")
			require.NoError(t, err)
			assert.Len(t, apps, 3) // Production Web, Staging Web, Development Web

			for _, app := range apps {
				assert.Contains(t, app.Name, "Web")
			}
		})
	})
}
