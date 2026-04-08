// Package testplan_test tests PlanStatus value object
package testplan_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/testplan"
)

func TestParsePlanStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    testplan.PlanStatus
		wantErr bool
	}{
		{
			name:  "draft status",
			input: "draft",
			want:  testplan.StatusDraft,
		},
		{
			name:  "active status",
			input: "active",
			want:  testplan.StatusActive,
		},
		{
			name:  "completed status",
			input: "completed",
			want:  testplan.StatusCompleted,
		},
		{
			name:  "archived status",
			input: "archived",
			want:  testplan.StatusArchived,
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
			got, err := testplan.ParsePlanStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePlanStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParsePlanStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPlanStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status testplan.PlanStatus
		want   string
	}{
		{
			name:   "draft string",
			status: testplan.StatusDraft,
			want:   "draft",
		},
		{
			name:   "active string",
			status: testplan.StatusActive,
			want:   "active",
		},
		{
			name:   "completed string",
			status: testplan.StatusCompleted,
			want:   "completed",
		},
		{
			name:   "archived string",
			status: testplan.StatusArchived,
			want:   "archived",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("PlanStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from testplan.PlanStatus
		to   testplan.PlanStatus
		want bool
	}{
		{
			name: "draft to active",
			from: testplan.StatusDraft,
			to:   testplan.StatusActive,
			want: true,
		},
		{
			name: "active to completed",
			from: testplan.StatusActive,
			to:   testplan.StatusCompleted,
			want: true,
		},
		{
			name: "active to archived",
			from: testplan.StatusActive,
			to:   testplan.StatusArchived,
			want: true,
		},
		{
			name: "completed to archived",
			from: testplan.StatusCompleted,
			to:   testplan.StatusArchived,
			want: true,
		},
		{
			name: "draft to completed - not allowed",
			from: testplan.StatusDraft,
			to:   testplan.StatusCompleted,
			want: false,
		},
		{
			name: "archived to active - not allowed",
			from: testplan.StatusArchived,
			to:   testplan.StatusActive,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("PlanStatus.CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
