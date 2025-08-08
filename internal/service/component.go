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

func getThreatAssignmentRepository() (*models.ThreatAssignmentRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatAssignmentRepository(db), nil
}

// AssignThreatToComponent creates a threat assignment linking a threat to a specific component.
// Returns the created threat assignment or an error if the assignment already exists or creation fails.
// This function also triggers automatic threat inheritance to descendant components.
// If the assignment has a Severity set, ResidualSeverity will be automatically set to match it.
func AssignThreatToComponent(componentID, threatID uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	assignment, err := threatAssignmentRepository.AssignThreatToComponent(nil, threatID, componentID)
	if err != nil {
		return nil, err
	}

	// If the assignment has a Severity but no ResidualSeverity, automatically set ResidualSeverity to match Severity
	if assignment.Severity != nil && assignment.ResidualSeverity == nil {
		err = threatAssignmentRepository.SetResidualSeverity(nil, assignment.ID, assignment.Severity)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-set residual severity on creation: %w", err)
		}
		// Update the assignment object to reflect the change
		assignment.ResidualSeverity = assignment.Severity
	}

	// Trigger threat inheritance for this assignment
	if err := onThreatAssignmentAssigned(assignment.ID, componentID); err != nil {
		// Log error but don't fail the main operation
		// In a production system, you might want to use a proper logger here
		fmt.Printf("Warning: Failed to process threat inheritance: %v\n", err)
	}

	return assignment, nil
}

// createThreatAssignmentWithoutInheritance creates a threat assignment without triggering inheritance
// This is used internally during inheritance processing to avoid infinite recursion
// If the assignment has a Severity set, ResidualSeverity will be automatically set to match it.
func createThreatAssignmentWithoutInheritance(componentID, threatID uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	assignment, err := threatAssignmentRepository.AssignThreatToComponent(nil, threatID, componentID)
	if err != nil {
		return nil, err
	}

	// If the assignment has a Severity but no ResidualSeverity, automatically set ResidualSeverity to match Severity
	if assignment.Severity != nil && assignment.ResidualSeverity == nil {
		err = threatAssignmentRepository.SetResidualSeverity(nil, assignment.ID, assignment.Severity)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-set residual severity on creation: %w", err)
		}
		// Update the assignment object to reflect the change
		assignment.ResidualSeverity = assignment.Severity
	}

	return assignment, nil
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

// GetThreatAssignmentById retrieves a threat assignment by its unique identifier.
// Returns the threat assignment if found, or an error if it does not exist or database access fails.
func GetThreatAssignmentById(id uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.GetByID(nil, id)
}

// SetThreatAssignmentSeverity updates the severity of a threat assignment.
// Returns an error if the assignment is not found, the severity is invalid, or database update fails.
// This function also triggers automatic threat inheritance if the assignment has a residual severity.
// If ResidualSeverity is nil, it will be automatically set to match the new Severity value.
func SetThreatAssignmentSeverity(id uuid.UUID, severity *models.ThreatSeverity) error {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return fmt.Errorf("error getting threat assignment repository: %w", err)
	}

	// Get the assignment to check its current state
	assignment, err := GetThreatAssignmentById(id)
	if err != nil {
		// If assignment doesn't exist, just try the update and return whatever error we get
		return threatAssignmentRepository.SetSeverity(nil, id, severity)
	}

	err = threatAssignmentRepository.SetSeverity(nil, id, severity)
	if err != nil {
		return err
	}

	// Auto-update ResidualSeverity logic:
	// 1. If ResidualSeverity is nil, set it to match the new Severity
	// 2. If ResidualSeverity currently matches the old Severity, update it to match the new Severity
	shouldUpdateResidualSeverity := false
	if assignment.ResidualSeverity == nil && severity != nil {
		shouldUpdateResidualSeverity = true
	} else if assignment.Severity != nil && assignment.ResidualSeverity != nil && severity != nil {
		// If ResidualSeverity currently matches the old Severity, update it to match the new Severity
		if *assignment.ResidualSeverity == *assignment.Severity {
			shouldUpdateResidualSeverity = true
		}
	}
	
	if shouldUpdateResidualSeverity {
		err = threatAssignmentRepository.SetResidualSeverity(nil, id, severity)
		if err != nil {
			return fmt.Errorf("failed to auto-set residual severity: %w", err)
		}
	}

	// Trigger threat inheritance after severity update
	if err := onThreatAssignmentAssigned(id, assignment.ComponentID); err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Warning: Failed to process threat inheritance after severity update: %v\n", err)
	}

	return nil
}

