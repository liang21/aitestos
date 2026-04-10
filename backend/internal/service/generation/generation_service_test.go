// Package generation provides test case generation services
package generation

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

// MockGenerationTaskRepository for testing
type MockGenTaskRepo struct {
	tasks      map[uuid.UUID]*generation.GenerationTask
	projectIdx map[uuid.UUID][]*generation.GenerationTask
	saveErr    error
	findErr    error
}

func NewMockGenTaskRepo() *MockGenTaskRepo {
	return &MockGenTaskRepo{
		tasks:      make(map[uuid.UUID]*generation.GenerationTask),
		projectIdx: make(map[uuid.UUID][]*generation.GenerationTask),
	}
}

func (m *MockGenTaskRepo) Save(ctx context.Context, task *generation.GenerationTask) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.tasks[task.ID()] = task
	m.projectIdx[task.ProjectID()] = append(m.projectIdx[task.ProjectID()], task)
	return nil
}

func (m *MockGenTaskRepo) FindByID(ctx context.Context, id uuid.UUID) (*generation.GenerationTask, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	task, ok := m.tasks[id]
	if !ok {
		return nil, generation.ErrTaskNotFound
	}
	return task, nil
}

func (m *MockGenTaskRepo) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts generation.QueryOptions) ([]*generation.GenerationTask, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	tasks, ok := m.projectIdx[projectID]
	if !ok {
		return []*generation.GenerationTask{}, nil
	}
	return tasks, nil
}

func (m *MockGenTaskRepo) FindByStatus(ctx context.Context, status generation.TaskStatus, opts generation.QueryOptions) ([]*generation.GenerationTask, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	tasks := make([]*generation.GenerationTask, 0)
	for _, t := range m.tasks {
		if t.Status() == status {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (m *MockGenTaskRepo) Update(ctx context.Context, task *generation.GenerationTask) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.tasks[task.ID()] = task
	return nil
}

func (m *MockGenTaskRepo) FindByUserID(ctx context.Context, userID uuid.UUID, opts generation.QueryOptions) ([]*generation.GenerationTask, error) {
	return []*generation.GenerationTask{}, nil
}

func (m *MockGenTaskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// MockCaseDraftRepository for testing
type MockDraftRepo struct {
	drafts  map[uuid.UUID]*generation.GeneratedCaseDraft
	taskIdx map[uuid.UUID][]*generation.GeneratedCaseDraft
	saveErr error
	findErr error
}

func NewMockDraftRepo() *MockDraftRepo {
	return &MockDraftRepo{
		drafts:  make(map[uuid.UUID]*generation.GeneratedCaseDraft),
		taskIdx: make(map[uuid.UUID][]*generation.GeneratedCaseDraft),
	}
}

func (m *MockDraftRepo) Save(ctx context.Context, draft *generation.GeneratedCaseDraft) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.drafts[draft.ID()] = draft
	m.taskIdx[draft.TaskID()] = append(m.taskIdx[draft.TaskID()], draft)
	return nil
}

func (m *MockDraftRepo) FindByID(ctx context.Context, id uuid.UUID) (*generation.GeneratedCaseDraft, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	draft, ok := m.drafts[id]
	if !ok {
		return nil, generation.ErrDraftNotFound
	}
	return draft, nil
}

func (m *MockDraftRepo) FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*generation.GeneratedCaseDraft, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	drafts, ok := m.taskIdx[taskID]
	if !ok {
		return []*generation.GeneratedCaseDraft{}, nil
	}
	return drafts, nil
}

func (m *MockDraftRepo) FindByStatus(ctx context.Context, status generation.DraftStatus, opts generation.QueryOptions) ([]*generation.GeneratedCaseDraft, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	drafts := make([]*generation.GeneratedCaseDraft, 0)
	for _, d := range m.drafts {
		if d.Status() == status {
			drafts = append(drafts, d)
		}
	}
	return drafts, nil
}

func (m *MockDraftRepo) Update(ctx context.Context, draft *generation.GeneratedCaseDraft) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.drafts[draft.ID()] = draft
	return nil
}

func (m *MockDraftRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.drafts, id)
	return nil
}

func (m *MockDraftRepo) FindByTaskIDAndStatus(ctx context.Context, taskID uuid.UUID, status generation.DraftStatus) ([]*generation.GeneratedCaseDraft, error) {
	return []*generation.GeneratedCaseDraft{}, nil
}

