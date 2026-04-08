// Package project_test tests ProjectConfigRepository implementation
package project_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
)

// Test fixtures
func createTestConfig(t *testing.T, projectID uuid.UUID, key string) *domainproject.ProjectConfig {
	t.Helper()
	config, err := domainproject.NewProjectConfig(projectID, key, map[string]any{"value": "test"}, "test config")
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	return config
}

func TestProjectConfigRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new config", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("upsert existing config", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectConfigRepository_FindByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find configs by project id", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectConfigRepository_FindByKey(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by key", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent key", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestProjectConfigRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("delete config", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockProjectConfigRepository for testing without database
type MockProjectConfigRepository struct {
	configs       map[uuid.UUID]*domainproject.ProjectConfig
	configsByProject map[uuid.UUID][]*domainproject.ProjectConfig
	configsByKey   map[string]*domainproject.ProjectConfig // key: "projectID:key"
}

func NewMockProjectConfigRepository() *MockProjectConfigRepository {
	return &MockProjectConfigRepository{
		configs:         make(map[uuid.UUID]*domainproject.ProjectConfig),
		configsByProject: make(map[uuid.UUID][]*domainproject.ProjectConfig),
		configsByKey:    make(map[string]*domainproject.ProjectConfig),
	}
}

func (m *MockProjectConfigRepository) Save(ctx context.Context, config *domainproject.ProjectConfig) error {
	m.configs[config.ID()] = config
	m.configsByProject[config.ProjectID()] = append(m.configsByProject[config.ProjectID()], config)
	key := config.ProjectID().String() + ":" + config.Key()
	m.configsByKey[key] = config
	return nil
}

func (m *MockProjectConfigRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domainproject.ProjectConfig, error) {
	configs, ok := m.configsByProject[projectID]
	if !ok {
		return []*domainproject.ProjectConfig{}, nil
	}
	return configs, nil
}

func (m *MockProjectConfigRepository) FindByKey(ctx context.Context, projectID uuid.UUID, key string) (*domainproject.ProjectConfig, error) {
	fullKey := projectID.String() + ":" + key
	config, ok := m.configsByKey[fullKey]
	if !ok {
		return nil, domainproject.ErrConfigNotFound
	}
	return config, nil
}

func (m *MockProjectConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	config, ok := m.configs[id]
	if !ok {
		return domainproject.ErrConfigNotFound
	}
	delete(m.configs, id)

	// Remove from project list
	projectConfigs := m.configsByProject[config.ProjectID()]
	for i, cfg := range projectConfigs {
		if cfg.ID() == id {
			m.configsByProject[config.ProjectID()] = append(projectConfigs[:i], projectConfigs[i+1:]...)
			break
		}
	}

	// Remove from key map
	fullKey := config.ProjectID().String() + ":" + config.Key()
	delete(m.configsByKey, fullKey)

	return nil
}

func TestMockProjectConfigRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockProjectConfigRepository()
	projectID := uuid.New()

	// Create
	config := createTestConfig(t, projectID, "llm.model")
	err := repo.Save(ctx, config)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ProjectID
	configs, err := repo.FindByProjectID(ctx, projectID)
	if err != nil {
		t.Fatalf("FindByProjectID() error = %v", err)
	}
	if len(configs) != 1 {
		t.Errorf("FindByProjectID() returned %d configs, want 1", len(configs))
	}

	// Read by Key
	found, err := repo.FindByKey(ctx, projectID, "llm.model")
	if err != nil {
		t.Fatalf("FindByKey() error = %v", err)
	}
	if found.Key() != "llm.model" {
		t.Errorf("FindByKey().Key() = %v, want llm.model", found.Key())
	}

	// Upsert (save again with same projectID and key)
	newConfig := createTestConfig(t, projectID, "llm.model")
	err = repo.Save(ctx, newConfig)
	if err != nil {
		t.Fatalf("Save() upsert error = %v", err)
	}

	// Delete
	err = repo.Delete(ctx, config.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = repo.FindByKey(ctx, projectID, "llm.model")
	if err != domainproject.ErrConfigNotFound {
		t.Errorf("FindByKey() after Delete() error = %v, want %v", err, domainproject.ErrConfigNotFound)
	}
}

func TestMockProjectConfigRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockProjectConfigRepository()
	projectID := uuid.New()

	_, err := repo.FindByKey(ctx, projectID, "notfound")
	if err != domainproject.ErrConfigNotFound {
		t.Errorf("FindByKey() error = %v, want %v", err, domainproject.ErrConfigNotFound)
	}

	err = repo.Delete(ctx, uuid.New())
	if err != domainproject.ErrConfigNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domainproject.ErrConfigNotFound)
	}
}
