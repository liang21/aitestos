// Package project_test tests Module entity
package project_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/liang21/aitestos/internal/domain/project"
)

func TestNewModule(t *testing.T) {
	projectID := uuid.New()

	tests := []struct {
		name        string
		projectID   uuid.UUID
		moduleName  string
		description string
		wantErr     bool
	}{
		{
			name:        "valid module",
			projectID:   projectID,
			moduleName:  "User Management",
			description: "User management module",
			wantErr:     false,
		},
		{
			name:        "empty name",
			projectID:   projectID,
			moduleName:  "",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "nil project ID",
			projectID:   uuid.Nil,
			moduleName:  "Test Module",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "empty description is allowed",
			projectID:   projectID,
			moduleName:  "Test Module",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := project.NewModule(tt.projectID, tt.moduleName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewModule() returned nil module")
					return
				}
				if got.Name() != tt.moduleName {
					t.Errorf("Module.Name() = %v, want %v", got.Name(), tt.moduleName)
				}
				if got.ProjectID() != tt.projectID {
					t.Errorf("Module.ProjectID() = %v, want %v", got.ProjectID(), tt.projectID)
				}
			}
		})
	}
}

func TestModule_Accessors(t *testing.T) {
	projectID := uuid.New()

	m, err := project.NewModule(projectID, "User Module", "User management")
	if err != nil {
		t.Fatalf("Failed to create module: %v", err)
	}

	// Test ID accessor
	if m.ID() == uuid.Nil {
		t.Error("Module.ID() should not be nil")
	}

	// Test ProjectID accessor
	if m.ProjectID() != projectID {
		t.Errorf("Module.ProjectID() = %v, want %v", m.ProjectID(), projectID)
	}

	// Test Name accessor
	if m.Name() != "User Module" {
		t.Errorf("Module.Name() = %v, want User Module", m.Name())
	}

	// Test Description accessor
	if m.Description() != "User management" {
		t.Errorf("Module.Description() = %v, want User management", m.Description())
	}

	// Test CreatedAt accessor
	if m.CreatedAt().IsZero() {
		t.Error("Module.CreatedAt() should not be zero")
	}

	// Test UpdatedAt accessor
	if m.UpdatedAt().IsZero() {
		t.Error("Module.UpdatedAt() should not be zero")
	}
}

func TestModule_UpdateDescription(t *testing.T) {
	projectID := uuid.New()

	m, err := project.NewModule(projectID, "Test Module", "Original")
	if err != nil {
		t.Fatalf("Failed to create module: %v", err)
	}

	originalUpdatedAt := m.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	newDesc := "Updated Description"
	m.UpdateDescription(newDesc)

	if m.Description() != newDesc {
		t.Errorf("Module.Description() = %v, want %v", m.Description(), newDesc)
	}
	if !m.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Module.UpdatedAt() should be updated")
	}
}
