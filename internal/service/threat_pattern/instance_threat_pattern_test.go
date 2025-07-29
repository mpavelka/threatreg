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
		webProduct, err := service.CreateProduct("Web Application", "A web application")
		require.NoError(t, err)

		dbProduct, err := service.CreateProduct("Database", "A database system")
		require.NoError(t, err)

		apiProduct, err := service.CreateProduct("API Service", "An API service")
		require.NoError(t, err)

		// Create instances
		webServer, err := service.CreateInstance("Web Server", webProduct.ID)
		require.NoError(t, err)

		database, err := service.CreateInstance("Database", dbProduct.ID)
		require.NoError(t, err)

		apiServer, err := service.CreateInstance("API Server", apiProduct.ID)
		require.NoError(t, err)

		adminPanel, err := service.CreateInstance("Admin Panel", webProduct.ID)
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
		err = service.AssignTagToInstance(internetFacingTag.ID, webServer.ID)
		require.NoError(t, err)

		err = service.AssignTagToInstance(databaseTag.ID, database.ID)
		require.NoError(t, err)

		err = service.AssignTagToInstance(privilegedTag.ID, adminPanel.ID)
		require.NoError(t, err)

		// Assign tags to products
		err = service.AssignTagToProduct(highSecurityTag.ID, dbProduct.ID)
		require.NoError(t, err)

		// Create relationships
		err = service.AddRelationship(&webServer.ID, nil, &database.ID, nil, "connects_to", "")
		require.NoError(t, err)

		err = service.AddRelationship(&adminPanel.ID, nil, &database.ID, nil, "reads_from", "")
		require.NoError(t, err)

		err = service.AddRelationship(&apiServer.ID, nil, &database.ID, nil, "connects_to", "")
		require.NoError(t, err)

		// Create threats
		sqlInjectionThreat, err := service.CreateThreat("SQL Injection", "SQL injection vulnerability")
		require.NoError(t, err)

		privilegeEscalationThreat, err := service.CreateThreat("Privilege Escalation", "Privilege escalation threat")
		require.NoError(t, err)

		dataExposureThreat, err := service.CreateThreat("Data Exposure", "Data exposure threat")
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
		matches, err := GetAllInstancesThreatsByThreatPattern()
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
		isolatedProduct, err := service.CreateProduct("Isolated Service", "An isolated service")
		require.NoError(t, err)

		isolatedInstance, err := service.CreateInstance("Isolated Instance", isolatedProduct.ID)
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
		matches, err := GetAllInstancesThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify no matches for isolated instance
		_, exists := matches[isolatedInstance.ID]
		assert.False(t, exists, "Isolated instance should have no matches")
	})

	t.Run("InactivePattern", func(t *testing.T) {
		// Create a product and instance
		testProduct, err := service.CreateProduct("Test Product", "A test product")
		require.NoError(t, err)

		testInstance, err := service.CreateInstance("Test Instance", testProduct.ID)
		require.NoError(t, err)

		// Create a tag and assign it to the instance
		testTag, err := service.CreateTag("test-tag", "Test tag", "#FFFFFF")
		require.NoError(t, err)

		err = service.AssignTagToInstance(testTag.ID, testInstance.ID)
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
		matches, err := GetAllInstancesThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify no matches for instance with inactive pattern
		_, exists := matches[testInstance.ID]
		assert.False(t, exists, "Instance should have no matches from inactive pattern")
	})

	t.Run("TagExistsCondition", func(t *testing.T) {
		// Create product and instances
		product, err := service.CreateProduct("Tag Exists Product", "Product for testing tag existence")
		require.NoError(t, err)

		taggedInstance, err := service.CreateInstance("Tagged Instance", product.ID)
		require.NoError(t, err)

		untaggedInstance, err := service.CreateInstance("Untagged Instance", product.ID)
		require.NoError(t, err)

		// Create and assign a tag to one instance
		tag, err := service.CreateTag("any-tag", "Any tag", "#000000")
		require.NoError(t, err)

		err = service.AssignTagToInstance(tag.ID, taggedInstance.ID)
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
		matches, err := GetAllInstancesThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify tagged instance matches but untagged doesn't
		_, taggedExists := matches[taggedInstance.ID]
		assert.True(t, taggedExists, "Tagged instance should match")

		_, untaggedExists := matches[untaggedInstance.ID]
		assert.False(t, untaggedExists, "Untagged instance should not match")
	})

	t.Run("RelationshipExistsCondition", func(t *testing.T) {
		// Create products and instances
		product1, err := service.CreateProduct("Product 1", "First product")
		require.NoError(t, err)

		product2, err := service.CreateProduct("Product 2", "Second product")
		require.NoError(t, err)

		instance1, err := service.CreateInstance("Instance 1", product1.ID)
		require.NoError(t, err)

		instance2, err := service.CreateInstance("Instance 2", product2.ID)
		require.NoError(t, err)

		isolatedInstance, err := service.CreateInstance("Isolated Instance", product1.ID)
		require.NoError(t, err)

		// Create relationship between instance1 and instance2
		err = service.AddRelationship(&instance1.ID, nil, &instance2.ID, nil, "depends_on", "")
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
		matches, err := GetAllInstancesThreatsByThreatPattern()
		require.NoError(t, err)

		// Verify instance1 matches but isolated instance doesn't
		_, instance1Exists := matches[instance1.ID]
		assert.True(t, instance1Exists, "Instance1 should match (has depends_on relationship)")

		_, isolatedExists := matches[isolatedInstance.ID]
		assert.False(t, isolatedExists, "Isolated instance should not match (no depends_on relationship)")
	})
}

