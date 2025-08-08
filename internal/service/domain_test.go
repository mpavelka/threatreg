package service

import (
	"fmt"
	"testing"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDomainService_Integration(t *testing.T) {

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

	t.Run("AddComponentToDomain", func(t *testing.T) {
		domain, err := CreateDomain("Test Domain", "Test Description")
		require.NoError(t, err)

		component, err := CreateComponent("Test Component", "Test Component Description", models.ComponentTypeInstance)
		require.NoError(t, err)

		err = AddComponentToDomain(domain.ID, component.ID)
		require.NoError(t, err)

		components, err := GetComponentsByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, components, 1)
		assert.Equal(t, component.ID, components[0].ID)
		assert.Equal(t, component.Name, components[0].Name)

		domains, err := GetDomainsByComponent(component.ID)
		require.NoError(t, err)
		assert.Len(t, domains, 1)
		assert.Equal(t, domain.ID, domains[0].ID)
		assert.Equal(t, domain.Name, domains[0].Name)
	})

	t.Run("AddComponentToDomain_Multiple", func(t *testing.T) {
		domain, err := CreateDomain("Multi Test Domain", "Test Description")
		require.NoError(t, err)

		var components []*models.Component
		for i := 0; i < 3; i++ {
			component, err := CreateComponent(fmt.Sprintf("Test Component %d", i+1), fmt.Sprintf("Test Component %d Description", i+1), models.ComponentTypeInstance)
			require.NoError(t, err)
			components = append(components, component)

			err = AddComponentToDomain(domain.ID, component.ID)
			require.NoError(t, err)
		}

		retrievedComponents, err := GetComponentsByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedComponents, 3)

		componentMap := make(map[uuid.UUID]models.Component)
		for _, comp := range retrievedComponents {
			componentMap[comp.ID] = comp
		}

		for _, created := range components {
			retrieved, exists := componentMap[created.ID]
			assert.True(t, exists, "Created component should exist in domain")
			assert.Equal(t, created.Name, retrieved.Name)
		}
	})

	t.Run("RemoveComponentFromDomain", func(t *testing.T) {
		domain, err := CreateDomain("Remove Test Domain", "Test Description")
		require.NoError(t, err)

		component1, err := CreateComponent("Remove Test Component 1", "Test Component 1 Description", models.ComponentTypeInstance)
		require.NoError(t, err)

		component2, err := CreateComponent("Remove Test Component 2", "Test Component 2 Description", models.ComponentTypeInstance)
		require.NoError(t, err)

		err = AddComponentToDomain(domain.ID, component1.ID)
		require.NoError(t, err)

		err = AddComponentToDomain(domain.ID, component2.ID)
		require.NoError(t, err)

		components, err := GetComponentsByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, components, 2)

		err = RemoveComponentFromDomain(domain.ID, component1.ID)
		require.NoError(t, err)

		components, err = GetComponentsByDomainId(domain.ID)
		require.NoError(t, err)
		assert.Len(t, components, 1)
		assert.Equal(t, component2.ID, components[0].ID)

		domains, err := GetDomainsByComponent(component1.ID)
		require.NoError(t, err)
		assert.Len(t, domains, 0)

		domains, err = GetDomainsByComponent(component2.ID)
		require.NoError(t, err)
		assert.Len(t, domains, 1)
		assert.Equal(t, domain.ID, domains[0].ID)
	})

	t.Run("MultipleDomainsPerComponent", func(t *testing.T) {
		component, err := CreateComponent("Multi Domain Component", "Test Component Description", models.ComponentTypeInstance)
		require.NoError(t, err)

		var domains []*models.Domain
		for i := 0; i < 3; i++ {
			domain, err := CreateDomain(fmt.Sprintf("Domain %d", i+1), fmt.Sprintf("Description %d", i+1))
			require.NoError(t, err)
			domains = append(domains, domain)

			err = AddComponentToDomain(domain.ID, component.ID)
			require.NoError(t, err)
		}

		retrievedDomains, err := GetDomainsByComponent(component.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedDomains, 3)

		domainMap := make(map[uuid.UUID]models.Domain)
		for _, dom := range retrievedDomains {
			domainMap[dom.ID] = dom
		}

		for _, created := range domains {
			retrieved, exists := domainMap[created.ID]
			assert.True(t, exists, "Created domain should contain the component")
			assert.Equal(t, created.Name, retrieved.Name)
		}
	})

	t.Run("AddComponentToDomain_NonExistentDomain", func(t *testing.T) {
		component, err := CreateComponent("Non-Existent Domain Component", "Test Component Description", models.ComponentTypeInstance)
		require.NoError(t, err)

		nonExistentDomainID := uuid.New()
		err = AddComponentToDomain(nonExistentDomainID, component.ID)
		assert.Error(t, err)
	})

	t.Run("AddComponentToDomain_NonExistentComponent", func(t *testing.T) {
		domain, err := CreateDomain("Non-Existent Component Domain", "Test Description")
		require.NoError(t, err)

		nonExistentComponentID := uuid.New()
		err = AddComponentToDomain(domain.ID, nonExistentComponentID)
		assert.Error(t, err)
	})
}

