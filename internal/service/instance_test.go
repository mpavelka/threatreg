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

		for i := range 2 {
			app1Name := fmt.Sprintf("Test Product App %d", i+1)
			CreateInstance(app1Name, testProduct.ID)
			app2Name := fmt.Sprintf("Another Product App %d", i+1)
			CreateInstance(app2Name, anotherProduct.ID)
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

		for _, ta := range testApps {
			CreateInstance(ta.name, ta.productID)
			require.NoError(t, err)
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

	t.Run("AssignThreatToInstance", func(t *testing.T) {
		// Create an instance first
		instance, err := CreateInstance("Test Instance for Threat", testProduct.ID)
		require.NoError(t, err)

		// Create a threat first
		threat, err := CreateThreat("Test Instance Threat", "A test threat for instance assignment")
		require.NoError(t, err)

		// Assign threat to instance
		assignment, err := AssignThreatToInstance(instance.ID, threat.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.NotEqual(t, 0, assignment.ID)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, instance.ID, assignment.InstanceID)
		assert.Equal(t, uuid.Nil, assignment.ProductID) // Should be nil for instance assignment

		// Verify assignment was saved to database
		db := database.GetDB()
		var dbAssignment models.ThreatAssignment
		err = db.First(&dbAssignment, "id = ?", assignment.ID).Error
		require.NoError(t, err)
		assert.Equal(t, assignment.ThreatID, dbAssignment.ThreatID)
		assert.Equal(t, assignment.InstanceID, dbAssignment.InstanceID)
		assert.Equal(t, uuid.Nil, dbAssignment.ProductID)
	})

	t.Run("AssignThreatToInstance_Duplicate", func(t *testing.T) {
		// Create an instance and threat
		instance, err := CreateInstance("Duplicate Test Instance", testProduct.ID)
		require.NoError(t, err)

		threat, err := CreateThreat("Duplicate Instance Threat", "A test threat for duplicate assignment")
		require.NoError(t, err)

		// Assign threat to instance first time
		assignment1, err := AssignThreatToInstance(instance.ID, threat.ID)
		require.NoError(t, err)
		require.NotNil(t, assignment1)

		// Try to assign the same threat to the same instance again
		assignment2, err := AssignThreatToInstance(instance.ID, threat.ID)

		// Should return the existing assignment, not create a new one
		require.NoError(t, err)
		assert.NotNil(t, assignment2)
		assert.Equal(t, assignment1.ID, assignment2.ID)
		assert.Equal(t, assignment1.ThreatID, assignment2.ThreatID)
		assert.Equal(t, assignment1.InstanceID, assignment2.InstanceID)

		// Verify that only one assignment exists in the database for this threat/instance combination
		db := database.GetDB()
		var count int64
		err = db.Model(&models.ThreatAssignment{}).
			Where("threat_id = ? AND instance_id = ? AND (product_id IS NULL OR product_id = ?)", threat.ID, instance.ID, uuid.Nil).
			Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Should only have one assignment for this threat/instance combination")
	})

	t.Run("AssignThreatToInstance_InvalidThreatID", func(t *testing.T) {
		// Create an instance
		instance, err := CreateInstance("Invalid Threat Test Instance", testProduct.ID)
		require.NoError(t, err)

		// Try to assign non-existent threat
		nonExistentThreatID := uuid.New()
		assignment, err := AssignThreatToInstance(instance.ID, nonExistentThreatID)

		// Should succeed (foreign key constraint allows it, but relationship won't load)
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, nonExistentThreatID, assignment.ThreatID)
		assert.Equal(t, instance.ID, assignment.InstanceID)
	})

	t.Run("AssignThreatToInstance_InvalidInstanceID", func(t *testing.T) {
		// Create a threat
		threat, err := CreateThreat("Invalid Instance Test Threat", "A test threat")
		require.NoError(t, err)

		// Try to assign to non-existent instance
		nonExistentInstanceID := uuid.New()
		assignment, err := AssignThreatToInstance(nonExistentInstanceID, threat.ID)

		// Should succeed (foreign key constraint allows it, but relationship won't load)
		require.NoError(t, err)
		assert.NotNil(t, assignment)
		assert.Equal(t, threat.ID, assignment.ThreatID)
		assert.Equal(t, nonExistentInstanceID, assignment.InstanceID)
	})

	t.Run("ListThreatAssignmentsByInstanceID", func(t *testing.T) {
		// Create an instance for testing
		instance, err := CreateInstance("Test Instance for Assignments", testProduct.ID)
		require.NoError(t, err)

		// Create multiple threats for testing
		threat1, err := CreateThreat("Threat 1", "First test threat")
		require.NoError(t, err)
		threat2, err := CreateThreat("Threat 2", "Second test threat")
		require.NoError(t, err)
		threat3, err := CreateThreat("Threat 3", "Third test threat")
		require.NoError(t, err)

		// Assign threats to the instance
		assignment1, err := AssignThreatToInstance(instance.ID, threat1.ID)
		require.NoError(t, err)
		assignment2, err := AssignThreatToInstance(instance.ID, threat2.ID)
		require.NoError(t, err)

		// Assign one threat to a different instance to ensure filtering works
		otherInstance, err := CreateInstance("Other Instance", testProduct.ID)
		require.NoError(t, err)
		_, err = AssignThreatToInstance(otherInstance.ID, threat3.ID)
		require.NoError(t, err)

		// List threat assignments for our test instance
		assignments, err := ListThreatAssignmentsByInstanceID(instance.ID)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, assignments, 2)

		// Check that we got the correct assignments
		assignmentIDs := []int{assignments[0].ID, assignments[1].ID}
		assert.Contains(t, assignmentIDs, assignment1.ID)
		assert.Contains(t, assignmentIDs, assignment2.ID)

		// Verify threat relationships are loaded
		for _, assignment := range assignments {
			assert.NotEmpty(t, assignment.Threat.Title)
			assert.Equal(t, instance.ID, assignment.InstanceID)
			assert.Equal(t, uuid.Nil, assignment.ProductID) // Should be nil for instance assignments
		}
	})

	t.Run("ListThreatAssignmentsByInstanceID_Empty", func(t *testing.T) {
		// Create an instance with no threat assignments
		instance, err := CreateInstance("Empty Assignments Instance", testProduct.ID)
		require.NoError(t, err)

		// List threat assignments for this instance
		assignments, err := ListThreatAssignmentsByInstanceID(instance.ID)

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, assignments, 0)
	})

	t.Run("ListThreatAssignmentsByInstanceID_InvalidInstanceID", func(t *testing.T) {
		// Try to list assignments for non-existent instance
		nonExistentInstanceID := uuid.New()
		assignments, err := ListThreatAssignmentsByInstanceID(nonExistentInstanceID)

		// Should succeed but return empty slice
		require.NoError(t, err)
		assert.Len(t, assignments, 0)
	})
}

