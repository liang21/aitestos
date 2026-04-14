// Package testcase defines testcase domain errors
package testcase

import "errors"

var (
	// ErrCaseNotFound indicates test case does not exist
	ErrCaseNotFound = errors.New("test case not found")
	// ErrCaseNumberDuplicate indicates case number already exists
	ErrCaseNumberDuplicate = errors.New("case number already exists")
	// ErrInvalidCaseNumber indicates case number format is invalid
	ErrInvalidCaseNumber = errors.New("invalid case number format")
	// ErrEmptySteps indicates test case has no steps
	ErrEmptySteps = errors.New("test case steps cannot be empty")
	// ErrInvalidPriority indicates priority is invalid
	ErrInvalidPriority = errors.New("invalid priority")
)
