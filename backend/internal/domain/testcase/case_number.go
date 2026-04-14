// Package testcase defines CaseNumber value object
package testcase

import (
	"fmt"
	"regexp"
	"time"
)

// CaseNumber is a value object representing a unique test case identifier
// Format: {ProjectPrefix}-{ModuleAbbreviation}-{Date}-{Sequence}
// Example: ECO-USR-20260402-001
type CaseNumber string

var caseNumberRegex = regexp.MustCompile(`^[A-Z]{2,4}-[A-Z]{2,4}-\d{8}-\d{3}$`)

// ParseCaseNumber validates and creates a CaseNumber
func ParseCaseNumber(s string) (CaseNumber, error) {
	if !caseNumberRegex.MatchString(s) {
		return "", ErrInvalidCaseNumber
	}
	return CaseNumber(s), nil
}

// GenerateCaseNumber creates a new CaseNumber
func GenerateCaseNumber(projectPrefix, moduleAbbrev string, seq int) CaseNumber {
	date := time.Now().Format("20060102")
	return CaseNumber(fmt.Sprintf("%s-%s-%s-%03d", projectPrefix, moduleAbbrev, date, seq))
}

// String returns the string representation
func (n CaseNumber) String() string {
	return string(n)
}

// ProjectPrefix extracts the project prefix from the case number
func (n CaseNumber) ProjectPrefix() string {
	return string(n[:3])
}

// ModuleAbbrev extracts the module abbreviation from the case number
func (n CaseNumber) ModuleAbbrev() string {
	s := string(n)
	parts := splitByHyphen(s)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// Date extracts the date from the case number
func (n CaseNumber) Date() string {
	s := string(n)
	parts := splitByHyphen(s)
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

// Sequence extracts the sequence number from the case number
func (n CaseNumber) Sequence() int {
	s := string(n)
	parts := splitByHyphen(s)
	if len(parts) >= 4 {
		var seq int
		_, _ = fmt.Sscanf(parts[3], "%d", &seq)
		return seq
	}
	return 0
}

// Equal checks if two case numbers are equal
func (n CaseNumber) Equal(other CaseNumber) bool {
	return n == other
}

func splitByHyphen(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '-' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}
