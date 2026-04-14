// Package app provides application initialization
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/liang21/aitestos/internal/config"
	httptransport "github.com/liang21/aitestos/internal/transport/http"
	"github.com/liang21/aitestos/internal/transport/http/handler"
	httpmiddleware "github.com/liang21/aitestos/internal/transport/http/middleware"

	// Repository imports
	generationRepo "github.com/liang21/aitestos/internal/repository/generation"
	identityRepo "github.com/liang21/aitestos/internal/repository/identity"
	knowledgeRepo "github.com/liang21/aitestos/internal/repository/knowledge"
	projectRepo "github.com/liang21/aitestos/internal/repository/project"
	testcaseRepo "github.com/liang21/aitestos/internal/repository/testcase"
	testplanRepo "github.com/liang21/aitestos/internal/repository/testplan"

	// Service imports
	generationSvc "github.com/liang21/aitestos/internal/service/generation"
	identitySvc "github.com/liang21/aitestos/internal/service/identity"
	identitydomain "github.com/liang21/aitestos/internal/domain/identity"
	knowledgeSvc "github.com/liang21/aitestos/internal/service/knowledge"
	projectSvc "github.com/liang21/aitestos/internal/service/project"
	testcaseSvc "github.com/liang21/aitestos/internal/service/testcase"
	testplanSvc "github.com/liang21/aitestos/internal/service/testplan"

	// Domain imports
	domaingeneration "github.com/liang21/aitestos/internal/domain/generation"
	domainKnowledge "github.com/liang21/aitestos/internal/domain/knowledge"
	domainProject "github.com/liang21/aitestos/internal/domain/project"
	domaintestcase "github.com/liang21/aitestos/internal/domain/testcase"

	// Infrastructure imports
	"github.com/liang21/aitestos/internal/infrastructure/llm"
	"github.com/liang21/aitestos/internal/infrastructure/cache"
	redispkg "github.com/liang21/aitestos/internal/infrastructure/redis"
	"github.com/liang21/aitestos/internal/infrastructure/milvus"
	"github.com/liang21/aitestos/internal/infrastructure/rag"
	"github.com/liang21/aitestos/internal/infrastructure/vector"
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

	// 1.5. Connect to Redis (for refresh token storage)
	var redisClient identitydomain.TokenStore
	if cfg.Redis.Host != "" {
		redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
		client, err := redispkg.NewClient(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to connect to Redis, refresh tokens will be stored in-memory")
			// Fall back to mock store
			redisClient = redispkg.NewMockTokenStore()
		} else {
			redisClient = redispkg.NewTokenStore(client)
			log.Info().Msg("Redis connected successfully")
		}
	} else {
		log.Info().Msg("Redis not configured, using in-memory token store")
		redisClient = redispkg.NewMockTokenStore()
	}

	// 1.6. Initialize cache client
	var cacheClient cache.Cache
	logger := zerolog.New(os.Stdout).With().Str("service_name", "aitestos").Logger()
	if cfg.Redis.Host != "" {
		redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
		redisCacheClient, err := redispkg.NewClient(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to connect to Redis for caching, using in-memory cache")
			cacheClient = cache.NewMemoryCache()
		} else {
			cacheClient = cache.NewRedisCache(redisCacheClient, &logger)
			log.Info().Msg("Redis cache initialized successfully")
		}
	} else {
		log.Info().Msg("Redis not configured, using in-memory cache")
		cacheClient = cache.NewMemoryCache()
	}

	// 2. Initialize Repositories
	userRepo := identityRepo.NewUserRepository(db)
	baseProjectRepository := projectRepo.NewProjectRepository(db)
	projectRepository := projectRepo.NewCachedProjectRepository(baseProjectRepository, cacheClient)
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
	projectRepoAdapter := &projectRepoAdapterImpl{repo: baseProjectRepository}

	// For generation service, use direct repositories (they match the expected interface)
	genModuleRepoAdapter := &genModuleRepoAdapterImpl{repo: moduleRepo}
	genProjectRepoAdapter := &genProjectRepoAdapterImpl{repo: baseProjectRepository}

	// 4. Initialize external infrastructure services
	// NOTE: These are placeholder implementations that return "not yet implemented" errors.
	// TODO: Replace with real implementations when LLM/Milvus are available.

	// Initialize Milvus client for vector operations
	var milvusClient *milvus.Client
	var vectorRepo *vector.Repository
	if cfg.Milvus.Host != "" {
		milvusClient, err = milvus.NewClient(&cfg.Milvus)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to initialize Milvus client, vector search will be unavailable")
		} else {
			// Ensure collection exists
			if err := milvusClient.EnsureCollection(context.Background()); err != nil {
				log.Warn().Err(err).Msg("Failed to ensure Milvus collection, vector search may be unavailable")
			}
			vectorRepo = vector.NewRepository(milvusClient)
		}
	}

	// Initialize LLM client
	var llmClient *llm.Client
	if cfg.LLM.APIKey != "" {
		llmClient, err = llm.NewClient(&cfg.LLM)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to initialize LLM client, AI features will be unavailable")
		}
	}

	// Initialize RAG service (requires Vector repository)
	var ragSvc *rag.Service
	if vectorRepo != nil {
		ragSvc, err = rag.NewService(&cfg.Milvus, vectorRepo)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to initialize RAG service, semantic search will be unavailable")
		} else {
			// Wire up dependencies for RAG service
			ragSvc = ragSvc.WithDocumentRepo(documentRepo).
				WithChunkRepo(chunkRepo)
			// LLM client will be set after generation service is created
		}
	}

	// Create vector repo adapter for knowledge service
	var vectorRepoAdapter domainKnowledge.VectorRepository
	if vectorRepo != nil {
		vectorRepoAdapter = &vectorRepoAdapterImpl{repo: vectorRepo}
	} else {
		// Fallback to nil adapter - document operations will work but vector search will fail
		vectorRepoAdapter = nil
		log.Warn().Msg("Vector repository not available, knowledge features limited to storage only")
	}

	// 5. Initialize Services
	authService := identitySvc.NewAuthService(userRepo, cfg.JWT.Secret, redisClient)
	projectService := projectSvc.NewProjectService(projectRepository, moduleRepo, configRepo)
	caseService := testcaseSvc.NewCaseService(caseRepo, moduleRepoAdapter, projectRepoAdapter)
	planService := testplanSvc.NewPlanService(planRepo, resultRepo, caseRepo)
	documentService := knowledgeSvc.NewDocumentService(documentRepo, chunkRepo, vectorRepoAdapter)

	// Generation service requires LLM and RAG
	var generationService generationSvc.GenerationService
	if llmClient != nil && ragSvc != nil {
		generationService = generationSvc.NewGenerationService(
			taskRepo,
			draftRepo,
			ragSvc,
			llmClient,
			genModuleRepoAdapter,
			genProjectRepoAdapter,
			caseRepo,
		)
	} else {
		log.Warn().Msg("LLM or RAG service unavailable, generation features disabled")
		// Create a placeholder service that returns proper errors
		generationService = &placeholderGenerationService{}
	}

	// 6. Initialize Handlers
	identityHandler := handler.NewIdentityHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	caseHandler := handler.NewTestCaseHandler(caseService)
	planHandler := handler.NewTestPlanHandler(planService)
	generationHandler := handler.NewGenerationHandler(generationService)
	knowledgeHandler := handler.NewKnowledgeHandler(documentService)

	// 7. Create HTTP handlers struct
	handlers := &httptransport.Handlers{
		Identity:   identityHandler,
		Project:    projectHandler,
		TestCase:   caseHandler,
		TestPlan:   planHandler,
		Generation: generationHandler,
		Knowledge:  knowledgeHandler,
	}

	// 8. Create observability middleware
	logger = zerolog.New(os.Stdout).With().
		Str("service_name", "aitestos").
		Logger()
	logger = logger.With().Timestamp().Logger()

	metrics := httpmiddleware.NewMetrics("aitestos")
	prometheus.MustRegister(metrics.RequestsTotal, metrics.RequestDuration)

	// 9. Create router with authentication middleware
	router := httptransport.NewRouterWithMiddleware(handlers, cfg.JWT.Secret, logger, metrics)

	// 10. Create HTTP server
	httpServer := NewHTTPServerFromConfig(cfg, router)

	// 11. Create shutdown manager and register closers
	shutdownMgr := NewShutdownManagerWithTimeout(cfg.Server.ShutdownTimeout)
	shutdownMgr.Register(dbCloser)
	if milvusClient != nil {
		shutdownMgr.Register(&milvusCloser{client: milvusClient})
	}

	// 12. Create application
	app := New(cfg, httpServer, shutdownMgr)

	// 13. Return cleanup function
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

