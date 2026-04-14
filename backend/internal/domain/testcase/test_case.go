// Package testcase defines TestCase aggregate
package testcase

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Preconditions represents test case preconditions
type Preconditions []string

// Steps represents test case execution steps
type Steps []string

// ExpectedResult represents expected test result
type ExpectedResult map[string]any

// CaseType represents the type of test case
type CaseType string

const (
	// CaseTypeFunctionality represents functional testing
	CaseTypeFunctionality CaseType = "functionality"
	// CaseTypePerformance represents performance testing
	CaseTypePerformance CaseType = "performance"
	// CaseTypeAPI represents API testing
	CaseTypeAPI CaseType = "api"
	// CaseTypeUI represents UI testing
	CaseTypeUI CaseType = "ui"
	// CaseTypeSecurity represents security testing
	CaseTypeSecurity CaseType = "security"
)

// Priority represents test case priority
type Priority string

const (
	// PriorityP0 represents highest priority
	PriorityP0 Priority = "P0"
	// PriorityP1 represents high priority
	PriorityP1 Priority = "P1"
	// PriorityP2 represents medium priority
	PriorityP2 Priority = "P2"
	// PriorityP3 represents low priority
	PriorityP3 Priority = "P3"
)

// TestCase is the aggregate root for testcase context
type TestCase struct {
	id            uuid.UUID
	moduleID      uuid.UUID
	userID        uuid.UUID
	number        CaseNumber
	title         string
	preconditions Preconditions
	steps         Steps
	expected      ExpectedResult
	aiMetadata    *AiMetadata
	caseType      CaseType
	priority      Priority
	status        CaseStatus
	createdAt     time.Time
	updatedAt     time.Time
}

// NewTestCase creates a new test case
func NewTestCase(
	moduleID, userID uuid.UUID,
	number CaseNumber,
	title string,
	preconditions Preconditions,
	steps Steps,
	expected ExpectedResult,
	caseType CaseType,
	priority Priority,
) (*TestCase, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID cannot be nil")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be nil")
	}
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if len(steps) == 0 {
		return nil, ErrEmptySteps
	}

	now := time.Now()
	return &TestCase{
		id:            uuid.New(),
		moduleID:      moduleID,
		userID:        userID,
		number:        number,
		title:         title,
		preconditions: preconditions,
		steps:         steps,
		expected:      expected,
		caseType:      caseType,
		priority:      priority,
		status:        StatusUnexecuted,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// ID returns the test case's unique identifier
func (tc *TestCase) ID() uuid.UUID {
	return tc.id
}

// ModuleID returns the associated module's ID
func (tc *TestCase) ModuleID() uuid.UUID {
	return tc.moduleID
}

// UserID returns the creator's user ID
func (tc *TestCase) UserID() uuid.UUID {
	return tc.userID
}

// Number returns the test case number
func (tc *TestCase) Number() CaseNumber {
	return tc.number
}

// Title returns the test case title
func (tc *TestCase) Title() string {
	return tc.title
}

// Preconditions returns the test case preconditions
func (tc *TestCase) Preconditions() Preconditions {
	return tc.preconditions
}

// Steps returns the test case execution steps
func (tc *TestCase) Steps() Steps {
	return tc.steps
}

// ExpectedResult returns the expected result
func (tc *TestCase) ExpectedResult() ExpectedResult {
	return tc.expected
}

// AiMetadata returns the AI metadata (if AI generated)
func (tc *TestCase) AiMetadata() *AiMetadata {
	return tc.aiMetadata
}

// CaseType returns the test case type
func (tc *TestCase) CaseType() CaseType {
	return tc.caseType
}

// Priority returns the test case priority
func (tc *TestCase) Priority() Priority {
	return tc.priority
}

// Status returns the test case status
func (tc *TestCase) Status() CaseStatus {
	return tc.status
}

// CreatedAt returns the creation timestamp
func (tc *TestCase) CreatedAt() time.Time {
	return tc.createdAt
}

// UpdatedAt returns the last update timestamp
func (tc *TestCase) UpdatedAt() time.Time {
	return tc.updatedAt
}

// UpdateStatus updates the test case status
func (tc *TestCase) UpdateStatus(status CaseStatus) {
	tc.status = status
	tc.updatedAt = time.Now()
}

// UpdateTitle updates the test case title
func (tc *TestCase) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	tc.title = title
	tc.updatedAt = time.Now()
	return nil
}

// UpdateSteps updates the test case steps
func (tc *TestCase) UpdateSteps(steps Steps) error {
	if len(steps) == 0 {
		return ErrEmptySteps
	}
	tc.steps = steps
	tc.updatedAt = time.Now()
	return nil
}

// SetAiMetadata sets the AI metadata
func (tc *TestCase) SetAiMetadata(metadata *AiMetadata) {
	tc.aiMetadata = metadata
	tc.updatedAt = time.Now()
}

// ReconstructTestCase reconstructs a TestCase from persistence layer.
// Used by repository implementations to hydrate domain objects from database rows.
func ReconstructTestCase(
	id uuid.UUID,
	moduleID uuid.UUID,
	userID uuid.UUID,
	number CaseNumber,
	title string,
	preconditions Preconditions,
	steps Steps,
	expected ExpectedResult,
	aiMetadata *AiMetadata,
	caseType CaseType,
	priority Priority,
	status CaseStatus,
	createdAt time.Time,
	updatedAt time.Time,
) *TestCase {
	return &TestCase{
		id:            id,
		moduleID:      moduleID,
		userID:        userID,
		number:        number,
		title:         title,
		preconditions: preconditions,
		steps:         steps,
		expected:      expected,
		aiMetadata:    aiMetadata,
		caseType:      caseType,
		priority:      priority,
		status:        status,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
