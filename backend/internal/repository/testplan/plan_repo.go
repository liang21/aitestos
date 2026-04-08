// Package testplan provides test plan repository implementation
package testplan

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domaintestplan "github.com/liang21/aitestos/internal/domain/testplan"
)

// TestPlanRepository implements domaintestplan.TestPlanRepository interface
type TestPlanRepository struct {
	db *sqlx.DB
}

// NewTestPlanRepository creates a new test plan repository
func NewTestPlanRepository(db *sqlx.DB) *TestPlanRepository {
	return &TestPlanRepository{db: db}
}

// Save persists a new test plan
func (r *TestPlanRepository) Save(ctx context.Context, plan *domaintestplan.TestPlan) error {
	query := `
		INSERT INTO test_plans (id, project_id, name, description, status, case_ids, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	caseIDsJSON, _ := toJSON(plan.CaseIDs())

	_, err := r.db.ExecContext(ctx, query,
		plan.ID(),
		plan.ProjectID(),
		plan.Name(),
		plan.Description(),
		string(plan.Status()),
		caseIDsJSON,
		plan.CreatedBy(),
		plan.CreatedAt(),
		plan.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save test plan: %w", err)
	}
	return nil
}

// FindByID retrieves a test plan by ID
func (r *TestPlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintestplan.TestPlan, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		ProjectID   uuid.UUID `db:"project_id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
		CaseIDs     string    `db:"case_ids"`
		CreatedBy   uuid.UUID `db:"created_by"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	query := `
		SELECT id, project_id, name, description, status, case_ids, created_by, created_at, updated_at
		FROM test_plans
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaintestplan.ErrPlanNotFound
		}
		return nil, fmt.Errorf("find test plan by id: %w", err)
	}

	var caseIDs []uuid.UUID
	if err := fromJSON(row.CaseIDs, &caseIDs); err != nil {
		return nil, fmt.Errorf("parse case ids: %w", err)
	}

	status, err := domaintestplan.ParsePlanStatus(row.Status)
	if err != nil {
		return nil, fmt.Errorf("parse plan status: %w", err)
	}

	return domaintestplan.Reconstruct(
		row.ID,
		row.ProjectID,
		row.Name,
		row.Description,
		status,
		caseIDs,
		row.CreatedBy,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// FindByProjectID retrieves all test plans for a project with pagination
func (r *TestPlanRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestPlan, error) {
	query := `
		SELECT id, project_id, name, description, status, case_ids, created_by, created_at, updated_at
		FROM test_plans
		WHERE project_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []struct {
		ID          uuid.UUID `db:"id"`
		ProjectID   uuid.UUID `db:"project_id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
		CaseIDs     string    `db:"case_ids"`
		CreatedBy   uuid.UUID `db:"created_by"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, projectID, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test plans by project id: %w", err)
	}

	plans := make([]*domaintestplan.TestPlan, 0, len(rows))
	for _, row := range rows {
		var caseIDs []uuid.UUID
		if err := fromJSON(row.CaseIDs, &caseIDs); err != nil {
			return nil, fmt.Errorf("parse case ids: %w", err)
		}

		status, err := domaintestplan.ParsePlanStatus(row.Status)
		if err != nil {
			return nil, fmt.Errorf("parse plan status: %w", err)
		}

		plan := domaintestplan.Reconstruct(
			row.ID,
			row.ProjectID,
			row.Name,
			row.Description,
			status,
			caseIDs,
			row.CreatedBy,
			parseTime(row.CreatedAt),
			parseTime(row.UpdatedAt),
		)
		plans = append(plans, plan)
	}

	return plans, nil
}

// Update updates an existing test plan
func (r *TestPlanRepository) Update(ctx context.Context, plan *domaintestplan.TestPlan) error {
	caseIDsJSON, _ := toJSON(plan.CaseIDs())

	query := `
		UPDATE test_plans
		SET name = $2, description = $3, status = $4, case_ids = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
		plan.ID(),
		plan.Name(),
		plan.Description(),
		string(plan.Status()),
		caseIDsJSON,
		plan.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("update test plan: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaintestplan.ErrPlanNotFound
	}
	return nil
}

// Delete removes a test plan (soft delete)
func (r *TestPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE test_plans SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete test plan: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaintestplan.ErrPlanNotFound
	}
	return nil
}
