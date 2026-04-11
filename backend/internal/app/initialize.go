// Package app provides application initialization
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/liang21/aitestos/internal/config"
	httptransport "github.com/liang21/aitestos/internal/transport/http"
	"github.com/liang21/aitestos/internal/transport/http/handler"

	// Repository imports
	identityRepo "github.com/liang21/aitestos/internal/repository/identity"
	projectRepo "github.com/liang21/aitestos/internal/repository/project"
	testcaseRepo "github.com/liang21/aitestos/internal/repository/testcase"
	testplanRepo "github.com/liang21/aitestos/internal/repository/testplan"
	knowledgeRepo "github.com/liang21/aitestos/internal/repository/knowledge"
	generationRepo "github.com/liang21/aitestos/internal/repository/generation"

	// Service imports
	identitySvc "github.com/liang21/aitestos/internal/service/identity"
	projectSvc "github.com/liang21/aitestos/internal/service/project"
	testcaseSvc "github.com/liang21/aitestos/internal/service/testcase"
	testplanSvc "github.com/liang21/aitestos/internal/service/testplan"
	knowledgeSvc "github.com/liang21/aitestos/internal/service/knowledge"
	generationSvc "github.com/liang21/aitestos/internal/service/generation"

	// Domain imports
	domainKnowledge "github.com/liang21/aitestos/internal/domain/knowledge"
	domainTestcase "github.com/liang21/aitestos/internal/domain/testcase"
	domainProject "github.com/liang21/aitestos/internal/domain/project"
)

