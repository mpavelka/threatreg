package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Instance struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	Name              string             `gorm:"type:varchar(255);index" json:"name"`
	InstanceOf        uuid.UUID          `gorm:"type:uuid" json:"instanceOf"`
	Product           Product            `gorm:"foreignKey:InstanceOf;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"product"`
	ThreatAssignments []ThreatAssignment `gorm:"foreignKey:InstanceID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"threatAssignments"`
	Domains           []Domain           `gorm:"many2many:domain_instances;" json:"domains"`
}

func (i *Instance) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}

type InstanceRepository struct {
	db *gorm.DB
}

func NewInstanceRepository(db *gorm.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) Create(tx *gorm.DB, instance *Instance) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(instance).Error
}

func (r *InstanceRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*Instance, error) {
	if tx == nil {
		tx = r.db
	}
	var instance Instance
	err := tx.Preload("Product").First(&instance, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *InstanceRepository) Update(tx *gorm.DB, instance *Instance) error {
	if tx == nil {
		tx = r.db
	}
	// Use Updates to explicitly update both name and foreign key field
	// GORM may not update foreign key fields with Save due to the relationship
	return tx.Model(instance).Select("name", "instance_of").Updates(instance).Error
}

func (r *InstanceRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&Instance{}, "id = ?", id).Error
}

func (r *InstanceRepository) List(tx *gorm.DB) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *InstanceRepository) ListByProductID(tx *gorm.DB, productID uuid.UUID) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").Where("instance_of = ?", productID).Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *InstanceRepository) Filter(tx *gorm.DB, instanceName, productName string) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	query := tx.Preload("Product")

	// Detect database dialect for case-insensitive comparison
	dialect := tx.Dialector.Name()

	if instanceName != "" {
		if dialect == "postgres" {
			query = query.Where("instances.name ILIKE ?", "%"+instanceName+"%")
		} else {
			query = query.Where("LOWER(instances.name) LIKE LOWER(?)", "%"+instanceName+"%")
		}
	}

	if productName != "" {
		if dialect == "postgres" {
			query = query.Where("Product.name ILIKE ?", "%"+productName+"%").Joins("Product")
		} else {
			query = query.Where("LOWER(Product.name) LIKE LOWER(?)", "%"+productName+"%").Joins("Product")
		}
	}

	err := query.Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *InstanceRepository) ListByDomainId(tx *gorm.DB, domainID uuid.UUID) ([]Instance, error) {
	if tx == nil {
		tx = r.db
	}

	var instances []Instance
	err := tx.Preload("Product").
		Joins("JOIN domain_instances ON instances.id = domain_instances.instance_id").
		Where("domain_instances.domain_id = ?", domainID).
		Find(&instances).Error

	if err != nil {
		return nil, err
	}
	return instances, nil
}

// InstanceWithThreatStats represents an instance with threat assignment statistics
type InstanceWithThreatStats struct {
	Instance
	UnresolvedThreatCount int `json:"unresolved_threat_count"`
}

func (r *InstanceRepository) ListByDomainIdWithThreatStats(tx *gorm.DB, domainID uuid.UUID) ([]InstanceWithThreatStats, error) {
	if tx == nil {
		tx = r.db
	}

	var results []InstanceWithThreatStats

	// Complex query to get instances with their unresolved threat counts
	// Count threat assignments that either:
	// 1. Reference the instance ID directly, OR
	// 2. Reference the product ID of the instance
	// But only count those that either have no resolution OR have resolution status not in ('accepted', 'resolved')

	query := `
		SELECT 
			i.id,
			i.name,
			i.instance_of,
			p.id as "Product__id",
			p.name as "Product__name", 
			p.description as "Product__description",
			COALESCE(
				(SELECT COUNT(DISTINCT ta.id)
				 FROM threat_assignments ta
				 LEFT JOIN threat_assignment_resolutions tar ON ta.id = tar.threat_assignment_id
				 WHERE (ta.instance_id = i.id OR ta.product_id = i.instance_of)
				   AND (tar.id IS NULL OR tar.status NOT IN ('accepted', 'resolved'))
				), 0
			) as unresolved_threat_count
		FROM instances i
		JOIN domain_instances di ON i.id = di.instance_id
		LEFT JOIN products p ON i.instance_of = p.id
		WHERE di.domain_id = ?
	`

	rows, err := tx.Raw(query, domainID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result InstanceWithThreatStats
		var productID, productName, productDescription *string

		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.InstanceOf,
			&productID,
			&productName,
			&productDescription,
			&result.UnresolvedThreatCount,
		)
		if err != nil {
			return nil, err
		}

		// Set product details if they exist
		if productID != nil {
			if err := result.Product.ID.UnmarshalText([]byte(*productID)); err != nil {
				return nil, err
			}
		}
		if productName != nil {
			result.Product.Name = *productName
		}
		if productDescription != nil {
			result.Product.Description = *productDescription
		}

		results = append(results, result)
	}

	return results, nil
}
