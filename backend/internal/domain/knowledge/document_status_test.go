// Package knowledge_test tests DocumentStatus value object
package knowledge_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/knowledge"
)

func TestParseDocumentStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    knowledge.DocumentStatus
		wantErr bool
	}{
		{
			name:  "pending status",
			input: "pending",
			want:  knowledge.StatusPending,
		},
		{
			name:  "processing status",
			input: "processing",
			want:  knowledge.StatusProcessing,
		},
		{
			name:  "completed status",
			input: "completed",
			want:  knowledge.StatusCompleted,
		},
		{
			name:  "failed status",
			input: "failed",
			want:  knowledge.StatusFailed,
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
			got, err := knowledge.ParseDocumentStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDocumentStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDocumentStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDocumentStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status knowledge.DocumentStatus
		want   string
	}{
		{
			name:   "pending string",
			status: knowledge.StatusPending,
			want:   "pending",
		},
		{
			name:   "processing string",
			status: knowledge.StatusProcessing,
			want:   "processing",
		},
		{
			name:   "completed string",
			status: knowledge.StatusCompleted,
			want:   "completed",
		},
		{
			name:   "failed string",
			status: knowledge.StatusFailed,
			want:   "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("DocumentStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from knowledge.DocumentStatus
		to   knowledge.DocumentStatus
		want bool
	}{
		{
			name: "pending to processing",
			from: knowledge.StatusPending,
			to:   knowledge.StatusProcessing,
			want: true,
		},
		{
			name: "processing to completed",
			from: knowledge.StatusProcessing,
			to:   knowledge.StatusCompleted,
			want: true,
		},
		{
			name: "processing to failed",
			from: knowledge.StatusProcessing,
			to:   knowledge.StatusFailed,
			want: true,
		},
		{
			name: "failed to pending - retry",
			from: knowledge.StatusFailed,
			to:   knowledge.StatusPending,
			want: true,
		},
		{
			name: "completed to pending - not allowed",
			from: knowledge.StatusCompleted,
			to:   knowledge.StatusPending,
			want: false,
		},
		{
			name: "pending to completed - not allowed",
			from: knowledge.StatusPending,
			to:   knowledge.StatusCompleted,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("DocumentStatus.CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name   string
		status knowledge.DocumentStatus
		want   bool
	}{
		{
			name:   "pending is not final",
			status: knowledge.StatusPending,
			want:   false,
		},
		{
			name:   "processing is not final",
			status: knowledge.StatusProcessing,
			want:   false,
		},
		{
			name:   "completed is final",
			status: knowledge.StatusCompleted,
			want:   true,
		},
		{
			name:   "failed is final",
			status: knowledge.StatusFailed,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsFinal(); got != tt.want {
				t.Errorf("DocumentStatus.IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}
