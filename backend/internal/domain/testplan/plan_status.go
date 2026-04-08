// Package testplan defines PlanStatus value object
package testplan

import "errors"

// PlanStatus is a value object representing the status of a test plan
type PlanStatus string

const (
	// StatusDraft means the plan is in draft state
	StatusDraft PlanStatus = "draft"
	// StatusActive means the plan is active
	StatusActive PlanStatus = "active"
	// StatusCompleted means the plan is completed
	StatusCompleted PlanStatus = "completed"
	// StatusArchived means the plan is archived
	StatusArchived PlanStatus = "archived"
)

// ParsePlanStatus validates and creates a PlanStatus
func ParsePlanStatus(s string) (PlanStatus, error) {
	switch PlanStatus(s) {
	case StatusDraft, StatusActive, StatusCompleted, StatusArchived:
		return PlanStatus(s), nil
	default:
		return "", errors.New("invalid plan status")
	}
}

// String returns the string representation
func (s PlanStatus) String() string {
	return string(s)
}

// CanTransitionTo checks if transition to another status is valid
func (s PlanStatus) CanTransitionTo(to PlanStatus) bool {
	transitions := map[PlanStatus]map[PlanStatus]bool{
		StatusDraft:     {StatusActive: true, StatusArchived: true},
		StatusActive:    {StatusCompleted: true, StatusArchived: true},
		StatusCompleted: {StatusArchived: true},
		StatusArchived:  {},
	}

	allowed, exists := transitions[s][to]
	return exists && allowed
}

// IsFinal returns true if the status is a final state
func (s PlanStatus) IsFinal() bool {
	return s == StatusCompleted || s == StatusArchived
}
