// Package generation provides case draft repository implementation
package generation

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domaingeneration "github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

// CaseDraftRepository implements domaingeneration.CaseDraftRepository interface
type CaseDraftRepository struct {
	db *sqlx.DB
}

// NewCaseDraftRepository creates a new case draft repository
func NewCaseDraftRepository(db *sqlx.DB) *CaseDraftRepository {
	return &CaseDraftRepository{db: db}
}

// Save persists a new case draft
func (r *CaseDraftRepository) Save(ctx context.Context, draft *domaingeneration.GeneratedCaseDraft) error {
	query := `
		INSERT INTO case_drafts (
			id, task_id, module_id, title, preconditions, steps, expected_result,
			case_type, priority, ai_metadata, status, feedback, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	preconditionsJSON, _ := toJSON(draft.Preconditions())
	stepsJSON, _ := toJSON(draft.Steps())
	expectedJSON, _ := toJSON(draft.ExpectedResult())
	var aiMetadataJSON []byte
	if draft.AiMetadata() != nil {
		aiMetadataJSON, _ = jsonMarshal(draft.AiMetadata())
	}

	var moduleID interface{}
	if draft.ModuleID() != nil {
		moduleID = *draft.ModuleID()
	}

	_, err := r.db.ExecContext(ctx, query,
		draft.ID(),
		draft.TaskID(),
		moduleID,
		draft.Title(),
		preconditionsJSON,
		stepsJSON,
		expectedJSON,
		string(draft.CaseType()),
		string(draft.Priority()),
		aiMetadataJSON,
		string(draft.Status()),
		draft.Feedback(),
		draft.CreatedAt(),
		draft.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save case draft: %w", err)
	}
	return nil
}

// FindByID retrieves a case draft by ID
func (r *CaseDraftRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaingeneration.GeneratedCaseDraft, error) {
	query := `
		SELECT id, task_id, module_id, title, preconditions, steps, expected_result,
			   case_type, priority, ai_metadata, status, feedback, created_at, updated_at
		FROM case_drafts
		WHERE id = $1
	`

	var row struct {
		ID            uuid.UUID  `db:"id"`
		TaskID        uuid.UUID  `db:"task_id"`
		ModuleID      *uuid.UUID `db:"module_id"`
		Title         string     `db:"title"`
		Preconditions string     `db:"preconditions"`
		Steps         string     `db:"steps"`
		Expected      string     `db:"expected_result"`
		CaseType      string     `db:"case_type"`
		Priority      string     `db:"priority"`
		AiMetadata    []byte     `db:"ai_metadata"`
		Status        string     `db:"status"`
		Feedback      string     `db:"feedback"`
		CreatedAt     string     `db:"created_at"`
		UpdatedAt     string     `db:"updated_at"`
	}

	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaingeneration.ErrDraftNotFound
		}
		return nil, fmt.Errorf("find case draft by id: %w", err)
	}

	return r.rowToDraft(&row)
}

// FindByTaskID retrieves all case drafts for a generation task
func (r *CaseDraftRepository) FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*domaingeneration.GeneratedCaseDraft, error) {
	query := `
		SELECT id, task_id, module_id, title, preconditions, steps, expected_result,
			   case_type, priority, ai_metadata, status, feedback, created_at, updated_at
		FROM case_drafts
		WHERE task_id = $1
		ORDER BY created_at ASC
	`

	var rows []struct {
		ID            uuid.UUID  `db:"id"`
		TaskID        uuid.UUID  `db:"task_id"`
		ModuleID      *uuid.UUID `db:"module_id"`
		Title         string     `db:"title"`
		Preconditions string     `db:"preconditions"`
		Steps         string     `db:"steps"`
		Expected      string     `db:"expected_result"`
		CaseType      string     `db:"case_type"`
		Priority      string     `db:"priority"`
		AiMetadata    []byte     `db:"ai_metadata"`
		Status        string     `db:"status"`
		Feedback      string     `db:"feedback"`
		CreatedAt     string     `db:"created_at"`
		UpdatedAt     string     `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, taskID); err != nil {
		return nil, fmt.Errorf("find case drafts by task id: %w", err)
	}

	return r.rowsToDrafts(rows)
}

