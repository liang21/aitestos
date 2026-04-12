// Package testutil provides shared utilities for integration tests
package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/liang21/aitestos/internal/domain/identity"
	domainKnowledge "github.com/liang21/aitestos/internal/domain/knowledge"
	domainProject "github.com/liang21/aitestos/internal/domain/project"
	domainTestcase "github.com/liang21/aitestos/internal/domain/testcase"
	generationRepo "github.com/liang21/aitestos/internal/repository/generation"
	identityRepo "github.com/liang21/aitestos/internal/repository/identity"
	knowledgeRepo "github.com/liang21/aitestos/internal/repository/knowledge"
	projectRepo "github.com/liang21/aitestos/internal/repository/project"
	testcaseRepo "github.com/liang21/aitestos/internal/repository/testcase"
	testplanRepo "github.com/liang21/aitestos/internal/repository/testplan"
	generationService "github.com/liang21/aitestos/internal/service/generation"
	knowledgeService "github.com/liang21/aitestos/internal/service/knowledge"
	projectSvc "github.com/liang21/aitestos/internal/service/project"
	testcaseSvc "github.com/liang21/aitestos/internal/service/testcase"
	testplanSvc "github.com/liang21/aitestos/internal/service/testplan"
	"github.com/liang21/aitestos/internal/transport/http/handler"
)

// Context key for user ID
type contextKey string

const userIDContextKey contextKey = "user_id"

// IntegrationTestSuite encapsulates integration test infrastructure
type IntegrationTestSuite struct {
	DB     *sqlx.DB
	Router chi.Router
}

// NewIntegrationTestSuite creates a test suite with real handlers wired through
func NewIntegrationTestSuite(db *sqlx.DB) *IntegrationTestSuite {
	suite := &IntegrationTestSuite{DB: db}

	// Initialize mocks for external dependencies
	ragSvc := NewMockRAGService()
	llmSvc := NewMockLLMService()
	vectorRepo := NewMockVectorRepository()

	// Initialize repositories
	_ = identityRepo.NewUserRepository(db) // available for future use
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

	// Initialize services
	// testcase service defines its own ModuleRepository/ProjectRepository interfaces
	// We use adapter wrappers to bridge domain repos to service interfaces
	projectSvc_ := projectSvc.NewProjectService(projectRepository, moduleRepo, configRepo)
	caseSvc := testcaseSvc.NewCaseService(caseRepo, &moduleRepoAdapter{repo: moduleRepo}, &projectRepoAdapter{repo: projectRepository})
	planSvc := testplanSvc.NewPlanService(planRepo, resultRepo, caseRepo)
	documentSvc := knowledgeService.NewDocumentService(documentRepo, chunkRepo, vectorRepo)
	generationSvc := generationService.NewGenerationService(taskRepo, draftRepo, ragSvc, llmSvc, moduleRepo, projectRepository, caseRepo)

	// Initialize handlers
	projectHandler := handler.NewProjectHandler(projectSvc_)
	testCaseHandler := handler.NewTestCaseHandler(caseSvc)
	testPlanHandler := handler.NewTestPlanHandler(planSvc)
	generationHandler := handler.NewGenerationHandler(generationSvc)
	knowledgeHandler := handler.NewKnowledgeHandler(documentSvc)

	// Build router
	r := chi.NewRouter()

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Project routes
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", projectHandler.ListProjects)
			r.Post("/", projectHandler.CreateProject)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", projectHandler.GetProject)
				r.Put("/", projectHandler.UpdateProject)
				r.Delete("/", projectHandler.DeleteProject)
				r.Get("/modules", projectHandler.ListModules)
				r.Post("/modules", projectHandler.CreateModule)
				r.Get("/configs", projectHandler.ListConfigs)
				r.Put("/configs/{key}", projectHandler.SetConfig)
			})
		})

		// Module routes
		r.Route("/modules", func(r chi.Router) {
			r.Delete("/{id}", projectHandler.DeleteModule)
		})

		// Test case routes
		r.Route("/testcases", func(r chi.Router) {
			r.Get("/", testCaseHandler.ListCases)
			r.Post("/", testCaseHandler.CreateCase)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", testCaseHandler.GetCase)
				r.Put("/", testCaseHandler.UpdateCase)
				r.Delete("/", testCaseHandler.DeleteCase)
			})
		})

		// Test plan routes
		r.Route("/plans", func(r chi.Router) {
			r.Get("/", testPlanHandler.ListPlans)
			r.Post("/", testPlanHandler.CreatePlan)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", testPlanHandler.GetPlan)
				r.Put("/", testPlanHandler.UpdatePlan)
				r.Delete("/", testPlanHandler.DeletePlan)
				r.Post("/cases", testPlanHandler.AddCases)
				r.Delete("/cases/{caseId}", testPlanHandler.RemoveCase)
				r.Post("/results", testPlanHandler.RecordResult)
				r.Get("/results", testPlanHandler.GetResults)
			})
		})

		// Generation routes
		r.Route("/generation", func(r chi.Router) {
			r.Route("/tasks", func(r chi.Router) {
				r.Post("/", generationHandler.CreateTask)
				r.Get("/{id}", generationHandler.GetTask)
				r.Get("/{id}/drafts", generationHandler.GetDrafts)
			})
			r.Route("/drafts", func(r chi.Router) {
				r.Post("/{id}/confirm", generationHandler.ConfirmDraft)
				r.Post("/{id}/reject", generationHandler.RejectDraft)
				r.Post("/batch-confirm", generationHandler.BatchConfirm)
			})
		})

		// Knowledge routes
		r.Route("/knowledge", func(r chi.Router) {
			r.Route("/documents", func(r chi.Router) {
				r.Get("/", knowledgeHandler.ListDocuments)
				r.Post("/", knowledgeHandler.UploadDocument)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", knowledgeHandler.GetDocument)
					r.Delete("/", knowledgeHandler.DeleteDocument)
					r.Get("/chunks", knowledgeHandler.GetChunks)
				})
			})
		})
	})

	suite.Router = r
	return suite
}

