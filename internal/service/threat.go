package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getThreatRepository() (*models.ThreatRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatRepository(db), nil
}

// GetThreat retrieves a threat by its unique identifier.
// Returns the threat if found, or an error if the threat does not exist or database access fails.
func GetThreat(id uuid.UUID) (*models.Threat, error) {
	threatRepository, err := getThreatRepository()
	if err != nil {
		return nil, err
	}

	return threatRepository.GetByID(nil, id)
}

// CreateThreat creates a new threat with the specified title and description.
// Returns the created threat with its assigned ID, or an error if creation fails.
func CreateThreat(title string, description string) (*models.Threat, error) {
	threat := &models.Threat{
		Title:       title,
		Description: description,
	}
	threatRepository, err := getThreatRepository()
	if err != nil {
		return nil, err
	}

	err = threatRepository.Create(nil, threat)
	if err != nil {
		fmt.Println("Error creating threat:", err)
		return nil, err
	}

	return threat, nil
}

// UpdateThreat updates an existing threat's title and/or description within a transaction.
// Only non-nil fields are updated. Returns the updated threat or an error if the update fails.
func UpdateThreat(
	id uuid.UUID,
	title *string,
	description *string,
) (*models.Threat, error) {
	var updatedThreat *models.Threat
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		threatRepository, err := getThreatRepository()
		if err != nil {
			return err
		}
		threat, err := threatRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// New values
		if title != nil {
			threat.Title = *title
		}
		if description != nil {
			threat.Description = *description
		}

		err = threatRepository.Update(tx, threat)
		if err != nil {
			return err
		}
		updatedThreat = threat
		return nil
	})

	return updatedThreat, err
}

// DeleteThreat removes a threat from the system by its unique identifier.
// Returns an error if the threat does not exist or if deletion fails.
func DeleteThreat(id uuid.UUID) error {
	threatRepository, err := getThreatRepository()
	if err != nil {
		return err
	}

	return threatRepository.Delete(nil, id)
}

// ListThreats retrieves all threats in the system.
// Returns a slice of threats or an error if database access fails.
func ListThreats() ([]models.Threat, error) {
	threatRepository, err := getThreatRepository()
	if err != nil {
		return nil, err
	}

	return threatRepository.List(nil)
}
