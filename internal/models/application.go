package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique"`
	InstanceOf        uuid.UUID          `gorm:"type:uuid"`
	Product           Product            `gorm:"foreignKey:InstanceOf;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:ApplicationID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

func (a *Application) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
