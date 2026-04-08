// Package testplan_test tests TestPlan aggregate
package testplan_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testplan"
)

func TestNewTestPlan(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name        string
		projectID   uuid.UUID
		planName    string
		description string
		userID      uuid.UUID
		wantErr     bool
	}{
		{
			name:        "valid test plan",
			projectID:   projectID,
			planName:    "Sprint 1 Test Plan",
			description: "Test plan for sprint 1",
			userID:      userID,
			wantErr:     false,
		},
		{
			name:        "empty name",
			projectID:   projectID,
			planName:    "",
			description: "Description",
			userID:      userID,
			wantErr:     true,
		},
		{
			name:        "nil project ID",
			projectID:   uuid.Nil,
			planName:    "Test Plan",
			description: "Description",
			userID:      userID,
			wantErr:     true,
		},
		{
			name:        "nil user ID",
			projectID:   projectID,
			planName:    "Test Plan",
			description: "Description",
			userID:      uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "empty description is allowed",
			projectID:   projectID,
			planName:    "Test Plan",
			description: "",
			userID:      userID,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testplan.NewTestPlan(tt.projectID, tt.planName, tt.description, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewTestPlan() returned nil plan")
					return
				}
				if got.Name() != tt.planName {
					t.Errorf("TestPlan.Name() = %v, want %v", got.Name(), tt.planName)
				}
				if got.ProjectID() != tt.projectID {
					t.Errorf("TestPlan.ProjectID() = %v, want %v", got.ProjectID(), tt.projectID)
				}
				if got.Status() != testplan.StatusDraft {
					t.Errorf("TestPlan.Status() = %v, want draft", got.Status())
				}
			}
		})
	}
}

func TestTestPlan_Accessors(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	tp, err := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	if err != nil {
		t.Fatalf("Failed to create test plan: %v", err)
	}

	if tp.ID() == uuid.Nil {
		t.Error("TestPlan.ID() should not be nil")
	}
	if tp.ProjectID() != projectID {
		t.Errorf("TestPlan.ProjectID() = %v, want %v", tp.ProjectID(), projectID)
	}
	if tp.Name() != "Test Plan" {
		t.Errorf("TestPlan.Name() = %v, want Test Plan", tp.Name())
	}
	if tp.Description() != "Description" {
		t.Errorf("TestPlan.Description() = %v, want Description", tp.Description())
	}
	if tp.Status() != testplan.StatusDraft {
		t.Errorf("TestPlan.Status() = %v, want draft", tp.Status())
	}
	if tp.CreatedAt().IsZero() {
		t.Error("TestPlan.CreatedAt() should not be zero")
	}
	if tp.UpdatedAt().IsZero() {
		t.Error("TestPlan.UpdatedAt() should not be zero")
	}
	if tp.CreatedBy() != userID {
		t.Errorf("TestPlan.CreatedBy() = %v, want %v", tp.CreatedBy(), userID)
	}
}

func TestTestPlan_AddCase(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	tp, err := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	if err != nil {
		t.Fatalf("Failed to create test plan: %v", err)
	}

	caseID := uuid.New()
	err = tp.AddCase(caseID)
	if err != nil {
		t.Errorf("TestPlan.AddCase() error = %v", err)
	}
	if !tp.HasCase(caseID) {
		t.Error("TestPlan.HasCase() should return true for added case")
	}

	// Test adding duplicate case should fail
	err = tp.AddCase(caseID)
	if err == nil {
		t.Error("TestPlan.AddCase() should return error for duplicate case")
	}
}

func TestTestPlan_RemoveCase(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	tp, err := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	if err != nil {
		t.Fatalf("Failed to create test plan: %v", err)
	}

	caseID := uuid.New()
	tp.AddCase(caseID)
	tp.RemoveCase(caseID)

	if tp.HasCase(caseID) {
		t.Error("TestPlan.HasCase() should return false for removed case")
	}
}

func TestTestPlan_UpdateStatus(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	tp, err := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	if err != nil {
		t.Fatalf("Failed to create test plan: %v", err)
	}

	originalUpdatedAt := tp.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	// Test valid transition: draft -> active
	err = tp.UpdateStatus(testplan.StatusActive)
	if err != nil {
		t.Errorf("TestPlan.UpdateStatus() error = %v", err)
	}
	if tp.Status() != testplan.StatusActive {
		t.Errorf("TestPlan.Status() = %v, want active", tp.Status())
	}
	if !tp.UpdatedAt().After(originalUpdatedAt) {
		t.Error("TestPlan.UpdatedAt() should be updated")
	}

	// Test invalid transition: active -> draft
	err = tp.UpdateStatus(testplan.StatusDraft)
	if err == nil {
		t.Error("TestPlan.UpdateStatus() should return error for invalid transition")
	}
}

func TestTestPlan_CaseIDs(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	tp, err := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	if err != nil {
		t.Fatalf("Failed to create test plan: %v", err)
	}

	case1 := uuid.New()
	case2 := uuid.New()
	case3 := uuid.New()

	tp.AddCase(case1)
	tp.AddCase(case2)
	tp.AddCase(case3)

	cases := tp.CaseIDs()
	if len(cases) != 3 {
		t.Errorf("TestPlan.CaseIDs() length = %v, want 3", len(cases))
	}
}
