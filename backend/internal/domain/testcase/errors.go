// Package testcase defines testcase domain errors
package testcase

import "errors"

var (
	// ErrCaseNotFound indicates test case does not exist
	ErrCaseNotFound = errors.New("test case not found")
	// ErrEmptySteps indicates test case has no steps
	ErrEmptySteps = errors.New("test case steps cannot be empty")
	// ErrInvalidPriority indicates priority is invalid
	ErrInvalidPriority = errors.New("invalid priority")
)
