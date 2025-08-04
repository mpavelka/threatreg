package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Relationship struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	Type            string    `gorm:"type:varchar(255);index" json:"type"`
	FromComponentID uuid.UUID `gorm:"type:uuid;index" json:"fromComponentId"`
	ToComponentID   uuid.UUID `gorm:"type:uuid;index" json:"toComponentId"`
	FromComponent   Component `gorm:"foreignKey:FromComponentID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"fromComponent"`
	ToComponent     Component `gorm:"foreignKey:ToComponentID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"toComponent"`
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
	err := tx.Preload("FromComponent").Preload("ToComponent").First(&relationship, "id = ?", id).Error
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
	err := tx.Preload("FromComponent").Preload("ToComponent").Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

func (r *RelationshipRepository) ListByFromComponentID(tx *gorm.DB, fromComponentID uuid.UUID) ([]Relationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationships []Relationship
	err := tx.Preload("FromComponent").Preload("ToComponent").Where("from_component_id = ?", fromComponentID).Find(&relationships).Error
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
	err := tx.Preload("FromComponent").Preload("ToComponent").Where("type = ?", relType).Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}
