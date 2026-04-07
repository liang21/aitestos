// Package generation provides AI generation services
package generation

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

// MockVectorRepository for RAG testing
type MockVectorRepo struct {
	results   []*knowledge.DocumentChunk
	searchErr error
}

func NewMockVectorRepo() *MockVectorRepo {
	return &MockVectorRepo{}
}

func (m *MockVectorRepo) Upsert(ctx context.Context, chunks []*knowledge.DocumentChunk) error {
	return nil
}

func (m *MockVectorRepo) Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*knowledge.DocumentChunk, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	return m.results, nil
}

func (m *MockVectorRepo) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	return nil
}

func (m *MockVectorRepo) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return int64(len(m.results)), nil
}

// MockChunkRepo for RAG testing
type MockChunkRepo struct {
	chunks map[uuid.UUID]*knowledge.DocumentChunk
}

func NewMockChunkRepo() *MockChunkRepo {
	return &MockChunkRepo{
		chunks: make(map[uuid.UUID]*knowledge.DocumentChunk),
	}
}

func (m *MockChunkRepo) SaveBatch(ctx context.Context, chunks []*knowledge.DocumentChunk) error {
	for _, c := range chunks {
		m.chunks[c.ID()] = c
	}
	return nil
}

func (m *MockChunkRepo) FindByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*knowledge.DocumentChunk, error) {
	result := make([]*knowledge.DocumentChunk, 0)
	for _, c := range m.chunks {
		if c.DocumentID() == documentID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *MockChunkRepo) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	for id, c := range m.chunks {
		if c.DocumentID() == documentID {
			delete(m.chunks, id)
		}
	}
	return nil
}

func (m *MockChunkRepo) CountByDocumentID(ctx context.Context, documentID uuid.UUID) (int64, error) {
	count := 0
	for _, c := range m.chunks {
		if c.DocumentID() == documentID {
			count++
		}
	}
	return int64(count), nil
}

// MockDocRepo for RAG testing
type MockDocRepo struct {
	docs map[uuid.UUID]*knowledge.Document
}

func NewMockDocRepo() *MockDocRepo {
	return &MockDocRepo{
		docs: make(map[uuid.UUID]*knowledge.Document),
	}
}

func (m *MockDocRepo) Save(ctx context.Context, doc *knowledge.Document) error {
	m.docs[doc.ID()] = doc
	return nil
}

func (m *MockDocRepo) FindByID(ctx context.Context, id uuid.UUID) (*knowledge.Document, error) {
	doc, ok := m.docs[id]
	if !ok {
		return nil, knowledge.ErrDocumentNotFound
	}
	return doc, nil
}

func (m *MockDocRepo) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts knowledge.QueryOptions) ([]*knowledge.Document, error) {
	result := make([]*knowledge.Document, 0)
	for _, d := range m.docs {
		if d.ProjectID() == projectID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (m *MockDocRepo) Update(ctx context.Context, doc *knowledge.Document) error {
	m.docs[doc.ID()] = doc
	return nil
}

func (m *MockDocRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status knowledge.DocumentStatus) error {
	_, ok := m.docs[id]
	if !ok {
		return knowledge.ErrDocumentNotFound
	}
	// Document doesn't have SetStatus method, so we need to recreate a new document
	// or use a different approach
	return nil
}

func (m *MockDocRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.docs, id)
	return nil
}

func (m *MockDocRepo) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	count := 0
	for _, d := range m.docs {
		if d.ProjectID() == projectID {
			count++
		}
	}
	return int64(count), nil
}

