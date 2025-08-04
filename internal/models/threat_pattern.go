package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PatternConditionType represents the type of pattern condition
type PatternConditionType int

const (
	ConditionTypeTag PatternConditionType = iota
	ConditionTypeRelationship
	ConditionTypeRelationshipTargetID
	ConditionTypeRelationshipTargetTag
)

// String returns the string representation of the condition type
func (ct PatternConditionType) String() string {
	switch ct {
	case ConditionTypeTag:
		return "TAG"
	case ConditionTypeRelationship:
		return "RELATIONSHIP"
	case ConditionTypeRelationshipTargetID:
		return "RELATIONSHIP_TARGET_ID"
	case ConditionTypeRelationshipTargetTag:
		return "RELATIONSHIP_TARGET_TAG"
	default:
		return ""
	}
}

// ParsePatternConditionType parses a string to PatternConditionType
func ParsePatternConditionType(s string) (PatternConditionType, bool) {
	switch s {
	case "TAG":
		return ConditionTypeTag, true
	case "RELATIONSHIP":
		return ConditionTypeRelationship, true
	case "RELATIONSHIP_TARGET_ID":
		return ConditionTypeRelationshipTargetID, true
	case "RELATIONSHIP_TARGET_TAG":
		return ConditionTypeRelationshipTargetTag, true
	default:
		return 0, false
	}
}

// PatternOperator represents the operator for pattern conditions
type PatternOperator int

const (
	OperatorEquals PatternOperator = iota
	OperatorContains
	OperatorNotContains
	OperatorNotEquals
	OperatorExists
	OperatorNotExists
	OperatorHasRelationshipWith
	OperatorNotHasRelationshipWith
)

// String returns the string representation of the operator
func (op PatternOperator) String() string {
	switch op {
	case OperatorEquals:
		return "EQUALS"
	case OperatorContains:
		return "CONTAINS"
	case OperatorNotContains:
		return "NOT_CONTAINS"
	case OperatorNotEquals:
		return "NOT_EQUALS"
	case OperatorExists:
		return "EXISTS"
	case OperatorNotExists:
		return "NOT_EXISTS"
	case OperatorHasRelationshipWith:
		return "HAS_RELATIONSHIP_WITH"
	case OperatorNotHasRelationshipWith:
		return "NOT_HAS_RELATIONSHIP_WITH"
	default:
		return ""
	}
}

// ParsePatternOperator parses a string to PatternOperator
func ParsePatternOperator(s string) (PatternOperator, bool) {
	switch s {
	case "EQUALS":
		return OperatorEquals, true
	case "CONTAINS":
		return OperatorContains, true
	case "NOT_CONTAINS":
		return OperatorNotContains, true
	case "NOT_EQUALS":
		return OperatorNotEquals, true
	case "EXISTS":
		return OperatorExists, true
	case "NOT_EXISTS":
		return OperatorNotExists, true
	case "HAS_RELATIONSHIP_WITH":
		return OperatorHasRelationshipWith, true
	case "NOT_HAS_RELATIONSHIP_WITH":
		return OperatorNotHasRelationshipWith, true
	default:
		return 0, false
	}
}

// GetAllConditionTypes returns all valid pattern condition types
func GetAllConditionTypes() []PatternConditionType {
	return []PatternConditionType{
		ConditionTypeTag,
		ConditionTypeRelationship,
		ConditionTypeRelationshipTargetID,
		ConditionTypeRelationshipTargetTag,
	}
}

// GetAllOperators returns all valid pattern operators
func GetAllOperators() []PatternOperator {
	return []PatternOperator{
		OperatorEquals,
		OperatorContains,
		OperatorNotContains,
		OperatorNotEquals,
		OperatorExists,
		OperatorNotExists,
		OperatorHasRelationshipWith,
		OperatorNotHasRelationshipWith,
	}
}

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
