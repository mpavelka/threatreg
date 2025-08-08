package service

import (
	"testing"
	"threatreg/internal/models"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetThreatAssignmentSeverity_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("SetValidSeverity", func(t *testing.T) {
		// Create test data
		component, err := CreateComponent("Severity Test Component", "Test component for severity", models.ComponentTypeProduct)
		require.NoError(t, err)
		
		threat, err := CreateThreat("Severity Test Threat", "Test threat for severity")
		require.NoError(t, err)
		
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)
		assert.Nil(t, assignment.Severity)

		// Set severity
		severity := models.ThreatSeverityHigh
		err = SetThreatAssignmentSeverity(assignment.ID, &severity)
		assert.NoError(t, err)

		// Verify severity was set
		updated, err := GetThreatAssignmentById(assignment.ID)
		require.NoError(t, err)
		assert.NotNil(t, updated.Severity)
		assert.Equal(t, severity, *updated.Severity)
	})

	t.Run("ClearSeverity", func(t *testing.T) {
		// Create test data
		component, err := CreateComponent("Severity Clear Test Component", "Test component for severity clearing", models.ComponentTypeProduct)
		require.NoError(t, err)
		
		threat, err := CreateThreat("Severity Clear Test Threat", "Test threat for severity clearing")
		require.NoError(t, err)
		
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)

		// Set initial severity
		initialSeverity := models.ThreatSeverityCritical
		err = SetThreatAssignmentSeverity(assignment.ID, &initialSeverity)
		require.NoError(t, err)

		// Clear severity
		err = SetThreatAssignmentSeverity(assignment.ID, nil)
		assert.NoError(t, err)

		// Verify severity was cleared
		updated, err := GetThreatAssignmentById(assignment.ID)
		require.NoError(t, err)
		assert.Nil(t, updated.Severity)
	})

	t.Run("SetInvalidSeverity", func(t *testing.T) {
		// Create test data
		component, err := CreateComponent("Invalid Severity Test Component", "Test component for invalid severity", models.ComponentTypeProduct)
		require.NoError(t, err)
		
		threat, err := CreateThreat("Invalid Severity Test Threat", "Test threat for invalid severity")
		require.NoError(t, err)
		
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)

		// Try to set invalid severity
		invalidSeverity := models.ThreatSeverity("invalid")
		err = SetThreatAssignmentSeverity(assignment.ID, &invalidSeverity)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid severity")
	})

	t.Run("SetSeverityNonExistentAssignment", func(t *testing.T) {
		nonExistentID := uuid.New()
		severity := models.ThreatSeverityLow
		
		err := SetThreatAssignmentSeverity(nonExistentID, &severity)
		// Service should not return an error for non-existent assignments
		// as GORM update with no matches doesn't error
		assert.NoError(t, err)
	})

	t.Run("SetAllSeverityLevels", func(t *testing.T) {
		severities := []models.ThreatSeverity{
			models.ThreatSeverityNone,
			models.ThreatSeverityLow,
			models.ThreatSeverityMedium,
			models.ThreatSeverityHigh,
			models.ThreatSeverityCritical,
		}

		for i, severity := range severities {
			t.Run(string(severity), func(t *testing.T) {
				// Create unique test data for each severity level
				component, err := CreateComponent(
					"Severity Level Test Component "+string(severity),
					"Test component for severity level "+string(severity),
					models.ComponentTypeProduct,
				)
				require.NoError(t, err)
				
				threat, err := CreateThreat(
					"Severity Level Test Threat "+string(severity),
					"Test threat for severity level "+string(severity),
				)
				require.NoError(t, err)
				
				assignment, err := AssignThreatToComponent(component.ID, threat.ID)
				require.NoError(t, err)

				// Set severity
				err = SetThreatAssignmentSeverity(assignment.ID, &severities[i])
				assert.NoError(t, err)

				// Verify severity was set correctly
				updated, err := GetThreatAssignmentById(assignment.ID)
				require.NoError(t, err)
				assert.NotNil(t, updated.Severity)
				assert.Equal(t, severities[i], *updated.Severity)
			})
		}
	})
}

