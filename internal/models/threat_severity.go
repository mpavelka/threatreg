package models

import (
	"fmt"
	"strings"
)

// ThreatSeverity represents the severity level of a threat
type ThreatSeverity string

const (
	ThreatSeverityNone     ThreatSeverity = "none"
	ThreatSeverityLow      ThreatSeverity = "low"
	ThreatSeverityMedium   ThreatSeverity = "medium"
	ThreatSeverityHigh     ThreatSeverity = "high"
	ThreatSeverityCritical ThreatSeverity = "critical"
)

// severityToInt maps threat severity to integer values for comparison
var severityToInt = map[ThreatSeverity]int{
	ThreatSeverityNone:     0,
	ThreatSeverityLow:      1,
	ThreatSeverityMedium:   2,
	ThreatSeverityHigh:     3,
	ThreatSeverityCritical: 4,
}

// intToSeverity maps integer values back to threat severity
var intToSeverity = map[int]ThreatSeverity{
	0: ThreatSeverityNone,
	1: ThreatSeverityLow,
	2: ThreatSeverityMedium,
	3: ThreatSeverityHigh,
	4: ThreatSeverityCritical,
}

// ToInt returns the integer representation of the severity level for mathematical comparisons
func (ts ThreatSeverity) ToInt() int {
	if val, exists := severityToInt[ts]; exists {
		return val
	}
	return -1 // Invalid severity
}

// FromString creates a ThreatSeverity from a string value
func ThreatSeverityFromString(s string) (ThreatSeverity, error) {
	severity := ThreatSeverity(strings.ToLower(strings.TrimSpace(s)))
	if !severity.IsValid() {
		return "", fmt.Errorf("invalid threat severity: %s. Valid values are: none, low, medium, high, critical", s)
	}
	return severity, nil
}

// FromInt creates a ThreatSeverity from an integer value
func ThreatSeverityFromInt(i int) (ThreatSeverity, error) {
	if severity, exists := intToSeverity[i]; exists {
		return severity, nil
	}
	return "", fmt.Errorf("invalid threat severity integer: %d. Valid values are: 0-4", i)
}

// String returns the string representation of the severity level
func (ts ThreatSeverity) String() string {
	return string(ts)
}

// IsValid checks if the severity value is one of the allowed values
func (ts ThreatSeverity) IsValid() bool {
	_, exists := severityToInt[ts]
	return exists
}

// Less returns true if this severity is lower than the other severity
func (ts ThreatSeverity) Less(other ThreatSeverity) bool {
	return ts.ToInt() < other.ToInt()
}

// Greater returns true if this severity is higher than the other severity
func (ts ThreatSeverity) Greater(other ThreatSeverity) bool {
	return ts.ToInt() > other.ToInt()
}

// Equal returns true if this severity is equal to the other severity
func (ts ThreatSeverity) Equal(other ThreatSeverity) bool {
	return ts.ToInt() == other.ToInt()
}

// LessOrEqual returns true if this severity is less than or equal to the other severity
func (ts ThreatSeverity) LessOrEqual(other ThreatSeverity) bool {
	return ts.ToInt() <= other.ToInt()
}

// GreaterOrEqual returns true if this severity is greater than or equal to the other severity
func (ts ThreatSeverity) GreaterOrEqual(other ThreatSeverity) bool {
	return ts.ToInt() >= other.ToInt()
}

// GetAllValidSeverities returns all valid severity values
func GetAllValidThreatSeverities() []ThreatSeverity {
	return []ThreatSeverity{
		ThreatSeverityNone,
		ThreatSeverityLow,
		ThreatSeverityMedium,
		ThreatSeverityHigh,
		ThreatSeverityCritical,
	}
}