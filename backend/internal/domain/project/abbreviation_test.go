// Package project_test tests ModuleAbbreviation value object
package project_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/project"
)

func TestParseModuleAbbreviation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    project.ModuleAbbreviation
		wantErr bool
	}{
		{
			name:  "valid 2 letter abbreviation",
			input: "US",
			want:  project.ModuleAbbreviation("US"),
		},
		{
			name:  "valid 3 letter abbreviation",
			input: "USR",
			want:  project.ModuleAbbreviation("USR"),
		},
		{
			name:  "valid 4 letter abbreviation",
			input: "USER",
			want:  project.ModuleAbbreviation("USER"),
		},
		{
			name:    "single letter - too short",
			input:   "U",
			wantErr: true,
		},
		{
			name:    "5 letters - too long",
			input:   "USERS",
			wantErr: true,
		},
		{
			name:    "lowercase letters",
			input:   "usr",
			wantErr: true,
		},
		{
			name:    "mixed case",
			input:   "UsR",
			wantErr: true,
		},
		{
			name:    "contains numbers",
			input:   "US1",
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
			got, err := project.ParseModuleAbbreviation(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseModuleAbbreviation(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseModuleAbbreviation(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestModuleAbbreviation_String(t *testing.T) {
	tests := []struct {
		name   string
		abbrev project.ModuleAbbreviation
		want   string
	}{
		{
			name:   "USR abbreviation",
			abbrev: project.ModuleAbbreviation("USR"),
			want:   "USR",
		},
		{
			name:   "API abbreviation",
			abbrev: project.ModuleAbbreviation("API"),
			want:   "API",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.abbrev.String(); got != tt.want {
				t.Errorf("ModuleAbbreviation.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleAbbreviation_Equal(t *testing.T) {
	tests := []struct {
		name string
		a1   project.ModuleAbbreviation
		a2   project.ModuleAbbreviation
		want bool
	}{
		{
			name: "same abbreviations",
			a1:   project.ModuleAbbreviation("USR"),
			a2:   project.ModuleAbbreviation("USR"),
			want: true,
		},
		{
			name: "different abbreviations",
			a1:   project.ModuleAbbreviation("USR"),
			a2:   project.ModuleAbbreviation("API"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a1.Equal(tt.a2); got != tt.want {
				t.Errorf("ModuleAbbreviation.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
