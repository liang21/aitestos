// Package testplan defines TestResult entity
package testplan

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ResultStatus is a value object representing the result of a test case execution
type ResultStatus string

const (
	// ResultPass means the test passed
	ResultPass ResultStatus = "pass"
	// ResultFail means the test failed
	ResultFail ResultStatus = "fail"
	// ResultBlock means the test is blocked
	ResultBlock ResultStatus = "block"
	// ResultSkip means the test was skipped
	ResultSkip ResultStatus = "skip"
)

// ParseResultStatus validates and creates a ResultStatus
func ParseResultStatus(s string) (ResultStatus, error) {
	switch ResultStatus(s) {
	case ResultPass, ResultFail, ResultBlock, ResultSkip:
		return ResultStatus(s), nil
	default:
		return "", errors.New("invalid result status")
	}
}

// String returns the string representation
func (s ResultStatus) String() string {
	return string(s)
}

// TestResult represents the execution result of a test case
type TestResult struct {
	id         uuid.UUID
	planID     uuid.UUID
	caseID     uuid.UUID
	executedBy uuid.UUID
	status     ResultStatus
	note       string
	executedAt time.Time
	updatedAt  time.Time
}

// NewTestResult creates a new test result
func NewTestResult(planID, caseID, userID uuid.UUID, status ResultStatus, note string) (*TestResult, error) {
	if planID == uuid.Nil {
		return nil, errors.New("plan ID cannot be nil")
	}
	if caseID == uuid.Nil {
		return nil, errors.New("case ID cannot be nil")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be nil")
	}

	now := time.Now()
	return &TestResult{
		id:         uuid.New(),
		planID:     planID,
		caseID:     caseID,
		executedBy: userID,
		status:     status,
		note:       note,
		executedAt: now,
		updatedAt:  now,
	}, nil
}

// ID returns the test result's unique identifier
func (r *TestResult) ID() uuid.UUID {
	return r.id
}

// PlanID returns the associated plan's ID
func (r *TestResult) PlanID() uuid.UUID {
	return r.planID
}

// CaseID returns the associated test case's ID
func (r *TestResult) CaseID() uuid.UUID {
	return r.caseID
}

// ExecutedBy returns the user who executed this result
func (r *TestResult) ExecutedBy() uuid.UUID {
	return r.executedBy
}

// Status returns the result status
func (r *TestResult) Status() ResultStatus {
	return r.status
}

// Note returns the result note
func (r *TestResult) Note() string {
	return r.note
}

// ExecutedAt returns the execution timestamp
func (r *TestResult) ExecutedAt() time.Time {
	return r.executedAt
}

// UpdatedAt returns the last update timestamp
func (r *TestResult) UpdatedAt() time.Time {
	return r.updatedAt
}

// UpdateNote updates the result note
func (r *TestResult) UpdateNote(note string) {
	r.note = note
	r.updatedAt = time.Now()
}

// ReconstructResult reconstructs a TestResult from stored data
func ReconstructResult(
	id uuid.UUID,
	planID uuid.UUID,
	caseID uuid.UUID,
	executedBy uuid.UUID,
	status ResultStatus,
	note string,
	executedAt time.Time,
	updatedAt time.Time,
) *TestResult {
	return &TestResult{
		id:         id,
		planID:     planID,
		caseID:     caseID,
		executedBy: executedBy,
		status:     status,
		note:       note,
		executedAt: executedAt,
		updatedAt:  updatedAt,
	}
}
