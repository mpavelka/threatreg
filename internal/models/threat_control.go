package models

import (
	"github.com/google/uuid"
)

type ThreatControl struct {
	ID        int       `gorm:"primaryKey;autoIncrement;not null;unique"`
	ThreatID  uuid.UUID `gorm:"type:uuid"`
	ControlID uuid.UUID `gorm:"type:uuid"`
	Threat    Threat    `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Control   Control   `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}
