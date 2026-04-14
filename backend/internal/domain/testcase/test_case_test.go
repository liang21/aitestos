// Package testcase_test tests TestCase aggregate
package testcase_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

func TestNewTestCase(t *testing.T) {
	moduleID := uuid.New()
	userID := uuid.New()
	number := testcase.GenerateCaseNumber("ECO", "USR", 1)

	tests := []struct {
		name          string
		moduleID      uuid.UUID
		userID        uuid.UUID
		number        testcase.CaseNumber
		title         string
		preconditions testcase.Preconditions
		steps         testcase.Steps
		expected      testcase.ExpectedResult
		caseType      testcase.CaseType
		priority      testcase.Priority
		wantErr       bool
	}{
		{
			name:          "valid test case",
			moduleID:      moduleID,
			userID:        userID,
			number:        number,
			title:         "Test Case Title",
			preconditions: testcase.Preconditions{"Precondition 1", "Precondition 2"},
			steps:         testcase.Steps{"Step 1", "Step 2"},
			expected:      testcase.ExpectedResult{"result": "success"},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       false,
		},
		{
			name:          "empty title",
			moduleID:      moduleID,
			userID:        userID,
			number:        number,
			title:         "",
			preconditions: testcase.Preconditions{},
			steps:         testcase.Steps{"Step"},
			expected:      testcase.ExpectedResult{},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       true,
		},
		{
			name:          "empty steps",
			moduleID:      moduleID,
			userID:        userID,
			number:        number,
			title:         "Test Case",
			preconditions: testcase.Preconditions{},
			steps:         testcase.Steps{},
			expected:      testcase.ExpectedResult{},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       true,
		},
		{
			name:          "nil module ID",
			moduleID:      uuid.Nil,
			userID:        userID,
			number:        number,
			title:         "Test Case",
			preconditions: testcase.Preconditions{},
			steps:         testcase.Steps{"Step"},
			expected:      testcase.ExpectedResult{},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       true,
		},
		{
			name:          "nil user ID",
			moduleID:      moduleID,
			userID:        uuid.Nil,
			number:        number,
			title:         "Test Case",
			preconditions: testcase.Preconditions{},
			steps:         testcase.Steps{"Step"},
			expected:      testcase.ExpectedResult{},
			caseType:      testcase.CaseTypeFunctionality,
			priority:      testcase.PriorityP0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testcase.NewTestCase(
				tt.moduleID,
				tt.userID,
				tt.number,
				tt.title,
				tt.preconditions,
				tt.steps,
				tt.expected,
				tt.caseType,
				tt.priority,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewTestCase() returned nil test case")
					return
				}
				if got.Title() != tt.title {
					t.Errorf("TestCase.Title() = %v, want %v", got.Title(), tt.title)
				}
				if got.Status() != testcase.StatusUnexecuted {
					t.Errorf("TestCase.Status() = %v, want unexecuted", got.Status())
				}
			}
		})
	}
}

func TestTestCase_Accessors(t *testing.T) {
	moduleID := uuid.New()
	userID := uuid.New()
	number := testcase.GenerateCaseNumber("ECO", "USR", 1)

	tc, err := testcase.NewTestCase(
		moduleID,
		userID,
		number,
		"Test Case Title",
		testcase.Preconditions{"Precondition 1"},
		testcase.Steps{"Step 1", "Step 2"},
		testcase.ExpectedResult{"result": "success"},
		testcase.CaseTypeFunctionality,
		testcase.PriorityP0,
	)
	if err != nil {
		t.Fatalf("Failed to create test case: %v", err)
	}

	// Test ID accessor
	if tc.ID() == uuid.Nil {
		t.Error("TestCase.ID() should not be nil")
	}

	// Test ModuleID accessor
	if tc.ModuleID() != moduleID {
		t.Errorf("TestCase.ModuleID() = %v, want %v", tc.ModuleID(), moduleID)
	}

	// Test UserID accessor
	if tc.UserID() != userID {
		t.Errorf("TestCase.UserID() = %v, want %v", tc.UserID(), userID)
	}

	// Test Number accessor
	if tc.Number().String() == "" {
		t.Error("TestCase.Number() should not be empty")
	}

	// Test Title accessor
	if tc.Title() != "Test Case Title" {
		t.Errorf("TestCase.Title() = %v, want Test Case Title", tc.Title())
	}

	// Test Steps accessor
	if len(tc.Steps()) != 2 {
		t.Errorf("TestCase.Steps() length = %v, want 2", len(tc.Steps()))
	}

	// Test CaseType accessor
	if tc.CaseType() != testcase.CaseTypeFunctionality {
		t.Errorf("TestCase.CaseType() = %v, want functionality", tc.CaseType())
	}

	// Test Priority accessor
	if tc.Priority() != testcase.PriorityP0 {
		t.Errorf("TestCase.Priority() = %v, want P0", tc.Priority())
	}

	// Test Status accessor
	if tc.Status() != testcase.StatusUnexecuted {
		t.Errorf("TestCase.Status() = %v, want unexecuted", tc.Status())
	}

	// Test CreatedAt accessor
	if tc.CreatedAt().IsZero() {
		t.Error("TestCase.CreatedAt() should not be zero")
	}

	// Test UpdatedAt accessor
	if tc.UpdatedAt().IsZero() {
		t.Error("TestCase.UpdatedAt() should not be zero")
	}
}