func (m *MockDraftRepo) BatchUpdateStatus(ctx context.Context, draftIDs []uuid.UUID, status generation.DraftStatus, moduleID uuid.UUID) error {
	return nil
}

func (m *MockDraftRepo) CountByTaskIDAndStatus(ctx context.Context, taskID uuid.UUID, status generation.DraftStatus) (int64, error) {
	return 0, nil
}

func (m *MockDraftRepo) DeleteByTaskID(ctx context.Context, taskID uuid.UUID) error {
	return nil
}

// MockGenModuleRepo for testing
type MockGenModuleRepo struct {
	modules map[uuid.UUID]*project.Module
}

func NewMockGenModuleRepo() *MockGenModuleRepo {
	return &MockGenModuleRepo{
		modules: make(map[uuid.UUID]*project.Module),
	}
}

func (m *MockGenModuleRepo) FindByID(ctx context.Context, id uuid.UUID) (*project.Module, error) {
	mod, ok := m.modules[id]
	if !ok {
		return nil, errors.New("module not found")
	}
	return mod, nil
}

func (m *MockGenModuleRepo) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*project.Module, error) {
	mods := make([]*project.Module, 0)
	for _, mod := range m.modules {
		if mod.ProjectID() == projectID {
			mods = append(mods, mod)
		}
	}
	return mods, nil
}

// MockGenProjectRepo for testing
type MockGenProjectRepo struct {
	projects map[uuid.UUID]*project.Project
}

func NewMockGenProjectRepo() *MockGenProjectRepo {
	return &MockGenProjectRepo{
		projects: make(map[uuid.UUID]*project.Project),
	}
}

func (m *MockGenProjectRepo) FindByID(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	proj, ok := m.projects[id]
	if !ok {
		return nil, errors.New("project not found")
	}
	return proj, nil
}

// MockCaseRepo for testing
type MockCaseRepo struct {
	cases map[uuid.UUID]*testcase.TestCase
}

func NewMockCaseRepo() *MockCaseRepo {
	return &MockCaseRepo{
		cases: make(map[uuid.UUID]*testcase.TestCase),
	}
}

func (m *MockCaseRepo) Save(ctx context.Context, tc *testcase.TestCase) error {
	m.cases[tc.ID()] = tc
	return nil
}

func (m *MockCaseRepo) FindByID(ctx context.Context, id uuid.UUID) (*testcase.TestCase, error) {
	tc, ok := m.cases[id]
	if !ok {
		return nil, errors.New("test case not found")
	}
	return tc, nil
}

func (m *MockCaseRepo) FindByNumber(ctx context.Context, number testcase.CaseNumber) (*testcase.TestCase, error) {
	for _, tc := range m.cases {
		if tc.Number() == number {
			return tc, nil
		}
	}
	return nil, errors.New("test case not found")
}

func (m *MockCaseRepo) FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts testcase.QueryOptions) ([]*testcase.TestCase, error) {
	result := make([]*testcase.TestCase, 0)
	for _, tc := range m.cases {
		if tc.ModuleID() == moduleID {
			result = append(result, tc)
		}
	}
	return result, nil
}

func (m *MockCaseRepo) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts testcase.QueryOptions) ([]*testcase.TestCase, error) {
	return []*testcase.TestCase{}, nil
}

func (m *MockCaseRepo) Update(ctx context.Context, tc *testcase.TestCase) error {
	m.cases[tc.ID()] = tc
	return nil
}

func (m *MockCaseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.cases, id)
	return nil
}

func (m *MockCaseRepo) CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error) {
	return 0, nil
}

