// Package project defines repository interfaces
package project

import (
	"context"

	"github.com/google/uuid"
)

// ProjectRepository defines the interface for project persistence
type ProjectRepository interface {
	// Save persists a new project
	Save(ctx context.Context, project *Project) error

	// FindByID retrieves a project by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Project, error)

	// FindByName retrieves a project by name
	FindByName(ctx context.Context, name string) (*Project, error)

	// FindByPrefix retrieves a project by prefix
	FindByPrefix(ctx context.Context, prefix ProjectPrefix) (*Project, error)

	// FindAll retrieves all projects with pagination
	FindAll(ctx context.Context, opts QueryOptions) ([]*Project, error)

	// Update updates an existing project
	Update(ctx context.Context, project *Project) error

	// Delete removes a project
	Delete(ctx context.Context, id uuid.UUID) error

	// GetStatistics retrieves aggregated project statistics
	GetStatistics(ctx context.Context, id uuid.UUID) (*ProjectStatistics, error)

	// SetStatistics stores statistics in cache (for warmup)
	SetStatistics(ctx context.Context, id uuid.UUID, stats *ProjectStatistics) error
}

// ModuleRepository defines the interface for module persistence
type ModuleRepository interface {
	// Save persists a new module
	Save(ctx context.Context, module *Module) error

	// FindByID retrieves a module by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Module, error)

	// FindByProjectID retrieves all modules for a project
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*Module, error)

	// FindByAbbreviation retrieves a module by abbreviation within a project
	FindByAbbreviation(ctx context.Context, projectID uuid.UUID, abbrev ModuleAbbreviation) (*Module, error)

	// Update updates an existing module
	Update(ctx context.Context, module *Module) error

	// Delete removes a module
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProjectConfigRepository defines the interface for project configuration persistence
type ProjectConfigRepository interface {
	// Save persists a project configuration
	Save(ctx context.Context, config *ProjectConfig) error

	// FindByProjectID retrieves all configurations for a project
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*ProjectConfig, error)

	// FindByKey retrieves a configuration by project ID and key
	FindByKey(ctx context.Context, projectID uuid.UUID, key string) (*ProjectConfig, error)

	// Delete removes a configuration
	Delete(ctx context.Context, id uuid.UUID) error

	// Update updates an existing configuration
	Update(ctx context.Context, config *ProjectConfig) error

	// BatchUpsert batch inserts or updates configurations
	BatchUpsert(ctx context.Context, configs []*ProjectConfig) error

	// ExportConfigs exports all configurations for a project
	ExportConfigs(ctx context.Context, projectID uuid.UUID) ([]map[string]any, error)
}

// QueryOptions holds pagination and filtering options
type QueryOptions struct {
	Offset   int
	Limit    int
	OrderBy  string
	Keywords string
}
