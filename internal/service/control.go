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

// GetControl retrieves a control by its unique identifier.
// Returns the control if found, or an error if the control does not exist or database access fails.
func GetControl(id uuid.UUID) (*models.Control, error) {
	controlRepository, err := getControlRepository()
	if err != nil {
		return nil, err
	}

	return controlRepository.GetByID(nil, id)
}

// CreateControl creates a new security control with the specified title and description.
// Returns the created control with its assigned ID, or an error if creation fails.
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

// UpdateControl updates an existing control's title and/or description within a transaction.
// Only non-nil fields are updated. Returns the updated control or an error if the update fails.
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

// DeleteControl removes a control from the system by its unique identifier.
// Returns an error if the control does not exist or if deletion fails.
func DeleteControl(id uuid.UUID) error {
	controlRepository, err := getControlRepository()
	if err != nil {
		return err
	}

	return controlRepository.Delete(nil, id)
}

// ListControls retrieves all security controls in the system.
// Returns a slice of controls or an error if database access fails.
func ListControls() ([]models.Control, error) {
	controlRepository, err := getControlRepository()
	if err != nil {
		return nil, err
	}

	return controlRepository.List(nil)
}
