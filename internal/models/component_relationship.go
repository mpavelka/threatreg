package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReservedComponentRelationshipLabels defines reserved labels for component relationships
type ReservedComponentRelationshipLabels string

const (
	// ReservedLabelInheritsThreatsFrom indicates that the child component inherits threats from the parent
	ReservedLabelInheritsThreatsFrom ReservedComponentRelationshipLabels = "__inherits_threats_from"
)

// ComponentRelationship represents a parent-child relationship between components
type ComponentRelationship struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;not null;unique" json:"id"`
	FromID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_component_relationship_unique" json:"fromId"`
	ToID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_component_relationship_unique" json:"toId"`
	Label  string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_component_relationship_unique" json:"label"`
	From   Component `gorm:"foreignKey:FromID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"from,omitempty"`
	To     Component `gorm:"foreignKey:ToID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"to,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (cr *ComponentRelationship) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == uuid.Nil {
		cr.ID = uuid.New()
	}
	return nil
}

// Validate validates the component relationship
func (cr *ComponentRelationship) Validate() error {
	if cr.FromID == uuid.Nil {
		return fmt.Errorf("from component ID is required")
	}
	if cr.ToID == uuid.Nil {
		return fmt.Errorf("to component ID is required")
	}
	if cr.FromID == cr.ToID {
		return fmt.Errorf("component cannot have a relationship to itself")
	}
	if strings.TrimSpace(cr.Label) == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}

// ValidateLabel validates that a label doesn't use reserved prefixes
func ValidateLabel(label string) error {
	if strings.HasPrefix(label, "__") {
		return fmt.Errorf("labels starting with '__' are reserved for system use")
	}
	return nil
}

// ComponentTreePath represents a path in the component hierarchy tree
type ComponentTreePath struct {
	ComponentID uuid.UUID   `json:"componentId"`
	Path        []uuid.UUID `json:"path"` // Path from root to this component (including this component)
	Depth       int         `json:"depth"`
}

type ComponentRelationshipRepository struct {
	db *gorm.DB
}

func NewComponentRelationshipRepository(db *gorm.DB) *ComponentRelationshipRepository {
	return &ComponentRelationshipRepository{db: db}
}

func (r *ComponentRelationshipRepository) Create(tx *gorm.DB, relationship *ComponentRelationship) error {
	if tx == nil {
		tx = r.db
	}

	if err := relationship.Validate(); err != nil {
		return err
	}

	return tx.Create(relationship).Error
}

func (r *ComponentRelationshipRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*ComponentRelationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationship ComponentRelationship
	err := tx.Preload("From").Preload("To").First(&relationship, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

func (r *ComponentRelationshipRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Delete(&ComponentRelationship{}, "id = ?", id).Error
}

func (r *ComponentRelationshipRepository) GetByFromAndTo(tx *gorm.DB, fromID, toID uuid.UUID) (*ComponentRelationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationship ComponentRelationship
	err := tx.Where("from_id = ? AND to_id = ?", fromID, toID).First(&relationship).Error
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

func (r *ComponentRelationshipRepository) ListByFrom(tx *gorm.DB, fromID uuid.UUID) ([]ComponentRelationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationships []ComponentRelationship
	err := tx.Preload("To").Where("from_id = ?", fromID).Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

func (r *ComponentRelationshipRepository) ListByTo(tx *gorm.DB, toID uuid.UUID) ([]ComponentRelationship, error) {
	if tx == nil {
		tx = r.db
	}

	var relationships []ComponentRelationship
	err := tx.Preload("From").Where("to_id = ?", toID).Find(&relationships).Error
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

func (r *ComponentRelationshipRepository) DeleteByFromAndTo(tx *gorm.DB, fromID, toID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Where("from_id = ? AND to_id = ?", fromID, toID).Delete(&ComponentRelationship{}).Error
}

// GetTreePaths retrieves all tree paths for a given component
// This includes all paths from root ancestors to descendants that pass through the given component
func (r *ComponentRelationshipRepository) GetTreePaths(tx *gorm.DB, componentID uuid.UUID) ([]ComponentTreePath, error) {
	if tx == nil {
		tx = r.db
	}

	// Get all component relationships
	var relationships []ComponentRelationship
	err := tx.Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	// Build adjacency maps
	children := make(map[uuid.UUID][]uuid.UUID)
	parents := make(map[uuid.UUID][]uuid.UUID)

	for _, relationship := range relationships {
		children[relationship.ToID] = append(children[relationship.ToID], relationship.FromID)
		parents[relationship.FromID] = append(parents[relationship.FromID], relationship.ToID)
	}

	var paths []ComponentTreePath

	// Get all paths that include the given component
	paths = append(paths, r.getPathsIncludingComponent(componentID, children, parents)...)

	return paths, nil
}

// getPathsIncludingComponent finds all tree paths that include the given component
func (r *ComponentRelationshipRepository) getPathsIncludingComponent(componentID uuid.UUID, children map[uuid.UUID][]uuid.UUID, parents map[uuid.UUID][]uuid.UUID) []ComponentTreePath {
	var paths []ComponentTreePath

	// Find all root paths to this component
	rootPaths := r.findPathsToRoot(componentID, parents)

	// For each root path, extend it with all descendant paths
	if len(rootPaths) == 0 {
		// Component is a root itself
		rootPaths = [][]uuid.UUID{{componentID}}
	}

	for _, rootPath := range rootPaths {
		// Add paths from this component to all its descendants
		r.extendPathsToDescendants(rootPath, componentID, children, &paths)
	}

	return paths
}

// findPathsToRoot finds all paths from the given component to its root ancestors
func (r *ComponentRelationshipRepository) findPathsToRoot(componentID uuid.UUID, parents map[uuid.UUID][]uuid.UUID) [][]uuid.UUID {
	var paths [][]uuid.UUID

	// Base case: if no parents, this is a root
	parentList, hasParents := parents[componentID]
	if !hasParents {
		return [][]uuid.UUID{{componentID}}
	}

	// Recursive case: get paths through each parent
	for _, parentID := range parentList {
		parentPaths := r.findPathsToRoot(parentID, parents)
		for _, parentPath := range parentPaths {
			// Extend parent path with current component
			fullPath := append(parentPath, componentID)
			paths = append(paths, fullPath)
		}
	}

	return paths
}

// extendPathsToDescendants extends paths from the current component to all its descendants
func (r *ComponentRelationshipRepository) extendPathsToDescendants(currentPath []uuid.UUID, currentComponent uuid.UUID, children map[uuid.UUID][]uuid.UUID, paths *[]ComponentTreePath) {
	// Add current path
	pathCopy := make([]uuid.UUID, len(currentPath))
	copy(pathCopy, currentPath)

	*paths = append(*paths, ComponentTreePath{
		ComponentID: currentComponent,
		Path:        pathCopy,
		Depth:       len(pathCopy) - 1,
	})

	// Extend to children
	if childList, hasChildren := children[currentComponent]; hasChildren {
		for _, childID := range childList {
			newPath := append(currentPath, childID)
			r.extendPathsToDescendants(newPath, childID, children, paths)
		}
	}
}

// GetAllTreePaths gets tree paths for all components in the system
func (r *ComponentRelationshipRepository) GetAllTreePaths(tx *gorm.DB) ([]ComponentTreePath, error) {
	if tx == nil {
		tx = r.db
	}

	// Get all component relationships
	var relationships []ComponentRelationship
	err := tx.Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	// Build adjacency maps
	children := make(map[uuid.UUID][]uuid.UUID)
	parents := make(map[uuid.UUID][]uuid.UUID)
	allComponents := make(map[uuid.UUID]bool)

	for _, relationship := range relationships {
		children[relationship.ToID] = append(children[relationship.ToID], relationship.FromID)
		parents[relationship.FromID] = append(parents[relationship.FromID], relationship.ToID)
		allComponents[relationship.ToID] = true
		allComponents[relationship.FromID] = true
	}

	var paths []ComponentTreePath
	visited := make(map[uuid.UUID]bool)

	// Find and traverse from all root components
	for componentID := range allComponents {
		if _, hasParent := parents[componentID]; !hasParent && !visited[componentID] {
			r.traverseFromRoot(componentID, []uuid.UUID{componentID}, 0, children, &paths, visited)
		}
	}

	return paths, nil
}

// traverseFromRoot recursively traverses the tree from a root component
func (r *ComponentRelationshipRepository) traverseFromRoot(componentID uuid.UUID, currentPath []uuid.UUID, depth int, children map[uuid.UUID][]uuid.UUID, paths *[]ComponentTreePath, visited map[uuid.UUID]bool) {
	// Add current path
	pathCopy := make([]uuid.UUID, len(currentPath))
	copy(pathCopy, currentPath)

	*paths = append(*paths, ComponentTreePath{
		ComponentID: componentID,
		Path:        pathCopy,
		Depth:       depth,
	})

	visited[componentID] = true

	// Traverse children
	if childComponents, hasChildren := children[componentID]; hasChildren {
		for _, child := range childComponents {
			newPath := append(currentPath, child)
			r.traverseFromRoot(child, newPath, depth+1, children, paths, visited)
		}
	}
}