func TestSetThreatAssignmentResidualSeverity_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	t.Run("SetValidResidualSeverity", func(t *testing.T) {
		// Create test data
		component, err := CreateComponent("Residual Severity Test Component", "Test component for residual severity", models.ComponentTypeProduct)
		require.NoError(t, err)
		
		threat, err := CreateThreat("Residual Severity Test Threat", "Test threat for residual severity")
		require.NoError(t, err)
		
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)
		assert.Nil(t, assignment.ResidualSeverity)

		// Set residual severity
		residualSeverity := models.ThreatSeverityLow
		err = SetThreatAssignmentResidualSeverity(assignment.ID, &residualSeverity)
		assert.NoError(t, err)

		// Verify residual severity was set
		updated, err := GetThreatAssignmentById(assignment.ID)
		require.NoError(t, err)
		assert.NotNil(t, updated.ResidualSeverity)
		assert.Equal(t, residualSeverity, *updated.ResidualSeverity)
	})

	t.Run("ClearResidualSeverity", func(t *testing.T) {
		// Create test data
		component, err := CreateComponent("Residual Severity Clear Test Component", "Test component for residual severity clearing", models.ComponentTypeProduct)
		require.NoError(t, err)
		
		threat, err := CreateThreat("Residual Severity Clear Test Threat", "Test threat for residual severity clearing")
		require.NoError(t, err)
		
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)

		// Set initial residual severity
		initialResidualSeverity := models.ThreatSeverityHigh
		err = SetThreatAssignmentResidualSeverity(assignment.ID, &initialResidualSeverity)
		require.NoError(t, err)

		// Clear residual severity
		err = SetThreatAssignmentResidualSeverity(assignment.ID, nil)
		assert.NoError(t, err)

		// Verify residual severity was cleared
		updated, err := GetThreatAssignmentById(assignment.ID)
		require.NoError(t, err)
		assert.Nil(t, updated.ResidualSeverity)
	})

	t.Run("SetInvalidResidualSeverity", func(t *testing.T) {
		// Create test data
		component, err := CreateComponent("Invalid Residual Severity Test Component", "Test component for invalid residual severity", models.ComponentTypeProduct)
		require.NoError(t, err)
		
		threat, err := CreateThreat("Invalid Residual Severity Test Threat", "Test threat for invalid residual severity")
		require.NoError(t, err)
		
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)

		// Try to set invalid residual severity
		invalidResidualSeverity := models.ThreatSeverity("invalid")
		err = SetThreatAssignmentResidualSeverity(assignment.ID, &invalidResidualSeverity)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid residual severity")
	})

	t.Run("SetBothSeverities", func(t *testing.T) {
		// Create test data
		component, err := CreateComponent("Both Severities Test Component", "Test component for both severities", models.ComponentTypeProduct)
		require.NoError(t, err)
		
		threat, err := CreateThreat("Both Severities Test Threat", "Test threat for both severities")
		require.NoError(t, err)
		
		assignment, err := AssignThreatToComponent(component.ID, threat.ID)
		require.NoError(t, err)

		// Set both severities
		severity := models.ThreatSeverityCritical
		residualSeverity := models.ThreatSeverityMedium

		err = SetThreatAssignmentSeverity(assignment.ID, &severity)
		require.NoError(t, err)
		
		err = SetThreatAssignmentResidualSeverity(assignment.ID, &residualSeverity)
		require.NoError(t, err)

		// Verify both severities were set correctly
		updated, err := GetThreatAssignmentById(assignment.ID)
		require.NoError(t, err)
		
		assert.NotNil(t, updated.Severity)
		assert.Equal(t, severity, *updated.Severity)
		
		assert.NotNil(t, updated.ResidualSeverity)
		assert.Equal(t, residualSeverity, *updated.ResidualSeverity)
	})
}

func TestSeverityComparisons_Integration(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create test data
	component1, err := CreateComponent("Comparison Test Component 1", "Test component 1 for comparisons", models.ComponentTypeProduct)
	require.NoError(t, err)
	
	component2, err := CreateComponent("Comparison Test Component 2", "Test component 2 for comparisons", models.ComponentTypeProduct)
	require.NoError(t, err)
	
	threat, err := CreateThreat("Comparison Test Threat", "Test threat for comparisons")
	require.NoError(t, err)
	
	assignment1, err := AssignThreatToComponent(component1.ID, threat.ID)
	require.NoError(t, err)
	
	assignment2, err := AssignThreatToComponent(component2.ID, threat.ID)
	require.NoError(t, err)

	// Set different severities
	highSeverity := models.ThreatSeverityHigh
	lowSeverity := models.ThreatSeverityLow

	err = SetThreatAssignmentSeverity(assignment1.ID, &highSeverity)
	require.NoError(t, err)
	
	err = SetThreatAssignmentSeverity(assignment2.ID, &lowSeverity)
	require.NoError(t, err)

	// Retrieve assignments
	updated1, err := GetThreatAssignmentById(assignment1.ID)
	require.NoError(t, err)
	
	updated2, err := GetThreatAssignmentById(assignment2.ID)
	require.NoError(t, err)

	// Test comparisons
	assert.True(t, updated1.Severity.Greater(*updated2.Severity))
	assert.False(t, updated1.Severity.Less(*updated2.Severity))
	assert.True(t, updated2.Severity.Less(*updated1.Severity))
	assert.False(t, updated2.Severity.Greater(*updated1.Severity))
	assert.False(t, updated1.Severity.Equal(*updated2.Severity))
}