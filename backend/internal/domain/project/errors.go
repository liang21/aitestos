// Package project defines project domain errors
package project

import "errors"

var (
	// ErrProjectNotFound indicates project does not exist
	ErrProjectNotFound = errors.New("project not found")
	// ErrProjectNameDuplicate indicates project name already exists
	ErrProjectNameDuplicate = errors.New("project name already exists")
	// ErrProjectPrefixDuplicate indicates project prefix already exists
	ErrProjectPrefixDuplicate = errors.New("project prefix already exists")
	// ErrInvalidProjectPrefix indicates project prefix format is invalid
	ErrInvalidProjectPrefix = errors.New("invalid project prefix: must be 2-4 uppercase letters")
	// ErrModuleNotFound indicates module does not exist
	ErrModuleNotFound = errors.New("module not found")
	// ErrModuleNameDuplicate indicates module name already exists in project
	ErrModuleNameDuplicate = errors.New("module name already exists in project")
	// ErrModuleAbbrevDuplicate indicates module abbreviation already exists in project
	ErrModuleAbbrevDuplicate = errors.New("module abbreviation already exists in project")
	// ErrInvalidModuleAbbrev indicates module abbreviation format is invalid
	ErrInvalidModuleAbbrev = errors.New("invalid module abbreviation: must be 2-4 uppercase letters")
	// ErrConfigNotFound indicates config not found
	ErrConfigNotFound = errors.New("config not found")
	// ErrConfigKeyDuplicate indicates config key already exists
	ErrConfigKeyDuplicate = errors.New("config key already exists")
)