// FindByStatus retrieves case drafts by status with pagination
func (r *CaseDraftRepository) FindByStatus(ctx context.Context, status domaingeneration.DraftStatus, opts domaingeneration.QueryOptions) ([]*domaingeneration.GeneratedCaseDraft, error) {
	query := `
		SELECT id, task_id, module_id, title, preconditions, steps, expected_result,
			   case_type, priority, ai_metadata, status, feedback, created_at, updated_at
		FROM case_drafts
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []struct {
		ID            uuid.UUID  `db:"id"`
		TaskID        uuid.UUID  `db:"task_id"`
		ModuleID      *uuid.UUID `db:"module_id"`
		Title         string     `db:"title"`
		Preconditions string     `db:"preconditions"`
		Steps         string     `db:"steps"`
		Expected      string     `db:"expected_result"`
		CaseType      string     `db:"case_type"`
		Priority      string     `db:"priority"`
		AiMetadata    []byte     `db:"ai_metadata"`
		Status        string     `db:"status"`
		Feedback      string     `db:"feedback"`
		CreatedAt     string     `db:"created_at"`
		UpdatedAt     string     `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, string(status), opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find case drafts by status: %w", err)
	}

	return r.rowsToDrafts(rows)
}

// FindByTaskIDAndStatus retrieves case drafts by task ID and status
func (r *CaseDraftRepository) FindByTaskIDAndStatus(ctx context.Context, taskID uuid.UUID, status domaingeneration.DraftStatus) ([]*domaingeneration.GeneratedCaseDraft, error) {
	query := `
		SELECT id, task_id, module_id, title, preconditions, steps, expected_result,
			   case_type, priority, ai_metadata, status, feedback, created_at, updated_at
		FROM case_drafts
		WHERE task_id = $1 AND status = $2
		ORDER BY created_at ASC
	`

	var rows []struct {
		ID            uuid.UUID  `db:"id"`
		TaskID        uuid.UUID  `db:"task_id"`
		ModuleID      *uuid.UUID `db:"module_id"`
		Title         string     `db:"title"`
		Preconditions string     `db:"preconditions"`
		Steps         string     `db:"steps"`
		Expected      string     `db:"expected_result"`
		CaseType      string     `db:"case_type"`
		Priority      string     `db:"priority"`
		AiMetadata    []byte     `db:"ai_metadata"`
		Status        string     `db:"status"`
		Feedback      string     `db:"feedback"`
		CreatedAt     string     `db:"created_at"`
		UpdatedAt     string     `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, taskID, string(status)); err != nil {
		return nil, fmt.Errorf("find case drafts by task id and status: %w", err)
	}

	return r.rowsToDrafts(rows)
}

// BatchUpdateStatus updates the status and module ID for multiple drafts in a transaction
func (r *CaseDraftRepository) BatchUpdateStatus(ctx context.Context, draftIDs []uuid.UUID, status domaingeneration.DraftStatus, moduleID uuid.UUID) error {
	if len(draftIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE case_drafts
		SET status = $2, module_id = $3, updated_at = NOW()
		WHERE id = $1
	`

	for _, id := range draftIDs {
		_, err := tx.ExecContext(ctx, query, id, string(status), moduleID)
		if err != nil {
			return fmt.Errorf("batch update draft status for %s: %w", id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit batch update: %w", err)
	}
	return nil
}

// CountByTaskIDAndStatus counts drafts by task ID and status
func (r *CaseDraftRepository) CountByTaskIDAndStatus(ctx context.Context, taskID uuid.UUID, status domaingeneration.DraftStatus) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM case_drafts WHERE task_id = $1 AND status = $2`
	err := r.db.GetContext(ctx, &count, query, taskID, string(status))
	if err != nil {
		return 0, fmt.Errorf("count case drafts by task id and status: %w", err)
	}
	return count, nil
}

// Update updates an existing case draft
func (r *CaseDraftRepository) Update(ctx context.Context, draft *domaingeneration.GeneratedCaseDraft) error {
	query := `
		UPDATE case_drafts
		SET module_id = $2, status = $3, feedback = $4, updated_at = $5
		WHERE id = $1
	`

	var moduleID interface{}
	if draft.ModuleID() != nil {
		moduleID = *draft.ModuleID()
	}

	result, err := r.db.ExecContext(ctx, query,
		draft.ID(),
		moduleID,
		string(draft.Status()),
		draft.Feedback(),
		draft.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("update case draft: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaingeneration.ErrDraftNotFound
	}
	return nil
}

// Delete removes a case draft
func (r *CaseDraftRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM case_drafts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete case draft: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaingeneration.ErrDraftNotFound
	}
	return nil
}

// DeleteByTaskID removes all case drafts for a task
func (r *CaseDraftRepository) DeleteByTaskID(ctx context.Context, taskID uuid.UUID) error {
	query := `DELETE FROM case_drafts WHERE task_id = $1`
	_, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("delete case drafts by task id: %w", err)
	}
	return nil
}

// Helper functions
func (r *CaseDraftRepository) rowToDraft(row *struct {
	ID            uuid.UUID  `db:"id"`
	TaskID        uuid.UUID  `db:"task_id"`
	ModuleID      *uuid.UUID `db:"module_id"`
	Title         string     `db:"title"`
	Preconditions string     `db:"preconditions"`
	Steps         string     `db:"steps"`
	Expected      string     `db:"expected_result"`
	CaseType      string     `db:"case_type"`
	Priority      string     `db:"priority"`
	AiMetadata    []byte     `db:"ai_metadata"`
	Status        string     `db:"status"`
	Feedback      string     `db:"feedback"`
	CreatedAt     string     `db:"created_at"`
	UpdatedAt     string     `db:"updated_at"`
}) (*domaingeneration.GeneratedCaseDraft, error) {
	status, err := domaingeneration.ParseDraftStatus(row.Status)
	if err != nil {
		return nil, fmt.Errorf("parse draft status: %w", err)
	}

	var preconditions testcase.Preconditions
	if row.Preconditions != "" {
		if err := fromJSON(row.Preconditions, &preconditions); err != nil {
			return nil, fmt.Errorf("parse preconditions: %w", err)
		}
	}

	var steps testcase.Steps
	if row.Steps != "" {
		if err := fromJSON(row.Steps, &steps); err != nil {
			return nil, fmt.Errorf("parse steps: %w", err)
		}
	}

	var expected testcase.ExpectedResult
	if row.Expected != "" {
		if err := fromJSON(row.Expected, &expected); err != nil {
			return nil, fmt.Errorf("parse expected result: %w", err)
		}
	}

	var aiMetadata *testcase.AiMetadata
	if len(row.AiMetadata) > 0 && string(row.AiMetadata) != "{}" {
		aiMetadata = &testcase.AiMetadata{}
		if err := json.Unmarshal(row.AiMetadata, aiMetadata); err != nil {
			return nil, fmt.Errorf("parse ai metadata: %w", err)
		}
	}

	return domaingeneration.ReconstructDraft(
		row.ID,
		row.TaskID,
		row.ModuleID,
		row.Title,
		preconditions,
		steps,
		expected,
		testcase.CaseType(row.CaseType),
		testcase.Priority(row.Priority),
		aiMetadata,
		status,
		row.Feedback,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

func (r *CaseDraftRepository) rowsToDrafts(rows []struct {
	ID            uuid.UUID  `db:"id"`
	TaskID        uuid.UUID  `db:"task_id"`
	ModuleID      *uuid.UUID `db:"module_id"`
	Title         string     `db:"title"`
	Preconditions string     `db:"preconditions"`
	Steps         string     `db:"steps"`
	Expected      string     `db:"expected_result"`
	CaseType      string     `db:"case_type"`
	Priority      string     `db:"priority"`
	AiMetadata    []byte     `db:"ai_metadata"`
	Status        string     `db:"status"`
	Feedback      string     `db:"feedback"`
	CreatedAt     string     `db:"created_at"`
	UpdatedAt     string     `db:"updated_at"`
}) ([]*domaingeneration.GeneratedCaseDraft, error) {
	drafts := make([]*domaingeneration.GeneratedCaseDraft, 0, len(rows))
	for _, row := range rows {
		draft, err := r.rowToDraft(&row)
		if err != nil {
			return nil, err
		}
		drafts = append(drafts, draft)
	}
	return drafts, nil
}
