package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreatAssignment struct {
	ID                 int                 `gorm:"primaryKey;autoIncrement;not null;unique" json:"id"`
	ThreatID           uuid.UUID           `gorm:"type:uuid;uniqueIndex:idx_threat_assignment" json:"threatID"`
	ProductID          uuid.UUID           `gorm:"type:uuid;uniqueIndex:idx_threat_assignment" json:"productID"`
	InstanceID         uuid.UUID           `gorm:"type:uuid;uniqueIndex:idx_threat_assignment" json:"instanceID"`
	Threat             Threat              `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threat"`
	Product            Product             `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"product"`
	Instance           Instance            `gorm:"foreignKey:InstanceID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"instance"`
	ControlAssignments []ControlAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"controlAssignments"`
}

// BeforeCreate ensures exactly one of ProductID or InstanceID is set
func (ta *ThreatAssignment) BeforeCreate(tx *gorm.DB) error {
	return ta.validateAssignment()
}

// BeforeUpdate ensures exactly one of ProductID or InstanceID is set
func (ta *ThreatAssignment) BeforeUpdate(tx *gorm.DB) error {
	return ta.validateAssignment()
}

// validateAssignment checks that exactly one of ProductID or InstanceID is not null/nil
func (ta *ThreatAssignment) validateAssignment() error {
	productIsSet := ta.ProductID != uuid.Nil
	instanceIsSet := ta.InstanceID != uuid.Nil

	if productIsSet && instanceIsSet {
		return errors.New("threat assignment cannot have both ProductID and InstanceID set")
	}

	if !productIsSet && !instanceIsSet {
		return errors.New("threat assignment must have either ProductID or InstanceID set")
	}

	return nil
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

	// Create new assignment - explicitly set InstanceID to NULL-equivalent
	assignment := &ThreatAssignment{
		ThreatID:   threatID,
		ProductID:  productID,
		InstanceID: uuid.Nil, // Explicitly set to nil UUID
	}

	err := tx.Create(assignment).Error
	if err != nil {
		// Check if this is a unique constraint violation
		if isUniqueConstraintError(err) {
			// Find and return the existing assignment
			var existing ThreatAssignment
			findErr := tx.Where("threat_id = ? AND product_id = ? AND (instance_id IS NULL OR instance_id = ?)", threatID, productID, uuid.Nil).First(&existing).Error
			if findErr == nil {
				return &existing, nil
			}
			// If we can't find the existing record, return the original error
			return nil, err
		}
		return nil, err
	}

	return assignment, nil
}

func (r *ThreatAssignmentRepository) AssignThreatToInstance(tx *gorm.DB, threatID, instanceID uuid.UUID) (*ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}

	// Create new assignment - explicitly set ProductID to NULL-equivalent
	assignment := &ThreatAssignment{
		ThreatID:   threatID,
		InstanceID: instanceID,
		ProductID:  uuid.Nil, // Explicitly set to nil UUID
	}

	err := tx.Create(assignment).Error
	if err != nil {
		// Check if this is a unique constraint violation
		if isUniqueConstraintError(err) {
			// Find and return the existing assignment
			var existing ThreatAssignment
			findErr := tx.Where("threat_id = ? AND instance_id = ? AND (product_id IS NULL OR product_id = ?)", threatID, instanceID, uuid.Nil).First(&existing).Error
			if findErr == nil {
				return &existing, nil
			}
			// If we can't find the existing record, return the original error
			return nil, err
		}
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

// ListWithResolutionByProductID retrieves threat assignments with resolution and delegation status for a product
// The resolutionInstanceID parameter filters which resolutions to include - only resolutions for that specific instance will be shown
func (r *ThreatAssignmentRepository) ListWithResolutionByProductID(tx *gorm.DB, productID, resolutionInstanceID uuid.UUID) ([]ThreatAssignmentWithResolution, error) {
	if tx == nil {
		tx = r.db
	}

	var results []ThreatAssignmentWithResolution

	// Use GORM to join with resolution and delegation tables, similar to ListByProductID
	// Filter resolutions to only include those for the specified resolutionInstanceID
	err := tx.Table("threat_assignments ta").
		Select(`ta.*, 
			tar.status as resolution_status,
			CASE WHEN tard.id IS NOT NULL THEN 1 ELSE 0 END as is_delegated`).
		Joins("LEFT JOIN threat_assignment_resolutions tar ON ta.id = tar.threat_assignment_id AND tar.instance_id = ?", resolutionInstanceID).
		Joins("LEFT JOIN threat_assignment_resolution_delegations tard ON tar.id = tard.delegated_by").
		Where("ta.product_id = ?", productID).
		Preload("Threat").
		Preload("Product").
		Preload("Instance").
		Preload("ControlAssignments").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// ListWithResolutionByInstanceID retrieves threat assignments with resolution and delegation status for an instance
// The resolutionInstanceID parameter filters which resolutions to include - only resolutions for that specific instance will be shown
func (r *ThreatAssignmentRepository) ListWithResolutionByInstanceID(tx *gorm.DB, instanceID, resolutionInstanceID uuid.UUID) ([]ThreatAssignmentWithResolution, error) {
	if tx == nil {
		tx = r.db
	}

	// First get the basic threat assignments for the instance
	var assignments []ThreatAssignment
	err := tx.Preload("Threat").Preload("Product").Preload("Instance").Preload("ControlAssignments").
		Where("instance_id = ?", instanceID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}

	// Convert to ThreatAssignmentWithResolution and add resolution info
	var results []ThreatAssignmentWithResolution
	for _, assignment := range assignments {
		result := ThreatAssignmentWithResolution{
			ThreatAssignment: assignment,
		}

		// Get resolution info for this assignment filtered by resolutionInstanceID
		var resolution struct {
			Status      *string
			IsDelegated bool
		}

		err := tx.Table("threat_assignment_resolutions tar").
			Select("tar.status, CASE WHEN tard.id IS NOT NULL THEN 1 ELSE 0 END as is_delegated").
			Joins("LEFT JOIN threat_assignment_resolution_delegations tard ON tar.id = tard.delegated_by").
			Where("tar.threat_assignment_id = ? AND tar.instance_id = ?", assignment.ID, resolutionInstanceID).
			First(&resolution).Error

		if err == nil && resolution.Status != nil {
			status := ThreatAssignmentResolutionStatus(*resolution.Status)
			result.ResolutionStatus = &status
			result.IsDelegated = resolution.IsDelegated
		}

		results = append(results, result)
	}

	return results, nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	// Check for SQLite unique constraint errors
	if strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "constraint failed") {
		return true
	}
	// Check for PostgreSQL unique constraint errors
	if strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "unique_violation") {
		return true
	}
	return false
}

// ThreatAssignmentWithResolution extends ThreatAssignment with resolution status and delegation info
type ThreatAssignmentWithResolution struct {
	ThreatAssignment
	// Additional fields for resolution and delegation status
	ResolutionStatus *ThreatAssignmentResolutionStatus `json:"resolutionStatus,omitempty"`
	IsDelegated      bool                              `json:"isDelegated"`
}
