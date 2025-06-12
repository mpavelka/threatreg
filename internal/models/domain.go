package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Domain struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;not null;unique"`
	Name        string     `gorm:"type:varchar(255);index"`
	Description string     `gorm:"type:text"`
	Instances   []Instance `gorm:"many2many:domain_instances;"`
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
	err := tx.Preload("Instances").Preload("Instances.Product").First(&domain, "id = ?", id).Error
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
	err := tx.Preload("Instances").Preload("Instances.Product").Find(&domains).Error
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func (r *DomainRepository) AddInstance(tx *gorm.DB, domainID, instanceID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	domain, err := r.GetByID(tx, domainID)
	if err != nil {
		return err
	}

	var instance Instance
	err = tx.First(&instance, "id = ?", instanceID).Error
	if err != nil {
		return err
	}

	return tx.Model(domain).Association("Instances").Append(&instance)
}

func (r *DomainRepository) RemoveInstance(tx *gorm.DB, domainID, instanceID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	domain, err := r.GetByID(tx, domainID)
	if err != nil {
		return err
	}

	var instance Instance
	err = tx.First(&instance, "id = ?", instanceID).Error
	if err != nil {
		return err
	}

	return tx.Model(domain).Association("Instances").Delete(&instance)
}

func (r *DomainRepository) GetInstancesByDomainID(tx *gorm.DB, domainID uuid.UUID) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").
		Joins("JOIN domain_instances ON instances.id = domain_instances.instance_id").
		Where("domain_instances.domain_id = ?", domainID).
		Find(&instances).Error

	return instances, err
}

func (r *DomainRepository) GetDomainsByInstanceID(tx *gorm.DB, instanceID uuid.UUID) ([]Domain, error) {
	if tx == nil {
		tx = r.db
	}

	var domains []Domain
	err := tx.
		Joins("JOIN domain_instances ON domains.id = domain_instances.domain_id").
		Where("domain_instances.instance_id = ?", instanceID).
		Find(&domains).Error

	return domains, err
}
