package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ControlAssignment struct {
	ID                 uuid.UUID        `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	ThreatAssignmentID uuid.UUID        `gorm:"type:uuid;not null" json:"threatAssignmentId"`
	ControlID          uuid.UUID        `gorm:"type:uuid" json:"controlId"`
	ThreatAssignment   ThreatAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threatAssignment"`
	Control            Control          `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"control"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (ca *ControlAssignment) BeforeCreate(tx *gorm.DB) error {
	if ca.ID == uuid.Nil {
		ca.ID = uuid.New()
	}
	return nil
}
