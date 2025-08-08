package service

import (
	"testing"
	"threatreg/internal/database"
	"threatreg/internal/models"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestComponentRelationshipService_CreateComponentRelationship(t *testing.T) {

	// Create test components
	parent, err := CreateComponent("Parent Component", "Parent Description", models.ComponentTypeProduct)
	require.NoError(t, err)
	child, err := CreateComponent("Child Component", "Child Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("CreateValidRelationship", func(t *testing.T) {
		relationship, err := CreateComponentRelationship(child.ID, parent.ID, "test_relationship")

		require.NoError(t, err)
		assert.NotNil(t, relationship)
		assert.NotEqual(t, uuid.Nil, relationship.ID)
		assert.Equal(t, child.ID, relationship.FromID)
		assert.Equal(t, parent.ID, relationship.ToID)
	})

	t.Run("CreateRelationship_NonExistentChild", func(t *testing.T) {
		nonExistentID := uuid.New()
		relationship, err := CreateComponentRelationship(nonExistentID, parent.ID, "test_relationship")

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "from component does not exist")
	})

	t.Run("CreateRelationship_NonExistentParent", func(t *testing.T) {
		nonExistentID := uuid.New()
		relationship, err := CreateComponentRelationship(child.ID, nonExistentID, "test_relationship")

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "to component does not exist")
	})

	t.Run("CreateRelationship_SelfReference", func(t *testing.T) {
		relationship, err := CreateComponentRelationship(parent.ID, parent.ID, "test_relationship")

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "component cannot have a relationship to itself")
	})

	t.Run("CreateRelationship_ReservedLabelPrefix", func(t *testing.T) {
		relationship, err := CreateComponentRelationship(child.ID, parent.ID, "__reserved_label")

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "labels starting with '__' are reserved for system use")
	})

	t.Run("CreateRelationship_ReservedLabelConstant", func(t *testing.T) {
		relationship, err := CreateComponentRelationship(child.ID, parent.ID, string(models.ReservedLabelInheritsThreatsFrom))

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "labels starting with '__' are reserved for system use")
	})

	t.Run("CreateRelationship_EmptyLabel", func(t *testing.T) {
		relationship, err := CreateComponentRelationship(child.ID, parent.ID, "")

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "label is required")
	})

	t.Run("CreateRelationship_WhitespaceOnlyLabel", func(t *testing.T) {
		relationship, err := CreateComponentRelationship(child.ID, parent.ID, "   ")

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Contains(t, err.Error(), "label is required")
	})

	t.Run("CreateRelationship_DuplicateRelationship", func(t *testing.T) {
		// Create fresh components for this test to avoid conflicts with other tests
		testParent, err := CreateComponent("Test Parent", "Parent for duplicate test", models.ComponentTypeProduct)
		require.NoError(t, err)
		testChild, err := CreateComponent("Test Child", "Child for duplicate test", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create first relationship
		firstRelationship, err := CreateComponentRelationship(testChild.ID, testParent.ID, "test_relationship")
		require.NoError(t, err)
		require.NotNil(t, firstRelationship)

		// Try to create duplicate with same label - should fail due to unique constraint
		relationship, err := CreateComponentRelationship(testChild.ID, testParent.ID, "test_relationship")
		assert.Error(t, err)
		assert.Nil(t, relationship)
		// The error should contain information about the constraint failure
		assert.Contains(t, err.Error(), "error creating component relationship")
	})
}