// TestRAGService_Retrieve tests document retrieval
func TestRAGService_Retrieve(t *testing.T) {
	ctx := context.Background()
	vectorRepo := NewMockVectorRepo()
	chunkRepo := NewMockChunkRepo()
	docRepo := NewMockDocRepo()

	projectID := uuid.New()

	// Create test documents and chunks
	doc1, _ := knowledge.NewDocument(projectID, "PRD.pdf", knowledge.TypePRD, "https://example.com/prd.pdf", uuid.New())
	doc2, _ := knowledge.NewDocument(projectID, "API Spec.pdf", knowledge.TypeAPISpec, "https://example.com/api.pdf", uuid.New())
	docRepo.docs[doc1.ID()] = doc1
	docRepo.docs[doc2.ID()] = doc2

	chunk1 := knowledge.NewDocumentChunk(doc1.ID(), 0, "User login feature description")
	chunk2 := knowledge.NewDocumentChunk(doc1.ID(), 1, "Password validation rules")
	chunk3 := knowledge.NewDocumentChunk(doc2.ID(), 0, "POST /api/login endpoint")
	chunkRepo.chunks[chunk1.ID()] = chunk1
	chunkRepo.chunks[chunk2.ID()] = chunk2
	chunkRepo.chunks[chunk3.ID()] = chunk3

	service := NewRAGService(vectorRepo, chunkRepo, docRepo)

	tests := []struct {
		name    string
		req     *RetrieveRequest
		setup   func()
		wantErr error
		wantLen int
	}{
		{
			name: "successful retrieval with results",
			req: &RetrieveRequest{
				ProjectID: projectID,
				Query:     "user login functionality",
				TopK:      3,
			},
			setup: func() {
				vectorRepo.results = []*knowledge.DocumentChunk{chunk1, chunk2, chunk3}
			},
			wantErr: nil,
			wantLen: 3,
		},
		{
			name: "successful retrieval with topK limit",
			req: &RetrieveRequest{
				ProjectID: projectID,
				Query:     "login test",
				TopK:      2,
			},
			setup: func() {
				vectorRepo.results = []*knowledge.DocumentChunk{chunk1, chunk2}
			},
			wantErr: nil,
			wantLen: 2,
		},
		{
			name: "empty knowledge base",
			req: &RetrieveRequest{
				ProjectID: uuid.New(), // Different project with no documents
				Query:     "no results",
				TopK:      5,
			},
			setup: func() {
				vectorRepo.results = []*knowledge.DocumentChunk{}
			},
			wantErr: nil,
			wantLen: 0,
		},
		{
			name: "query too short",
			req: &RetrieveRequest{
				ProjectID: projectID,
				Query:     "hi",
				TopK:      5,
			},
			setup:   func() {},
			wantErr: errors.New("query must be at least 5 characters"),
		},
		{
			name: "nil project ID",
			req: &RetrieveRequest{
				ProjectID: uuid.Nil,
				Query:     "valid query",
				TopK:      5,
			},
			setup:   func() {},
			wantErr: errors.New("project ID cannot be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := service.Retrieve(ctx, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Retrieve() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Retrieve() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Retrieve() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Retrieve() returned nil result")
				return
			}

			if len(result.Chunks) != tt.wantLen {
				t.Errorf("Retrieve() chunks count = %v, want %v", len(result.Chunks), tt.wantLen)
			}
		})
	}
}

// TestRAGService_CalculateConfidence tests confidence calculation
func TestRAGService_CalculateConfidence(t *testing.T) {
	vectorRepo := NewMockVectorRepo()
	chunkRepo := NewMockChunkRepo()
	docRepo := NewMockDocRepo()
	service := NewRAGService(vectorRepo, chunkRepo, docRepo)

	tests := []struct {
		name      string
		chunks    []*RetrievedChunk
		wantLevel testcase.Confidence
	}{
		{
			name: "high confidence - multiple high scores",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.92},
				{SimilarityScore: 0.88},
				{SimilarityScore: 0.85},
			},
			wantLevel: testcase.ConfidenceHigh,
		},
		{
			name: "high confidence - two very high scores",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.95},
				{SimilarityScore: 0.90},
			},
			wantLevel: testcase.ConfidenceHigh,
		},
		{
			name: "medium confidence - one good score",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.75},
				{SimilarityScore: 0.40},
			},
			wantLevel: testcase.ConfidenceMedium,
		},
		{
			name: "medium confidence - single medium score",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.60},
			},
			wantLevel: testcase.ConfidenceMedium,
		},
		{
			name: "low confidence - single low score",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.45},
			},
			wantLevel: testcase.ConfidenceLow,
		},
		{
			name: "low confidence - no chunks",
			chunks:    []*RetrievedChunk{},
			wantLevel: testcase.ConfidenceLow,
		},
		{
			name: "low confidence - very low scores",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.30},
				{SimilarityScore: 0.25},
			},
			wantLevel: testcase.ConfidenceLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateConfidence(tt.chunks)

			if result != tt.wantLevel {
				t.Errorf("CalculateConfidence() = %v, want %v", result, tt.wantLevel)
			}
		})
	}
}

// TestRAGService_ConfidenceThresholds tests edge cases in confidence calculation
func TestRAGService_ConfidenceThresholds(t *testing.T) {
	vectorRepo := NewMockVectorRepo()
	chunkRepo := NewMockChunkRepo()
	docRepo := NewMockDocRepo()
	service := NewRAGService(vectorRepo, chunkRepo, docRepo)

	tests := []struct {
		name      string
		chunks    []*RetrievedChunk
		wantLevel testcase.Confidence
	}{
		{
			name: "exactly 0.8 threshold with 2 chunks",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.81},
				{SimilarityScore: 0.50},
			},
			wantLevel: testcase.ConfidenceHigh,
		},
		{
			name: "just below 0.8 threshold",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.79},
				{SimilarityScore: 0.60},
			},
			wantLevel: testcase.ConfidenceMedium,
		},
		{
			name: "exactly 0.5 threshold",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.50},
			},
			wantLevel: testcase.ConfidenceMedium,
		},
		{
			name: "just below 0.5 threshold",
			chunks: []*RetrievedChunk{
				{SimilarityScore: 0.49},
			},
			wantLevel: testcase.ConfidenceLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateConfidence(tt.chunks)

			if result != tt.wantLevel {
				t.Errorf("CalculateConfidence() = %v, want %v", result, tt.wantLevel)
			}
		})
	}
}
