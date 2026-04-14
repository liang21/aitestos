// Package generation defines GenerationTask aggregate
package generation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// GenerationTask is the aggregate root for generation context
type GenerationTask struct {
	id            uuid.UUID
	projectID     uuid.UUID
	moduleID      uuid.UUID
	userID        uuid.UUID
	status        TaskStatus
	prompt        string
	resultSummary map[string]any
	errorMsg      string
	createdAt     time.Time
	updatedAt     time.Time
}

// NewGenerationTask creates a new generation task
func NewGenerationTask(projectID, moduleID uuid.UUID, prompt string, userID uuid.UUID) (*GenerationTask, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	if prompt == "" {
		return nil, errors.New("prompt cannot be empty")
	}
	if len(prompt) < 10 {
		return nil, errors.New("prompt must be at least 10 characters")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be nil")
	}

	now := time.Now()
	return &GenerationTask{
		id:        uuid.New(),
		projectID: projectID,
		moduleID:  moduleID,
		userID:    userID,
		status:    TaskPending,
		prompt:    prompt,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the task's unique identifier
func (t *GenerationTask) ID() uuid.UUID {
	return t.id
}

// ProjectID returns the associated project's ID
func (t *GenerationTask) ProjectID() uuid.UUID {
	return t.projectID
}

// ModuleID returns the associated module's ID
func (t *GenerationTask) ModuleID() uuid.UUID {
	return t.moduleID
}

// CreatedBy returns the user who created this task
func (t *GenerationTask) CreatedBy() uuid.UUID {
	return t.userID
}

// Status returns the task's status
func (t *GenerationTask) Status() TaskStatus {
	return t.status
}

// Prompt returns the generation prompt
func (t *GenerationTask) Prompt() string {
	return t.prompt
}

// ResultSummary returns the result summary after completion
func (t *GenerationTask) ResultSummary() map[string]any {
	return t.resultSummary
}

// ErrorMsg returns the error message if failed
func (t *GenerationTask) ErrorMsg() string {
	return t.errorMsg
}

// CreatedAt returns the creation timestamp
func (t *GenerationTask) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the last update timestamp
func (t *GenerationTask) UpdatedAt() time.Time {
	return t.updatedAt
}

// StartProcessing starts processing the task
func (t *GenerationTask) StartProcessing() error {
	if !t.status.CanTransitionTo(TaskProcessing) {
		return errors.New("cannot start processing from current status")
	}
	t.status = TaskProcessing
	t.updatedAt = time.Now()
	return nil
}

// Complete completes the task with a result summary
func (t *GenerationTask) Complete(summary map[string]any) {
	t.status = TaskCompleted
	t.resultSummary = summary
	t.updatedAt = time.Now()
}

// Fail marks the task as failed with an error message
func (t *GenerationTask) Fail(errMsg string) {
	t.status = TaskFailed
	t.errorMsg = errMsg
	t.updatedAt = time.Now()
}

// ReconstructTask reconstructs a GenerationTask from stored data
func ReconstructTask(
	id uuid.UUID,
	projectID uuid.UUID,
	moduleID uuid.UUID,
	userID uuid.UUID,
	prompt string,
	status TaskStatus,
	resultSummary map[string]any,
	errorMsg string,
	createdAt time.Time,
	updatedAt time.Time,
) *GenerationTask {
	return &GenerationTask{
		id:            id,
		projectID:     projectID,
		moduleID:      moduleID,
		userID:        userID,
		status:        status,
		prompt:        prompt,
		resultSummary: resultSummary,
		errorMsg:      errorMsg,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