func TestGetComponentsByDomainIdWithThreatStats(t *testing.T) {

	// Create shared test entities at top level

	var domain *models.Domain
	var component1, component2, component3 *models.Component
	var threat1, threat2, threat3, threat4, threat5 *models.Threat

	var setUp = func() {
		// Create test domain
		domain, _ = CreateDomain("Test Domain for Stats", "Domain for threat stats testing")
		// Create test components
		component1, _ = CreateComponent("Component 1", "First test component", models.ComponentTypeInstance)
		component2, _ = CreateComponent("Component 2", "Second test component", models.ComponentTypeInstance)
		component3, _ = CreateComponent("Component 3", "Third test component", models.ComponentTypeProduct)
		// Add components to domain
		AddComponentToDomain(domain.ID, component1.ID)
		AddComponentToDomain(domain.ID, component2.ID)
		AddComponentToDomain(domain.ID, component3.ID)
		// Create test threats
		threat1, _ = CreateThreat("Threat 1", "First test threat")
		threat2, _ = CreateThreat("Threat 2", "Second test threat")
		threat3, _ = CreateThreat("Threat 3", "Third test threat")
		threat4, _ = CreateThreat("Threat 4", "Fourth test threat")
		threat5, _ = CreateThreat("Threat 5", "Fifth test threat")
	}

	t.Run("MixedComponentThreats", func(t *testing.T) {
		setUp()
		// Component1 has direct threat assignments:
		// - threat1: no resolution (unresolved)
		// - threat2: awaiting status (unresolved)
		// - threat3: resolved status (resolved - NOT counted)
		AssignThreatToComponent(component1.ID, threat1.ID)
		assignment2, _ := AssignThreatToComponent(component1.ID, threat2.ID)
		assignment3, _ := AssignThreatToComponent(component1.ID, threat3.ID)

		// Create resolutions
		CreateThreatResolution(
			assignment2.ID,
			component1.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Awaiting resolution",
		)
		CreateThreatResolution(
			assignment3.ID,
			component1.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Resolved threat",
		)

		// Component2 and Component3 threats:
		// - threat4: assigned to component2, no resolution (unresolved)
		// - threat5: assigned to component3, accepted status (resolved - NOT counted)
		AssignThreatToComponent(component2.ID, threat4.ID)
		compAssignment, _ := AssignThreatToComponent(component3.ID, threat5.ID)
		CreateThreatResolution(
			compAssignment.ID,
			component3.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Accepted risk",
		)

		// Test the function
		components, err := GetComponentsByDomainIdWithThreatStats(domain.ID)
		require.NoError(t, err)
		assert.Len(t, components, 3)

		// Create a map for easier assertions
		componentMap := make(map[uuid.UUID]models.ComponentWithThreatStats)
		for _, comp := range components {
			componentMap[comp.ID] = comp
		}

		// Component1: 2 unresolved threats (threat1 + threat2)
		comp1, exists := componentMap[component1.ID]
		assert.True(t, exists)
		assert.Equal(t, "Component 1", comp1.Name)
		assert.Equal(t, 2, comp1.UnresolvedThreatCount)

		// Component2: 1 unresolved threat (threat4)
		comp2, exists := componentMap[component2.ID]
		assert.True(t, exists)
		assert.Equal(t, "Component 2", comp2.Name)
		assert.Equal(t, 1, comp2.UnresolvedThreatCount)

		// Component3: 0 unresolved threats (threat5 is accepted/resolved)
		comp3, exists := componentMap[component3.ID]
		assert.True(t, exists)
		assert.Equal(t, "Component 3", comp3.Name)
		assert.Equal(t, 0, comp3.UnresolvedThreatCount)
	})

	t.Run("EmptyDomain", func(t *testing.T) {
		// Test with domain that has no components
		emptyDomain, _ := CreateDomain("Empty Domain", "Domain with no components")
		components, _ := GetComponentsByDomainIdWithThreatStats(emptyDomain.ID)
		assert.Len(t, components, 0)
	})

	t.Run("NoThreats", func(t *testing.T) {
		setUp()
		// Test components with no threat assignments
		components, _ := GetComponentsByDomainIdWithThreatStats(domain.ID)
		assert.Len(t, components, 3)
		// All components should have 0 unresolved threats
		for _, comp := range components {
			assert.Equal(t, 0, comp.UnresolvedThreatCount)
		}
	})

	t.Run("AllThreatsResolved", func(t *testing.T) {
		// Assign threats and resolve all of them

		// Component threat assignments
		compAssignment1, _ := AssignThreatToComponent(component1.ID, threat1.ID)
		compAssignment2, _ := AssignThreatToComponent(component2.ID, threat2.ID)
		compAssignment3, _ := AssignThreatToComponent(component3.ID, threat3.ID)

		// Resolve all threats with "resolved" or "accepted" status
		CreateThreatResolution(
			compAssignment1.ID,
			component1.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Component threat resolved",
		)
		CreateThreatResolution(
			compAssignment2.ID,
			component2.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Component threat accepted",
		)
		CreateThreatResolution(
			compAssignment3.ID,
			component3.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Component threat resolved",
		)

		components, _ := GetComponentsByDomainIdWithThreatStats(domain.ID)
		assert.Len(t, components, 3)

		// All components should have 0 unresolved threats since all are resolved/accepted
		for _, comp := range components {
			assert.Equal(t, 0, comp.UnresolvedThreatCount)
		}
	})

	t.Run("NonExistentDomain", func(t *testing.T) {
		nonExistentDomainID := uuid.New()
		components, _ := GetComponentsByDomainIdWithThreatStats(nonExistentDomainID)
		assert.Len(t, components, 0)
	})

	t.Run("ComplexMixedResolutionStates", func(t *testing.T) {
		setUp()
		// Complex scenario testing various resolution states

		// Component threat assignments:
		// - component1 + threat1: resolved (should NOT count)
		// - component1 + threat2: accepted (should NOT count)
		// - component2 + threat3: no resolution (should count)
		// - component2 + threat4: awaiting (should count)
		// - component3 + threat5: no resolution (should count)
		compAssignment1, _ := AssignThreatToComponent(component1.ID, threat1.ID)
		compAssignment2, _ := AssignThreatToComponent(component1.ID, threat2.ID)
		_, _ = AssignThreatToComponent(component2.ID, threat3.ID)
		compAssignment4, _ := AssignThreatToComponent(component2.ID, threat4.ID)
		_, _ = AssignThreatToComponent(component3.ID, threat5.ID)

		// Create resolutions
		CreateThreatResolution(
			compAssignment1.ID,
			component1.ID,
			models.ThreatAssignmentResolutionStatusResolved,
			"Component1 threat1 resolved",
		)
		CreateThreatResolution(
			compAssignment2.ID,
			component1.ID,
			models.ThreatAssignmentResolutionStatusAccepted,
			"Component1 threat2 accepted",
		)
		CreateThreatResolution(
			compAssignment4.ID,
			component2.ID,
			models.ThreatAssignmentResolutionStatusAwaiting,
			"Component2 threat4 awaiting",
		)

		components, _ := GetComponentsByDomainIdWithThreatStats(domain.ID)
		assert.Len(t, components, 3)

		componentMap := make(map[uuid.UUID]models.ComponentWithThreatStats)
		for _, comp := range components {
			componentMap[comp.ID] = comp
		}

		// Component1: 0 unresolved (both threats resolved/accepted)
		comp1, exists := componentMap[component1.ID]
		assert.True(t, exists)
		assert.Equal(t, 0, comp1.UnresolvedThreatCount)

		// Component2: 2 unresolved (threat3 no resolution + threat4 awaiting)
		comp2, exists := componentMap[component2.ID]
		assert.True(t, exists)
		assert.Equal(t, 2, comp2.UnresolvedThreatCount)

		// Component3: 1 unresolved (threat5 no resolution)
		comp3, exists := componentMap[component3.ID]
		assert.True(t, exists)
		assert.Equal(t, 1, comp3.UnresolvedThreatCount)
	})

}
