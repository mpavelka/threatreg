package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReservedThreatAssignmentRelationshipLabels defines reserved labels for threat assignment relationships
type ReservedThreatAssignmentRelationshipLabels string

const (
	// ReservedInheritsFrom indicates that the child threat assignment inherits threats from the parent
	ReservedInheritsFrom ReservedThreatAssignmentRelationshipLabels = "__inherits_from"
)

// ThreatAssignmentInheritance represents a parent-child relationship between threat assignments
type ThreatAssignmentInheritance struct {
	ID     uuid.UUID        `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	FromID int              `gorm:"not null;index;uniqueIndex:idx_threat_assignment_inheritance_unique" json:"fromId"`
	ToID   int              `gorm:"not null;index;uniqueIndex:idx_threat_assignment_inheritance_unique" json:"toId"`
	From   ThreatAssignment `gorm:"foreignKey:FromID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"from,omitempty"`
	To     ThreatAssignment `gorm:"foreignKey:ToID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"to,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (tai *ThreatAssignmentInheritance) BeforeCreate(tx *gorm.DB) error {
	if tai.ID == uuid.Nil {
		tai.ID = uuid.New()
	}
	return nil
}

// Validate validates the threat assignment inheritance relationship
func (tai *ThreatAssignmentInheritance) Validate() error {
	if tai.FromID == 0 {
		return fmt.Errorf("child threat assignment ID is required")
	}
	if tai.ToID == 0 {
		return fmt.Errorf("parent threat assignment ID is required")
	}
	if tai.FromID == tai.ToID {
		return fmt.Errorf("threat assignment cannot be its own parent")
	}
	return nil
}

// ThreatAssignmentTreePath represents a path in the threat assignment hierarchy tree
type ThreatAssignmentTreePath struct {
	ThreatAssignmentID int   `json:"threatAssignmentId"`
	Path               []int `json:"path"` // Path from root to this threat assignment (including this assignment)
	Depth              int   `json:"depth"`
}

type ThreatAssignmentInheritanceRepository struct {
	db *gorm.DB
}

func NewThreatAssignmentInheritanceRepository(db *gorm.DB) *ThreatAssignmentInheritanceRepository {
	return &ThreatAssignmentInheritanceRepository{db: db}
}

func (r *ThreatAssignmentInheritanceRepository) Create(tx *gorm.DB, inheritance *ThreatAssignmentInheritance) error {
	if tx == nil {
		tx = r.db
	}

	if err := inheritance.Validate(); err != nil {
		return err
	}

	return tx.Create(inheritance).Error
}

func (r *ThreatAssignmentInheritanceRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*ThreatAssignmentInheritance, error) {
	if tx == nil {
		tx = r.db
	}

	var inheritance ThreatAssignmentInheritance
	err := tx.Preload("ChildThreatAssignment").Preload("ParentThreatAssignment").First(&inheritance, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &inheritance, nil
}

func (r *ThreatAssignmentInheritanceRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&ThreatAssignmentInheritance{}, "id = ?", id).Error
}

func (r *ThreatAssignmentInheritanceRepository) GetByChildAndParent(tx *gorm.DB, childID, parentID int) (*ThreatAssignmentInheritance, error) {
	if tx == nil {
		tx = r.db
	}

	var inheritance ThreatAssignmentInheritance
	err := tx.Where("child_threat_assignment_id = ? AND parent_threat_assignment_id = ?", childID, parentID).First(&inheritance).Error
	if err != nil {
		return nil, err
	}
	return &inheritance, nil
}

func (r *ThreatAssignmentInheritanceRepository) ListByChild(tx *gorm.DB, childID int) ([]ThreatAssignmentInheritance, error) {
	if tx == nil {
		tx = r.db
	}

	var inheritances []ThreatAssignmentInheritance
	err := tx.Preload("ParentThreatAssignment").Where("child_threat_assignment_id = ?", childID).Find(&inheritances).Error
	if err != nil {
		return nil, err
	}
	return inheritances, nil
}

func (r *ThreatAssignmentInheritanceRepository) ListByParent(tx *gorm.DB, parentID int) ([]ThreatAssignmentInheritance, error) {
	if tx == nil {
		tx = r.db
	}

	var inheritances []ThreatAssignmentInheritance
	err := tx.Preload("ChildThreatAssignment").Where("parent_threat_assignment_id = ?", parentID).Find(&inheritances).Error
	if err != nil {
		return nil, err
	}
	return inheritances, nil
}

func (r *ThreatAssignmentInheritanceRepository) DeleteByChildAndParent(tx *gorm.DB, childID, parentID int) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Where("child_threat_assignment_id = ? AND parent_threat_assignment_id = ?", childID, parentID).Delete(&ThreatAssignmentInheritance{}).Error
}

// GetTreePaths retrieves all tree paths for a given threat assignment
// This includes all paths from root ancestors to descendants that pass through the given assignment
func (r *ThreatAssignmentInheritanceRepository) GetTreePaths(tx *gorm.DB, threatAssignmentID int) ([]ThreatAssignmentTreePath, error) {
	if tx == nil {
		tx = r.db
	}

	// Get all inheritance relationships
	var inheritances []ThreatAssignmentInheritance
	err := tx.Find(&inheritances).Error
	if err != nil {
		return nil, err
	}

	// Build adjacency maps
	children := make(map[int][]int)
	parents := make(map[int][]int)

	for _, inheritance := range inheritances {
		children[inheritance.ToID] = append(children[inheritance.ToID], inheritance.FromID)
		parents[inheritance.FromID] = append(parents[inheritance.FromID], inheritance.ToID)
	}

	var paths []ThreatAssignmentTreePath

	// Get all paths that include the given threat assignment
	paths = append(paths, r.getPathsIncludingThreatAssignment(threatAssignmentID, children, parents)...)

	return paths, nil
}

// getPathsIncludingThreatAssignment finds all tree paths that include the given threat assignment
func (r *ThreatAssignmentInheritanceRepository) getPathsIncludingThreatAssignment(threatAssignmentID int, children map[int][]int, parents map[int][]int) []ThreatAssignmentTreePath {
	var paths []ThreatAssignmentTreePath

	// Find all root paths to this threat assignment
	rootPaths := r.findPathsToRoot(threatAssignmentID, parents)

	// For each root path, extend it with all descendant paths
	if len(rootPaths) == 0 {
		// Threat assignment is a root itself
		rootPaths = [][]int{{threatAssignmentID}}
	}

	for _, rootPath := range rootPaths {
		// Add paths from this threat assignment to all its descendants
		r.extendPathsToDescendants(rootPath, threatAssignmentID, children, &paths)
	}

	return paths
}

// findPathsToRoot finds all paths from the given threat assignment to its root ancestors
func (r *ThreatAssignmentInheritanceRepository) findPathsToRoot(threatAssignmentID int, parents map[int][]int) [][]int {
	var paths [][]int

	// Base case: if no parents, this is a root
	parentList, hasParents := parents[threatAssignmentID]
	if !hasParents {
		return [][]int{{threatAssignmentID}}
	}

	// Recursive case: get paths through each parent
	for _, parentID := range parentList {
		parentPaths := r.findPathsToRoot(parentID, parents)
		for _, parentPath := range parentPaths {
			// Extend parent path with current threat assignment
			fullPath := append(parentPath, threatAssignmentID)
			paths = append(paths, fullPath)
		}
	}

	return paths
}

// extendPathsToDescendants extends paths from the current threat assignment to all its descendants
func (r *ThreatAssignmentInheritanceRepository) extendPathsToDescendants(currentPath []int, currentAssignment int, children map[int][]int, paths *[]ThreatAssignmentTreePath) {
	// Add current path
	pathCopy := make([]int, len(currentPath))
	copy(pathCopy, currentPath)

	*paths = append(*paths, ThreatAssignmentTreePath{
		ThreatAssignmentID: currentAssignment,
		Path:               pathCopy,
		Depth:              len(pathCopy) - 1,
	})

	// Extend to children
	if childList, hasChildren := children[currentAssignment]; hasChildren {
		for _, childID := range childList {
			newPath := append(currentPath, childID)
			r.extendPathsToDescendants(newPath, childID, children, paths)
		}
	}
}

// GetAllTreePaths gets tree paths for all threat assignments in the system
func (r *ThreatAssignmentInheritanceRepository) GetAllTreePaths(tx *gorm.DB) ([]ThreatAssignmentTreePath, error) {
	if tx == nil {
		tx = r.db
	}

	// Get all inheritance relationships
	var inheritances []ThreatAssignmentInheritance
	err := tx.Find(&inheritances).Error
	if err != nil {
		return nil, err
	}

	// Build adjacency maps
	children := make(map[int][]int)
	parents := make(map[int][]int)
	allAssignments := make(map[int]bool)

	for _, inheritance := range inheritances {
		children[inheritance.ToID] = append(children[inheritance.ToID], inheritance.FromID)
		parents[inheritance.FromID] = append(parents[inheritance.FromID], inheritance.ToID)
		allAssignments[inheritance.ToID] = true
		allAssignments[inheritance.FromID] = true
	}

	var paths []ThreatAssignmentTreePath
	visited := make(map[int]bool)

	// Find and traverse from all root threat assignments
	for assignmentID := range allAssignments {
		if _, hasParent := parents[assignmentID]; !hasParent && !visited[assignmentID] {
			r.traverseFromRoot(assignmentID, []int{assignmentID}, 0, children, &paths, visited)
		}
	}

	return paths, nil
}

// traverseFromRoot recursively traverses the tree from a root threat assignment
func (r *ThreatAssignmentInheritanceRepository) traverseFromRoot(assignmentID int, currentPath []int, depth int, children map[int][]int, paths *[]ThreatAssignmentTreePath, visited map[int]bool) {
	// Add current path
	pathCopy := make([]int, len(currentPath))
	copy(pathCopy, currentPath)

	*paths = append(*paths, ThreatAssignmentTreePath{
		ThreatAssignmentID: assignmentID,
		Path:               pathCopy,
		Depth:              depth,
	})

	visited[assignmentID] = true

	// Traverse children
	if childAssignments, hasChildren := children[assignmentID]; hasChildren {
		for _, child := range childAssignments {
			newPath := append(currentPath, child)
			r.traverseFromRoot(child, newPath, depth+1, children, paths, visited)
		}
	}
}
