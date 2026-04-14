// Package project_test tests ProjectRepository implementation
package project_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
)

// Test fixtures
func createTestProject(t *testing.T, name, prefix string) *domainproject.Project {
	t.Helper()
	project, err := domainproject.NewProject(name, prefix, "test description")
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}
	return project
}

func TestProjectRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new project", func(t *testing.T) {
		// Placeholder for integration test
		// Real implementation would use testcontainers
	})
}

func TestProjectRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing project", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent project", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectRepository_FindByName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by name", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectRepository_FindByPrefix(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by prefix", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectRepository_FindAll(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("list all projects with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update project", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("soft delete project", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockProjectRepository for testing without database
type MockProjectRepository struct {
	projects         map[uuid.UUID]*domainproject.Project
	projectsByName   map[string]*domainproject.Project
	projectsByPrefix map[string]*domainproject.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects:         make(map[uuid.UUID]*domainproject.Project),
		projectsByName:   make(map[string]*domainproject.Project),
		projectsByPrefix: make(map[string]*domainproject.Project),
	}
}

func (m *MockProjectRepository) Save(ctx context.Context, project *domainproject.Project) error {
	m.projects[project.ID()] = project
	m.projectsByName[project.Name()] = project
	m.projectsByPrefix[project.Prefix().String()] = project
	return nil
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainproject.Project, error) {
	project, ok := m.projects[id]
	if !ok {
		return nil, domainproject.ErrProjectNotFound
	}
	return project, nil
}

func (m *MockProjectRepository) FindByName(ctx context.Context, name string) (*domainproject.Project, error) {
	project, ok := m.projectsByName[name]
	if !ok {
		return nil, domainproject.ErrProjectNotFound
	}
	return project, nil
}

func (m *MockProjectRepository) FindByPrefix(ctx context.Context, prefix domainproject.ProjectPrefix) (*domainproject.Project, error) {
	project, ok := m.projectsByPrefix[prefix.String()]
	if !ok {
		return nil, domainproject.ErrProjectNotFound
	}
	return project, nil
}

func (m *MockProjectRepository) FindAll(ctx context.Context, opts domainproject.QueryOptions) ([]*domainproject.Project, error) {
	projects := make([]*domainproject.Project, 0, len(m.projects))
	for _, project := range m.projects {
		projects = append(projects, project)
	}
	return projects, nil
}

func (m *MockProjectRepository) Update(ctx context.Context, project *domainproject.Project) error {
	if _, ok := m.projects[project.ID()]; !ok {
		return domainproject.ErrProjectNotFound
	}
	m.projects[project.ID()] = project
	m.projectsByName[project.Name()] = project
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	project, ok := m.projects[id]
	if !ok {
		return domainproject.ErrProjectNotFound
	}
	delete(m.projects, id)
	delete(m.projectsByName, project.Name())
	delete(m.projectsByPrefix, project.Prefix().String())
	return nil
}

func TestMockProjectRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockProjectRepository()

	// Create
	project := createTestProject(t, "Test Project", "TP")
	err := repo.Save(ctx, project)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, project.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != project.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), project.ID())
	}

	// Read by Name
	found, err = repo.FindByName(ctx, project.Name())
	if err != nil {
		t.Fatalf("FindByName() error = %v", err)
	}
	if found.Name() != project.Name() {
		t.Errorf("FindByName().Name() = %v, want %v", found.Name(), project.Name())
	}

	// Read by Prefix
	found, err = repo.FindByPrefix(ctx, project.Prefix())
	if err != nil {
		t.Fatalf("FindByPrefix() error = %v", err)
	}
	if found.Prefix() != project.Prefix() {
		t.Errorf("FindByPrefix().Prefix() = %v, want %v", found.Prefix(), project.Prefix())
	}

	// Update
	project.UpdateDescription("updated description")
	err = repo.Update(ctx, project)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	found, _ = repo.FindByID(ctx, project.ID())
	if found.Description() != "updated description" {
		t.Errorf("Update().Description() = %v, want updated description", found.Description())
	}

	// Delete
	err = repo.Delete(ctx, project.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = repo.FindByID(ctx, project.ID())
	if err != domainproject.ErrProjectNotFound {
		t.Errorf("FindByID() after Delete() error = %v, want %v", err, domainproject.ErrProjectNotFound)
	}
}

func TestMockProjectRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockProjectRepository()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != domainproject.ErrProjectNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domainproject.ErrProjectNotFound)
	}

	_, err = repo.FindByName(ctx, "notfound")
	if err != domainproject.ErrProjectNotFound {
		t.Errorf("FindByName() error = %v, want %v", err, domainproject.ErrProjectNotFound)
	}

	prefix, _ := domainproject.ParseProjectPrefix("NF")
	_, err = repo.FindByPrefix(ctx, prefix)
	if err != domainproject.ErrProjectNotFound {
		t.Errorf("FindByPrefix() error = %v, want %v", err, domainproject.ErrProjectNotFound)
	}

	err = repo.Update(ctx, createTestProject(t, "test", "TS"))
	if err != domainproject.ErrProjectNotFound {
		t.Errorf("Update() error = %v, want %v", err, domainproject.ErrProjectNotFound)
	}

	err = repo.Delete(ctx, uuid.New())
	if err != domainproject.ErrProjectNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domainproject.ErrProjectNotFound)
	}
}
