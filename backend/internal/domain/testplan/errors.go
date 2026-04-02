// Package testplan defines testplan domain errors
package testplan

import "errors"

var (
	// ErrPlanNotFound indicates test plan does not exist
	ErrPlanNotFound = errors.New("test plan not found")
	// ErrPlanNameDuplicate indicates plan name already exists
	ErrPlanNameDuplicate = errors.New("plan name already exists")
	// ErrPlanArchived indicates plan is archived and cannot be modified
	ErrPlanArchived = errors.New("plan is archived")
	// ErrResultNotFound indicates test result does not exist
	ErrResultNotFound = errors.New("test result not found")
	// ErrCaseNotInPlan indicates case is not assigned to plan
	ErrCaseNotInPlan = errors.New("case not assigned to plan")
	// ErrDuplicateExecution indicates duplicate execution for same case
	ErrDuplicateExecution = errors.New("duplicate execution for same case")
)
