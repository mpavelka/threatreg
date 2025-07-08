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

func GetPatternCondition(id uuid.UUID) (*models.PatternCondition, error) {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return nil, err
	}

	return conditionRepository.GetByID(nil, id)
}

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

func DeletePatternCondition(id uuid.UUID) error {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return err
	}

	return conditionRepository.Delete(nil, id)
}

func ListPatternConditionsByPatternID(patternID uuid.UUID) ([]models.PatternCondition, error) {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return nil, err
	}

	return conditionRepository.ListByPatternID(nil, patternID)
}

func DeletePatternConditionsByPatternID(patternID uuid.UUID) error {
	conditionRepository, err := getPatternConditionRepository()
	if err != nil {
		return err
	}

	return conditionRepository.DeleteByPatternID(nil, patternID)
}

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
		models.ConditionTypeProduct,
		models.ConditionTypeTag,
		models.ConditionTypeProductTag,
		models.ConditionTypeProductID,
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
