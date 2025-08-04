package threat_pattern

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getPatternConditionRepository() (*models.PatternConditionRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewPatternConditionRepository(db), nil
}

// CreatePatternCondition creates a new pattern condition with validation of enum values.
// Validates that conditionType and operator are valid enum values. Returns the created condition or an error.
func CreatePatternCondition(patternID uuid.UUID, conditionType, operator, value, relationshipType string) (*models.PatternCondition, error) {
	condition := &models.PatternCondition{
		PatternID:        patternID,
		ConditionType:    conditionType,
		Operator:         operator,
		Value:            value,
		RelationshipType: relationshipType,
	}

	var result *models.PatternCondition
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		conditionRepository, err := getPatternConditionRepository()
		if err != nil {
			return err
		}

		// Verify pattern exists
		patternRepository, err := getThreatPatternRepository()
		if err != nil {
			return err
		}
		_, err = patternRepository.GetByID(tx, patternID)
		if err != nil {
			return fmt.Errorf("threat pattern not found: %w", err)
		}

		// Validate condition
		if err := validatePatternCondition(condition); err != nil {
			return fmt.Errorf("invalid pattern condition: %w", err)
		}

		err = conditionRepository.Create(tx, condition)
		if err != nil {
			return fmt.Errorf("error creating pattern condition: %w", err)
		}

		result = condition
		return nil
	})

	return result, err
}

// GetPatternCondition retrieves a pattern condition by its unique identifier.
// Returns the condition if found, or an error if the condition does not exist or database access fails.
func GetPatternCondition(id uuid.UUID) (*models.PatternCondition, error) {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return nil, err
	}

	return conditionRepository.GetByID(nil, id)
}

// UpdatePatternCondition updates an existing pattern condition with validation within a transaction.
// Only non-nil fields are updated. Validates enum values before updating. Returns the updated condition or an error.
func UpdatePatternCondition(id uuid.UUID, conditionType, operator, value, relationshipType *string) (*models.PatternCondition, error) {
	var result *models.PatternCondition
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		conditionRepository, err := getPatternConditionRepository()
		if err != nil {
			return err
		}

		condition, err := conditionRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// Update fields if provided
		if conditionType != nil {
			condition.ConditionType = *conditionType
		}
		if operator != nil {
			condition.Operator = *operator
		}
		if value != nil {
			condition.Value = *value
		}
		if relationshipType != nil {
			condition.RelationshipType = *relationshipType
		}

		// Validate updated condition
		if err := validatePatternCondition(condition); err != nil {
			return fmt.Errorf("invalid pattern condition: %w", err)
		}

		err = conditionRepository.Update(tx, condition)
		if err != nil {
			return err
		}

		result = condition
		return nil
	})

	return result, err
}

// DeletePatternCondition removes a pattern condition from the system by its unique identifier.
// Returns an error if the condition does not exist or if deletion fails.
func DeletePatternCondition(id uuid.UUID) error {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return err
	}

	return conditionRepository.Delete(nil, id)
}

// ListPatternConditionsByPatternID retrieves all conditions associated with a specific threat pattern.
// Returns a slice of conditions or an error if database access fails.
func ListPatternConditionsByPatternID(patternID uuid.UUID) ([]models.PatternCondition, error) {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return nil, err
	}

	return conditionRepository.ListByPatternID(nil, patternID)
}

// DeletePatternConditionsByPatternID removes all conditions associated with a specific threat pattern.
// Used when deleting a pattern to clean up its conditions. Returns an error if deletion fails.
func DeletePatternConditionsByPatternID(patternID uuid.UUID) error {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return err
	}

	return conditionRepository.DeleteByPatternID(nil, patternID)
}

// ListAllPatternConditions retrieves all pattern conditions in the system.
// Returns a slice of all conditions or an error if database access fails.
func ListAllPatternConditions() ([]models.PatternCondition, error) {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return nil, err
	}

	return conditionRepository.List(nil)
}

func validatePatternCondition(condition *models.PatternCondition) error {
	// Validate required fields
	if condition.ConditionType == "" {
		return fmt.Errorf("condition_type is required")
	}
	if condition.Operator == "" {
		return fmt.Errorf("operator is required")
	}

	// Validate condition type using enum
	conditionType, validConditionType := models.ParsePatternConditionType(condition.ConditionType)
	if !validConditionType {
		return fmt.Errorf("invalid condition_type: %s", condition.ConditionType)
	}

	// Validate operator using enum
	operator, validOperator := models.ParsePatternOperator(condition.Operator)
	if !validOperator {
		return fmt.Errorf("invalid operator: %s", condition.Operator)
	}

	// Validate condition-specific requirements
	relationshipConditionTypes := []models.PatternConditionType{
		models.ConditionTypeRelationshipTargetID,
		models.ConditionTypeRelationshipTargetTag,
		models.ConditionTypeRelationship,
	}
	for _, rct := range relationshipConditionTypes {
		if conditionType == rct && condition.RelationshipType == "" {
			return fmt.Errorf("relationship_type is required for %s condition", condition.ConditionType)
		}
	}

	// Most conditions require a value (except EXISTS/NOT_EXISTS)
	valueRequiredTypes := []models.PatternConditionType{
		models.ConditionTypeTag,
		models.ConditionTypeRelationshipTargetID,
		models.ConditionTypeRelationshipTargetTag,
	}
	for _, vrt := range valueRequiredTypes {
		if conditionType == vrt {
			if operator != models.OperatorExists && operator != models.OperatorNotExists && condition.Value == "" {
				return fmt.Errorf("value is required for %s condition with %s operator", condition.ConditionType, condition.Operator)
			}
		}
	}

	return nil
}
