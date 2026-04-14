// Package testcase_test tests CaseNumber value object
package testcase_test

import (
	"strings"
	"testing"
	"time"

	"github.com/liang21/aitestos/internal/domain/testcase"
)

func TestParseCaseNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "valid case number",
			input: "ECO-USR-20260402-001",
		},
		{
			name:  "valid case number with different prefix",
			input: "AIT-API-20260101-999",
		},
		{
			name:    "invalid format - missing parts",
			input:   "ECO-USR-20260402",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong separator",
			input:   "ECO_USR_20260402_001",
			wantErr: true,
		},
		{
			name:    "invalid date format",
			input:   "ECO-USR-2026-04-02-001",
			wantErr: true,
		},
		{
			name:    "invalid sequence - not 3 digits",
			input:   "ECO-USR-20260402-01",
			wantErr: true,
		},
		{
			name:    "invalid prefix - lowercase",
			input:   "eco-usr-20260402-001",
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
			got, err := testcase.ParseCaseNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCaseNumber(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.input {
				t.Errorf("ParseCaseNumber(%q) = %v, want %v", tt.input, got, tt.input)
			}
		})
	}
}

func TestGenerateCaseNumber(t *testing.T) {
	tests := []struct {
		name          string
		projectPrefix string
		moduleAbbrev  string
		seq           int
	}{
		{
			name:          "generate first case",
			projectPrefix: "ECO",
			moduleAbbrev:  "USR",
			seq:           1,
		},
		{
			name:          "generate case 999",
			projectPrefix: "AIT",
			moduleAbbrev:  "API",
			seq:           999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testcase.GenerateCaseNumber(tt.projectPrefix, tt.moduleAbbrev, tt.seq)

			// Check format: PREFIX-ABBREV-DATE-SEQ
			gotStr := got.String()
			if !strings.HasPrefix(gotStr, tt.projectPrefix+"-"+tt.moduleAbbrev+"-") {
				t.Errorf("GenerateCaseNumber() = %v, want prefix %v-%v-", gotStr, tt.projectPrefix, tt.moduleAbbrev)
			}

			// Check date is today
			expectedDate := time.Now().Format("20060102")
			if !strings.Contains(gotStr, expectedDate) {
				t.Errorf("GenerateCaseNumber() should contain today's date %v, got %v", expectedDate, gotStr)
			}
		})
	}
}

func TestCaseNumber_String(t *testing.T) {
	num, _ := testcase.ParseCaseNumber("ECO-USR-20260402-001")
	if got := num.String(); got != "ECO-USR-20260402-001" {
		t.Errorf("CaseNumber.String() = %v, want ECO-USR-20260402-001", got)
	}
}

func TestCaseNumber_Equal(t *testing.T) {
	n1, _ := testcase.ParseCaseNumber("ECO-USR-20260402-001")
	n2, _ := testcase.ParseCaseNumber("ECO-USR-20260402-001")
	n3, _ := testcase.ParseCaseNumber("ECO-USR-20260402-002")

	if !n1.Equal(n2) {
		t.Error("CaseNumber.Equal() should return true for same case numbers")
	}
	if n1.Equal(n3) {
		t.Error("CaseNumber.Equal() should return false for different case numbers")
	}
}
