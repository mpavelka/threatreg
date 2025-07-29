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

// AddRelationshipConsumesApiOf creates bidirectional API consumption relationships between instances or products.
// Creates 'consumes_api_of' and 'provides_api_to' relationships. Exactly one pair (instances or products) must be provided.
func AddRelationshipConsumesApiOf(fromInstanceID, toInstanceID *uuid.UUID, fromProductID, toProductID *uuid.UUID) error {
	// Validate that relationship is between two instances or two products, but not mixed
	if fromInstanceID != nil && toInstanceID != nil {
		// Instance to instance relationship
		if fromProductID != nil || toProductID != nil {
			return fmt.Errorf("cannot mix instance and product relationships")
		}
	} else if fromProductID != nil && toProductID != nil {
		// Product to product relationship
		if fromInstanceID != nil || toInstanceID != nil {
			return fmt.Errorf("cannot mix instance and product relationships")
		}
	} else {
		return fmt.Errorf("invalid relationship: must be between two instances or two products")
	}

	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		relationshipRepository, err := getRelationshipRepository()
		if err != nil {
			return err
		}

		// Create the CONSUMES_API_OF relationship
		consumesRel := &models.Relationship{
			Type:           "CONSUMES_API_OF",
			FromInstanceID: fromInstanceID,
			FromProductID:  fromProductID,
			ToInstanceID:   toInstanceID,
			ToProductID:    toProductID,
		}

		err = relationshipRepository.Create(tx, consumesRel)
		if err != nil {
			return fmt.Errorf("error creating CONSUMES_API_OF relationship: %w", err)
		}

		// Create the vice-versa API_CONSUMED_BY relationship
		consumedByRel := &models.Relationship{
			Type:           "API_CONSUMED_BY",
			FromInstanceID: toInstanceID,
			FromProductID:  toProductID,
			ToInstanceID:   fromInstanceID,
			ToProductID:    fromProductID,
		}

		err = relationshipRepository.Create(tx, consumedByRel)
		if err != nil {
			return fmt.Errorf("error creating API_CONSUMED_BY relationship: %w", err)
		}

		return nil
	})
}

// AddRelationship creates relationships between instances or products with optional reverse relationship.
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

	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		relationshipRepository, err := getRelationshipRepository()
		if err != nil {
			return err
		}

		// Create the primary relationship
		primaryRel := &models.Relationship{
			Type:           relType,
			FromInstanceID: fromInstanceID,
			FromProductID:  fromProductID,
			ToInstanceID:   toInstanceID,
			ToProductID:    toProductID,
		}

		err = relationshipRepository.Create(tx, primaryRel)
		if err != nil {
			return fmt.Errorf("error creating primary relationship: %w", err)
		}

		// Create the vice-versa relationship if viceVersaType is provided
		if viceVersaType != "" {
			viceVersaRel := &models.Relationship{
				Type:           viceVersaType,
				FromInstanceID: toInstanceID,
				FromProductID:  toProductID,
				ToInstanceID:   fromInstanceID,
				ToProductID:    fromProductID,
			}

			err = relationshipRepository.Create(tx, viceVersaRel)
			if err != nil {
				return fmt.Errorf("error creating vice-versa relationship: %w", err)
			}
		}

		return nil
	})
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

// ListRelationshipsByFromInstanceID retrieves all relationships originating from a specific instance.
// Returns a slice of relationships or an error if database access fails.
func ListRelationshipsByFromInstanceID(fromInstanceID uuid.UUID) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByFromInstanceID(nil, fromInstanceID)
}

// ListRelationshipsByFromProductID retrieves all relationships originating from a specific product.
// Returns a slice of relationships or an error if database access fails.
func ListRelationshipsByFromProductID(fromProductID uuid.UUID) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByFromProductID(nil, fromProductID)
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
