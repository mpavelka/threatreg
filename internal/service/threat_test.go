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

	t.Run("ListByDomainWithUnresolvedByInstancesCount_ProductLevelThreats", func(t *testing.T) {
		// Scenario 1: Multiple threats assigned to a product, instance added to domain
		// Expected: Multiple records with UnresolvedByInstancesCount = 1

		// Create domain
		domain, _ := CreateDomain("Test Domain 1", "Domain for testing product-level threats")

		// Create product
		product, _ := CreateProduct("Test Product 1", "Product with threats")

		// Create instance of the product
		instance, _ := CreateInstance("Test Instance 1", product.ID)

		// Add instance to domain
		_ = AddInstanceToDomain(domain.ID, instance.ID)

		// Create multiple threats
		threat1, _ := CreateThreat("SQL Injection", "Database query manipulation threat")
		threat2, _ := CreateThreat("XSS", "Cross-site scripting threat")
		threat3, _ := CreateThreat("CSRF", "Cross-site request forgery threat")

		// Assign threats to product
		_, _ = AssignThreatToProduct(product.ID, threat1.ID)
		_, _ = AssignThreatToProduct(product.ID, threat2.ID)
		_, _ = AssignThreatToProduct(product.ID, threat3.ID)

		// Test the function
		results, err := ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)

		// Should return 3 threats, each with UnresolvedByInstancesCount = 1
		assert.Len(t, results, 3)

		// Create a map for easier assertion
		threatCountMap := make(map[uuid.UUID]int)
		for _, result := range results {
			threatCountMap[result.ID] = result.UnresolvedByInstancesCount
		}

		assert.Equal(t, 1, threatCountMap[threat1.ID], "Threat1 should have 1 unresolved instance")
		assert.Equal(t, 1, threatCountMap[threat2.ID], "Threat2 should have 1 unresolved instance")
		assert.Equal(t, 1, threatCountMap[threat3.ID], "Threat3 should have 1 unresolved instance")
	})

	t.Run("ListByDomainWithUnresolvedByInstancesCount_MultipleInstances", func(t *testing.T) {
		// Scenario 2: Multiple threats assigned to multiple products, instances added to domain
		// Expected: Multiple records with UnresolvedByInstancesCount = 2

		// Create domain
		domain, _ := CreateDomain("Test Domain 2", "Domain for testing multiple instances")

		// Create two products
		product1, _ := CreateProduct("Test Product 2A", "First product with threats")
		product2, _ := CreateProduct("Test Product 2B", "Second product with threats")

		// Create instances for both products
		instance1, _ := CreateInstance("Test Instance 2A", product1.ID)
		instance2, _ := CreateInstance("Test Instance 2B", product2.ID)

		// Add both instances to domain
		_ = AddInstanceToDomain(domain.ID, instance1.ID)
		_ = AddInstanceToDomain(domain.ID, instance2.ID)

		// Create threats
		threat1, _ := CreateThreat("Buffer Overflow", "Memory corruption vulnerability")
		threat2, _ := CreateThreat("Path Traversal", "Directory traversal vulnerability")

		// Assign threats to both products
		_, _ = AssignThreatToProduct(product1.ID, threat1.ID)
		_, _ = AssignThreatToProduct(product1.ID, threat2.ID)
		_, _ = AssignThreatToProduct(product2.ID, threat1.ID)
		_, _ = AssignThreatToProduct(product2.ID, threat2.ID)

		// Test the function
		results, err := ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)

		// Should return 2 threats, each with UnresolvedByInstancesCount = 2
		assert.Len(t, results, 2)

		// Create a map for easier assertion
		threatCountMap := make(map[uuid.UUID]int)
		for _, result := range results {
			threatCountMap[result.ID] = result.UnresolvedByInstancesCount
		}

		assert.Equal(t, 2, threatCountMap[threat1.ID], "Threat1 should have 2 unresolved instances")
		assert.Equal(t, 2, threatCountMap[threat2.ID], "Threat2 should have 2 unresolved instances")
	})

	t.Run("ListByDomainWithUnresolvedByInstancesCount_WithResolution", func(t *testing.T) {
		// Scenario 3: Following up scenario 2, mark one threat as resolved by one instance
		// Expected: That threat should have UnresolvedByInstancesCount = 1, other should remain 2

		// Create domain
		domain, _ := CreateDomain("Test Domain 3", "Domain for testing with resolutions")

		// Create two products
		product1, _ := CreateProduct("Test Product 3A", "First product with threats")
		product2, _ := CreateProduct("Test Product 3B", "Second product with threats")

		// Create instances for both products
		instance1, _ := CreateInstance("Test Instance 3A", product1.ID)
		instance2, _ := CreateInstance("Test Instance 3B", product2.ID)

		// Add both instances to domain
		_ = AddInstanceToDomain(domain.ID, instance1.ID)
		_ = AddInstanceToDomain(domain.ID, instance2.ID)

		// Create threats
		threat1, _ := CreateThreat("Privilege Escalation", "Gaining elevated access")
		threat2, _ := CreateThreat("Information Disclosure", "Unauthorized information access")

		// Assign threats to both products
		assignment1_1, _ := AssignThreatToProduct(product1.ID, threat1.ID)
		_, _ = AssignThreatToProduct(product1.ID, threat2.ID)
		_, _ = AssignThreatToProduct(product2.ID, threat1.ID)
		_, _ = AssignThreatToProduct(product2.ID, threat2.ID)

		// Check initial state - both threats should have 2 unresolved instances
		results, err := ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)
		assert.Len(t, results, 2)

		initialThreatCountMap := make(map[uuid.UUID]int)
		for _, result := range results {
			initialThreatCountMap[result.ID] = result.UnresolvedByInstancesCount
		}

		assert.Equal(t, 2, initialThreatCountMap[threat1.ID], "Threat1 should initially have 2 unresolved instances")
		assert.Equal(t, 2, initialThreatCountMap[threat2.ID], "Threat2 should initially have 2 unresolved instances")

		// Resolve threat1 for instance1 (mark as "resolved")
		_, _ = CreateThreatResolution(
			assignment1_1.ID,
			&instance1.ID,
			nil,
			models.ThreatAssignmentResolutionStatusResolved,
			"Fixed the privilege escalation vulnerability in instance 1",
		)

		// Test the function again after resolution
		results, err = ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)

		// Should still return 2 threats, but threat1 should now have UnresolvedByInstancesCount = 1
		assert.Len(t, results, 2)

		finalThreatCountMap := make(map[uuid.UUID]int)
		for _, result := range results {
			finalThreatCountMap[result.ID] = result.UnresolvedByInstancesCount
		}

		assert.Equal(t, 1, finalThreatCountMap[threat1.ID], "Threat1 should have 1 unresolved instance after resolution")
		assert.Equal(t, 2, finalThreatCountMap[threat2.ID], "Threat2 should still have 2 unresolved instances")
	})

	t.Run("ListByDomainWithUnresolvedByInstancesCount_ProductAndInstanceLevelAssignments", func(t *testing.T) {
		// Scenario 4: Product gets threat assignment, then instance gets same threat assigned
		// Expected: UnresolvedByInstancesCount should show 1 (not 2, since it's the same instance)
		// When product-level assignment is resolved, count should remain 1
		// Only when instance-level assignment is resolved, count should become 0

		// Create domain
		domain, _ := CreateDomain("Test Domain 4", "Domain for testing overlapping assignments")

		// Create product
		product, _ := CreateProduct("Test Product 4", "Product with overlapping threat assignments")

		// Create instance
		instance, _ := CreateInstance("Test Instance 4", product.ID)

		// Add instance to domain
		_ = AddInstanceToDomain(domain.ID, instance.ID)

		// Create threat
		threat, _ := CreateThreat("Code Injection", "Arbitrary code execution vulnerability")

		// Assign threat to product first
		productAssignment, _ := AssignThreatToProduct(product.ID, threat.ID)

		// Assign same threat to instance
		instanceAssignment, _ := AssignThreatToInstance(instance.ID, threat.ID)

		// Test initial state - should show 1 unresolved instance (not 2)
		results, err := ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, threat.ID, results[0].ID)
		assert.Equal(t, 1, results[0].UnresolvedByInstancesCount, "Should show 1 unresolved instance despite both product and instance assignments")

		// Resolve the product-level assignment
		_, _ = CreateThreatResolution(
			productAssignment.ID,
			&instance.ID,
			nil,
			models.ThreatAssignmentResolutionStatusResolved,
			"Resolved product-level assignment for this instance",
		)

		// Test after product-level resolution - should still show 1 unresolved instance
		results, err = ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, threat.ID, results[0].ID)
		assert.Equal(t, 1, results[0].UnresolvedByInstancesCount, "Should still show 1 unresolved instance after product-level resolution")

		// Resolve the instance-level assignment
		_, _ = CreateThreatResolution(
			instanceAssignment.ID,
			&instance.ID,
			nil,
			models.ThreatAssignmentResolutionStatusResolved,
			"Resolved instance-level assignment",
		)

		// Test after instance-level resolution - should show 0 results (threat fully resolved)
		results, err = ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)
		assert.Len(t, results, 0, "Should show no unresolved threats after both assignments are resolved")
	})

	t.Run("ListByDomainWithUnresolvedByInstancesCount_EmptyDomain", func(t *testing.T) {
		// Test empty domain
		domain, _ := CreateDomain("Empty Domain", "Domain with no instances")

		results, err := ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)

		// Should return empty slice
		assert.Len(t, results, 0)
	})

	t.Run("ListByDomainWithUnresolvedByInstancesCount_NoThreats", func(t *testing.T) {
		// Test domain with instances but no threat assignments
		domain, _ := CreateDomain("No Threats Domain", "Domain with instances but no threats")

		product, _ := CreateProduct("No Threats Product", "Product without threats")

		instance, _ := CreateInstance("No Threats Instance", product.ID)

		_ = AddInstanceToDomain(domain.ID, instance.ID)

		results, err := ListByDomainWithUnresolvedByInstancesCount(domain.ID)
		require.NoError(t, err)

		// Should return empty slice
		assert.Len(t, results, 0)
	})
}
