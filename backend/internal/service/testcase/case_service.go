// Package testcase provides test case management services
package testcase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	domaintestcase "github.com/liang21/aitestos/internal/domain/testcase"
)

// CreateCaseRequest contains test case creation data
type CreateCaseRequest struct {
	ModuleID      uuid.UUID   `json:"module_id" validate:"required"`
	Title         string      `json:"title" validate:"required,min=2,max=500"`
	Preconditions []string    `json:"preconditions"`
	Steps         []string    `json:"steps" validate:"required,min=1"`
	Expected      map[string]any `json:"expected" validate:"required"`
	CaseType      string      `json:"case_type" validate:"required,oneof=functionality performance api ui security"`
	Priority      string      `json:"priority" validate:"required,oneof=P0 P1 P2 P3"`
}

// UpdateCaseRequest contains test case update data
type UpdateCaseRequest struct {
	Title         *string        `json:"title,omitempty" validate:"omitempty,min=2,max=500"`
	Preconditions []string       `json:"preconditions,omitempty"`
	Steps         []string       `json:"steps,omitempty" validate:"omitempty,min=1"`
	Expected      map[string]any `json:"expected,omitempty"`
	CaseType      *string        `json:"case_type,omitempty" validate:"omitempty,oneof=functionality performance api ui security"`
	Priority      *string        `json:"priority,omitempty" validate:"omitempty,oneof=P0 P1 P2 P3"`
	Status        *string        `json:"status,omitempty" validate:"omitempty,oneof=unexecuted pass fail block"`
}

// CaseDetail contains test case details with related info
type CaseDetail struct {
	*domaintestcase.TestCase
	ModuleName    string `json:"module_name"`
	ProjectName   string `json:"project_name"`
	ProjectPrefix string `json:"project_prefix"`
	CreatedByName string `json:"created_by_name"`
}

// CaseListOptions contains filtering options for listing cases
type CaseListOptions struct {
	ModuleID  uuid.UUID `json:"module_id"`
	ProjectID uuid.UUID `json:"project_id"`
	Status    string    `json:"status"`
	CaseType  string    `json:"case_type"`
	Priority  string    `json:"priority"`
	Keywords  string    `json:"keywords"`
	Offset    int       `json:"offset"`
	Limit     int       `json:"limit"`
}

// CaseService provides test case management operations
type CaseService interface {
	// CreateCase creates a new test case with auto-generated number
	CreateCase(ctx context.Context, req *CreateCaseRequest, userID uuid.UUID) (*domaintestcase.TestCase, error)

	// UpdateCase updates an existing test case
	UpdateCase(ctx context.Context, id uuid.UUID, req *UpdateCaseRequest) (*domaintestcase.TestCase, error)

	// GetCaseDetail retrieves test case with related info
	GetCaseDetail(ctx context.Context, id uuid.UUID) (*CaseDetail, error)

	// GetCaseByNumber retrieves test case by case number
	GetCaseByNumber(ctx context.Context, number domaintestcase.CaseNumber) (*CaseDetail, error)

	// ListByModule lists test cases by module with pagination
	ListByModule(ctx context.Context, moduleID uuid.UUID, opts CaseListOptions) ([]*domaintestcase.TestCase, int64, error)

	// ListByProject lists test cases by project with pagination
	ListByProject(ctx context.Context, projectID uuid.UUID, opts CaseListOptions) ([]*domaintestcase.TestCase, int64, error)

	// DeleteCase soft deletes a test case
	DeleteCase(ctx context.Context, id uuid.UUID) error

	// GenerateCaseNumber generates a new case number for a module
	GenerateCaseNumber(ctx context.Context, moduleID uuid.UUID) (domaintestcase.CaseNumber, error)
}

// CaseServiceImpl implements CaseService
type CaseServiceImpl struct {
	caseRepo   domaintestcase.TestCaseRepository
	moduleRepo ModuleRepository
	projectRepo ProjectRepository
}

// ModuleRepository interface for testcase service (to avoid import cycle)
type ModuleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (Module, error)
}

// ProjectRepository interface for testcase service
type ProjectRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (Project, error)
}

// Module represents a project module (minimal interface)
type Module interface {
	ID() uuid.UUID
	ProjectID() uuid.UUID
	Name() string
	Abbreviation() string
}

// Project represents a project (minimal interface)
type Project interface {
	ID() uuid.UUID
	Name() string
	Prefix() string
}

// ModuleWrapper wraps *domainproject.Module to implement Module interface
type ModuleWrapper struct {
	*domainproject.Module
}

// Abbreviation implements Module interface by converting ModuleAbbreviation to string
func (w ModuleWrapper) Abbreviation() string {
	return string(w.Module.Abbreviation())
}

// ProjectWrapper wraps *domainproject.Project to implement Project interface
type ProjectWrapper struct {
	*domainproject.Project
}

// Prefix implements Project interface by converting ProjectPrefix to string
func (w ProjectWrapper) Prefix() string {
	return string(w.Project.Prefix())
}

// NewCaseService creates a new CaseService instance
func NewCaseService(
	caseRepo domaintestcase.TestCaseRepository,
	moduleRepo ModuleRepository,
	projectRepo ProjectRepository,
) CaseService {
	return &CaseServiceImpl{
		caseRepo:    caseRepo,
		moduleRepo:  moduleRepo,
		projectRepo: projectRepo,
	}
}

