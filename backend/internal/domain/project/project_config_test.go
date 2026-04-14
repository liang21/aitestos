// Package project_test tests ProjectConfig entity
package project_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/project"
)

func TestNewProjectConfig(t *testing.T) {
	projectID := uuid.New()
	key := "llm.model"
	value := map[string]any{"name": "deepseek", "version": "1.0"}
	description := "LLM model configuration"

	tests := []struct {
		name        string
		projectID   uuid.UUID
		key         string
		value       map[string]any
		description string
		wantErr     bool
	}{
		{
			name:        "valid config",
			projectID:   projectID,
			key:         key,
			value:       value,
			description: description,
			wantErr:     false,
		},
		{
			name:        "valid config without description",
			projectID:   projectID,
			key:         "cache.ttl",
			value:       map[string]any{"seconds": 3600},
			description: "",
			wantErr:     false,
		},
		{
			name:        "empty key",
			projectID:   projectID,
			key:         "",
			value:       value,
			description: description,
			wantErr:     true,
		},
		{
			name:        "nil project ID",
			projectID:   uuid.Nil,
			key:         key,
			value:       value,
			description: description,
			wantErr:     true,
		},
		{
			name:        "nil value",
			projectID:   projectID,
			key:         key,
			value:       nil,
			description: description,
			wantErr:     true,
		},
		{
			name:        "empty value map",
			projectID:   projectID,
			key:         key,
			value:       map[string]any{},
			description: description,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := project.NewProjectConfig(tt.projectID, tt.key, tt.value, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewProjectConfig() returned nil config")
					return
				}
				if got.Key() != tt.key {
					t.Errorf("ProjectConfig.Key() = %v, want %v", got.Key(), tt.key)
				}
				if got.ProjectID() != tt.projectID {
					t.Errorf("ProjectConfig.ProjectID() = %v, want %v", got.ProjectID(), tt.projectID)
				}
				if got.Description() != tt.description {
					t.Errorf("ProjectConfig.Description() = %v, want %v", got.Description(), tt.description)
				}
			}
		})
	}
}

func TestProjectConfig_Accessors(t *testing.T) {
	projectID := uuid.New()
	key := "test.key"
	value := map[string]any{"field": "value", "number": 42}
	description := "test description"

	config, err := project.NewProjectConfig(projectID, key, value, description)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	if config.ID() == uuid.Nil {
		t.Error("ProjectConfig.ID() should not be nil")
	}
	if config.ProjectID() != projectID {
		t.Errorf("ProjectConfig.ProjectID() = %v, want %v", config.ProjectID(), projectID)
	}
	if config.Key() != key {
		t.Errorf("ProjectConfig.Key() = %v, want %v", config.Key(), key)
	}
	if config.Description() != description {
		t.Errorf("ProjectConfig.Description() = %v, want %v", config.Description(), description)
	}
	if config.CreatedAt().IsZero() {
		t.Error("ProjectConfig.CreatedAt() should not be zero")
	}
	if config.UpdatedAt().IsZero() {
		t.Error("ProjectConfig.UpdatedAt() should not be zero")
	}

	// Verify value map
	configValue := config.Value()
	if configValue == nil {
		t.Error("ProjectConfig.Value() should not be nil")
	}
	if configValue["field"] != "value" {
		t.Errorf("ProjectConfig.Value()[\"field\"] = %v, want value", configValue["field"])
	}
	if configValue["number"] != 42 {
		t.Errorf("ProjectConfig.Value()[\"number\"] = %v, want 42", configValue["number"])
	}
}

func TestProjectConfig_UpdateValue(t *testing.T) {
	projectID := uuid.New()
	key := "test.key"
	initialValue := map[string]any{"version": "1.0"}

	config, err := project.NewProjectConfig(projectID, key, initialValue, "")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	originalUpdatedAt := config.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	newValue := map[string]any{"version": "2.0", "new_field": "added"}
	err = config.UpdateValue(newValue)
	if err != nil {
		t.Errorf("ProjectConfig.UpdateValue() error = %v", err)
	}

	if config.Value()["version"] != "2.0" {
		t.Errorf("ProjectConfig.Value()[\"version\"] = %v, want 2.0", config.Value()["version"])
	}
	if config.Value()["new_field"] != "added" {
		t.Errorf("ProjectConfig.Value()[\"new_field\"] = %v, want added", config.Value()["new_field"])
	}
	if !config.UpdatedAt().After(originalUpdatedAt) {
		t.Error("ProjectConfig.UpdatedAt() should be updated after UpdateValue()")
	}
}

func TestProjectConfig_UpdateValue_NilValue(t *testing.T) {
	projectID := uuid.New()
	key := "test.key"
	initialValue := map[string]any{"version": "1.0"}

	config, err := project.NewProjectConfig(projectID, key, initialValue, "")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	err = config.UpdateValue(nil)
	if err == nil {
		t.Error("ProjectConfig.UpdateValue(nil) should return error")
	}
}

func TestProjectConfig_UpdateDescription(t *testing.T) {
	projectID := uuid.New()
	key := "test.key"
	value := map[string]any{"test": "data"}

	config, err := project.NewProjectConfig(projectID, key, value, "old description")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	originalUpdatedAt := config.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	newDescription := "updated description"
	err = config.UpdateDescription(newDescription)
	if err != nil {
		t.Errorf("ProjectConfig.UpdateDescription() error = %v", err)
	}

	if config.Description() != newDescription {
		t.Errorf("ProjectConfig.Description() = %v, want %v", config.Description(), newDescription)
	}
	if !config.UpdatedAt().After(originalUpdatedAt) {
		t.Error("ProjectConfig.UpdatedAt() should be updated after UpdateDescription()")
	}
}

func TestProjectConfig_Equal(t *testing.T) {
	projectID := uuid.New()
	key := "test.key"
	value := map[string]any{"test": "data"}

	config1, err := project.NewProjectConfig(projectID, key, value, "description")
	if err != nil {
		t.Fatalf("Failed to create config1: %v", err)
	}

	// Same projectID and key should be considered equal
	config2, err := project.NewProjectConfig(projectID, key, map[string]any{"other": "value"}, "other")
	if err != nil {
		t.Fatalf("Failed to create config2: %v", err)
	}

	tests := []struct {
		name     string
		config   *project.ProjectConfig
		other    *project.ProjectConfig
		expected bool
	}{
		{
			name:     "same projectID and key",
			config:   config1,
			other:    config2,
			expected: true,
		},
		{
			name:   "different projectID",
			config: config1,
			other: func() *project.ProjectConfig {
				c, _ := project.NewProjectConfig(uuid.New(), key, value, "")
				return c
			}(),
			expected: false,
		},
		{
			name:   "different key",
			config: config1,
			other: func() *project.ProjectConfig {
				c, _ := project.NewProjectConfig(projectID, "other.key", value, "")
				return c
			}(),
			expected: false,
		},
		{
			name:     "nil comparison",
			config:   config1,
			other:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.Equal(tt.other); got != tt.expected {
				t.Errorf("ProjectConfig.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}
