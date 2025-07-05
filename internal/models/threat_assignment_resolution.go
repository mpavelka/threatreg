package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreatAssignmentResolutionStatus string

const (
	ThreatAssignmentResolutionStatusResolved ThreatAssignmentResolutionStatus = "resolved"
	ThreatAssignmentResolutionStatusAwaiting ThreatAssignmentResolutionStatus = "awaiting"
	ThreatAssignmentResolutionStatusAccepted ThreatAssignmentResolutionStatus = "accepted"
)

type ThreatAssignmentResolution struct {
	ID                 uuid.UUID                        `gorm:"type:uuid;primaryKey"`
	ThreatAssignmentID int                              `gorm:"not null;uniqueIndex:idx_threat_assignment_resolution"`
	InstanceID         uuid.UUID                        `gorm:"type:uuid;uniqueIndex:idx_threat_assignment_resolution"`
	ProductID          uuid.UUID                        `gorm:"type:uuid;uniqueIndex:idx_threat_assignment_resolution"`
	Status             ThreatAssignmentResolutionStatus `gorm:"type:varchar(20);not null"`
	Description        string                           `gorm:"type:text"`
	ThreatAssignment   ThreatAssignment                 `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Instance           Instance                         `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Product            Product                          `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// BeforeCreate ensures exactly one of ProductID or InstanceID is set and validates status
func (tar *ThreatAssignmentResolution) BeforeCreate(tx *gorm.DB) error {
	if tar.ID == uuid.Nil {
		tar.ID = uuid.New()
	}
	return tar.validateResolution()
}

// BeforeUpdate ensures exactly one of ProductID or InstanceID is set and validates status
func (tar *ThreatAssignmentResolution) BeforeUpdate(tx *gorm.DB) error {
	return tar.validateResolution()
}

// validateResolution checks that exactly one of ProductID or InstanceID is not null/nil and validates status
func (tar *ThreatAssignmentResolution) validateResolution() error {
	productIsSet := tar.ProductID != uuid.Nil
	instanceIsSet := tar.InstanceID != uuid.Nil

	if productIsSet && instanceIsSet {
		return errors.New("threat assignment resolution cannot have both ProductID and InstanceID set")
	}

	if !productIsSet && !instanceIsSet {
		return errors.New("threat assignment resolution must have either ProductID or InstanceID set")
	}

	// Validate status
	if !tar.isValidStatus() {
		return errors.New("invalid status: must be 'resolved', 'awaiting', 'accepted', or 'delegated'")
	}

	return nil
}

// isValidStatus checks if the status is one of the allowed values
func (tar *ThreatAssignmentResolution) isValidStatus() bool {
	switch tar.Status {
	case ThreatAssignmentResolutionStatusResolved,
		ThreatAssignmentResolutionStatusAwaiting,
		ThreatAssignmentResolutionStatusAccepted:
		return true
	default:
		return false
	}
}

type ThreatAssignmentResolutionRepository struct {
	db *gorm.DB
}

func NewThreatAssignmentResolutionRepository(db *gorm.DB) *ThreatAssignmentResolutionRepository {
	return &ThreatAssignmentResolutionRepository{db: db}
}

func (r *ThreatAssignmentResolutionRepository) Create(tx *gorm.DB, resolution *ThreatAssignmentResolution) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(resolution).Error
}

func (r *ThreatAssignmentResolutionRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}
	var resolution ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Instance").Preload("Product").Preload("ThreatAssignment.Threat").
		First(&resolution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &resolution, nil
}

func (r *ThreatAssignmentResolutionRepository) Update(tx *gorm.DB, resolution *ThreatAssignmentResolution) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Save(resolution).Error
}

func (r *ThreatAssignmentResolutionRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Delete(&ThreatAssignmentResolution{}, "id = ?", id).Error
}

func (r *ThreatAssignmentResolutionRepository) GetOneByThreatAssignmentID(tx *gorm.DB, threatAssignmentID int) (*ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}
	var resolution ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Instance").Preload("Product").Preload("ThreatAssignment.Threat").
		First(&resolution, "threat_assignment_id = ?", threatAssignmentID).Error
	if err != nil {
		return nil, err
	}
	return &resolution, nil
}

func (r *ThreatAssignmentResolutionRepository) GetOneByThreatAssignmentIDAndInstanceID(tx *gorm.DB, threatAssignmentID int, instanceID uuid.UUID) (*ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}

	var resolution ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Instance").Preload("Product").Preload("ThreatAssignment.Threat").
		First(&resolution, "threat_assignment_id = ? AND instance_id = ?", threatAssignmentID, instanceID).Error
	if err != nil {
		return nil, err
	}
	return &resolution, nil
}

func (r *ThreatAssignmentResolutionRepository) GetOneByAssignmentIDAndProductID(tx *gorm.DB, threatAssignmentID int, productID uuid.UUID) (*ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}

	var resolution ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Instance").Preload("Product").Preload("ThreatAssignment.Threat").
		First(&resolution, "threat_assignment_id = ? AND product_id = ?", threatAssignmentID, productID).Error
	if err != nil {
		return nil, err
	}
	return &resolution, nil
}

func (r *ThreatAssignmentResolutionRepository) ListByProductID(tx *gorm.DB, productID uuid.UUID) ([]ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}

	var resolutions []ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Instance").Preload("Product").Preload("ThreatAssignment.Threat").
		Where("product_id = ?", productID).Find(&resolutions).Error
	if err != nil {
		return nil, err
	}
	return resolutions, nil
}

func (r *ThreatAssignmentResolutionRepository) ListByInstanceID(tx *gorm.DB, instanceID uuid.UUID) ([]ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}

	var resolutions []ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Instance").Preload("Product").Preload("ThreatAssignment.Threat").
		Where("instance_id = ?", instanceID).Find(&resolutions).Error
	if err != nil {
		return nil, err
	}
	return resolutions, nil
}

// ThreatAssignmentResolutionWithDelegation extends ThreatAssignmentResolution with delegation info
type ThreatAssignmentResolutionWithDelegation struct {
	Resolution ThreatAssignmentResolution
	Delegation *ThreatAssignmentResolutionDelegation
}
