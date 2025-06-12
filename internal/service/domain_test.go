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

		instances, err := GetInstancesByDomain(domain.ID)
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

		retrievedInstances, err := GetInstancesByDomain(domain.ID)
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

		instances, err := GetInstancesByDomain(domain.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 2)

		err = RemoveInstanceFromDomain(domain.ID, instance1.ID)
		require.NoError(t, err)

		instances, err = GetInstancesByDomain(domain.ID)
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
