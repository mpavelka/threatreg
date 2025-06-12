package models

import (
	"github.com/google/uuid"
)

type ThreatAssignment struct {
	ID                 int                 `gorm:"primaryKey;autoIncrement;not null;unique"`
	ThreatID           uuid.UUID           `gorm:"type:uuid"`
	ProductID          uuid.UUID           `gorm:"type:uuid"`
	InstanceID         uuid.UUID           `gorm:"type:uuid"`
	Threat             Threat              `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Product            Product             `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Instance           Instance            `gorm:"foreignKey:InstanceID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ControlAssignments []ControlAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}
