package threat_pattern

import (
	"testing"
	"threatreg/internal/database"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestThreatPatternService_CreateThreatPattern(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Create a threat first
		threat, err := service.CreateThreat("Test Threat", "A test threat for pattern testing")
		require.NoError(t, err)

		// Test data
		name := "Test Threat Pattern"
		description := "A test threat pattern"
		isActive := true

		// Create threat pattern
		pattern, err := CreateThreatPattern(name, description, threat.ID, isActive)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, pattern)
		assert.NotEqual(t, uuid.Nil, pattern.ID)
		assert.Equal(t, name, pattern.Name)
		assert.Equal(t, description, pattern.Description)
		assert.Equal(t, threat.ID, pattern.ThreatID)
		assert.Equal(t, isActive, pattern.IsActive)

		// Verify pattern was saved to database
		db := database.GetDB()
		var dbPattern models.ThreatPattern
		err = db.First(&dbPattern, "id = ?", pattern.ID).Error
		require.NoError(t, err)
		assert.Equal(t, pattern.ID, dbPattern.ID)
		assert.Equal(t, name, dbPattern.Name)
	})

	t.Run("InvalidThreatID", func(t *testing.T) {
		// Try to create pattern with non-existent threat ID
		nonExistentThreatID := uuid.New()
		pattern, err := CreateThreatPattern("Test Pattern", "Description", nonExistentThreatID, true)

		assert.Error(t, err)
		assert.Nil(t, pattern)
		assert.Contains(t, err.Error(), "threat not found")
	})
}

func TestThreatPatternService_GetThreatPattern(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := service.CreateThreat("Get Test Threat", "Threat for get test")
		require.NoError(t, err)

		createdPattern, err := CreateThreatPattern("Get Test Pattern", "Pattern for get test", threat.ID, true)
		require.NoError(t, err)

		// Get the pattern
		retrievedPattern, err := GetThreatPattern(createdPattern.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedPattern)
		assert.Equal(t, createdPattern.ID, retrievedPattern.ID)
		assert.Equal(t, createdPattern.Name, retrievedPattern.Name)
		assert.Equal(t, createdPattern.Description, retrievedPattern.Description)
		assert.Equal(t, createdPattern.ThreatID, retrievedPattern.ThreatID)
		assert.Equal(t, createdPattern.IsActive, retrievedPattern.IsActive)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to get non-existent pattern
		nonExistentID := uuid.New()
		pattern, err := GetThreatPattern(nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, pattern)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestThreatPatternService_UpdateThreatPattern(t *testing.T) {

	t.Run("FullUpdate", func(t *testing.T) {
		// Create test data
		threat1, err := service.CreateThreat("Update Test Threat 1", "First threat")
		require.NoError(t, err)

		createdPattern, err := CreateThreatPattern("Original Pattern", "Original description", threat1.ID, true)
		require.NoError(t, err)

		// Update the pattern (without changing threat ID for now)
		newName := "Updated Pattern"
		newDescription := "Updated description"
		newIsActive := false
		updatedPattern, err := UpdateThreatPattern(createdPattern.ID, &newName, &newDescription, nil, &newIsActive)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedPattern)
		assert.Equal(t, createdPattern.ID, updatedPattern.ID)
		assert.Equal(t, newName, updatedPattern.Name)
		assert.Equal(t, newDescription, updatedPattern.Description)
		assert.Equal(t, threat1.ID, updatedPattern.ThreatID) // Should remain threat1.ID since we didn't update it
		assert.Equal(t, newIsActive, updatedPattern.IsActive)

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbPattern models.ThreatPattern
		err = db.First(&dbPattern, "id = ?", createdPattern.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbPattern.Name)
		assert.Equal(t, newDescription, dbPattern.Description)
		assert.Equal(t, threat1.ID, dbPattern.ThreatID) // Should remain threat1.ID since we didn't update it
		assert.Equal(t, newIsActive, dbPattern.IsActive)
	})

	t.Run("PartialUpdate", func(t *testing.T) {
		// Create test data
		threat, err := service.CreateThreat("Partial Test Threat", "Threat for partial update")
		require.NoError(t, err)

		originalName := "Partial Pattern"
		originalDescription := "Original description"
		originalIsActive := true
		createdPattern, err := CreateThreatPattern(originalName, originalDescription, threat.ID, originalIsActive)
		require.NoError(t, err)

		// Update only the name
		newName := "New Name Only"
		updatedPattern, err := UpdateThreatPattern(createdPattern.ID, &newName, nil, nil, nil)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, newName, updatedPattern.Name)
		assert.Equal(t, originalDescription, updatedPattern.Description)
		assert.Equal(t, threat.ID, updatedPattern.ThreatID)
		assert.Equal(t, originalIsActive, updatedPattern.IsActive)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to update non-existent pattern
		nonExistentID := uuid.New()
		newName := "New Name"
		pattern, err := UpdateThreatPattern(nonExistentID, &newName, nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, pattern)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("InvalidThreatID", func(t *testing.T) {
		// Create test pattern
		threat, err := service.CreateThreat("Invalid Update Threat", "Threat for invalid update test")
		require.NoError(t, err)

		createdPattern, err := CreateThreatPattern("Invalid Update Pattern", "Pattern for invalid update", threat.ID, true)
		require.NoError(t, err)

		// Try to update with non-existent threat ID
		nonExistentThreatID := uuid.New()
		pattern, err := UpdateThreatPattern(createdPattern.ID, nil, nil, &nonExistentThreatID, nil)

		assert.Error(t, err)
		assert.Nil(t, pattern)
		assert.Contains(t, err.Error(), "threat not found")
	})
}

func TestThreatPatternService_DeleteThreatPattern(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := service.CreateThreat("Delete Test Threat", "Threat for delete test")
		require.NoError(t, err)

		createdPattern, err := CreateThreatPattern("Delete Test Pattern", "Pattern to be deleted", threat.ID, true)
		require.NoError(t, err)

		// Delete the pattern
		err = DeleteThreatPattern(createdPattern.ID)
		require.NoError(t, err)

		// Verify deletion
		db := database.GetDB()
		var dbPattern models.ThreatPattern
		err = db.First(&dbPattern, "id = ?", createdPattern.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to delete non-existent pattern
		nonExistentID := uuid.New()
		err := DeleteThreatPattern(nonExistentID)

		// Should succeed (GORM behavior)
		assert.NoError(t, err)
	})
}

