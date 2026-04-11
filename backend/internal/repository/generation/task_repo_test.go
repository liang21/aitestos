// Package generation_test tests GenerationTaskRepository implementation
package generation_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domaingeneration "github.com/liang21/aitestos/internal/domain/generation"
)

// Test fixtures
func createTestTask(t *testing.T, projectID uuid.UUID, prompt string) *domaingeneration.GenerationTask {
	t.Helper()
	task, err := domaingeneration.NewGenerationTask(projectID, uuid.New(), prompt, uuid.New())
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	return task
}

func TestGenerationTaskRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new task", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestGenerationTaskRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing task", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent task", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestGenerationTaskRepository_FindByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by project id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestGenerationTaskRepository_FindByStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by status with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestGenerationTaskRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update task", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockGenerationTaskRepository for testing without database
type MockGenerationTaskRepository struct {
	tasks          map[uuid.UUID]*domaingeneration.GenerationTask
	tasksByProject map[uuid.UUID][]*domaingeneration.GenerationTask
	tasksByStatus  map[domaingeneration.TaskStatus][]*domaingeneration.GenerationTask
}

func NewMockGenerationTaskRepository() *MockGenerationTaskRepository {
	return &MockGenerationTaskRepository{
		tasks:          make(map[uuid.UUID]*domaingeneration.GenerationTask),
		tasksByProject: make(map[uuid.UUID][]*domaingeneration.GenerationTask),
		tasksByStatus:  make(map[domaingeneration.TaskStatus][]*domaingeneration.GenerationTask),
	}
}

func (m *MockGenerationTaskRepository) Save(ctx context.Context, task *domaingeneration.GenerationTask) error {
	m.tasks[task.ID()] = task
	m.tasksByProject[task.ProjectID()] = append(m.tasksByProject[task.ProjectID()], task)
	m.tasksByStatus[task.Status()] = append(m.tasksByStatus[task.Status()], task)
	return nil
}

func (m *MockGenerationTaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaingeneration.GenerationTask, error) {
	task, ok := m.tasks[id]
	if !ok {
		return nil, domaingeneration.ErrTaskNotFound
	}
	return task, nil
}

func (m *MockGenerationTaskRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domaingeneration.QueryOptions) ([]*domaingeneration.GenerationTask, error) {
	tasks, ok := m.tasksByProject[projectID]
	if !ok {
		return []*domaingeneration.GenerationTask{}, nil
	}
	return tasks, nil
}

func (m *MockGenerationTaskRepository) FindByStatus(ctx context.Context, status domaingeneration.TaskStatus, opts domaingeneration.QueryOptions) ([]*domaingeneration.GenerationTask, error) {
	tasks, ok := m.tasksByStatus[status]
	if !ok {
		return []*domaingeneration.GenerationTask{}, nil
	}
	return tasks, nil
}

func (m *MockGenerationTaskRepository) Update(ctx context.Context, task *domaingeneration.GenerationTask) error {
	if _, ok := m.tasks[task.ID()]; !ok {
		return domaingeneration.ErrTaskNotFound
	}
	m.tasks[task.ID()] = task
	return nil
}

func TestMockGenerationTaskRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockGenerationTaskRepository()
	projectID := uuid.New()

	// Create
	task := createTestTask(t, projectID, "Generate test cases for login feature")
	err := repo.Save(ctx, task)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, task.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != task.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), task.ID())
	}

	// Read by ProjectID
	tasks, err := repo.FindByProjectID(ctx, projectID, domaingeneration.QueryOptions{})
	if err != nil {
		t.Fatalf("FindByProjectID() error = %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("FindByProjectID() returned %d tasks, want 1", len(tasks))
	}

	// Read by Status
	tasks, err = repo.FindByStatus(ctx, domaingeneration.TaskPending, domaingeneration.QueryOptions{})
	if err != nil {
		t.Fatalf("FindByStatus() error = %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("FindByStatus() returned %d tasks, want 1", len(tasks))
	}

	// Update - Start processing
	task.StartProcessing()
	err = repo.Update(ctx, task)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	found, _ = repo.FindByID(ctx, task.ID())
	if found.Status() != domaingeneration.TaskProcessing {
		t.Errorf("Update().Status() = %v, want processing", found.Status())
	}
}

func TestMockGenerationTaskRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockGenerationTaskRepository()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != domaingeneration.ErrTaskNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domaingeneration.ErrTaskNotFound)
	}

	task := createTestTask(t, uuid.New(), "test prompt")
	err = repo.Update(ctx, task)
	if err != domaingeneration.ErrTaskNotFound {
		t.Errorf("Update() error = %v, want %v", err, domaingeneration.ErrTaskNotFound)
	}
}
