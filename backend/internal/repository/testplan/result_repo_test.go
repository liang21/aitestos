// Package testplan_test tests TestResultRepository implementation
package testplan_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domaintestplan "github.com/liang21/aitestos/internal/domain/testplan"
)

// Test fixtures
func createTestResult(t *testing.T, planID, caseID uuid.UUID, status domaintestplan.ResultStatus) *domaintestplan.TestResult {
	t.Helper()
	result, err := domaintestplan.NewTestResult(planID, caseID, uuid.New(), status, "test note")
	if err != nil {
		t.Fatalf("Failed to create test result: %v", err)
	}
	return result
}

func TestTestResultRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new test result", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestResultRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing test result", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent test result", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestResultRepository_FindByPlanID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by plan id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestResultRepository_FindByCaseID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by case id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestResultRepository_FindByExecutorID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by executor id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestResultRepository_DeleteByPlanID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("delete all results for a plan", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestResultRepository_CountByPlanID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("count results for a plan", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockTestResultRepository for testing without database
type MockTestResultRepository struct {
	results       map[uuid.UUID]*domaintestplan.TestResult
	resultsByPlan map[uuid.UUID][]*domaintestplan.TestResult
	resultsByCase map[uuid.UUID][]*domaintestplan.TestResult
}

func NewMockTestResultRepository() *MockTestResultRepository {
	return &MockTestResultRepository{
		results:       make(map[uuid.UUID]*domaintestplan.TestResult),
		resultsByPlan: make(map[uuid.UUID][]*domaintestplan.TestResult),
		resultsByCase: make(map[uuid.UUID][]*domaintestplan.TestResult),
	}
}

func (m *MockTestResultRepository) Save(ctx context.Context, result *domaintestplan.TestResult) error {
	m.results[result.ID()] = result
	m.resultsByPlan[result.PlanID()] = append(m.resultsByPlan[result.PlanID()], result)
	m.resultsByCase[result.CaseID()] = append(m.resultsByCase[result.CaseID()], result)
	return nil
}

func (m *MockTestResultRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintestplan.TestResult, error) {
	result, ok := m.results[id]
	if !ok {
		return nil, domaintestplan.ErrResultNotFound
	}
	return result, nil
}

func (m *MockTestResultRepository) FindByPlanID(ctx context.Context, planID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestResult, error) {
	results, ok := m.resultsByPlan[planID]
	if !ok {
		return []*domaintestplan.TestResult{}, nil
	}
	return results, nil
}

func (m *MockTestResultRepository) FindByCaseID(ctx context.Context, caseID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestResult, error) {
	results, ok := m.resultsByCase[caseID]
	if !ok {
		return []*domaintestplan.TestResult{}, nil
	}
	return results, nil
}

func (m *MockTestResultRepository) FindByExecutorID(ctx context.Context, executorID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestResult, error) {
	var results []*domaintestplan.TestResult
	for _, result := range m.results {
		if result.ExecutedBy() == executorID {
			results = append(results, result)
		}
	}
	return results, nil
}

func (m *MockTestResultRepository) DeleteByPlanID(ctx context.Context, planID uuid.UUID) error {
	results := m.resultsByPlan[planID]
	for _, result := range results {
		delete(m.results, result.ID())
		// Remove from case list
		caseResults := m.resultsByCase[result.CaseID()]
		for i, r := range caseResults {
			if r.ID() == result.ID() {
				m.resultsByCase[result.CaseID()] = append(caseResults[:i], caseResults[i+1:]...)
				break
			}
		}
	}
	delete(m.resultsByPlan, planID)
	return nil
}

func (m *MockTestResultRepository) CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error) {
	return int64(len(m.resultsByPlan[planID])), nil
}

func TestMockTestResultRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockTestResultRepository()
	planID := uuid.New()
	caseID := uuid.New()

	// Create
	result := createTestResult(t, planID, caseID, domaintestplan.ResultPass)
	err := repo.Save(ctx, result)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, result.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != result.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), result.ID())
	}

	// Read by PlanID
	results, err := repo.FindByPlanID(ctx, planID, domaintestplan.QueryOptions{})
	if err != nil {
		t.Fatalf("FindByPlanID() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("FindByPlanID() returned %d results, want 1", len(results))
	}

	// Read by CaseID
	results, err = repo.FindByCaseID(ctx, caseID, domaintestplan.QueryOptions{})
	if err != nil {
		t.Fatalf("FindByCaseID() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("FindByCaseID() returned %d results, want 1", len(results))
	}

	// Count by PlanID
	count, err := repo.CountByPlanID(ctx, planID)
	if err != nil {
		t.Fatalf("CountByPlanID() error = %v", err)
	}
	if count != 1 {
		t.Errorf("CountByPlanID() = %d, want 1", count)
	}

	// Delete by PlanID
	err = repo.DeleteByPlanID(ctx, planID)
	if err != nil {
		t.Fatalf("DeleteByPlanID() error = %v", err)
	}
	_, err = repo.FindByID(ctx, result.ID())
	if err != domaintestplan.ErrResultNotFound {
		t.Errorf("FindByID() after DeleteByPlanID() error = %v, want %v", err, domaintestplan.ErrResultNotFound)
	}
}

func TestMockTestResultRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockTestResultRepository()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != domaintestplan.ErrResultNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domaintestplan.ErrResultNotFound)
	}
}
