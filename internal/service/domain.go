package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getDomainRepository() (*models.DomainRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewDomainRepository(db), nil
}

// GetDomain retrieves a domain by its unique identifier.
// Returns the domain if found, or an error if the domain does not exist or database access fails.
func GetDomain(id uuid.UUID) (*models.Domain, error) {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	return domainRepository.GetByID(nil, id)
}

// CreateDomain creates a new domain with the specified name and description.
// Returns the created domain with its assigned ID, or an error if creation fails.
func CreateDomain(
	Name string,
	Description string,
) (*models.Domain, error) {

	domain := &models.Domain{
		Name:        Name,
		Description: Description,
	}
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	err = domainRepository.Create(nil, domain)
	if err != nil {
		fmt.Println("Error creating domain:", err)
		return nil, err
	}

	return domain, nil
}

// UpdateDomain updates an existing domain's name and/or description within a transaction.
// Only non-nil fields are updated. Returns the updated domain or an error if the update fails.
func UpdateDomain(
	id uuid.UUID,
	Name *string,
	Description *string,
) (*models.Domain, error) {
	var updatedDomain *models.Domain
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		domainRepository, err := getDomainRepository()
		if err != nil {
			return err
		}
		domain, err := domainRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		if Name != nil {
			domain.Name = *Name
		}
		if Description != nil {
			domain.Description = *Description
		}

		err = domainRepository.Update(tx, domain)
		if err != nil {
			return err
		}
		updatedDomain = domain
		return nil
	})

	return updatedDomain, err
}

// DeleteDomain removes a domain from the system by its unique identifier.
// Returns an error if the domain does not exist or if deletion fails.
func DeleteDomain(id uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.Delete(nil, id)
}

// ListDomains retrieves all domains in the system.
// Returns a slice of domains or an error if database access fails.
func ListDomains() ([]models.Domain, error) {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	return domainRepository.List(nil)
}

// GetComponentsByDomainId retrieves all components that belong to a specific domain.
// Returns a slice of components with their information or an error if database access fails.
func GetComponentsByDomainId(domainID uuid.UUID) ([]models.Component, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	return componentRepository.ListByDomainId(nil, domainID)
}

// AddComponentToDomain associates a component with a domain.
// Returns an error if the domain or component does not exist, or if the association fails.
func AddComponentToDomain(domainID, componentID uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.AddComponent(nil, domainID, componentID)
}

// RemoveComponentFromDomain removes the association between a component and a domain.
// Returns an error if the association does not exist or if removal fails.
func RemoveComponentFromDomain(domainID, componentID uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.RemoveComponent(nil, domainID, componentID)
}

// GetDomainsByComponent retrieves all domains that contain a specific component.
// Returns a slice of domains or an error if database access fails.
func GetDomainsByComponent(componentID uuid.UUID) ([]models.Domain, error) {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	return domainRepository.GetDomainsByComponentID(nil, componentID)
}

// GetComponentsByDomainIdWithThreatStats retrieves all components in a domain with their threat statistics.
// Returns components with counts of total, resolved, and unresolved threats, or an error if database access fails.
func GetComponentsByDomainIdWithThreatStats(domainID uuid.UUID) ([]models.ComponentWithThreatStats, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	return componentRepository.ListByDomainIdWithThreatStats(nil, domainID)
}

// DEPRECATED FUNCTIONS - These redirect to component-based implementations for backward compatibility

// AddInstanceToDomain adds an instance (component) to a domain.
// DEPRECATED: Use AddComponentToDomain instead.
func AddInstanceToDomain(domainID, instanceID uuid.UUID) error {
	return AddComponentToDomain(domainID, instanceID)
}

// GetInstancesByDomainId retrieves all instances (components) in a domain.
// DEPRECATED: Use GetComponentsByDomainId instead.
func GetInstancesByDomainId(domainID uuid.UUID) ([]models.Component, error) {
	return GetComponentsByDomainId(domainID)
}

// GetInstancesByDomainIdWithThreatStats retrieves all instances (components) in a domain with threat stats.
// DEPRECATED: Use GetComponentsByDomainIdWithThreatStats instead.
func GetInstancesByDomainIdWithThreatStats(domainID uuid.UUID) ([]models.ComponentWithThreatStats, error) {
	return GetComponentsByDomainIdWithThreatStats(domainID)
}

// RemoveInstanceFromDomain removes an instance (component) from a domain.
// DEPRECATED: Use RemoveComponentFromDomain instead.
func RemoveInstanceFromDomain(domainID, instanceID uuid.UUID) error {
	return RemoveComponentFromDomain(domainID, instanceID)
}

// GetDomainsByInstance retrieves all domains that contain a specific instance (component).
// DEPRECATED: Use GetDomainsByComponent instead.
func GetDomainsByInstance(instanceID uuid.UUID) ([]models.Domain, error) {
	return GetDomainsByComponent(instanceID)
}
