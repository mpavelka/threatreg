package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a product entity in the domain
type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	gorm.Model
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}