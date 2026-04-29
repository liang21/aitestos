// Package generation provides test case generation services
package generation

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

// CreateTaskRequest contains generation task creation data
type CreateTaskRequest struct {
	ProjectID  uuid.UUID `json:"project_id" validate:"required"`
	ModuleID   uuid.UUID `json:"module_id" validate:"required"`
	Prompt     string    `json:"prompt" validate:"required,min=10"`
	CaseCount  int       `json:"case_count"`  // default 5, max 20
	SceneTypes []string  `json:"scene_types"` // positive, negative, boundary
	Priority   string    `json:"priority"`    // P0-P3
	CaseType   string    `json:"case_type"`   // functionality, performance, api, ui, security
}

// ListTaskOptions contains options for listing generation tasks
type ListTaskOptions struct {
	Offset   int
	Limit    int
	Status   string
	Keywords string
	ModuleID uuid.UUID
}

// ConfirmDraftRequest contains draft confirmation data
type ConfirmDraftRequest struct {
	DraftID  uuid.UUID `json:"draft_id" validate:"required"`
	ModuleID uuid.UUID `json:"module_id" validate:"required"`
}

// RejectDraftRequest contains draft rejection data
type RejectDraftRequest struct {
	DraftID  uuid.UUID                  `json:"draft_id" validate:"required"`
	Reason   generation.RejectionReason `json:"reason" validate:"required"`
	Feedback string                     `json:"feedback"`
}

// BatchConfirmRequest contains batch confirmation data
type BatchConfirmRequest struct {
	DraftIDs []uuid.UUID `json:"draft_ids" validate:"required,min=1,max=50"`
	ModuleID uuid.UUID   `json:"module_id" validate:"required"`
}

// BatchConfirmResult contains batch confirmation results
type BatchConfirmResult struct {
	SuccessCount int
	FailedCount  int
	SuccessIDs   []uuid.UUID
	FailedIDs    []uuid.UUID
	Errors       []string
}

// ListAllDraftsOptions contains options for listing all drafts
type ListAllDraftsOptions struct {
	Offset    int
	Limit     int
	Status    string
	ProjectID uuid.UUID
	TaskID    uuid.UUID
}

// GenerationService provides test case generation operations
type GenerationService interface {
	// CreateTask creates a new generation task
	CreateTask(ctx context.Context, req *CreateTaskRequest, userID uuid.UUID) (*generation.GenerationTask, error)

	// GetTask retrieves a generation task by ID
	GetTask(ctx context.Context, id uuid.UUID) (*generation.GenerationTask, error)

	// ListTasks lists generation tasks with pagination
	ListTasks(ctx context.Context, projectID uuid.UUID, opts ListTaskOptions) ([]*generation.GenerationTask, int64, error)

	// GetDrafts retrieves all drafts for a task
	GetDrafts(ctx context.Context, taskID uuid.UUID) ([]*generation.GeneratedCaseDraft, error)

	// ListAllDrafts lists all drafts with filters
	ListAllDrafts(ctx context.Context, opts ListAllDraftsOptions) ([]*generation.GeneratedCaseDraft, int64, error)

	// ConfirmDraft confirms a draft and creates a test case
	ConfirmDraft(ctx context.Context, req *ConfirmDraftRequest, userID uuid.UUID) (*testcase.TestCase, error)

	// RejectDraft rejects a draft with reason
	RejectDraft(ctx context.Context, req *RejectDraftRequest) error

	// BatchConfirm confirms multiple drafts at once
	BatchConfirm(ctx context.Context, req *BatchConfirmRequest, userID uuid.UUID) (*BatchConfirmResult, error)

	// ProcessTask executes the generation workflow
	ProcessTask(ctx context.Context, taskID uuid.UUID) error
}

// GenerationServiceImpl implements GenerationService
type GenerationServiceImpl struct {
	taskRepo    generation.GenerationTaskRepository
	draftRepo   generation.CaseDraftRepository
	ragService  RAGService
	llmService  LLMService
	moduleRepo  ModuleRepository
	projectRepo ProjectRepository
	caseRepo    testcase.TestCaseRepository
}

// ModuleRepository interface for generation service
type ModuleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*project.Module, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*project.Module, error)
}

// ProjectRepository interface for generation service
type ProjectRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*project.Project, error)
}

// NewGenerationService creates a new GenerationService instance
func NewGenerationService(
	taskRepo generation.GenerationTaskRepository,
	draftRepo generation.CaseDraftRepository,
	ragService RAGService,
	llmService LLMService,
	moduleRepo ModuleRepository,
	projectRepo ProjectRepository,
	caseRepo testcase.TestCaseRepository,
) GenerationService {
	return &GenerationServiceImpl{
		taskRepo:    taskRepo,
		draftRepo:   draftRepo,
		ragService:  ragService,
		llmService:  llmService,
		moduleRepo:  moduleRepo,
		projectRepo: projectRepo,
		caseRepo:    caseRepo,
	}
}

