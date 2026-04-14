// Package project provides project management services
package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/project"
)

// CreateProjectRequest contains project creation data
type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Prefix      string `json:"prefix" validate:"required"`
	Description string `json:"description"`
}

// UpdateProjectRequest contains project update data
type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty"`
}

// CreateModuleRequest contains module creation data
type CreateModuleRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=255"`
	Abbreviation string `json:"abbreviation" validate:"required"`
	Description  string `json:"description"`
}

// ProjectDetail contains project info with statistics
type ProjectDetail struct {
	*project.Project
	ModuleCount   int64 `json:"module_count"`
	CaseCount     int64 `json:"case_count"`
	DocumentCount int64 `json:"document_count"`
}

// ListOptions contains pagination and filtering options
type ListOptions struct {
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
	Keywords string `json:"keywords"`
}

// ProjectService provides project management operations
type ProjectService interface {
	// Project management
	CreateProject(ctx context.Context, req *CreateProjectRequest, userID uuid.UUID) (*project.Project, error)
	GetProject(ctx context.Context, id uuid.UUID) (*ProjectDetail, error)
	ListProjects(ctx context.Context, opts ListOptions) ([]*project.Project, int64, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req *UpdateProjectRequest) (*project.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error

	// Statistics
	GetProjectStatistics(ctx context.Context, id uuid.UUID) (*project.ProjectStatistics, error)

	// Module management
	CreateModule(ctx context.Context, projectID uuid.UUID, req *CreateModuleRequest, userID uuid.UUID) (*project.Module, error)
	ListModules(ctx context.Context, projectID uuid.UUID) ([]*project.Module, error)
	GetModule(ctx context.Context, id uuid.UUID) (*project.Module, error)
	DeleteModule(ctx context.Context, id uuid.UUID) error

	// Config management
	SetConfig(ctx context.Context, projectID uuid.UUID, key string, value map[string]any) error
	GetConfig(ctx context.Context, projectID uuid.UUID, key string) (*project.ProjectConfig, error)
	ListConfigs(ctx context.Context, projectID uuid.UUID) ([]*project.ProjectConfig, error)
	ImportConfigs(ctx context.Context, projectID uuid.UUID, req *ImportConfigsRequest) (*ImportConfigsResult, error)
	ExportConfigs(ctx context.Context, projectID uuid.UUID) ([]map[string]any, error)
}

// ProjectServiceImpl implements ProjectService
type ProjectServiceImpl struct {
	projectRepo project.ProjectRepository
	moduleRepo  project.ModuleRepository
	configRepo  project.ProjectConfigRepository
}

// NewProjectService creates a new ProjectService instance
func NewProjectService(
	projectRepo project.ProjectRepository,
	moduleRepo project.ModuleRepository,
	configRepo project.ProjectConfigRepository,
) ProjectService {
	return &ProjectServiceImpl{
		projectRepo: projectRepo,
		moduleRepo:  moduleRepo,
		configRepo:  configRepo,
	}
}

// CreateProject creates a new project
func (s *ProjectServiceImpl) CreateProject(ctx context.Context, req *CreateProjectRequest, userID uuid.UUID) (*project.Project, error) {
	// Check if name already exists
	existing, err := s.projectRepo.FindByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, project.ErrProjectNameDuplicate
	}
	if err != nil && !errors.Is(err, project.ErrProjectNotFound) {
		return nil, fmt.Errorf("check project name: %w", err)
	}

	// Check if prefix already exists
	existingPrefix, err := s.projectRepo.FindByPrefix(ctx, project.ProjectPrefix(req.Prefix))
	if err == nil && existingPrefix != nil {
		return nil, project.ErrProjectPrefixDuplicate
	}
	if err != nil && !errors.Is(err, project.ErrProjectNotFound) {
		return nil, fmt.Errorf("check project prefix: %w", err)
	}

	// Create new project
	proj, err := project.NewProject(req.Name, req.Prefix, req.Description)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	// Save project
	if err := s.projectRepo.Save(ctx, proj); err != nil {
		return nil, fmt.Errorf("save project: %w", err)
	}

	return proj, nil
}

// GetProject retrieves project details with statistics
func (s *ProjectServiceImpl) GetProject(ctx context.Context, id uuid.UUID) (*ProjectDetail, error) {
	proj, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find project: %w", err)
	}

	// Get module count
	modules, err := s.moduleRepo.FindByProjectID(ctx, id)
	moduleCount := int64(0)
	if err == nil {
		moduleCount = int64(len(modules))
	}

	return &ProjectDetail{
		Project:       proj,
		ModuleCount:   moduleCount,
		CaseCount:     0, // Will be populated by testcase service integration
		DocumentCount: 0, // Will be populated by knowledge service integration
	}, nil
}

// ListProjects lists projects with pagination
func (s *ProjectServiceImpl) ListProjects(ctx context.Context, opts ListOptions) ([]*project.Project, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Limit > 100 {
		opts.Limit = 100
	}

	queryOpts := project.QueryOptions{
		Offset:   opts.Offset,
		Limit:    opts.Limit,
		Keywords: opts.Keywords,
	}

	projects, err := s.projectRepo.FindAll(ctx, queryOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}

	return projects, int64(len(projects)), nil
}

// UpdateProject updates project information
func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, id uuid.UUID, req *UpdateProjectRequest) (*project.Project, error) {
	proj, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find project: %w", err)
	}

	if req.Name != nil {
		// Check if new name already exists for another project
		existing, err := s.projectRepo.FindByName(ctx, *req.Name)
		if err == nil && existing != nil && existing.ID() != id {
			return nil, project.ErrProjectNameDuplicate
		}
		if err := proj.UpdateName(*req.Name); err != nil {
			return nil, fmt.Errorf("update name: %w", err)
		}
	}

	if req.Description != nil {
		proj.UpdateDescription(*req.Description)
	}

	if err := s.projectRepo.Update(ctx, proj); err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	return proj, nil
}

