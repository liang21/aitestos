// Package project provides project management services
package project

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/project"
)

// MockProjectRepository implements project.ProjectRepository for testing
type MockProjectRepository struct {
	projects   map[uuid.UUID]*project.Project
	nameIndex  map[string]*project.Project
	prefixIndex map[string]*project.Project
	saveErr    error
	findErr    error
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects:    make(map[uuid.UUID]*project.Project),
		nameIndex:   make(map[string]*project.Project),
		prefixIndex: make(map[string]*project.Project),
	}
}

func (m *MockProjectRepository) Save(ctx context.Context, p *project.Project) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.projects[p.ID()] = p
	m.nameIndex[p.Name()] = p
	m.prefixIndex[string(p.Prefix())] = p
	return nil
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	p, ok := m.projects[id]
	if !ok {
		return nil, project.ErrProjectNotFound
	}
	return p, nil
}

func (m *MockProjectRepository) FindByName(ctx context.Context, name string) (*project.Project, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	p, ok := m.nameIndex[name]
	if !ok {
		return nil, project.ErrProjectNotFound
	}
	return p, nil
}

func (m *MockProjectRepository) FindByPrefix(ctx context.Context, prefix project.ProjectPrefix) (*project.Project, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	p, ok := m.prefixIndex[string(prefix)]
	if !ok {
		return nil, project.ErrProjectNotFound
	}
	return p, nil
}

func (m *MockProjectRepository) CountByOwnerID(ctx context.Context, ownerID uuid.UUID) (int64, error) {
	count := 0
	for range m.projects {
		// Assuming owner tracking - for now just count all
		count++
	}
	return int64(count), nil
}

func (m *MockProjectRepository) FindAll(ctx context.Context, opts project.QueryOptions) ([]*project.Project, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}

	// Collect all projects
	allProjects := make([]*project.Project, 0, len(m.projects))
	for _, p := range m.projects {
		// Apply keyword filter if specified
		if opts.Keywords != "" {
			if !containsIgnoreCase(p.Name(), opts.Keywords) && !containsIgnoreCase(p.Description(), opts.Keywords) {
				continue
			}
		}
		allProjects = append(allProjects, p)
	}

	// Apply pagination
	start := opts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(allProjects) {
		return []*project.Project{}, nil
	}

	end := start + opts.Limit
	if end > len(allProjects) || opts.Limit <= 0 {
		end = len(allProjects)
	}

	return allProjects[start:end], nil
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	sLower := make([]byte, len(s))
	substrLower := make([]byte, len(substr))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		sLower[i] = c
	}
	for i := 0; i < len(substr); i++ {
		c := substr[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		substrLower[i] = c
	}

	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		match := true
		for j := 0; j < len(substrLower); j++ {
			if sLower[i+j] != substrLower[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func (m *MockProjectRepository) Update(ctx context.Context, p *project.Project) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.projects[p.ID()] = p
	m.nameIndex[p.Name()] = p
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	delete(m.projects, id)
	return nil
}

// MockModuleRepository implements project.ModuleRepository for testing
type MockModuleRepository struct {
	modules    map[uuid.UUID]*project.Module
	abbrevKeys map[string]*project.Module // projectID:abbrev -> module
	saveErr    error
	findErr    error
}

func NewMockModuleRepository() *MockModuleRepository {
	return &MockModuleRepository{
		modules:    make(map[uuid.UUID]*project.Module),
		abbrevKeys: make(map[string]*project.Module),
	}
}

func (m *MockModuleRepository) Save(ctx context.Context, mod *project.Module) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.modules[mod.ID()] = mod
	key := mod.ProjectID().String() + ":" + string(mod.Abbreviation())
	m.abbrevKeys[key] = mod
	return nil
}

func (m *MockModuleRepository) FindByID(ctx context.Context, id uuid.UUID) (*project.Module, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	mod, ok := m.modules[id]
	if !ok {
		return nil, project.ErrModuleNotFound
	}
	return mod, nil
}

func (m *MockModuleRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*project.Module, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	modules := make([]*project.Module, 0)
	for _, mod := range m.modules {
		if mod.ProjectID() == projectID {
			modules = append(modules, mod)
		}
	}
	return modules, nil
}

func (m *MockModuleRepository) FindByAbbreviation(ctx context.Context, projectID uuid.UUID, abbrev project.ModuleAbbreviation) (*project.Module, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	key := projectID.String() + ":" + string(abbrev)
	mod, ok := m.abbrevKeys[key]
	if !ok {
		return nil, project.ErrModuleNotFound
	}
	return mod, nil
}

func (m *MockModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	delete(m.modules, id)
	return nil
}

// MockProjectConfigRepository implements project.ProjectConfigRepository for testing
type MockProjectConfigRepository struct {
	configs  map[string]*project.ProjectConfig // projectID:key -> config
	saveErr  error
	findErr  error
}

func NewMockProjectConfigRepository() *MockProjectConfigRepository {
	return &MockProjectConfigRepository{
		configs: make(map[string]*project.ProjectConfig),
	}
}

func (m *MockProjectConfigRepository) Save(ctx context.Context, cfg *project.ProjectConfig) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	key := cfg.ProjectID().String() + ":" + cfg.Key()
	m.configs[key] = cfg
	return nil
}

