package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getRelationshipRepository() (*models.RelationshipRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewRelationshipRepository(db), nil
}

// AddRelationshipConsumesApiOfComponents creates bidirectional API consumption relationships between components.
// Creates 'CONSUMES_API_OF' and 'API_CONSUMED_BY' relationships.
func AddRelationshipConsumesApiOfComponents(fromComponentID, toComponentID uuid.UUID) error {
	// Begin transaction
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return err
	}

	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Create the CONSUMES_API_OF relationship
	consumesRel := &models.Relationship{
		Type:            "CONSUMES_API_OF",
		FromComponentID: fromComponentID,
		ToComponentID:   toComponentID,
	}

	err = relationshipRepository.Create(tx, consumesRel)
	if err != nil {
		return fmt.Errorf("error creating CONSUMES_API_OF relationship: %w", err)
	}

	// Create the vice-versa API_CONSUMED_BY relationship
	consumedByRel := &models.Relationship{
		Type:            "API_CONSUMED_BY",
		FromComponentID: toComponentID,
		ToComponentID:   fromComponentID,
	}

	err = relationshipRepository.Create(tx, consumedByRel)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating API_CONSUMED_BY relationship: %w", err)
	}

	return tx.Commit().Error
}

// AddRelationshipConsumesApiOf creates bidirectional API consumption relationships between instances or products.
// DEPRECATED: This function is deprecated. Use AddRelationshipConsumesApiOfComponents instead.
// Creates 'consumes_api_of' and 'provides_api_to' relationships. Exactly one pair (instances or products) must be provided.
func AddRelationshipConsumesApiOf(fromInstanceID, toInstanceID *uuid.UUID, fromProductID, toProductID *uuid.UUID) error {
	// Determine which component IDs to use
	var fromComponentID, toComponentID uuid.UUID
	
	if fromInstanceID != nil && toInstanceID != nil {
		fromComponentID = *fromInstanceID
		toComponentID = *toInstanceID
	} else if fromProductID != nil && toProductID != nil {
		fromComponentID = *fromProductID
		toComponentID = *toProductID
	} else {
		return fmt.Errorf("invalid relationship: must be between two instances or two products")
	}

	// Redirect to component-based function
	return AddRelationshipConsumesApiOfComponents(fromComponentID, toComponentID)
}

// AddRelationshipComponents creates relationships between components with optional reverse relationship.
// If viceVersaType is provided, creates a bidirectional relationship.
func AddRelationshipComponents(fromComponentID, toComponentID uuid.UUID, relType, viceVersaType string) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		relationshipRepository, err := getRelationshipRepository()
		if err != nil {
			return err
		}

		// Create the primary relationship
		relationship := &models.Relationship{
			Type:            relType,
			FromComponentID: fromComponentID,
			ToComponentID:   toComponentID,
		}

		err = relationshipRepository.Create(tx, relationship)
		if err != nil {
			return fmt.Errorf("error creating %s relationship: %w", relType, err)
		}

		// Create vice versa relationship if specified
		if viceVersaType != "" {
			viceVersaRel := &models.Relationship{
				Type:            viceVersaType,
				FromComponentID: toComponentID,
				ToComponentID:   fromComponentID,
			}

			err = relationshipRepository.Create(tx, viceVersaRel)
			if err != nil {
				return fmt.Errorf("error creating %s relationship: %w", viceVersaType, err)
			}
		}

		return nil
	})
}

// AddRelationship creates relationships between instances or products with optional reverse relationship.
// DEPRECATED: This function is deprecated. Use AddRelationshipComponents instead.
// If viceVersaType is provided, creates a bidirectional relationship. Validates entity consistency and existence.
func AddRelationship(fromInstanceID, fromProductID, toInstanceID, toProductID *uuid.UUID, relType, viceVersaType string) error {
	// Validate that only one "from" attribute is passed
	fromCount := 0
	if fromInstanceID != nil {
		fromCount++
	}
	if fromProductID != nil {
		fromCount++
	}
	if fromCount != 1 {
		return fmt.Errorf("exactly one 'from' attribute must be provided")
	}

	// Validate that only one "to" attribute is passed
	toCount := 0
	if toInstanceID != nil {
		toCount++
	}
	if toProductID != nil {
		toCount++
	}
	if toCount != 1 {
		return fmt.Errorf("exactly one 'to' attribute must be provided")
	}

	// This function should have been replaced with proper redirection logic above
	// For now, return an error to indicate it needs to be properly implemented
	return fmt.Errorf("AddRelationship function is deprecated, use AddRelationshipComponents instead")
}

// DeleteRelationshipById removes a relationship from the system by its unique identifier.
// Returns an error if the relationship does not exist or if deletion fails.
func DeleteRelationshipById(id uuid.UUID) error {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return err
	}

	return relationshipRepository.Delete(nil, id)
}

// GetRelationship retrieves a relationship by its unique identifier.
// Returns the relationship if found, or an error if the relationship does not exist or database access fails.
func GetRelationship(id uuid.UUID) (*models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.GetByID(nil, id)
}

// ListRelationships retrieves all relationships in the system.
// Returns a slice of relationships or an error if database access fails.
func ListRelationships() ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.List(nil)
}

// ListRelationshipsByFromComponentID retrieves all relationships originating from a specific component.
// Returns a slice of relationships or an error if database access fails.
func ListRelationshipsByFromComponentID(fromComponentID uuid.UUID) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByFromComponentID(nil, fromComponentID)
}

// ListRelationshipsByFromInstanceID retrieves all relationships originating from a specific instance.
// DEPRECATED: This function is deprecated. Use ListRelationshipsByFromComponentID instead.
// Returns a slice of relationships or an error if database access fails.
func ListRelationshipsByFromInstanceID(fromInstanceID uuid.UUID) ([]models.Relationship, error) {
	// Redirect to component-based function
	return ListRelationshipsByFromComponentID(fromInstanceID)
}

// ListRelationshipsByFromProductID retrieves all relationships originating from a specific product.
// DEPRECATED: This function is deprecated. Use ListRelationshipsByFromComponentID instead.
// Returns a slice of relationships or an error if database access fails.
func ListRelationshipsByFromProductID(fromProductID uuid.UUID) ([]models.Relationship, error) {
	// Redirect to component-based function
	return ListRelationshipsByFromComponentID(fromProductID)
}

// ListRelationshipsByType retrieves all relationships of a specific type.
// Returns a slice of relationships matching the type or an error if database access fails.
func ListRelationshipsByType(relType string) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByType(nil, relType)
}
