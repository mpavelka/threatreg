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

func getThreatAssignmentResolutionRepository() (*models.ThreatAssignmentResolutionRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatAssignmentResolutionRepository(db), nil
}

func CreateThreatResolution(
	threatAssignmentID int,
	instanceID *uuid.UUID,
	productID *uuid.UUID,
	status models.ThreatAssignmentResolutionStatus,
	description string,
) (*models.ThreatAssignmentResolution, error) {
	resolution := &models.ThreatAssignmentResolution{
		ThreatAssignmentID: threatAssignmentID,
		Status:             status,
		Description:        description,
	}

	// Set exactly one of InstanceID or ProductID
	if instanceID != nil {
		resolution.InstanceID = *instanceID
		resolution.ProductID = uuid.Nil
	} else if productID != nil {
		resolution.ProductID = *productID
		resolution.InstanceID = uuid.Nil
	} else {
		return nil, fmt.Errorf("either instanceID or productID must be provided")
	}

	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	err = resolutionRepository.Create(nil, resolution)
	if err != nil {
		return nil, fmt.Errorf("error creating threat resolution: %w", err)
	}

	return resolution, nil
}

func UpdateThreatResolution(
	id uuid.UUID,
	status *models.ThreatAssignmentResolutionStatus,
	description *string,
) (*models.ThreatAssignmentResolution, error) {
	var updatedResolution *models.ThreatAssignmentResolution
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		resolutionRepository, err := getThreatAssignmentResolutionRepository()
		if err != nil {
			return err
		}

		resolution, err := resolutionRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// Update fields if provided
		if status != nil {
			resolution.Status = *status
		}
		if description != nil {
			resolution.Description = *description
		}

		err = resolutionRepository.Update(tx, resolution)
		if err != nil {
			return err
		}
		updatedResolution = resolution
		return nil
	})

	return updatedResolution, err
}

func GetThreatResolution(id uuid.UUID) (*models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.GetByID(nil, id)
}

func GetThreatResolutionByThreatAssignmentID(threatAssignmentID int) (*models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.GetByThreatAssignmentID(nil, threatAssignmentID)
}

func ListThreatResolutionsByProductID(productID uuid.UUID) ([]models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.ListByProductID(nil, productID)
}

func ListThreatResolutionsByInstanceID(instanceID uuid.UUID) ([]models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.ListByInstanceID(nil, instanceID)
}

func DeleteThreatResolution(id uuid.UUID) error {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return err
	}

	return resolutionRepository.Delete(nil, id)
}