// DeleteProject deletes a project
func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, id uuid.UUID) error {
	// Check if project exists
	_, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find project: %w", err)
	}

	// Delete project
	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}

	return nil
}

// CreateModule creates a new module within a project
func (s *ProjectServiceImpl) CreateModule(ctx context.Context, projectID uuid.UUID, req *CreateModuleRequest, userID uuid.UUID) (*project.Module, error) {
	// Verify project exists
	_, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("find project: %w", err)
	}

	// Check if abbreviation already exists in project
	existing, err := s.moduleRepo.FindByAbbreviation(ctx, projectID, project.ModuleAbbreviation(req.Abbreviation))
	if err == nil && existing != nil {
		return nil, project.ErrModuleAbbrevDuplicate
	}
	if err != nil && !errors.Is(err, project.ErrModuleNotFound) {
		return nil, fmt.Errorf("check module abbreviation: %w", err)
	}

	// Create new module
	module, err := project.NewModule(projectID, req.Name, req.Abbreviation, req.Description, userID)
	if err != nil {
		return nil, fmt.Errorf("create module: %w", err)
	}

	// Save module
	if err := s.moduleRepo.Save(ctx, module); err != nil {
		return nil, fmt.Errorf("save module: %w", err)
	}

	return module, nil
}

// ListModules lists all modules within a project
func (s *ProjectServiceImpl) ListModules(ctx context.Context, projectID uuid.UUID) ([]*project.Module, error) {
	modules, err := s.moduleRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list modules: %w", err)
	}
	return modules, nil
}

// GetModule retrieves a module by ID
func (s *ProjectServiceImpl) GetModule(ctx context.Context, id uuid.UUID) (*project.Module, error) {
	module, err := s.moduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find module: %w", err)
	}
	return module, nil
}

// DeleteModule deletes a module
func (s *ProjectServiceImpl) DeleteModule(ctx context.Context, id uuid.UUID) error {
	// Check if module exists
	_, err := s.moduleRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find module: %w", err)
	}

	// Delete module
	if err := s.moduleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete module: %w", err)
	}

	return nil
}

// SetConfig sets a project configuration
func (s *ProjectServiceImpl) SetConfig(ctx context.Context, projectID uuid.UUID, key string, value map[string]any) error {
	// Verify project exists
	_, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("find project: %w", err)
	}

	// Check if config already exists
	existing, err := s.configRepo.FindByKey(ctx, projectID, key)
	if err == nil && existing != nil {
		// Update existing config and save
		existing.UpdateValue(value) //nolint:errcheck // UpdateValue only sets internal state
		return s.configRepo.Save(ctx, existing)
	}

	// Create new config
	config, err := project.NewProjectConfig(projectID, key, value, "")
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}

	return s.configRepo.Save(ctx, config)
}

// GetConfig retrieves a project configuration
func (s *ProjectServiceImpl) GetConfig(ctx context.Context, projectID uuid.UUID, key string) (*project.ProjectConfig, error) {
	config, err := s.configRepo.FindByKey(ctx, projectID, key)
	if err != nil {
		return nil, fmt.Errorf("find config: %w", err)
	}
	return config, nil
}

// ListConfigs lists all project configurations
func (s *ProjectServiceImpl) ListConfigs(ctx context.Context, projectID uuid.UUID) ([]*project.ProjectConfig, error) {
	configs, err := s.configRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list configs: %w", err)
	}
	return configs, nil
}

// ImportConfigsRequest represents batch import request
type ImportConfigsRequest struct {
	Configs []struct {
		Key         string                 `json:"key" validate:"required"`
		Value       map[string]interface{} `json:"value" validate:"required"`
		Description string                 `json:"description"`
	} `json:"configs" validate:"required"`
}

// ImportConfigsResult represents import result
type ImportConfigsResult struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

// ImportConfigs imports configurations from JSON array
func (s *ProjectServiceImpl) ImportConfigs(ctx context.Context, projectID uuid.UUID, req *ImportConfigsRequest) (*ImportConfigsResult, error) {
	// Validate project exists
	_, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Convert to domain objects
	configs := make([]*project.ProjectConfig, len(req.Configs))
	for i, cfg := range req.Configs {
		config, err := project.NewProjectConfig(projectID, cfg.Key, cfg.Value, cfg.Description)
		if err != nil {
			return nil, fmt.Errorf("invalid config %s: %w", cfg.Key, err)
		}
		configs[i] = config
	}

	// Batch save
	if err := s.configRepo.BatchUpsert(ctx, configs); err != nil {
		return nil, fmt.Errorf("batch upsert configs: %w", err)
	}

	return &ImportConfigsResult{Imported: len(req.Configs)}, nil
}

// ExportConfigs exports project configurations as JSON array
func (s *ProjectServiceImpl) ExportConfigs(ctx context.Context, projectID uuid.UUID) ([]map[string]any, error) {
	// Validate project exists
	_, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return s.configRepo.ExportConfigs(ctx, projectID)
}

// GetProjectStatistics retrieves project statistics
func (s *ProjectServiceImpl) GetProjectStatistics(ctx context.Context, id uuid.UUID) (*project.ProjectStatistics, error) {
	// Validate project exists
	_, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get statistics from repository
	stats, err := s.projectRepo.GetStatistics(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project statistics: %w", err)
	}

	return stats, nil
}
