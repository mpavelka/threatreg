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

func GetTag(id uuid.UUID) (*models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.GetByID(nil, id)
}

func GetTagByName(name string) (*models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.GetByName(nil, name)
}

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

func DeleteTag(id uuid.UUID) error {
	tagRepository, err := getTagRepository()
	if err != nil {
		return err
	}

	return tagRepository.Delete(nil, id)
}

func ListTags() ([]models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.List(nil)
}

func ListTagsByProductID(productID uuid.UUID) ([]models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.ListByProductID(nil, productID)
}

func ListTagsByInstanceID(instanceID uuid.UUID) ([]models.Tag, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.ListByInstanceID(nil, instanceID)
}

func AssignTagToProduct(tagID, productID uuid.UUID) error {
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

		// Verify product exists
		productRepository, err := getProductRepository()
		if err != nil {
			return err
		}
		_, err = productRepository.GetByID(tx, productID)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		return tagRepository.AssignToProduct(tx, tagID, productID)
	})
}

func UnassignTagFromProduct(tagID, productID uuid.UUID) error {
	tagRepository, err := getTagRepository()
	if err != nil {
		return err
	}

	return tagRepository.UnassignFromProduct(nil, tagID, productID)
}

func AssignTagToInstance(tagID, instanceID uuid.UUID) error {
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

		// Verify instance exists
		instanceRepository, err := getInstanceRepository()
		if err != nil {
			return err
		}
		_, err = instanceRepository.GetByID(tx, instanceID)
		if err != nil {
			return fmt.Errorf("instance not found: %w", err)
		}

		return tagRepository.AssignToInstance(tx, tagID, instanceID)
	})
}

func UnassignTagFromInstance(tagID, instanceID uuid.UUID) error {
	tagRepository, err := getTagRepository()
	if err != nil {
		return err
	}

	return tagRepository.UnassignFromInstance(nil, tagID, instanceID)
}

func ListProductsByTagID(tagID uuid.UUID) ([]models.Product, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.ListProductsByTagID(nil, tagID)
}

func ListInstancesByTagID(tagID uuid.UUID) ([]models.Instance, error) {
	tagRepository, err := getTagRepository()
	if err != nil {
		return nil, err
	}

	return tagRepository.ListInstancesByTagID(nil, tagID)
}

func AssignTagToProductByName(tagName string, productID uuid.UUID) error {
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

		// Verify product exists
		productRepository, err := getProductRepository()
		if err != nil {
			return err
		}
		_, err = productRepository.GetByID(tx, productID)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		return tagRepository.AssignToProduct(tx, tag.ID, productID)
	})
}

func AssignTagToInstanceByName(tagName string, instanceID uuid.UUID) error {
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

		// Verify instance exists
		instanceRepository, err := getInstanceRepository()
		if err != nil {
			return err
		}
		_, err = instanceRepository.GetByID(tx, instanceID)
		if err != nil {
			return fmt.Errorf("instance not found: %w", err)
		}

		return tagRepository.AssignToInstance(tx, tag.ID, instanceID)
	})
}
