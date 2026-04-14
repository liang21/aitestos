// Package testcase_test tests CaseStatus value object
package testcase_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/testcase"
)

func TestParseCaseStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    testcase.CaseStatus
		wantErr bool
	}{
		{
			name:  "unexecuted status",
			input: "unexecuted",
			want:  testcase.StatusUnexecuted,
		},
		{
			name:  "pass status",
			input: "pass",
			want:  testcase.StatusPass,
		},
		{
			name:  "fail status",
			input: "fail",
			want:  testcase.StatusFail,
		},
		{
			name:  "block status",
			input: "block",
			want:  testcase.StatusBlock,
		},
		{
			name:    "invalid status",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty status",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testcase.ParseCaseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCaseStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseCaseStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCaseStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status testcase.CaseStatus
		want   string
	}{
		{
			name:   "unexecuted string",
			status: testcase.StatusUnexecuted,
			want:   "unexecuted",
		},
		{
			name:   "pass string",
			status: testcase.StatusPass,
			want:   "pass",
		},
		{
			name:   "fail string",
			status: testcase.StatusFail,
			want:   "fail",
		},
		{
			name:   "block string",
			status: testcase.StatusBlock,
			want:   "block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("CaseStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCaseStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name   string
		status testcase.CaseStatus
		want   bool
	}{
		{
			name:   "unexecuted is not final",
			status: testcase.StatusUnexecuted,
			want:   false,
		},
		{
			name:   "pass is final",
			status: testcase.StatusPass,
			want:   true,
		},
		{
			name:   "fail is final",
			status: testcase.StatusFail,
			want:   true,
		},
		{
			name:   "block is final",
			status: testcase.StatusBlock,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsFinal(); got != tt.want {
				t.Errorf("CaseStatus.IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}
