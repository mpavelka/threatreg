package models

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComponentAttributeType string

const (
	ComponentAttributeTypeString    ComponentAttributeType = "string"
	ComponentAttributeTypeText      ComponentAttributeType = "text"
	ComponentAttributeTypeNumber    ComponentAttributeType = "number"
	ComponentAttributeTypeComponent ComponentAttributeType = "component"
)

// ComponentAttribute represents a dynamic attribute assigned to a component
type ComponentAttribute struct {
	ID          uuid.UUID              `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	ComponentID uuid.UUID              `gorm:"type:uuid;not null;index" json:"componentId"`
	Name        string                 `gorm:"type:varchar(255);not null;index" json:"name"`
	Type        ComponentAttributeType `gorm:"type:varchar(20);not null" json:"type"`
	Value       string                 `gorm:"type:text;not null" json:"value"`
	Component   Component              `gorm:"foreignKey:ComponentID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"component,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (ca *ComponentAttribute) BeforeCreate(tx *gorm.DB) error {
	if ca.ID == uuid.Nil {
		ca.ID = uuid.New()
	}
	return nil
}

// Validate validates the component attribute based on its type
func (ca *ComponentAttribute) Validate() error {
	if ca.ComponentID == uuid.Nil {
		return fmt.Errorf("component ID is required")
	}
	if ca.Name == "" {
		return fmt.Errorf("attribute name is required")
	}
	if ca.Value == "" {
		return fmt.Errorf("attribute value is required")
	}
	
	switch ca.Type {
	case ComponentAttributeTypeString, ComponentAttributeTypeText:
		// Basic validation - value can be any non-empty string
		return nil
	case ComponentAttributeTypeNumber:
		// For number type, value must be parsable as integer or float
		if _, err := strconv.Atoi(ca.Value); err != nil {
			if _, err := strconv.ParseFloat(ca.Value, 64); err != nil {
				return fmt.Errorf("number attribute value must be a valid integer or float")
			}
		}
		return nil
	case ComponentAttributeTypeComponent:
		// For component type, value must be a valid UUID
		_, err := uuid.Parse(ca.Value)
		if err != nil {
			return fmt.Errorf("component attribute value must be a valid UUID for type 'component'")
		}
		return nil
	default:
		return fmt.Errorf("invalid attribute type: %s", ca.Type)
	}
}

type ComponentAttributeRepository struct {
	db *gorm.DB
}

func NewComponentAttributeRepository(db *gorm.DB) *ComponentAttributeRepository {
	return &ComponentAttributeRepository{db: db}
}

func (r *ComponentAttributeRepository) Create(tx *gorm.DB, attribute *ComponentAttribute) error {
	if tx == nil {
		tx = r.db
	}
	
	if err := attribute.Validate(); err != nil {
		return err
	}
	
	return tx.Create(attribute).Error
}

func (r *ComponentAttributeRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*ComponentAttribute, error) {
	if tx == nil {
		tx = r.db
	}
	
	var attribute ComponentAttribute
	err := tx.Preload("Component").First(&attribute, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &attribute, nil
}

func (r *ComponentAttributeRepository) Update(tx *gorm.DB, attribute *ComponentAttribute) error {
	if tx == nil {
		tx = r.db
	}
	
	if err := attribute.Validate(); err != nil {
		return err
	}
	
	return tx.Save(attribute).Error
}

func (r *ComponentAttributeRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&ComponentAttribute{}, "id = ?", id).Error
}

func (r *ComponentAttributeRepository) ListByComponentID(tx *gorm.DB, componentID uuid.UUID) ([]ComponentAttribute, error) {
	if tx == nil {
		tx = r.db
	}
	
	var attributes []ComponentAttribute
	err := tx.Where("component_id = ?", componentID).Find(&attributes).Error
	if err != nil {
		return nil, err
	}
	return attributes, nil
}

func (r *ComponentAttributeRepository) GetByComponentIDAndName(tx *gorm.DB, componentID uuid.UUID, name string) (*ComponentAttribute, error) {
	if tx == nil {
		tx = r.db
	}
	
	var attribute ComponentAttribute
	err := tx.Where("component_id = ? AND name = ?", componentID, name).First(&attribute).Error
	if err != nil {
		return nil, err
	}
	return &attribute, nil
}

func (r *ComponentAttributeRepository) FindComponentsByAttribute(tx *gorm.DB, name string, value string) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}
	
	var components []Component
	err := tx.Joins("JOIN component_attributes ON components.id = component_attributes.component_id").
		Where("component_attributes.name = ? AND component_attributes.value = ?", name, value).
		Find(&components).Error
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *ComponentAttributeRepository) FindComponentsByAttributeAndType(tx *gorm.DB, name string, value string, attributeType ComponentAttributeType) ([]Component, error) {
	if tx == nil {
		tx = r.db
	}
	
	var components []Component
	err := tx.Joins("JOIN component_attributes ON components.id = component_attributes.component_id").
		Where("component_attributes.name = ? AND component_attributes.value = ? AND component_attributes.type = ?", name, value, attributeType).
		Find(&components).Error
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *ComponentAttributeRepository) ComponentHasAttribute(tx *gorm.DB, componentID uuid.UUID, name string) (bool, error) {
	if tx == nil {
		tx = r.db
	}
	
	var count int64
	err := tx.Model(&ComponentAttribute{}).
		Where("component_id = ? AND name = ?", componentID, name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ComponentAttributeRepository) ComponentHasAttributeWithValue(tx *gorm.DB, componentID uuid.UUID, name string, value string) (bool, error) {
	if tx == nil {
		tx = r.db
	}
	
	var count int64
	err := tx.Model(&ComponentAttribute{}).
		Where("component_id = ? AND name = ? AND value = ?", componentID, name, value).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ComponentAttributeRepository) DeleteByComponentIDAndName(tx *gorm.DB, componentID uuid.UUID, name string) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Where("component_id = ? AND name = ?", componentID, name).Delete(&ComponentAttribute{}).Error
}