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

func GetThreat(id uuid.UUID) (*models.Threat, error) {
	threatRepository, err := getThreatRepository()
	if err != nil {
		return nil, err
	}

	return threatRepository.GetByID(nil, id)
}

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

func DeleteThreat(id uuid.UUID) error {
	threatRepository, err := getThreatRepository()
	if err != nil {
		return err
	}

	return threatRepository.Delete(nil, id)
}

func ListThreats() ([]models.Threat, error) {
	threatRepository, err := getThreatRepository()
	if err != nil {
		return nil, err
	}

	return threatRepository.List(nil)
}