// Package testcase provides test case repository implementation
package testcase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domaintestcase "github.com/liang21/aitestos/internal/domain/testcase"
)

// TestCaseRepository implements domaintestcase.TestCaseRepository interface
type TestCaseRepository struct {
	db *sqlx.DB
}

// NewTestCaseRepository creates a new test case repository
func NewTestCaseRepository(db *sqlx.DB) *TestCaseRepository {
	return &TestCaseRepository{db: db}
}

// Save persists a new test case
func (r *TestCaseRepository) Save(ctx context.Context, tc *domaintestcase.TestCase) error {
	query := `
		INSERT INTO test_cases (
			id, module_id, user_id, number, title, preconditions, steps,
			expected_result, case_type, priority, status, ai_metadata,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	preconditionsJSON, _ := toJSON(tc.Preconditions())
	stepsJSON, _ := toJSON(tc.Steps())
	expectedJSON, _ := toJSON(tc.ExpectedResult())
	var aiMetadataJSON []byte
	if tc.AiMetadata() != nil {
		aiMetadataJSON, _ = jsonMarshal(tc.AiMetadata())
	}

	_, err := r.db.ExecContext(ctx, query,
		tc.ID(),
		tc.ModuleID(),
		tc.UserID(),
		tc.Number().String(),
		tc.Title(),
		preconditionsJSON,
		stepsJSON,
		expectedJSON,
		string(tc.CaseType()),
		string(tc.Priority()),
		string(tc.Status()),
		aiMetadataJSON,
		tc.CreatedAt(),
		tc.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save test case: %w", err)
	}
	return nil
}

// FindByID retrieves a test case by ID
func (r *TestCaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintestcase.TestCase, error) {
	var row struct {
		ID            uuid.UUID `db:"id"`
		ModuleID      uuid.UUID `db:"module_id"`
		UserID        uuid.UUID `db:"user_id"`
		Number        string    `db:"number"`
		Title         string    `db:"title"`
		Preconditions string    `db:"preconditions"`
		Steps         string    `db:"steps"`
		Expected      string    `db:"expected_result"`
		CaseType      string    `db:"case_type"`
		Priority      string    `db:"priority"`
		Status        string    `db:"status"`
		AiMetadata    []byte    `db:"ai_metadata"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}

	query := `
		SELECT id, module_id, user_id, number, title, preconditions, steps,
			   expected_result, case_type, priority, status, ai_metadata,
			   created_at, updated_at
		FROM test_cases
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaintestcase.ErrCaseNotFound
		}
		return nil, fmt.Errorf("find test case by id: %w", err)
	}

	return r.rowToTestCase(&row)
}

// FindByNumber retrieves a test case by case number
func (r *TestCaseRepository) FindByNumber(ctx context.Context, number domaintestcase.CaseNumber) (*domaintestcase.TestCase, error) {
	var row struct {
		ID            uuid.UUID `db:"id"`
		ModuleID      uuid.UUID `db:"module_id"`
		UserID        uuid.UUID `db:"user_id"`
		Number        string    `db:"number"`
		Title         string    `db:"title"`
		Preconditions string    `db:"preconditions"`
		Steps         string    `db:"steps"`
		Expected      string    `db:"expected_result"`
		CaseType      string    `db:"case_type"`
		Priority      string    `db:"priority"`
		Status        string    `db:"status"`
		AiMetadata    []byte    `db:"ai_metadata"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}

	query := `
		SELECT id, module_id, user_id, number, title, preconditions, steps,
			   expected_result, case_type, priority, status, ai_metadata,
			   created_at, updated_at
		FROM test_cases
		WHERE number = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, number.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaintestcase.ErrCaseNotFound
		}
		return nil, fmt.Errorf("find test case by number: %w", err)
	}

	return r.rowToTestCase(&row)
}

// FindByModuleID retrieves all test cases for a module with pagination
func (r *TestCaseRepository) FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts domaintestcase.QueryOptions) ([]*domaintestcase.TestCase, error) {
	query := `
		SELECT id, module_id, user_id, number, title, preconditions, steps,
			   expected_result, case_type, priority, status, ai_metadata,
			   created_at, updated_at
		FROM test_cases
		WHERE module_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []struct {
		ID            uuid.UUID `db:"id"`
		ModuleID      uuid.UUID `db:"module_id"`
		UserID        uuid.UUID `db:"user_id"`
		Number        string    `db:"number"`
		Title         string    `db:"title"`
		Preconditions string    `db:"preconditions"`
		Steps         string    `db:"steps"`
		Expected      string    `db:"expected_result"`
		CaseType      string    `db:"case_type"`
		Priority      string    `db:"priority"`
		Status        string    `db:"status"`
		AiMetadata    []byte    `db:"ai_metadata"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, moduleID, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test cases by module id: %w", err)
	}

	return r.rowsToTestCases(rows)
}

// Update updates an existing test case
func (r *TestCaseRepository) Update(ctx context.Context, tc *domaintestcase.TestCase) error {
	preconditionsJSON, _ := toJSON(tc.Preconditions())
	stepsJSON, _ := toJSON(tc.Steps())
	expectedJSON, _ := toJSON(tc.ExpectedResult())

	query := `
		UPDATE test_cases
		SET title = $2, preconditions = $3, steps = $4, expected_result = $5,
			status = $6, updated_at = $7
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
		tc.ID(),
		tc.Title(),
		preconditionsJSON,
		stepsJSON,
		expectedJSON,
		string(tc.Status()),
		tc.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("update test case: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaintestcase.ErrCaseNotFound
	}
	return nil
}

// Delete removes a test case (soft delete)
func (r *TestCaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE test_cases SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete test case: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domaintestcase.ErrCaseNotFound
	}
	return nil
}

// CountByDate counts test cases created on a specific date for a module
func (r *TestCaseRepository) CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM test_cases
		WHERE module_id = $1
		  AND DATE(created_at) = DATE($2)
		  AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &count, query, moduleID, date)
	if err != nil {
		return 0, fmt.Errorf("count test cases by date: %w", err)
	}
	return count, nil
}

// Helper functions
func (r *TestCaseRepository) rowToTestCase(row any) (*domaintestcase.TestCase, error) {
	// Implementation depends on row structure
	return nil, fmt.Errorf("not implemented")
}

func (r *TestCaseRepository) rowsToTestCases(rows any) ([]*domaintestcase.TestCase, error) {
	return nil, fmt.Errorf("not implemented")
}

func toJSON(v any) ([]byte, error) {
	return jsonMarshal(v)
}

func jsonMarshal(v any) ([]byte, error) {
	// Simple implementation - in production use encoding/json
	return []byte("{}"), nil
}
