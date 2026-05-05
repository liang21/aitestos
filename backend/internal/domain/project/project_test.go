// Package project_test tests Project aggregate
package project_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/liang21/aitestos/internal/domain/project"
)

func TestNewProject(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		description string
		wantErr     bool
	}{
		{
			name:        "valid project",
			projectName: "E-Commerce Platform",
			description: "E-commerce testing platform",
			wantErr:     false,
		},
		{
			name:        "empty name",
			projectName: "",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "empty description is allowed",
			projectName: "Test Project",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := project.NewProject(tt.projectName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewProject() returned nil project")
					return
				}
				if got.Name() != tt.projectName {
					t.Errorf("Project.Name() = %v, want %v", got.Name(), tt.projectName)
				}
				if got.Description() != tt.description {
					t.Errorf("Project.Description() = %v, want %v", got.Description(), tt.description)
				}
				if got.ID() == uuid.Nil {
					t.Error("Project.ID() should not be nil UUID")
				}
			}
		})
	}
}

func TestProject_Accessors(t *testing.T) {
	p, err := project.NewProject("Test Project", "Test Description")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Test ID accessor
	if p.ID() == uuid.Nil {
		t.Error("Project.ID() should not be nil")
	}

	// Test Name accessor
	if p.Name() != "Test Project" {
		t.Errorf("Project.Name() = %v, want Test Project", p.Name())
	}

	// Test Description accessor
	if p.Description() != "Test Description" {
		t.Errorf("Project.Description() = %v, want Test Description", p.Description())
	}

	// Test CreatedAt accessor
	if p.CreatedAt().IsZero() {
		t.Error("Project.CreatedAt() should not be zero")
	}

	// Test UpdatedAt accessor
	if p.UpdatedAt().IsZero() {
		t.Error("Project.UpdatedAt() should not be zero")
	}
}

func TestProject_UpdateDescription(t *testing.T) {
	p, err := project.NewProject("Test Project", "Original Description")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	originalUpdatedAt := p.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	newDesc := "Updated Description"
	p.UpdateDescription(newDesc)

	if p.Description() != newDesc {
		t.Errorf("Project.Description() = %v, want %v", p.Description(), newDesc)
	}
	if !p.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Project.UpdatedAt() should be updated")
	}
}

func TestProject_UpdateName(t *testing.T) {
	p, err := project.NewProject("Original Name", "Description")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	originalUpdatedAt := p.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	// Test valid name update
	err = p.UpdateName("New Name")
	if err != nil {
		t.Errorf("Project.UpdateName() error = %v", err)
	}
	if p.Name() != "New Name" {
		t.Errorf("Project.Name() = %v, want New Name", p.Name())
	}
	if !p.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Project.UpdatedAt() should be updated")
	}

	// Test empty name update should fail
	err = p.UpdateName("")
	if err == nil {
		t.Error("Project.UpdateName() should return error for empty name")
	}
}

func TestProject_MarshalJSON(t *testing.T) {
	p, err := project.NewProject("Test Project", "A test project")
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("failed to marshal project: %v", err)
	}

	t.Logf("JSON output: %s", string(data))

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify required fields
	if result["id"] == nil {
		t.Error("id field is missing")
	}
	if result["name"] == nil || result["name"] != "Test Project" {
		t.Errorf("name field is missing or wrong: got %v", result["name"])
	}
	if result["description"] == nil || result["description"] != "A test project" {
		t.Errorf("description field is missing or wrong: got %v", result["description"])
	}
	// Verify prefix field is NOT present
	if result["prefix"] != nil {
		t.Error("prefix field should not be present")
	}
}
