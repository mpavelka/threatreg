package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
)

func getComponentRepository() (*models.ComponentRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewComponentRepository(db), nil
}

// CreateComponent creates a new component with the provided details.
// Returns the created component or an error if creation fails.
func CreateComponent(name, description string, componentType models.ComponentType) (*models.Component, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	component := &models.Component{
		Name:        name,
		Description: description,
		Type:        componentType,
	}

	err = componentRepository.Create(nil, component)
	if err != nil {
		return nil, fmt.Errorf("error creating component: %w", err)
	}

	return component, nil
}

// GetComponent retrieves a component by its ID.
// Returns the component if found, or an error if the component does not exist or database access fails.
func GetComponent(id uuid.UUID) (*models.Component, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	return componentRepository.GetByID(nil, id)
}

// UpdateComponent updates an existing component with the provided fields.
// Only non-nil fields are updated. Returns the updated component or an error if update fails.
func UpdateComponent(id uuid.UUID, name, description *string) (*models.Component, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	component, err := componentRepository.GetByID(nil, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != nil {
		component.Name = *name
	}
	if description != nil {
		component.Description = *description
	}

	err = componentRepository.Update(nil, component)
	if err != nil {
		return nil, fmt.Errorf("error updating component: %w", err)
	}

	return component, nil
}

// DeleteComponent removes a component by its ID.
// Returns an error if deletion fails or the component does not exist.
func DeleteComponent(id uuid.UUID) error {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return err
	}

	return componentRepository.Delete(nil, id)
}

// ListComponents retrieves all components from the database.
// Returns a slice of components or an error if database access fails.
func ListComponents() ([]models.Component, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	return componentRepository.List(nil)
}

// ListComponentsByType retrieves all components of a specific type.
// Returns a slice of components or an error if database access fails.
func ListComponentsByType(componentType models.ComponentType) ([]models.Component, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	return componentRepository.ListByType(nil, componentType)
}

// FilterComponents searches for components by name.
// Returns a slice of matching components or an error if database access fails.
func FilterComponents(componentName string) ([]models.Component, error) {
	componentRepository, err := getComponentRepository()
	if err != nil {
		return nil, err
	}

	return componentRepository.Filter(nil, componentName)
}

// AssignThreatToComponent creates a threat assignment linking a threat to a specific component.
// Returns the created threat assignment or an error if the assignment already exists or creation fails.
func AssignThreatToComponent(componentID, threatID uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.AssignThreatToComponent(nil, threatID, componentID)
}

// ListThreatAssignmentsByComponentID retrieves all threat assignments for a specific component.
// Returns a slice of threat assignments or an error if database access fails.
func ListThreatAssignmentsByComponentID(componentID uuid.UUID) ([]models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.ListByComponentID(nil, componentID)
}

// ListThreatAssignmentsByComponentIDWithResolutionByComponentID retrieves threat assignments with resolution and delegation status for a component.
// The resolutionComponentID parameter filters which resolutions to include - only resolutions for that specific component will be shown.
// Returns a slice of threat assignments with resolution status or an error if database access fails.
func ListThreatAssignmentsByComponentIDWithResolutionByComponentID(componentID, resolutionComponentID uuid.UUID) ([]models.ThreatAssignmentWithResolution, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.ListWithResolutionByComponentID(nil, componentID, resolutionComponentID)
}