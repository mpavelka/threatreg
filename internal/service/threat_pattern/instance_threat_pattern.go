package threat_pattern

import (
	"fmt"
	"strings"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
)

// ThreatPatternMatch represents a matched threat pattern for an component
type ThreatPatternMatch struct {
	ComponentID uuid.UUID
	ThreatID    uuid.UUID
	PatternID   uuid.UUID
	Pattern     models.ThreatPattern
}

// GetAllComponentsThreatsByThreatPattern evaluates all components against threat patterns
// and returns a map of component IDs to their matching threat patterns
func GetAllComponentsThreatsByThreatPattern() (map[uuid.UUID][]ThreatPatternMatch, error) {
	// Get all components
	components, err := ListComponents()
	if err != nil {
		return nil, fmt.Errorf("error getting components: %w", err)
	}

	// Get all active threat patterns
	activePatterns, err := ListActiveThreatPatterns()
	if err != nil {
		return nil, fmt.Errorf("error getting active patterns: %w", err)
	}

	result := make(map[uuid.UUID][]ThreatPatternMatch)

	// Evaluate each component against all active patterns
	for _, component := range components {
		var matches []ThreatPatternMatch

		for _, pattern := range activePatterns {
			if evaluatePattern(component, pattern) {
				matches = append(matches, ThreatPatternMatch{
					ComponentID: component.ID,
					ThreatID:    pattern.ThreatID,
					PatternID:   pattern.ID,
					Pattern:     pattern,
				})
			}
		}

		if len(matches) > 0 {
			result[component.ID] = matches
		}
	}

	return result, nil
}

// GetComponentThreatsByThreatPattern evaluates a single component against a single threat pattern
// and returns the threat pattern matches if the component matches the pattern
func GetComponentThreatsByThreatPattern(component models.Component, pattern models.ThreatPattern) ([]ThreatPatternMatch, error) {
	var matches []ThreatPatternMatch

	if evaluatePattern(component, pattern) {
		matches = append(matches, ThreatPatternMatch{
			ComponentID: component.ID,
			ThreatID:    pattern.ThreatID,
			PatternID:   pattern.ID,
			Pattern:     pattern,
		})
	}

	return matches, nil
}

// GetComponentThreatsByExistingThreatPatterns evaluates a single component against all active threat patterns
// and returns all matching threat patterns for that component
func GetComponentThreatsByExistingThreatPatterns(component models.Component) ([]ThreatPatternMatch, error) {
	// Get all active threat patterns
	activePatterns, err := ListActiveThreatPatterns()
	if err != nil {
		return nil, fmt.Errorf("error getting active patterns: %w", err)
	}

	var matches []ThreatPatternMatch

	// Evaluate the component against all active patterns
	for _, pattern := range activePatterns {
		if evaluatePattern(component, pattern) {
			matches = append(matches, ThreatPatternMatch{
				ComponentID: component.ID,
				ThreatID:    pattern.ThreatID,
				PatternID:   pattern.ID,
				Pattern:     pattern,
			})
		}
	}

	return matches, nil
}

// evaluatePattern checks if an component matches a threat pattern
func evaluatePattern(component models.Component, pattern models.ThreatPattern) bool {
	if !pattern.IsActive {
		return false
	}

	// All conditions must be met for the pattern to match
	for _, condition := range pattern.Conditions {
		if !evaluateCondition(component, condition) {
			return false
		}
	}

	return true
}

// evaluateCondition evaluates a single pattern condition against an component
func evaluateCondition(component models.Component, condition models.PatternCondition) bool {
	conditionType, _ := models.ParsePatternConditionType(condition.ConditionType)
	operator, _ := models.ParsePatternOperator(condition.Operator)

	switch conditionType {
	case models.ConditionTypeTag:
		componentTags, err := getComponentTagNames(component.ID)
		if err != nil {
			return false
		}
		return evaluateTagCondition(componentTags, operator, condition.Value)

	case models.ConditionTypeRelationship:
		relationships, err := getComponentRelationships(component.ID)
		if err != nil {
			return false
		}
		return evaluateRelationshipCondition(relationships, operator, condition.RelationshipType, condition.Value)

	case models.ConditionTypeRelationshipTargetID:
		relationships, err := getComponentRelationships(component.ID)
		if err != nil {
			return false
		}
		return evaluateRelationshipTargetIDCondition(relationships, operator, condition.RelationshipType, condition.Value)

	case models.ConditionTypeRelationshipTargetTag:
		relationships, err := getComponentRelationships(component.ID)
		if err != nil {
			return false
		}
		return evaluateRelationshipTargetTagCondition(relationships, operator, condition.RelationshipType, condition.Value)

	default:
		return false
	}
}

func getComponentTagNames(componentID uuid.UUID) ([]string, error) {
	tags, err := service.ListTagsByComponentID(componentID)
	if err != nil {
		return nil, err
	}
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames, nil
}

func getComponentRelationships(componentID uuid.UUID) ([]models.Relationship, error) {
	relationships, err := service.ListRelationshipsByFromComponentID(componentID)
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
			// Check if target component has the specified tag
			var targetComponentID uuid.UUID
			if rel.ToComponentID != uuid.Nil {
				targetComponentID = rel.ToComponentID
			} else {
				continue // Skip if no target component
			}

			targetTags, err := getComponentTagNames(targetComponentID)
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
			if rel.ToComponentID != uuid.Nil && rel.ToComponentID.String() == targetID {
				return true
			}
		}
	}
	return false
}

// Helper functions to call service functions from the main service package

// ListComponents is a wrapper around the main service package function to list all components.
// Returns all components with their product information or an error if database access fails.
func ListComponents() ([]models.Component, error) {
	return service.ListComponents()
}

// GetProduct is a wrapper around the main service package function to retrieve a product by ID.
// Returns the product if found, or an error if the product does not exist or database access fails.
func GetProduct(productID uuid.UUID) (*models.Product, error) {
	return service.GetProduct(productID)
}

// ListTagsByProductID is a wrapper around the main service package function to list tags by product.
// Returns all tags assigned to the specified product or an error if database access fails.
func ListTagsByProductID(productID uuid.UUID) ([]models.Tag, error) {
	return service.ListTagsByProductID(productID)
}

// ListTagsByComponentID is a wrapper around the main service package function to list tags by component.
// Returns all tags assigned to the specified component or an error if database access fails.
func ListTagsByComponentID(componentID uuid.UUID) ([]models.Tag, error) {
	return service.ListTagsByComponentID(componentID)
}

// ListRelationshipsByFromComponentID is a wrapper around the main service package function to list relationships.
// Returns all relationships originating from the specified component or an error if database access fails.
func ListRelationshipsByFromComponentID(componentID uuid.UUID) ([]models.Relationship, error) {
	return service.ListRelationshipsByFromComponentID(componentID)
}
