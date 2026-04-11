// Package rag provides Retrieval-Augmented Generation service implementations
package rag

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/config"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	domaintestcase "github.com/liang21/aitestos/internal/domain/testcase"
	generationservice "github.com/liang21/aitestos/internal/service/generation"
)

// Service implements RAG service operations
type Service struct {
	vectorRepo   VectorRepository
	chunkRepo    ChunkRepository
	documentRepo DocumentRepository
	llmService   LLMService
}

// VectorRepository defines the interface for vector storage operations
type VectorRepository interface {
	Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*knowledge.DocumentChunk, error)
}

// ChunkRepository defines document chunk repository operations
type ChunkRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*knowledge.DocumentChunk, error)
}

// DocumentRepository defines document repository operations
type DocumentRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*knowledge.Document, error)
}

// LLMService defines embedding generation operations
type LLMService interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

// NewService creates a new RAG service
func NewService(cfg *config.MilvusConfig, vectorRepo VectorRepository) (*Service, error) {
	if vectorRepo == nil {
		return nil, fmt.Errorf("vector repository is required")
	}

	return &Service{
		vectorRepo: vectorRepo,
	}, nil
}

// WithChunkRepo sets the chunk repository
func (s *Service) WithChunkRepo(repo ChunkRepository) *Service {
	s.chunkRepo = repo
	return s
}

// WithDocumentRepo sets the document repository
func (s *Service) WithDocumentRepo(repo DocumentRepository) *Service {
	s.documentRepo = repo
	return s
}

// WithLLMService sets the LLM service for embedding generation
func (s *Service) WithLLMService(llm LLMService) *Service {
	s.llmService = llm
	return s
}

// Retrieve retrieves relevant document chunks for a query using vector similarity search
func (s *Service) Retrieve(ctx context.Context, req *generationservice.RetrieveRequest) (*generationservice.RetrieveResult, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// Generate embedding for the query if LLM service is available
	var queryVector []float32
	var err error

	if s.llmService != nil {
		queryVector, err = s.llmService.GenerateEmbedding(ctx, req.Query)
		if err != nil {
			return nil, fmt.Errorf("generate query embedding: %w", err)
		}
	} else {
		// If no LLM service, RAG cannot work
		return nil, fmt.Errorf("LLM service required for RAG retrieval but not configured")
	}

	// Build filter
	filter := make(map[string]any)
	if req.ProjectID != uuid.Nil {
		filter["project_id"] = req.ProjectID.String()
	}

	// Perform vector search
	topK := req.TopK
	if topK <= 0 {
		topK = 5 // Default to top 5
	}

	chunks, err := s.vectorRepo.Search(ctx, queryVector, topK, filter)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}

	// Fetch document names for chunks
	retrievedChunks := make([]*generationservice.RetrievedChunk, len(chunks))
	for i, chunk := range chunks {
		// Default similarity score (would be returned by Milvus in a real implementation)
		similarity := 0.8

		// Try to get document name
		documentName := "Unknown"
		if s.documentRepo != nil {
			if doc, err := s.documentRepo.FindByID(ctx, chunk.DocumentID()); err == nil {
				documentName = doc.Name()
			}
		}

		retrievedChunks[i] = &generationservice.RetrievedChunk{
			ChunkID:         chunk.ID(),
			DocumentID:      chunk.DocumentID(),
			DocumentName:    documentName,
			Content:         chunk.Content(),
			SimilarityScore: similarity,
		}
	}

	return &generationservice.RetrieveResult{
		Chunks: retrievedChunks,
		Query:  req.Query,
	}, nil
}

// CalculateConfidence calculates confidence score for retrieved chunks
func (s *Service) CalculateConfidence(chunks []*generationservice.RetrievedChunk) domaintestcase.Confidence {
	if len(chunks) == 0 {
		return domaintestcase.ConfidenceLow
	}

	// Check average similarity score
	if len(chunks) > 0 {
		totalScore := 0.0
		for _, chunk := range chunks {
			totalScore += chunk.SimilarityScore
		}
		avgScore := totalScore / float64(len(chunks))

		// Confidence based on both count and similarity
		if avgScore >= 0.8 && len(chunks) >= 3 {
			return domaintestcase.ConfidenceHigh
		} else if avgScore >= 0.6 && len(chunks) >= 1 {
			return domaintestcase.ConfidenceMedium
		}
	}

	// Fallback to count-based confidence
	switch {
	case len(chunks) >= 5:
		return domaintestcase.ConfidenceHigh
	case len(chunks) >= 2:
		return domaintestcase.ConfidenceMedium
	default:
		return domaintestcase.ConfidenceLow
	}
}