// Initialize creates a fully initialized application
func Initialize(cfg *config.Config) (*App, func(), error) {
	// 1. Connect to database
	db, err := NewDB(&cfg.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("initialize database: %w", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}
	log.Info().Msg("Database connected successfully")

	// 2. Initialize Repositories
	userRepo := identityRepo.NewUserRepository(db)
	projectRepository := projectRepo.NewProjectRepository(db)
	moduleRepo := projectRepo.NewModuleRepository(db)
	configRepo := projectRepo.NewProjectConfigRepository(db)
	caseRepo := testcaseRepo.NewTestCaseRepository(db)
	planRepo := testplanRepo.NewTestPlanRepository(db)
	resultRepo := testplanRepo.NewTestResultRepository(db)
	documentRepo := knowledgeRepo.NewDocumentRepository(db)
	chunkRepo := knowledgeRepo.NewDocumentChunkRepository(db)
	taskRepo := generationRepo.NewGenerationTaskRepository(db)
	draftRepo := generationRepo.NewCaseDraftRepository(db)

	// Register database for cleanup
	dbCloser := NewDBCloser(db, "database")

	// 3. Create repository adapters for services
	moduleRepoAdapter := &moduleRepoAdapterImpl{repo: moduleRepo}
	projectRepoAdapter := &projectRepoAdapterImpl{repo: projectRepository}

	// For generation service, use direct repositories (they match the expected interface)
	genModuleRepoAdapter := &genModuleRepoAdapterImpl{repo: moduleRepo}
	genProjectRepoAdapter := &genProjectRepoAdapterImpl{repo: projectRepository}

	// 4. Initialize Services (using mocks for external dependencies)
	mockRAGSvc := &mockRAGService{}
	mockLLMSvc := &mockLLMService{}
	mockVectorRepo := &mockVectorRepository{}

	authService := identitySvc.NewAuthService(userRepo, cfg.JWT.Secret)
	projectService := projectSvc.NewProjectService(projectRepository, moduleRepo, configRepo)
	caseService := testcaseSvc.NewCaseService(caseRepo, moduleRepoAdapter, projectRepoAdapter)
	planService := testplanSvc.NewPlanService(planRepo, resultRepo, caseRepo)
	documentService := knowledgeSvc.NewDocumentService(documentRepo, chunkRepo, mockVectorRepo)
	generationService := generationSvc.NewGenerationService(taskRepo, draftRepo, mockRAGSvc, mockLLMSvc, genModuleRepoAdapter, genProjectRepoAdapter, caseRepo)

	// 5. Initialize Handlers
	identityHandler := handler.NewIdentityHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	caseHandler := handler.NewTestCaseHandler(caseService)
	planHandler := handler.NewTestPlanHandler(planService)
	generationHandler := handler.NewGenerationHandler(generationService)
	knowledgeHandler := handler.NewKnowledgeHandler(documentService)

	// 6. Create HTTP handlers struct
	handlers := &httptransport.Handlers{
		Identity:   identityHandler,
		Project:    projectHandler,
		TestCase:   caseHandler,
		TestPlan:   planHandler,
		Generation: generationHandler,
		Knowledge:  knowledgeHandler,
	}

	// 7. Create router with authentication middleware
	router := httptransport.NewRouterWithMiddleware(handlers, cfg.JWT.Secret)

	// 8. Create HTTP server
	httpServer := NewHTTPServerFromConfig(cfg, router)

	// 9. Create shutdown manager and register closers
	shutdownMgr := NewShutdownManagerWithTimeout(cfg.Server.ShutdownTimeout)
	shutdownMgr.Register(dbCloser)

	// 10. Create application
	app := New(cfg, httpServer, shutdownMgr)

	// 11. Return cleanup function
	cleanup := func() {
		log.Info().Msg("Running cleanup...")
	}

	log.Info().Msg("Application initialized successfully")
	return app, cleanup, nil
}

// NewHTTPServerFromConfig creates an HTTP server from config
func NewHTTPServerFromConfig(cfg *config.Config, handler http.Handler) HTTPServer {
	return &httpServerWrapper{
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// httpServerWrapper wraps http.Server to implement HTTPServer interface
type httpServerWrapper struct {
	*http.Server
}

// Addr returns the server address
func (s *httpServerWrapper) Addr() string {
	return s.Server.Addr
}

// NewShutdownManagerWithTimeout creates a shutdown manager with specific timeout
func NewShutdownManagerWithTimeout(timeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		closers: make([]Closer, 0),
		timeout: timeout,
	}
}

// Repository adapters for testcase service

type moduleRepoAdapterImpl struct {
	repo *projectRepo.ModuleRepository
}

func (a *moduleRepoAdapterImpl) FindByID(ctx context.Context, id uuid.UUID) (testcaseSvc.Module, error) {
	m, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &moduleWrapper{Module: m}, nil
}

type moduleWrapper struct {
	*domainProject.Module
}

func (w *moduleWrapper) ID() uuid.UUID { return w.Module.ID() }
func (w *moduleWrapper) ProjectID() uuid.UUID { return w.Module.ProjectID() }
func (w *moduleWrapper) Name() string { return w.Module.Name() }
func (w *moduleWrapper) Abbreviation() string { return w.Module.Abbreviation().String() }

type projectRepoAdapterImpl struct {
	repo *projectRepo.ProjectRepository
}

func (a *projectRepoAdapterImpl) FindByID(ctx context.Context, id uuid.UUID) (testcaseSvc.Project, error) {
	p, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &projectWrapper{Project: p}, nil
}

type projectWrapper struct {
	*domainProject.Project
}

func (w *projectWrapper) ID() uuid.UUID { return w.Project.ID() }
func (w *projectWrapper) Name() string { return w.Project.Name() }
func (w *projectWrapper) Prefix() string { return w.Project.Prefix().String() }

// Repository adapters for generation service

type genModuleRepoAdapterImpl struct {
	repo *projectRepo.ModuleRepository
}

func (a *genModuleRepoAdapterImpl) FindByID(ctx context.Context, id uuid.UUID) (*domainProject.Module, error) {
	return a.repo.FindByID(ctx, id)
}

func (a *genModuleRepoAdapterImpl) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domainProject.Module, error) {
	return a.repo.FindByProjectID(ctx, projectID)
}

type genProjectRepoAdapterImpl struct {
	repo *projectRepo.ProjectRepository
}

func (a *genProjectRepoAdapterImpl) FindByID(ctx context.Context, id uuid.UUID) (*domainProject.Project, error) {
	return a.repo.FindByID(ctx, id)
}

// Mock implementations for external dependencies

type mockRAGService struct{}

func (m *mockRAGService) Retrieve(ctx context.Context, req *generationSvc.RetrieveRequest) (*generationSvc.RetrieveResult, error) {
	return &generationSvc.RetrieveResult{
		Chunks: []*generationSvc.RetrievedChunk{},
		Query:  req.Query,
	}, nil
}

func (m *mockRAGService) CalculateConfidence(chunks []*generationSvc.RetrievedChunk) domainTestcase.Confidence {
	return domainTestcase.ConfidenceMedium
}

type mockLLMService struct{}

func (m *mockLLMService) GenerateCases(ctx context.Context, req *generationSvc.GenerateCasesRequest) (*generationSvc.GenerateCasesResponse, error) {
	return &generationSvc.GenerateCasesResponse{
		Cases:        []*generationSvc.GeneratedCase{},
		ModelVersion: "mock-v1.0",
		TokensUsed:   1000,
	}, nil
}

func (m *mockLLMService) GenerateEmbedding(ctx context.Context, req *generationSvc.GenerateEmbeddingRequest) (*generationSvc.GenerateEmbeddingResponse, error) {
	embedding := make([]byte, 1536*4)
	return &generationSvc.GenerateEmbeddingResponse{
		Embedding: embedding,
		Model:     "mock-embedding-v1",
	}, nil
}

func (m *mockLLMService) GetModelVersion() string {
	return "mock-v1.0"
}

type mockVectorRepository struct{}

func (m *mockVectorRepository) Upsert(ctx context.Context, chunks []*domainKnowledge.DocumentChunk) error {
	return nil
}

func (m *mockVectorRepository) Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*domainKnowledge.DocumentChunk, error) {
	return []*domainKnowledge.DocumentChunk{}, nil
}

func (m *mockVectorRepository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	return nil
}

func (m *mockVectorRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return 0, nil
}
