// Package testcase defines CaseStatus value object
package testcase

import "errors"

// CaseStatus is a value object representing the execution status of a test case
type CaseStatus string

const (
	// StatusUnexecuted means the case has not been executed
	StatusUnexecuted CaseStatus = "unexecuted"
	// StatusPass means the case passed
	StatusPass CaseStatus = "pass"
	// StatusFail means the case failed
	StatusFail CaseStatus = "fail"
	// StatusBlock means the case is blocked
	StatusBlock CaseStatus = "block"
)

// ParseCaseStatus validates and creates a CaseStatus
func ParseCaseStatus(s string) (CaseStatus, error) {
	switch CaseStatus(s) {
	case StatusUnexecuted, StatusPass, StatusFail, StatusBlock:
		return CaseStatus(s), nil
	default:
		return "", errors.New("invalid case status")
	}
}

// String returns the string representation
func (s CaseStatus) String() string {
	return string(s)
}

// IsFinal returns true if the status is a final state
func (s CaseStatus) IsFinal() bool {
	return s == StatusPass || s == StatusFail || s == StatusBlock
}
