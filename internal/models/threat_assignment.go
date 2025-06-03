package models

import (
	"github.com/google/uuid"
)

type ThreatAssignment struct {
	ID                 int                 `gorm:"primaryKey;autoIncrement;not null;unique"`
	ThreatID           uuid.UUID           `gorm:"type:uuid"`
	ProductID          uuid.UUID           `gorm:"type:uuid"`
	ApplicationID      uuid.UUID           `gorm:"type:uuid"`
	Threat             Threat              `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Product            Product             `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Application        Application         `gorm:"foreignKey:ApplicationID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ControlAssignments []ControlAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}
