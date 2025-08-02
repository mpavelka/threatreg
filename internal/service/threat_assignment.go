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
