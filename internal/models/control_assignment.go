package models

import (
	"github.com/google/uuid"
)

type ControlAssignment struct {
	ID                 int              `gorm:"primaryKey;autoIncrement;not null;unique" json:"id"`
	ThreatAssignmentID int              `gorm:"not null" json:"threatAssignmentID"`
	ControlID          uuid.UUID        `gorm:"type:uuid" json:"controlID"`
	ThreatAssignment   ThreatAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threatAssignment"`
	Control            Control          `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"control"`
}
