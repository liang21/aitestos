// Package project_test tests Project aggregate
package project_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/liang21/aitestos/internal/domain/project"
)

func TestNewProject(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		prefix      string
		description string
		wantErr     bool
	}{
		{
			name:        "valid project",
			projectName: "E-Commerce Platform",
			prefix:      "ECO",
			description: "E-commerce testing platform",
			wantErr:     false,
		},
		{
			name:        "empty name",
			projectName: "",
			prefix:      "ECO",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "invalid prefix",
			projectName: "Test Project",
			prefix:      "abc",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "empty description is allowed",
			projectName: "Test Project",
			prefix:      "TST",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := project.NewProject(tt.projectName, tt.prefix, tt.description)
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
				if got.Prefix().String() != tt.prefix {
					t.Errorf("Project.Prefix() = %v, want %v", got.Prefix(), tt.prefix)
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
	p, err := project.NewProject("Test Project", "TST", "Test Description")
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

	// Test Prefix accessor
	if p.Prefix().String() != "TST" {
		t.Errorf("Project.Prefix() = %v, want TST", p.Prefix())
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
	p, err := project.NewProject("Test Project", "TST", "Original Description")
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
	p, err := project.NewProject("Original Name", "TST", "Description")
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
