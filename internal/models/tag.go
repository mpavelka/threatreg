package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	Name        string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string      `gorm:"type:text" json:"description"`
	Color       string      `gorm:"type:varchar(7)" json:"color"` // Hex color code like #FF0000
	Components  []Component `gorm:"many2many:component_tags;" json:"components"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tx *gorm.DB, tag *Tag) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(tag).Error
}

func (r *TagRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Tag, error) {
	if tx == nil {
		tx = r.db
	}
	var tag Tag
	err := tx.Preload("Components").First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) GetByName(tx *gorm.DB, name string) (*Tag, error) {
	if tx == nil {
		tx = r.db
	}
	var tag Tag
	err := tx.Preload("Components").First(&tag, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) Update(tx *gorm.DB, tag *Tag) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(tag).Error
}

func (r *TagRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Tag{}, "id = ?", id).Error
}

func (r *TagRepository) List(tx *gorm.DB) ([]Tag, error) {
	if tx == nil {
		tx = r.db
	}

	var tags []Tag
	err := tx.Preload("Components").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) ListByComponentID(tx *gorm.DB, componentID uuid.UUID) ([]Tag, error) {
	if tx == nil {
		tx = r.db
	}

	var tags []Tag
	err := tx.Preload("Components").
		Joins("JOIN component_tags ON tags.id = component_tags.tag_id").
		Where("component_tags.component_id = ?", componentID).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) AssignToComponent(tx *gorm.DB, tagID, componentID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	// Check if association already exists
	var count int64
	err := tx.Table("component_tags").
		Where("tag_id = ? AND component_id = ?", tagID, componentID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Association already exists
	}

	// Create association
	return tx.Exec("INSERT INTO component_tags (tag_id, component_id) VALUES (?, ?)", tagID, componentID).Error
}

func (r *TagRepository) UnassignFromComponent(tx *gorm.DB, tagID, componentID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Exec("DELETE FROM component_tags WHERE tag_id = ? AND component_id = ?", tagID, componentID).Error
}

func (r *TagRepository) ListComponentsByTagID(tx *gorm.DB, tagID uuid.UUID) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}

	var components []Component
	err := tx.Joins("JOIN component_tags ON components.id = component_tags.component_id").
		Where("component_tags.tag_id = ?", tagID).
		Find(&components).Error
	if err != nil {
		return nil, err
	}
	return components, nil
}