// SetThreatAssignmentResidualSeverity updates the residual severity of a threat assignment.
// Returns an error if the assignment is not found, the severity is invalid, or database update fails.
// This function also triggers automatic threat inheritance to propagate the new residual severity to descendants.
func SetThreatAssignmentResidualSeverity(id uuid.UUID, residualSeverity *models.ThreatSeverity) error {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return fmt.Errorf("error getting threat assignment repository: %w", err)
	}

	// Get the assignment to find its component ID for inheritance
	assignment, err := GetThreatAssignmentById(id)
	if err != nil {
		return fmt.Errorf("failed to get threat assignment for inheritance: %w", err)
	}

	err = threatAssignmentRepository.SetResidualSeverity(nil, id, residualSeverity)
	if err != nil {
		return err
	}

	// Trigger threat inheritance after residual severity update
	if err := onThreatAssignmentAssigned(id, assignment.ComponentID); err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Warning: Failed to process threat inheritance after residual severity update: %v\n", err)
	}

	return nil
}


// CreateInheritsThreatsRelationship creates a reserved relationship ReservedLabelInheritsThreatsFrom
// from childComponentId to parentComponentId.
func CreateInheritsThreatsRelationship(childComponentID, parentComponentID uuid.UUID) (*models.ComponentRelationship, error) {
	return createComponentRelationship(childComponentID, parentComponentID, string(models.ReservedLabelInheritsThreatsFrom), true)
}

// RemoveInheritsThreatsRelationship removes the reserved relationship ReservedLabelInheritsThreatsFrom
// from childComponentId to parentComponentId.
func RemoveInheritsThreatsRelationship(childComponentID, parentComponentID uuid.UUID) error {
	return DeleteComponentRelationshipByChildAndParent(childComponentID, parentComponentID)
}

func getThreatAssignmentRelationshipRepository() (*models.ThreatAssignmentRelationshipRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatAssignmentRelationshipRepository(db), nil
}

// onThreatAssignmentAssigned handles the automatic threat inheritance when a threat assignment is created or updated.
// This function propagates threat assignments to all descendant components that have an inherits-threats relationship.
func onThreatAssignmentAssigned(assignmentID, componentID uuid.UUID) error {
	// Get the updated/created assignment to access its ResidualSeverity
	parentAssignment, err := GetThreatAssignmentById(assignmentID)
	if err != nil {
		return fmt.Errorf("failed to get parent threat assignment: %w", err)
	}

	// Skip if parent has no residual severity to inherit
	if parentAssignment.ResidualSeverity == nil {
		return nil
	}

	// Get only descendant components that have the inherits-threats relationship
	relationshipRepository, err := getComponentRelationshipRepository()
	if err != nil {
		return err
	}

	// Only get descendants through the specific inherits-threats relationship
	descendantPaths, err := relationshipRepository.GetDescendantsOfComponentByLabel(
		nil, 
		componentID, 
		string(models.ReservedLabelInheritsThreatsFrom),
	)
	if err != nil {
		return fmt.Errorf("failed to get descendant components with inherits-threats relationship: %w", err)
	}

	// Process each descendant path recursively
	for _, path := range descendantPaths {
		err := processDescendantThreatInheritance(parentAssignment, path)
		if err != nil {
			return fmt.Errorf("failed to process descendant threat inheritance for component %s: %w", path.ComponentID, err)
		}
	}

	return nil
}

