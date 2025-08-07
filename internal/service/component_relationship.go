package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
)

func getComponentRelationshipRepository() (*models.ComponentRelationshipRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewComponentRelationshipRepository(db), nil
}

// CreateComponentRelationship creates a new relationship between components.
// Validates that both components exist and that the label is valid before creating the relationship.
func CreateComponentRelationship(fromComponentID, toComponentID uuid.UUID, label string) (*models.ComponentRelationship, error) {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return nil, err
	}

	// Validate label first
	if err := models.ValidateLabel(label); err != nil {
		return nil, err
	}

	// Validate that both components exist
	_, err = GetComponent(fromComponentID)
	if err != nil {
		return nil, fmt.Errorf("from component does not exist: %w", err)
	}

	_, err = GetComponent(toComponentID)
	if err != nil {
		return nil, fmt.Errorf("to component does not exist: %w", err)
	}

	relationship := &models.ComponentRelationship{
		FromID: fromComponentID,
		ToID:   toComponentID,
		Label:  label,
	}

	err = relationshipRepository.Create(nil, relationship)
	if err != nil {
		return nil, fmt.Errorf("error creating component relationship: %w", err)
	}

	return relationship, nil
}

// GetComponentRelationship retrieves a component relationship relationship by its ID.
func GetComponentRelationship(id uuid.UUID) (*models.ComponentRelationship, error) {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.GetByID(nil, id)
}

// DeleteComponentRelationship removes a component relationship relationship by its ID.
func DeleteComponentRelationship(id uuid.UUID) error {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return err
	}

	return relationshipRepository.Delete(nil, id)
}

// GetComponentRelationshipByChildAndParent retrieves a specific inheritance relationship.
func GetComponentRelationshipByChildAndParent(childID, parentID uuid.UUID) (*models.ComponentRelationship, error) {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.GetByFromAndTo(nil, childID, parentID)
}

// ListComponentParents retrieves all parent components for a given child component.
func ListComponentParents(childID uuid.UUID) ([]models.ComponentRelationship, error) {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByFrom(nil, childID)
}

// ListComponentChildren retrieves all child components for a given parent component.
func ListComponentChildren(parentID uuid.UUID) ([]models.ComponentRelationship, error) {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByTo(nil, parentID)
}

// DeleteComponentRelationshipByChildAndParent removes a specific inheritance relationship.
func DeleteComponentRelationshipByChildAndParent(childID, parentID uuid.UUID) error {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return err
	}

	return relationshipRepository.DeleteByFromAndTo(nil, childID, parentID)
}

// GetComponentTreePaths retrieves all tree paths for a given component.
// This includes all paths from root ancestors to descendants that pass through the component.
func GetComponentTreePaths(componentID uuid.UUID) ([]models.ComponentTreePath, error) {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.GetTreePaths(nil, componentID)
}

// GetAllComponentTreePaths retrieves tree paths for all components in the system.
func GetAllComponentTreePaths() ([]models.ComponentTreePath, error) {
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.GetAllTreePaths(nil)
}
