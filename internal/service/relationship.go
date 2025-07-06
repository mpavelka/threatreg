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

func DeleteRelationshipById(id uuid.UUID) error {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return err
	}

	return relationshipRepository.Delete(nil, id)
}

func GetRelationship(id uuid.UUID) (*models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.GetByID(nil, id)
}

func ListRelationships() ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.List(nil)
}

func ListRelationshipsByFromInstanceID(fromInstanceID uuid.UUID) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByFromInstanceID(nil, fromInstanceID)
}

func ListRelationshipsByFromProductID(fromProductID uuid.UUID) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByFromProductID(nil, fromProductID)
}

func ListRelationshipsByType(relType string) ([]models.Relationship, error) {
	relationshipRepository, err := getRelationshipRepository()
	if err != nil {
		return nil, err
	}

	return relationshipRepository.ListByType(nil, relType)
}