func (m *MockProjectConfigRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*project.ProjectConfig, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	configs := make([]*project.ProjectConfig, 0)
	prefix := projectID.String() + ":"
	for k, cfg := range m.configs {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			configs = append(configs, cfg)
		}
	}
	return configs, nil
}

func (m *MockProjectConfigRepository) FindByKey(ctx context.Context, projectID uuid.UUID, key string) (*project.ProjectConfig, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	k := projectID.String() + ":" + key
	cfg, ok := m.configs[k]
	if !ok {
		return nil, project.ErrConfigNotFound
	}
	return cfg, nil
}

func (m *MockProjectConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	return nil
}

func (m *MockProjectConfigRepository) Update(ctx context.Context, cfg *project.ProjectConfig) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	key := cfg.ProjectID().String() + ":" + cfg.Key()
	m.configs[key] = cfg
	return nil
}

// TestProjectService_CreateProject tests project creation
func TestProjectService_CreateProject(t *testing.T) {
	ctx := context.Background()
	projectRepo := NewMockProjectRepository()
	moduleRepo := NewMockModuleRepository()
	configRepo := NewMockProjectConfigRepository()
	service := NewProjectService(projectRepo, moduleRepo, configRepo)

	userID := uuid.New()

	tests := []struct {
		name    string
		req     *CreateProjectRequest
		setup   func()
		wantErr error
	}{
		{
			name: "successful creation",
			req: &CreateProjectRequest{
				Name:        "Test Project",
				Prefix:      "TEST",
				Description: "A test project",
			},
			setup:   func() {},
			wantErr: nil,
		},
		{
			name: "project name already exists",
			req: &CreateProjectRequest{
				Name:        "Existing Project",
				Prefix:      "NEWP",
				Description: "New project",
			},
			setup: func() {
				existing, _ := project.NewProject("Existing Project", "EXST", "Existing")
				projectRepo.projects[existing.ID()] = existing
				projectRepo.nameIndex["Existing Project"] = existing
				projectRepo.prefixIndex["EXST"] = existing
			},
			wantErr: project.ErrProjectNameDuplicate,
		},
		{
			name: "project prefix already exists",
			req: &CreateProjectRequest{
				Name:        "Another Project",
				Prefix:      "EXST",
				Description: "Another project",
			},
			setup: func() {
				// The EXST prefix was already added in the previous test case
				// No additional setup needed
			},
			wantErr: project.ErrProjectPrefixDuplicate,
		},
		{
			name: "invalid prefix format - too short",
			req: &CreateProjectRequest{
				Name:        "Invalid Prefix Project",
				Prefix:      "A",
				Description: "Invalid prefix",
			},
			setup:   func() {},
			wantErr: project.ErrInvalidProjectPrefix,
		},
		{
			name: "invalid prefix format - lowercase",
			req: &CreateProjectRequest{
				Name:        "Lowercase Prefix Project",
				Prefix:      "test",
				Description: "Lowercase prefix",
			},
			setup:   func() {},
			wantErr: project.ErrInvalidProjectPrefix,
		},
		{
			name: "empty name",
			req: &CreateProjectRequest{
				Name:        "",
				Prefix:      "EMPN",
				Description: "Empty name",
			},
			setup: func() {
				// Use a fresh repo to avoid state pollution
				projectRepo.projects = make(map[uuid.UUID]*project.Project)
				projectRepo.nameIndex = make(map[string]*project.Project)
				projectRepo.prefixIndex = make(map[string]*project.Project)
			},
			wantErr: errors.New("create project: project name cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			proj, err := service.CreateProject(ctx, tt.req, userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreateProject() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("CreateProject() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateProject() unexpected error: %v", err)
				return
			}

			if proj == nil {
				t.Error("CreateProject() returned nil project")
				return
			}

			if proj.Name() != tt.req.Name {
				t.Errorf("CreateProject() name = %v, want %v", proj.Name(), tt.req.Name)
			}
			if string(proj.Prefix()) != tt.req.Prefix {
				t.Errorf("CreateProject() prefix = %v, want %v", proj.Prefix(), tt.req.Prefix)
			}
		})
	}
}

