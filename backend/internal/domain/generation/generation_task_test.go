// Package generation_test tests GenerationTask aggregate
package generation_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
)

func TestNewGenerationTask(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		projectID uuid.UUID
		moduleID  uuid.UUID
		prompt    string
		userID    uuid.UUID
		wantErr   bool
	}{
		{
			name:      "valid task",
			projectID: projectID,
			moduleID:  uuid.New(),
			prompt:    "Generate test cases for user login feature",
			userID:    userID,
			wantErr:   false,
		},
		{
			name:      "empty prompt",
			projectID: projectID,
			moduleID:  uuid.New(),
			prompt:    "",
			userID:    userID,
			wantErr:   true,
		},
		{
			name:      "prompt too short",
			projectID: projectID,
			moduleID:  uuid.New(),
			prompt:    "short",
			userID:    userID,
			wantErr:   true,
		},
		{
			name:      "nil project ID",
			projectID: uuid.Nil,
			moduleID:  uuid.New(),
			prompt:    "Generate test cases",
			userID:    userID,
			wantErr:   true,
		},
		{
			name:      "nil user ID",
			projectID: projectID,
			moduleID:  uuid.New(),
			prompt:    "Generate test cases",
			userID:    uuid.Nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generation.NewGenerationTask(tt.projectID, tt.moduleID, tt.prompt, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGenerationTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewGenerationTask() returned nil task")
					return
				}
				if got.Prompt() != tt.prompt {
					t.Errorf("GenerationTask.Prompt() = %v, want %v", got.Prompt(), tt.prompt)
				}
				if got.Status() != generation.TaskPending {
					t.Errorf("GenerationTask.Status() = %v, want pending", got.Status())
				}
			}
		})
	}
}

func TestGenerationTask_Accessors(t *testing.T) {
	projectID := uuid.New()
	moduleID := uuid.New()
	userID := uuid.New()

	task, err := generation.NewGenerationTask(projectID, moduleID, "Generate test cases for login", userID)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if task.ID() == uuid.Nil {
		t.Error("GenerationTask.ID() should not be nil")
	}
	if task.ProjectID() != projectID {
		t.Errorf("GenerationTask.ProjectID() = %v, want %v", task.ProjectID(), projectID)
	}
	if task.ModuleID() != moduleID {
		t.Errorf("GenerationTask.ModuleID() = %v, want %v", task.ModuleID(), moduleID)
	}
	if task.Prompt() != "Generate test cases for login" {
		t.Errorf("GenerationTask.Prompt() = %v, want Generate test cases for login", task.Prompt())
	}
	if task.Status() != generation.TaskPending {
		t.Errorf("GenerationTask.Status() = %v, want pending", task.Status())
	}
	if task.CreatedAt().IsZero() {
		t.Error("GenerationTask.CreatedAt() should not be zero")
	}
	if task.UpdatedAt().IsZero() {
		t.Error("GenerationTask.UpdatedAt() should not be zero")
	}
	if task.CreatedBy() != userID {
		t.Errorf("GenerationTask.CreatedBy() = %v, want %v", task.CreatedBy(), userID)
	}
}

func TestGenerationTask_StartProcessing(t *testing.T) {
	projectID := uuid.New()
	moduleID := uuid.New()
	userID := uuid.New()

	task, err := generation.NewGenerationTask(projectID, moduleID, "Generate test cases", userID)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	originalUpdatedAt := task.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	task.StartProcessing()

	if task.Status() != generation.TaskProcessing {
		t.Errorf("GenerationTask.Status() = %v, want processing", task.Status())
	}
	if !task.UpdatedAt().After(originalUpdatedAt) {
		t.Error("GenerationTask.UpdatedAt() should be updated")
	}

	// Cannot start processing again
	err = task.StartProcessing()
	if err == nil {
		t.Error("StartProcessing() should fail for already processing task")
	}
}

func TestGenerationTask_Complete(t *testing.T) {
	projectID := uuid.New()
	moduleID := uuid.New()
	userID := uuid.New()

	task, err := generation.NewGenerationTask(projectID, moduleID, "Generate test cases", userID)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	task.StartProcessing()

	summary := map[string]any{
		"total_drafts":  5,
		"passed_review": 4,
	}
	task.Complete(summary)

	if task.Status() != generation.TaskCompleted {
		t.Errorf("GenerationTask.Status() = %v, want completed", task.Status())
	}
	if task.ResultSummary()["total_drafts"] != 5 {
		t.Errorf("GenerationTask.ResultSummary()[total_drafts] = %v, want 5", task.ResultSummary()["total_drafts"])
	}
}

func TestGenerationTask_Fail(t *testing.T) {
	projectID := uuid.New()
	moduleID := uuid.New()
	userID := uuid.New()

	task, err := generation.NewGenerationTask(projectID, moduleID, "Generate test cases", userID)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	task.StartProcessing()
	task.Fail("LLM API timeout")

	if task.Status() != generation.TaskFailed {
		t.Errorf("GenerationTask.Status() = %v, want failed", task.Status())
	}
	if task.ErrorMsg() != "LLM API timeout" {
		t.Errorf("GenerationTask.ErrorMsg() = %v, want LLM API timeout", task.ErrorMsg())
	}
}
