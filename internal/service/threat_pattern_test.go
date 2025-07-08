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

func TestThreatPatternService_CreateThreatPattern(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create a threat first
		threat, err := CreateThreat("Test Threat", "A test threat for pattern testing")
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
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("Get Test Threat", "Threat for get test")
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
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("FullUpdate", func(t *testing.T) {
		// Create test data
		threat1, err := CreateThreat("Update Test Threat 1", "First threat")
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
		threat, err := CreateThreat("Partial Test Threat", "Threat for partial update")
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
		threat, err := CreateThreat("Invalid Update Threat", "Threat for invalid update test")
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
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("Delete Test Threat", "Threat for delete test")
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
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("WithPatterns", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("List Test Threat", "Threat for list test")
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
		threat, err := CreateThreat("Active List Test Threat", "Threat for active list test")
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

func TestThreatPatternService_CreatePatternCondition(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("Condition Test Threat", "Threat for condition test")
		require.NoError(t, err)

		pattern, err := CreateThreatPattern("Condition Test Pattern", "Pattern for condition test", threat.ID, true)
		require.NoError(t, err)

		// Create pattern condition
		conditionType := models.ConditionTypeTag.String()
		operator := models.OperatorContains.String()
		value := "test-tag"
		relationshipType := ""

		condition, err := CreatePatternCondition(pattern.ID, conditionType, operator, value, relationshipType)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, condition)
		assert.NotEqual(t, uuid.Nil, condition.ID)
		assert.Equal(t, pattern.ID, condition.PatternID)
		assert.Equal(t, conditionType, condition.ConditionType)
		assert.Equal(t, operator, condition.Operator)
		assert.Equal(t, value, condition.Value)
		assert.Equal(t, relationshipType, condition.RelationshipType)
	})

	t.Run("RelationshipCondition", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("Relationship Condition Threat", "Threat for relationship condition")
		require.NoError(t, err)

		pattern, err := CreateThreatPattern("Relationship Pattern", "Pattern for relationship test", threat.ID, true)
		require.NoError(t, err)

		// Create relationship-based condition
		condition, err := CreatePatternCondition(
			pattern.ID,
			models.ConditionTypeRelationshipTargetTag.String(),
			models.OperatorHasRelationshipWith.String(),
			"privileged",
			"connects_to",
		)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, models.ConditionTypeRelationshipTargetTag.String(), condition.ConditionType)
		assert.Equal(t, models.OperatorHasRelationshipWith.String(), condition.Operator)
		assert.Equal(t, "privileged", condition.Value)
		assert.Equal(t, "connects_to", condition.RelationshipType)
	})

	t.Run("InvalidPattern", func(t *testing.T) {
		// Try to create condition for non-existent pattern
		nonExistentPatternID := uuid.New()
		condition, err := CreatePatternCondition(nonExistentPatternID, models.ConditionTypeTag.String(), models.OperatorContains.String(), "test", "")

		assert.Error(t, err)
		assert.Nil(t, condition)
		assert.Contains(t, err.Error(), "threat pattern not found")
	})

	t.Run("InvalidCondition", func(t *testing.T) {
		// Create test pattern
		threat, err := CreateThreat("Invalid Condition Threat", "Threat for invalid condition")
		require.NoError(t, err)

		pattern, err := CreateThreatPattern("Invalid Condition Pattern", "Pattern for invalid condition", threat.ID, true)
		require.NoError(t, err)

		// Try to create invalid condition (missing relationship_type for relationship condition)
		condition, err := CreatePatternCondition(pattern.ID, models.ConditionTypeRelationshipTargetTag.String(), models.OperatorHasRelationshipWith.String(), "test", "")

		assert.Error(t, err)
		assert.Nil(t, condition)
		assert.Contains(t, err.Error(), "relationship_type is required")
	})
}

