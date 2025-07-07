package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreatPattern struct {
	ID          uuid.UUID          `gorm:"type:uuid;primaryKey;not null;unique"`
	Name        string             `gorm:"type:varchar(255);not null"`
	Description string             `gorm:"type:text"`
	ThreatID    uuid.UUID          `gorm:"type:uuid;not null"`
	IsActive    bool               `gorm:"not null"`
	Threat      Threat             `gorm:"foreignKey:ThreatID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Conditions  []PatternCondition `gorm:"foreignKey:PatternID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (tp *ThreatPattern) BeforeCreate(tx *gorm.DB) error {
	if tp.ID == uuid.Nil {
		tp.ID = uuid.New()
	}
	return nil
}

type PatternCondition struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;not null;unique"`
	PatternID        uuid.UUID `gorm:"type:uuid;not null"`
	ConditionType    string    `gorm:"type:varchar(50);not null"` // 'PRODUCT', 'TAG', 'RELATIONSHIP', 'RELATIONSHIP_TARGET_ID', 'RELATIONSHIP_TARGET_TAG', 'PRODUCT_TAG', 'PRODUCT_ID'
	Operator         string    `gorm:"type:varchar(20);not null"` // 'EQUALS', 'CONTAINS', 'NOT_EQUALS', 'EXISTS', 'NOT_EXISTS', 'HAS_RELATIONSHIP_WITH', 'NOT_HAS_RELATIONSHIP_WITH'
	Value            string    `gorm:"type:varchar(255)"`         // the value to match against (tag name, entity ID, product name, etc.)
	RelationshipType string    `gorm:"type:varchar(100)"`         // relationship type when condition involves relationships
}

func (pc *PatternCondition) BeforeCreate(tx *gorm.DB) error {
	if pc.ID == uuid.Nil {
		pc.ID = uuid.New()
	}
	return nil
}

type ThreatPatternRepository struct {
	db *gorm.DB
}

func NewThreatPatternRepository(db *gorm.DB) *ThreatPatternRepository {
	return &ThreatPatternRepository{db: db}
}

func (r *ThreatPatternRepository) Create(tx *gorm.DB, pattern *ThreatPattern) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(pattern).Error
}

func (r *ThreatPatternRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*ThreatPattern, error) {
	if tx == nil {
		tx = r.db
	}
	var pattern ThreatPattern
	err := tx.Preload("Threat").Preload("Conditions").First(&pattern, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pattern, nil
}

func (r *ThreatPatternRepository) Update(tx *gorm.DB, pattern *ThreatPattern) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(pattern).Error
}

func (r *ThreatPatternRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&ThreatPattern{}, "id = ?", id).Error
}

func (r *ThreatPatternRepository) List(tx *gorm.DB) ([]ThreatPattern, error) {
	if tx == nil {
		tx = r.db
	}

	var patterns []ThreatPattern
	err := tx.Preload("Threat").Preload("Conditions").Find(&patterns).Error
	if err != nil {
		return nil, err
	}
	return patterns, nil
}

func (r *ThreatPatternRepository) ListActive(tx *gorm.DB) ([]ThreatPattern, error) {
	if tx == nil {
		tx = r.db
	}

	var patterns []ThreatPattern
	err := tx.Preload("Threat").Preload("Conditions").Where("is_active = ?", true).Find(&patterns).Error
	if err != nil {
		return nil, err
	}
	return patterns, nil
}

func (r *ThreatPatternRepository) ListByThreatID(tx *gorm.DB, threatID uuid.UUID) ([]ThreatPattern, error) {
	if tx == nil {
		tx = r.db
	}

	var patterns []ThreatPattern
	err := tx.Preload("Threat").Preload("Conditions").Where("threat_id = ?", threatID).Find(&patterns).Error
	if err != nil {
		return nil, err
	}
	return patterns, nil
}

func (r *ThreatPatternRepository) SetActive(tx *gorm.DB, id uuid.UUID, isActive bool) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Model(&ThreatPattern{}).Where("id = ?", id).Update("is_active", isActive).Error
}

type PatternConditionRepository struct {
	db *gorm.DB
}

func NewPatternConditionRepository(db *gorm.DB) *PatternConditionRepository {
	return &PatternConditionRepository{db: db}
}

func (r *PatternConditionRepository) Create(tx *gorm.DB, condition *PatternCondition) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(condition).Error
}

func (r *PatternConditionRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*PatternCondition, error) {
	if tx == nil {
		tx = r.db
	}
	var condition PatternCondition
	err := tx.First(&condition, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &condition, nil
}

func (r *PatternConditionRepository) Update(tx *gorm.DB, condition *PatternCondition) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(condition).Error
}

func (r *PatternConditionRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&PatternCondition{}, "id = ?", id).Error
}

func (r *PatternConditionRepository) ListByPatternID(tx *gorm.DB, patternID uuid.UUID) ([]PatternCondition, error) {
	if tx == nil {
		tx = r.db
	}

	var conditions []PatternCondition
	err := tx.Where("pattern_id = ?", patternID).Find(&conditions).Error
	if err != nil {
		return nil, err
	}
	return conditions, nil
}

func (r *PatternConditionRepository) DeleteByPatternID(tx *gorm.DB, patternID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&PatternCondition{}, "pattern_id = ?", patternID).Error
}

func (r *PatternConditionRepository) List(tx *gorm.DB) ([]PatternCondition, error) {
	if tx == nil {
		tx = r.db
	}

	var conditions []PatternCondition
	err := tx.Find(&conditions).Error
	if err != nil {
		return nil, err
	}
	return conditions, nil
}
