// Package testcase provides test case repository implementation
package testcase

import (
	"context"
	"database/sql"
	"encoding/json"
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

// caseRow maps SQL columns for test_case table
type caseRow struct {
	ID            uuid.UUID `db:"id"`
	ModuleID      uuid.UUID `db:"module_id"`
	UserID        uuid.UUID `db:"user_id"`
	Number        string    `db:"number"`
	Title         string    `db:"title"`
	Preconditions []byte    `db:"preconditions"`
	Steps         []byte    `db:"steps"`
	Expected      []byte    `db:"expected"`
	AiMetadata    []byte    `db:"ai_metadata"`
	CaseType      string    `db:"case_type"`
	Priority      string    `db:"priority"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// toTestCase converts a database row to domain TestCase
func (r *TestCaseRepository) toTestCase(row *caseRow) (*domaintestcase.TestCase, error) {
	number, err := domaintestcase.ParseCaseNumber(row.Number)
	if err != nil {
		return nil, fmt.Errorf("parse case number: %w", err)
	}

	var preconditions domaintestcase.Preconditions
	if err := json.Unmarshal(row.Preconditions, &preconditions); err != nil {
		return nil, fmt.Errorf("unmarshal preconditions: %w", err)
	}

	var steps domaintestcase.Steps
	if err := json.Unmarshal(row.Steps, &steps); err != nil {
		return nil, fmt.Errorf("unmarshal steps: %w", err)
	}

	var expected domaintestcase.ExpectedResult
	if err := json.Unmarshal(row.Expected, &expected); err != nil {
		return nil, fmt.Errorf("unmarshal expected result: %w", err)
	}

	var aiMetadata *domaintestcase.AiMetadata
	if len(row.AiMetadata) > 0 && string(row.AiMetadata) != "{}" {
		aiMetadata = &domaintestcase.AiMetadata{}
		if err := json.Unmarshal(row.AiMetadata, aiMetadata); err != nil {
			return nil, fmt.Errorf("unmarshal ai metadata: %w", err)
		}
	}

	return domaintestcase.ReconstructTestCase(
		row.ID,
		row.ModuleID,
		row.UserID,
		number,
		row.Title,
		preconditions,
		steps,
		expected,
		aiMetadata,
		domaintestcase.CaseType(row.CaseType),
		domaintestcase.Priority(row.Priority),
		domaintestcase.CaseStatus(row.Status),
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

const caseColumns = `id, module_id, user_id, number, title, preconditions, steps, expected, ai_metadata, case_type, priority, status, created_at, updated_at`

// Save persists a new test case
func (r *TestCaseRepository) Save(ctx context.Context, tc *domaintestcase.TestCase) error {
	preconditionsJSON, err := json.Marshal(tc.Preconditions())
	if err != nil {
		return fmt.Errorf("marshal preconditions: %w", err)
	}

	stepsJSON, err := json.Marshal(tc.Steps())
	if err != nil {
		return fmt.Errorf("marshal steps: %w", err)
	}

	expectedJSON, err := json.Marshal(tc.ExpectedResult())
	if err != nil {
		return fmt.Errorf("marshal expected result: %w", err)
	}

	var aiMetadataJSON []byte
	if tc.AiMetadata() != nil {
		aiMetadataJSON, err = json.Marshal(tc.AiMetadata())
		if err != nil {
			return fmt.Errorf("marshal ai metadata: %w", err)
		}
	}

	query := `
		INSERT INTO test_case (id, module_id, user_id, number, title, preconditions, steps, expected, ai_metadata, case_type, priority, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err = r.db.ExecContext(ctx, query,
		tc.ID(),
		tc.ModuleID(),
		tc.UserID(),
		tc.Number().String(),
		tc.Title(),
		preconditionsJSON,
		stepsJSON,
		expectedJSON,
		aiMetadataJSON,
		string(tc.CaseType()),
		string(tc.Priority()),
		string(tc.Status()),
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
	var row caseRow
	query := fmt.Sprintf(`SELECT %s FROM test_case WHERE id = $1 AND deleted_at IS NULL`, caseColumns)
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaintestcase.ErrCaseNotFound
		}
		return nil, fmt.Errorf("find test case by id: %w", err)
	}
	return r.toTestCase(&row)
}

// FindByNumber retrieves a test case by case number
func (r *TestCaseRepository) FindByNumber(ctx context.Context, number domaintestcase.CaseNumber) (*domaintestcase.TestCase, error) {
	var row caseRow
	query := fmt.Sprintf(`SELECT %s FROM test_case WHERE number = $1 AND deleted_at IS NULL`, caseColumns)
	err := r.db.GetContext(ctx, &row, query, number.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domaintestcase.ErrCaseNotFound
		}
		return nil, fmt.Errorf("find test case by number: %w", err)
	}
	return r.toTestCase(&row)
}

