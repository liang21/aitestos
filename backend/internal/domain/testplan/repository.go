// Package testplan defines repository interfaces
package testplan

import (
	"context"

	"github.com/google/uuid"
)

// TestPlanRepository defines the interface for test plan persistence
type TestPlanRepository interface {
	// Save persists a new test plan
	Save(ctx context.Context, plan *TestPlan) error

	// FindByID retrieves a test plan by ID
	FindByID(ctx context.Context, id uuid.UUID) (*TestPlan, error)

	// FindByProjectID retrieves all test plans for a project with pagination
	FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*TestPlan, error)

	// FindByStatus retrieves test plans by status with pagination
	FindByStatus(ctx context.Context, status PlanStatus, opts QueryOptions) ([]*TestPlan, error)

	// Update updates an existing test plan
	Update(ctx context.Context, plan *TestPlan) error

	// Delete removes a test plan
	Delete(ctx context.Context, id uuid.UUID) error

	// AddCase adds a test case to a plan
	AddCase(ctx context.Context, planID, caseID uuid.UUID) error

	// RemoveCase removes a test case from a plan
	RemoveCase(ctx context.Context, planID, caseID uuid.UUID) error

	// GetCaseIDs retrieves all case IDs for a plan
	GetCaseIDs(ctx context.Context, planID uuid.UUID) ([]uuid.UUID, error)

	// UpdateStatus updates the plan's status
	UpdateStatus(ctx context.Context, planID uuid.UUID, status PlanStatus) error
}

// TestResultRepository defines the interface for test result persistence
type TestResultRepository interface {
	// Save persists a test result
	Save(ctx context.Context, result *TestResult) error

	// FindByID retrieves a test result by ID
	FindByID(ctx context.Context, id uuid.UUID) (*TestResult, error)

	// FindByPlanID retrieves all test results for a plan with pagination
	FindByPlanID(ctx context.Context, planID uuid.UUID, opts QueryOptions) ([]*TestResult, error)

	// FindByCaseID retrieves all test results for a test case with pagination
	FindByCaseID(ctx context.Context, caseID uuid.UUID, opts QueryOptions) ([]*TestResult, error)

	// FindByExecutorID retrieves all test results by executor with pagination
	FindByExecutorID(ctx context.Context, executorID uuid.UUID, opts QueryOptions) ([]*TestResult, error)

	// FindLatestByCaseID retrieves the most recent test result for a case
	FindLatestByCaseID(ctx context.Context, caseID uuid.UUID) (*TestResult, error)

	// FindByPlanIDAndCaseID retrieves test results by plan and case ID
	FindByPlanIDAndCaseID(ctx context.Context, planID, caseID uuid.UUID) ([]*TestResult, error)

	// DeleteByPlanID removes all test results for a plan
	DeleteByPlanID(ctx context.Context, planID uuid.UUID) error

	// CountByPlanID counts test results for a plan
	CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error)

	// CountByStatus counts results grouped by status for a plan
	CountByStatus(ctx context.Context, planID uuid.UUID) (map[ResultStatus]int, error)
}

// QueryOptions holds pagination and filtering options
type QueryOptions struct {
	Offset   int
	Limit    int
	OrderBy  string
	Keywords string
	Status   ResultStatus
}
