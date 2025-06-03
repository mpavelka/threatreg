package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Control struct {
	ID                 uuid.UUID           `gorm:"type:uuid;primaryKey;not null;unique"`
	Title              string              `gorm:"type:varchar(255)"`
	Description        string              `gorm:"type:text"`
	ControlAssignments []ControlAssignment `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ThreatControls     []ThreatControl     `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

func (c *Control) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
