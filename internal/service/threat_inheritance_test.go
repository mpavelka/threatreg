package service

import (
	"testing"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThreatInheritance_Integration(t *testing.T) {

	t.Run("BasicInheritance", func(t *testing.T) {
		// Create test components: Parent -> Child
		parent, err := CreateComponent("Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child, err := CreateComponent("Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat
		threat, err := CreateThreat("Test Threat", "Test threat for inheritance")
		require.NoError(t, err)

		// Create inherits-threats relationship from child to parent
		_, err = CreateInheritsThreatsRelationship(child.ID, parent.ID)
		require.NoError(t, err)

		// Assign threat to parent and set severity (ResidualSeverity should auto-match)
		parentAssignment, err := AssignThreatToComponent(parent.ID, threat.ID)
		require.NoError(t, err)

		severity := models.ThreatSeverityMedium
		err = SetThreatAssignmentSeverity(parentAssignment.ID, &severity)
		require.NoError(t, err)

		// Verify parent has both Severity and ResidualSeverity set to the same value
		updatedParent, err := GetThreatAssignmentById(parentAssignment.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedParent.Severity)
		require.NotNil(t, updatedParent.ResidualSeverity)
		assert.Equal(t, severity, *updatedParent.Severity)
		assert.Equal(t, severity, *updatedParent.ResidualSeverity)

		// Verify child automatically got the threat assignment with parent's residual severity
		childAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		require.Len(t, childAssignments, 1)

		childAssignment := childAssignments[0]
		assert.Equal(t, threat.ID, childAssignment.ThreatID)
		assert.Equal(t, child.ID, childAssignment.ComponentID)
		require.NotNil(t, childAssignment.Severity)
		assert.Equal(t, severity, *childAssignment.Severity)

		// Verify inheritance relationship exists
		threatRelRepo, err := getThreatAssignmentRelationshipRepository()
		require.NoError(t, err)

		inheritanceRels, err := threatRelRepo.ListByFromAndToAndLabel(
			nil,
			childAssignment.ID,
			parentAssignment.ID,
			string(models.ReservedInheritsFrom),
		)
		require.NoError(t, err)
		assert.Len(t, inheritanceRels, 1)
	})

	t.Run("FullTreePropagationMultipleChildren", func(t *testing.T) {
		// Create complex tree structure:
		//                   Root
		//                /        \
		//           Child1          Child2
		//          /      \            \
		//    Grandchild1  Grandchild2   Grandchild3
		//                              /
		//                      GreatGrandchild1

		root, err := CreateComponent("Root Component", "Root component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child1, err := CreateComponent("Child1 Component", "First child component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child2, err := CreateComponent("Child2 Component", "Second child component", models.ComponentTypeProduct)
		require.NoError(t, err)

		grandchild1, err := CreateComponent("Grandchild1 Component", "First grandchild", models.ComponentTypeInstance)
		require.NoError(t, err)

		grandchild2, err := CreateComponent("Grandchild2 Component", "Second grandchild", models.ComponentTypeInstance)
		require.NoError(t, err)

		grandchild3, err := CreateComponent("Grandchild3 Component", "Third grandchild", models.ComponentTypeInstance)
		require.NoError(t, err)

		greatGrandchild1, err := CreateComponent("GreatGrandchild1 Component", "Great grandchild", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat
		threat, err := CreateThreat("Tree Propagation Threat", "Test threat for full tree propagation")
		require.NoError(t, err)

		// Create inherits-threats relationships to build the tree
		_, err = CreateInheritsThreatsRelationship(child1.ID, root.ID)
		require.NoError(t, err)
		_, err = CreateInheritsThreatsRelationship(child2.ID, root.ID)
		require.NoError(t, err)
		_, err = CreateInheritsThreatsRelationship(grandchild1.ID, child1.ID)
		require.NoError(t, err)
		_, err = CreateInheritsThreatsRelationship(grandchild2.ID, child1.ID)
		require.NoError(t, err)
		_, err = CreateInheritsThreatsRelationship(grandchild3.ID, child2.ID)
		require.NoError(t, err)
		_, err = CreateInheritsThreatsRelationship(greatGrandchild1.ID, grandchild3.ID)
		require.NoError(t, err)

		// Assign threat to root with severity
		rootAssignment, err := AssignThreatToComponent(root.ID, threat.ID)
		require.NoError(t, err)

		rootSeverity := models.ThreatSeverityHigh
		err = SetThreatAssignmentSeverity(rootAssignment.ID, &rootSeverity)
		require.NoError(t, err)

		// Verify root has correct severity and residual severity
		updatedRoot, err := GetThreatAssignmentById(rootAssignment.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedRoot.Severity)
		require.NotNil(t, updatedRoot.ResidualSeverity)
		assert.Equal(t, rootSeverity, *updatedRoot.Severity)
		assert.Equal(t, rootSeverity, *updatedRoot.ResidualSeverity)

		// Verify both children got the threat assignment with root's residual severity
		child1Assignments, err := ListThreatAssignmentsByComponentID(child1.ID)
		require.NoError(t, err)
		require.Len(t, child1Assignments, 1)
		child1Assignment := child1Assignments[0]
		require.NotNil(t, child1Assignment.Severity)
		require.NotNil(t, child1Assignment.ResidualSeverity)
		assert.Equal(t, rootSeverity, *child1Assignment.Severity)
		assert.Equal(t, rootSeverity, *child1Assignment.ResidualSeverity)

		child2Assignments, err := ListThreatAssignmentsByComponentID(child2.ID)
		require.NoError(t, err)
		require.Len(t, child2Assignments, 1)
		child2Assignment := child2Assignments[0]
		require.NotNil(t, child2Assignment.Severity)
		require.NotNil(t, child2Assignment.ResidualSeverity)
		assert.Equal(t, rootSeverity, *child2Assignment.Severity)
		assert.Equal(t, rootSeverity, *child2Assignment.ResidualSeverity)

		// Verify all grandchildren got the threat assignment
		grandchild1Assignments, err := ListThreatAssignmentsByComponentID(grandchild1.ID)
		require.NoError(t, err)
		require.Len(t, grandchild1Assignments, 1)
		grandchild1Assignment := grandchild1Assignments[0]
		require.NotNil(t, grandchild1Assignment.Severity)
		require.NotNil(t, grandchild1Assignment.ResidualSeverity)
		assert.Equal(t, rootSeverity, *grandchild1Assignment.Severity)
		assert.Equal(t, rootSeverity, *grandchild1Assignment.ResidualSeverity)

		grandchild2Assignments, err := ListThreatAssignmentsByComponentID(grandchild2.ID)
		require.NoError(t, err)
		require.Len(t, grandchild2Assignments, 1)
		grandchild2Assignment := grandchild2Assignments[0]
		require.NotNil(t, grandchild2Assignment.Severity)
		require.NotNil(t, grandchild2Assignment.ResidualSeverity)
		assert.Equal(t, rootSeverity, *grandchild2Assignment.Severity)
		assert.Equal(t, rootSeverity, *grandchild2Assignment.ResidualSeverity)

		grandchild3Assignments, err := ListThreatAssignmentsByComponentID(grandchild3.ID)
		require.NoError(t, err)
		require.Len(t, grandchild3Assignments, 1)
		grandchild3Assignment := grandchild3Assignments[0]
		require.NotNil(t, grandchild3Assignment.Severity)
		require.NotNil(t, grandchild3Assignment.ResidualSeverity)
		assert.Equal(t, rootSeverity, *grandchild3Assignment.Severity)
		assert.Equal(t, rootSeverity, *grandchild3Assignment.ResidualSeverity)

		// Verify great grandchild got the threat assignment
		greatGrandchild1Assignments, err := ListThreatAssignmentsByComponentID(greatGrandchild1.ID)
		require.NoError(t, err)
		require.Len(t, greatGrandchild1Assignments, 1)
		greatGrandchild1Assignment := greatGrandchild1Assignments[0]
		require.NotNil(t, greatGrandchild1Assignment.Severity)
		require.NotNil(t, greatGrandchild1Assignment.ResidualSeverity)
		assert.Equal(t, rootSeverity, *greatGrandchild1Assignment.Severity)
		assert.Equal(t, rootSeverity, *greatGrandchild1Assignment.ResidualSeverity)

		// Now test that updating root severity propagates throughout the entire tree
		updatedRootSeverity := models.ThreatSeverityLow
		err = SetThreatAssignmentSeverity(rootAssignment.ID, &updatedRootSeverity)
		require.NoError(t, err)

		// Verify all components now have the updated severity
		componentsToCheck := []uuid.UUID{child1.ID, child2.ID, grandchild1.ID, grandchild2.ID, grandchild3.ID, greatGrandchild1.ID}
		
		for _, componentID := range componentsToCheck {
			assignments, err := ListThreatAssignmentsByComponentID(componentID)
			require.NoError(t, err)
			require.Len(t, assignments, 1)
			require.NotNil(t, assignments[0].Severity)
			assert.Equal(t, updatedRootSeverity, *assignments[0].Severity, "Component %s should have updated severity", componentID)
		}
	})

	t.Run("MultiLevelInheritance", func(t *testing.T) {
		// Create test components: Grandparent -> Parent -> Child
		grandparent, err := CreateComponent("Grandparent Component", "Test grandparent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		parent, err := CreateComponent("Parent Component", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child, err := CreateComponent("Child Component", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat
		threat, err := CreateThreat("Multilevel Test Threat", "Test threat for multilevel inheritance")
		require.NoError(t, err)

		// Create inherits-threats relationships
		_, err = CreateInheritsThreatsRelationship(parent.ID, grandparent.ID)
		require.NoError(t, err)
		_, err = CreateInheritsThreatsRelationship(child.ID, parent.ID)
		require.NoError(t, err)

		// Assign threat to grandparent with severity
		grandparentAssignment, err := AssignThreatToComponent(grandparent.ID, threat.ID)
		require.NoError(t, err)

		severity := models.ThreatSeverityHigh
		err = SetThreatAssignmentSeverity(grandparentAssignment.ID, &severity)
		require.NoError(t, err)

		// Verify parent got the threat assignment with grandparent's residual severity
		parentAssignments, err := ListThreatAssignmentsByComponentID(parent.ID)
		require.NoError(t, err)
		require.Len(t, parentAssignments, 1)

		parentAssignment := parentAssignments[0]
		assert.Equal(t, threat.ID, parentAssignment.ThreatID)
		require.NotNil(t, parentAssignment.Severity)
		assert.Equal(t, severity, *parentAssignment.Severity)
		// Parent's ResidualSeverity should also be set automatically
		require.NotNil(t, parentAssignment.ResidualSeverity)
		assert.Equal(t, severity, *parentAssignment.ResidualSeverity)

		// Verify child got the threat assignment (from parent's residual severity)
		childAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		require.Len(t, childAssignments, 1)

		childAssignment := childAssignments[0]
		assert.Equal(t, threat.ID, childAssignment.ThreatID)
		require.NotNil(t, childAssignment.Severity)
		// Child gets parent's ResidualSeverity as its Severity
		assert.Equal(t, *parentAssignment.ResidualSeverity, *childAssignment.Severity)
	})

	t.Run("NoInheritanceWithoutRelationship", func(t *testing.T) {
		// Create test components without inherits-threats relationship
		parent, err := CreateComponent("No Inherit Parent", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child, err := CreateComponent("No Inherit Child", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat
		threat, err := CreateThreat("No Inherit Threat", "Test threat for no inheritance")
		require.NoError(t, err)

		// Assign threat to parent with severity (but no relationship)
		parentAssignment, err := AssignThreatToComponent(parent.ID, threat.ID)
		require.NoError(t, err)

		severity := models.ThreatSeverityCritical
		err = SetThreatAssignmentSeverity(parentAssignment.ID, &severity)
		require.NoError(t, err)

		// Verify child did NOT get the threat assignment
		childAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		assert.Len(t, childAssignments, 0)
	})

	t.Run("InheritanceWithSeverityAlwaysPropagates", func(t *testing.T) {
		// Test that severity propagates even without explicit ResidualSeverity setting
		parent, err := CreateComponent("Auto Propagate Parent", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child, err := CreateComponent("Auto Propagate Child", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat
		threat, err := CreateThreat("Auto Propagate Threat", "Test threat for auto propagation")
		require.NoError(t, err)

		// Create inherits-threats relationship
		_, err = CreateInheritsThreatsRelationship(child.ID, parent.ID)
		require.NoError(t, err)

		// Assign threat to parent and set severity (ResidualSeverity will be auto-set)
		parentAssignment, err := AssignThreatToComponent(parent.ID, threat.ID)
		require.NoError(t, err)

		severity := models.ThreatSeverityLow
		err = SetThreatAssignmentSeverity(parentAssignment.ID, &severity)
		require.NoError(t, err)

		// Verify child automatically got the threat assignment
		childAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		require.Len(t, childAssignments, 1)

		childAssignment := childAssignments[0]
		require.NotNil(t, childAssignment.Severity)
		assert.Equal(t, severity, *childAssignment.Severity)
	})

	t.Run("UpdateExistingInheritedAssignment", func(t *testing.T) {
		// Create test components
		parent, err := CreateComponent("Update Parent", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child, err := CreateComponent("Update Child", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat
		threat, err := CreateThreat("Update Threat", "Test threat for updates")
		require.NoError(t, err)

		// Create inherits-threats relationship
		_, err = CreateInheritsThreatsRelationship(child.ID, parent.ID)
		require.NoError(t, err)

		// Assign threat to parent with initial severity
		parentAssignment, err := AssignThreatToComponent(parent.ID, threat.ID)
		require.NoError(t, err)

		initialSeverity := models.ThreatSeverityLow
		err = SetThreatAssignmentSeverity(parentAssignment.ID, &initialSeverity)
		require.NoError(t, err)

		// Verify child got initial severity
		childAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		require.Len(t, childAssignments, 1)
		require.NotNil(t, childAssignments[0].Severity)
		assert.Equal(t, initialSeverity, *childAssignments[0].Severity)

		// Update parent's severity (ResidualSeverity should also be updated since it matched)
		updatedSeverity := models.ThreatSeverityCritical
		err = SetThreatAssignmentSeverity(parentAssignment.ID, &updatedSeverity)
		require.NoError(t, err)

		// Verify parent's ResidualSeverity was updated too
		updatedParent, err := GetThreatAssignmentById(parentAssignment.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedParent.ResidualSeverity)
		assert.Equal(t, updatedSeverity, *updatedParent.ResidualSeverity)

		// Verify child's severity was updated
		updatedChildAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		require.Len(t, updatedChildAssignments, 1)
		require.NotNil(t, updatedChildAssignments[0].Severity)
		assert.Equal(t, updatedSeverity, *updatedChildAssignments[0].Severity)
	})

	t.Run("CustomResidualSeverityOverride", func(t *testing.T) {
		// Test that custom ResidualSeverity doesn't get overwritten by Severity updates
		parent, err := CreateComponent("Custom Residual Parent", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child, err := CreateComponent("Custom Residual Child", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create threat
		threat, err := CreateThreat("Custom Residual Threat", "Test threat for custom residual")
		require.NoError(t, err)

		// Create inherits-threats relationship
		_, err = CreateInheritsThreatsRelationship(child.ID, parent.ID)
		require.NoError(t, err)

		// Assign threat to parent with severity
		parentAssignment, err := AssignThreatToComponent(parent.ID, threat.ID)
		require.NoError(t, err)

		severity := models.ThreatSeverityCritical
		err = SetThreatAssignmentSeverity(parentAssignment.ID, &severity)
		require.NoError(t, err)

		// Explicitly set a different ResidualSeverity
		customResidualSeverity := models.ThreatSeverityLow
		err = SetThreatAssignmentResidualSeverity(parentAssignment.ID, &customResidualSeverity)
		require.NoError(t, err)

		// Update parent's severity - ResidualSeverity should NOT change
		updatedSeverity := models.ThreatSeverityHigh
		err = SetThreatAssignmentSeverity(parentAssignment.ID, &updatedSeverity)
		require.NoError(t, err)

		// Verify parent still has custom ResidualSeverity
		updatedParent, err := GetThreatAssignmentById(parentAssignment.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedParent.Severity)
		require.NotNil(t, updatedParent.ResidualSeverity)
		assert.Equal(t, updatedSeverity, *updatedParent.Severity)
		assert.Equal(t, customResidualSeverity, *updatedParent.ResidualSeverity)

		// Verify child got the custom ResidualSeverity, not the Severity
		childAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		require.Len(t, childAssignments, 1)
		require.NotNil(t, childAssignments[0].Severity)
		assert.Equal(t, customResidualSeverity, *childAssignments[0].Severity)
	})

	t.Run("MultipleThreatsInheritance", func(t *testing.T) {
		// Create test components
		parent, err := CreateComponent("Multi Threat Parent", "Test parent component", models.ComponentTypeProduct)
		require.NoError(t, err)

		child, err := CreateComponent("Multi Threat Child", "Test child component", models.ComponentTypeInstance)
		require.NoError(t, err)

		// Create multiple threats
		threat1, err := CreateThreat("Multi Threat 1", "First threat for multi inheritance")
		require.NoError(t, err)

		threat2, err := CreateThreat("Multi Threat 2", "Second threat for multi inheritance")
		require.NoError(t, err)

		// Create inherits-threats relationship
		_, err = CreateInheritsThreatsRelationship(child.ID, parent.ID)
		require.NoError(t, err)

		// Assign both threats to parent with different severities
		parentAssignment1, err := AssignThreatToComponent(parent.ID, threat1.ID)
		require.NoError(t, err)
		severity1 := models.ThreatSeverityHigh
		err = SetThreatAssignmentSeverity(parentAssignment1.ID, &severity1)
		require.NoError(t, err)

		parentAssignment2, err := AssignThreatToComponent(parent.ID, threat2.ID)
		require.NoError(t, err)
		severity2 := models.ThreatSeverityMedium
		err = SetThreatAssignmentSeverity(parentAssignment2.ID, &severity2)
		require.NoError(t, err)

		// Verify child got both threat assignments with correct severities
		childAssignments, err := ListThreatAssignmentsByComponentID(child.ID)
		require.NoError(t, err)
		require.Len(t, childAssignments, 2)

		// Create map for easy lookup
		childAssignmentMap := make(map[uuid.UUID]*models.ThreatAssignment)
		for i := range childAssignments {
			childAssignmentMap[childAssignments[i].ThreatID] = &childAssignments[i]
		}

		// Verify first threat assignment
		childAssignment1, exists := childAssignmentMap[threat1.ID]
		require.True(t, exists)
		require.NotNil(t, childAssignment1.Severity)
		assert.Equal(t, severity1, *childAssignment1.Severity)

		// Verify second threat assignment
		childAssignment2, exists := childAssignmentMap[threat2.ID]
		require.True(t, exists)
		require.NotNil(t, childAssignment2.Severity)
		assert.Equal(t, severity2, *childAssignment2.Severity)
	})

	t.Run("AutoResidualSeverityOnCreation", func(t *testing.T) {
		// Test that ResidualSeverity is automatically set when creating assignment with Severity
		component, err := CreateComponent("Auto Residual Component", "Test component", models.ComponentTypeProduct)
		require.NoError(t, err)

		threat, err := CreateThreat("Auto Residual Threat", "Test threat")
		require.NoError(t, err)

		// Create assignment
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)

		// Set severity
		severity := models.ThreatSeverityMedium
		err = SetThreatAssignmentSeverity(assignment.ID, &severity)
		require.NoError(t, err)

		// Verify ResidualSeverity was automatically set to match Severity
		updated, err := GetThreatAssignmentById(assignment.ID)
		require.NoError(t, err)
		require.NotNil(t, updated.Severity)
		require.NotNil(t, updated.ResidualSeverity)
		assert.Equal(t, severity, *updated.Severity)
		assert.Equal(t, severity, *updated.ResidualSeverity)
	})
}