func TestListThreatAssignmentsByInstanceIDWithResolutionByInstanceID(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("Happy Flow", func(t *testing.T) {
		// Create product1 and product2
		product1, err := CreateProduct("Product 1", "First test product")
		require.NoError(t, err)
		product2, err := CreateProduct("Product 2", "Second test product")
		require.NoError(t, err)

		// Create instance1 (for product1) and instance2 (for product2)
		instance1, err := CreateInstance("Instance 1", product1.ID)
		require.NoError(t, err)
		instance2, err := CreateInstance("Instance 2", product2.ID)
		require.NoError(t, err)

		// Create threat1 and threat2
		threat1, err := CreateThreat("Threat 1", "First test threat")
		require.NoError(t, err)
		threat2, err := CreateThreat("Threat 2", "Second test threat")
		require.NoError(t, err)

		// Assign threat1 to instance1 (note: AssignThreatToInstance takes instanceID first, then threatID)
		assignment1, err := AssignThreatToInstance(instance1.ID, threat1.ID)
		require.NoError(t, err)
		require.NotNil(t, assignment1)
		t.Logf("Assignment1: ID=%d, ThreatID=%s, InstanceID=%s, ProductID=%s", 
			assignment1.ID, assignment1.ThreatID, assignment1.InstanceID, assignment1.ProductID)

		// Assign threat2 to instance2
		assignment2, err := AssignThreatToInstance(instance2.ID, threat2.ID)
		require.NoError(t, err)
		require.NotNil(t, assignment2)
		t.Logf("Assignment2: ID=%d, ThreatID=%s, InstanceID=%s, ProductID=%s", 
			assignment2.ID, assignment2.ThreatID, assignment2.InstanceID, assignment2.ProductID)

		// Create threat resolutions for both instance1 and instance2
		resolution1, err := CreateThreatResolution(
			assignment1.ID,
			&instance1.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution for instance1",
		)
		require.NoError(t, err)

		resolution2, err := CreateThreatResolution(
			assignment2.ID,
			&instance2.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Resolution for instance2",
		)
		require.NoError(t, err)

		// Delegate resolution1 to resolution2
		err = DelegateResolution(*resolution1, *resolution2)
		require.NoError(t, err)

		// Mark resolution2 as resolved
		resolvedStatus := models.ThreatAssignmentResolutionStatusResolved
		_, err = UpdateThreatResolution(resolution2.ID, &resolvedStatus, nil)
		require.NoError(t, err)

		// First verify basic ListThreatAssignmentsByInstanceID works
		basicResults, err := ListThreatAssignmentsByInstanceID(instance1.ID)
		require.NoError(t, err)
		require.Len(t, basicResults, 1, "Basic function should return one assignment")

		// Test ListThreatAssignmentsByInstanceIDWithResolutionByInstanceID for instance1 filtered by instance1
		results, err := ListThreatAssignmentsByInstanceIDWithResolutionByInstanceID(instance1.ID, instance1.ID)
		require.NoError(t, err)

		// Assertions
		assert.Len(t, results, 1, "Should return one threat assignment for instance1")
		
		result := results[0]
		assert.Equal(t, assignment1.ID, result.ID)
		assert.Equal(t, threat1.ID, result.ThreatID)
		assert.Equal(t, uuid.Nil, result.ProductID) // Instance assignment should have nil product ID
		assert.Equal(t, instance1.ID, result.InstanceID)

		// Verify threat relationship is loaded
		assert.Equal(t, threat1.ID, result.Threat.ID)
		assert.Equal(t, "Threat 1", result.Threat.Title)

		// Verify instance relationship is loaded
		assert.Equal(t, instance1.ID, result.Instance.ID)
		assert.Equal(t, "Instance 1", result.Instance.Name)

		// Verify resolution status - should show resolved because resolution1 was delegated to resolution2 which is resolved
		assert.NotNil(t, result.ResolutionStatus)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, *result.ResolutionStatus)

		// Verify delegation status
		assert.True(t, result.IsDelegated, "Should show as delegated")

		// Test with different instance filter - should return same assignment but no resolution info
		resultsOtherInstance, err := ListThreatAssignmentsByInstanceIDWithResolutionByInstanceID(instance1.ID, instance2.ID)
		require.NoError(t, err)

		assert.Len(t, resultsOtherInstance, 1, "Should still return the threat assignment for instance1")
		otherResult := resultsOtherInstance[0]
		assert.Equal(t, assignment1.ID, otherResult.ID)
		
		// But resolution info should be nil since we filtered by instance2 but resolution1 is for instance1
		assert.Nil(t, otherResult.ResolutionStatus, "Should not have resolution status for different instance")
		assert.False(t, otherResult.IsDelegated, "Should not show as delegated for different instance")

		// Test with instance2 - should return the assignment for instance2
		resultsInstance2, err := ListThreatAssignmentsByInstanceIDWithResolutionByInstanceID(instance2.ID, instance2.ID)
		require.NoError(t, err)
		assert.Len(t, resultsInstance2, 1, "Should return one assignment for instance2")
		
		instance2Result := resultsInstance2[0]
		assert.Equal(t, assignment2.ID, instance2Result.ID)
		assert.Equal(t, threat2.ID, instance2Result.ThreatID)
		assert.Equal(t, instance2.ID, instance2Result.InstanceID)
		
		// This should show resolved status since resolution2 is resolved
		assert.NotNil(t, instance2Result.ResolutionStatus)
		assert.Equal(t, models.ThreatAssignmentResolutionStatusResolved, *instance2Result.ResolutionStatus)
		assert.False(t, instance2Result.IsDelegated, "Should not show as delegated since resolution2 is not delegated")
	})
}
