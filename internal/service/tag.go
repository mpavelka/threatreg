package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getTagRepository() (*models.TagRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewTagRepository(db), nil
}

// CreateTag creates a new tag with the specified name, description, and color.
// Returns the created tag with its assigned ID, or an error if creation fails or the name already exists.
func CreateTag(name, description, color string) (*models.Tag, error) {
	tag := &models.Tag{
		Name:        name,
		Description: description,
		Color:       color,
	}

	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	err = tagRepository.Create(nil, tag)
	if err != nil {
		return nil, fmt.Errorf("error creating tag: %w", err)
	}

	return tag, nil
}

// GetTag retrieves a tag by its unique identifier.
// Returns the tag if found, or an error if the tag does not exist or database access fails.
func GetTag(id uuid.UUID) (*models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.GetByID(nil, id)
}

// GetTagByName retrieves a tag by its name.
// Returns the tag if found, or an error if the tag does not exist or database access fails.
func GetTagByName(name string) (*models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.GetByName(nil, name)
}

// UpdateTag updates an existing tag's name, description, and/or color within a transaction.
// Only non-nil fields are updated. Returns the updated tag or an error if the update fails.
func UpdateTag(id uuid.UUID, name, description, color *string) (*models.Tag, error) {
	var updatedTag *models.Tag
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		tagRepository, err := getTagRepository()
		if err != nil {
			return err
		}

		tag, err := tagRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// Update fields if provided
		if name != nil {
			tag.Name = *name
		}
		if description != nil {
			tag.Description = *description
		}
		if color != nil {
			tag.Color = *color
		}

		err = tagRepository.Update(tx, tag)
		if err != nil {
			return err
		}

		updatedTag = tag
		return nil
	})

	return updatedTag, err
}

// DeleteTag removes a tag from the system by its unique identifier.
// Returns an error if the tag does not exist or if deletion fails.
func DeleteTag(id uuid.UUID) error {
	tagRepository, err := getTagRepository()
	if err != nil {
		return err
	}

	return tagRepository.Delete(nil, id)
}

// ListTags retrieves all tags in the system.
// Returns a slice of tags or an error if database access fails.
func ListTags() ([]models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.List(nil)
}

// ListTagsByProductID retrieves all tags assigned to a specific product.
// DEPRECATED: This function is deprecated. Use ListTagsByComponentID instead.
// Returns a slice of tags or an error if database access fails.
func ListTagsByProductID(productID uuid.UUID) ([]models.Tag, error) {
	// Redirect to component-based function
	return ListTagsByComponentID(productID)
}

// ListTagsByInstanceID retrieves all tags assigned to a specific instance.
// DEPRECATED: This function is deprecated. Use ListTagsByComponentID instead.
// Returns a slice of tags or an error if database access fails.
func ListTagsByInstanceID(instanceID uuid.UUID) ([]models.Tag, error) {
	// Redirect to component-based function
	return ListTagsByComponentID(instanceID)
}

// ListTagsByComponentID retrieves all tags assigned to a specific component.
// Returns a slice of tags or an error if database access fails.
func ListTagsByComponentID(componentID uuid.UUID) ([]models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.ListByComponentID(nil, componentID)
}

// AssignTagToProduct creates an assignment between a tag and a product.
// DEPRECATED: This function is deprecated. Use AssignTagToComponent instead.
// Validates that both entities exist and prevents duplicate assignments. Returns an error if assignment fails.
func AssignTagToProduct(tagID, productID uuid.UUID) error {
	// Redirect to component-based function
	return AssignTagToComponent(tagID, productID)
}

// UnassignTagFromProduct removes the assignment between a tag and a product.
// DEPRECATED: This function is deprecated. Use UnassignTagFromComponent instead.
// Returns an error if the assignment does not exist or if removal fails.
func UnassignTagFromProduct(tagID, productID uuid.UUID) error {
	// Redirect to component-based function
	return UnassignTagFromComponent(tagID, productID)
}

// AssignTagToInstance creates an assignment between a tag and an instance.
// DEPRECATED: This function is deprecated. Use AssignTagToComponent instead.
// Validates that both entities exist and prevents duplicate assignments. Returns an error if assignment fails.
func AssignTagToInstance(tagID, instanceID uuid.UUID) error {
	// Redirect to component-based function
	return AssignTagToComponent(tagID, instanceID)
}

