package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;not null;unique"`
	Name        string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string     `gorm:"type:text"`
	Color       string     `gorm:"type:varchar(7)"` // Hex color code like #FF0000
	Products    []Product  `gorm:"many2many:product_tags;"`
	Instances   []Instance `gorm:"many2many:instance_tags;"`
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
	err := tx.Preload("Products").Preload("Instances").First(&tag, "id = ?", id).Error
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
	err := tx.Preload("Products").Preload("Instances").First(&tag, "name = ?", name).Error
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
	err := tx.Preload("Products").Preload("Instances").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) ListByProductID(tx *gorm.DB, productID uuid.UUID) ([]Tag, error) {
	if tx == nil {
		tx = r.db
	}

	var tags []Tag
	err := tx.Preload("Products").Preload("Instances").
		Joins("JOIN product_tags ON tags.id = product_tags.tag_id").
		Where("product_tags.product_id = ?", productID).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) ListByInstanceID(tx *gorm.DB, instanceID uuid.UUID) ([]Tag, error) {
	if tx == nil {
		tx = r.db
	}

	var tags []Tag
	err := tx.Preload("Products").Preload("Instances").
		Joins("JOIN instance_tags ON tags.id = instance_tags.tag_id").
		Where("instance_tags.instance_id = ?", instanceID).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) AssignToProduct(tx *gorm.DB, tagID, productID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	// Check if association already exists
	var count int64
	err := tx.Table("product_tags").
		Where("tag_id = ? AND product_id = ?", tagID, productID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Association already exists
	}

	// Create association
	return tx.Exec("INSERT INTO product_tags (tag_id, product_id) VALUES (?, ?)", tagID, productID).Error
}

func (r *TagRepository) UnassignFromProduct(tx *gorm.DB, tagID, productID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Exec("DELETE FROM product_tags WHERE tag_id = ? AND product_id = ?", tagID, productID).Error
}

func (r *TagRepository) AssignToInstance(tx *gorm.DB, tagID, instanceID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	// Check if association already exists
	var count int64
	err := tx.Table("instance_tags").
		Where("tag_id = ? AND instance_id = ?", tagID, instanceID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Association already exists
	}

	// Create association
	return tx.Exec("INSERT INTO instance_tags (tag_id, instance_id) VALUES (?, ?)", tagID, instanceID).Error
}

func (r *TagRepository) UnassignFromInstance(tx *gorm.DB, tagID, instanceID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Exec("DELETE FROM instance_tags WHERE tag_id = ? AND instance_id = ?", tagID, instanceID).Error
}

func (r *TagRepository) ListProductsByTagID(tx *gorm.DB, tagID uuid.UUID) ([]Product, error) {
	if tx == nil {
		tx = r.db
	}

	var products []Product
	err := tx.Joins("JOIN product_tags ON products.id = product_tags.product_id").
		Where("product_tags.tag_id = ?", tagID).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *TagRepository) ListInstancesByTagID(tx *gorm.DB, tagID uuid.UUID) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").
		Joins("JOIN instance_tags ON instances.id = instance_tags.instance_id").
		Where("instance_tags.tag_id = ?", tagID).
		Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}
