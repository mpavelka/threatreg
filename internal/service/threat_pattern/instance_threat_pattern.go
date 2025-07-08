package threat_pattern

import (
	"fmt"
	"strings"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
)

// ThreatPatternMatch represents a matched threat pattern for an instance
type ThreatPatternMatch struct {
	InstanceID uuid.UUID
	ThreatID   uuid.UUID
	PatternID  uuid.UUID
	Pattern    models.ThreatPattern
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

// Helper functions to call service functions from the main service package
func ListInstances() ([]models.Instance, error) {
	return service.ListInstances()
}

func GetProduct(productID uuid.UUID) (*models.Product, error) {
	return service.GetProduct(productID)
}

func ListTagsByProductID(productID uuid.UUID) ([]models.Tag, error) {
	return service.ListTagsByProductID(productID)
}

func ListTagsByInstanceID(instanceID uuid.UUID) ([]models.Tag, error) {
	return service.ListTagsByInstanceID(instanceID)
}

func ListRelationshipsByFromInstanceID(instanceID uuid.UUID) ([]models.Relationship, error) {
	return service.ListRelationshipsByFromInstanceID(instanceID)
}
