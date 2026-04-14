// Package project_test tests ProjectPrefix value object
package project_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/project"
)

func TestParseProjectPrefix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    project.ProjectPrefix
		wantErr bool
	}{
		{
			name:  "valid 2 letter prefix",
			input: "EC",
			want:  project.ProjectPrefix("EC"),
		},
		{
			name:  "valid 3 letter prefix",
			input: "AIT",
			want:  project.ProjectPrefix("AIT"),
		},
		{
			name:  "valid 4 letter prefix",
			input: "TEST",
			want:  project.ProjectPrefix("TEST"),
		},
		{
			name:    "single letter - too short",
			input:   "A",
			wantErr: true,
		},
		{
			name:    "5 letters - too long",
			input:   "ABCDE",
			wantErr: true,
		},
		{
			name:    "lowercase letters",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "mixed case",
			input:   "AbC",
			wantErr: true,
		},
		{
			name:    "contains numbers",
			input:   "AB1",
			wantErr: true,
		},
		{
			name:    "contains special characters",
			input:   "A-B",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := project.ParseProjectPrefix(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProjectPrefix(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseProjectPrefix(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestProjectPrefix_String(t *testing.T) {
	tests := []struct {
		name   string
		prefix project.ProjectPrefix
		want   string
	}{
		{
			name:   "EC prefix",
			prefix: project.ProjectPrefix("EC"),
			want:   "EC",
		},
		{
			name:   "TEST prefix",
			prefix: project.ProjectPrefix("TEST"),
			want:   "TEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prefix.String(); got != tt.want {
				t.Errorf("ProjectPrefix.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectPrefix_Equal(t *testing.T) {
	tests := []struct {
		name string
		p1   project.ProjectPrefix
		p2   project.ProjectPrefix
		want bool
	}{
		{
			name: "same prefixes",
			p1:   project.ProjectPrefix("EC"),
			p2:   project.ProjectPrefix("EC"),
			want: true,
		},
		{
			name: "different prefixes",
			p1:   project.ProjectPrefix("EC"),
			p2:   project.ProjectPrefix("TEST"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p1.Equal(tt.p2); got != tt.want {
				t.Errorf("ProjectPrefix.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
