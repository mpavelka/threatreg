package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreatSeverity_ToInt(t *testing.T) {
	tests := []struct {
		name     string
		severity ThreatSeverity
		expected int
	}{
		{"None", ThreatSeverityNone, 0},
		{"Low", ThreatSeverityLow, 1},
		{"Medium", ThreatSeverityMedium, 2},
		{"High", ThreatSeverityHigh, 3},
		{"Critical", ThreatSeverityCritical, 4},
		{"Invalid", ThreatSeverity("invalid"), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.severity.ToInt()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestThreatSeverityFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    ThreatSeverity
		expectError bool
	}{
		{"ValidNone", "none", ThreatSeverityNone, false},
		{"ValidLow", "low", ThreatSeverityLow, false},
		{"ValidMedium", "medium", ThreatSeverityMedium, false},
		{"ValidHigh", "high", ThreatSeverityHigh, false},
		{"ValidCritical", "critical", ThreatSeverityCritical, false},
		{"ValidCaseInsensitive", "HIGH", ThreatSeverityHigh, false},
		{"ValidWithSpaces", " medium ", ThreatSeverityMedium, false},
		{"Invalid", "invalid", "", true},
		{"Empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ThreatSeverityFromString(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestThreatSeverityFromInt(t *testing.T) {
	tests := []struct {
		name        string
		input       int
		expected    ThreatSeverity
		expectError bool
	}{
		{"Valid0", 0, ThreatSeverityNone, false},
		{"Valid1", 1, ThreatSeverityLow, false},
		{"Valid2", 2, ThreatSeverityMedium, false},
		{"Valid3", 3, ThreatSeverityHigh, false},
		{"Valid4", 4, ThreatSeverityCritical, false},
		{"Invalid-1", -1, "", true},
		{"Invalid5", 5, "", true},
		{"Invalid100", 100, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ThreatSeverityFromInt(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestThreatSeverity_String(t *testing.T) {
	tests := []struct {
		name     string
		severity ThreatSeverity
		expected string
	}{
		{"None", ThreatSeverityNone, "none"},
		{"Low", ThreatSeverityLow, "low"},
		{"Medium", ThreatSeverityMedium, "medium"},
		{"High", ThreatSeverityHigh, "high"},
		{"Critical", ThreatSeverityCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.severity.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestThreatSeverity_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		severity ThreatSeverity
		expected bool
	}{
		{"ValidNone", ThreatSeverityNone, true},
		{"ValidLow", ThreatSeverityLow, true},
		{"ValidMedium", ThreatSeverityMedium, true},
		{"ValidHigh", ThreatSeverityHigh, true},
		{"ValidCritical", ThreatSeverityCritical, true},
		{"Invalid", ThreatSeverity("invalid"), false},
		{"Empty", ThreatSeverity(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.severity.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestThreatSeverity_Comparisons(t *testing.T) {
	low := ThreatSeverityLow
	medium := ThreatSeverityMedium
	high := ThreatSeverityHigh
	
	t.Run("Less", func(t *testing.T) {
		assert.True(t, low.Less(medium))
		assert.True(t, medium.Less(high))
		assert.False(t, high.Less(medium))
		assert.False(t, medium.Less(low))
		assert.False(t, medium.Less(medium))
	})
	
	t.Run("Greater", func(t *testing.T) {
		assert.True(t, high.Greater(medium))
		assert.True(t, medium.Greater(low))
		assert.False(t, low.Greater(medium))
		assert.False(t, medium.Greater(high))
		assert.False(t, medium.Greater(medium))
	})
	
	t.Run("Equal", func(t *testing.T) {
		assert.True(t, medium.Equal(medium))
		assert.True(t, low.Equal(low))
		assert.False(t, low.Equal(medium))
		assert.False(t, medium.Equal(high))
	})
	
	t.Run("LessOrEqual", func(t *testing.T) {
		assert.True(t, low.LessOrEqual(medium))
		assert.True(t, medium.LessOrEqual(medium))
		assert.True(t, medium.LessOrEqual(high))
		assert.False(t, high.LessOrEqual(medium))
	})
	
	t.Run("GreaterOrEqual", func(t *testing.T) {
		assert.True(t, high.GreaterOrEqual(medium))
		assert.True(t, medium.GreaterOrEqual(medium))
		assert.True(t, medium.GreaterOrEqual(low))
		assert.False(t, low.GreaterOrEqual(medium))
	})
}

func TestThreatSeverity_ComparisonOrder(t *testing.T) {
	severities := []ThreatSeverity{
		ThreatSeverityNone,
		ThreatSeverityLow,
		ThreatSeverityMedium,
		ThreatSeverityHigh,
		ThreatSeverityCritical,
	}

	// Test that each severity is less than all subsequent severities
	for i, current := range severities {
		for j, other := range severities {
			if i < j {
				assert.True(t, current.Less(other), 
					"Expected %s to be less than %s", current, other)
				assert.False(t, current.Greater(other),
					"Expected %s not to be greater than %s", current, other)
			} else if i > j {
				assert.True(t, current.Greater(other),
					"Expected %s to be greater than %s", current, other)
				assert.False(t, current.Less(other),
					"Expected %s not to be less than %s", current, other)
			} else {
				assert.True(t, current.Equal(other),
					"Expected %s to be equal to %s", current, other)
			}
		}
	}
}

func TestGetAllValidThreatSeverities(t *testing.T) {
	validSeverities := GetAllValidThreatSeverities()
	
	expected := []ThreatSeverity{
		ThreatSeverityNone,
		ThreatSeverityLow,
		ThreatSeverityMedium,
		ThreatSeverityHigh,
		ThreatSeverityCritical,
	}
	
	assert.Equal(t, expected, validSeverities)
	assert.Len(t, validSeverities, 5)
	
	// Verify all returned severities are valid
	for _, severity := range validSeverities {
		assert.True(t, severity.IsValid(), "All returned severities should be valid")
	}
}