// CreateCase creates a new test case with auto-generated number
func (s *CaseServiceImpl) CreateCase(ctx context.Context, req *CreateCaseRequest, userID uuid.UUID) (*domaintestcase.TestCase, error) {
	// Validate steps
	if len(req.Steps) == 0 {
		return nil, domaintestcase.ErrEmptySteps
	}

	// Get module
	module, err := s.moduleRepo.FindByID(ctx, req.ModuleID)
	if err != nil {
		return nil, errors.New("module not found")
	}

	// Get project
	project, err := s.projectRepo.FindByID(ctx, module.ProjectID())
	if err != nil {
		return nil, errors.New("project not found")
	}

	// Generate case number
	caseNumber, err := s.GenerateCaseNumber(ctx, req.ModuleID)
	if err != nil {
		return nil, fmt.Errorf("generate case number: %w", err)
	}

	// Create test case
	tc, err := domaintestcase.NewTestCase(
		req.ModuleID,
		userID,
		caseNumber,
		req.Title,
		req.Preconditions,
		req.Steps,
		req.Expected,
		domaintestcase.CaseType(req.CaseType),
		domaintestcase.Priority(req.Priority),
	)
	if err != nil {
		return nil, fmt.Errorf("create test case: %w", err)
	}

	// Save test case
	if err := s.caseRepo.Save(ctx, tc); err != nil {
		return nil, fmt.Errorf("save test case: %w", err)
	}

	_ = project // Used for future reference

	return tc, nil
}

// UpdateCase updates an existing test case
func (s *CaseServiceImpl) UpdateCase(ctx context.Context, id uuid.UUID, req *UpdateCaseRequest) (*domaintestcase.TestCase, error) {
	// Get test case
	tc, err := s.caseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find test case: %w", err)
	}

	// Update fields
	if req.Title != nil {
		if err := tc.UpdateTitle(*req.Title); err != nil {
			return nil, fmt.Errorf("update title: %w", err)
		}
	}

	if req.Steps != nil {
		if err := tc.UpdateSteps(req.Steps); err != nil {
			return nil, fmt.Errorf("update steps: %w", err)
		}
	}

	if req.Status != nil {
		tc.UpdateStatus(domaintestcase.CaseStatus(*req.Status))
	}

	// Save changes
	if err := s.caseRepo.Update(ctx, tc); err != nil {
		return nil, fmt.Errorf("update test case: %w", err)
	}

	return tc, nil
}

// GetCaseDetail retrieves test case with related info
func (s *CaseServiceImpl) GetCaseDetail(ctx context.Context, id uuid.UUID) (*CaseDetail, error) {
	tc, err := s.caseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find test case: %w", err)
	}

	// Get module and project info
	module, _ := s.moduleRepo.FindByID(ctx, tc.ModuleID())
	detail := &CaseDetail{
		TestCase: tc,
	}

	if module != nil {
		detail.ModuleName = module.Name()
		project, _ := s.projectRepo.FindByID(ctx, module.ProjectID())
		if project != nil {
			detail.ProjectName = project.Name()
			detail.ProjectPrefix = project.Prefix()
		}
	}

	return detail, nil
}

// GetCaseByNumber retrieves test case by case number
func (s *CaseServiceImpl) GetCaseByNumber(ctx context.Context, number domaintestcase.CaseNumber) (*CaseDetail, error) {
	tc, err := s.caseRepo.FindByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("find test case by number: %w", err)
	}

	return s.GetCaseDetail(ctx, tc.ID())
}

// ListByModule lists test cases by module with pagination
func (s *CaseServiceImpl) ListByModule(ctx context.Context, moduleID uuid.UUID, opts CaseListOptions) ([]*domaintestcase.TestCase, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}

	queryOpts := domaintestcase.QueryOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	}

	cases, err := s.caseRepo.FindByModuleID(ctx, moduleID, queryOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("list cases by module: %w", err)
	}

	return cases, int64(len(cases)), nil
}

// ListByProject lists test cases by project with pagination
func (s *CaseServiceImpl) ListByProject(ctx context.Context, projectID uuid.UUID, opts CaseListOptions) ([]*domaintestcase.TestCase, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}

	queryOpts := domaintestcase.QueryOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	}

	cases, err := s.caseRepo.FindByProjectID(ctx, projectID, queryOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("list cases by project: %w", err)
	}

	return cases, int64(len(cases)), nil
}

// DeleteCase soft deletes a test case
func (s *CaseServiceImpl) DeleteCase(ctx context.Context, id uuid.UUID) error {
	if err := s.caseRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete test case: %w", err)
	}
	return nil
}

// GenerateCaseNumber generates a new case number for a module
func (s *CaseServiceImpl) GenerateCaseNumber(ctx context.Context, moduleID uuid.UUID) (domaintestcase.CaseNumber, error) {
	// Get module
	module, err := s.moduleRepo.FindByID(ctx, moduleID)
	if err != nil {
		return "", errors.New("module not found")
	}

	// Get project
	project, err := s.projectRepo.FindByID(ctx, module.ProjectID())
	if err != nil {
		return "", errors.New("project not found")
	}

	// Get count for today
	today := time.Now()
	count, err := s.caseRepo.CountByDate(ctx, moduleID, today)
	if err != nil {
		count = 0
	}

	// Generate case number
	caseNumber := domaintestcase.GenerateCaseNumber(
		project.Prefix(),
		module.Abbreviation(),
		int(count)+1,
	)

	return caseNumber, nil
}
