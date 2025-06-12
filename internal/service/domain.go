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

func GetDomain(id uuid.UUID) (*models.Domain, error) {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	return domainRepository.GetByID(nil, id)
}

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

func DeleteDomain(id uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.Delete(nil, id)
}

func ListDomains() ([]models.Domain, error) {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	return domainRepository.List(nil)
}

func GetInstancesByDomainId(domainID uuid.UUID) ([]models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.ListByDomainId(nil, domainID)
}

func AddInstanceToDomain(domainID, instanceID uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.AddInstance(nil, domainID, instanceID)
}

func RemoveInstanceFromDomain(domainID, instanceID uuid.UUID) error {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return err
	}

	return domainRepository.RemoveInstance(nil, domainID, instanceID)
}

func GetDomainsByInstance(instanceID uuid.UUID) ([]models.Domain, error) {
	domainRepository, err := getDomainRepository()
	if err != nil {
		return nil, err
	}

	return domainRepository.GetDomainsByInstanceID(nil, instanceID)
}
