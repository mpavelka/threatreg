package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Threat struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	Title             string             `gorm:"type:varchar(255)" json:"title"`
	Description       string             `gorm:"type:text" json:"description"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threatAssignments"`
	ThreatControls    []ThreatControl    `gorm:"foreignKey:ThreatID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threatControls"`
}

// ThreatWithUnresolvedComponentCount represents a threat with count of unresolved components in a domain
type ThreatWithUnresolvedByComponentsCount struct {
	ID                           uuid.UUID `json:"id"`
	Title                        string    `json:"title"`
	Description                  string    `json:"description"`
	UnresolvedByComponentsCount  int       `json:"unresolvedByComponentsCount"`
}

// BeforeCreate is a GORM hook that is triggered before a new record is created.
func (t *Threat) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

type ThreatRepository struct {
	db *gorm.DB
}

func NewThreatRepository(db *gorm.DB) *ThreatRepository {
	return &ThreatRepository{db: db}
}

func (r *ThreatRepository) Create(tx *gorm.DB, threat *Threat) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(threat).Error
}

func (r *ThreatRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Threat, error) {
	if tx == nil {
		tx = r.db
	}
	var threat Threat
	err := tx.First(&threat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &threat, nil
}

func (r *ThreatRepository) Update(tx *gorm.DB, threat *Threat) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(threat).Error
}

func (r *ThreatRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Threat{}, "id = ?", id).Error
}

func (r *ThreatRepository) List(tx *gorm.DB) ([]Threat, error) {
	if tx == nil {
		tx = r.db
	}

	var threats []Threat
	err := tx.Find(&threats).Error
	if err != nil {
		return nil, err
	}
	return threats, nil
}

// ListByDomainWithUnresolvedByComponentsCount returns threats assigned to components in a domain
// with count of components that haven't resolved the threat
func (r *ThreatRepository) ListByDomainWithUnresolvedByComponentsCount(tx *gorm.DB, domainID uuid.UUID) ([]ThreatWithUnresolvedByComponentsCount, error) {
	if tx == nil {
		tx = r.db
	}

	var results []ThreatWithUnresolvedByComponentsCount

	query := `
		WITH domain_components AS (
			-- Get all components in the domain
			SELECT c.id as component_id
			FROM components c
			INNER JOIN domain_components dc ON c.id = dc.component_id
			WHERE dc.domain_id = ?
		),
		threat_assignments_in_domain AS (
			-- Get all threat assignments for components in domain
			SELECT DISTINCT 
				ta.threat_id,
				dc.component_id,
				ta.id as assignment_id
			FROM threat_assignments ta
			INNER JOIN domain_components dc ON ta.component_id = dc.component_id
		),
		component_threat_resolution_status AS (
			-- Determine if each component has resolved each threat
			SELECT 
				tad.threat_id,
				tad.component_id,
				CASE 
					WHEN COUNT(CASE WHEN tar.id IS NULL OR tar.status NOT IN ('resolved', 'accepted') THEN 1 END) = 0 
					THEN 0  -- All assignments resolved
					ELSE 1  -- At least one assignment unresolved
				END as is_unresolved
			FROM threat_assignments_in_domain tad
			LEFT JOIN threat_assignment_resolutions tar ON (
				tar.threat_assignment_id = tad.assignment_id 
				AND tar.component_id = tad.component_id
			)
			GROUP BY tad.threat_id, tad.component_id
		)
		SELECT 
			t.id,
			t.title,
			t.description,
			COALESCE(SUM(ctrs.is_unresolved), 0) as unresolved_by_components_count
		FROM threats t
		INNER JOIN component_threat_resolution_status ctrs ON t.id = ctrs.threat_id
		GROUP BY t.id, t.title, t.description
		HAVING SUM(ctrs.is_unresolved) > 0
		ORDER BY t.title
	`

	err := tx.Raw(query, domainID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
