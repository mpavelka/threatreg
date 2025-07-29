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

// CreateThreatPattern creates a new threat pattern with validation that the associated threat exists.
// Returns the created pattern with its assigned ID, or an error if creation fails or the threat doesn't exist.
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

// GetThreatPattern retrieves a threat pattern by its unique identifier.
// Returns the pattern with its conditions, or an error if the pattern does not exist or database access fails.
func GetThreatPattern(id uuid.UUID) (*models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.GetByID(nil, id)
}

// UpdateThreatPattern updates an existing threat pattern's fields within a transaction.
// Only non-nil fields are updated. Validates threat existence if threatID is provided. Returns the updated pattern.
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

// DeleteThreatPattern removes a threat pattern from the system by its unique identifier.
// Returns an error if the pattern does not exist or if deletion fails.
func DeleteThreatPattern(id uuid.UUID) error {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return err
	}

	return patternRepository.Delete(nil, id)
}

// ListThreatPatterns retrieves all threat patterns in the system.
// Returns a slice of patterns with their conditions, or an error if database access fails.
func ListThreatPatterns() ([]models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.List(nil)
}

// ListActiveThreatPatterns retrieves only active threat patterns in the system.
// Returns a slice of active patterns with their conditions, or an error if database access fails.
func ListActiveThreatPatterns() ([]models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.ListActive(nil)
}

// ListThreatPatternsByThreatID retrieves all threat patterns associated with a specific threat.
// Returns a slice of patterns for the given threat, or an error if database access fails.
func ListThreatPatternsByThreatID(threatID uuid.UUID) ([]models.ThreatPattern, error) {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return nil, err
	}

	return patternRepository.ListByThreatID(nil, threatID)
}

// SetThreatPatternActive updates the active status of a threat pattern.
// Returns an error if the pattern does not exist or if the update fails.
func SetThreatPatternActive(id uuid.UUID, isActive bool) error {
	patternRepository, err := getThreatPatternRepository()
	if err != nil {
		return err
	}

	return patternRepository.SetActive(nil, id, isActive)
}

// CreateThreatPatternWithConditions creates a threat pattern and its conditions in a single transaction.
// Validates all conditions and creates them atomically with the pattern. Returns the created pattern or an error.
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
