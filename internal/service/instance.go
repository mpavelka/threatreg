package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getInstanceRepository() (*models.InstanceRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewInstanceRepository(db), nil
}

func GetInstance(id uuid.UUID) (*models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.GetByID(nil, id)
}

func CreateInstance(name string, instanceOf uuid.UUID) (*models.Instance, error) {
	instance := &models.Instance{
		Name:       name,
		InstanceOf: instanceOf,
	}
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	err = instanceRepository.Create(nil, instance)
	if err != nil {
		fmt.Println("Error creating instance:", err)
		return nil, err
	}

	return instance, nil
}

func UpdateInstance(
	id uuid.UUID,
	name *string,
	instanceOf *uuid.UUID,
) (*models.Instance, error) {
	var updatedInstance *models.Instance
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		instanceRepository, err := getInstanceRepository()
		if err != nil {
			return err
		}
		instance, err := instanceRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// New values
		if name != nil {
			instance.Name = *name
		}
		if instanceOf != nil {
			instance.InstanceOf = *instanceOf
		}

		err = instanceRepository.Update(tx, instance)
		if err != nil {
			return err
		}

		// Reload to get updated Product relationship
		updatedInstance, err = instanceRepository.GetByID(tx, id)
		return err
	})

	return updatedInstance, err
}

func DeleteInstance(id uuid.UUID) error {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return err
	}

	return instanceRepository.Delete(nil, id)
}

func ListInstances() ([]models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.List(nil)
}

func ListInstancesByProductID(productID uuid.UUID) ([]models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.ListByProductID(nil, productID)
}

func FilterInstances(instanceName, productName string) ([]models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.Filter(nil, instanceName, productName)
}

func getThreatAssignmentRepositoryForInstance() (*models.ThreatAssignmentRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatAssignmentRepository(db), nil
}

func AssignThreatToInstance(instanceID, threatID uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepositoryForInstance()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.AssignThreatToInstance(nil, threatID, instanceID)
}
