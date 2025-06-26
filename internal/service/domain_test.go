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

func TestDomainService_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("CreateDomain", func(t *testing.T) {
		name := "Test Domain"
		description := "A test domain description"

		domain, err := CreateDomain(name, description)

		require.NoError(t, err)
		assert.NotNil(t, domain)
		assert.NotEqual(t, uuid.Nil, domain.ID)
		assert.Equal(t, name, domain.Name)
		assert.Equal(t, description, domain.Description)

		db := database.GetDB()
		var dbDomain models.Domain
		err = db.First(&dbDomain, "id = ?", domain.ID).Error
		require.NoError(t, err)
		assert.Equal(t, domain.ID, dbDomain.ID)
		assert.Equal(t, name, dbDomain.Name)
		assert.Equal(t, description, dbDomain.Description)
	})

	t.Run("GetDomain", func(t *testing.T) {
		name := "Get Test Domain"
		description := "Domain for get test"
		createdDomain, err := CreateDomain(name, description)
		require.NoError(t, err)

		retrievedDomain, err := GetDomain(createdDomain.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedDomain)
		assert.Equal(t, createdDomain.ID, retrievedDomain.ID)
		assert.Equal(t, name, retrievedDomain.Name)
		assert.Equal(t, description, retrievedDomain.Description)
	})

	t.Run("GetDomain_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()
		domain, err := GetDomain(nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, domain)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateDomain", func(t *testing.T) {
		originalName := "Original Domain"
		originalDescription := "Original description"
		createdDomain, err := CreateDomain(originalName, originalDescription)
		require.NoError(t, err)

		newName := "Updated Domain"
		newDescription := "Updated description"
		updatedDomain, err := UpdateDomain(createdDomain.ID, &newName, &newDescription)

		require.NoError(t, err)
		assert.NotNil(t, updatedDomain)
		assert.Equal(t, createdDomain.ID, updatedDomain.ID)
		assert.Equal(t, newName, updatedDomain.Name)
		assert.Equal(t, newDescription, updatedDomain.Description)

		db := database.GetDB()
		var dbDomain models.Domain
		err = db.First(&dbDomain, "id = ?", createdDomain.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbDomain.Name)
		assert.Equal(t, newDescription, dbDomain.Description)
	})

	t.Run("UpdateDomain_PartialUpdate", func(t *testing.T) {
		originalName := "Partial Update Domain"
		originalDescription := "Original description"
		createdDomain, err := CreateDomain(originalName, originalDescription)
		require.NoError(t, err)

		newName := "New Name Only"
		updatedDomain, err := UpdateDomain(createdDomain.ID, &newName, nil)

		require.NoError(t, err)
		assert.NotNil(t, updatedDomain)
		assert.Equal(t, createdDomain.ID, updatedDomain.ID)
		assert.Equal(t, newName, updatedDomain.Name)
		assert.Equal(t, originalDescription, updatedDomain.Description)

		db := database.GetDB()
		var dbDomain models.Domain
		err = db.First(&dbDomain, "id = ?", createdDomain.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbDomain.Name)
		assert.Equal(t, originalDescription, dbDomain.Description)
	})

	t.Run("UpdateDomain_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()
		newName := "New Name"
		domain, err := UpdateDomain(nonExistentID, &newName, nil)

		assert.Error(t, err)
		assert.Nil(t, domain)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteDomain", func(t *testing.T) {
		name := "Delete Test Domain"
		description := "Domain to be deleted"
		createdDomain, err := CreateDomain(name, description)
		require.NoError(t, err)

		err = DeleteDomain(createdDomain.ID)

		require.NoError(t, err)

		db := database.GetDB()
		var dbDomain models.Domain
		err = db.First(&dbDomain, "id = ?", createdDomain.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteDomain_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := DeleteDomain(nonExistentID)

		assert.NoError(t, err)
	})

	t.Run("ListDomains", func(t *testing.T) {
		db := database.GetDB()
		db.Exec("DELETE FROM domains")

		domains := []struct {
			name        string
			description string
		}{
			{"Domain 1", "Description 1"},
			{"Domain 2", "Description 2"},
			{"Domain 3", "Description 3"},
		}

		var createdDomains []*models.Domain
		for _, d := range domains {
			domain, err := CreateDomain(d.name, d.description)
			require.NoError(t, err)
			createdDomains = append(createdDomains, domain)
		}

		retrievedDomains, err := ListDomains()

		require.NoError(t, err)
		assert.Len(t, retrievedDomains, len(domains))

		domainMap := make(map[uuid.UUID]models.Domain)
		for _, d := range retrievedDomains {
			domainMap[d.ID] = d
		}

		for _, created := range createdDomains {
			retrieved, exists := domainMap[created.ID]
			assert.True(t, exists, "Created domain should exist in list")
			assert.Equal(t, created.Name, retrieved.Name)
			assert.Equal(t, created.Description, retrieved.Description)
		}
	})

	t.Run("ListDomains_Empty", func(t *testing.T) {
		db := database.GetDB()
		db.Exec("DELETE FROM domains")

		domains, err := ListDomains()

		require.NoError(t, err)
		assert.Len(t, domains, 0)
	})

	t.Run("AddInstanceToDomain", func(t *testing.T) {
		domain, err := CreateDomain("Test Domain", "Test Description")
		require.NoError(t, err)

		product, err := CreateProduct("Test Product", "Test Product Description")
		require.NoError(t, err)

		instance, err := CreateInstance("Test Instance", product.ID)
		require.NoError(t, err)

		err = AddInstanceToDomain(domain.ID, instance.ID)
		require.NoError(t, err)

		instances, err := GetInstancesByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 1)
		assert.Equal(t, instance.ID, instances[0].ID)
		assert.Equal(t, instance.Name, instances[0].Name)

		domains, err := GetDomainsByInstance(instance.ID)
		require.NoError(t, err)
		assert.Len(t, domains, 1)
		assert.Equal(t, domain.ID, domains[0].ID)
		assert.Equal(t, domain.Name, domains[0].Name)
	})

	t.Run("AddInstanceToDomain_Multiple", func(t *testing.T) {
		domain, err := CreateDomain("Multi Test Domain", "Test Description")
		require.NoError(t, err)

		product, err := CreateProduct("Multi Test Product", "Test Product Description")
		require.NoError(t, err)

		var instances []*models.Instance
		for i := 0; i < 3; i++ {
			instance, err := CreateInstance(fmt.Sprintf("Test Instance %d", i+1), product.ID)
			require.NoError(t, err)
			instances = append(instances, instance)

			err = AddInstanceToDomain(domain.ID, instance.ID)
			require.NoError(t, err)
		}

		retrievedInstances, err := GetInstancesByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedInstances, 3)

		instanceMap := make(map[uuid.UUID]models.Instance)
		for _, inst := range retrievedInstances {
			instanceMap[inst.ID] = inst
		}

		for _, created := range instances {
			retrieved, exists := instanceMap[created.ID]
			assert.True(t, exists, "Created instance should exist in domain")
			assert.Equal(t, created.Name, retrieved.Name)
		}
	})

	t.Run("RemoveInstanceFromDomain", func(t *testing.T) {
		domain, err := CreateDomain("Remove Test Domain", "Test Description")
		require.NoError(t, err)

		product, err := CreateProduct("Remove Test Product", "Test Product Description")
		require.NoError(t, err)

		instance1, err := CreateInstance("Remove Test Instance 1", product.ID)
		require.NoError(t, err)

		instance2, err := CreateInstance("Remove Test Instance 2", product.ID)
		require.NoError(t, err)

		err = AddInstanceToDomain(domain.ID, instance1.ID)
		require.NoError(t, err)

		err = AddInstanceToDomain(domain.ID, instance2.ID)
		require.NoError(t, err)

		instances, err := GetInstancesByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 2)

		err = RemoveInstanceFromDomain(domain.ID, instance1.ID)
		require.NoError(t, err)

		instances, err = GetInstancesByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 1)
		assert.Equal(t, instance2.ID, instances[0].ID)

		domains, err := GetDomainsByInstance(instance1.ID)
		require.NoError(t, err)
		assert.Len(t, domains, 0)

		domains, err = GetDomainsByInstance(instance2.ID)
		require.NoError(t, err)
		assert.Len(t, domains, 1)
		assert.Equal(t, domain.ID, domains[0].ID)
	})

	t.Run("MultipleDomainsPerInstance", func(t *testing.T) {
		product, err := CreateProduct("Multi Domain Product", "Test Product Description")
		require.NoError(t, err)

		instance, err := CreateInstance("Multi Domain Instance", product.ID)
		require.NoError(t, err)

		var domains []*models.Domain
		for i := 0; i < 3; i++ {
			domain, err := CreateDomain(fmt.Sprintf("Domain %d", i+1), fmt.Sprintf("Description %d", i+1))
			require.NoError(t, err)
			domains = append(domains, domain)

			err = AddInstanceToDomain(domain.ID, instance.ID)
			require.NoError(t, err)
		}

		retrievedDomains, err := GetDomainsByInstance(instance.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedDomains, 3)

		domainMap := make(map[uuid.UUID]models.Domain)
		for _, dom := range retrievedDomains {
			domainMap[dom.ID] = dom
		}

		for _, created := range domains {
			retrieved, exists := domainMap[created.ID]
			assert.True(t, exists, "Created domain should contain the instance")
			assert.Equal(t, created.Name, retrieved.Name)
		}
	})

	t.Run("AddInstanceToDomain_NonExistentDomain", func(t *testing.T) {
		product, err := CreateProduct("Non-Existent Domain Product", "Test Product Description")
		require.NoError(t, err)

		instance, err := CreateInstance("Non-Existent Domain Instance", product.ID)
		require.NoError(t, err)

		nonExistentDomainID := uuid.New()
		err = AddInstanceToDomain(nonExistentDomainID, instance.ID)
		assert.Error(t, err)
	})

	t.Run("AddInstanceToDomain_NonExistentInstance", func(t *testing.T) {
		domain, err := CreateDomain("Non-Existent Instance Domain", "Test Description")
		require.NoError(t, err)

		nonExistentInstanceID := uuid.New()
		err = AddInstanceToDomain(domain.ID, nonExistentInstanceID)
		assert.Error(t, err)
	})
}

