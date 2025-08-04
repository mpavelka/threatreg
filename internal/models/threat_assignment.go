package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreatAssignment struct {
	ID                 int                 `gorm:"primaryKey;autoIncrement;not null;unique" json:"id"`
	ThreatID           uuid.UUID           `gorm:"type:uuid;uniqueIndex:idx_threat_assignment" json:"threatId"`
	ComponentID        uuid.UUID           `gorm:"type:uuid;uniqueIndex:idx_threat_assignment" json:"componentId"`
	Threat             Threat              `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threat"`
	Component          Component           `gorm:"foreignKey:ComponentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"component"`
	ControlAssignments []ControlAssignment `gorm:"foreignKey:ThreatAssignmentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"controlAssignments"`
}

// BeforeCreate ensures ComponentID is set
func (ta *ThreatAssignment) BeforeCreate(tx *gorm.DB) error {
	return ta.validateAssignment()
}

// BeforeUpdate ensures ComponentID is set
func (ta *ThreatAssignment) BeforeUpdate(tx *gorm.DB) error {
	return ta.validateAssignment()
}

// validateAssignment checks that ComponentID is not null/nil
func (ta *ThreatAssignment) validateAssignment() error {
	if ta.ComponentID == uuid.Nil {
		return errors.New("threat assignment must have ComponentID set")
	}

	return nil
}

type ThreatAssignmentRepository struct {
	db *gorm.DB
}

func NewThreatAssignmentRepository(db *gorm.DB) *ThreatAssignmentRepository {
	return &ThreatAssignmentRepository{db: db}
}

func (r *ThreatAssignmentRepository) AssignThreatToComponent(tx *gorm.DB, threatID, componentID uuid.UUID) (*ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}

	// Create new assignment
	assignment := &ThreatAssignment{
		ThreatID:    threatID,
		ComponentID: componentID,
	}

	err := tx.Create(assignment).Error
	if err != nil {
		// Check if this is a unique constraint violation
		if isUniqueConstraintError(err) {
			// Find and return the existing assignment
			var existing ThreatAssignment
			findErr := tx.Where("threat_id = ? AND component_id = ?", threatID, componentID).First(&existing).Error
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
	err := tx.Preload("Threat").Preload("Component").First(&assignment, "id = ?", id).Error
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

func (r *ThreatAssignmentRepository) ListByComponentID(tx *gorm.DB, componentID uuid.UUID) ([]ThreatAssignment, error) {
	if tx == nil {
		tx = r.db
	}

	var assignments []ThreatAssignment
	err := tx.Preload("Threat").Preload("Component").
		Where("component_id = ?", componentID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// ListWithResolutionByComponentID retrieves threat assignments with resolution and delegation status for a component
// The resolutionComponentID parameter filters which resolutions to include - only resolutions for that specific component will be shown
func (r *ThreatAssignmentRepository) ListWithResolutionByComponentID(tx *gorm.DB, componentID, resolutionComponentID uuid.UUID) ([]ThreatAssignmentWithResolution, error) {
	if tx == nil {
		tx = r.db
	}

	// First get the basic threat assignments for the component
	var assignments []ThreatAssignment
	err := tx.Preload("Threat").Preload("Component").Preload("ControlAssignments").
		Where("component_id = ?", componentID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}

	// Convert to ThreatAssignmentWithResolution and add resolution info
	var results []ThreatAssignmentWithResolution
	for _, assignment := range assignments {
		result := ThreatAssignmentWithResolution{
			ThreatAssignment: assignment,
		}

		// Get resolution info for this assignment filtered by resolutionComponentID
		var resolution struct {
			Status      *string
			IsDelegated bool
		}

		err := tx.Table("threat_assignment_resolutions tar").
			Select("tar.status, CASE WHEN tard.id IS NOT NULL THEN 1 ELSE 0 END as is_delegated").
			Joins("LEFT JOIN threat_assignment_resolution_delegations tard ON tar.id = tard.delegated_by").
			Where("tar.threat_assignment_id = ? AND tar.component_id = ?", assignment.ID, resolutionComponentID).
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
