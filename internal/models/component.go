package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComponentType string

const (
	ComponentTypeProduct  ComponentType = "product"
	ComponentTypeInstance ComponentType = "instance"
)

// Component represents a unified model for products and instances
type Component struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	Name              string             `gorm:"type:varchar(255);index" json:"name"`
	Description       string             `gorm:"type:text" json:"description"`
	Type              ComponentType      `gorm:"type:varchar(20);not null;index" json:"type"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:ComponentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threatAssignments,omitempty"`
	Domains           []Domain           `gorm:"many2many:domain_components;" json:"domains,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Component) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type ComponentRepository struct {
	db *gorm.DB
}

func NewComponentRepository(db *gorm.DB) *ComponentRepository {
	return &ComponentRepository{db: db}
}

func (r *ComponentRepository) Create(tx *gorm.DB, component *Component) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(component).Error
}

func (r *ComponentRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Component, error) {
	if tx == nil {
		tx = r.db
	}
	var component Component
	err := tx.First(&component, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *ComponentRepository) Update(tx *gorm.DB, component *Component) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(component).Error
}

func (r *ComponentRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Component{}, "id = ?", id).Error
}

func (r *ComponentRepository) List(tx *gorm.DB) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}

	var components []Component
	err := tx.Find(&components).Error
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *ComponentRepository) ListByType(tx *gorm.DB, componentType ComponentType) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}

	var components []Component
	err := tx.Where("type = ?", componentType).Find(&components).Error
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *ComponentRepository) Filter(tx *gorm.DB, componentName string) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}

	var components []Component
	query := tx

	// Detect database dialect for case-insensitive comparison
	dialect := tx.Dialector.Name()

	if componentName != "" {
		if dialect == "postgres" {
			query = query.Where("components.name ILIKE ?", "%"+componentName+"%")
		} else {
			query = query.Where("LOWER(components.name) LIKE LOWER(?)", "%"+componentName+"%")
		}
	}

	err := query.Find(&components).Error
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *ComponentRepository) ListByDomainId(tx *gorm.DB, domainID uuid.UUID) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}

	var components []Component
	err := tx.Joins("JOIN domain_components ON components.id = domain_components.component_id").
		Where("domain_components.domain_id = ?", domainID).
		Find(&components).Error

	if err != nil {
		return nil, err
	}
	return components, nil
}

// ComponentWithThreatStats represents a component with threat assignment statistics
type ComponentWithThreatStats struct {
	Component
	UnresolvedThreatCount int `json:"unresolvedThreatCount"`
}

func (r *ComponentRepository) ListByDomainIdWithThreatStats(tx *gorm.DB, domainID uuid.UUID) ([]ComponentWithThreatStats, error) {
	if tx == nil {
		tx = r.db
	}

	var results []ComponentWithThreatStats

	// Query to get components with their unresolved threat counts
	// Count threat assignments that either have no resolution OR have resolution status not in ('accepted', 'resolved')

	query := `
		SELECT 
			c.id,
			c.name,
			c.description,
			c.type,
			COALESCE(
				(SELECT COUNT(DISTINCT ta.id)
				 FROM threat_assignments ta
				 LEFT JOIN threat_assignment_resolutions tar ON ta.id = tar.threat_assignment_id
				 WHERE ta.component_id = c.id
				   AND (tar.id IS NULL OR tar.status NOT IN ('accepted', 'resolved'))
				), 0
			) as unresolved_threat_count
		FROM components c
		JOIN domain_components dc ON c.id = dc.component_id
		WHERE dc.domain_id = ?
	`

	rows, err := tx.Raw(query, domainID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result ComponentWithThreatStats

		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.Description,
			&result.Type,
			&result.UnresolvedThreatCount,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}