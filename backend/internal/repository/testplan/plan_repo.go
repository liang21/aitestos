// Package testplan provides test plan repository implementation
package testplan

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

const planColumns = `id, project_id, user_id, name, description, status, created_at, updated_at`

// planRow maps SQL columns for test_plan table
type planRow struct {
	ID          uuid.UUID `db:"id"`
	ProjectID   uuid.UUID `db:"project_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Save persists a new test plan
func (r *TestPlanRepository) Save(ctx context.Context, plan *domaintestplan.TestPlan) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO test_plan (id, project_id, user_id, name, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, query,
		plan.ID(),
		plan.ProjectID(),
		plan.CreatedBy(), // maps to user_id column
		plan.Name(),
		plan.Description(),
		string(plan.Status()),
		plan.CreatedAt(),
		plan.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save test plan: %w", err)
	}

	// Insert case associations into plan_cases junction table
	for _, caseID := range plan.CaseIDs() {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO plan_cases (plan_id, case_id) VALUES ($1, $2)`,
			plan.ID(), caseID)
		if err != nil {
			return fmt.Errorf("save plan case association: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a test plan by ID
func (r *TestPlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintestplan.TestPlan, error) {
	var row planRow
	query := fmt.Sprintf(`SELECT %s FROM test_plan WHERE id = $1 AND deleted_at IS NULL`, planColumns)
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaintestplan.ErrPlanNotFound
		}
		return nil, fmt.Errorf("find test plan by id: %w", err)
	}

	caseIDs, err := r.findCaseIDs(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find case ids: %w", err)
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
		row.UserID, // maps to createdBy
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

// FindByProjectID retrieves all test plans for a project with pagination
func (r *TestPlanRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestPlan, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	var rows []planRow
	query := fmt.Sprintf(
		`SELECT %s FROM test_plan WHERE project_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		planColumns)
	if err := r.db.SelectContext(ctx, &rows, query, projectID, limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test plans by project id: %w", err)
	}

	plans := make([]*domaintestplan.TestPlan, 0, len(rows))
	for _, row := range rows {
		caseIDs, err := r.findCaseIDs(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("find case ids for plan %s: %w", row.ID, err)
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
			row.UserID,
			row.CreatedAt,
			row.UpdatedAt,
		)
		plans = append(plans, plan)
	}

	return plans, nil
}

// Update updates an existing test plan
func (r *TestPlanRepository) Update(ctx context.Context, plan *domaintestplan.TestPlan) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE test_plan
		SET name = $2, description = $3, status = $4, updated_at = $5
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(ctx, query,
		plan.ID(),
		plan.Name(),
		plan.Description(),
		string(plan.Status()),
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

	// Reconcile case associations: delete all, then re-insert
	if _, err := tx.ExecContext(ctx, `DELETE FROM plan_cases WHERE plan_id = $1`, plan.ID()); err != nil {
		return fmt.Errorf("delete plan cases: %w", err)
	}
	for _, caseID := range plan.CaseIDs() {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO plan_cases (plan_id, case_id) VALUES ($1, $2)`,
			plan.ID(), caseID); err != nil {
			return fmt.Errorf("insert plan case: %w", err)
		}
	}

	return tx.Commit()
}

// Delete removes a test plan (soft delete)
func (r *TestPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE test_plan SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
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

// FindByStatus retrieves test plans by status with pagination
func (r *TestPlanRepository) FindByStatus(ctx context.Context, status domaintestplan.PlanStatus, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestPlan, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	var rows []planRow
	query := fmt.Sprintf(
		`SELECT %s FROM test_plan WHERE status = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		planColumns)
	if err := r.db.SelectContext(ctx, &rows, query, string(status), limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test plans by status: %w", err)
	}

	plans := make([]*domaintestplan.TestPlan, 0, len(rows))
	for _, row := range rows {
		caseIDs, err := r.findCaseIDs(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("find case ids for plan %s: %w", row.ID, err)
		}

		planStatus, err := domaintestplan.ParsePlanStatus(row.Status)
		if err != nil {
			return nil, fmt.Errorf("parse plan status: %w", err)
		}

		plan := domaintestplan.Reconstruct(
			row.ID,
			row.ProjectID,
			row.Name,
			row.Description,
			planStatus,
			caseIDs,
			row.UserID,
			row.CreatedAt,
			row.UpdatedAt,
		)
		plans = append(plans, plan)
	}

	return plans, nil
}

// AddCase adds a test case to a plan
func (r *TestPlanRepository) AddCase(ctx context.Context, planID, caseID uuid.UUID) error {
	query := `INSERT INTO plan_cases (plan_id, case_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, planID, caseID)
	if err != nil {
		return fmt.Errorf("add case to plan: %w", err)
	}
	return nil
}

// RemoveCase removes a test case from a plan
func (r *TestPlanRepository) RemoveCase(ctx context.Context, planID, caseID uuid.UUID) error {
	query := `DELETE FROM plan_cases WHERE plan_id = $1 AND case_id = $2`
	_, err := r.db.ExecContext(ctx, query, planID, caseID)
	if err != nil {
		return fmt.Errorf("remove case from plan: %w", err)
	}
	return nil
}

// GetCaseIDs retrieves all case IDs for a plan
func (r *TestPlanRepository) GetCaseIDs(ctx context.Context, planID uuid.UUID) ([]uuid.UUID, error) {
	caseIDs, err := r.findCaseIDs(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("get case ids: %w", err)
	}
	return caseIDs, nil
}

// UpdateStatus updates the plan's status
func (r *TestPlanRepository) UpdateStatus(ctx context.Context, planID uuid.UUID, status domaintestplan.PlanStatus) error {
	query := `UPDATE test_plan SET status = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, planID, string(status))
	if err != nil {
		return fmt.Errorf("update plan status: %w", err)
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

// findCaseIDs retrieves all case IDs for a plan from the junction table
func (r *TestPlanRepository) findCaseIDs(ctx context.Context, planID uuid.UUID) ([]uuid.UUID, error) {
	var caseIDs []uuid.UUID
	query := `SELECT case_id FROM plan_cases WHERE plan_id = $1`
	if err := r.db.SelectContext(ctx, &caseIDs, query, planID); err != nil {
		return nil, fmt.Errorf("find case ids for plan: %w", err)
	}
	if caseIDs == nil {
		caseIDs = []uuid.UUID{}
	}
	return caseIDs, nil
}
