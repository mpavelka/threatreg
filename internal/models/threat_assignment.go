package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreatAssignment struct {
	ID                 int                 `gorm:"primaryKey;autoIncrement;not null;unique"`
	ThreatID           uuid.UUID           `gorm:"type:uuid"`
	ProductID          uuid.UUID           `gorm:"type:uuid"`
	InstanceID         uuid.UUID           `gorm:"type:uuid"`
	Threat             Threat              `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Product            Product             `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Instance           Instance            `gorm:"foreignKey:InstanceID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	ControlAssignments []ControlAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

type ThreatAssignmentRepository struct {
	db *gorm.DB
}

func NewThreatAssignmentRepository(db *gorm.DB) *ThreatAssignmentRepository {
	return &ThreatAssignmentRepository{db: db}
}

func (r *ThreatAssignmentRepository) AssignThreatToProduct(tx *gorm.DB, threatID, productID uuid.UUID) (*ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}

	// Check if assignment already exists - explicitly check for uuid.Nil or NULL
	var existing ThreatAssignment
	err := tx.Where("threat_id = ? AND product_id = ? AND (instance_id IS NULL OR instance_id = ?)", threatID, productID, uuid.Nil).First(&existing).Error
	if err == nil {
		return &existing, nil // Assignment already exists
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err // Other error occurred
	}

	// Create new assignment - explicitly set InstanceID to NULL-equivalent
	assignment := &ThreatAssignment{
		ThreatID:   threatID,
		ProductID:  productID,
		InstanceID: uuid.Nil, // Explicitly set to nil UUID
	}

	err = tx.Create(assignment).Error
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (r *ThreatAssignmentRepository) AssignThreatToInstance(tx *gorm.DB, threatID, instanceID uuid.UUID) (*ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}

	// Check if assignment already exists
	var existing ThreatAssignment
	err := tx.Where("threat_id = ? AND instance_id = ? AND (product_id IS NULL OR product_id = ?)", threatID, instanceID, uuid.Nil).First(&existing).Error
	if err == nil {
		return &existing, nil // Assignment already exists
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err // Other error occurred
	}

	// Create new assignment - explicitly set ProductID to NULL-equivalent
	assignment := &ThreatAssignment{
		ThreatID:   threatID,
		InstanceID: instanceID,
		ProductID:  uuid.Nil, // Explicitly set to nil UUID
	}

	err = tx.Create(assignment).Error
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (r *ThreatAssignmentRepository) GetByID(tx *gorm.DB, id int) (*ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}
	var assignment ThreatAssignment
	err := tx.Preload("Threat").Preload("Product").Preload("Instance").First(&assignment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *ThreatAssignmentRepository) Delete(tx *gorm.DB, id int) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&ThreatAssignment{}, "id = ?", id).Error
}

func (r *ThreatAssignmentRepository) ListByProductID(tx *gorm.DB, productID uuid.UUID) ([]ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}

	var assignments []ThreatAssignment
	err := tx.Preload("Threat").Preload("Product").Preload("Instance").
		Where("product_id = ?", productID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *ThreatAssignmentRepository) ListByInstanceID(tx *gorm.DB, instanceID uuid.UUID) ([]ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}

	var assignments []ThreatAssignment
	err := tx.Preload("Threat").Preload("Product").Preload("Instance").
		Where("instance_id = ?", instanceID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}
