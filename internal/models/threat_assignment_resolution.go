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
	ID                 uuid.UUID                        `gorm:"type:uuid;primaryKey" json:"id"`
	ThreatAssignmentID uuid.UUID                        `gorm:"type:uuid;not null;uniqueIndex:idx_threat_assignment_resolution" json:"threatAssignmentId"`
	ComponentID        uuid.UUID                        `gorm:"type:uuid;uniqueIndex:idx_threat_assignment_resolution" json:"componentId"`
	Status             ThreatAssignmentResolutionStatus `gorm:"type:varchar(20);not null" json:"status"`
	Description        string                           `gorm:"type:text" json:"description"`
	ThreatAssignment   ThreatAssignment                 `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"threatAssignment,omitempty"`
	Component          Component                        `gorm:"foreignKey:ComponentID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"component,omitempty"`
}

// BeforeCreate ensures ComponentID is set and validates status
func (tar *ThreatAssignmentResolution) BeforeCreate(tx *gorm.DB) error {
	if tar.ID == uuid.Nil {
		tar.ID = uuid.New()
	}
	return tar.validateResolution()
}

// BeforeUpdate ensures ComponentID is set and validates status
func (tar *ThreatAssignmentResolution) BeforeUpdate(tx *gorm.DB) error {
	return tar.validateResolution()
}

// validateResolution checks that ComponentID is not null/nil and validates status
func (tar *ThreatAssignmentResolution) validateResolution() error {
	if tar.ComponentID == uuid.Nil {
		return errors.New("threat assignment resolution must have ComponentID set")
	}

	// Validate status
	if !tar.isValidStatus() {
		return errors.New("invalid status: must be 'resolved', 'awaiting', or 'accepted'")
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
	err := tx.Preload("ThreatAssignment").Preload("Component").Preload("ThreatAssignment.Threat").
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

func (r *ThreatAssignmentResolutionRepository) GetOneByThreatAssignmentID(tx *gorm.DB, threatAssignmentID uuid.UUID) (*ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}
	var resolution ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Component").Preload("ThreatAssignment.Threat").
		First(&resolution, "threat_assignment_id = ?", threatAssignmentID).Error
	if err != nil {
		return nil, err
	}
	return &resolution, nil
}

func (r *ThreatAssignmentResolutionRepository) GetOneByThreatAssignmentIDAndComponentID(tx *gorm.DB, threatAssignmentID uuid.UUID, componentID uuid.UUID) (*ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}

	var resolution ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Component").Preload("ThreatAssignment.Threat").
		First(&resolution, "threat_assignment_id = ? AND component_id = ?", threatAssignmentID, componentID).Error
	if err != nil {
		return nil, err
	}
	return &resolution, nil
}

func (r *ThreatAssignmentResolutionRepository) ListByComponentID(tx *gorm.DB, componentID uuid.UUID) ([]ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = r.db
	}

	var resolutions []ThreatAssignmentResolution
	err := tx.Preload("ThreatAssignment").Preload("Component").Preload("ThreatAssignment.Threat").
		Where("component_id = ?", componentID).Find(&resolutions).Error
	if err != nil {
		return nil, err
	}
	return resolutions, nil
}

// ThreatAssignmentResolutionWithDelegation extends ThreatAssignmentResolution with delegation info
type ThreatAssignmentResolutionWithDelegation struct {
	Resolution ThreatAssignmentResolution               `json:"resolution"`
	Delegation *ThreatAssignmentResolutionDelegation    `json:"delegation,omitempty"`
}