func TestGetInstanceThreatsByThreatPattern(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("InstanceMatchesPattern", func(t *testing.T) {
		// Create a product and instance
		product, err := service.CreateProduct("Web Application", "A web application")
		require.NoError(t, err)

		instance, err := service.CreateInstance("Web Server", product.ID)
		require.NoError(t, err)

		// Create a tag and assign it to the instance
		tag, err := service.CreateTag("internet-facing", "Internet facing component", "#FF0000")
		require.NoError(t, err)

		err = service.AssignTagToInstance(tag.ID, instance.ID)
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
		matches, err := GetInstanceThreatsByThreatPattern(*instance, *pattern)
		require.NoError(t, err)

		// Verify the match
		require.Len(t, matches, 1, "Should have exactly one match")
		assert.Equal(t, instance.ID, matches[0].InstanceID)
		assert.Equal(t, threat.ID, matches[0].ThreatID)
		assert.Equal(t, pattern.ID, matches[0].PatternID)
		assert.Equal(t, *pattern, matches[0].Pattern)
	})
}

func TestGetInstanceThreatsByExistingThreatPatterns(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
		&models.Threat{},
		&models.ThreatPattern{},
		&models.PatternCondition{},
	)
	defer cleanup()

	t.Run("InstanceMatchesMultiplePatterns", func(t *testing.T) {
		// Create a product and instance
		product, err := service.CreateProduct("Database System", "A database system")
		require.NoError(t, err)

		instance, err := service.CreateInstance("Main Database", product.ID)
		require.NoError(t, err)

		// Create tags and assign them to the instance
		criticalTag, err := service.CreateTag("critical", "Critical component", "#FF0000")
		require.NoError(t, err)

		databaseTag, err := service.CreateTag("database", "Database component", "#00FF00")
		require.NoError(t, err)

		err = service.AssignTagToInstance(criticalTag.ID, instance.ID)
		require.NoError(t, err)

		err = service.AssignTagToInstance(databaseTag.ID, instance.ID)
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
		matches, err := GetInstanceThreatsByExistingThreatPatterns(*instance)
		require.NoError(t, err)

		// Verify the matches
		require.Len(t, matches, 2, "Should have exactly two matches")

		// Check that both patterns are matched
		patternIDs := make(map[uuid.UUID]bool)
		for _, match := range matches {
			patternIDs[match.PatternID] = true
			assert.Equal(t, instance.ID, match.InstanceID)
		}

		assert.True(t, patternIDs[criticalPattern.ID], "Should match critical pattern")
		assert.True(t, patternIDs[databasePattern.ID], "Should match database pattern")
	})
}
