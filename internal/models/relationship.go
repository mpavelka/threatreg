package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Relationship struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;not null;unique"`
	Type           string     `gorm:"type:varchar(255);index"`
	FromInstanceID *uuid.UUID `gorm:"type:uuid;index"`
	FromProductID  *uuid.UUID `gorm:"type:uuid;index"`
	ToInstanceID   *uuid.UUID `gorm:"type:uuid;index"`
	ToProductID    *uuid.UUID `gorm:"type:uuid;index"`
	FromInstance   *Instance  `gorm:"foreignKey:FromInstanceID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	FromProduct    *Product   `gorm:"foreignKey:FromProductID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ToInstance     *Instance  `gorm:"foreignKey:ToInstanceID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ToProduct      *Product   `gorm:"foreignKey:ToProductID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (r *Relationship) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type RelationshipRepository struct {
	db *gorm.DB
}

func NewRelationshipRepository(db *gorm.DB) *RelationshipRepository {
	return &RelationshipRepository{db: db}
}

func (r *RelationshipRepository) Create(tx *gorm.DB, relationship *Relationship) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(relationship).Error
}

func (r *RelationshipRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Relationship, error) {
	if tx == nil {
		tx = r.db
	}
	var relationship Relationship
	err := tx.Preload("FromInstance").Preload("FromProduct").Preload("ToInstance").Preload("ToProduct").First(&relationship, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

func (r *RelationshipRepository) Update(tx *gorm.DB, relationship *Relationship) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(relationship).Error
}

func (r *RelationshipRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Relationship{}, "id = ?", id).Error
}

func (r *RelationshipRepository) List(tx *gorm.DB) ([]Relationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationships []Relationship
	err := tx.Preload("FromInstance").Preload("FromProduct").Preload("ToInstance").Preload("ToProduct").Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

func (r *RelationshipRepository) ListByFromInstanceID(tx *gorm.DB, fromInstanceID uuid.UUID) ([]Relationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationships []Relationship
	err := tx.Preload("FromInstance").Preload("FromProduct").Preload("ToInstance").Preload("ToProduct").Where("from_instance_id = ?", fromInstanceID).Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

func (r *RelationshipRepository) ListByFromProductID(tx *gorm.DB, fromProductID uuid.UUID) ([]Relationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationships []Relationship
	err := tx.Preload("FromInstance").Preload("FromProduct").Preload("ToInstance").Preload("ToProduct").Where("from_product_id = ?", fromProductID).Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

func (r *RelationshipRepository) ListByType(tx *gorm.DB, relType string) ([]Relationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationships []Relationship
	err := tx.Preload("FromInstance").Preload("FromProduct").Preload("ToInstance").Preload("ToProduct").Where("type = ?", relType).Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}
