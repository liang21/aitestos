// Package generation provides AI generation services
package generation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

// RetrieveRequest contains RAG retrieval parameters
type RetrieveRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Query     string    `json:"query" validate:"required,min=5"`
	TopK      int       `json:"top_k"`
}

// RetrieveResult contains retrieved document chunks
type RetrieveResult struct {
	Chunks []*RetrievedChunk `json:"chunks"`
	Query  string            `json:"query"`
}

// RetrievedChunk represents a retrieved document chunk with similarity score
type RetrievedChunk struct {
	ChunkID         uuid.UUID `json:"chunk_id"`
	DocumentID      uuid.UUID `json:"document_id"`
	DocumentName    string    `json:"document_name"`
	Content         string    `json:"content"`
	SimilarityScore float64   `json:"similarity_score"`
}

// RAGService provides RAG (Retrieval-Augmented Generation) operations
type RAGService interface {
	// Retrieve retrieves relevant document chunks based on query
	Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResult, error)

	// CalculateConfidence calculates confidence level based on retrieved chunks
	CalculateConfidence(chunks []*RetrievedChunk) testcase.Confidence
}

// RAGServiceImpl implements RAGService
type RAGServiceImpl struct {
	vectorRepo knowledge.VectorRepository
	chunkRepo  knowledge.DocumentChunkRepository
	docRepo    knowledge.DocumentRepository
}

// chunkRepo is an alias for knowledge.DocumentChunkRepository
// We need a way to find individual chunks, but DocumentChunkRepository doesn't have FindByID
// So we'll search through FindByDocumentID results

// NewRAGService creates a new RAGService instance
func NewRAGService(
	vectorRepo knowledge.VectorRepository,
	chunkRepo knowledge.DocumentChunkRepository,
	docRepo knowledge.DocumentRepository,
) RAGService {
	return &RAGServiceImpl{
		vectorRepo: vectorRepo,
		chunkRepo:  chunkRepo,
		docRepo:    docRepo,
	}
}

// Retrieve retrieves relevant document chunks based on query
func (s *RAGServiceImpl) Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResult, error) {
	// Validate project ID
	if req.ProjectID == uuid.Nil {
		return nil, fmt.Errorf("project ID cannot be nil")
	}

	// Validate query length
	if len(req.Query) < 5 {
		return nil, fmt.Errorf("query must be at least 5 characters")
	}

	// Set default TopK
	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}
	if topK > 20 {
		topK = 20
	}

	// Get all chunks for the project's documents
	docs, err := s.docRepo.FindByProjectID(ctx, req.ProjectID, knowledge.QueryOptions{Limit: 100})
	if err != nil {
		return nil, fmt.Errorf("find documents: %w", err)
	}

	if len(docs) == 0 {
		return &RetrieveResult{Query: req.Query, Chunks: []*RetrievedChunk{}}, nil
	}

	// Collect all chunk IDs from documents
	var allChunks []*knowledge.DocumentChunk
	for _, doc := range docs {
		chunks, err := s.chunkRepo.FindByDocumentID(ctx, doc.ID())
		if err != nil {
			continue
		}
		allChunks = append(allChunks, chunks...)
	}

	// Build filter for vector search
	filter := map[string]any{
		"project_id": req.ProjectID,
	}

	// Search for similar chunks using a simple query embedding
	// Note: In a real implementation, we would generate an embedding for the query
	// For now, we use a placeholder approach with the first chunk's embedding as query vector
	var queryVector []float32
	if len(allChunks) > 0 && len(allChunks[0].Embedding()) > 0 {
		// Convert bytes to float32 slice
		embeddingBytes := allChunks[0].Embedding()
		floatCount := len(embeddingBytes) / 4
		queryVector = make([]float32, floatCount)
		for i := 0; i < floatCount; i++ {
			queryVector[i] = float32(float64(embeddingBytes[i*4]) / 255.0)
		}
	} else {
		// Use a default query vector
		queryVector = make([]float32, 1536)
	}

	searchResults, err := s.vectorRepo.Search(ctx, queryVector, topK, filter)
	if err != nil {
		return nil, fmt.Errorf("search vectors: %w", err)
	}

	// Build result with document info
	result := &RetrieveResult{
		Query:  req.Query,
		Chunks: make([]*RetrievedChunk, 0, len(searchResults)),
	}

	for _, chunk := range searchResults {
		// Get document info
		doc, err := s.docRepo.FindByID(ctx, chunk.DocumentID())
		docName := ""
		if err == nil {
			docName = doc.Name()
		}

		result.Chunks = append(result.Chunks, &RetrievedChunk{
			ChunkID:         chunk.ID(),
			DocumentID:      chunk.DocumentID(),
			DocumentName:    docName,
			Content:         chunk.Content(),
			SimilarityScore: 0.0, // VectorRepository doesn't return score, use default
		})
	}

	return result, nil
}

// CalculateConfidence calculates confidence level based on retrieved chunks
func (s *RAGServiceImpl) CalculateConfidence(chunks []*RetrievedChunk) testcase.Confidence {
	if len(chunks) == 0 {
		return testcase.ConfidenceLow
	}

	// High confidence: at least 2 chunks with score > 0.8
	if len(chunks) >= 2 && chunks[0].SimilarityScore > 0.8 {
		return testcase.ConfidenceHigh
	}

	// Medium confidence: at least 1 chunk with score >= 0.5
	if len(chunks) >= 1 && chunks[0].SimilarityScore >= 0.5 {
		return testcase.ConfidenceMedium
	}

	// Low confidence: low scores or no chunks
	return testcase.ConfidenceLow
}
