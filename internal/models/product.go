package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a product entity in the domain
type Product struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique"`
	Name              string             `gorm:"type:varchar(255)"`
	Description       string             `gorm:"type:text"`
	Applications      []Application      `gorm:"foreignKey:InstanceOf;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
