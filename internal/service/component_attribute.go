package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
)

func getComponentAttributeRepository() (*models.ComponentAttributeRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewComponentAttributeRepository(db), nil
}

// CreateComponentAttribute creates a new component attribute with validation.
// For component type attributes, validates that the referenced component exists.
func CreateComponentAttribute(componentID uuid.UUID, name string, attributeType models.ComponentAttributeType, value string) (*models.ComponentAttribute, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	// Additional validation for component type - ensure referenced component exists
	if attributeType == models.ComponentAttributeTypeComponent {
		referencedComponentID, err := uuid.Parse(value)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID for component attribute: %w", err)
		}
		
		// Verify the referenced component exists
		_, err = GetComponent(referencedComponentID)
		if err != nil {
			return nil, fmt.Errorf("referenced component does not exist: %w", err)
		}
	}

	attribute := &models.ComponentAttribute{
		ComponentID: componentID,
		Name:        name,
		Type:        attributeType,
		Value:       value,
	}

	err = attributeRepository.Create(nil, attribute)
	if err != nil {
		return nil, fmt.Errorf("error creating component attribute: %w", err)
	}

	return attribute, nil
}

// GetComponentAttribute retrieves a component attribute by its ID.
func GetComponentAttribute(id uuid.UUID) (*models.ComponentAttribute, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	return attributeRepository.GetByID(nil, id)
}

// UpdateComponentAttribute updates an existing component attribute with validation.
// For component type attributes, validates that the referenced component exists.
func UpdateComponentAttribute(id uuid.UUID, name *string, attributeType *models.ComponentAttributeType, value *string) (*models.ComponentAttribute, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	attribute, err := attributeRepository.GetByID(nil, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != nil {
		attribute.Name = *name
	}
	if attributeType != nil {
		attribute.Type = *attributeType
	}
	if value != nil {
		attribute.Value = *value
	}

	// Additional validation for component type - ensure referenced component exists
	if attribute.Type == models.ComponentAttributeTypeComponent {
		referencedComponentID, err := uuid.Parse(attribute.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID for component attribute: %w", err)
		}
		
		// Verify the referenced component exists
		_, err = GetComponent(referencedComponentID)
		if err != nil {
			return nil, fmt.Errorf("referenced component does not exist: %w", err)
		}
	}

	err = attributeRepository.Update(nil, attribute)
	if err != nil {
		return nil, fmt.Errorf("error updating component attribute: %w", err)
	}

	return attribute, nil
}

// DeleteComponentAttribute removes a component attribute by its ID.
func DeleteComponentAttribute(id uuid.UUID) error {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return err
	}

	return attributeRepository.Delete(nil, id)
}

// GetComponentAttributes retrieves all attributes for a specific component.
func GetComponentAttributes(componentID uuid.UUID) ([]models.ComponentAttribute, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	return attributeRepository.ListByComponentID(nil, componentID)
}

// GetComponentAttributeByName retrieves a specific attribute of a component by name.
func GetComponentAttributeByName(componentID uuid.UUID, name string) (*models.ComponentAttribute, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	return attributeRepository.GetByComponentIDAndName(nil, componentID, name)
}

// FindComponentsByAttribute finds all components that have a specific attribute with a given value.
func FindComponentsByAttribute(name string, value string) ([]models.Component, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	return attributeRepository.FindComponentsByAttribute(nil, name, value)
}

// FindComponentsByAttributeAndType finds all components that have a specific attribute with a given value and type.
func FindComponentsByAttributeAndType(name string, value string, attributeType models.ComponentAttributeType) ([]models.Component, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	return attributeRepository.FindComponentsByAttributeAndType(nil, name, value, attributeType)
}

// ComponentHasAttribute checks if a component has a specific attribute (by name).
func ComponentHasAttribute(componentID uuid.UUID, name string) (bool, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return false, err
	}

	return attributeRepository.ComponentHasAttribute(nil, componentID, name)
}

// ComponentHasAttributeWithValue checks if a component has a specific attribute with a given value.
func ComponentHasAttributeWithValue(componentID uuid.UUID, name string, value string) (bool, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return false, err
	}

	return attributeRepository.ComponentHasAttributeWithValue(nil, componentID, name, value)
}

// DeleteComponentAttributeByName removes a component attribute by component ID and attribute name.
func DeleteComponentAttributeByName(componentID uuid.UUID, name string) error {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return err
	}

	return attributeRepository.DeleteByComponentIDAndName(nil, componentID, name)
}

// SetComponentAttribute creates or updates a component attribute.
// If an attribute with the same name already exists, it updates it; otherwise, it creates a new one.
func SetComponentAttribute(componentID uuid.UUID, name string, attributeType models.ComponentAttributeType, value string) (*models.ComponentAttribute, error) {
	attributeRepository, err := getComponentAttributeRepository()
	if err != nil {
		return nil, err
	}

	// Check if attribute already exists
	existingAttribute, err := attributeRepository.GetByComponentIDAndName(nil, componentID, name)
	if err == nil {
		// Attribute exists, update it
		return UpdateComponentAttribute(existingAttribute.ID, nil, &attributeType, &value)
	}

	// Attribute doesn't exist, create new one
	return CreateComponentAttribute(componentID, name, attributeType, value)
}