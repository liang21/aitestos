// Package generation provides generation task repository implementation
package generation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domaingeneration "github.com/liang21/aitestos/internal/domain/generation"
)

// GenerationTaskRepository implements domaingeneration.GenerationTaskRepository interface
type GenerationTaskRepository struct {
	db *sqlx.DB
}

// NewGenerationTaskRepository creates a new generation task repository
func NewGenerationTaskRepository(db *sqlx.DB) *GenerationTaskRepository {
	return &GenerationTaskRepository{db: db}
}

// Save persists a new generation task
func (r *GenerationTaskRepository) Save(ctx context.Context, task *domaingeneration.GenerationTask) error {
	query := `
		INSERT INTO generation_tasks (id, project_id, module_id, user_id, prompt, status, result_summary, error_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	resultSummaryJSON, _ := toJSON(task.ResultSummary())

	_, err := r.db.ExecContext(ctx, query,
		task.ID(),
		task.ProjectID(),
		task.ModuleID(),
		task.CreatedBy(),
		task.Prompt(),
		string(task.Status()),
		resultSummaryJSON,
		task.ErrorMsg(),
		task.CreatedAt(),
		task.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save generation task: %w", err)
	}
	return nil
}

// FindByID retrieves a generation task by ID
func (r *GenerationTaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaingeneration.GenerationTask, error) {
	var row struct {
		ID            uuid.UUID `db:"id"`
		ProjectID     uuid.UUID `db:"project_id"`
		ModuleID      uuid.UUID `db:"module_id"`
		UserID        uuid.UUID `db:"user_id"`
		Prompt        string    `db:"prompt"`
		Status        string    `db:"status"`
		ResultSummary string    `db:"result_summary"`
		ErrorMessage  string    `db:"error_message"`
		CreatedAt     string    `db:"created_at"`
		UpdatedAt     string    `db:"updated_at"`
	}

	query := `
		SELECT id, project_id, module_id, user_id, prompt, status, result_summary, error_message, created_at, updated_at
		FROM generation_tasks
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaingeneration.ErrTaskNotFound
		}
		return nil, fmt.Errorf("find generation task by id: %w", err)
	}

	status, err := domaingeneration.ParseTaskStatus(row.Status)
	if err != nil {
		return nil, fmt.Errorf("parse task status: %w", err)
	}

	var resultSummary map[string]any
	if row.ResultSummary != "" {
		if err := fromJSON(row.ResultSummary, &resultSummary); err != nil {
			return nil, fmt.Errorf("parse result summary: %w", err)
		}
	}

	return domaingeneration.ReconstructTask(
		row.ID,
		row.ProjectID,
		row.ModuleID,
		row.UserID,
		row.Prompt,
		status,
		resultSummary,
		row.ErrorMessage,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// FindByProjectID retrieves all generation tasks for a project with pagination
func (r *GenerationTaskRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domaingeneration.QueryOptions) ([]*domaingeneration.GenerationTask, error) {
	query := `
		SELECT id, project_id, module_id, user_id, prompt, status, result_summary, error_message, created_at, updated_at
		FROM generation_tasks
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	return r.findTasks(ctx, query, projectID, opts.Limit, opts.Offset)
}

// FindByStatus retrieves generation tasks by status with pagination
func (r *GenerationTaskRepository) FindByStatus(ctx context.Context, status domaingeneration.TaskStatus, opts domaingeneration.QueryOptions) ([]*domaingeneration.GenerationTask, error) {
	query := `
		SELECT id, project_id, module_id, user_id, prompt, status, result_summary, error_message, created_at, updated_at
		FROM generation_tasks
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	return r.findTasks(ctx, query, string(status), opts.Limit, opts.Offset)
}

// FindByUserID retrieves all generation tasks for a user with pagination
func (r *GenerationTaskRepository) FindByUserID(ctx context.Context, userID uuid.UUID, opts domaingeneration.QueryOptions) ([]*domaingeneration.GenerationTask, error) {
	query := `
		SELECT id, project_id, module_id, user_id, prompt, status, result_summary, error_message, created_at, updated_at
		FROM generation_tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	return r.findTasks(ctx, query, userID, opts.Limit, opts.Offset)
}

// Delete removes a generation task
func (r *GenerationTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM generation_tasks WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete generation task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaingeneration.ErrTaskNotFound
	}
	return nil
}

// Update updates an existing generation task
func (r *GenerationTaskRepository) Update(ctx context.Context, task *domaingeneration.GenerationTask) error {
	resultSummaryJSON, _ := toJSON(task.ResultSummary())

	query := `
		UPDATE generation_tasks
		SET status = $2, result_summary = $3, error_message = $4, updated_at = $5
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		task.ID(),
		string(task.Status()),
		resultSummaryJSON,
		task.ErrorMsg(),
		task.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("update generation task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaingeneration.ErrTaskNotFound
	}
	return nil
}

// Helper functions
func (r *GenerationTaskRepository) findTasks(ctx context.Context, query string, args ...interface{}) ([]*domaingeneration.GenerationTask, error) {
	var rows []struct {
		ID            uuid.UUID `db:"id"`
		ProjectID     uuid.UUID `db:"project_id"`
		ModuleID      uuid.UUID `db:"module_id"`
		UserID        uuid.UUID `db:"user_id"`
		Prompt        string    `db:"prompt"`
		Status        string    `db:"status"`
		ResultSummary string    `db:"result_summary"`
		ErrorMessage  string    `db:"error_message"`
		CreatedAt     string    `db:"created_at"`
		UpdatedAt     string    `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("find generation tasks: %w", err)
	}

	tasks := make([]*domaingeneration.GenerationTask, 0, len(rows))
	for _, row := range rows {
		status, err := domaingeneration.ParseTaskStatus(row.Status)
		if err != nil {
			return nil, fmt.Errorf("parse task status: %w", err)
		}

		var resultSummary map[string]any
		if row.ResultSummary != "" {
			if err := fromJSON(row.ResultSummary, &resultSummary); err != nil {
				return nil, fmt.Errorf("parse result summary: %w", err)
			}
		}

		task := domaingeneration.ReconstructTask(
			row.ID,
			row.ProjectID,
			row.ModuleID,
			row.UserID,
			row.Prompt,
			status,
			resultSummary,
			row.ErrorMessage,
			parseTime(row.CreatedAt),
			parseTime(row.UpdatedAt),
		)
		tasks = append(tasks, task)
	}

	return tasks, nil
}
