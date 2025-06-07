package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getControlRepository() (*models.ControlRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewControlRepository(db), nil
}

func GetControl(id uuid.UUID) (*models.Control, error) {
	controlRepository, err := getControlRepository()
	if err != nil {
		return nil, err
	}

	return controlRepository.GetByID(nil, id)
}

func CreateControl(
	Title string,
	Description string,
) (*models.Control, error) {

	control := &models.Control{
		Title:       Title,
		Description: Description,
	}
	controlRepository, err := getControlRepository()
	if err != nil {
		return nil, err
	}

	err = controlRepository.Create(nil, control)
	if err != nil {
		fmt.Println("Error creating control:", err)
		return nil, err
	}

	return control, nil
}

func UpdateControl(
	id uuid.UUID,
	Title *string,
	Description *string,
) (*models.Control, error) {
	var updatedControl *models.Control
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		controlRepository, err := getControlRepository()
		if err != nil {
			return err
		}
		control, err := controlRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// New values
		if Title != nil {
			control.Title = *Title
		}
		if Description != nil {
			control.Description = *Description
		}

		err = controlRepository.Update(tx, control)
		if err != nil {
			return err
		}
		updatedControl = control
		return nil
	})

	return updatedControl, err
}

func DeleteControl(id uuid.UUID) error {
	controlRepository, err := getControlRepository()
	if err != nil {
		return err
	}

	return controlRepository.Delete(nil, id)
}

func ListControls() ([]models.Control, error) {
	controlRepository, err := getControlRepository()
	if err != nil {
		return nil, err
	}

	return controlRepository.List(nil)
}
