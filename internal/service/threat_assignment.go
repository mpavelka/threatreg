package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"
)

func getThreatAssignmentRepository() (*models.ThreatAssignmentRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatAssignmentRepository(db), nil
}

// GetThreatAssignmentById retrieves a threat assignment by its unique identifier.
// Returns the threat assignment if found, or an error if the assignment does not exist or database access fails.
func GetThreatAssignmentById(id int) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.GetById(nil, id)
}