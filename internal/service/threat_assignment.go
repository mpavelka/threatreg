package service

import (
	"threatreg/internal/models"
)

// GetThreatAssignmentById retrieves a threat assignment by its unique identifier.
// Returns the threat assignment if found, or an error if the assignment does not exist or database access fails.
func GetThreatAssignmentById(id int) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.GetByID(nil, id)
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
