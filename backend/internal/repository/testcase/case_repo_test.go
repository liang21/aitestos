// Package testcase_test tests TestCaseRepository implementation
package testcase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domaintestcase "github.com/liang21/aitestos/internal/domain/testcase"
)

// Test fixtures
func createTestCase(t *testing.T, moduleID uuid.UUID, number string) *domaintestcase.TestCase {
	t.Helper()
	caseNumber, err := domaintestcase.ParseCaseNumber(number)
	if err != nil {
		t.Fatalf("Failed to parse case number: %v", err)
	}
	tc, err := domaintestcase.NewTestCase(
		moduleID,
		uuid.New(),
		caseNumber,
		"Test Case Title",
		domaintestcase.Preconditions{"User is logged in"},
		domaintestcase.Steps{"Step 1", "Step 2"},
		domaintestcase.ExpectedResult{"status": "success"},
		domaintestcase.CaseTypeFunctionality,
		domaintestcase.PriorityP0,
	)
	if err != nil {
		t.Fatalf("Failed to create test case: %v", err)
	}
	return tc
}

func TestTestCaseRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new test case", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestCaseRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing test case", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent test case", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestCaseRepository_FindByNumber(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by case number", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestCaseRepository_FindByModuleID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by module id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestCaseRepository_FindByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by project id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestCaseRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update test case", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestCaseRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("soft delete test case", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestCaseRepository_CountByDate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("count test cases by date", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockTestCaseRepository for testing without database
type MockTestCaseRepository struct {
	cases       map[uuid.UUID]*domaintestcase.TestCase
	casesByNumber map[string]*domaintestcase.TestCase
	casesByModule map[uuid.UUID][]*domaintestcase.TestCase
}

func NewMockTestCaseRepository() *MockTestCaseRepository {
	return &MockTestCaseRepository{
		cases:        make(map[uuid.UUID]*domaintestcase.TestCase),
		casesByNumber: make(map[string]*domaintestcase.TestCase),
		casesByModule: make(map[uuid.UUID][]*domaintestcase.TestCase),
	}
}

func (m *MockTestCaseRepository) Save(ctx context.Context, tc *domaintestcase.TestCase) error {
	m.cases[tc.ID()] = tc
	m.casesByNumber[tc.Number().String()] = tc
	m.casesByModule[tc.ModuleID()] = append(m.casesByModule[tc.ModuleID()], tc)
	return nil
}

func (m *MockTestCaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintestcase.TestCase, error) {
	tc, ok := m.cases[id]
	if !ok {
		return nil, domaintestcase.ErrCaseNotFound
	}
	return tc, nil
}

func (m *MockTestCaseRepository) FindByNumber(ctx context.Context, number domaintestcase.CaseNumber) (*domaintestcase.TestCase, error) {
	tc, ok := m.casesByNumber[number.String()]
	if !ok {
		return nil, domaintestcase.ErrCaseNotFound
	}
	return tc, nil
}

func (m *MockTestCaseRepository) FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts domaintestcase.QueryOptions) ([]*domaintestcase.TestCase, error) {
	cases, ok := m.casesByModule[moduleID]
	if !ok {
		return []*domaintestcase.TestCase{}, nil
	}
	return cases, nil
}

func (m *MockTestCaseRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domaintestcase.QueryOptions) ([]*domaintestcase.TestCase, error) {
	// Placeholder - would need project->module relationship
	return []*domaintestcase.TestCase{}, nil
}

func (m *MockTestCaseRepository) Update(ctx context.Context, tc *domaintestcase.TestCase) error {
	if _, ok := m.cases[tc.ID()]; !ok {
		return domaintestcase.ErrCaseNotFound
	}
	m.cases[tc.ID()] = tc
	return nil
}

func (m *MockTestCaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tc, ok := m.cases[id]
	if !ok {
		return domaintestcase.ErrCaseNotFound
	}
	delete(m.cases, id)
	delete(m.casesByNumber, tc.Number().String())
	return nil
}

func (m *MockTestCaseRepository) CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error) {
	cases := m.casesByModule[moduleID]
	var count int64
	for _, tc := range cases {
		if tc.CreatedAt().Format("2006-01-02") == date.Format("2006-01-02") {
			count++
		}
	}
	return count, nil
}

func TestMockTestCaseRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockTestCaseRepository()
	moduleID := uuid.New()

	// Create
	tc := createTestCase(t, moduleID, "ECO-USR-20260403-001")
	err := repo.Save(ctx, tc)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, tc.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != tc.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), tc.ID())
	}

	// Read by Number
	found, err = repo.FindByNumber(ctx, tc.Number())
	if err != nil {
		t.Fatalf("FindByNumber() error = %v", err)
	}
	if found.Number() != tc.Number() {
		t.Errorf("FindByNumber().Number() = %v, want %v", found.Number(), tc.Number())
	}

	// Read by ModuleID
	cases, err := repo.FindByModuleID(ctx, moduleID, domaintestcase.QueryOptions{})
	if err != nil {
		t.Fatalf("FindByModuleID() error = %v", err)
	}
	if len(cases) != 1 {
		t.Errorf("FindByModuleID() returned %d cases, want 1", len(cases))
	}

	// Update
	tc.UpdateStatus(domaintestcase.StatusPass)
	err = repo.Update(ctx, tc)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	found, _ = repo.FindByID(ctx, tc.ID())
	if found.Status() != domaintestcase.StatusPass {
		t.Errorf("Update().Status() = %v, want pass", found.Status())
	}

	// Count by Date
	count, err := repo.CountByDate(ctx, moduleID, time.Now())
	if err != nil {
		t.Fatalf("CountByDate() error = %v", err)
	}
	if count != 1 {
		t.Errorf("CountByDate() = %d, want 1", count)
	}

	// Delete
	err = repo.Delete(ctx, tc.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = repo.FindByID(ctx, tc.ID())
	if err != domaintestcase.ErrCaseNotFound {
		t.Errorf("FindByID() after Delete() error = %v, want %v", err, domaintestcase.ErrCaseNotFound)
	}
}

func TestMockTestCaseRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockTestCaseRepository()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != domaintestcase.ErrCaseNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domaintestcase.ErrCaseNotFound)
	}

	number, _ := domaintestcase.ParseCaseNumber("ECO-USR-20260403-001")
	_, err = repo.FindByNumber(ctx, number)
	if err != domaintestcase.ErrCaseNotFound {
		t.Errorf("FindByNumber() error = %v, want %v", err, domaintestcase.ErrCaseNotFound)
	}

	tc := createTestCase(t, uuid.New(), "ECO-USR-20260403-002")
	err = repo.Update(ctx, tc)
	if err != domaintestcase.ErrCaseNotFound {
		t.Errorf("Update() error = %v, want %v", err, domaintestcase.ErrCaseNotFound)
	}

	err = repo.Delete(ctx, uuid.New())
	if err != domaintestcase.ErrCaseNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domaintestcase.ErrCaseNotFound)
	}
}
