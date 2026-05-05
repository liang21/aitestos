// Package project defines Module entity
package project

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Module represents a module within a project
type Module struct {
	id          uuid.UUID
	projectID   uuid.UUID
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewModule creates a new module
func NewModule(projectID uuid.UUID, name, description string) (*Module, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	if name == "" {
		return nil, errors.New("module name cannot be empty")
	}

	now := time.Now()
	return &Module{
		id:          uuid.New(),
		projectID:   projectID,
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ID returns the module's unique identifier
func (m *Module) ID() uuid.UUID {
	return m.id
}

// ProjectID returns the associated project's ID
func (m *Module) ProjectID() uuid.UUID {
	return m.projectID
}

// Name returns the module's name
func (m *Module) Name() string {
	return m.name
}

// Description returns the module's description
func (m *Module) Description() string {
	return m.description
}

// CreatedAt returns the creation timestamp
func (m *Module) CreatedAt() time.Time {
	return m.createdAt
}

// UpdatedAt returns the last update timestamp
func (m *Module) UpdatedAt() time.Time {
	return m.updatedAt
}

// UpdateDescription updates the module's description
func (m *Module) UpdateDescription(description string) {
	m.description = description
	m.updatedAt = time.Now()
}

// UpdateName updates the module's name
func (m *Module) UpdateName(name string) {
	m.name = name
	m.updatedAt = time.Now()
}

// ReconstructModule reconstructs a Module from stored data
func ReconstructModule(
	id uuid.UUID,
	projectID uuid.UUID,
	name string,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
) *Module {
	return &Module{
		id:          id,
		projectID:   projectID,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