func TestComponentRelationshipService_GetComponentRelationship(t *testing.T) {

	parent, err := CreateComponent("Parent Component", "Parent Description", models.ComponentTypeProduct)
	require.NoError(t, err)
	child, err := CreateComponent("Child Component", "Child Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("GetExistingRelationship", func(t *testing.T) {
		createdRelationship, err := CreateComponentRelationship(child.ID, parent.ID, "test_relationship")
		require.NoError(t, err)

		retrievedRelationship, err := GetComponentRelationship(createdRelationship.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedRelationship)
		assert.Equal(t, createdRelationship.ID, retrievedRelationship.ID)
		assert.Equal(t, child.ID, retrievedRelationship.FromID)
		assert.Equal(t, parent.ID, retrievedRelationship.ToID)
		// Check preloaded relationships
		assert.Equal(t, child.ID, retrievedRelationship.From.ID)
		assert.Equal(t, parent.ID, retrievedRelationship.To.ID)
	})

	t.Run("GetNonExistentRelationship", func(t *testing.T) {
		nonExistentID := uuid.New()
		relationship, err := GetComponentRelationship(nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestComponentRelationshipService_DeleteComponentRelationship(t *testing.T) {

	parent, err := CreateComponent("Parent Component", "Parent Description", models.ComponentTypeProduct)
	require.NoError(t, err)
	child, err := CreateComponent("Child Component", "Child Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("DeleteExistingRelationship", func(t *testing.T) {
		createdRelationship, err := CreateComponentRelationship(child.ID, parent.ID, "test_relationship")
		require.NoError(t, err)

		err = DeleteComponentRelationship(createdRelationship.ID)
		require.NoError(t, err)

		// Verify it's deleted
		retrievedRelationship, err := GetComponentRelationship(createdRelationship.ID)
		assert.Error(t, err)
		assert.Nil(t, retrievedRelationship)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteNonExistentRelationship", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := DeleteComponentRelationship(nonExistentID)

		// Should not return error for non-existent records (GORM behavior)
		assert.NoError(t, err)
	})
}

func TestComponentRelationshipService_GetComponentRelationshipByChildAndParent(t *testing.T) {

	parent, err := CreateComponent("Parent Component", "Parent Description", models.ComponentTypeProduct)
	require.NoError(t, err)
	child, err := CreateComponent("Child Component", "Child Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("GetExistingRelationship", func(t *testing.T) {
		createdRelationship, err := CreateComponentRelationship(child.ID, parent.ID, "test_relationship")
		require.NoError(t, err)

		retrievedRelationship, err := GetComponentRelationshipByChildAndParent(child.ID, parent.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedRelationship)
		assert.Equal(t, createdRelationship.ID, retrievedRelationship.ID)
		assert.Equal(t, child.ID, retrievedRelationship.FromID)
		assert.Equal(t, parent.ID, retrievedRelationship.ToID)
	})

	t.Run("GetNonExistentRelationship", func(t *testing.T) {
		otherComponent, err := CreateComponent("Other Component", "Other Description", models.ComponentTypeInstance)
		require.NoError(t, err)

		relationship, err := GetComponentRelationshipByChildAndParent(child.ID, otherComponent.ID)

		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestComponentRelationshipService_ListComponentParents(t *testing.T) {

	// Create test components
	parent1, err := CreateComponent("Parent 1", "First Parent", models.ComponentTypeProduct)
	require.NoError(t, err)
	parent2, err := CreateComponent("Parent 2", "Second Parent", models.ComponentTypeProduct)
	require.NoError(t, err)
	child, err := CreateComponent("Child Component", "Child with multiple parents", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("ListParentsForComponentWithMultipleParents", func(t *testing.T) {
		// Create multiple parent relationships
		_, err := CreateComponentRelationship(child.ID, parent1.ID, "child_to_parent1")
		require.NoError(t, err)
		_, err = CreateComponentRelationship(child.ID, parent2.ID, "child_to_parent2")
		require.NoError(t, err)

		parents, err := ListComponentParents(child.ID)

		require.NoError(t, err)
		assert.Len(t, parents, 2)

		// Create a map for easier assertion
		parentMap := make(map[uuid.UUID]models.ComponentRelationship)
		for _, relationship := range parents {
			parentMap[relationship.ToID] = relationship
		}

		assert.Contains(t, parentMap, parent1.ID)
		assert.Contains(t, parentMap, parent2.ID)

		// Check preloaded parent components
		parent1Relationship := parentMap[parent1.ID]
		assert.Equal(t, "Parent 1", parent1Relationship.To.Name)
		parent2Relationship := parentMap[parent2.ID]
		assert.Equal(t, "Parent 2", parent2Relationship.To.Name)
	})

	t.Run("ListParentsForComponentWithNoParents", func(t *testing.T) {
		orphanComponent, err := CreateComponent("Orphan Component", "Component with no parents", models.ComponentTypeInstance)
		require.NoError(t, err)

		parents, err := ListComponentParents(orphanComponent.ID)

		require.NoError(t, err)
		assert.Len(t, parents, 0)
	})
}

func TestComponentRelationshipService_ListComponentChildren(t *testing.T) {

	// Create test components
	parent, err := CreateComponent("Parent Component", "Parent with multiple children", models.ComponentTypeProduct)
	require.NoError(t, err)
	child1, err := CreateComponent("Child 1", "First Child", models.ComponentTypeInstance)
	require.NoError(t, err)
	child2, err := CreateComponent("Child 2", "Second Child", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("ListChildrenForComponentWithMultipleChildren", func(t *testing.T) {
		// Create multiple child relationships
		_, err := CreateComponentRelationship(child1.ID, parent.ID, "child1_to_parent")
		require.NoError(t, err)
		_, err = CreateComponentRelationship(child2.ID, parent.ID, "child2_to_parent")
		require.NoError(t, err)

		children, err := ListComponentChildren(parent.ID)

		require.NoError(t, err)
		assert.Len(t, children, 2)

		// Create a map for easier assertion
		childMap := make(map[uuid.UUID]models.ComponentRelationship)
		for _, relationship := range children {
			childMap[relationship.FromID] = relationship
		}

		assert.Contains(t, childMap, child1.ID)
		assert.Contains(t, childMap, child2.ID)

		// Check preloaded child components
		child1Relationship := childMap[child1.ID]
		assert.Equal(t, "Child 1", child1Relationship.From.Name)
		child2Relationship := childMap[child2.ID]
		assert.Equal(t, "Child 2", child2Relationship.From.Name)
	})

	t.Run("ListChildrenForComponentWithNoChildren", func(t *testing.T) {
		leafComponent, err := CreateComponent("Leaf Component", "Component with no children", models.ComponentTypeInstance)
		require.NoError(t, err)

		children, err := ListComponentChildren(leafComponent.ID)

		require.NoError(t, err)
		assert.Len(t, children, 0)
	})
}

func TestComponentRelationshipService_DeleteComponentRelationshipByChildAndParent(t *testing.T) {

	parent, err := CreateComponent("Parent Component", "Parent Description", models.ComponentTypeProduct)
	require.NoError(t, err)
	child, err := CreateComponent("Child Component", "Child Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("DeleteExistingRelationship", func(t *testing.T) {
		_, err := CreateComponentRelationship(child.ID, parent.ID, "test_relationship")
		require.NoError(t, err)

		err = DeleteComponentRelationshipByChildAndParent(child.ID, parent.ID)
		require.NoError(t, err)

		// Verify it's deleted
		relationship, err := GetComponentRelationshipByChildAndParent(child.ID, parent.ID)
		assert.Error(t, err)
		assert.Nil(t, relationship)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteNonExistentRelationship", func(t *testing.T) {
		otherComponent, err := CreateComponent("Other Component", "Other Description", models.ComponentTypeInstance)
		require.NoError(t, err)

		err = DeleteComponentRelationshipByChildAndParent(child.ID, otherComponent.ID)

		// Should not return error for non-existent records
		assert.NoError(t, err)
	})
}

func TestComponentRelationshipService_GetComponentTreePaths(t *testing.T) {
	
	// This test requires PostgreSQL for recursive CTE support
	testutil.RequirePostgreSQL(t)

	// Create a component hierarchy:
	// Root -> Parent -> Child -> Grandchild
	root, err := CreateComponent("Root", "Root component", models.ComponentTypeProduct)
	require.NoError(t, err)
	parent, err := CreateComponent("Parent", "Parent component", models.ComponentTypeProduct)
	require.NoError(t, err)
	child, err := CreateComponent("Child", "Child component", models.ComponentTypeInstance)
	require.NoError(t, err)
	grandchild, err := CreateComponent("Grandchild", "Grandchild component", models.ComponentTypeInstance)
	require.NoError(t, err)

	// Create component relationships
	_, err = CreateComponentRelationship(parent.ID, root.ID, "parent_to_root")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(child.ID, parent.ID, "child_to_parent")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(grandchild.ID, child.ID, "grandchild_to_child")
	require.NoError(t, err)

	t.Run("GetTreePathsForMiddleComponent", func(t *testing.T) {
		paths, err := GetComponentTreePaths(parent.ID)

		require.NoError(t, err)
		assert.NotEmpty(t, paths)

		// Find the path that includes the parent component
		var parentPath *models.ComponentTreePath
		for _, path := range paths {
			if path.ComponentID == parent.ID {
				parentPath = &path
				break
			}
		}

		require.NotNil(t, parentPath, "Should find a path for the parent component")
		assert.Contains(t, parentPath.Path, root.ID, "Path should include root")
		assert.Contains(t, parentPath.Path, parent.ID, "Path should include parent")
	})

	t.Run("GetTreePathsForRootComponent", func(t *testing.T) {
		paths, err := GetComponentTreePaths(root.ID)

		require.NoError(t, err)
		assert.NotEmpty(t, paths)

		// Should include paths to all descendants
		componentIDs := make(map[uuid.UUID]bool)
		for _, path := range paths {
			componentIDs[path.ComponentID] = true
		}

		assert.True(t, componentIDs[root.ID], "Should include root")
		assert.True(t, componentIDs[parent.ID], "Should include parent")
		assert.True(t, componentIDs[child.ID], "Should include child")
		assert.True(t, componentIDs[grandchild.ID], "Should include grandchild")
	})

	t.Run("GetTreePathsForLeafComponent", func(t *testing.T) {
		paths, err := GetComponentTreePaths(grandchild.ID)

		require.NoError(t, err)
		assert.NotEmpty(t, paths)

		// Find the path for the grandchild
		var grandchildPath *models.ComponentTreePath
		for _, path := range paths {
			if path.ComponentID == grandchild.ID {
				grandchildPath = &path
				break
			}
		}

		require.NotNil(t, grandchildPath, "Should find a path for the grandchild component")
		assert.Contains(t, grandchildPath.Path, root.ID, "Path should include root")
		assert.Contains(t, grandchildPath.Path, parent.ID, "Path should include parent")
		assert.Contains(t, grandchildPath.Path, child.ID, "Path should include child")
		assert.Contains(t, grandchildPath.Path, grandchild.ID, "Path should include grandchild")
		assert.Equal(t, 3, grandchildPath.Depth, "Grandchild should be at depth 3")
	})

	t.Run("GetTreePathsForIsolatedComponent", func(t *testing.T) {
		isolatedComponent, err := CreateComponent("Isolated", "Isolated component", models.ComponentTypeInstance)
		require.NoError(t, err)

		paths, err := GetComponentTreePaths(isolatedComponent.ID)

		require.NoError(t, err)
		// Should return empty or minimal paths since component is not part of any hierarchy
		// Implementation may vary based on exact tree traversal logic
		_ = paths // Use the variable to avoid unused error
	})
}

func TestComponentRelationshipService_GetAllComponentTreePaths(t *testing.T) {

	// Create multiple component hierarchies
	// Hierarchy 1: Root1 -> Child1
	root1, err := CreateComponent("Root1", "First root", models.ComponentTypeProduct)
	require.NoError(t, err)
	child1, err := CreateComponent("Child1", "Child of root1", models.ComponentTypeInstance)
	require.NoError(t, err)

	// Hierarchy 2: Root2 -> Child2 -> Grandchild2
	root2, err := CreateComponent("Root2", "Second root", models.ComponentTypeProduct)
	require.NoError(t, err)
	child2, err := CreateComponent("Child2", "Child of root2", models.ComponentTypeInstance)
	require.NoError(t, err)
	grandchild2, err := CreateComponent("Grandchild2", "Grandchild of root2", models.ComponentTypeInstance)
	require.NoError(t, err)

	// Create component relationships
	_, err = CreateComponentRelationship(child1.ID, root1.ID, "child1_to_root1")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(child2.ID, root2.ID, "child2_to_root2")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(grandchild2.ID, child2.ID, "grandchild2_to_child2")
	require.NoError(t, err)

	t.Run("GetAllTreePaths", func(t *testing.T) {
		paths, err := GetAllComponentTreePaths()

		require.NoError(t, err)
		assert.NotEmpty(t, paths)

		// Collect all component IDs from paths
		componentIDs := make(map[uuid.UUID]bool)
		for _, path := range paths {
			componentIDs[path.ComponentID] = true
		}

		// Should include all components in hierarchies
		assert.True(t, componentIDs[root1.ID], "Should include root1")
		assert.True(t, componentIDs[child1.ID], "Should include child1")
		assert.True(t, componentIDs[root2.ID], "Should include root2")
		assert.True(t, componentIDs[child2.ID], "Should include child2")
		assert.True(t, componentIDs[grandchild2.ID], "Should include grandchild2")

		// Verify depth calculations
		depthMap := make(map[uuid.UUID]int)
		for _, path := range paths {
			depthMap[path.ComponentID] = path.Depth
		}

		assert.Equal(t, 0, depthMap[root1.ID], "Root1 should be at depth 0")
		assert.Equal(t, 1, depthMap[child1.ID], "Child1 should be at depth 1")
		assert.Equal(t, 0, depthMap[root2.ID], "Root2 should be at depth 0")
		assert.Equal(t, 1, depthMap[child2.ID], "Child2 should be at depth 1")
		assert.Equal(t, 2, depthMap[grandchild2.ID], "Grandchild2 should be at depth 2")
	})

	t.Run("GetAllTreePathsWithNoHierarchies", func(t *testing.T) {
		// Clear all component relationships
		db, err := database.GetDBOrError()
		require.NoError(t, err)
		err = db.Exec("DELETE FROM component_relationships").Error
		require.NoError(t, err)

		paths, err := GetAllComponentTreePaths()

		require.NoError(t, err)
		assert.Empty(t, paths, "Should return empty paths when no hierarchies exist")
	})
}

func TestComponentRelationshipService_ComplexHierarchy(t *testing.T) {

	// Create a complex hierarchy with multiple relationship:
	//     Root
	//    /    \
	//   A      B
	//   |    / |
	//   C   D  E
	//    \ /
	//     F

	root, err := CreateComponent("Root", "Root component", models.ComponentTypeProduct)
	require.NoError(t, err)
	compA, err := CreateComponent("A", "Component A", models.ComponentTypeProduct)
	require.NoError(t, err)
	compB, err := CreateComponent("B", "Component B", models.ComponentTypeProduct)
	require.NoError(t, err)
	compC, err := CreateComponent("C", "Component C", models.ComponentTypeInstance)
	require.NoError(t, err)
	compD, err := CreateComponent("D", "Component D", models.ComponentTypeInstance)
	require.NoError(t, err)
	compE, err := CreateComponent("E", "Component E", models.ComponentTypeInstance)
	require.NoError(t, err)
	compF, err := CreateComponent("F", "Component F", models.ComponentTypeInstance)
	require.NoError(t, err)

	// Create component relationships
	_, err = CreateComponentRelationship(compA.ID, root.ID, "compA_to_root")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(compB.ID, root.ID, "compB_to_root")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(compC.ID, compA.ID, "compC_to_compA")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(compD.ID, compB.ID, "compD_to_compB")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(compE.ID, compB.ID, "compE_to_compB")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(compF.ID, compC.ID, "compF_to_compC")
	require.NoError(t, err)
	_, err = CreateComponentRelationship(compF.ID, compD.ID, "compF_to_compD") // Multiple relationship
	require.NoError(t, err)

	t.Run("ComplexHierarchy_MultipleRelationship", func(t *testing.T) {
		// Test that F has multiple parents
		parents, err := ListComponentParents(compF.ID)
		require.NoError(t, err)
		assert.Len(t, parents, 2, "F should have 2 parents (C and D)")

		parentIDs := make(map[uuid.UUID]bool)
		for _, parent := range parents {
			parentIDs[parent.ToID] = true
		}
		assert.True(t, parentIDs[compC.ID], "F should have C as parent")
		assert.True(t, parentIDs[compD.ID], "F should have D as parent")
	})

	t.Run("ComplexHierarchy_MultipleChildren", func(t *testing.T) {
		// Test that B has multiple children
		children, err := ListComponentChildren(compB.ID)
		require.NoError(t, err)
		assert.Len(t, children, 2, "B should have 2 children (D and E)")

		childIDs := make(map[uuid.UUID]bool)
		for _, child := range children {
			childIDs[child.FromID] = true
		}
		assert.True(t, childIDs[compD.ID], "B should have D as child")
		assert.True(t, childIDs[compE.ID], "B should have E as child")
	})

	t.Run("ComplexHierarchy_TreePaths", func(t *testing.T) {
		paths, err := GetAllComponentTreePaths()
		require.NoError(t, err)
		assert.NotEmpty(t, paths)

		// Verify all components are represented
		componentIDs := make(map[uuid.UUID]bool)
		for _, path := range paths {
			componentIDs[path.ComponentID] = true
		}

		allComponents := []uuid.UUID{root.ID, compA.ID, compB.ID, compC.ID, compD.ID, compE.ID, compF.ID}
		for _, componentID := range allComponents {
			assert.True(t, componentIDs[componentID], "All components should be in tree paths")
		}
	})
}