// milvusCloser wraps milvus.Client to implement Closer interface
type milvusCloser struct {
	client *milvus.Client
}

func (m *milvusCloser) Name() string {
	return "milvus"
}

func (m *milvusCloser) Close(ctx context.Context) error {
	return m.client.Close()
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

func (w *moduleWrapper) ID() uuid.UUID        { return w.Module.ID() }
func (w *moduleWrapper) ProjectID() uuid.UUID { return w.Module.ProjectID() }
func (w *moduleWrapper) Name() string         { return w.Module.Name() }
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

func (w *projectWrapper) ID() uuid.UUID  { return w.Project.ID() }
func (w *projectWrapper) Name() string   { return w.Project.Name() }
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

// Vector repository adapter for knowledge service
type vectorRepoAdapterImpl struct {
	repo *vector.Repository
}

func (a *vectorRepoAdapterImpl) Upsert(ctx context.Context, chunks []*domainKnowledge.DocumentChunk) error {
	return a.repo.Upsert(ctx, chunks)
}

func (a *vectorRepoAdapterImpl) Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*domainKnowledge.DocumentChunk, error) {
	return a.repo.Search(ctx, queryVector, topK, filter)
}

func (a *vectorRepoAdapterImpl) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	return a.repo.DeleteByDocumentID(ctx, documentID)
}

