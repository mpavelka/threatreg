package service

import (
	"fmt"
	"strings"
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

func getPatternConditionRepository() (*models.PatternConditionRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewPatternConditionRepository(db), nil
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

// GetInstanceThreatsByThreatPattern evaluates all instances against threat patterns
// and returns a map of instance IDs to their matching threat patterns
func GetInstanceThreatsByThreatPattern() (map[uuid.UUID][]ThreatPatternMatch, error) {
	// Get all instances
	instances, err := ListInstances()
	if err != nil {
		return nil, fmt.Errorf("error getting instances: %w", err)
	}

	// Get all active threat patterns
	activePatterns, err := ListActiveThreatPatterns()
	if err != nil {
		return nil, fmt.Errorf("error getting active patterns: %w", err)
	}

	result := make(map[uuid.UUID][]ThreatPatternMatch)

	// Evaluate each instance against all active patterns
	for _, instance := range instances {
		var matches []ThreatPatternMatch

		for _, pattern := range activePatterns {
			if evaluatePattern(instance, pattern) {
				matches = append(matches, ThreatPatternMatch{
					InstanceID: instance.ID,
					ThreatID:   pattern.ThreatID,
					PatternID:  pattern.ID,
					Pattern:    pattern,
				})
			}
		}

		if len(matches) > 0 {
			result[instance.ID] = matches
		}
	}

	return result, nil
}

// ThreatPatternMatch represents a matched threat pattern for an instance
type ThreatPatternMatch struct {
	InstanceID uuid.UUID
	ThreatID   uuid.UUID
	PatternID  uuid.UUID
	Pattern    models.ThreatPattern
}

// evaluatePattern checks if an instance matches a threat pattern
func evaluatePattern(instance models.Instance, pattern models.ThreatPattern) bool {
	if !pattern.IsActive {
		return false
	}

	// All conditions must be met for the pattern to match
	for _, condition := range pattern.Conditions {
		if !evaluateCondition(instance, condition) {
			return false
		}
	}

	return true
}

// evaluateCondition evaluates a single pattern condition against an instance
func evaluateCondition(instance models.Instance, condition models.PatternCondition) bool {
	conditionType, _ := models.ParsePatternConditionType(condition.ConditionType)
	operator, _ := models.ParsePatternOperator(condition.Operator)

	switch conditionType {
	case models.ConditionTypeProduct:
		productName, err := getProductName(instance.InstanceOf)
		if err != nil {
			return false
		}
		return applyOperator(productName, operator, condition.Value)

	case models.ConditionTypeProductID:
		return applyOperator(instance.InstanceOf.String(), operator, condition.Value)

	case models.ConditionTypeProductTag:
		productTags, err := getProductTagNames(instance.InstanceOf)
		if err != nil {
			return false
		}
		return evaluateTagCondition(productTags, operator, condition.Value)

	case models.ConditionTypeTag:
		instanceTags, err := getInstanceTagNames(instance.ID)
		if err != nil {
			return false
		}
		return evaluateTagCondition(instanceTags, operator, condition.Value)

	case models.ConditionTypeRelationship:
		relationships, err := getInstanceRelationships(instance.ID)
		if err != nil {
			return false
		}
		return evaluateRelationshipCondition(relationships, operator, condition.RelationshipType, condition.Value)

	case models.ConditionTypeRelationshipTargetID:
		relationships, err := getInstanceRelationships(instance.ID)
		if err != nil {
			return false
		}
		return evaluateRelationshipTargetIDCondition(relationships, operator, condition.RelationshipType, condition.Value)

	case models.ConditionTypeRelationshipTargetTag:
		relationships, err := getInstanceRelationships(instance.ID)
		if err != nil {
			return false
		}
		return evaluateRelationshipTargetTagCondition(relationships, operator, condition.RelationshipType, condition.Value)

	default:
		return false
	}
}

// Helper functions for getting data
func getProductName(productID uuid.UUID) (string, error) {
	product, err := GetProduct(productID)
	if err != nil {
		return "", err
	}
	return product.Name, nil
}

func getProductTagNames(productID uuid.UUID) ([]string, error) {
	tags, err := ListTagsByProductID(productID)
	if err != nil {
		return nil, err
	}
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames, nil
}

func getInstanceTagNames(instanceID uuid.UUID) ([]string, error) {
	tags, err := ListTagsByInstanceID(instanceID)
	if err != nil {
		return nil, err
	}
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames, nil
}

func getInstanceRelationships(instanceID uuid.UUID) ([]models.Relationship, error) {
	relationships, err := ListRelationshipsByFromInstanceID(instanceID)
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

// Condition evaluation helper functions
func evaluateTagCondition(tags []string, operator models.PatternOperator, value string) bool {
	switch operator {
	case models.OperatorContains:
		return containsString(tags, value)
	case models.OperatorNotContains:
		return !containsString(tags, value)
	case models.OperatorExists:
		return len(tags) > 0
	case models.OperatorNotExists:
		return len(tags) == 0
	default:
		return false
	}
}

func evaluateRelationshipCondition(relationships []models.Relationship, operator models.PatternOperator, relationshipType, value string) bool {
	switch operator {
	case models.OperatorExists:
		return hasRelationshipType(relationships, relationshipType)
	case models.OperatorNotExists:
		return !hasRelationshipType(relationships, relationshipType)
	case models.OperatorEquals:
		return hasRelationshipToTarget(relationships, relationshipType, value)
	default:
		return false
	}
}

func evaluateRelationshipTargetIDCondition(relationships []models.Relationship, operator models.PatternOperator, relationshipType, targetID string) bool {
	hasRelationship := hasRelationshipToTarget(relationships, relationshipType, targetID)
	switch operator {
	case models.OperatorHasRelationshipWith:
		return hasRelationship
	case models.OperatorNotHasRelationshipWith:
		return !hasRelationship
	default:
		return false
	}
}

func evaluateRelationshipTargetTagCondition(relationships []models.Relationship, operator models.PatternOperator, relationshipType, targetTag string) bool {
	for _, rel := range relationships {
		if rel.Type == relationshipType {
			// Check if target instance has the specified tag
			var targetInstanceID uuid.UUID
			if rel.ToInstanceID != nil {
				targetInstanceID = *rel.ToInstanceID
			} else {
				continue // Skip if no target instance
			}

			targetTags, err := getInstanceTagNames(targetInstanceID)
			if err != nil {
				continue
			}

			hasTag := containsString(targetTags, targetTag)
			switch operator {
			case models.OperatorHasRelationshipWith:
				if hasTag {
					return true
				}
			case models.OperatorNotHasRelationshipWith:
				if hasTag {
					return false
				}
			}
		}
	}

	// If we get here and operator is NOT_HAS_RELATIONSHIP_WITH, return true
	return operator == models.OperatorNotHasRelationshipWith
}

func applyOperator(actualValue string, operator models.PatternOperator, expectedValue string) bool {
	switch operator {
	case models.OperatorEquals:
		return actualValue == expectedValue
	case models.OperatorNotEquals:
		return actualValue != expectedValue
	case models.OperatorContains:
		return strings.Contains(actualValue, expectedValue)
	case models.OperatorNotContains:
		return !strings.Contains(actualValue, expectedValue)
	default:
		return false
	}
}

// Utility functions
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func hasRelationshipType(relationships []models.Relationship, relationshipType string) bool {
	for _, rel := range relationships {
		if rel.Type == relationshipType {
			return true
		}
	}
	return false
}

func hasRelationshipToTarget(relationships []models.Relationship, relationshipType, targetID string) bool {
	for _, rel := range relationships {
		if rel.Type == relationshipType {
			if rel.ToInstanceID != nil && rel.ToInstanceID.String() == targetID {
				return true
			}
			if rel.ToProductID != nil && rel.ToProductID.String() == targetID {
				return true
			}
		}
	}
	return false
}
