// Package generation_test tests DraftStatus value object
package generation_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/generation"
)

func TestParseDraftStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    generation.DraftStatus
		wantErr bool
	}{
		{
			name:  "pending status",
			input: "pending",
			want:  generation.DraftPending,
		},
		{
			name:  "confirmed status",
			input: "confirmed",
			want:  generation.DraftConfirmed,
		},
		{
			name:  "rejected status",
			input: "rejected",
			want:  generation.DraftRejected,
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
			got, err := generation.ParseDraftStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDraftStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDraftStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDraftStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status generation.DraftStatus
		want   string
	}{
		{
			name:   "pending string",
			status: generation.DraftPending,
			want:   "pending",
		},
		{
			name:   "confirmed string",
			status: generation.DraftConfirmed,
			want:   "confirmed",
		},
		{
			name:   "rejected string",
			status: generation.DraftRejected,
			want:   "rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("DraftStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDraftStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from generation.DraftStatus
		to   generation.DraftStatus
		want bool
	}{
		{
			name: "pending to confirmed",
			from: generation.DraftPending,
			to:   generation.DraftConfirmed,
			want: true,
		},
		{
			name: "pending to rejected",
			from: generation.DraftPending,
			to:   generation.DraftRejected,
			want: true,
		},
		{
			name: "confirmed to rejected - not allowed",
			from: generation.DraftConfirmed,
			to:   generation.DraftRejected,
			want: false,
		},
		{
			name: "rejected to confirmed - not allowed",
			from: generation.DraftRejected,
			to:   generation.DraftConfirmed,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("DraftStatus.CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDraftStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name   string
		status generation.DraftStatus
		want   bool
	}{
		{
			name:   "pending is not final",
			status: generation.DraftPending,
			want:   false,
		},
		{
			name:   "confirmed is final",
			status: generation.DraftConfirmed,
			want:   true,
		},
		{
			name:   "rejected is final",
			status: generation.DraftRejected,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsFinal(); got != tt.want {
				t.Errorf("DraftStatus.IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRejectionReason(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    generation.RejectionReason
		wantErr bool
	}{
		{
			name:  "duplicate reason",
			input: "duplicate",
			want:  generation.ReasonDuplicate,
		},
		{
			name:  "irrelevant reason",
			input: "irrelevant",
			want:  generation.ReasonIrrelevant,
		},
		{
			name:  "low quality reason",
			input: "low_quality",
			want:  generation.ReasonLowQuality,
		},
		{
			name:  "other reason",
			input: "other",
			want:  generation.ReasonOther,
		},
		{
			name:    "invalid reason",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generation.ParseRejectionReason(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRejectionReason(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseRejectionReason(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