// TestGenerationService_CreateTask tests task creation
func TestGenerationService_CreateTask(t *testing.T) {
	ctx := context.Background()
	taskRepo := NewMockGenTaskRepo()
	draftRepo := NewMockDraftRepo()
	ragSvc := NewMockRAGService()
	llmSvc := NewMockLLMService()
	moduleRepo := NewMockGenModuleRepo()
	projectRepo := NewMockGenProjectRepo()
	caseRepo := NewMockCaseRepo()

	service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo)

	projectID := uuid.New()
	moduleID := uuid.New()
	userID := uuid.New()

	// Create test module
	testModule, _ := project.NewModule(projectID, "User Module", "USER", "Description", userID)
	moduleRepo.modules[moduleID] = testModule

	tests := []struct {
		name    string
		req     *CreateTaskRequest
		setup   func()
		wantErr error
	}{
		{
			name: "successful task creation",
			req: &CreateTaskRequest{
				ProjectID:  projectID,
				ModuleID:   moduleID,
				Prompt:     "Generate test cases for user login functionality including positive and negative scenarios",
				CaseCount:  3,
				SceneTypes: []string{"positive", "negative"},
				Priority:   "P1",
			},
			setup:   func() {},
			wantErr: nil,
		},
		{
			name: "prompt too short",
			req: &CreateTaskRequest{
				ProjectID: projectID,
				ModuleID:  moduleID,
				Prompt:    "short",
			},
			setup:   func() {},
			wantErr: errors.New("prompt must be at least 10 characters"),
		},
		{
			name: "nil project ID",
			req: &CreateTaskRequest{
				ProjectID: uuid.Nil,
				ModuleID:  moduleID,
				Prompt:    "Generate test cases",
			},
			setup:   func() {},
			wantErr: errors.New("project ID cannot be nil"),
		},
		{
			name: "invalid case count - too many",
			req: &CreateTaskRequest{
				ProjectID: projectID,
				ModuleID:  moduleID,
				Prompt:    "Generate test cases",
				CaseCount: 25,
			},
			setup:   func() {},
			wantErr: errors.New("case count must be between 1 and 20"),
		},
		{
			name: "invalid scene type",
			req: &CreateTaskRequest{
				ProjectID:  projectID,
				ModuleID:   moduleID,
				Prompt:     "Generate test cases",
				SceneTypes: []string{"invalid_type"},
			},
			setup:   func() {},
			wantErr: errors.New("invalid scene type: invalid_type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			task, err := service.CreateTask(ctx, tt.req, userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreateTask() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("CreateTask() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateTask() unexpected error: %v", err)
				return
			}

			if task == nil {
				t.Error("CreateTask() returned nil task")
				return
			}

			if task.Prompt() != tt.req.Prompt {
				t.Errorf("CreateTask() prompt = %v, want %v", task.Prompt(), tt.req.Prompt)
			}
			if task.Status() != generation.TaskPending {
				t.Errorf("CreateTask() status = %v, want %v", task.Status(), generation.TaskPending)
			}
		})
	}
}

// TestGenerationService_ConfirmDraft tests draft confirmation
func TestGenerationService_ConfirmDraft(t *testing.T) {
	moduleID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		setup   func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft)
		req     func(draft *generation.GeneratedCaseDraft) *ConfirmDraftRequest
		wantErr error
	}{
		{
			name: "successful confirmation",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				// Create test module and project with matching IDs
				testProject, _ := project.NewProject("Test Project", "TEST", "Description")
				testModule, _ := project.NewModule(testProject.ID(), "Test Module", "TEST", "Description", userID)
				moduleRepo.modules[moduleID] = testModule
				projectRepo.projects[testProject.ID()] = testProject

				// Create test draft
				task, _ := generation.NewGenerationTask(testProject.ID(), moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft, _ := generation.NewGeneratedCaseDraft(
					task.ID(),
					"Test Case Title",
					testcase.Preconditions{"Precondition 1"},
					testcase.Steps{"Step 1", "Step 2"},
					testcase.ExpectedResult{"result": "success"},
					testcase.CaseTypeFunctionality,
					testcase.PriorityP1,
				)
				draftRepo.drafts[draft.ID()] = draft

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, draft
			},
			req: func(draft *generation.GeneratedCaseDraft) *ConfirmDraftRequest {
				return &ConfirmDraftRequest{
					DraftID:  draft.ID(),
					ModuleID: moduleID,
				}
			},
			wantErr: nil,
		},
		{
			name: "draft not found",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				testProject, _ := project.NewProject("Test Project", "TEST", "Description")
				projectRepo.projects[testProject.ID()] = testProject

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, nil
			},
			req: func(draft *generation.GeneratedCaseDraft) *ConfirmDraftRequest {
				return &ConfirmDraftRequest{
					DraftID:  uuid.New(),
					ModuleID: moduleID,
				}
			},
			wantErr: generation.ErrDraftNotFound,
		},
		{
			name: "draft already confirmed",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				testProject, _ := project.NewProject("Test Project", "TEST", "Description")
				testModule, _ := project.NewModule(testProject.ID(), "Test Module", "TEST", "Description", userID)
				moduleRepo.modules[moduleID] = testModule
				projectRepo.projects[testProject.ID()] = testProject

				task, _ := generation.NewGenerationTask(testProject.ID(), moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft, _ := generation.NewGeneratedCaseDraft(
					task.ID(),
					"Test Case Title",
					testcase.Preconditions{"Precondition 1"},
					testcase.Steps{"Step 1", "Step 2"},
					testcase.ExpectedResult{"result": "success"},
					testcase.CaseTypeFunctionality,
					testcase.PriorityP1,
				)
				draft.Confirm(moduleID) // Pre-confirm the draft
				draftRepo.drafts[draft.ID()] = draft

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, draft
			},
			req: func(draft *generation.GeneratedCaseDraft) *ConfirmDraftRequest {
				return &ConfirmDraftRequest{
					DraftID:  draft.ID(),
					ModuleID: moduleID,
				}
			},
			wantErr: generation.ErrDraftAlreadyConfirmed,
		},
		{
			name: "module not found",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				testProject, _ := project.NewProject("Test Project", "TEST", "Description")
				projectRepo.projects[testProject.ID()] = testProject

				task, _ := generation.NewGenerationTask(testProject.ID(), moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft, _ := generation.NewGeneratedCaseDraft(
					task.ID(),
					"Test Case Title",
					testcase.Preconditions{"Precondition 1"},
					testcase.Steps{"Step 1", "Step 2"},
					testcase.ExpectedResult{"result": "success"},
					testcase.CaseTypeFunctionality,
					testcase.PriorityP1,
				)
				draftRepo.drafts[draft.ID()] = draft

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, draft
			},
			req: func(draft *generation.GeneratedCaseDraft) *ConfirmDraftRequest {
				return &ConfirmDraftRequest{
					DraftID:  draft.ID(),
					ModuleID: uuid.New(), // Non-existent module
				}
			},
			wantErr: errors.New("module not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, draft := tt.setup()
			req := tt.req(draft)

			tc, err := service.ConfirmDraft(context.Background(), req, userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ConfirmDraft() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("ConfirmDraft() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ConfirmDraft() unexpected error: %v", err)
				return
			}

			if tc == nil {
				t.Error("ConfirmDraft() returned nil test case")
				return
			}

			if tc.ModuleID() != req.ModuleID {
				t.Errorf("ConfirmDraft() moduleID = %v, want %v", tc.ModuleID(), req.ModuleID)
			}
		})
	}
}