// MakeRequest executes an HTTP request against the test router
func (s *IntegrationTestSuite) MakeRequest(method, path string, body interface{}, userID uuid.UUID) *httptest.ResponseRecorder {
	var reqBody bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&reqBody).Encode(body)
	}

	req := httptest.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Inject user ID into context for authentication
	ctx := context.WithValue(req.Context(), userIDContextKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)

	return w
}

// MakeRequestWithoutAuth executes an HTTP request without user context
func (s *IntegrationTestSuite) MakeRequestWithoutAuth(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&reqBody).Encode(body)
	}

	req := httptest.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)

	return w
}

// Service interface adapters
// testcase service defines its own ModuleRepository/ProjectRepository interfaces
// that return interface types instead of concrete domain pointers.
// These adapters bridge the gap.

type moduleRepoAdapter struct {
	repo *projectRepo.ModuleRepository
}

func (a *moduleRepoAdapter) FindByID(ctx context.Context, id uuid.UUID) (testcaseSvc.Module, error) {
	m, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return testcaseSvc.ModuleWrapper{Module: m}, nil
}

type projectRepoAdapter struct {
	repo *projectRepo.ProjectRepository
}

func (a *projectRepoAdapter) FindByID(ctx context.Context, id uuid.UUID) (testcaseSvc.Project, error) {
	p, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return testcaseSvc.ProjectWrapper{Project: p}, nil
}

