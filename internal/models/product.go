package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a product entity in the domain
type Product struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique"`
	Name              string             `gorm:"type:varchar(255);index"`
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

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(tx *gorm.DB, product *Product) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(product).Error
}

func (r *ProductRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Product, error) {
	if tx == nil {
		tx = r.db
	}
	var product Product
	err := tx.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(tx *gorm.DB, product *Product) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(product).Error
}

func (r *ProductRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Product{}, "id = ?", id).Error
}

func (r *ProductRepository) List(tx *gorm.DB) ([]Product, error) {
	if tx == nil {
		tx = r.db
	}

	var products []Product
	err := tx.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