func TestThreatPatternService_UpdatePatternCondition(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("Update Condition Threat", "Threat for update condition test")
		require.NoError(t, err)

		pattern, err := CreateThreatPattern("Update Condition Pattern", "Pattern for update condition test", threat.ID, true)
		require.NoError(t, err)

		createdCondition, err := CreatePatternCondition(pattern.ID, models.ConditionTypeTag.String(), models.OperatorContains.String(), "original-tag", "")
		require.NoError(t, err)

		// Update the condition
		newConditionType := models.ConditionTypeProductTag.String()
		newOperator := models.OperatorNotContains.String()
		newValue := "updated-tag"
		newRelationshipType := ""

		updatedCondition, err := UpdatePatternCondition(createdCondition.ID, &newConditionType, &newOperator, &newValue, &newRelationshipType)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, createdCondition.ID, updatedCondition.ID)
		assert.Equal(t, newConditionType, updatedCondition.ConditionType)
		assert.Equal(t, newOperator, updatedCondition.Operator)
		assert.Equal(t, newValue, updatedCondition.Value)
		assert.Equal(t, newRelationshipType, updatedCondition.RelationshipType)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to update non-existent condition
		nonExistentID := uuid.New()
		newValue := "new-value"
		condition, err := UpdatePatternCondition(nonExistentID, nil, nil, &newValue, nil)

		assert.Error(t, err)
		assert.Nil(t, condition)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestThreatPatternService_CreateThreatPatternWithConditions(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test threat
		threat, err := CreateThreat("Complex Pattern Threat", "Threat for complex pattern test")
		require.NoError(t, err)

		// Define conditions
		conditions := []models.PatternCondition{
			{
				ConditionType: models.ConditionTypeTag.String(),
				Operator:      models.OperatorContains.String(),
				Value:         "internet-facing",
			},
			{
				ConditionType:    models.ConditionTypeRelationshipTargetTag.String(),
				Operator:         models.OperatorHasRelationshipWith.String(),
				Value:            "database",
				RelationshipType: "connects_to",
			},
			{
				ConditionType: models.ConditionTypeProductTag.String(),
				Operator:      models.OperatorNotContains.String(),
				Value:         "high-security",
			},
		}

		// Create pattern with conditions
		pattern, err := CreateThreatPatternWithConditions(
			"Complex Pattern",
			"A complex pattern with multiple conditions",
			threat.ID,
			true,
			conditions,
		)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, pattern)
		assert.Equal(t, "Complex Pattern", pattern.Name)
		assert.Equal(t, threat.ID, pattern.ThreatID)
		assert.True(t, pattern.IsActive)
		assert.Len(t, pattern.Conditions, 3)

		// Verify conditions
		for i, condition := range pattern.Conditions {
			assert.Equal(t, pattern.ID, condition.PatternID)
			assert.Equal(t, conditions[i].ConditionType, condition.ConditionType)
			assert.Equal(t, conditions[i].Operator, condition.Operator)
			assert.Equal(t, conditions[i].Value, condition.Value)
			assert.Equal(t, conditions[i].RelationshipType, condition.RelationshipType)
		}
	})

	t.Run("InvalidCondition", func(t *testing.T) {
		// Create test threat
		threat, err := CreateThreat("Invalid Complex Threat", "Threat for invalid complex pattern")
		require.NoError(t, err)

		// Define conditions with one invalid condition
		conditions := []models.PatternCondition{
			{
				ConditionType: models.ConditionTypeTag.String(),
				Operator:      models.OperatorContains.String(),
				Value:         "valid-tag",
			},
			{
				ConditionType:    models.ConditionTypeRelationshipTargetTag.String(),
				Operator:         models.OperatorHasRelationshipWith.String(),
				Value:            "test-tag",
				RelationshipType: "", // Missing required relationship_type
			},
		}

		// Try to create pattern
		pattern, err := CreateThreatPatternWithConditions(
			"Invalid Complex Pattern",
			"A pattern with invalid conditions",
			threat.ID,
			true,
			conditions,
		)

		assert.Error(t, err)
		assert.Nil(t, pattern)
		assert.Contains(t, err.Error(), "relationship_type is required")
	})
}

func TestThreatPatternService_ListPatternConditions(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("List Conditions Threat", "Threat for list conditions test")
		require.NoError(t, err)

		pattern, err := CreateThreatPattern("List Conditions Pattern", "Pattern for list conditions test", threat.ID, true)
		require.NoError(t, err)

		// Create multiple conditions
		conditions := []struct {
			conditionType string
			operator      string
			value         string
		}{
			{models.ConditionTypeTag.String(), models.OperatorContains.String(), "tag1"},
			{models.ConditionTypeProduct.String(), models.OperatorEquals.String(), "product1"},
			{models.ConditionTypeProductTag.String(), models.OperatorNotContains.String(), "tag2"},
		}

		var createdConditions []*models.PatternCondition
		for _, c := range conditions {
			condition, err := CreatePatternCondition(pattern.ID, c.conditionType, c.operator, c.value, "")
			require.NoError(t, err)
			createdConditions = append(createdConditions, condition)
		}

		// List conditions for pattern
		retrievedConditions, err := ListPatternConditionsByPatternID(pattern.ID)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, retrievedConditions, len(conditions))

		// Verify all conditions are returned
		conditionMap := make(map[uuid.UUID]models.PatternCondition)
		for _, c := range retrievedConditions {
			conditionMap[c.ID] = c
		}

		for _, created := range createdConditions {
			retrieved, exists := conditionMap[created.ID]
			assert.True(t, exists, "Created condition should exist in list")
			assert.Equal(t, created.ConditionType, retrieved.ConditionType)
			assert.Equal(t, created.Operator, retrieved.Operator)
			assert.Equal(t, created.Value, retrieved.Value)
		}
	})
}

