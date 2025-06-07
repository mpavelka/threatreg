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

func TestControlService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("CreateControl", func(t *testing.T) {
		// Test data
		title := "Test Control"
		description := "A test control description"

		// Create control
		control, err := CreateControl(title, description)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, control)
		assert.NotEqual(t, uuid.Nil, control.ID)
		assert.Equal(t, title, control.Title)
		assert.Equal(t, description, control.Description)

		// Verify control was actually saved to database
		db := database.GetDB()
		var dbControl models.Control
		err = db.First(&dbControl, "id = ?", control.ID).Error
		require.NoError(t, err)
		assert.Equal(t, control.ID, dbControl.ID)
		assert.Equal(t, title, dbControl.Title)
		assert.Equal(t, description, dbControl.Description)
	})

	t.Run("GetControl", func(t *testing.T) {
		// Create a control first
		title := "Get Test Control"
		description := "Control for get test"
		createdControl, err := CreateControl(title, description)
		require.NoError(t, err)

		// Get the control
		retrievedControl, err := GetControl(createdControl.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedControl)
		assert.Equal(t, createdControl.ID, retrievedControl.ID)
		assert.Equal(t, title, retrievedControl.Title)
		assert.Equal(t, description, retrievedControl.Description)
	})

	t.Run("GetControl_NotFound", func(t *testing.T) {
		// Try to get a non-existent control
		nonExistentID := uuid.New()
		control, err := GetControl(nonExistentID)

		// Should return error and nil control
		assert.Error(t, err)
		assert.Nil(t, control)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateControl", func(t *testing.T) {
		// Create a control first
		originalTitle := "Original Control"
		originalDescription := "Original description"
		createdControl, err := CreateControl(originalTitle, originalDescription)
		require.NoError(t, err)

		// Update the control
		newTitle := "Updated Control"
		newDescription := "Updated description"
		updatedControl, err := UpdateControl(createdControl.ID, &newTitle, &newDescription)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedControl)
		assert.Equal(t, createdControl.ID, updatedControl.ID)
		assert.Equal(t, newTitle, updatedControl.Title)
		assert.Equal(t, newDescription, updatedControl.Description)

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbControl models.Control
		err = db.First(&dbControl, "id = ?", createdControl.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newTitle, dbControl.Title)
		assert.Equal(t, newDescription, dbControl.Description)
	})

	t.Run("UpdateControl_PartialUpdate", func(t *testing.T) {
		// Create a control first
		originalTitle := "Partial Update Control"
		originalDescription := "Original description"
		createdControl, err := CreateControl(originalTitle, originalDescription)
		require.NoError(t, err)

		// Update only the title
		newTitle := "New Title Only"
		updatedControl, err := UpdateControl(createdControl.ID, &newTitle, nil)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedControl)
		assert.Equal(t, createdControl.ID, updatedControl.ID)
		assert.Equal(t, newTitle, updatedControl.Title)
		assert.Equal(t, originalDescription, updatedControl.Description) // Should remain unchanged

		// Verify the partial update was persisted
		db := database.GetDB()
		var dbControl models.Control
		err = db.First(&dbControl, "id = ?", createdControl.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newTitle, dbControl.Title)
		assert.Equal(t, originalDescription, dbControl.Description)
	})

	t.Run("UpdateControl_NotFound", func(t *testing.T) {
		// Try to update a non-existent control
		nonExistentID := uuid.New()
		newTitle := "New Title"
		control, err := UpdateControl(nonExistentID, &newTitle, nil)

		// Should return error and nil control
		assert.Error(t, err)
		assert.Nil(t, control)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteControl", func(t *testing.T) {
		// Create a control first
		title := "Delete Test Control"
		description := "Control to be deleted"
		createdControl, err := CreateControl(title, description)
		require.NoError(t, err)

		// Delete the control
		err = DeleteControl(createdControl.ID)

		// Assertions
		require.NoError(t, err)

		// Verify the control was actually deleted from database
		db := database.GetDB()
		var dbControl models.Control
		err = db.First(&dbControl, "id = ?", createdControl.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteControl_NotFound", func(t *testing.T) {
		// Try to delete a non-existent control
		nonExistentID := uuid.New()
		err := DeleteControl(nonExistentID)

		// Delete should succeed even if control doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})

	t.Run("ListControls", func(t *testing.T) {
		// Clear any existing controls first
		db := database.GetDB()
		db.Exec("DELETE FROM controls")

		// Create multiple controls
		controls := []struct {
			title       string
			description string
		}{
			{"Control 1", "Description 1"},
			{"Control 2", "Description 2"},
			{"Control 3", "Description 3"},
		}

		var createdControls []*models.Control
		for _, c := range controls {
			control, err := CreateControl(c.title, c.description)
			require.NoError(t, err)
			createdControls = append(createdControls, control)
		}

		// List all controls
		retrievedControls, err := ListControls()

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedControls, len(controls))

		// Verify all created controls are in the list
		controlMap := make(map[uuid.UUID]models.Control)
		for _, c := range retrievedControls {
			controlMap[c.ID] = c
		}

		for _, created := range createdControls {
			retrieved, exists := controlMap[created.ID]
			assert.True(t, exists, "Created control should exist in list")
			assert.Equal(t, created.Title, retrieved.Title)
			assert.Equal(t, created.Description, retrieved.Description)
		}
	})

	t.Run("ListControls_Empty", func(t *testing.T) {
		// Clear all controls
		db := database.GetDB()
		db.Exec("DELETE FROM controls")

		// List controls
		controls, err := ListControls()

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, controls, 0)
	})
}
