// Package project defines ProjectPrefix value object
package project

import (
	"regexp"
)

// ProjectPrefix is a value object representing a unique project identifier prefix
type ProjectPrefix string

var prefixRegex = regexp.MustCompile(`^[A-Z]{2,4}$`)

// ParseProjectPrefix validates and creates a ProjectPrefix
func ParseProjectPrefix(s string) (ProjectPrefix, error) {
	if !prefixRegex.MatchString(s) {
		return "", ErrInvalidProjectPrefix
	}
	return ProjectPrefix(s), nil
}

// String returns the string representation
func (p ProjectPrefix) String() string {
	return string(p)
}

// Equal checks if two prefixes are equal
func (p ProjectPrefix) Equal(other ProjectPrefix) bool {
	return p == other
}