func TestTestCase_UpdateStatus(t *testing.T) {
	tc := createValidTestCase(t)

	originalUpdatedAt := tc.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	tc.UpdateStatus(testcase.StatusPass)

	if tc.Status() != testcase.StatusPass {
		t.Errorf("TestCase.Status() = %v, want pass", tc.Status())
	}
	if !tc.UpdatedAt().After(originalUpdatedAt) {
		t.Error("TestCase.UpdatedAt() should be updated")
	}
}

func TestTestCase_UpdateTitle(t *testing.T) {
	tc := createValidTestCase(t)

	originalUpdatedAt := tc.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err := tc.UpdateTitle("New Title")
	if err != nil {
		t.Errorf("TestCase.UpdateTitle() error = %v", err)
	}
	if tc.Title() != "New Title" {
		t.Errorf("TestCase.Title() = %v, want New Title", tc.Title())
	}
	if !tc.UpdatedAt().After(originalUpdatedAt) {
		t.Error("TestCase.UpdatedAt() should be updated")
	}

	// Test empty title should fail
	err = tc.UpdateTitle("")
	if err == nil {
		t.Error("TestCase.UpdateTitle() should return error for empty title")
	}
}

func TestTestCase_UpdateSteps(t *testing.T) {
	tc := createValidTestCase(t)

	newSteps := testcase.Steps{"New Step 1", "New Step 2", "New Step 3"}
	err := tc.UpdateSteps(newSteps)
	if err != nil {
		t.Errorf("TestCase.UpdateSteps() error = %v", err)
	}
	if len(tc.Steps()) != 3 {
		t.Errorf("TestCase.Steps() length = %v, want 3", len(tc.Steps()))
	}

	// Test empty steps should fail
	err = tc.UpdateSteps(testcase.Steps{})
	if err == nil {
		t.Error("TestCase.UpdateSteps() should return error for empty steps")
	}
}

func TestTestCase_SetAiMetadata(t *testing.T) {
	tc := createValidTestCase(t)

	taskID := uuid.New()
	metadata := testcase.NewAiMetadata(taskID, testcase.ConfidenceHigh, nil, "deepseek-v3")

	tc.SetAiMetadata(metadata)

	if tc.AiMetadata() == nil {
		t.Error("TestCase.AiMetadata() should not be nil")
	}
	if !tc.AiMetadata().IsAIGenerated() {
		t.Error("TestCase.AiMetadata().IsAIGenerated() should be true")
	}
}

func createValidTestCase(t *testing.T) *testcase.TestCase {
	t.Helper()
	moduleID := uuid.New()
	userID := uuid.New()
	number := testcase.GenerateCaseNumber("ECO", "USR", 1)

	tc, err := testcase.NewTestCase(
		moduleID,
		userID,
		number,
		"Test Case Title",
		testcase.Preconditions{"Precondition 1"},
		testcase.Steps{"Step 1", "Step 2"},
		testcase.ExpectedResult{"result": "success"},
		testcase.CaseTypeFunctionality,
		testcase.PriorityP0,
	)
	if err != nil {
		t.Fatalf("Failed to create test case: %v", err)
	}
	return tc
}
