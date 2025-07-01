package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreatAssignmentResolutionDelegation struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	DelegatedBy uuid.UUID `gorm:"type:uuid;not null"`
	DelegatedTo uuid.UUID `gorm:"type:uuid;not null"`
}

// BeforeCreate generates a UUID for the ID if not set
func (tard *ThreatAssignmentResolutionDelegation) BeforeCreate(tx *gorm.DB) error {
	if tard.ID == uuid.Nil {
		tard.ID = uuid.New()
	}
	return nil
}

type ThreatAssignmentResolutionDelegationRepository struct {
	db *gorm.DB
}

func NewThreatAssignmentResolutionDelegationRepository(db *gorm.DB) *ThreatAssignmentResolutionDelegationRepository {
	return &ThreatAssignmentResolutionDelegationRepository{db: db}
}

func (r *ThreatAssignmentResolutionDelegationRepository) CreateThreatAssignmentResolutionDelegation(tx *gorm.DB, delegation *ThreatAssignmentResolutionDelegation) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(delegation).Error
}

func (r *ThreatAssignmentResolutionDelegationRepository) DeleteThreatAssignmentResolutionDelegation(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&ThreatAssignmentResolutionDelegation{}, "id = ?", id).Error
}

func (r *ThreatAssignmentResolutionDelegationRepository) GetThreatAssignmentResolutionDelegationById(tx *gorm.DB, id uuid.UUID) (*ThreatAssignmentResolutionDelegation, error) {
	if tx == nil {
		tx = r.db
	}
	var delegation ThreatAssignmentResolutionDelegation
	err := tx.First(&delegation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &delegation, nil
}

func (r *ThreatAssignmentResolutionDelegationRepository) GetThreatAssignmentResolutionDelegations(tx *gorm.DB, delegatedBy *uuid.UUID, delegatedFrom *uuid.UUID) ([]ThreatAssignmentResolutionDelegation, error) {
	if tx == nil {
		tx = r.db
	}

	// Check that at least one parameter is provided
	if delegatedBy == nil && delegatedFrom == nil {
		return nil, errors.New("at least one of delegatedBy or delegatedFrom must be provided")
	}

	var delegations []ThreatAssignmentResolutionDelegation
	query := tx

	// Build query based on provided parameters
	if delegatedBy != nil && delegatedFrom != nil {
		// Both parameters provided - search by both
		query = query.Where("delegated_by = ? AND delegated_to = ?", *delegatedBy, *delegatedFrom)
	} else if delegatedBy != nil {
		// Only delegatedBy provided
		query = query.Where("delegated_by = ?", *delegatedBy)
	} else {
		// Only delegatedFrom provided
		query = query.Where("delegated_to = ?", *delegatedFrom)
	}

	err := query.Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}