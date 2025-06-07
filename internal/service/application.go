package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getApplicationRepository() (*models.ApplicationRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewApplicationRepository(db), nil
}

func GetApplication(id uuid.UUID) (*models.Application, error) {
	applicationRepository, err := getApplicationRepository()
	if err != nil {
		return nil, err
	}

	return applicationRepository.GetByID(nil, id)
}

func CreateApplication(name string, instanceOf uuid.UUID) (*models.Application, error) {
	application := &models.Application{
		Name:       name,
		InstanceOf: instanceOf,
	}
	applicationRepository, err := getApplicationRepository()
	if err != nil {
		return nil, err
	}

	err = applicationRepository.Create(nil, application)
	if err != nil {
		fmt.Println("Error creating application:", err)
		return nil, err
	}

	return application, nil
}

func UpdateApplication(
	id uuid.UUID,
	name *string,
	instanceOf *uuid.UUID,
) (*models.Application, error) {
	var updatedApplication *models.Application
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		applicationRepository, err := getApplicationRepository()
		if err != nil {
			return err
		}
		application, err := applicationRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// New values
		if name != nil {
			application.Name = *name
		}
		if instanceOf != nil {
			application.InstanceOf = *instanceOf
		}

		err = applicationRepository.Update(tx, application)
		if err != nil {
			return err
		}

		// Reload to get updated Product relationship
		updatedApplication, err = applicationRepository.GetByID(tx, id)
		return err
	})

	return updatedApplication, err
}

func DeleteApplication(id uuid.UUID) error {
	applicationRepository, err := getApplicationRepository()
	if err != nil {
		return err
	}

	return applicationRepository.Delete(nil, id)
}

func ListApplications() ([]models.Application, error) {
	applicationRepository, err := getApplicationRepository()
	if err != nil {
		return nil, err
	}

	return applicationRepository.List(nil)
}

func ListApplicationsByProductID(productID uuid.UUID) ([]models.Application, error) {
	applicationRepository, err := getApplicationRepository()
	if err != nil {
		return nil, err
	}

	return applicationRepository.ListByProductID(nil, productID)
}
