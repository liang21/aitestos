// Package generation_test tests GeneratedCaseDraft entity
package generation_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

func TestNewGeneratedCaseDraft(t *testing.T) {
	taskID := uuid.New()

	tests := []struct {
		name          string
		taskID        uuid.UUID
		title         string
		preconditions testcase.Preconditions
		steps         testcase.Steps
		expected      testcase.ExpectedResult
		caseType      testcase.CaseType
		priority      testcase.Priority
		wantErr       bool
	}{
		{
			name:          "valid draft",
			taskID:        taskID,
			title:         "Test User Login",
			preconditions: testcase.Preconditions{"User exists", "Browser ready"},
			steps:         testcase.Steps{"Open login page", "Enter credentials", "Click submit"},
			expected:      testcase.ExpectedResult{"status": "success"},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       false,
		},
		{
			name:          "empty title",
			taskID:        taskID,
			title:         "",
			preconditions: testcase.Preconditions{"Precondition"},
			steps:         testcase.Steps{"Step 1"},
			expected:      testcase.ExpectedResult{"result": "pass"},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       true,
		},
		{
			name:          "empty steps",
			taskID:        taskID,
			title:         "Test Case",
			preconditions: testcase.Preconditions{},
			steps:         testcase.Steps{},
			expected:      testcase.ExpectedResult{},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       true,
		},
		{
			name:          "nil task ID",
			taskID:        uuid.Nil,
			title:         "Test Case",
			preconditions: testcase.Preconditions{"Pre"},
			steps:         testcase.Steps{"Step"},
			expected:      testcase.ExpectedResult{"result": "pass"},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generation.NewGeneratedCaseDraft(
				tt.taskID,
				tt.title,
				tt.preconditions,
				tt.steps,
				tt.expected,
				tt.caseType,
				tt.priority,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGeneratedCaseDraft() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewGeneratedCaseDraft() returned nil draft")
					return
				}
				if got.Title() != tt.title {
					t.Errorf("GeneratedCaseDraft.Title() = %v, want %v", got.Title(), tt.title)
				}
				if got.Status() != generation.DraftPending {
					t.Errorf("GeneratedCaseDraft.Status() = %v, want pending", got.Status())
				}
			}
		})
	}
}

func TestGeneratedCaseDraft_Accessors(t *testing.T) {
	taskID := uuid.New()
	draft := createValidDraft(t, taskID)

	if draft.ID() == uuid.Nil {
		t.Error("GeneratedCaseDraft.ID() should not be nil")
	}
	if draft.TaskID() != taskID {
		t.Errorf("GeneratedCaseDraft.TaskID() = %v, want %v", draft.TaskID(), taskID)
	}
	if draft.Title() != "Test User Login" {
		t.Errorf("GeneratedCaseDraft.Title() = %v, want Test User Login", draft.Title())
	}
	if draft.Status() != generation.DraftPending {
		t.Errorf("GeneratedCaseDraft.Status() = %v, want pending", draft.Status())
	}
	if draft.CreatedAt().IsZero() {
		t.Error("GeneratedCaseDraft.CreatedAt() should not be zero")
	}
	if draft.UpdatedAt().IsZero() {
		t.Error("GeneratedCaseDraft.UpdatedAt() should not be zero")
	}
}

func TestGeneratedCaseDraft_Confirm(t *testing.T) {
	taskID := uuid.New()
	draft := createValidDraft(t, taskID)
	moduleID := uuid.New()

	originalUpdatedAt := draft.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	draft.Confirm(moduleID)

	if draft.Status() != generation.DraftConfirmed {
		t.Errorf("GeneratedCaseDraft.Status() = %v, want confirmed", draft.Status())
	}
	if draft.ModuleID() == nil || *draft.ModuleID() != moduleID {
		t.Errorf("GeneratedCaseDraft.ModuleID() = %v, want %v", draft.ModuleID(), moduleID)
	}
	if !draft.UpdatedAt().After(originalUpdatedAt) {
		t.Error("GeneratedCaseDraft.UpdatedAt() should be updated")
	}

	// Cannot confirm already confirmed draft
	err := draft.Confirm(moduleID)
	if err == nil {
		t.Error("Confirm() should fail for already confirmed draft")
	}
}

func TestGeneratedCaseDraft_Reject(t *testing.T) {
	taskID := uuid.New()
	draft := createValidDraft(t, taskID)

	originalUpdatedAt := draft.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	draft.Reject(generation.ReasonLowQuality, "Steps are not detailed enough")

	if draft.Status() != generation.DraftRejected {
		t.Errorf("GeneratedCaseDraft.Status() = %v, want rejected", draft.Status())
	}
	if draft.Feedback() == "" {
		t.Error("GeneratedCaseDraft.Feedback() should not be empty after rejection")
	}
	if !draft.UpdatedAt().After(originalUpdatedAt) {
		t.Error("GeneratedCaseDraft.UpdatedAt() should be updated")
	}

	// Cannot reject already rejected draft
	err := draft.Reject(generation.ReasonOther, "Another reason")
	if err == nil {
		t.Error("Reject() should fail for already rejected draft")
	}
}

func TestGeneratedCaseDraft_SetAiMetadata(t *testing.T) {
	taskID := uuid.New()
	draft := createValidDraft(t, taskID)

	taskID2 := uuid.New()
	metadata := testcase.NewAiMetadata(taskID2, testcase.ConfidenceHigh, nil, "v1.0")

	draft.SetAiMetadata(metadata)

	if draft.AiMetadata() == nil {
		t.Error("GeneratedCaseDraft.AiMetadata() should not be nil")
	}
	if draft.AiMetadata().Confidence() != testcase.ConfidenceHigh {
		t.Errorf("GeneratedCaseDraft.AiMetadata().Confidence() = %v, want high", draft.AiMetadata().Confidence())
	}
}

func createValidDraft(t *testing.T, taskID uuid.UUID) *generation.GeneratedCaseDraft {
	draft, err := generation.NewGeneratedCaseDraft(
		taskID,
		"Test User Login",
		testcase.Preconditions{"User exists"},
		testcase.Steps{"Open page", "Login"},
		testcase.ExpectedResult{"status": "success"},
		testcase.CaseTypeFunctionality,
		testcase.PriorityP0,
	)
	if err != nil {
		t.Fatalf("Failed to create draft: %v", err)
	}
	return draft
}