func TestThreatPatternService_DeletePatternCondition(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		threat, err := CreateThreat("Delete Condition Threat", "Threat for delete condition test")
		require.NoError(t, err)

		pattern, err := CreateThreatPattern("Delete Condition Pattern", "Pattern for delete condition test", threat.ID, true)
		require.NoError(t, err)

		condition, err := CreatePatternCondition(pattern.ID, models.ConditionTypeTag.String(), models.OperatorContains.String(), "delete-test", "")
		require.NoError(t, err)

		// Delete the condition
		err = DeletePatternCondition(condition.ID)
		require.NoError(t, err)

		// Verify deletion
		db := database.GetDB()
		var dbCondition models.PatternCondition
		err = db.First(&dbCondition, "id = ?", condition.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestThreatPatternService_GetInstanceThreatsByThreatPattern(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
		&models.Relationship{},
	)
	defer cleanup()

	t.Run("ComplexPatternMatching", func(t *testing.T) {
		// Create products
		webProduct, err := CreateProduct("Web Application", "A web application")
		require.NoError(t, err)

		dbProduct, err := CreateProduct("Database", "A database system")
		require.NoError(t, err)

		apiProduct, err := CreateProduct("API Service", "An API service")
		require.NoError(t, err)

		// Create instances
		webServer, err := CreateInstance("Web Server", webProduct.ID)
		require.NoError(t, err)

		database, err := CreateInstance("Database", dbProduct.ID)
		require.NoError(t, err)

		apiServer, err := CreateInstance("API Server", apiProduct.ID)
		require.NoError(t, err)

		adminPanel, err := CreateInstance("Admin Panel", webProduct.ID)
		require.NoError(t, err)

		// Create tags
		internetFacingTag, err := CreateTag("internet-facing", "Internet facing component", "#FF0000")
		require.NoError(t, err)

		databaseTag, err := CreateTag("database", "Database component", "#00FF00")
		require.NoError(t, err)

		privilegedTag, err := CreateTag("privileged", "Privileged component", "#0000FF")
		require.NoError(t, err)

		highSecurityTag, err := CreateTag("high-security", "High security component", "#FFA500")
		require.NoError(t, err)

		// Assign tags to instances
		err = AssignTagToInstance(internetFacingTag.ID, webServer.ID)
		require.NoError(t, err)

		err = AssignTagToInstance(databaseTag.ID, database.ID)
		require.NoError(t, err)

		err = AssignTagToInstance(privilegedTag.ID, adminPanel.ID)
		require.NoError(t, err)

		// Assign tags to products
		err = AssignTagToProduct(highSecurityTag.ID, dbProduct.ID)
		require.NoError(t, err)

		// Create relationships
		err = AddRelationship(&webServer.ID, nil, &database.ID, nil, "connects_to", "")
		require.NoError(t, err)

		err = AddRelationship(&adminPanel.ID, nil, &database.ID, nil, "reads_from", "")
		require.NoError(t, err)

		err = AddRelationship(&apiServer.ID, nil, &database.ID, nil, "connects_to", "")
		require.NoError(t, err)

		// Create threats
		sqlInjectionThreat, err := CreateThreat("SQL Injection", "SQL injection vulnerability")
		require.NoError(t, err)

		privilegeEscalationThreat, err := CreateThreat("Privilege Escalation", "Privilege escalation threat")
		require.NoError(t, err)

		dataExposureThreat, err := CreateThreat("Data Exposure", "Data exposure threat")
		require.NoError(t, err)

		// Create threat patterns

		// Pattern 1: Internet-facing instances that connect to database
		pattern1, err := CreateThreatPatternWithConditions(
			"Internet-facing DB Connection",
			"Internet-facing components that connect to database",
			sqlInjectionThreat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "internet-facing",
				},
				{
					ConditionType:    models.ConditionTypeRelationshipTargetTag.String(),
					Operator:         models.OperatorHasRelationshipWith.String(),
					Value:            "database",
					RelationshipType: "connects_to",
				},
			},
		)
		require.NoError(t, err)

		// Pattern 2: Privileged instances
		pattern2, err := CreateThreatPatternWithConditions(
			"Privileged Components",
			"Components with privileged access",
			privilegeEscalationThreat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "privileged",
				},
			},
		)
		require.NoError(t, err)

		// Pattern 3: Instances of products with high-security tag
		pattern3, err := CreateThreatPatternWithConditions(
			"High Security Products",
			"Instances of products with high security requirements",
			dataExposureThreat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeProductTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "high-security",
				},
			},
		)
		require.NoError(t, err)

		// Pattern 4: Specific product by name
		pattern4, err := CreateThreatPatternWithConditions(
			"Web Application Product",
			"Instances of web application product",
			dataExposureThreat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeProduct.String(),
					Operator:      models.OperatorEquals.String(),
					Value:         "Web Application",
				},
			},
		)
		require.NoError(t, err)

		// Pattern 5: Specific instance by relationship to specific target ID
		pattern5, err := CreateThreatPatternWithConditions(
			"Direct Database Connection",
			"Instances that connect directly to the database instance",
			sqlInjectionThreat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType:    models.ConditionTypeRelationshipTargetID.String(),
					Operator:         models.OperatorHasRelationshipWith.String(),
					Value:            database.ID.String(),
					RelationshipType: "connects_to",
				},
			},
		)
		require.NoError(t, err)

		// Execute pattern matching
		matches, err := GetInstanceThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify matches

		// Web Server should match pattern1 (internet-facing + connects to database) and pattern4 (web application product) and pattern5 (connects to database)
		webServerMatches, exists := matches[webServer.ID]
		require.True(t, exists, "Web Server should have matches")
		require.Len(t, webServerMatches, 3, "Web Server should have 3 matches")

		matchedPatterns := make(map[uuid.UUID]bool)
		for _, match := range webServerMatches {
			matchedPatterns[match.PatternID] = true
		}
		assert.True(t, matchedPatterns[pattern1.ID], "Web Server should match pattern1")
		assert.True(t, matchedPatterns[pattern4.ID], "Web Server should match pattern4")
		assert.True(t, matchedPatterns[pattern5.ID], "Web Server should match pattern5")

		// Database should match pattern3 (high-security product)
		databaseMatches, exists := matches[database.ID]
		require.True(t, exists, "Database should have matches")
		require.Len(t, databaseMatches, 1, "Database should have 1 match")
		assert.Equal(t, pattern3.ID, databaseMatches[0].PatternID, "Database should match pattern3")

		// Admin Panel should match pattern2 (privileged) and pattern4 (web application product)
		adminMatches, exists := matches[adminPanel.ID]
		require.True(t, exists, "Admin Panel should have matches")
		require.Len(t, adminMatches, 2, "Admin Panel should have 2 matches")

		adminMatchedPatterns := make(map[uuid.UUID]bool)
		for _, match := range adminMatches {
			adminMatchedPatterns[match.PatternID] = true
		}
		assert.True(t, adminMatchedPatterns[pattern2.ID], "Admin Panel should match pattern2")
		assert.True(t, adminMatchedPatterns[pattern4.ID], "Admin Panel should match pattern4")

		// API Server should match pattern5 (connects to database)
		apiMatches, exists := matches[apiServer.ID]
		require.True(t, exists, "API Server should have matches")
		require.Len(t, apiMatches, 1, "API Server should have 1 match")
		assert.Equal(t, pattern5.ID, apiMatches[0].PatternID, "API Server should match pattern5")
	})

	t.Run("NoMatches", func(t *testing.T) {
		// Create a product and instance with no matching patterns
		isolatedProduct, err := CreateProduct("Isolated Service", "An isolated service")
		require.NoError(t, err)

		isolatedInstance, err := CreateInstance("Isolated Instance", isolatedProduct.ID)
		require.NoError(t, err)

		// Create a threat and pattern that won't match
		threat, err := CreateThreat("Isolated Threat", "A threat for isolated components")
		require.NoError(t, err)

		_, err = CreateThreatPatternWithConditions(
			"Non-matching Pattern",
			"A pattern that won't match anything",
			threat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "non-existent-tag",
				},
			},
		)
		require.NoError(t, err)

		// Execute pattern matching
		matches, err := GetInstanceThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify no matches for isolated instance
		_, exists := matches[isolatedInstance.ID]
		assert.False(t, exists, "Isolated instance should have no matches")
	})

	t.Run("InactivePattern", func(t *testing.T) {
		// Create a product and instance
		testProduct, err := CreateProduct("Test Product", "A test product")
		require.NoError(t, err)

		testInstance, err := CreateInstance("Test Instance", testProduct.ID)
		require.NoError(t, err)

		// Create a tag and assign it to the instance
		testTag, err := CreateTag("test-tag", "Test tag", "#FFFFFF")
		require.NoError(t, err)

		err = AssignTagToInstance(testTag.ID, testInstance.ID)
		require.NoError(t, err)

		// Create a threat and inactive pattern
		threat, err := CreateThreat("Test Threat", "A test threat")
		require.NoError(t, err)

		_, err = CreateThreatPatternWithConditions(
			"Inactive Pattern",
			"An inactive pattern",
			threat.ID,
			false, // Inactive
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "test-tag",
				},
			},
		)
		require.NoError(t, err)

		// Execute pattern matching
		matches, err := GetInstanceThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify no matches for instance with inactive pattern
		_, exists := matches[testInstance.ID]
		assert.False(t, exists, "Instance should have no matches from inactive pattern")
	})

	t.Run("TagExistsCondition", func(t *testing.T) {
		// Create product and instances
		product, err := CreateProduct("Tag Exists Product", "Product for testing tag existence")
		require.NoError(t, err)

		taggedInstance, err := CreateInstance("Tagged Instance", product.ID)
		require.NoError(t, err)

		untaggedInstance, err := CreateInstance("Untagged Instance", product.ID)
		require.NoError(t, err)

		// Create and assign a tag to one instance
		tag, err := CreateTag("any-tag", "Any tag", "#000000")
		require.NoError(t, err)

		err = AssignTagToInstance(tag.ID, taggedInstance.ID)
		require.NoError(t, err)

		// Create threat and pattern that matches instances with any tag
		threat, err := CreateThreat("Tagged Threat", "Threat for tagged instances")
		require.NoError(t, err)

		_, err = CreateThreatPatternWithConditions(
			"Any Tag Pattern",
			"Pattern that matches instances with any tag",
			threat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorExists.String(),
					Value:         "", // Value not used for EXISTS
				},
			},
		)
		require.NoError(t, err)

		// Execute pattern matching
		matches, err := GetInstanceThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify tagged instance matches but untagged doesn't
		_, taggedExists := matches[taggedInstance.ID]
		assert.True(t, taggedExists, "Tagged instance should match")

		_, untaggedExists := matches[untaggedInstance.ID]
		assert.False(t, untaggedExists, "Untagged instance should not match")
	})

	t.Run("RelationshipExistsCondition", func(t *testing.T) {
		// Create products and instances
		product1, err := CreateProduct("Product 1", "First product")
		require.NoError(t, err)

		product2, err := CreateProduct("Product 2", "Second product")
		require.NoError(t, err)

		instance1, err := CreateInstance("Instance 1", product1.ID)
		require.NoError(t, err)

		instance2, err := CreateInstance("Instance 2", product2.ID)
		require.NoError(t, err)

		isolatedInstance, err := CreateInstance("Isolated Instance", product1.ID)
		require.NoError(t, err)

		// Create relationship between instance1 and instance2
		err = AddRelationship(&instance1.ID, nil, &instance2.ID, nil, "depends_on", "")
		require.NoError(t, err)

		// Create threat and pattern that matches instances with depends_on relationship
		threat, err := CreateThreat("Dependency Threat", "Threat for instances with dependencies")
		require.NoError(t, err)

		_, err = CreateThreatPatternWithConditions(
			"Dependency Pattern",
			"Pattern that matches instances with depends_on relationship",
			threat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType:    models.ConditionTypeRelationship.String(),
					Operator:         models.OperatorExists.String(),
					Value:            "", // Value not used for EXISTS
					RelationshipType: "depends_on",
				},
			},
		)
		require.NoError(t, err)

		// Execute pattern matching
		matches, err := GetInstanceThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify instance1 matches but isolated instance doesn't
		_, instance1Exists := matches[instance1.ID]
		assert.True(t, instance1Exists, "Instance1 should match (has depends_on relationship)")

		_, isolatedExists := matches[isolatedInstance.ID]
		assert.False(t, isolatedExists, "Isolated instance should not match (no depends_on relationship)")
	})
}
