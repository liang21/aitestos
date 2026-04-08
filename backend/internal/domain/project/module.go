// Package project defines Module entity
package project

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Module represents a module within a project
type Module struct {
	id           uuid.UUID
	projectID    uuid.UUID
	name         string
	abbreviation ModuleAbbreviation
	description  string
	createdBy    uuid.UUID
	createdAt    time.Time
	updatedAt    time.Time
}

// NewModule creates a new module
func NewModule(projectID uuid.UUID, name, abbrevStr, description string, userID uuid.UUID) (*Module, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	if name == "" {
		return nil, errors.New("module name cannot be empty")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be nil")
	}

	abbrev, err := ParseModuleAbbreviation(abbrevStr)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Module{
		id:           uuid.New(),
		projectID:    projectID,
		name:         name,
		abbreviation: abbrev,
		description:  description,
		createdBy:    userID,
		createdAt:    now,
		updatedAt:    now,
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

// Abbreviation returns the module's abbreviation
func (m *Module) Abbreviation() ModuleAbbreviation {
	return m.abbreviation
}

// Description returns the module's description
func (m *Module) Description() string {
	return m.description
}

// CreatedBy returns the user who created this module
func (m *Module) CreatedBy() uuid.UUID {
	return m.createdBy
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

// ReconstructModule reconstructs a Module from stored data
func ReconstructModule(
	id uuid.UUID,
	projectID uuid.UUID,
	name string,
	abbreviation ModuleAbbreviation,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
) *Module {
	return &Module{
		id:           id,
		projectID:    projectID,
		name:         name,
		abbreviation: abbreviation,
		description:  description,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
