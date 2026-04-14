// Package project defines Project aggregate
package project

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Project is the aggregate root for project context
type Project struct {
	id          uuid.UUID
	name        string
	prefix      ProjectPrefix
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProject creates a new project
func NewProject(name, prefixStr, description string) (*Project, error) {
	if name == "" {
		return nil, errors.New("project name cannot be empty")
	}

	prefix, err := ParseProjectPrefix(prefixStr)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Project{
		id:          uuid.New(),
		name:        name,
		prefix:      prefix,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ID returns the project's unique identifier
func (p *Project) ID() uuid.UUID {
	return p.id
}

// Name returns the project's name
func (p *Project) Name() string {
	return p.name
}

// Prefix returns the project's prefix
func (p *Project) Prefix() ProjectPrefix {
	return p.prefix
}

// Description returns the project's description
func (p *Project) Description() string {
	return p.description
}

// CreatedAt returns the creation timestamp
func (p *Project) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last update timestamp
func (p *Project) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdateDescription updates the project's description
func (p *Project) UpdateDescription(description string) {
	p.description = description
	p.updatedAt = time.Now()
}

// UpdateName updates the project's name
func (p *Project) UpdateName(name string) error {
	if name == "" {
		return errors.New("project name cannot be empty")
	}
	p.name = name
	p.updatedAt = time.Now()
	return nil
}

// Reconstruct reconstructs a Project from stored data
func Reconstruct(
	id uuid.UUID,
	name string,
	prefix ProjectPrefix,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
) *Project {
	return &Project{
		id:          id,
		name:        name,
		prefix:      prefix,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
