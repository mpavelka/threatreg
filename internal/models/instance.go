package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Instance struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique"`
	Name              string             `gorm:"type:varchar(255);index"`
	InstanceOf        uuid.UUID          `gorm:"type:uuid"`
	Product           Product            `gorm:"foreignKey:InstanceOf;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:InstanceID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Domains           []Domain           `gorm:"many2many:domain_instances;"`
}

func (i *Instance) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}

type InstanceRepository struct {
	db *gorm.DB
}

func NewInstanceRepository(db *gorm.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) Create(tx *gorm.DB, instance *Instance) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(instance).Error
}

func (r *InstanceRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Instance, error) {
	if tx == nil {
		tx = r.db
	}
	var instance Instance
	err := tx.Preload("Product").First(&instance, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *InstanceRepository) Update(tx *gorm.DB, instance *Instance) error {
	if tx == nil {
		tx = r.db
	}
	// Use Updates to explicitly update both name and foreign key field
	// GORM may not update foreign key fields with Save due to the relationship
	return tx.Model(instance).Select("name", "instance_of").Updates(instance).Error
}

func (r *InstanceRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Instance{}, "id = ?", id).Error
}

func (r *InstanceRepository) List(tx *gorm.DB) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *InstanceRepository) ListByProductID(tx *gorm.DB, productID uuid.UUID) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").Where("instance_of = ?", productID).Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *InstanceRepository) Filter(tx *gorm.DB, instanceName, productName string) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	query := tx.Preload("Product")

	// Detect database dialect for case-insensitive comparison
	dialect := tx.Dialector.Name()

	if instanceName != "" {
		if dialect == "postgres" {
			query = query.Where("instances.name ILIKE ?", "%"+instanceName+"%")
		} else {
			query = query.Where("LOWER(instances.name) LIKE LOWER(?)", "%"+instanceName+"%")
		}
	}

	if productName != "" {
		if dialect == "postgres" {
			query = query.Where("Product.name ILIKE ?", "%"+productName+"%").Joins("Product")
		} else {
			query = query.Where("LOWER(Product.name) LIKE LOWER(?)", "%"+productName+"%").Joins("Product")
		}
	}

	err := query.Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *InstanceRepository) ListByDomainId(tx *gorm.DB, domainID uuid.UUID) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").
		Joins("JOIN domain_instances ON instances.id = domain_instances.instance_id").
		Where("domain_instances.domain_id = ?", domainID).
		Find(&instances).Error

	if err != nil {
		return nil, err
	}
	return instances, nil
}
