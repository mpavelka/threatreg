package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique"`
	Name              string             `gorm:"type:varchar(255)"`
	InstanceOf        uuid.UUID          `gorm:"type:uuid"`
	Product           Product            `gorm:"foreignKey:InstanceOf;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:ApplicationID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

func (a *Application) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(tx *gorm.DB, application *Application) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(application).Error
}

func (r *ApplicationRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Application, error) {
	if tx == nil {
		tx = r.db
	}
	var application Application
	err := tx.Preload("Product").First(&application, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *ApplicationRepository) Update(tx *gorm.DB, application *Application) error {
	if tx == nil {
		tx = r.db
	}
	// Use Updates to explicitly update both name and foreign key field
	// GORM may not update foreign key fields with Save due to the relationship
	return tx.Model(application).Select("name", "instance_of").Updates(application).Error
}

func (r *ApplicationRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Application{}, "id = ?", id).Error
}

func (r *ApplicationRepository) List(tx *gorm.DB) ([]Application, error) {
	if tx == nil {
		tx = r.db
	}

	var applications []Application
	err := tx.Preload("Product").Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) ListByProductID(tx *gorm.DB, productID uuid.UUID) ([]Application, error) {
	if tx == nil {
		tx = r.db
	}

	var applications []Application
	err := tx.Preload("Product").Where("instance_of = ?", productID).Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}
