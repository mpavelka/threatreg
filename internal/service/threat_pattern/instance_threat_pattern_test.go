package threat_pattern

import (
	"testing"
	"threatreg/internal/models"
	"threatreg/internal/service"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThreatPatternService_GetComponentThreatsByThreatPattern(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.Tag{},
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
		&models.Relationship{},
	)
	defer cleanup()

	t.Run("ComplexPatternMatching", func(t *testing.T) {
		// Create products
		dbProduct, err := service.CreateComponent("Database", "A database system", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create instances
		webServer, err := service.CreateComponent("Web Server", "A web server", models.ComponentTypeProduct)
		require.NoError(t, err)

		database, err := service.CreateComponent("Database", "A database", models.ComponentTypeProduct)
		require.NoError(t, err)

		apiServer, err := service.CreateComponent("API Server", "An API server", models.ComponentTypeProduct)
		require.NoError(t, err)

		adminPanel, err := service.CreateComponent("Admin Panel", "An admin panel", models.ComponentTypeProduct)
		require.NoError(t, err)

		// Create tags
		internetFacingTag, err := service.CreateTag("internet-facing", "Internet facing component", "#FF0000")
		require.NoError(t, err)

		databaseTag, err := service.CreateTag("database", "Database component", "#00FF00")
		require.NoError(t, err)

		privilegedTag, err := service.CreateTag("privileged", "Privileged component", "#0000FF")
		require.NoError(t, err)

		highSecurityTag, err := service.CreateTag("high-security", "High security component", "#FFA500")
		require.NoError(t, err)

		// Assign tags to instances
		err = service.AssignTagToComponent(internetFacingTag.ID, webServer.ID)
		require.NoError(t, err)

		err = service.AssignTagToComponent(databaseTag.ID, database.ID)
		require.NoError(t, err)

		err = service.AssignTagToComponent(privilegedTag.ID, adminPanel.ID)
		require.NoError(t, err)

		// Assign tags to products
		err = service.AssignTagToComponent(highSecurityTag.ID, dbProduct.ID)
		require.NoError(t, err)

		// Create relationships
		err = service.AddRelationship(webServer.ID, database.ID, "connects_to", "")
		require.NoError(t, err)

		err = service.AddRelationship(adminPanel.ID, database.ID, "reads_from", "")
		require.NoError(t, err)

		err = service.AddRelationship(apiServer.ID, database.ID, "connects_to", "")
		require.NoError(t, err)

		// Create threats
		sqlInjectionThreat, err := service.CreateThreat("SQL Injection", "SQL injection vulnerability")
		require.NoError(t, err)

		privilegeEscalationThreat, err := service.CreateThreat("Privilege Escalation", "Privilege escalation threat")
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

		// Pattern 5: Specific instance by relationship to specific target ID
		pattern5, err := CreateThreatPatternWithConditions(
			"Direct Database Connection",
			"Components that connect directly to the database instance",
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
		matches, err := GetAllComponentsThreatsByThreatPattern()
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
		assert.True(t, matchedPatterns[pattern5.ID], "Web Server should match pattern5")

		// Database should match pattern3 (high-security product)
		databaseMatches, exists := matches[database.ID]
		require.True(t, exists, "Database should have matches")
		require.Len(t, databaseMatches, 1, "Database should have 1 match")

		// Admin Panel should match pattern2 (privileged) and pattern4 (web application product)
		adminMatches, exists := matches[adminPanel.ID]
		require.True(t, exists, "Admin Panel should have matches")
		require.Len(t, adminMatches, 2, "Admin Panel should have 2 matches")

		adminMatchedPatterns := make(map[uuid.UUID]bool)
		for _, match := range adminMatches {
			adminMatchedPatterns[match.PatternID] = true
		}
		assert.True(t, adminMatchedPatterns[pattern2.ID], "Admin Panel should match pattern2")

		// API Server should match pattern5 (connects to database)
		apiMatches, exists := matches[apiServer.ID]
		require.True(t, exists, "API Server should have matches")
		require.Len(t, apiMatches, 1, "API Server should have 1 match")
		assert.Equal(t, pattern5.ID, apiMatches[0].PatternID, "API Server should match pattern5")
	})

	t.Run("NoMatches", func(t *testing.T) {

		isolatedComponent, err := service.CreateComponent("Isolated Component", "An isolated component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create a threat and pattern that won't match
		threat, err := service.CreateThreat("Isolated Threat", "A threat for isolated components")
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
		matches, err := GetAllComponentsThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify no matches for isolated instance
		_, exists := matches[isolatedComponent.ID]
		assert.False(t, exists, "Isolated instance should have no matches")
	})

	t.Run("InactivePattern", func(t *testing.T) {

		testComponent, err := service.CreateComponent("Test Component", "A test component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create a tag and assign it to the instance
		testTag, err := service.CreateTag("test-tag", "Test tag", "#FFFFFF")
		require.NoError(t, err)

		err = service.AssignTagToComponent(testTag.ID, testComponent.ID)
		require.NoError(t, err)

		// Create a threat and inactive pattern
		threat, err := service.CreateThreat("Test Threat", "A test threat")
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
		matches, err := GetAllComponentsThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify no matches for instance with inactive pattern
		_, exists := matches[testComponent.ID]
		assert.False(t, exists, "Component should have no matches from inactive pattern")
	})

	t.Run("TagExistsCondition", func(t *testing.T) {

		taggedComponent, err := service.CreateComponent("Tagged Component", "A tagged component", models.ComponentTypeInstance)
		require.NoError(t, err)

		untaggedComponent, err := service.CreateComponent("Untagged Component", "An untagged component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create and assign a tag to one instance
		tag, err := service.CreateTag("any-tag", "Any tag", "#000000")
		require.NoError(t, err)

		err = service.AssignTagToComponent(tag.ID, taggedComponent.ID)
		require.NoError(t, err)

		// Create threat and pattern that matches instances with any tag
		threat, err := service.CreateThreat("Tagged Threat", "Threat for tagged instances")
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
		matches, err := GetAllComponentsThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify tagged instance matches but untagged doesn't
		_, taggedExists := matches[taggedComponent.ID]
		assert.True(t, taggedExists, "Tagged instance should match")

		_, untaggedExists := matches[untaggedComponent.ID]
		assert.False(t, untaggedExists, "Untagged instance should not match")
	})

	t.Run("RelationshipExistsCondition", func(t *testing.T) {

		instance1, err := service.CreateComponent("Component 1", "First component", models.ComponentTypeInstance)
		require.NoError(t, err)

		instance2, err := service.CreateComponent("Component 2", "Second component", models.ComponentTypeInstance)
		require.NoError(t, err)

		isolatedComponent, err := service.CreateComponent("Isolated Component", "An isolated component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create relationship between instance1 and instance2
		err = service.AddRelationship(instance1.ID, instance2.ID, "depends_on", "")
		require.NoError(t, err)

		// Create threat and pattern that matches instances with depends_on relationship
		threat, err := service.CreateThreat("Dependency Threat", "Threat for instances with dependencies")
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
		matches, err := GetAllComponentsThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify instance1 matches but isolated instance doesn't
		_, instance1Exists := matches[instance1.ID]
		assert.True(t, instance1Exists, "Component1 should match (has depends_on relationship)")

		_, isolatedExists := matches[isolatedComponent.ID]
		assert.False(t, isolatedExists, "Isolated instance should not match (no depends_on relationship)")
	})
}

func TestGetComponentThreatsByThreatPattern(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.Tag{},
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("ComponentMatchesPattern", func(t *testing.T) {

		instance, err := service.CreateComponent("Web Server", "A web server", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create a tag and assign it to the instance
		tag, err := service.CreateTag("internet-facing", "Internet facing component", "#FF0000")
		require.NoError(t, err)

		err = service.AssignTagToComponent(tag.ID, instance.ID)
		require.NoError(t, err)

		// Create a threat
		threat, err := service.CreateThreat("SQL Injection", "SQL injection vulnerability")
		require.NoError(t, err)

		// Create a threat pattern that should match the instance
		pattern, err := CreateThreatPatternWithConditions(
			"Internet-facing Components",
			"Components facing the internet",
			threat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "internet-facing",
				},
			},
		)
		require.NoError(t, err)

		// Test the function
		matches, err := GetComponentThreatsByThreatPattern(*instance, *pattern)
		require.NoError(t, err)

		// Verify the match
		require.Len(t, matches, 1, "Should have exactly one match")
		assert.Equal(t, instance.ID, matches[0].ComponentID)
		assert.Equal(t, threat.ID, matches[0].ThreatID)
		assert.Equal(t, pattern.ID, matches[0].PatternID)
		assert.Equal(t, *pattern, matches[0].Pattern)
	})
}

func TestGetComponentThreatsByExistingThreatPatterns(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.Tag{},
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("ComponentMatchesMultiplePatterns", func(t *testing.T) {
		// Create a product and instance
		instance, err := service.CreateComponent("Main Database", "The main database instance", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create tags and assign them to the instance
		criticalTag, err := service.CreateTag("critical", "Critical component", "#FF0000")
		require.NoError(t, err)

		databaseTag, err := service.CreateTag("database", "Database component", "#00FF00")
		require.NoError(t, err)

		err = service.AssignTagToComponent(criticalTag.ID, instance.ID)
		require.NoError(t, err)

		err = service.AssignTagToComponent(databaseTag.ID, instance.ID)
		require.NoError(t, err)

		// Create threats
		dataBreachThreat, err := service.CreateThreat("Data Breach", "Data breach threat")
		require.NoError(t, err)

		sqlInjectionThreat, err := service.CreateThreat("SQL Injection", "SQL injection threat")
		require.NoError(t, err)

		// Create threat patterns
		criticalPattern, err := CreateThreatPatternWithConditions(
			"Critical Components",
			"Critical components pattern",
			dataBreachThreat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "critical",
				},
			},
		)
		require.NoError(t, err)

		databasePattern, err := CreateThreatPatternWithConditions(
			"Database Components",
			"Database components pattern",
			sqlInjectionThreat.ID,
			true,
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "database",
				},
			},
		)
		require.NoError(t, err)

		// Create an inactive pattern that shouldn't match
		_, err = CreateThreatPatternWithConditions(
			"Inactive Pattern",
			"This pattern is inactive",
			dataBreachThreat.ID,
			false, // inactive
			[]models.PatternCondition{
				{
					ConditionType: models.ConditionTypeTag.String(),
					Operator:      models.OperatorContains.String(),
					Value:         "database",
				},
			},
		)
		require.NoError(t, err)

		// Test the function
		matches, err := GetComponentThreatsByExistingThreatPatterns(*instance)
		require.NoError(t, err)

		// Verify the matches
		require.Len(t, matches, 2, "Should have exactly two matches")

		// Check that both patterns are matched
		patternIDs := make(map[uuid.UUID]bool)
		for _, match := range matches {
			patternIDs[match.PatternID] = true
			assert.Equal(t, instance.ID, match.ComponentID)
		}

		assert.True(t, patternIDs[criticalPattern.ID], "Should match critical pattern")
		assert.True(t, patternIDs[databasePattern.ID], "Should match database pattern")
	})
}
