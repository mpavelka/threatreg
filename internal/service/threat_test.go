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

func TestThreatService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("CreateThreat", func(t *testing.T) {
		// Test data
		title := "SQL Injection"
		description := "Attackers can execute arbitrary SQL code"

		// Create threat
		threat, err := CreateThreat(title, description)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, threat)
		assert.NotEqual(t, uuid.Nil, threat.ID)
		assert.Equal(t, title, threat.Title)
		assert.Equal(t, description, threat.Description)

		// Verify threat was actually saved to database
		db := database.GetDB()
		var dbThreat models.Threat
		err = db.First(&dbThreat, "id = ?", threat.ID).Error
		require.NoError(t, err)
		assert.Equal(t, threat.ID, dbThreat.ID)
		assert.Equal(t, title, dbThreat.Title)
		assert.Equal(t, description, dbThreat.Description)
	})

	t.Run("GetThreat", func(t *testing.T) {
		// Create a threat first
		title := "Cross-Site Scripting"
		description := "Injection of malicious scripts into web pages"
		createdThreat, err := CreateThreat(title, description)
		require.NoError(t, err)

		// Get the threat
		retrievedThreat, err := GetThreat(createdThreat.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedThreat)
		assert.Equal(t, createdThreat.ID, retrievedThreat.ID)
		assert.Equal(t, title, retrievedThreat.Title)
		assert.Equal(t, description, retrievedThreat.Description)
	})

	t.Run("GetThreat_NotFound", func(t *testing.T) {
		// Try to get a non-existent threat
		nonExistentID := uuid.New()
		threat, err := GetThreat(nonExistentID)

		// Should return error and nil threat
		assert.Error(t, err)
		assert.Nil(t, threat)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateThreat", func(t *testing.T) {
		// Create a threat first
		originalTitle := "Buffer Overflow"
		originalDescription := "Memory corruption vulnerability"
		createdThreat, err := CreateThreat(originalTitle, originalDescription)
		require.NoError(t, err)

		// Update the threat
		newTitle := "Stack Buffer Overflow"
		newDescription := "Stack-based memory corruption vulnerability"
		updatedThreat, err := UpdateThreat(createdThreat.ID, &newTitle, &newDescription)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedThreat)
		assert.Equal(t, createdThreat.ID, updatedThreat.ID)
		assert.Equal(t, newTitle, updatedThreat.Title)
		assert.Equal(t, newDescription, updatedThreat.Description)

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbThreat models.Threat
		err = db.First(&dbThreat, "id = ?", createdThreat.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newTitle, dbThreat.Title)
		assert.Equal(t, newDescription, dbThreat.Description)
	})

	t.Run("UpdateThreat_PartialUpdate", func(t *testing.T) {
		// Create a threat first
		originalTitle := "CSRF Attack"
		originalDescription := "Cross-Site Request Forgery vulnerability"
		createdThreat, err := CreateThreat(originalTitle, originalDescription)
		require.NoError(t, err)

		// Update only the title
		newTitle := "Cross-Site Request Forgery"
		updatedThreat, err := UpdateThreat(createdThreat.ID, &newTitle, nil)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedThreat)
		assert.Equal(t, createdThreat.ID, updatedThreat.ID)
		assert.Equal(t, newTitle, updatedThreat.Title)
		assert.Equal(t, originalDescription, updatedThreat.Description) // Should remain unchanged

		// Verify the partial update was persisted
		db := database.GetDB()
		var dbThreat models.Threat
		err = db.First(&dbThreat, "id = ?", createdThreat.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newTitle, dbThreat.Title)
		assert.Equal(t, originalDescription, dbThreat.Description)
	})

	t.Run("UpdateThreat_NotFound", func(t *testing.T) {
		// Try to update a non-existent threat
		nonExistentID := uuid.New()
		newTitle := "Non-existent Threat"
		threat, err := UpdateThreat(nonExistentID, &newTitle, nil)

		// Should return error and nil threat
		assert.Error(t, err)
		assert.Nil(t, threat)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteThreat", func(t *testing.T) {
		// Create a threat first
		title := "Privilege Escalation"
		description := "Gaining higher access privileges than intended"
		createdThreat, err := CreateThreat(title, description)
		require.NoError(t, err)

		// Delete the threat
		err = DeleteThreat(createdThreat.ID)

		// Assertions
		require.NoError(t, err)

		// Verify the threat was actually deleted from database
		db := database.GetDB()
		var dbThreat models.Threat
		err = db.First(&dbThreat, "id = ?", createdThreat.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteThreat_NotFound", func(t *testing.T) {
		// Try to delete a non-existent threat
		nonExistentID := uuid.New()
		err := DeleteThreat(nonExistentID)

		// Delete should succeed even if threat doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})

	t.Run("ListThreats", func(t *testing.T) {
		// Clear any existing threats first
		db := database.GetDB()
		db.Exec("DELETE FROM threats")

		// Create multiple threats
		threats := []struct {
			title       string
			description string
		}{
			{"SQL Injection", "Database query manipulation"},
			{"XSS", "Cross-site scripting attack"},
			{"CSRF", "Cross-site request forgery"},
		}

		var createdThreats []*models.Threat
		for _, threatData := range threats {
			threat, err := CreateThreat(threatData.title, threatData.description)
			require.NoError(t, err)
			createdThreats = append(createdThreats, threat)
		}

		// List all threats
		retrievedThreats, err := ListThreats()

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedThreats, len(threats))

		// Verify all created threats are in the list
		threatMap := make(map[uuid.UUID]models.Threat)
		for _, threat := range retrievedThreats {
			threatMap[threat.ID] = threat
		}

		for _, created := range createdThreats {
			retrieved, exists := threatMap[created.ID]
			assert.True(t, exists, "Created threat should exist in list")
			assert.Equal(t, created.Title, retrieved.Title)
			assert.Equal(t, created.Description, retrieved.Description)
		}
	})

	t.Run("ListThreats_Empty", func(t *testing.T) {
		// Clear all threats
		db := database.GetDB()
		db.Exec("DELETE FROM threats")

		// List threats
		threats, err := ListThreats()

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, threats, 0)
	})
}