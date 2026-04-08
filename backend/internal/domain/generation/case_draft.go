// Package generation defines GeneratedCaseDraft entity
package generation

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

// GeneratedCaseDraft is an entity representing a generated test case draft
type GeneratedCaseDraft struct {
	id            uuid.UUID
	taskID        uuid.UUID
	moduleID      *uuid.UUID
	title         string
	preconditions testcase.Preconditions
	steps         testcase.Steps
	expected      testcase.ExpectedResult
	caseType      testcase.CaseType
	priority      testcase.Priority
	aiMetadata    *testcase.AiMetadata
	status        DraftStatus
	feedback      string
	createdAt     time.Time
	updatedAt     time.Time
}

// NewGeneratedCaseDraft creates a new generated case draft
func NewGeneratedCaseDraft(
	taskID uuid.UUID,
	title string,
	preconditions testcase.Preconditions,
	steps testcase.Steps,
	expected testcase.ExpectedResult,
	caseType testcase.CaseType,
	priority testcase.Priority,
) (*GeneratedCaseDraft, error) {
	if taskID == uuid.Nil {
		return nil, errors.New("task ID cannot be nil")
	}
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if len(steps) == 0 {
		return nil, errors.New("steps cannot be empty")
	}

	now := time.Now()
	return &GeneratedCaseDraft{
		id:            uuid.New(),
		taskID:        taskID,
		title:         title,
		preconditions: preconditions,
		steps:         steps,
		expected:      expected,
		caseType:      caseType,
		priority:      priority,
		status:        DraftPending,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// ID returns the draft's unique identifier
func (d *GeneratedCaseDraft) ID() uuid.UUID {
	return d.id
}

// TaskID returns the associated task's ID
func (d *GeneratedCaseDraft) TaskID() uuid.UUID {
	return d.taskID
}

// ModuleID returns the associated module's ID (nil if not yet confirmed)
func (d *GeneratedCaseDraft) ModuleID() *uuid.UUID {
	return d.moduleID
}

// Title returns the draft's title
func (d *GeneratedCaseDraft) Title() string {
	return d.title
}

// Preconditions returns the preconditions
func (d *GeneratedCaseDraft) Preconditions() testcase.Preconditions {
	return d.preconditions
}

// Steps returns the test steps
func (d *GeneratedCaseDraft) Steps() testcase.Steps {
	return d.steps
}

// ExpectedResult returns the expected result
func (d *GeneratedCaseDraft) ExpectedResult() testcase.ExpectedResult {
	return d.expected
}

// CaseType returns the case type
func (d *GeneratedCaseDraft) CaseType() testcase.CaseType {
	return d.caseType
}

// Priority returns the priority
func (d *GeneratedCaseDraft) Priority() testcase.Priority {
	return d.priority
}

// AiMetadata returns the AI metadata
func (d *GeneratedCaseDraft) AiMetadata() *testcase.AiMetadata {
	return d.aiMetadata
}

// Status returns the draft's status
func (d *GeneratedCaseDraft) Status() DraftStatus {
	return d.status
}

// Feedback returns the rejection feedback
func (d *GeneratedCaseDraft) Feedback() string {
	return d.feedback
}

// CreatedAt returns the creation timestamp
func (d *GeneratedCaseDraft) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns the last update timestamp
func (d *GeneratedCaseDraft) UpdatedAt() time.Time {
	return d.updatedAt
}

// Confirm confirms the draft with a module ID
func (d *GeneratedCaseDraft) Confirm(moduleID uuid.UUID) error {
	if !d.status.CanTransitionTo(DraftConfirmed) {
		return errors.New("cannot confirm draft from current status")
	}
	d.moduleID = &moduleID
	d.status = DraftConfirmed
	d.updatedAt = time.Now()
	return nil
}

// Reject rejects the draft with a reason and feedback
func (d *GeneratedCaseDraft) Reject(reason RejectionReason, detail string) error {
	if !d.status.CanTransitionTo(DraftRejected) {
		return errors.New("cannot reject draft from current status")
	}
	d.status = DraftRejected
	d.feedback = string(reason) + ": " + detail
	d.updatedAt = time.Now()
	return nil
}

// SetAiMetadata sets the AI metadata
func (d *GeneratedCaseDraft) SetAiMetadata(metadata *testcase.AiMetadata) {
	d.aiMetadata = metadata
	d.updatedAt = time.Now()
}

// ReconstructDraft reconstructs a GeneratedCaseDraft from stored data
func ReconstructDraft(
	id uuid.UUID,
	taskID uuid.UUID,
	moduleID *uuid.UUID,
	title string,
	preconditions testcase.Preconditions,
	steps testcase.Steps,
	expected testcase.ExpectedResult,
	caseType testcase.CaseType,
	priority testcase.Priority,
	aiMetadata *testcase.AiMetadata,
	status DraftStatus,
	feedback string,
	createdAt time.Time,
	updatedAt time.Time,
) *GeneratedCaseDraft {
	return &GeneratedCaseDraft{
		id:            id,
		taskID:        taskID,
		moduleID:      moduleID,
		title:         title,
		preconditions: preconditions,
		steps:         steps,
		expected:      expected,
		caseType:      caseType,
		priority:      priority,
		aiMetadata:    aiMetadata,
		status:        status,
		feedback:      feedback,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
