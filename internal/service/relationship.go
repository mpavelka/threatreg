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

// AddRelationshipConsumesApiOf creates bidirectional API consumption relationships between components.
// Creates 'CONSUMES_API_OF' and 'API_CONSUMED_BY' relationships.
func AddRelationshipConsumesApiOf(fromComponentID, toComponentID uuid.UUID) error {
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

// AddRelationship creates relationships between components with optional reverse relationship.
// If viceVersaType is provided, creates a bidirectional relationship.
func AddRelationship(fromComponentID, toComponentID uuid.UUID, relType, viceVersaType string) error {
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

// ListRelationshipsByType retrieves all relationships of a specific type.
// Returns a slice of relationships matching the type or an error if database access fails.
func ListRelationshipsByType(relType string) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByType(nil, relType)
}
