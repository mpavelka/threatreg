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