// CreateTask creates a new generation task
func (s *GenerationServiceImpl) CreateTask(ctx context.Context, req *CreateTaskRequest, userID uuid.UUID) (*generation.GenerationTask, error) {
	// Validate project ID
	if req.ProjectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}

	// Validate prompt
	if len(req.Prompt) < 10 {
		return nil, errors.New("prompt must be at least 10 characters")
	}

	// Validate case count
	caseCount := req.CaseCount
	if caseCount <= 0 {
		caseCount = 5
	}
	if caseCount > 20 {
		return nil, errors.New("case count must be between 1 and 20")
	}

	// Validate scene types
	validSceneTypes := map[string]bool{
		"positive": true, "negative": true, "boundary": true,
	}
	for _, st := range req.SceneTypes {
		if !validSceneTypes[st] {
			return nil, fmt.Errorf("invalid scene type: %s", st)
		}
	}

	// Create task
	task, err := generation.NewGenerationTask(req.ProjectID, req.ModuleID, req.Prompt, userID)
	if err != nil {
		return nil, fmt.Errorf("create generation task: %w", err)
	}

	// Save task
	if err := s.taskRepo.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("save generation task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a generation task by ID
func (s *GenerationServiceImpl) GetTask(ctx context.Context, id uuid.UUID) (*generation.GenerationTask, error) {
	task, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find generation task: %w", err)
	}
	return task, nil
}

// ListTasks lists generation tasks with pagination
func (s *GenerationServiceImpl) ListTasks(ctx context.Context, projectID uuid.UUID, opts ListTaskOptions) ([]*generation.GenerationTask, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}

	queryOpts := generation.QueryOptions{
		Offset:   opts.Offset,
		Limit:    opts.Limit,
		Keywords: opts.Keywords,
	}

	tasks, err := s.taskRepo.FindByProjectID(ctx, projectID, queryOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("list generation tasks: %w", err)
	}

	// Filter by ModuleID if specified
	if opts.ModuleID != uuid.Nil {
		filtered := make([]*generation.GenerationTask, 0)
		for _, task := range tasks {
			if task.ModuleID() == opts.ModuleID {
				filtered = append(filtered, task)
			}
		}
		tasks = filtered
	}

	// Filter by Status if specified
	if opts.Status != "" {
		status, statusErr := generation.ParseTaskStatus(opts.Status)
		if statusErr == nil {
			filtered := make([]*generation.GenerationTask, 0)
			for _, task := range tasks {
				if task.Status() == status {
					filtered = append(filtered, task)
				}
			}
			tasks = filtered
		}
	}

	return tasks, int64(len(tasks)), nil
}

// GetDrafts retrieves all drafts for a task
func (s *GenerationServiceImpl) GetDrafts(ctx context.Context, taskID uuid.UUID) ([]*generation.GeneratedCaseDraft, error) {
	drafts, err := s.draftRepo.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("find drafts: %w", err)
	}
	return drafts, nil
}

// ListAllDrafts lists all drafts with filters
func (s *GenerationServiceImpl) ListAllDrafts(ctx context.Context, opts ListAllDraftsOptions) ([]*generation.GeneratedCaseDraft, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}

	var drafts []*generation.GeneratedCaseDraft
	var err error

	// Query based on filters
	if opts.TaskID != uuid.Nil {
		// If task_id is specified, get drafts for that task
		drafts, err = s.draftRepo.FindByTaskID(ctx, opts.TaskID)
		if err != nil {
			return nil, 0, fmt.Errorf("find drafts by task: %w", err)
		}
	} else if opts.ProjectID != uuid.Nil {
		// If project_id is specified, get all tasks for project, then their drafts
		tasks, taskErr := s.taskRepo.FindByProjectID(ctx, opts.ProjectID, generation.QueryOptions{})
		if taskErr != nil {
			return nil, 0, fmt.Errorf("find tasks for project: %w", taskErr)
		}

		// Collect all drafts from all tasks
		for _, task := range tasks {
			taskDrafts, draftErr := s.draftRepo.FindByTaskID(ctx, task.ID())
			if draftErr == nil {
				drafts = append(drafts, taskDrafts...)
			}
		}
	} else {
		// No filter specified, return empty or error
		return []*generation.GeneratedCaseDraft{}, 0, nil
	}

	// Filter by status if specified
	if opts.Status != "" {
		status, statusErr := generation.ParseDraftStatus(opts.Status)
		if statusErr == nil {
			filtered := make([]*generation.GeneratedCaseDraft, 0)
			for _, draft := range drafts {
				if draft.Status() == status {
					filtered = append(filtered, draft)
				}
			}
			drafts = filtered
		}
	}

	// Apply pagination
	total := int64(len(drafts))
	start := opts.Offset
	if start > len(drafts) {
		start = len(drafts)
	}
	end := start + opts.Limit
	if end > len(drafts) {
		end = len(drafts)
	}

	if start >= end {
		return []*generation.GeneratedCaseDraft{}, total, nil
	}

	return drafts[start:end], total, nil
}

