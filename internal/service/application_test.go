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

func TestApplicationService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create a test product first for foreign key reference
	testProduct, err := CreateProduct("Test Product", "A test product for applications")
	require.NoError(t, err)

	t.Run("CreateApplication", func(t *testing.T) {
		// Test data
		instanceOf := testProduct.ID

		// Create application
		application, err := CreateApplication(instanceOf)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, application)
		assert.NotEqual(t, uuid.Nil, application.ID)
		assert.Equal(t, instanceOf, application.InstanceOf)

		// Verify application was actually saved to database
		db := database.GetDB()
		var dbApplication models.Application
		err = db.First(&dbApplication, "id = ?", application.ID).Error
		require.NoError(t, err)
		assert.Equal(t, application.ID, dbApplication.ID)
		assert.Equal(t, instanceOf, dbApplication.InstanceOf)
	})

	t.Run("GetApplication", func(t *testing.T) {
		// Create an application first
		instanceOf := testProduct.ID
		createdApplication, err := CreateApplication(instanceOf)
		require.NoError(t, err)

		// Get the application
		retrievedApplication, err := GetApplication(createdApplication.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedApplication)
		assert.Equal(t, createdApplication.ID, retrievedApplication.ID)
		assert.Equal(t, instanceOf, retrievedApplication.InstanceOf)
		assert.Equal(t, testProduct.ID, retrievedApplication.Product.ID)
		assert.Equal(t, testProduct.Name, retrievedApplication.Product.Name)
	})

	t.Run("GetApplication_NotFound", func(t *testing.T) {
		// Try to get a non-existent application
		nonExistentID := uuid.New()
		application, err := GetApplication(nonExistentID)

		// Should return error and nil application
		assert.Error(t, err)
		assert.Nil(t, application)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateApplication", func(t *testing.T) {
		// Create another product for update test
		anotherProduct, err := CreateProduct("Another Product", "Another test product")
		require.NoError(t, err)

		// Create an application first
		originalInstanceOf := testProduct.ID
		createdApplication, err := CreateApplication(originalInstanceOf)
		require.NoError(t, err)

		// Update the application
		newInstanceOf := anotherProduct.ID
		updatedApplication, err := UpdateApplication(createdApplication.ID, &newInstanceOf)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedApplication)
		assert.Equal(t, createdApplication.ID, updatedApplication.ID)
		assert.Equal(t, newInstanceOf, updatedApplication.InstanceOf)
		
		// Also verify the Product relationship is loaded correctly
		assert.Equal(t, anotherProduct.ID, updatedApplication.Product.ID)
		assert.Equal(t, anotherProduct.Name, updatedApplication.Product.Name)

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbApplication models.Application
		err = db.First(&dbApplication, "id = ?", createdApplication.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newInstanceOf, dbApplication.InstanceOf)
	})

	t.Run("UpdateApplication_NoChange", func(t *testing.T) {
		// Create an application first
		originalInstanceOf := testProduct.ID
		createdApplication, err := CreateApplication(originalInstanceOf)
		require.NoError(t, err)

		// Update with nil (no change)
		updatedApplication, err := UpdateApplication(createdApplication.ID, nil)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedApplication)
		assert.Equal(t, createdApplication.ID, updatedApplication.ID)
		assert.Equal(t, originalInstanceOf, updatedApplication.InstanceOf) // Should remain unchanged

		// Verify no change was persisted
		db := database.GetDB()
		var dbApplication models.Application
		err = db.First(&dbApplication, "id = ?", createdApplication.ID).Error
		require.NoError(t, err)
		assert.Equal(t, originalInstanceOf, dbApplication.InstanceOf)
	})

	t.Run("UpdateApplication_NotFound", func(t *testing.T) {
		// Try to update a non-existent application
		nonExistentID := uuid.New()
		newInstanceOf := testProduct.ID
		application, err := UpdateApplication(nonExistentID, &newInstanceOf)

		// Should return error and nil application
		assert.Error(t, err)
		assert.Nil(t, application)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteApplication", func(t *testing.T) {
		// Create an application first
		instanceOf := testProduct.ID
		createdApplication, err := CreateApplication(instanceOf)
		require.NoError(t, err)

		// Delete the application
		err = DeleteApplication(createdApplication.ID)

		// Assertions
		require.NoError(t, err)

		// Verify the application was actually deleted from database
		db := database.GetDB()
		var dbApplication models.Application
		err = db.First(&dbApplication, "id = ?", createdApplication.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteApplication_NotFound", func(t *testing.T) {
		// Try to delete a non-existent application
		nonExistentID := uuid.New()
		err := DeleteApplication(nonExistentID)

		// Delete should succeed even if application doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})

	t.Run("ListApplications", func(t *testing.T) {
		// Clear any existing applications first
		db := database.GetDB()
		db.Exec("DELETE FROM applications")

		// Create multiple applications
		var createdApplications []*models.Application
		for i := 0; i < 3; i++ {
			application, err := CreateApplication(testProduct.ID)
			require.NoError(t, err)
			createdApplications = append(createdApplications, application)
		}

		// List all applications
		retrievedApplications, err := ListApplications()

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedApplications, 3)

		// Verify all created applications are in the list
		applicationMap := make(map[uuid.UUID]models.Application)
		for _, a := range retrievedApplications {
			applicationMap[a.ID] = a
		}

		for _, created := range createdApplications {
			retrieved, exists := applicationMap[created.ID]
			assert.True(t, exists, "Created application should exist in list")
			assert.Equal(t, created.InstanceOf, retrieved.InstanceOf)
			assert.Equal(t, testProduct.ID, retrieved.Product.ID)
		}
	})

	t.Run("ListApplications_Empty", func(t *testing.T) {
		// Clear all applications
		db := database.GetDB()
		db.Exec("DELETE FROM applications")

		// List applications
		applications, err := ListApplications()

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, applications, 0)
	})

	t.Run("ListApplicationsByProductID", func(t *testing.T) {
		// Clear all applications first
		db := database.GetDB()
		db.Exec("DELETE FROM applications")

		// Create another product
		anotherProduct, err := CreateProduct("Another Product", "Another test product")
		require.NoError(t, err)

		// Create applications for both products
		var testProductApps []*models.Application
		var anotherProductApps []*models.Application

		for i := 0; i < 2; i++ {
			app1, err := CreateApplication(testProduct.ID)
			require.NoError(t, err)
			testProductApps = append(testProductApps, app1)

			app2, err := CreateApplication(anotherProduct.ID)
			require.NoError(t, err)
			anotherProductApps = append(anotherProductApps, app2)
		}

		// List applications by test product ID
		retrievedApps, err := ListApplicationsByProductID(testProduct.ID)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedApps, 2)

		for _, app := range retrievedApps {
			assert.Equal(t, testProduct.ID, app.InstanceOf)
			assert.Equal(t, testProduct.ID, app.Product.ID)
		}

		// List applications by another product ID
		retrievedApps2, err := ListApplicationsByProductID(anotherProduct.ID)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedApps2, 2)

		for _, app := range retrievedApps2 {
			assert.Equal(t, anotherProduct.ID, app.InstanceOf)
			assert.Equal(t, anotherProduct.ID, app.Product.ID)
		}
	})

	t.Run("ListApplicationsByProductID_Empty", func(t *testing.T) {
		// Clear all applications
		db := database.GetDB()
		db.Exec("DELETE FROM applications")

		// List applications for a product with no applications
		applications, err := ListApplicationsByProductID(testProduct.ID)

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, applications, 0)
	})
}