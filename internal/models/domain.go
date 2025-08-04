package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Domain struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	Name        string      `gorm:"type:varchar(255);index" json:"name"`
	Description string      `gorm:"type:text" json:"description"`
	Components  []Component `gorm:"many2many:domain_components;" json:"components"`
}

func (d *Domain) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

func (r *DomainRepository) Create(tx *gorm.DB, domain *Domain) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(domain).Error
}

func (r *DomainRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Domain, error) {
	if tx == nil {
		tx = r.db
	}
	var domain Domain
	err := tx.First(&domain, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (r *DomainRepository) Update(tx *gorm.DB, domain *Domain) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(domain).Error
}

func (r *DomainRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Domain{}, "id = ?", id).Error
}

func (r *DomainRepository) List(tx *gorm.DB) ([]Domain, error) {
	if tx == nil {
		tx = r.db
	}

	var domains []Domain
	err := tx.Find(&domains).Error
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func (r *DomainRepository) AddComponent(tx *gorm.DB, domainID, componentID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	domain, err := r.GetByID(tx, domainID)
	if err != nil {
		return err
	}

	var component Component
	err = tx.First(&component, "id = ?", componentID).Error
	if err != nil {
		return err
	}

	return tx.Model(domain).Association("Components").Append(&component)
}

func (r *DomainRepository) RemoveComponent(tx *gorm.DB, domainID, componentID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	domain, err := r.GetByID(tx, domainID)
	if err != nil {
		return err
	}

	var component Component
	err = tx.First(&component, "id = ?", componentID).Error
	if err != nil {
		return err
	}

	return tx.Model(domain).Association("Components").Delete(&component)
}

func (r *DomainRepository) GetComponentsByDomainID(tx *gorm.DB, domainID uuid.UUID) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}

	var components []Component
	err := tx.Joins("JOIN domain_components ON components.id = domain_components.component_id").
		Where("domain_components.domain_id = ?", domainID).
		Find(&components).Error

	return components, err
}

func (r *DomainRepository) GetDomainsByComponentID(tx *gorm.DB, componentID uuid.UUID) ([]Domain, error) {
	if tx == nil {
		tx = r.db
	}

	var domains []Domain
	err := tx.
		Joins("JOIN domain_components ON domains.id = domain_components.domain_id").
		Where("domain_components.component_id = ?", componentID).
		Find(&domains).Error

	return domains, err
}
