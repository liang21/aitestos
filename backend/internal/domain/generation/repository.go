// Package generation defines repository interfaces
package generation

import (
	"context"

	"github.com/google/uuid"
)

// GenerationTaskRepository defines the interface for generation task persistence
type GenerationTaskRepository interface {
	// Save persists a new generation task
	Save(ctx context.Context, task *GenerationTask) error

	// FindByID retrieves a generation task by ID
	FindByID(ctx context.Context, id uuid.UUID) (*GenerationTask, error)

	// FindByProjectID retrieves all generation tasks for a project with pagination
	FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*GenerationTask, error)

	// FindByStatus retrieves generation tasks by status with pagination
	FindByStatus(ctx context.Context, status TaskStatus, opts QueryOptions) ([]*GenerationTask, error)

	// FindByUserID retrieves all generation tasks for a user with pagination
	FindByUserID(ctx context.Context, userID uuid.UUID, opts QueryOptions) ([]*GenerationTask, error)

	// Update updates an existing generation task
	Update(ctx context.Context, task *GenerationTask) error

	// Delete removes a generation task
	Delete(ctx context.Context, id uuid.UUID) error
}

// CaseDraftRepository defines the interface for case draft persistence
type CaseDraftRepository interface {
	// Save persists a new case draft
	Save(ctx context.Context, draft *GeneratedCaseDraft) error

	// FindByID retrieves a case draft by ID
	FindByID(ctx context.Context, id uuid.UUID) (*GeneratedCaseDraft, error)

	// FindByTaskID retrieves all case drafts for a generation task
	FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*GeneratedCaseDraft, error)

	// FindByStatus retrieves case drafts by status with pagination
	FindByStatus(ctx context.Context, status DraftStatus, opts QueryOptions) ([]*GeneratedCaseDraft, error)

	// FindByTaskIDAndStatus retrieves case drafts by task ID and status
	FindByTaskIDAndStatus(ctx context.Context, taskID uuid.UUID, status DraftStatus) ([]*GeneratedCaseDraft, error)

	// BatchUpdateStatus updates the status and module ID for multiple drafts
	BatchUpdateStatus(ctx context.Context, draftIDs []uuid.UUID, status DraftStatus, moduleID uuid.UUID) error

	// CountByTaskIDAndStatus counts drafts by task ID and status
	CountByTaskIDAndStatus(ctx context.Context, taskID uuid.UUID, status DraftStatus) (int64, error)

	// Update updates an existing case draft
	Update(ctx context.Context, draft *GeneratedCaseDraft) error

	// Delete removes a case draft
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByTaskID removes all case drafts for a task
	DeleteByTaskID(ctx context.Context, taskID uuid.UUID) error
}

// QueryOptions holds pagination and filtering options
type QueryOptions struct {
	Offset   int
	Limit    int
	OrderBy  string
	Keywords string
}
