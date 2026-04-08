// Package project defines ProjectConfig entity
package project

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ProjectConfig is an entity for storing project-level key-value configuration
type ProjectConfig struct {
	id          uuid.UUID
	projectID   uuid.UUID
	key         string
	value       map[string]any
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProjectConfig creates a new project configuration
func NewProjectConfig(projectID uuid.UUID, key string, value map[string]any, description string) (*ProjectConfig, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}
	if value == nil {
		return nil, errors.New("value cannot be nil")
	}

	now := time.Now()
	return &ProjectConfig{
		id:          uuid.New(),
		projectID:   projectID,
		key:         key,
		value:       value,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ID returns the config's unique identifier
func (c *ProjectConfig) ID() uuid.UUID {
	return c.id
}

// ProjectID returns the associated project's ID
func (c *ProjectConfig) ProjectID() uuid.UUID {
	return c.projectID
}

// Key returns the configuration key
func (c *ProjectConfig) Key() string {
	return c.key
}

// Value returns the configuration value
func (c *ProjectConfig) Value() map[string]any {
	return c.value
}

// Description returns the configuration description
func (c *ProjectConfig) Description() string {
	return c.description
}

// CreatedAt returns the creation timestamp
func (c *ProjectConfig) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the last update timestamp
func (c *ProjectConfig) UpdatedAt() time.Time {
	return c.updatedAt
}

// UpdateValue updates the configuration value
func (c *ProjectConfig) UpdateValue(value map[string]any) error {
	if value == nil {
		return errors.New("value cannot be nil")
	}
	c.value = value
	c.updatedAt = time.Now()
	return nil
}

// UpdateDescription updates the configuration description
func (c *ProjectConfig) UpdateDescription(description string) error {
	c.description = description
	c.updatedAt = time.Now()
	return nil
}

// Equal checks if two configs have the same projectID and key
func (c *ProjectConfig) Equal(other *ProjectConfig) bool {
	if other == nil {
		return false
	}
	return c.projectID == other.projectID && c.key == other.key
}

// ReconstructProjectConfig reconstructs a ProjectConfig from stored data
func ReconstructProjectConfig(
	id uuid.UUID,
	projectID uuid.UUID,
	key string,
	value map[string]any,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
) (*ProjectConfig, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}
	if value == nil {
		return nil, errors.New("value cannot be nil")
	}

	return &ProjectConfig{
		id:          id,
		projectID:   projectID,
		key:         key,
		value:       value,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}
