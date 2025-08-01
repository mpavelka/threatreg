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

// ThreatWithUnresolvedInstanceCount represents a threat with count of unresolved instances in a domain
type ThreatWithUnresolvedByInstancesCount struct {
	ID                         uuid.UUID `json:"id"`
	Title                      string    `json:"title"`
	Description                string    `json:"description"`
	UnresolvedByInstancesCount int       `json:"unresolved_by_instance_count"`
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

// ListByDomainWithUnresolvedByInstancesCount returns threats assigned to instances in a domain
// with count of instances that haven't resolved the threat
func (r *ThreatRepository) ListByDomainWithUnresolvedByInstancesCount(tx *gorm.DB, domainID uuid.UUID) ([]ThreatWithUnresolvedByInstancesCount, error) {
	if tx == nil {
		tx = r.db
	}

	var results []ThreatWithUnresolvedByInstancesCount

	query := `
		WITH domain_instance_mapping AS (
			-- Get all instances in the domain
			SELECT i.id as instance_id, i.instance_of as product_id
			FROM instances i
			INNER JOIN domain_instances di ON i.id = di.instance_id
			WHERE di.domain_id = ?
		),
		threat_assignments_in_domain AS (
			-- Get all threat assignments for instances in domain (instance-level)
			SELECT DISTINCT 
				ta.threat_id,
				dim.instance_id,
				ta.id as assignment_id,
				dim.product_id,
				ta.product_id as assignment_product_id
			FROM threat_assignments ta
			INNER JOIN domain_instance_mapping dim ON ta.instance_id = dim.instance_id
			WHERE ta.instance_id IS NOT NULL AND ta.instance_id != '00000000-0000-0000-0000-000000000000'
			
			UNION
			
			-- Get all threat assignments for products referenced by instances in domain (product-level)
			SELECT DISTINCT 
				ta.threat_id,
				dim.instance_id,
				ta.id as assignment_id,
				dim.product_id,
				ta.product_id as assignment_product_id
			FROM threat_assignments ta
			INNER JOIN domain_instance_mapping dim ON ta.product_id = dim.product_id
			WHERE ta.product_id IS NOT NULL AND ta.product_id != '00000000-0000-0000-0000-000000000000'
		),
		instance_threat_resolution_status AS (
			-- Determine if each instance has resolved each threat (considering all assignments)
			SELECT 
				tad.threat_id,
				tad.instance_id,
				-- An instance is considered to have resolved a threat if ALL assignments 
				-- of that threat to that instance are resolved
				CASE 
					WHEN COUNT(CASE WHEN tar.id IS NULL OR tar.status NOT IN ('resolved', 'accepted') THEN 1 END) = 0 
					THEN 0  -- All assignments resolved
					ELSE 1  -- At least one assignment unresolved
				END as is_unresolved
			FROM threat_assignments_in_domain tad
			LEFT JOIN threat_assignment_resolutions tar ON (
				tar.threat_assignment_id = tad.assignment_id 
				AND (
					tar.instance_id = tad.instance_id  -- Instance-level resolution
					OR (tar.product_id = tad.assignment_product_id AND (tar.instance_id IS NULL OR tar.instance_id = '00000000-0000-0000-0000-000000000000'))  -- Product-level resolution
				)
			)
			GROUP BY tad.threat_id, tad.instance_id
		)
		SELECT 
			t.id,
			t.title,
			t.description,
			COALESCE(SUM(itrs.is_unresolved), 0) as unresolved_by_instances_count
		FROM threats t
		INNER JOIN instance_threat_resolution_status itrs ON t.id = itrs.threat_id
		GROUP BY t.id, t.title, t.description
		HAVING SUM(itrs.is_unresolved) > 0
		ORDER BY t.title
	`

	err := tx.Raw(query, domainID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
