// Package project_test tests ModuleRepository implementation
package project_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
)

// Test fixtures
func createTestModule(t *testing.T, projectID uuid.UUID, name, abbrev string) *domainproject.Module {
	t.Helper()
	module, err := domainproject.NewModule(projectID, name, "test description", uuid.New())
	if err != nil {
		t.Fatalf("Failed to create test module: %v", err)
	}
	return module
}

func TestModuleRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new module", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestModuleRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing module", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent module", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestModuleRepository_FindByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find modules by project id", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestModuleRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("soft delete module", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockModuleRepository for testing without database
type MockModuleRepository struct {
	modules          map[uuid.UUID]*domainproject.Module
	modulesByProject map[uuid.UUID][]*domainproject.Module
}

func NewMockModuleRepository() *MockModuleRepository {
	return &MockModuleRepository{
		modules:          make(map[uuid.UUID]*domainproject.Module),
		modulesByProject: make(map[uuid.UUID][]*domainproject.Module),
	}
}

func (m *MockModuleRepository) Save(ctx context.Context, module *domainproject.Module) error {
	m.modules[module.ID()] = module
	m.modulesByProject[module.ProjectID()] = append(m.modulesByProject[module.ProjectID()], module)
	return nil
}

func (m *MockModuleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainproject.Module, error) {
	module, ok := m.modules[id]
	if !ok {
		return nil, domainproject.ErrModuleNotFound
	}
	return module, nil
}

func (m *MockModuleRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domainproject.Module, error) {
	modules, ok := m.modulesByProject[projectID]
	if !ok {
		return []*domainproject.Module{}, nil
	}
	return modules, nil
}

func (m *MockModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	module, ok := m.modules[id]
	if !ok {
		return domainproject.ErrModuleNotFound
	}
	delete(m.modules, id)

	// Remove from project list
	projectModules := m.modulesByProject[module.ProjectID()]
	for i, mod := range projectModules {
		if mod.ID() == id {
			m.modulesByProject[module.ProjectID()] = append(projectModules[:i], projectModules[i+1:]...)
			break
		}
	}
	return nil
}

func TestMockModuleRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockModuleRepository()
	projectID := uuid.New()

	// Create
	module := createTestModule(t, projectID, "User Module", "USR")
	err := repo.Save(ctx, module)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, module.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != module.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), module.ID())
	}

	// Read by ProjectID
	modules, err := repo.FindByProjectID(ctx, projectID)
	if err != nil {
		t.Fatalf("FindByProjectID() error = %v", err)
	}
	if len(modules) != 1 {
		t.Errorf("FindByProjectID() returned %d modules, want 1", len(modules))
	}

}
