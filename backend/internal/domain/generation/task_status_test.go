// Package generation_test tests TaskStatus value object
package generation_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/generation"
)

func TestParseTaskStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    generation.TaskStatus
		wantErr bool
	}{
		{
			name:  "pending status",
			input: "pending",
			want:  generation.TaskPending,
		},
		{
			name:  "processing status",
			input: "processing",
			want:  generation.TaskProcessing,
		},
		{
			name:  "completed status",
			input: "completed",
			want:  generation.TaskCompleted,
		},
		{
			name:  "failed status",
			input: "failed",
			want:  generation.TaskFailed,
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
			got, err := generation.ParseTaskStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTaskStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseTaskStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTaskStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status generation.TaskStatus
		want   string
	}{
		{
			name:   "pending string",
			status: generation.TaskPending,
			want:   "pending",
		},
		{
			name:   "processing string",
			status: generation.TaskProcessing,
			want:   "processing",
		},
		{
			name:   "completed string",
			status: generation.TaskCompleted,
			want:   "completed",
		},
		{
			name:   "failed string",
			status: generation.TaskFailed,
			want:   "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("TaskStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from generation.TaskStatus
		to   generation.TaskStatus
		want bool
	}{
		{
			name: "pending to processing",
			from: generation.TaskPending,
			to:   generation.TaskProcessing,
			want: true,
		},
		{
			name: "processing to completed",
			from: generation.TaskProcessing,
			to:   generation.TaskCompleted,
			want: true,
		},
		{
			name: "processing to failed",
			from: generation.TaskProcessing,
			to:   generation.TaskFailed,
			want: true,
		},
		{
			name: "completed to processing - not allowed",
			from: generation.TaskCompleted,
			to:   generation.TaskProcessing,
			want: false,
		},
		{
			name: "failed to pending - not allowed",
			from: generation.TaskFailed,
			to:   generation.TaskPending,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("TaskStatus.CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name   string
		status generation.TaskStatus
		want   bool
	}{
		{
			name:   "pending is not final",
			status: generation.TaskPending,
			want:   false,
		},
		{
			name:   "processing is not final",
			status: generation.TaskProcessing,
			want:   false,
		},
		{
			name:   "completed is final",
			status: generation.TaskCompleted,
			want:   true,
		},
		{
			name:   "failed is final",
			status: generation.TaskFailed,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsFinal(); got != tt.want {
				t.Errorf("TaskStatus.IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}
