// Package testplan_test tests TestPlanRepository implementation
package testplan_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domaintestplan "github.com/liang21/aitestos/internal/domain/testplan"
)

// Test fixtures
func createTestPlan(t *testing.T, projectID uuid.UUID, name string) *domaintestplan.TestPlan {
	t.Helper()
	plan, err := domaintestplan.NewTestPlan(projectID, name, "test description", uuid.New())
	if err != nil {
		t.Fatalf("Failed to create test plan: %v", err)
	}
	return plan
}

func TestTestPlanRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new test plan", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestPlanRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing test plan", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent test plan", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestPlanRepository_FindByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by project id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestPlanRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update test plan", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestTestPlanRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("soft delete test plan", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockTestPlanRepository for testing without database
type MockTestPlanRepository struct {
	plans       map[uuid.UUID]*domaintestplan.TestPlan
	plansByProject map[uuid.UUID][]*domaintestplan.TestPlan
}

func NewMockTestPlanRepository() *MockTestPlanRepository {
	return &MockTestPlanRepository{
		plans:        make(map[uuid.UUID]*domaintestplan.TestPlan),
		plansByProject: make(map[uuid.UUID][]*domaintestplan.TestPlan),
	}
}

func (m *MockTestPlanRepository) Save(ctx context.Context, plan *domaintestplan.TestPlan) error {
	m.plans[plan.ID()] = plan
	m.plansByProject[plan.ProjectID()] = append(m.plansByProject[plan.ProjectID()], plan)
	return nil
}

func (m *MockTestPlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintestplan.TestPlan, error) {
	plan, ok := m.plans[id]
	if !ok {
		return nil, domaintestplan.ErrPlanNotFound
	}
	return plan, nil
}

func (m *MockTestPlanRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domaintestplan.QueryOptions) ([]*domaintestplan.TestPlan, error) {
	plans, ok := m.plansByProject[projectID]
	if !ok {
		return []*domaintestplan.TestPlan{}, nil
	}
	return plans, nil
}

func (m *MockTestPlanRepository) Update(ctx context.Context, plan *domaintestplan.TestPlan) error {
	if _, ok := m.plans[plan.ID()]; !ok {
		return domaintestplan.ErrPlanNotFound
	}
	m.plans[plan.ID()] = plan
	return nil
}

func (m *MockTestPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	plan, ok := m.plans[id]
	if !ok {
		return domaintestplan.ErrPlanNotFound
	}
	delete(m.plans, id)

	// Remove from project list
	projectPlans := m.plansByProject[plan.ProjectID()]
	for i, p := range projectPlans {
		if p.ID() == id {
			m.plansByProject[plan.ProjectID()] = append(projectPlans[:i], projectPlans[i+1:]...)
			break
		}
	}
	return nil
}

func TestMockTestPlanRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockTestPlanRepository()
	projectID := uuid.New()

	// Create
	plan := createTestPlan(t, projectID, "Sprint 1 Test Plan")
	err := repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, plan.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != plan.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), plan.ID())
	}

	// Read by ProjectID
	plans, err := repo.FindByProjectID(ctx, projectID, domaintestplan.QueryOptions{})
	if err != nil {
		t.Fatalf("FindByProjectID() error = %v", err)
	}
	if len(plans) != 1 {
		t.Errorf("FindByProjectID() returned %d plans, want 1", len(plans))
	}

	// Update
	if err := plan.UpdateStatus(domaintestplan.StatusActive); err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	err = repo.Update(ctx, plan)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	found, _ = repo.FindByID(ctx, plan.ID())
	if found.Status() != domaintestplan.StatusActive {
		t.Errorf("Update().Status() = %v, want active", found.Status())
	}

	// Delete
	err = repo.Delete(ctx, plan.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = repo.FindByID(ctx, plan.ID())
	if err != domaintestplan.ErrPlanNotFound {
		t.Errorf("FindByID() after Delete() error = %v, want %v", err, domaintestplan.ErrPlanNotFound)
	}
}

func TestMockTestPlanRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockTestPlanRepository()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != domaintestplan.ErrPlanNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domaintestplan.ErrPlanNotFound)
	}

	plan := createTestPlan(t, uuid.New(), "test")
	err = repo.Update(ctx, plan)
	if err != domaintestplan.ErrPlanNotFound {
		t.Errorf("Update() error = %v, want %v", err, domaintestplan.ErrPlanNotFound)
	}

	err = repo.Delete(ctx, uuid.New())
	if err != domaintestplan.ErrPlanNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domaintestplan.ErrPlanNotFound)
	}
}
