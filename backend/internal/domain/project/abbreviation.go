// Package project defines ModuleAbbreviation value object
package project

import (
	"regexp"
)

// ModuleAbbreviation is a value object representing a unique module abbreviation within a project
type ModuleAbbreviation string

var abbrevRegex = regexp.MustCompile(`^[A-Z]{2,4}$`)

// ParseModuleAbbreviation validates and creates a ModuleAbbreviation
func ParseModuleAbbreviation(s string) (ModuleAbbreviation, error) {
	if !abbrevRegex.MatchString(s) {
		return "", ErrInvalidModuleAbbrev
	}
	return ModuleAbbreviation(s), nil
}

// String returns the string representation
func (a ModuleAbbreviation) String() string {
	return string(a)
}

// Equal checks if two abbreviations are equal
func (a ModuleAbbreviation) Equal(other ModuleAbbreviation) bool {
	return a == other
}