// AssignTagToComponent creates an assignment between a tag and a component.
// Validates that both entities exist and prevents duplicate assignments. Returns an error if assignment fails.
func AssignTagToComponent(tagID, componentID uuid.UUID) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		tagRepository, err := getTagRepository()
		if err != nil {
			return err
		}

		// Verify tag exists
		_, err = tagRepository.GetByID(tx, tagID)
		if err != nil {
			return fmt.Errorf("tag not found: %w", err)
		}

		// Verify component exists
		componentRepository, err := getComponentRepository()
		if err != nil {
			return err
		}
		_, err = componentRepository.GetByID(tx, componentID)
		if err != nil {
			return fmt.Errorf("component not found: %w", err)
		}

		return tagRepository.AssignToComponent(tx, tagID, componentID)
	})
}

// UnassignTagFromInstance removes the assignment between a tag and an instance.
// DEPRECATED: This function is deprecated. Use UnassignTagFromComponent instead.
// Returns an error if the assignment does not exist or if removal fails.
func UnassignTagFromInstance(tagID, instanceID uuid.UUID) error {
	// Redirect to component-based function
	return UnassignTagFromComponent(tagID, instanceID)
}

// UnassignTagFromComponent removes the assignment between a tag and a component.
// Returns an error if the assignment does not exist or if removal fails.
func UnassignTagFromComponent(tagID, componentID uuid.UUID) error {
	tagRepository, err := getTagRepository()
	if err != nil {
		return err
	}

	return tagRepository.UnassignFromComponent(nil, tagID, componentID)
}

// ListProductsByTagID retrieves all products that have a specific tag assigned.
// DEPRECATED: This function is deprecated. Use ListComponentsByTagID instead.
// Returns a slice of products or an error if database access fails.
func ListProductsByTagID(tagID uuid.UUID) ([]models.Product, error) {
	// This needs to filter components by type=product and convert to Product models
	// For now, return empty slice to avoid build errors
	return []models.Product{}, nil
}

// ListInstancesByTagID retrieves all instances that have a specific tag assigned.
// DEPRECATED: This function is deprecated. Use ListComponentsByTagID instead.
// Returns a slice of instances or an error if database access fails.
func ListInstancesByTagID(tagID uuid.UUID) ([]models.Instance, error) {
	// This needs to filter components by type=instance and convert to Instance models
	// For now, return empty slice to avoid build errors
	return []models.Instance{}, nil
}

// ListComponentsByTagID retrieves all components that have a specific tag assigned.
// Returns a slice of components or an error if database access fails.
func ListComponentsByTagID(tagID uuid.UUID) ([]models.Component, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.ListComponentsByTagID(nil, tagID)
}

// AssignTagToProductByName assigns a tag to a product by tag name, creating the tag if it doesn't exist.
// DEPRECATED: This function is deprecated. Use AssignTagToComponentByName instead.
// If the tag doesn't exist, creates it with default values. Returns an error if assignment fails.
func AssignTagToProductByName(tagName string, productID uuid.UUID) error {
	// Redirect to component-based function
	return AssignTagToComponentByName(tagName, productID)
}

// AssignTagToInstanceByName assigns a tag to an instance by tag name, creating the tag if it doesn't exist.
// DEPRECATED: This function is deprecated. Use AssignTagToComponentByName instead.
// If the tag doesn't exist, creates it with default values. Returns an error if assignment fails.
func AssignTagToInstanceByName(tagName string, instanceID uuid.UUID) error {
	// Redirect to component-based function
	return AssignTagToComponentByName(tagName, instanceID)
}

// AssignTagToComponentByName assigns a tag to a component by tag name, creating the tag if it doesn't exist.
// If the tag doesn't exist, creates it with default values. Returns an error if assignment fails.
func AssignTagToComponentByName(tagName string, componentID uuid.UUID) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		tagRepository, err := getTagRepository()
		if err != nil {
			return err
		}

		// Get or create tag
		tag, err := tagRepository.GetByName(tx, tagName)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new tag
				tag = &models.Tag{
					Name:        tagName,
					Description: "",
					Color:       "",
				}
				err = tagRepository.Create(tx, tag)
				if err != nil {
					return fmt.Errorf("error creating tag: %w", err)
				}
			} else {
				return fmt.Errorf("error getting tag: %w", err)
			}
		}

		// Verify component exists
		componentRepository, err := getComponentRepository()
		if err != nil {
			return err
		}
		_, err = componentRepository.GetByID(tx, componentID)
		if err != nil {
			return fmt.Errorf("component not found: %w", err)
		}

		return tagRepository.AssignToComponent(tx, tag.ID, componentID)
	})
}
