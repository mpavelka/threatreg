package models

import (
	"github.com/google/uuid"
)

type ControlAssignment struct {
	ID                 int              `gorm:"primaryKey;autoIncrement;not null;unique"`
	ThreatAssignmentID int              `gorm:"not null"`
	ControlID          uuid.UUID        `gorm:"type:uuid"`
	ThreatAssignment   ThreatAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Control            Control          `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}