// Test Data Builders

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *sqlx.DB) *identity.User {
	t.Helper()

	ctx := context.Background()
	username := "testuser_" + uuid.New().String()[:8]
	email := "test_" + uuid.New().String()[:8] + "@example.com"
	password := "Password123!"
	role := identity.RoleNormal

	user, err := identity.NewUser(username, email, password, role)
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}

	query := `
		INSERT INTO users (id, username, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = db.ExecContext(ctx, query,
		user.ID(),
		user.Username(),
		user.Email(),
		password,
		user.Role(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		t.Fatalf("save test user: %v", err)
	}

	return user
}

// CreateTestProject creates a test project in the database
func CreateTestProject(t *testing.T, db *sqlx.DB) *domainProject.Project {
	t.Helper()

	ctx := context.Background()
	prefix := "TP" + uuid.New().String()[:2]
	name := "Test Project " + uuid.New().String()[:8]

	project, err := domainProject.NewProject(name, prefix, "Test project description")
	if err != nil {
		t.Fatalf("create test project: %v", err)
	}

	query := `
		INSERT INTO project (id, name, prefix, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = db.ExecContext(ctx, query,
		project.ID(),
		project.Name(),
		project.Prefix().String(),
		project.Description(),
		project.CreatedAt(),
		project.UpdatedAt(),
	)
	if err != nil {
		t.Fatalf("save test project: %v", err)
	}

	return project
}

// CreateTestModule creates a test module in the database
func CreateTestModule(t *testing.T, db *sqlx.DB, projectID uuid.UUID, userID uuid.UUID) *domainProject.Module {
	t.Helper()

	ctx := context.Background()
	abbrev := "TM" + uuid.New().String()[:2]
	name := "Test Module " + uuid.New().String()[:8]

	module, err := domainProject.NewModule(projectID, name, abbrev, "Test module description", userID)
	if err != nil {
		t.Fatalf("create test module: %v", err)
	}

	query := `
		INSERT INTO module (id, project_id, name, abbreviation, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = db.ExecContext(ctx, query,
		module.ID(),
		module.ProjectID(),
		module.Name(),
		module.Abbreviation().String(),
		module.Description(),
		module.CreatedAt(),
		module.UpdatedAt(),
	)
	if err != nil {
		t.Fatalf("save test module: %v", err)
	}

	return module
}

// CreateTestCase creates a test case in the database
func CreateTestCase(t *testing.T, db *sqlx.DB, moduleID uuid.UUID, userID uuid.UUID) *domainTestcase.TestCase {
	t.Helper()

	ctx := context.Background()
	number := domainTestcase.GenerateCaseNumber("TP", "TM", 1)
	title := "Test Case " + uuid.New().String()[:8]

	tc, err := domainTestcase.NewTestCase(
		moduleID,
		userID,
		number,
		title,
		domainTestcase.Preconditions{"Precondition 1"},
		domainTestcase.Steps{"Step 1", "Step 2"},
		domainTestcase.ExpectedResult{"status": "success"},
		domainTestcase.CaseTypeFunctionality,
		domainTestcase.PriorityP2,
	)
	if err != nil {
		t.Fatalf("create test case: %v", err)
	}

	query := `
		INSERT INTO test_case (id, module_id, user_id, number, title, preconditions, steps, expected, case_type, priority, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err = db.ExecContext(ctx, query,
		tc.ID(),
		tc.ModuleID(),
		tc.UserID(),
		tc.Number().String(),
		tc.Title(),
		tc.Preconditions(),
		tc.Steps(),
		tc.ExpectedResult(),
		tc.CaseType(),
		tc.Priority(),
		tc.Status(),
		tc.CreatedAt(),
		tc.UpdatedAt(),
	)
	if err != nil {
		t.Fatalf("save test case: %v", err)
	}

	return tc
}

// WaitForCondition waits for a condition to be true with timeout
func WaitForCondition(t *testing.T, timeout time.Duration, condition func() bool, message string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for condition: %s", message)
}

// ParseJSONResponse parses JSON response body into target
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), target); err != nil {
		t.Fatalf("parse JSON response: %v", err)
	}
}

// Mock implementations for external dependencies

// MockRAGService implements generationService.RAGService for testing
type MockRAGService struct{}

// NewMockRAGService creates a new mock RAG service
func NewMockRAGService() *MockRAGService {
	return &MockRAGService{}
}