// processDescendantThreatInheritance processes threat inheritance for a single descendant path
// Note: descendant is already guaranteed to have inherits-threats relationship with parent
func processDescendantThreatInheritance(parentAssignment *models.ThreatAssignment, descendantPath models.ComponentTreePath) error {
	descendantComponentID := descendantPath.ComponentID

	// Get threat assignment relationship repository
	threatRelRepository, err := getThreatAssignmentRelationshipRepository()
	if err != nil {
		return err
	}

	// Find existing child assignment for this threat and component combination
	var childAssignment *models.ThreatAssignment
	threatAssignmentRepo, err := getThreatAssignmentRepository()
	if err != nil {
		return err
	}

	// Look for existing threat assignment on the descendant component for the same threat
	existingAssignments, err := threatAssignmentRepo.ListByComponentID(nil, descendantComponentID)
	if err != nil {
		return fmt.Errorf("failed to get threat assignments for descendant component: %w", err)
	}

	for _, assignment := range existingAssignments {
		if assignment.ThreatID == parentAssignment.ThreatID {
			// Check if this assignment has the ReservedInheritsFrom relationship to parent
			inheritanceRels, err := threatRelRepository.ListByFromAndToAndLabel(
				nil, 
				assignment.ID, 
				parentAssignment.ID, 
				string(models.ReservedInheritsFrom),
			)
			if err != nil {
				return fmt.Errorf("failed to check inheritance relationship: %w", err)
			}
			
			if len(inheritanceRels) > 0 {
				childAssignment = &assignment
				break
			}
		}
	}

	if childAssignment != nil {
		// Update existing child assignment's severity to parent's residual severity
		// Use direct repository access to avoid triggering inheritance recursion
		threatAssignmentRepo, err := getThreatAssignmentRepository()
		if err != nil {
			return fmt.Errorf("failed to get threat assignment repository: %w", err)
		}
		
		err = threatAssignmentRepo.SetSeverity(nil, childAssignment.ID, parentAssignment.ResidualSeverity)
		if err != nil {
			return fmt.Errorf("failed to update child threat assignment severity: %w", err)
		}
		
		// Also update ResidualSeverity if it was nil (auto-set behavior)
		if childAssignment.ResidualSeverity == nil && parentAssignment.ResidualSeverity != nil {
			err = threatAssignmentRepo.SetResidualSeverity(nil, childAssignment.ID, parentAssignment.ResidualSeverity)
			if err != nil {
				return fmt.Errorf("failed to auto-set child residual severity: %w", err)
			}
		}
		
		// Update the childAssignment object for recursive call
		childAssignment.Severity = parentAssignment.ResidualSeverity
		if childAssignment.ResidualSeverity == nil {
			childAssignment.ResidualSeverity = parentAssignment.ResidualSeverity
		}
	} else {
		// Create new child threat assignment without triggering inheritance to avoid recursion
		newChildAssignment, err := createThreatAssignmentWithoutInheritance(descendantComponentID, parentAssignment.ThreatID)
		if err != nil {
			return fmt.Errorf("failed to create child threat assignment: %w", err)
		}

		// Set severity to parent's residual severity
		err = SetThreatAssignmentSeverity(newChildAssignment.ID, parentAssignment.ResidualSeverity)
		if err != nil {
			return fmt.Errorf("failed to set child threat assignment severity: %w", err)
		}

		// Create the inheritance relationship
		relationship := &models.ThreatAssignmentRelationship{
			FromID: newChildAssignment.ID,
			ToID:   parentAssignment.ID,
			Label:  string(models.ReservedInheritsFrom),
		}

		err = threatRelRepository.Create(nil, relationship)
		if err != nil {
			return fmt.Errorf("failed to create threat assignment inheritance relationship: %w", err)
		}

		childAssignment = newChildAssignment
	}

	// Recursive call: propagate the child assignment to its descendants
	return onThreatAssignmentAssigned(childAssignment.ID, descendantComponentID)
}