// TestGenerationService_RejectDraft tests draft rejection
func TestGenerationService_RejectDraft(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		setup   func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft)
		req     func(draft *generation.GeneratedCaseDraft) *RejectDraftRequest
		wantErr error
	}{
		{
			name: "successful rejection",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				projectID := uuid.New()
				moduleID := uuid.New()

				task, _ := generation.NewGenerationTask(projectID, moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft, _ := generation.NewGeneratedCaseDraft(
					task.ID(),
					"Draft to Reject",
					testcase.Preconditions{"Precondition"},
					testcase.Steps{"Step"},
					testcase.ExpectedResult{},
					testcase.CaseTypeFunctionality,
					testcase.PriorityP1,
				)
				draftRepo.drafts[draft.ID()] = draft

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, draft
			},
			req: func(draft *generation.GeneratedCaseDraft) *RejectDraftRequest {
				return &RejectDraftRequest{
					DraftID:  draft.ID(),
					Reason:   generation.ReasonDuplicate,
					Feedback: "This case duplicates an existing one",
				}
			},
			wantErr: nil,
		},
		{
			name: "draft not found",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, nil
			},
			req: func(draft *generation.GeneratedCaseDraft) *RejectDraftRequest {
				return &RejectDraftRequest{
					DraftID:  uuid.New(),
					Reason:   generation.ReasonOther,
					Feedback: "Not found",
				}
			},
			wantErr: generation.ErrDraftNotFound,
		},
		{
			name: "draft already confirmed",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				projectID := uuid.New()
				moduleID := uuid.New()

				task, _ := generation.NewGenerationTask(projectID, moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft, _ := generation.NewGeneratedCaseDraft(
					task.ID(),
					"Draft Already Confirmed",
					testcase.Preconditions{"Precondition"},
					testcase.Steps{"Step"},
					testcase.ExpectedResult{},
					testcase.CaseTypeFunctionality,
					testcase.PriorityP1,
				)
				draft.Confirm(moduleID) // Pre-confirm the draft
				draftRepo.drafts[draft.ID()] = draft

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, draft
			},
			req: func(draft *generation.GeneratedCaseDraft) *RejectDraftRequest {
				return &RejectDraftRequest{
					DraftID:  draft.ID(),
					Reason:   generation.ReasonOther,
					Feedback: "Already processed",
				}
			},
			wantErr: generation.ErrDraftAlreadyConfirmed,
		},
		{
			name: "invalid reason",
			setup: func() (*GenerationServiceImpl, *generation.GeneratedCaseDraft) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				projectID := uuid.New()
				moduleID := uuid.New()

				task, _ := generation.NewGenerationTask(projectID, moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft, _ := generation.NewGeneratedCaseDraft(
					task.ID(),
					"Draft for Invalid Reason Test",
					testcase.Preconditions{"Precondition"},
					testcase.Steps{"Step"},
					testcase.ExpectedResult{},
					testcase.CaseTypeFunctionality,
					testcase.PriorityP1,
				)
				draftRepo.drafts[draft.ID()] = draft

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo).(*GenerationServiceImpl)
				return service, draft
			},
			req: func(draft *generation.GeneratedCaseDraft) *RejectDraftRequest {
				return &RejectDraftRequest{
					DraftID:  draft.ID(),
					Reason:   "invalid_reason",
					Feedback: "Invalid",
				}
			},
			wantErr: errors.New("invalid rejection reason"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, draft := tt.setup()
			req := tt.req(draft)

			err := service.RejectDraft(context.Background(), req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("RejectDraft() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("RejectDraft() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("RejectDraft() unexpected error: %v", err)
				return
			}
		})
	}
}