func TestGetInstancesByDomainIdWithThreatStats(t *testing.T) {

	// Create shared test entities at top level
	cleanupShared := testutil.SetupTestDatabase(t)
	defer cleanupShared()

	var err error
	var domain *models.Domain
	var product1, product2 *models.Product
	var instance1, instance2, instance3 *models.Instance
	var threat1, threat2, threat3, threat4, threat5 *models.Threat

	var setUp = func() {
		// Create test domain
		domain, err = CreateDomain("Test Domain for Stats", "Domain for threat stats testing")
		require.NoError(t, err)

		// Create test products
		product1, err = CreateProduct("Product 1", "First test product")
		require.NoError(t, err)

		product2, err = CreateProduct("Product 2", "Second test product")
		require.NoError(t, err)

		// Create test instances
		instance1, err = CreateInstance("Instance 1", product1.ID)
		require.NoError(t, err)

		instance2, err = CreateInstance("Instance 2", product1.ID)
		require.NoError(t, err)

		instance3, err = CreateInstance("Instance 3", product2.ID)
		require.NoError(t, err)

		// Add instances to domain
		err = AddInstanceToDomain(domain.ID, instance1.ID)
		require.NoError(t, err)
		err = AddInstanceToDomain(domain.ID, instance2.ID)
		require.NoError(t, err)
		err = AddInstanceToDomain(domain.ID, instance3.ID)
		require.NoError(t, err)

		// Create test threats
		threat1, err = CreateThreat("Threat 1", "First test threat")
		require.NoError(t, err)

		threat2, err = CreateThreat("Threat 2", "Second test threat")
		require.NoError(t, err)

		threat3, err = CreateThreat("Threat 3", "Third test threat")
		require.NoError(t, err)

		threat4, err = CreateThreat("Threat 4", "Fourth test threat")
		require.NoError(t, err)

		threat5, err = CreateThreat("Threat 5", "Fifth test threat")
		require.NoError(t, err)
	}

	t.Run("MixedInstanceAndProductThreats", func(t *testing.T) {
		setUp()
		// Instance1 has direct threat assignments:
		// - threat1: no resolution (unresolved)
		// - threat2: awaiting status (unresolved)
		// - threat3: resolved status (resolved - NOT counted)
		_, err := AssignThreatToInstance(instance1.ID, threat1.ID)
		require.NoError(t, err)

		assignment2, err := AssignThreatToInstance(instance1.ID, threat2.ID)
		require.NoError(t, err)

		assignment3, err := AssignThreatToInstance(instance1.ID, threat3.ID)
		require.NoError(t, err)

		// Create resolutions
		_, err = CreateThreatResolution(
			assignment2.ID,
			&instance1.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Awaiting resolution",
		)
		require.NoError(t, err)

		_, err = CreateThreatResolution(
			assignment3.ID,
			&instance1.ID,
			nil,
			models.ThreatAssignmentResolutionStatusResolved,
			"Resolved threat",
		)
		require.NoError(t, err)

		// Product1 threats (inherited by instance1 and instance2):
		// - threat4: no resolution (unresolved)
		// - threat5: accepted status (resolved - NOT counted)
		_, err = AssignThreatToProduct(product1.ID, threat4.ID)
		require.NoError(t, err)

		prodAssignment, err := AssignThreatToProduct(product1.ID, threat5.ID)
		require.NoError(t, err)

		_, err = CreateThreatResolution(
			prodAssignment.ID,
			nil,
			&product1.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Accepted risk",
		)
		require.NoError(t, err)

		// Instance3 has both instance and product threats:
		// - threat1: instance-level, no resolution (unresolved)
		// - threat4: product-level, no resolution (unresolved)
		_, err = AssignThreatToInstance(instance3.ID, threat1.ID)
		require.NoError(t, err)

		_, err = AssignThreatToProduct(product2.ID, threat4.ID)
		require.NoError(t, err)

		// Test the function
		instances, err := GetInstancesByDomainIdWithThreatStats(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 3)

		// Create a map for easier assertions
		instanceMap := make(map[uuid.UUID]models.InstanceWithThreatStats)
		for _, inst := range instances {
			instanceMap[inst.ID] = inst
		}

		// Instance1: 3 unresolved threats (threat1 + threat2 + threat4)
		inst1, exists := instanceMap[instance1.ID]
		assert.True(t, exists)
		assert.Equal(t, "Instance 1", inst1.Name)
		assert.Equal(t, 3, inst1.UnresolvedThreatCount)

		// Instance2: 1 unresolved threat (threat4 from product)
		inst2, exists := instanceMap[instance2.ID]
		assert.True(t, exists)
		assert.Equal(t, "Instance 2", inst2.Name)
		assert.Equal(t, 1, inst2.UnresolvedThreatCount)

		// Instance3: 2 unresolved threats (threat1 + threat4)
		inst3, exists := instanceMap[instance3.ID]
		assert.True(t, exists)
		assert.Equal(t, "Instance 3", inst3.Name)
		assert.Equal(t, 2, inst3.UnresolvedThreatCount)

		// Verify product information is properly loaded
		assert.Equal(t, product1.Name, inst1.Product.Name)
		assert.Equal(t, product1.Name, inst2.Product.Name)
		assert.Equal(t, product2.Name, inst3.Product.Name)
	})

	t.Run("EmptyDomain", func(t *testing.T) {
		// Test with domain that has no instances
		emptyDomain, err := CreateDomain("Empty Domain", "Domain with no instances")
		require.NoError(t, err)

		instances, err := GetInstancesByDomainIdWithThreatStats(emptyDomain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 0)
	})

	t.Run("NoThreats", func(t *testing.T) {
		setUp()
		// Test instances with no threat assignments
		instances, err := GetInstancesByDomainIdWithThreatStats(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 3)

		// All instances should have 0 unresolved threats
		for _, inst := range instances {
			assert.Equal(t, 0, inst.UnresolvedThreatCount)
		}
	})

	t.Run("AllThreatsResolved", func(t *testing.T) {
		// Assign threats and resolve all of them

		// Instance threat assignments
		instAssignment1, err := AssignThreatToInstance(instance1.ID, threat1.ID)
		require.NoError(t, err)
		instAssignment2, err := AssignThreatToInstance(instance2.ID, threat2.ID)
		require.NoError(t, err)

		// Product threat assignments
		prodAssignment1, err := AssignThreatToProduct(product1.ID, threat3.ID)
		require.NoError(t, err)
		prodAssignment2, err := AssignThreatToProduct(product2.ID, threat4.ID)
		require.NoError(t, err)

		// Resolve all threats with "resolved" or "accepted" status
		_, err = CreateThreatResolution(
			instAssignment1.ID,
			&instance1.ID,
			nil,
			models.ThreatAssignmentResolutionStatusResolved,
			"Instance threat resolved",
		)
		require.NoError(t, err)

		_, err = CreateThreatResolution(
			instAssignment2.ID,
			&instance2.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Instance threat accepted",
		)
		require.NoError(t, err)

		_, err = CreateThreatResolution(
			prodAssignment1.ID,
			nil,
			&product1.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Product threat resolved",
		)
		require.NoError(t, err)

		_, err = CreateThreatResolution(
			prodAssignment2.ID,
			nil,
			&product2.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Product threat accepted",
		)
		require.NoError(t, err)

		instances, err := GetInstancesByDomainIdWithThreatStats(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 3)

		// All instances should have 0 unresolved threats since all are resolved/accepted
		for _, inst := range instances {
			assert.Equal(t, 0, inst.UnresolvedThreatCount)
		}
	})

	t.Run("NonExistentDomain", func(t *testing.T) {
		nonExistentDomainID := uuid.New()
		instances, err := GetInstancesByDomainIdWithThreatStats(nonExistentDomainID)
		require.NoError(t, err)
		assert.Len(t, instances, 0)
	})

	t.Run("ComplexMixedResolutionStates", func(t *testing.T) {
		setUp()
		// Complex scenario testing various resolution states and inheritance patterns

		// Product1 threats (affect instance1 and instance2):
		// - threat1: resolved (should NOT count)
		// - threat2: accepted (should NOT count)
		prodAssignment1, err := AssignThreatToProduct(product1.ID, threat1.ID)
		require.NoError(t, err)
		prodAssignment2, err := AssignThreatToProduct(product1.ID, threat2.ID)
		require.NoError(t, err)

		// Product2 threats (affect instance3):
		// - threat3: no resolution (should count)
		_, err = AssignThreatToProduct(product2.ID, threat3.ID)
		require.NoError(t, err)

		// Instance-specific threats:
		// - instance1 + threat4: no resolution (should count)
		// - instance3 + threat5: awaiting (should count)
		_, err = AssignThreatToInstance(instance1.ID, threat4.ID)
		require.NoError(t, err)
		instAssignment, err := AssignThreatToInstance(instance3.ID, threat5.ID)
		require.NoError(t, err)

		// Create resolutions
		_, err = CreateThreatResolution(
			prodAssignment1.ID,
			nil,
			&product1.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Product1 threat1 resolved",
		)
		require.NoError(t, err)

		_, err = CreateThreatResolution(
			prodAssignment2.ID,
			nil,
			&product1.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Product1 threat2 accepted",
		)
		require.NoError(t, err)

		_, err = CreateThreatResolution(
			instAssignment.ID,
			&instance3.ID,
			nil,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Instance3 threat5 awaiting",
		)
		require.NoError(t, err)

		instances, err := GetInstancesByDomainIdWithThreatStats(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 3)

		instanceMap := make(map[uuid.UUID]models.InstanceWithThreatStats)
		for _, inst := range instances {
			instanceMap[inst.ID] = inst
		}

		// Instance1: 1 unresolved (threat4 instance-level)
		// Product threats (threat1, threat2) are resolved/accepted
		inst1, exists := instanceMap[instance1.ID]
		assert.True(t, exists)
		assert.Equal(t, 1, inst1.UnresolvedThreatCount)

		// Instance2: 0 unresolved
		// Product threats (threat1, threat2) are resolved/accepted, no instance-specific threats
		inst2, exists := instanceMap[instance2.ID]
		assert.True(t, exists)
		assert.Equal(t, 0, inst2.UnresolvedThreatCount)

		// Instance3: 2 unresolved
		// 1 product threat (threat3) + 1 instance threat (threat5 awaiting)
		inst3, exists := instanceMap[instance3.ID]
		assert.True(t, exists)
		assert.Equal(t, 2, inst3.UnresolvedThreatCount)
	})

}
