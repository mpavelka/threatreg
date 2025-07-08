package threat_pattern

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getThreatPatternRepository() (*models.ThreatPatternRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatPatternRepository(db), nil
}

func CreateThreatPattern(name, description string, threatID uuid.UUID, isActive bool) (*models.ThreatPattern, error) {
	pattern := &models.ThreatPattern{
		Name:        name,
		Description: description,
		ThreatID:    threatID,
		IsActive:    isActive,
	}

	var result *models.ThreatPattern
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		patternRepository, err := getThreatPatternRepository()
		if err != nil {
			return err
		}

		// Verify threat exists
		threatRepository, err := getThreatRepository()
		if err != nil {
			return err
		}
		_, err = threatRepository.GetByID(tx, threatID)
		if err != nil {
			return fmt.Errorf("threat not found: %w", err)
		}

		err = patternRepository.Create(tx, pattern)
		if err != nil {
			return fmt.Errorf("error creating threat pattern: %w", err)
		}

		result = pattern
		return nil
	})

	return result, err
}

func GetThreatPattern(id uuid.UUID) (*models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.GetByID(nil, id)
}

func UpdateThreatPattern(id uuid.UUID, name, description *string, threatID *uuid.UUID, isActive *bool) (*models.ThreatPattern, error) {
	var result *models.ThreatPattern
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		patternRepository, err := getThreatPatternRepository()
		if err != nil {
			return err
		}

		pattern, err := patternRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// Update fields if provided
		if name != nil {
			pattern.Name = *name
		}
		if description != nil {
			pattern.Description = *description
		}
		if threatID != nil {
			// Verify threat exists
			threatRepository, err := getThreatRepository()
			if err != nil {
				return err
			}
			_, err = threatRepository.GetByID(tx, *threatID)
			if err != nil {
				return fmt.Errorf("threat not found: %w", err)
			}
			pattern.ThreatID = *threatID
		}
		if isActive != nil {
			pattern.IsActive = *isActive
		}

		err = patternRepository.Update(tx, pattern)
		if err != nil {
			return err
		}

		// Reload the pattern to get the updated relationships
		updatedPattern, err := patternRepository.GetByID(tx, pattern.ID)
		if err != nil {
			return err
		}

		result = updatedPattern
		return nil
	})

	return result, err
}

func DeleteThreatPattern(id uuid.UUID) error {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return err
	}

	return patternRepository.Delete(nil, id)
}

func ListThreatPatterns() ([]models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.List(nil)
}

func ListActiveThreatPatterns() ([]models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.ListActive(nil)
}

func ListThreatPatternsByThreatID(threatID uuid.UUID) ([]models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.ListByThreatID(nil, threatID)
}

func SetThreatPatternActive(id uuid.UUID, isActive bool) error {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return err
	}

	return patternRepository.SetActive(nil, id, isActive)
}

func CreateThreatPatternWithConditions(name, description string, threatID uuid.UUID, isActive bool, conditions []models.PatternCondition) (*models.ThreatPattern, error) {
	var result *models.ThreatPattern
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create the pattern first
		pattern := &models.ThreatPattern{
			Name:        name,
			Description: description,
			ThreatID:    threatID,
			IsActive:    isActive,
		}

		patternRepository, err := getThreatPatternRepository()
		if err != nil {
			return err
		}

		// Verify threat exists
		threatRepository, err := getThreatRepository()
		if err != nil {
			return err
		}
		_, err = threatRepository.GetByID(tx, threatID)
		if err != nil {
			return fmt.Errorf("threat not found: %w", err)
		}

		err = patternRepository.Create(tx, pattern)
		if err != nil {
			return fmt.Errorf("error creating threat pattern: %w", err)
		}

		// Create conditions
		conditionRepository, err := getPatternConditionRepository()
		if err != nil {
			return err
		}

		for i := range conditions {
			conditions[i].PatternID = pattern.ID

			// Validate condition
			if err := validatePatternCondition(&conditions[i]); err != nil {
				return fmt.Errorf("invalid pattern condition %d: %w", i, err)
			}

			err = conditionRepository.Create(tx, &conditions[i])
			if err != nil {
				return fmt.Errorf("error creating pattern condition %d: %w", i, err)
			}
		}

		// Load the pattern with its conditions
		fullPattern, err := patternRepository.GetByID(tx, pattern.ID)
		if err != nil {
			return err
		}

		result = fullPattern
		return nil
	})

	return result, err
}

// Helper function to get threat repository - we need to import this from the main service package
func getThreatRepository() (*models.ThreatRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatRepository(db), nil
}