// TestGenerationService_BatchConfirm tests batch confirmation
func TestGenerationService_BatchConfirm(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		setup   func() (GenerationService, []uuid.UUID, uuid.UUID)
		req     func(draftIDs []uuid.UUID, moduleID uuid.UUID) *BatchConfirmRequest
		wantErr error
		check   func(result *BatchConfirmResult) error
	}{
		{
			name: "successful batch confirmation",
			setup: func() (GenerationService, []uuid.UUID, uuid.UUID) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				testProject, _ := project.NewProject("Test Project", "TEST", "Description")
				moduleID := uuid.New()
				testModule, _ := project.NewModule(testProject.ID(), "Test Module", "TEST", "Description", userID)
				moduleRepo.modules[moduleID] = testModule
				projectRepo.projects[testProject.ID()] = testProject

				task, _ := generation.NewGenerationTask(testProject.ID(), moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft1, _ := generation.NewGeneratedCaseDraft(task.ID(), "Draft 1", nil, testcase.Steps{"Step"}, testcase.ExpectedResult{}, testcase.CaseTypeFunctionality, testcase.PriorityP1)
				draft2, _ := generation.NewGeneratedCaseDraft(task.ID(), "Draft 2", nil, testcase.Steps{"Step"}, testcase.ExpectedResult{}, testcase.CaseTypeFunctionality, testcase.PriorityP1)
				draftRepo.drafts[draft1.ID()] = draft1
				draftRepo.drafts[draft2.ID()] = draft2

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo)
				return service, []uuid.UUID{draft1.ID(), draft2.ID()}, moduleID
			},
			req: func(draftIDs []uuid.UUID, moduleID uuid.UUID) *BatchConfirmRequest {
				return &BatchConfirmRequest{
					DraftIDs: draftIDs,
					ModuleID: moduleID,
				}
			},
			wantErr: nil,
			check: func(result *BatchConfirmResult) error {
				if result.SuccessCount != 2 {
					return fmt.Errorf("expected 2 successful confirmations, got %d", result.SuccessCount)
				}
				if len(result.FailedIDs) != 0 {
					return fmt.Errorf("expected no failed IDs, got %d", len(result.FailedIDs))
				}
				return nil
			},
		},
		{
			name: "partial success - one draft already confirmed",
			setup: func() (GenerationService, []uuid.UUID, uuid.UUID) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				testProject, _ := project.NewProject("Test Project", "TEST", "Description")
				moduleID := uuid.New()
				testModule, _ := project.NewModule(testProject.ID(), "Test Module", "TEST", "Description", userID)
				moduleRepo.modules[moduleID] = testModule
				projectRepo.projects[testProject.ID()] = testProject

				task, _ := generation.NewGenerationTask(testProject.ID(), moduleID, "Test prompt", userID)
				taskRepo.tasks[task.ID()] = task

				draft1, _ := generation.NewGeneratedCaseDraft(task.ID(), "Draft 1", nil, testcase.Steps{"Step"}, testcase.ExpectedResult{}, testcase.CaseTypeFunctionality, testcase.PriorityP1)
				draft2, _ := generation.NewGeneratedCaseDraft(task.ID(), "Draft 2", nil, testcase.Steps{"Step"}, testcase.ExpectedResult{}, testcase.CaseTypeFunctionality, testcase.PriorityP1)
				draft1.Confirm(moduleID) // Pre-confirm draft1
				draftRepo.drafts[draft1.ID()] = draft1
				draftRepo.drafts[draft2.ID()] = draft2

				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo)
				return service, []uuid.UUID{draft1.ID(), draft2.ID()}, moduleID
			},
			req: func(draftIDs []uuid.UUID, moduleID uuid.UUID) *BatchConfirmRequest {
				return &BatchConfirmRequest{
					DraftIDs: draftIDs,
					ModuleID: moduleID,
				}
			},
			wantErr: nil,
			check: func(result *BatchConfirmResult) error {
				if result.SuccessCount != 1 {
					return fmt.Errorf("expected 1 successful confirmation, got %d", result.SuccessCount)
				}
				if len(result.FailedIDs) != 1 {
					return fmt.Errorf("expected 1 failed ID, got %d", len(result.FailedIDs))
				}
				return nil
			},
		},
		{
			name: "empty draft list",
			setup: func() (GenerationService, []uuid.UUID, uuid.UUID) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				moduleID := uuid.New()
				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo)
				return service, []uuid.UUID{}, moduleID
			},
			req: func(draftIDs []uuid.UUID, moduleID uuid.UUID) *BatchConfirmRequest {
				return &BatchConfirmRequest{
					DraftIDs: draftIDs,
					ModuleID: moduleID,
				}
			},
			wantErr: errors.New("draft IDs cannot be empty"),
		},
		{
			name: "too many drafts",
			setup: func() (GenerationService, []uuid.UUID, uuid.UUID) {
				taskRepo := NewMockGenTaskRepo()
				draftRepo := NewMockDraftRepo()
				ragSvc := NewMockRAGService()
				llmSvc := NewMockLLMService()
				moduleRepo := NewMockGenModuleRepo()
				projectRepo := NewMockGenProjectRepo()
				caseRepo := NewMockCaseRepo()

				moduleID := uuid.New()
				service := NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepo, caseRepo)
				return service, make([]uuid.UUID, 51), moduleID
			},
			req: func(draftIDs []uuid.UUID, moduleID uuid.UUID) *BatchConfirmRequest {
				return &BatchConfirmRequest{
					DraftIDs: draftIDs,
					ModuleID: moduleID,
				}
			},
			wantErr: errors.New("cannot confirm more than 50 drafts at once"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, draftIDs, moduleID := tt.setup()
			req := tt.req(draftIDs, moduleID)

			result, err := service.BatchConfirm(context.Background(), req, userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("BatchConfirm() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("BatchConfirm() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("BatchConfirm() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("BatchConfirm() returned nil result")
				return
			}

			if tt.check != nil {
				if checkErr := tt.check(result); checkErr != nil {
					t.Errorf("BatchConfirm() validation failed: %v", checkErr)
				}
			}
		})
	}
}

