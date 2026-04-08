// Package knowledge defines repository interfaces
package knowledge

import (
	"context"

	"github.com/google/uuid"
)

// DocumentRepository defines the interface for document persistence
type DocumentRepository interface {
	// Save persists a new document
	Save(ctx context.Context, doc *Document) error

	// FindByID retrieves a document by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Document, error)

	// FindByProjectID retrieves all documents for a project with pagination
	FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*Document, error)

	// Update updates an existing document
	Update(ctx context.Context, doc *Document) error

	// UpdateStatus updates the document processing status
	UpdateStatus(ctx context.Context, id uuid.UUID, status DocumentStatus) error

	// Delete removes a document
	Delete(ctx context.Context, id uuid.UUID) error

	// CountByProjectID counts documents for a project
	CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error)
}

// DocumentChunkRepository defines the interface for document chunk persistence
type DocumentChunkRepository interface {
	// SaveBatch persists multiple document chunks
	SaveBatch(ctx context.Context, chunks []*DocumentChunk) error

	// FindByDocumentID retrieves all chunks for a document
	FindByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*DocumentChunk, error)

	// DeleteByDocumentID removes all chunks for a document
	DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error

	// CountByDocumentID counts chunks for a document
	CountByDocumentID(ctx context.Context, documentID uuid.UUID) (int64, error)
}

// VectorRepository defines the interface for vector storage operations
type VectorRepository interface {
	// Upsert inserts or updates vectors for document chunks
	Upsert(ctx context.Context, chunks []*DocumentChunk) error

	// Search performs vector similarity search
	Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*DocumentChunk, error)

	// DeleteByDocumentID removes all vectors for a document
	DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error

	// CountByProjectID counts vectors for a project
	CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error)
}

// QueryOptions holds pagination and filtering options
type QueryOptions struct {
	Offset   int
	Limit    int
	OrderBy  string
	Keywords string
	Type     DocumentType
	Status   DocumentStatus
}
