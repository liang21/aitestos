// Package testplan_test tests TestResult entity
package testplan_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testplan"
)

func TestNewTestResult(t *testing.T) {
	planID := uuid.New()
	caseID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		planID  uuid.UUID
		caseID  uuid.UUID
		userID  uuid.UUID
		status  testplan.ResultStatus
		note    string
		wantErr bool
	}{
		{
			name:    "valid pass result",
			planID:  planID,
			caseID:  caseID,
			userID:  userID,
			status:  testplan.ResultPass,
			note:    "All steps passed",
			wantErr: false,
		},
		{
			name:    "valid fail result",
			planID:  planID,
			caseID:  caseID,
			userID:  userID,
			status:  testplan.ResultFail,
			note:    "Step 3 failed",
			wantErr: false,
		},
		{
			name:    "valid block result",
			planID:  planID,
			caseID:  caseID,
			userID:  userID,
			status:  testplan.ResultBlock,
			note:    "Blocked by dependency",
			wantErr: false,
		},
		{
			name:    "valid skip result",
			planID:  planID,
			caseID:  caseID,
			userID:  userID,
			status:  testplan.ResultSkip,
			note:    "Not applicable",
			wantErr: false,
		},
		{
			name:    "nil plan ID",
			planID:  uuid.Nil,
			caseID:  caseID,
			userID:  userID,
			status:  testplan.ResultPass,
			wantErr: true,
		},
		{
			name:    "nil case ID",
			planID:  planID,
			caseID:  uuid.Nil,
			userID:  userID,
			status:  testplan.ResultPass,
			wantErr: true,
		},
		{
			name:    "nil user ID",
			planID:  planID,
			caseID:  caseID,
			userID:  uuid.Nil,
			status:  testplan.ResultPass,
			wantErr: true,
		},
		{
			name:    "empty note is allowed",
			planID:  planID,
			caseID:  caseID,
			userID:  userID,
			status:  testplan.ResultPass,
			note:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testplan.NewTestResult(tt.planID, tt.caseID, tt.userID, tt.status, tt.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewTestResult() returned nil result")
					return
				}
				if got.Status() != tt.status {
					t.Errorf("TestResult.Status() = %v, want %v", got.Status(), tt.status)
				}
				if got.Note() != tt.note {
					t.Errorf("TestResult.Note() = %v, want %v", got.Note(), tt.note)
				}
			}
		})
	}
}

func TestTestResult_Accessors(t *testing.T) {
	planID := uuid.New()
	caseID := uuid.New()
	userID := uuid.New()

	tr, err := testplan.NewTestResult(planID, caseID, userID, testplan.ResultPass, "Test note")
	if err != nil {
		t.Fatalf("Failed to create test result: %v", err)
	}

	if tr.ID() == uuid.Nil {
		t.Error("TestResult.ID() should not be nil")
	}
	if tr.PlanID() != planID {
		t.Errorf("TestResult.PlanID() = %v, want %v", tr.PlanID(), planID)
	}
	if tr.CaseID() != caseID {
		t.Errorf("TestResult.CaseID() = %v, want %v", tr.CaseID(), caseID)
	}
	if tr.ExecutedBy() != userID {
		t.Errorf("TestResult.ExecutedBy() = %v, want %v", tr.ExecutedBy(), userID)
	}
	if tr.Status() != testplan.ResultPass {
		t.Errorf("TestResult.Status() = %v, want pass", tr.Status())
	}
	if tr.Note() != "Test note" {
		t.Errorf("TestResult.Note() = %v, want Test note", tr.Note())
	}
	if tr.ExecutedAt().IsZero() {
		t.Error("TestResult.ExecutedAt() should not be zero")
	}
}

func TestResultStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status testplan.ResultStatus
		want   string
	}{
		{
			name:   "pass string",
			status: testplan.ResultPass,
			want:   "pass",
		},
		{
			name:   "fail string",
			status: testplan.ResultFail,
			want:   "fail",
		},
		{
			name:   "block string",
			status: testplan.ResultBlock,
			want:   "block",
		},
		{
			name:   "skip string",
			status: testplan.ResultSkip,
			want:   "skip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("ResultStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseResultStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    testplan.ResultStatus
		wantErr bool
	}{
		{
			name:  "pass status",
			input: "pass",
			want:  testplan.ResultPass,
		},
		{
			name:  "fail status",
			input: "fail",
			want:  testplan.ResultFail,
		},
		{
			name:  "block status",
			input: "block",
			want:  testplan.ResultBlock,
		},
		{
			name:  "skip status",
			input: "skip",
			want:  testplan.ResultSkip,
		},
		{
			name:    "invalid status",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testplan.ParseResultStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResultStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseResultStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTestResult_UpdateNote(t *testing.T) {
	planID := uuid.New()
	caseID := uuid.New()
	userID := uuid.New()

	tr, err := testplan.NewTestResult(planID, caseID, userID, testplan.ResultPass, "Original note")
	if err != nil {
		t.Fatalf("Failed to create test result: %v", err)
	}

	originalUpdatedAt := tr.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	tr.UpdateNote("Updated note")

	if tr.Note() != "Updated note" {
		t.Errorf("TestResult.Note() = %v, want Updated note", tr.Note())
	}
	if !tr.UpdatedAt().After(originalUpdatedAt) {
		t.Error("TestResult.UpdatedAt() should be updated")
	}
}