// MockRAGService for testing
type MockRAGService struct{}

func NewMockRAGService() *MockRAGService {
	return &MockRAGService{}
}

func (m *MockRAGService) Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResult, error) {
	return &RetrieveResult{Chunks: []*RetrievedChunk{}}, nil
}

func (m *MockRAGService) CalculateConfidence(chunks []*RetrievedChunk) testcase.Confidence {
	return testcase.ConfidenceMedium
}

// MockLLMService for testing
type MockLLMService struct{}

func NewMockLLMService() *MockLLMService {
	return &MockLLMService{}
}

func (m *MockLLMService) GenerateCases(ctx context.Context, req *GenerateCasesRequest) (*GenerateCasesResponse, error) {
	return &GenerateCasesResponse{
		Cases: []*GeneratedCase{
			{
				Title:         "Generated Test Case",
				Preconditions: []string{"System is running"},
				Steps:         []string{"Step 1", "Step 2"},
				Expected:      map[string]any{"status": "success"},
				CaseType:      "functionality",
				Priority:      "P1",
			},
		},
	}, nil
}

func (m *MockLLMService) GenerateEmbedding(ctx context.Context, req *GenerateEmbeddingRequest) (*GenerateEmbeddingResponse, error) {
	 embedding := make([]byte, 12)
    for i := 0; i < 12; i++ {
        embedding[i] = byte(i % 256)
    }
    return &GenerateEmbeddingResponse{Embedding: embedding}, nil
}

func (m *MockLLMService) GetModelVersion() string {
	return "mock-model-v1"
}
