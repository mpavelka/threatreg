package models

import (
	"github.com/google/uuid"
)

type ThreatControl struct {
	ID        int       `gorm:"primaryKey;autoIncrement;not null;unique" json:"id"`
	ThreatID  uuid.UUID `gorm:"type:uuid" json:"threatID"`
	ControlID uuid.UUID `gorm:"type:uuid" json:"controlID"`
	Threat    Threat    `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threat"`
	Control   Control   `gorm:"foreignKey:ControlID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"control"`
}