// FindByModuleID retrieves all test cases for a module with pagination
func (r *TestCaseRepository) FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts domaintestcase.QueryOptions) ([]*domaintestcase.TestCase, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	var rows []caseRow
	query := fmt.Sprintf(`SELECT %s FROM test_case WHERE module_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, caseColumns)
	if err := r.db.SelectContext(ctx, &rows, query, moduleID, limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test cases by module id: %w", err)
	}

	return r.toTestCases(rows)
}

// CountByModuleID counts total test cases for a module
func (r *TestCaseRepository) CountByModuleID(ctx context.Context, moduleID uuid.UUID) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM test_case WHERE module_id = $1 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &count, query, moduleID)
	if err != nil {
		return 0, fmt.Errorf("count test cases by module id: %w", err)
	}
	return count, nil
}

// FindByProjectID retrieves all test cases for a project with pagination
// Joins module table to filter by project_id
func (r *TestCaseRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domaintestcase.QueryOptions) ([]*domaintestcase.TestCase, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	var rows []caseRow
	query := `
		SELECT tc.id, tc.module_id, tc.user_id, tc.number, tc.title,
		       tc.preconditions, tc.steps, tc.expected, tc.ai_metadata,
		       tc.case_type, tc.priority, tc.status, tc.created_at, tc.updated_at
		FROM test_case tc
		INNER JOIN module m ON tc.module_id = m.id
		WHERE m.project_id = $1 AND tc.deleted_at IS NULL
		ORDER BY tc.created_at DESC
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &rows, query, projectID, limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find test cases by project id: %w", err)
	}

	return r.toTestCases(rows)
}

// CountByProjectID counts total test cases for a project
func (r *TestCaseRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM test_case tc
		INNER JOIN module m ON tc.module_id = m.id
		WHERE m.project_id = $1 AND tc.deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &count, query, projectID)
	if err != nil {
		return 0, fmt.Errorf("count test cases by project id: %w", err)
	}
	return count, nil
}

// Update updates an existing test case
func (r *TestCaseRepository) Update(ctx context.Context, tc *domaintestcase.TestCase) error {
	preconditionsJSON, err := json.Marshal(tc.Preconditions())
	if err != nil {
		return fmt.Errorf("marshal preconditions: %w", err)
	}

	stepsJSON, err := json.Marshal(tc.Steps())
	if err != nil {
		return fmt.Errorf("marshal steps: %w", err)
	}

	expectedJSON, err := json.Marshal(tc.ExpectedResult())
	if err != nil {
		return fmt.Errorf("marshal expected result: %w", err)
	}

	query := `
		UPDATE test_case
		SET title = $2, preconditions = $3, steps = $4, expected = $5,
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
	query := `UPDATE test_case SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
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
		FROM test_case
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

// toTestCases converts multiple rows to domain TestCases
func (r *TestCaseRepository) toTestCases(rows []caseRow) ([]*domaintestcase.TestCase, error) {
	cases := make([]*domaintestcase.TestCase, 0, len(rows))
	for i := range rows {
		tc, err := r.toTestCase(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}
