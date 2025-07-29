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

// GetInstancesByDomainId retrieves all instances that belong to a specific domain.
// Returns a slice of instances with their product information or an error if database access fails.
func GetInstancesByDomainId(domainID uuid.UUID) ([]models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.ListByDomainId(nil, domainID)
}

// AddInstanceToDomain associates an instance with a domain.
// Returns an error if the domain or instance does not exist, or if the association fails.
func AddInstanceToDomain(domainID, instanceID uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.AddInstance(nil, domainID, instanceID)
}

// RemoveInstanceFromDomain removes the association between an instance and a domain.
// Returns an error if the association does not exist or if removal fails.
func RemoveInstanceFromDomain(domainID, instanceID uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.RemoveInstance(nil, domainID, instanceID)
}

// GetDomainsByInstance retrieves all domains that contain a specific instance.
// Returns a slice of domains or an error if database access fails.
func GetDomainsByInstance(instanceID uuid.UUID) ([]models.Domain, error) {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	return domainRepository.GetDomainsByInstanceID(nil, instanceID)
}

// GetInstancesByDomainIdWithThreatStats retrieves all instances in a domain with their threat statistics.
// Returns instances with counts of total, resolved, and unresolved threats, or an error if database access fails.
func GetInstancesByDomainIdWithThreatStats(domainID uuid.UUID) ([]models.InstanceWithThreatStats, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.ListByDomainIdWithThreatStats(nil, domainID)
}