// TestProjectService_GetProject tests project retrieval
func TestProjectService_GetProject(t *testing.T) {
	ctx := context.Background()
	projectRepo := NewMockProjectRepository()
	moduleRepo := NewMockModuleRepository()
	configRepo := NewMockProjectConfigRepository()
	service := NewProjectService(projectRepo, moduleRepo, configRepo)

	// Create test project
	testProject, _ := project.NewProject("Test Project", "TEST", "Description")
	projectRepo.projects[testProject.ID()] = testProject

	tests := []struct {
		name      string
		projectID uuid.UUID
		wantErr   error
	}{
		{
			name:      "successful retrieval",
			projectID: testProject.ID(),
			wantErr:   nil,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			wantErr:   project.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := service.GetProject(ctx, tt.projectID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetProject() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetProject() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetProject() unexpected error: %v", err)
				return
			}

			if detail == nil {
				t.Error("GetProject() returned nil detail")
				return
			}

			if detail.ID() != tt.projectID {
				t.Errorf("GetProject() ID = %v, want %v", detail.ID(), tt.projectID)
			}
		})
	}
}

// TestProjectService_ListProjects tests project listing
func TestProjectService_ListProjects(t *testing.T) {
	ctx := context.Background()
	projectRepo := NewMockProjectRepository()
	moduleRepo := NewMockModuleRepository()
	configRepo := NewMockProjectConfigRepository()
	service := NewProjectService(projectRepo, moduleRepo, configRepo)

	// Create test projects
	for i := 0; i < 5; i++ {
		prefixes := []string{"TEST", "PROJ", "CORE", "APIE", "AUTH"}
		proj, _ := project.NewProject("Project "+string(rune('A'+i)), prefixes[i], "Description")
		projectRepo.projects[proj.ID()] = proj
	}

	tests := []struct {
		name      string
		opts      ListOptions
		wantCount int
		wantErr   error
	}{
		{
			name: "list all projects",
			opts: ListOptions{
				Offset: 0,
				Limit:  10,
			},
			wantCount: 5,
			wantErr:   nil,
		},
		{
			name: "list with pagination",
			opts: ListOptions{
				Offset: 0,
				Limit:  2,
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name: "list with keyword search",
			opts: ListOptions{
				Offset:   0,
				Limit:    10,
				Keywords: "Project A",
			},
			wantCount: 1,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projects, total, err := service.ListProjects(ctx, tt.opts)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ListProjects() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ListProjects() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ListProjects() unexpected error: %v", err)
				return
			}

			if len(projects) != tt.wantCount {
				t.Errorf("ListProjects() count = %v, want %v", len(projects), tt.wantCount)
			}

			if total < int64(tt.wantCount) {
				t.Errorf("ListProjects() total = %v, want at least %v", total, tt.wantCount)
			}
		})
	}
}

// TestProjectService_UpdateProject tests project update
func TestProjectService_UpdateProject(t *testing.T) {
	ctx := context.Background()
	projectRepo := NewMockProjectRepository()
	moduleRepo := NewMockModuleRepository()
	configRepo := NewMockProjectConfigRepository()
	service := NewProjectService(projectRepo, moduleRepo, configRepo)

	// Create test project
	testProject, _ := project.NewProject("Original Name", "ORIG", "Original Description")
	projectRepo.projects[testProject.ID()] = testProject

	newName := "Updated Name"
	newDesc := "Updated Description"

	tests := []struct {
		name      string
		projectID uuid.UUID
		req       *UpdateProjectRequest
		wantErr   error
	}{
		{
			name:      "successful update name",
			projectID: testProject.ID(),
			req: &UpdateProjectRequest{
				Name: &newName,
			},
			wantErr: nil,
		},
		{
			name:      "successful update description",
			projectID: testProject.ID(),
			req: &UpdateProjectRequest{
				Description: &newDesc,
			},
			wantErr: nil,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			req: &UpdateProjectRequest{
				Name: &newName,
			},
			wantErr: project.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj, err := service.UpdateProject(ctx, tt.projectID, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UpdateProject() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("UpdateProject() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateProject() unexpected error: %v", err)
				return
			}

			if proj == nil {
				t.Error("UpdateProject() returned nil project")
				return
			}

			if tt.req.Name != nil && proj.Name() != *tt.req.Name {
				t.Errorf("UpdateProject() name = %v, want %v", proj.Name(), *tt.req.Name)
			}
		})
	}
}