func TestThreatPatternService_ListThreatPatterns(t *testing.T) {

	t.Run("WithPatterns", func(t *testing.T) {
		// Create test data
		threat, err := service.CreateThreat("List Test Threat", "Threat for list test")
		require.NoError(t, err)

		patterns := []struct {
			name     string
			isActive bool
		}{
			{"Pattern 1", true},
			{"Pattern 2", false},
			{"Pattern 3", true},
		}

		var createdPatterns []*models.ThreatPattern
		for _, p := range patterns {
			pattern, err := CreateThreatPattern(p.name, "Description", threat.ID, p.isActive)
			require.NoError(t, err)
			createdPatterns = append(createdPatterns, pattern)
		}

		// List all patterns
		retrievedPatterns, err := ListThreatPatterns()

		// Assertions
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(retrievedPatterns), len(patterns))

		// Verify created patterns are in the list
		patternMap := make(map[uuid.UUID]models.ThreatPattern)
		for _, p := range retrievedPatterns {
			patternMap[p.ID] = p
		}

		for _, created := range createdPatterns {
			retrieved, exists := patternMap[created.ID]
			assert.True(t, exists, "Created pattern should exist in list")
			assert.Equal(t, created.Name, retrieved.Name)
			assert.Equal(t, created.IsActive, retrieved.IsActive)
		}
	})

	t.Run("ActiveOnly", func(t *testing.T) {
		// Create test data with mixed active status
		threat, err := service.CreateThreat("Active List Test Threat", "Threat for active list test")
		require.NoError(t, err)

		activePattern, err := CreateThreatPattern("Active Pattern", "Active pattern", threat.ID, true)
		require.NoError(t, err)

		inactivePattern, err := CreateThreatPattern("Inactive Pattern", "Inactive pattern", threat.ID, false)
		require.NoError(t, err)

		// List only active patterns
		activePatterns, err := ListActiveThreatPatterns()

		// Assertions
		require.NoError(t, err)

		// Check that our specific active pattern is included and our specific inactive is not
		activeFound := false
		inactiveFound := false
		for _, p := range activePatterns {
			if p.ID == activePattern.ID {
				activeFound = true
				assert.True(t, p.IsActive, "Pattern in active list should have IsActive=true")
			}
			if p.ID == inactivePattern.ID {
				inactiveFound = true
			}
		}
		assert.True(t, activeFound, "Active pattern should be in active list")
		assert.False(t, inactiveFound, "Inactive pattern should not be in active list")
	})
}

func TestThreatPatternService_CreateThreatPatternWithConditions(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Create a threat first
		threat, err := service.CreateThreat("Pattern with Conditions", "Threat for pattern with conditions")
		require.NoError(t, err)

		// Create conditions
		conditions := []models.PatternCondition{
			{
				ConditionType: models.ConditionTypeTag.String(),
				Operator:      models.OperatorContains.String(),
				Value:         "test-tag",
			},
			{
				ConditionType:    models.ConditionTypeRelationshipTargetTag.String(),
				Operator:         models.OperatorEquals.String(),
				Value:            "test-product",
				RelationshipType: "RELATED_TO",
			},
		}

		// Create pattern with conditions
		pattern, err := CreateThreatPatternWithConditions(
			"Pattern with Conditions",
			"A pattern with multiple conditions",
			threat.ID,
			true,
			conditions,
		)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, pattern)
		assert.Equal(t, "Pattern with Conditions", pattern.Name)
		assert.Equal(t, threat.ID, pattern.ThreatID)
		assert.True(t, pattern.IsActive)
		assert.Len(t, pattern.Conditions, 2)

		// Verify conditions were created with correct pattern ID
		for _, condition := range pattern.Conditions {
			assert.Equal(t, pattern.ID, condition.PatternID)
		}
	})

	t.Run("InvalidThreatID", func(t *testing.T) {
		// Try to create pattern with non-existent threat ID
		nonExistentThreatID := uuid.New()
		conditions := []models.PatternCondition{
			{
				ConditionType: models.ConditionTypeTag.String(),
				Operator:      models.OperatorContains.String(),
				Value:         "test-tag",
			},
		}

		pattern, err := CreateThreatPatternWithConditions(
			"Invalid Pattern",
			"Pattern with invalid threat ID",
			nonExistentThreatID,
			true,
			conditions,
		)

		assert.Error(t, err)
		assert.Nil(t, pattern)
		assert.Contains(t, err.Error(), "threat not found")
	})
}
