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

	// Update updates an existing test plan
	Update(ctx context.Context, plan *TestPlan) error

	// Delete removes a test plan
	Delete(ctx context.Context, id uuid.UUID) error
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

	// DeleteByPlanID removes all test results for a plan
	DeleteByPlanID(ctx context.Context, planID uuid.UUID) error

	// CountByPlanID counts test results for a plan
	CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error)
}

// QueryOptions holds pagination and filtering options
type QueryOptions struct {
	Offset   int
	Limit    int
	OrderBy  string
	Keywords string
	Status   PlanStatus
}