func (a *vectorRepoAdapterImpl) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return a.repo.CountByProjectID(ctx, projectID)
}

// placeholderGenerationService is used when LLM/RAG services are unavailable
type placeholderGenerationService struct{}

func (s *placeholderGenerationService) CreateTask(ctx context.Context, req *generationSvc.CreateTaskRequest, userID uuid.UUID) (*domaingeneration.GenerationTask, error) {
	return nil, fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}

func (s *placeholderGenerationService) GetTask(ctx context.Context, id uuid.UUID) (*domaingeneration.GenerationTask, error) {
	return nil, fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}

func (s *placeholderGenerationService) ListTasks(ctx context.Context, projectID uuid.UUID, opts generationSvc.ListTaskOptions) ([]*domaingeneration.GenerationTask, int64, error) {
	return nil, 0, fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}

func (s *placeholderGenerationService) GetDrafts(ctx context.Context, taskID uuid.UUID) ([]*domaingeneration.GeneratedCaseDraft, error) {
	return nil, fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}

func (s *placeholderGenerationService) ConfirmDraft(ctx context.Context, req *generationSvc.ConfirmDraftRequest, userID uuid.UUID) (*domaintestcase.TestCase, error) {
	return nil, fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}

func (s *placeholderGenerationService) RejectDraft(ctx context.Context, req *generationSvc.RejectDraftRequest) error {
	return fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}

func (s *placeholderGenerationService) BatchConfirm(ctx context.Context, req *generationSvc.BatchConfirmRequest, userID uuid.UUID) (*generationSvc.BatchConfirmResult, error) {
	return nil, fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}

func (s *placeholderGenerationService) ProcessTask(ctx context.Context, taskID uuid.UUID) error {
	return fmt.Errorf("generation service unavailable: LLM and RAG services not configured")
}
