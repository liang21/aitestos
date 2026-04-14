// Package testcase defines repository interfaces
package testcase

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TestCaseRepository defines the interface for test case persistence
type TestCaseRepository interface {
	// Save persists a new test case
	Save(ctx context.Context, tc *TestCase) error

	// FindByID retrieves a test case by ID
	FindByID(ctx context.Context, id uuid.UUID) (*TestCase, error)

	// FindByNumber retrieves a test case by case number
	FindByNumber(ctx context.Context, number CaseNumber) (*TestCase, error)

	// FindByModuleID retrieves all test cases for a module with pagination
	FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts QueryOptions) ([]*TestCase, error)

	// CountByModuleID counts total test cases for a module
	CountByModuleID(ctx context.Context, moduleID uuid.UUID) (int64, error)

	// FindByProjectID retrieves all test cases for a project with pagination
	FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*TestCase, error)

	// CountByProjectID counts total test cases for a project
	CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error)

	// Update updates an existing test case
	Update(ctx context.Context, tc *TestCase) error

	// Delete removes a test case
	Delete(ctx context.Context, id uuid.UUID) error

	// CountByDate counts test cases created on a specific date for a module
	CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error)
}

// QueryOptions holds pagination and filtering options
type QueryOptions struct {
	Offset   int
	Limit    int
	OrderBy  string
	Keywords string
	Status   CaseStatus
	Priority Priority
}
