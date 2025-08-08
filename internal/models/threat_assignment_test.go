package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestThreatAssignment_ValidateAssignment(t *testing.T) {
	t.Run("ValidWithNilSeverities", func(t *testing.T) {
		assignment := &ThreatAssignment{
			ComponentID:      uuid.New(),
			ThreatID:         uuid.New(),
			Severity:         nil,
			ResidualSeverity: nil,
		}
		
		err := assignment.validateAssignment()
		assert.NoError(t, err)
	})

	t.Run("ValidWithValidSeverities", func(t *testing.T) {
		severity := ThreatSeverityHigh
		residualSeverity := ThreatSeverityLow
		
		assignment := &ThreatAssignment{
			ComponentID:      uuid.New(),
			ThreatID:         uuid.New(),
			Severity:         &severity,
			ResidualSeverity: &residualSeverity,
		}
		
		err := assignment.validateAssignment()
		assert.NoError(t, err)
	})

	t.Run("InvalidNilComponentID", func(t *testing.T) {
		assignment := &ThreatAssignment{
			ComponentID: uuid.Nil,
			ThreatID:    uuid.New(),
		}
		
		err := assignment.validateAssignment()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ComponentID")
	})

	t.Run("InvalidSeverity", func(t *testing.T) {
		invalidSeverity := ThreatSeverity("invalid")
		
		assignment := &ThreatAssignment{
			ComponentID: uuid.New(),
			ThreatID:    uuid.New(),
			Severity:    &invalidSeverity,
		}
		
		err := assignment.validateAssignment()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid severity")
	})

	t.Run("InvalidResidualSeverity", func(t *testing.T) {
		invalidResidualSeverity := ThreatSeverity("invalid")
		
		assignment := &ThreatAssignment{
			ComponentID:      uuid.New(),
			ThreatID:         uuid.New(),
			ResidualSeverity: &invalidResidualSeverity,
		}
		
		err := assignment.validateAssignment()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid residual severity")
	})
}