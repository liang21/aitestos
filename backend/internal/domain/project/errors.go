// Package project defines project domain errors
package project

import "errors"

var (
	// ErrProjectNotFound indicates project does not exist
	ErrProjectNotFound = errors.New("project not found")
	// ErrProjectNameDuplicate indicates project name already exists
	ErrProjectNameDuplicate = errors.New("project name already exists")
	// ErrModuleNotFound indicates module does not exist
	ErrModuleNotFound = errors.New("module not found")
	// ErrModuleNameDuplicate indicates module name already exists in project
	ErrModuleNameDuplicate = errors.New("module name already exists in project")
	// ErrConfigNotFound indicates config not found
	ErrConfigNotFound = errors.New("config not found")
	// ErrConfigKeyDuplicate indicates config key already exists
	ErrConfigKeyDuplicate = errors.New("config key already exists")
)
