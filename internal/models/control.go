package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Control struct {
	ID                 uuid.UUID           `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	Title              string              `gorm:"type:varchar(255)" json:"title"`
	Description        string              `gorm:"type:text" json:"description"`
	ControlAssignments []ControlAssignment `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"controlAssignments"`
	ThreatControls     []ThreatControl     `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threatControls"`
}

func (c *Control) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

type ControlRepository struct {
	db *gorm.DB
}

func NewControlRepository(db *gorm.DB) *ControlRepository {
	return &ControlRepository{db: db}
}

func (r *ControlRepository) Create(tx *gorm.DB, control *Control) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(control).Error
}

func (r *ControlRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Control, error) {
	if tx == nil {
		tx = r.db
	}
	var control Control
	err := tx.First(&control, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &control, nil
}

func (r *ControlRepository) Update(tx *gorm.DB, control *Control) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(control).Error
}

func (r *ControlRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Control{}, "id = ?", id).Error
}

func (r *ControlRepository) List(tx *gorm.DB) ([]Control, error) {
	if tx == nil {
		tx = r.db
	}

	var controls []Control
	err := tx.Find(&controls).Error
	if err != nil {
		return nil, err
	}
	return controls, nil
}