// Retrieve returns mock retrieved chunks
func (m *MockRAGService) Retrieve(ctx context.Context, req *generationService.RetrieveRequest) (*generationService.RetrieveResult, error) {
	return &generationService.RetrieveResult{
		Chunks: []*generationService.RetrievedChunk{
			{
				ChunkID:         uuid.New(),
				DocumentID:      uuid.New(),
				DocumentName:    "Mock Document",
				Content:         "Mock content for testing",
				SimilarityScore: 0.85,
			},
		},
		Query: req.Query,
	}, nil
}

// CalculateConfidence calculates confidence based on chunks
func (m *MockRAGService) CalculateConfidence(chunks []*generationService.RetrievedChunk) domainTestcase.Confidence {
	if len(chunks) >= 2 && chunks[0].SimilarityScore > 0.8 {
		return domainTestcase.ConfidenceHigh
	}
	if len(chunks) >= 1 && chunks[0].SimilarityScore >= 0.5 {
		return domainTestcase.ConfidenceMedium
	}
	return domainTestcase.ConfidenceLow
}

// MockLLMService implements generationService.LLMService for testing
type MockLLMService struct{}

// NewMockLLMService creates a new mock LLM service
func NewMockLLMService() *MockLLMService {
	return &MockLLMService{}
}

// GenerateCases returns mock generated test cases
func (m *MockLLMService) GenerateCases(ctx context.Context, req *generationService.GenerateCasesRequest) (*generationService.GenerateCasesResponse, error) {
	cases := make([]*generationService.GeneratedCase, 0, req.CaseCount)
	for i := 0; i < req.CaseCount && i < 3; i++ {
		cases = append(cases, &generationService.GeneratedCase{
			Title:         "Mock Test Case " + string(rune('A'+i)),
			Preconditions: []string{"Precondition 1"},
			Steps:         []string{"Step 1", "Step 2"},
			Expected:      map[string]interface{}{"status": "success"},
			CaseType:      req.CaseType,
			Priority:      req.Priority,
			Reasoning:     "Mock reasoning for testing",
		})
	}
	return &generationService.GenerateCasesResponse{
		Cases:        cases,
		ModelVersion: "mock-v1.0",
		TokensUsed:   1000,
	}, nil
}

// GenerateEmbedding returns mock embedding
func (m *MockLLMService) GenerateEmbedding(ctx context.Context, req *generationService.GenerateEmbeddingRequest) (*generationService.GenerateEmbeddingResponse, error) {
	embedding := make([]byte, 1536*4)
	return &generationService.GenerateEmbeddingResponse{
		Embedding: embedding,
		Model:     "mock-embedding-v1",
	}, nil
}

// GetModelVersion returns mock model version
func (m *MockLLMService) GetModelVersion() string {
	return "mock-v1.0"
}

// MockVectorRepository implements domainKnowledge.VectorRepository for testing
type MockVectorRepository struct{}

// NewMockVectorRepository creates a new mock vector repository
func NewMockVectorRepository() *MockVectorRepository {
	return &MockVectorRepository{}
}

// Upsert inserts or updates vectors (mock implementation)
func (m *MockVectorRepository) Upsert(ctx context.Context, chunks []*domainKnowledge.DocumentChunk) error {
	return nil
}

// Search performs vector similarity search (mock implementation)
func (m *MockVectorRepository) Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*domainKnowledge.DocumentChunk, error) {
	// Return mock chunk with projectID
	chunk := domainKnowledge.ReconstructDocumentChunk(
		uuid.New(),
		uuid.New(),
		uuid.New(), // projectID
		0,
		"mock content",
		nil,
		time.Time{},
	)
	return []*domainKnowledge.DocumentChunk{chunk}, nil
}

// DeleteByDocumentID deletes vectors by document ID (mock implementation)
func (m *MockVectorRepository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	return nil
}

// CountByProjectID counts documents by project (mock implementation)
func (m *MockVectorRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return 0, nil
}