// TestProjectService_DeleteProject tests project deletion
func TestProjectService_DeleteProject(t *testing.T) {
	ctx := context.Background()
	projectRepo := NewMockProjectRepository()
	moduleRepo := NewMockModuleRepository()
	configRepo := NewMockProjectConfigRepository()
	service := NewProjectService(projectRepo, moduleRepo, configRepo)

	// Create test project
	testProject, _ := project.NewProject("To Delete", "DEL", "Will be deleted")
	projectRepo.projects[testProject.ID()] = testProject

	tests := []struct {
		name      string
		projectID uuid.UUID
		wantErr   error
	}{
		{
			name:      "successful deletion",
			projectID: testProject.ID(),
			wantErr:   nil,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			wantErr:   project.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteProject(ctx, tt.projectID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("DeleteProject() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("DeleteProject() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteProject() unexpected error: %v", err)
				return
			}
		})
	}
}

// TestProjectService_CreateModule tests module creation
func TestProjectService_CreateModule(t *testing.T) {
	ctx := context.Background()
	projectRepo := NewMockProjectRepository()
	moduleRepo := NewMockModuleRepository()
	configRepo := NewMockProjectConfigRepository()
	service := NewProjectService(projectRepo, moduleRepo, configRepo)

	// Create test project
	testProject, _ := project.NewProject("Test Project", "TEST", "Description")
	projectRepo.projects[testProject.ID()] = testProject

	userID := uuid.New()

	tests := []struct {
		name      string
		projectID uuid.UUID
		req       *CreateModuleRequest
		setup     func()
		wantErr   error
	}{
		{
			name:      "successful creation",
			projectID: testProject.ID(),
			req: &CreateModuleRequest{
				Name:         "User Module",
				Abbreviation: "USER",
				Description:  "User management module",
			},
			setup:   func() {},
			wantErr: nil,
		},
		{
			name:      "abbreviation already exists",
			projectID: testProject.ID(),
			req: &CreateModuleRequest{
				Name:         "Another User Module",
				Abbreviation: "USER",
				Description:  "Duplicate abbreviation",
			},
			setup:   func() {},
			wantErr: project.ErrModuleAbbrevDuplicate,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			req: &CreateModuleRequest{
				Name:         "Orphan Module",
				Abbreviation: "ORPH",
				Description:  "No parent project",
			},
			setup:   func() {},
			wantErr: project.ErrProjectNotFound,
		},
		{
			name:      "invalid abbreviation format",
			projectID: testProject.ID(),
			req: &CreateModuleRequest{
				Name:         "Invalid Module",
				Abbreviation: "toolong",
				Description:  "Invalid abbreviation",
			},
			setup:   func() {},
			wantErr: project.ErrInvalidModuleAbbrev,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			mod, err := service.CreateModule(ctx, tt.projectID, tt.req, userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreateModule() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("CreateModule() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateModule() unexpected error: %v", err)
				return
			}

			if mod == nil {
				t.Error("CreateModule() returned nil module")
				return
			}

			if mod.Name() != tt.req.Name {
				t.Errorf("CreateModule() name = %v, want %v", mod.Name(), tt.req.Name)
			}
			if string(mod.Abbreviation()) != tt.req.Abbreviation {
				t.Errorf("CreateModule() abbreviation = %v, want %v", mod.Abbreviation(), tt.req.Abbreviation)
			}
		})
	}
}

// TestProjectService_SetConfig tests configuration setting
func TestProjectService_SetConfig(t *testing.T) {
	ctx := context.Background()
	projectRepo := NewMockProjectRepository()
	moduleRepo := NewMockModuleRepository()
	configRepo := NewMockProjectConfigRepository()
	service := NewProjectService(projectRepo, moduleRepo, configRepo)

	// Create test project
	testProject, _ := project.NewProject("Test Project", "TEST", "Description")
	projectRepo.projects[testProject.ID()] = testProject

	tests := []struct {
		name      string
		projectID uuid.UUID
		key       string
		value     map[string]any
		wantErr   error
	}{
		{
			name:      "successful set",
			projectID: testProject.ID(),
			key:       "llm_config",
			value: map[string]any{
				"model":       "deepseek-chat",
				"temperature": 0.7,
			},
			wantErr: nil,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			key:       "llm_config",
			value:     map[string]any{},
			wantErr:   project.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SetConfig(ctx, tt.projectID, tt.key, tt.value)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("SetConfig() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("SetConfig() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("SetConfig() unexpected error: %v", err)
				return
			}

			// Verify config was saved
			cfg, err := service.GetConfig(ctx, tt.projectID, tt.key)
			if err != nil {
				t.Errorf("GetConfig() error after SetConfig: %v", err)
				return
			}
			if cfg.Key() != tt.key {
				t.Errorf("GetConfig() key = %v, want %v", cfg.Key(), tt.key)
			}
		})
	}
}
