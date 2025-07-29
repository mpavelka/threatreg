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

// GetInstance retrieves an instance by its unique identifier.
// Returns the instance with its associated product information, or an error if not found.
func GetInstance(id uuid.UUID) (*models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.GetByID(nil, id)
}

// CreateInstance creates a new instance with the specified name and product reference.
// The instanceOf parameter must reference an existing product. Returns the created instance or an error.
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

// UpdateInstance updates an existing instance's name and/or product reference within a transaction.
// Only non-nil fields are updated. Returns the updated instance with refreshed relationships or an error.
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

// DeleteInstance removes an instance from the system by its unique identifier.
// Returns an error if the instance does not exist or if deletion fails.
func DeleteInstance(id uuid.UUID) error {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return err
	}

	return instanceRepository.Delete(nil, id)
}

// ListInstances retrieves all instances in the system with their associated product information.
// Returns a slice of instances or an error if database access fails.
func ListInstances() ([]models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.List(nil)
}

// ListInstancesByProductID retrieves all instances that belong to a specific product.
// Returns a slice of instances for the given product or an error if database access fails.
func ListInstancesByProductID(productID uuid.UUID) ([]models.Instance, error) {
	instanceRepository, err := getInstanceRepository()
	if err != nil {
		return nil, err
	}

	return instanceRepository.ListByProductID(nil, productID)
}

// FilterInstances searches for instances by name and/or product name using case-insensitive partial matching.
// Both parameters support partial matches. Returns matching instances or an error if database access fails.
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

// AssignThreatToInstance creates a threat assignment linking a threat to a specific instance.
// Returns the created threat assignment or an error if the assignment already exists or creation fails.
func AssignThreatToInstance(instanceID, threatID uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepositoryForInstance()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.AssignThreatToInstance(nil, threatID, instanceID)
}

// ListThreatAssignmentsByInstanceID retrieves all threat assignments for a specific instance.
// Returns a slice of threat assignments with threat details or an error if database access fails.
func ListThreatAssignmentsByInstanceID(instanceID uuid.UUID) ([]models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepositoryForInstance()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.ListByInstanceID(nil, instanceID)
}
