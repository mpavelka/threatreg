package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Threat struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique"`
	Title             string             `gorm:"type:varchar(255)"`
	Description       string             `gorm:"type:text"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ThreatControls    []ThreatControl    `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

// BeforeCreate is a GORM hook that is triggered before a new record is created.
func (t *Threat) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

type ThreatRepository struct {
	db *gorm.DB
}

func NewThreatRepository(db *gorm.DB) *ThreatRepository {
	return &ThreatRepository{db: db}
}

func (r *ThreatRepository) Create(tx *gorm.DB, threat *Threat) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(threat).Error
}

func (r *ThreatRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Threat, error) {
	if tx == nil {
		tx = r.db
	}
	var threat Threat
	err := tx.First(&threat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &threat, nil
}

func (r *ThreatRepository) Update(tx *gorm.DB, threat *Threat) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(threat).Error
}

func (r *ThreatRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Threat{}, "id = ?", id).Error
}

func (r *ThreatRepository) List(tx *gorm.DB) ([]Threat, error) {
	if tx == nil {
		tx = r.db
	}

	var threats []Threat
	err := tx.Find(&threats).Error
	if err != nil {
		return nil, err
	}
	return threats, nil
}