// ConfirmDraft confirms a draft and creates a test case
func (s *GenerationServiceImpl) ConfirmDraft(ctx context.Context, req *ConfirmDraftRequest, userID uuid.UUID) (*testcase.TestCase, error) {
	// Get draft
	draft, err := s.draftRepo.FindByID(ctx, req.DraftID)
	if err != nil {
		return nil, fmt.Errorf("find draft: %w", err)
	}

	// Check if already confirmed
	if draft.Status() == generation.DraftConfirmed {
		return nil, generation.ErrDraftAlreadyConfirmed
	}

	// Check if already rejected
	if draft.Status() == generation.DraftRejected {
		return nil, generation.ErrDraftAlreadyRejected
	}

	// Get module
	module, err := s.moduleRepo.FindByID(ctx, req.ModuleID)
	if err != nil {
		return nil, errors.New("module not found")
	}

	// Confirm draft
	if err := draft.Confirm(req.ModuleID); err != nil {
		return nil, fmt.Errorf("confirm draft: %w", err)
	}

	// Generate case number
	caseNumber, err := s.generateCaseNumber(ctx, module)
	if err != nil {
		return nil, fmt.Errorf("generate case number: %w", err)
	}

	// Create test case from draft
	tc, err := testcase.NewTestCase(
		req.ModuleID,
		userID,
		caseNumber,
		draft.Title(),
		draft.Preconditions(),
		draft.Steps(),
		draft.ExpectedResult(),
		draft.CaseType(),
		draft.Priority(),
	)
	if err != nil {
		return nil, fmt.Errorf("create test case: %w", err)
	}

	// Set AI metadata if available
	if draft.AiMetadata() != nil {
		tc.SetAiMetadata(draft.AiMetadata())
	}

	// Update draft
	if err := s.draftRepo.Update(ctx, draft); err != nil {
		return nil, fmt.Errorf("update draft: %w", err)
	}

	// Save test case (if caseRepo is available)
	if s.caseRepo != nil {
		if err := s.caseRepo.Save(ctx, tc); err != nil {
			return nil, fmt.Errorf("save test case: %w", err)
		}
	}

	return tc, nil
}

// RejectDraft rejects a draft with reason
func (s *GenerationServiceImpl) RejectDraft(ctx context.Context, req *RejectDraftRequest) error {
	// Get draft
	draft, err := s.draftRepo.FindByID(ctx, req.DraftID)
	if err != nil {
		return fmt.Errorf("find draft: %w", err)
	}

	// Check if already confirmed
	if draft.Status() == generation.DraftConfirmed {
		return generation.ErrDraftAlreadyConfirmed
	}

	// Check if already rejected
	if draft.Status() == generation.DraftRejected {
		return generation.ErrDraftAlreadyRejected
	}

	// Validate rejection reason
	validReasons := map[generation.RejectionReason]bool{
		generation.ReasonDuplicate:  true,
		generation.ReasonIrrelevant: true,
		generation.ReasonLowQuality: true,
		generation.ReasonOther:      true,
	}
	if !validReasons[req.Reason] {
		return errors.New("invalid rejection reason")
	}

	// Reject draft
	if err := draft.Reject(req.Reason, req.Feedback); err != nil {
		return fmt.Errorf("reject draft: %w", err)
	}

	// Update draft
	if err := s.draftRepo.Update(ctx, draft); err != nil {
		return fmt.Errorf("update draft: %w", err)
	}

	return nil
}

// BatchConfirm confirms multiple drafts at once
func (s *GenerationServiceImpl) BatchConfirm(ctx context.Context, req *BatchConfirmRequest, userID uuid.UUID) (*BatchConfirmResult, error) {
	// Validate draft IDs
	if len(req.DraftIDs) == 0 {
		return nil, errors.New("draft IDs cannot be empty")
	}
	if len(req.DraftIDs) > 50 {
		return nil, errors.New("cannot confirm more than 50 drafts at once")
	}

	result := &BatchConfirmResult{
		SuccessIDs: make([]uuid.UUID, 0),
		FailedIDs:  make([]uuid.UUID, 0),
		Errors:     make([]string, 0),
	}

	for _, draftID := range req.DraftIDs {
		confirmReq := &ConfirmDraftRequest{
			DraftID:  draftID,
			ModuleID: req.ModuleID,
		}

		_, err := s.ConfirmDraft(ctx, confirmReq, userID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, draftID)
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.SuccessCount++
			result.SuccessIDs = append(result.SuccessIDs, draftID)
		}
	}

	return result, nil
}

