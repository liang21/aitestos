// Package testplan provides test result repository implementation
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

// TestResultRepository implements domaintestplan.TestResultRepository interface
type TestResultRepository struct {
	db *sqlx.DB
}

// NewTestResultRepository creates a new test result repository
func NewTestResultRepository(db *sqlx.DB) *TestResultRepository {
	return &TestResultRepository{db: db}
}

// Save persists a test result
func (r *TestResultRepository) Save(ctx context.Context, result *domaintestplan.TestResult) error {
	query := `
		INSERT INTO test_results (
			id, plan_id, case_id, executor_id, status, note, executed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		result.ID(),
		result.PlanID(),
		result.CaseID(),
		result.ExecutedBy(),
		string(result.Status()),
		result.Note(),
		result.ExecutedAt(),
		result.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save test result: %w", err)
	}
	return nil
}

// FindByID retrieves a test result by ID
func (r *TestResultRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintestplan.TestResult, error) {
	var row struct {
		ID         uuid.UUID `db:"id"`
		PlanID     uuid.UUID `db:"plan_id"`
		CaseID     uuid.UUID `db:"case_id"`
		ExecutorID uuid.UUID `db:"executor_id"`
		Status     string    `db:"status"`
		Note       string    `db:"note"`
		ExecutedAt time.Time `db:"executed_at"`
		CreatedAt  time.Time `db:"created_at"`
	}

	query := `
		SELECT id, plan_id, case_id, executor_id, status, note, executed_at, created_at
		FROM test_results
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaintestplan.ErrResultNotFound
		}
		return nil, fmt.Errorf("find test result by id: %w", err)
	}

	status, err := domaintestplan.ParseResultStatus(row.Status)
	if err != nil {
		return nil, fmt.Errorf("parse result status: %w", err)
	}

	return domaintestplan.ReconstructResult(
		row.ID,
		row.PlanID,
		row.CaseID,
		row.ExecutorID,
		status,
		row.Note,
		row.ExecutedAt,
		row.CreatedAt,
	), nil
}

// FindByPlanID retrieves all test results for a plan with pagination
func (r *TestResultRepository) FindByPlanID(ctx context.Context, planID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestResult, error) {
	query := `
		SELECT id, plan_id, case_id, executor_id, status, note, executed_at, created_at
		FROM test_results
		WHERE plan_id = $1
		ORDER BY executed_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []struct {
		ID         uuid.UUID `db:"id"`
		PlanID     uuid.UUID `db:"plan_id"`
		CaseID     uuid.UUID `db:"case_id"`
		ExecutorID uuid.UUID `db:"executor_id"`
		Status     string    `db:"status"`
		Note       string    `db:"note"`
		ExecutedAt time.Time `db:"executed_at"`
		CreatedAt  time.Time `db:"created_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, planID, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test results by plan id: %w", err)
	}

	return r.rowsToResults(rows)
}

// FindByCaseID retrieves all test results for a test case with pagination
func (r *TestResultRepository) FindByCaseID(ctx context.Context, caseID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestResult, error) {
	query := `
		SELECT id, plan_id, case_id, executor_id, status, note, executed_at, created_at
		FROM test_results
		WHERE case_id = $1
		ORDER BY executed_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []struct {
		ID         uuid.UUID `db:"id"`
		PlanID     uuid.UUID `db:"plan_id"`
		CaseID     uuid.UUID `db:"case_id"`
		ExecutorID uuid.UUID `db:"executor_id"`
		Status     string    `db:"status"`
		Note       string    `db:"note"`
		ExecutedAt time.Time `db:"executed_at"`
		CreatedAt  time.Time `db:"created_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, caseID, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test results by case id: %w", err)
	}

	return r.rowsToResults(rows)
}

// DeleteByPlanID removes all test results for a plan
func (r *TestResultRepository) DeleteByPlanID(ctx context.Context, planID uuid.UUID) error {
	query := `DELETE FROM test_results WHERE plan_id = $1`
	_, err := r.db.ExecContext(ctx, query, planID)
	if err != nil {
		return fmt.Errorf("delete test results by plan id: %w", err)
	}
	return nil
}

// CountByPlanID counts test results for a plan
func (r *TestResultRepository) CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM test_results WHERE plan_id = $1`
	err := r.db.GetContext(ctx, &count, query, planID)
	if err != nil {
		return 0, fmt.Errorf("count test results by plan id: %w", err)
	}
	return count, nil
}

// Helper functions
func (r *TestResultRepository) rowsToResults(rows []struct {
	ID         uuid.UUID `db:"id"`
	PlanID     uuid.UUID `db:"plan_id"`
	CaseID     uuid.UUID `db:"case_id"`
	ExecutorID uuid.UUID `db:"executor_id"`
	Status     string    `db:"status"`
	Note       string    `db:"note"`
	ExecutedAt time.Time `db:"executed_at"`
	CreatedAt  time.Time `db:"created_at"`
}) ([]*domaintestplan.TestResult, error) {
	results := make([]*domaintestplan.TestResult, 0, len(rows))
	for _, row := range rows {
		status, err := domaintestplan.ParseResultStatus(row.Status)
		if err != nil {
			return nil, fmt.Errorf("parse result status: %w", err)
		}

		result := domaintestplan.ReconstructResult(
			row.ID,
			row.PlanID,
			row.CaseID,
			row.ExecutorID,
			status,
			row.Note,
			row.ExecutedAt,
			row.CreatedAt,
		)
		results = append(results, result)
	}
	return results, nil
}