// ProcessTask executes the generation workflow
func (s *GenerationServiceImpl) ProcessTask(ctx context.Context, taskID uuid.UUID) error {
	// Get task
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("find task: %w", err)
	}

	// Start processing
	if err := task.StartProcessing(); err != nil {
		return fmt.Errorf("start processing: %w", err)
	}
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("update task status: %w", err)
	}

	// Retrieve context via RAG
	ragReq := &RetrieveRequest{
		ProjectID: task.ProjectID(),
		Query:     task.Prompt(),
		TopK:      10,
	}

	ragResult, err := s.ragService.Retrieve(ctx, ragReq)
	if err != nil {
		task.Fail(fmt.Sprintf("RAG retrieval failed: %v", err))
		_ = s.taskRepo.Update(ctx, task)
		return fmt.Errorf("rag retrieval: %w", err)
	}

	// Calculate confidence
	confidence := s.ragService.CalculateConfidence(ragResult.Chunks)

	// Build context for LLM
	contextStr := s.buildContextFromChunks(ragResult.Chunks)

	// Generate cases via LLM
	llmReq := &GenerateCasesRequest{
		Prompt:    task.Prompt(),
		Context:   contextStr,
		CaseCount: 5,
		CaseType:  "functionality",
		Priority:  "P2",
	}

	llmResult, err := s.llmService.GenerateCases(ctx, llmReq)
	if err != nil {
		task.Fail(fmt.Sprintf("LLM generation failed: %v", err))
		_ = s.taskRepo.Update(ctx, task)
		return fmt.Errorf("llm generation: %w", err)
	}

	// Create drafts from generated cases
	for _, genCase := range llmResult.Cases {
		draft, err := generation.NewGeneratedCaseDraft(
			taskID,
			genCase.Title,
			testcase.Preconditions(genCase.Preconditions),
			testcase.Steps(genCase.Steps),
			testcase.ExpectedResult(genCase.Expected),
			testcase.CaseType(genCase.CaseType),
			testcase.Priority(genCase.Priority),
		)
		if err != nil {
			continue // Skip invalid drafts
		}

		// Set AI metadata
		chunks := s.convertChunksToReferenced(ragResult.Chunks)
		aiMetadata := testcase.NewAiMetadata(
			taskID,
			confidence,
			chunks,
			s.llmService.GetModelVersion(),
		)
		draft.SetAiMetadata(aiMetadata)

		// Save draft
		if err := s.draftRepo.Save(ctx, draft); err != nil {
			continue // Skip failed saves
		}
	}

	// Complete task
	task.Complete(map[string]any{
		"generated_count": len(llmResult.Cases),
		"model_version":   llmResult.ModelVersion,
		"confidence":      string(confidence),
	})

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("complete task: %w", err)
	}

	return nil
}

// generateCaseNumber generates a new case number for a module
func (s *GenerationServiceImpl) generateCaseNumber(ctx context.Context, module *project.Module) (testcase.CaseNumber, error) {
	// Get project for prefix
	if s.projectRepo == nil {
		return "", errors.New("project repository not available")
	}

	projectEntity, err := s.projectRepo.FindByID(ctx, module.ProjectID())
	if err != nil {
		return "", fmt.Errorf("find project: %w", err)
	}

	// Generate case number
	caseNumber := testcase.GenerateCaseNumber(
		projectEntity.Prefix().String(),
		module.Abbreviation().String(),
		1, // TODO: Get actual count from repository
	)

	return caseNumber, nil
}

// buildContextFromChunks builds context string from retrieved chunks
func (s *GenerationServiceImpl) buildContextFromChunks(chunks []*RetrievedChunk) string {
	if len(chunks) == 0 {
		return ""
	}

	context := "Relevant documentation:\n\n"
	for i, chunk := range chunks {
		context += fmt.Sprintf("[%d] %s\n\n", i+1, chunk.Content)
	}
	return context
}

// convertChunksToReferenced converts RetrievedChunk to ReferencedChunk
func (s *GenerationServiceImpl) convertChunksToReferenced(chunks []*RetrievedChunk) []*testcase.ReferencedChunk {
	result := make([]*testcase.ReferencedChunk, 0, len(chunks))
	for _, c := range chunks {
		result = append(result, testcase.NewReferencedChunk(
			c.ChunkID,
			c.DocumentID,
			c.DocumentName,
			c.SimilarityScore,
		))
	}
	return result
}
