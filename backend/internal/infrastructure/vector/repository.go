// Package vector provides vector repository implementation using Milvus
package vector

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/infrastructure/milvus"
)

// Repository implements vector storage operations
type Repository struct {
	client     *milvus.Client
	collection string
}

// NewRepository creates a new vector repository
func NewRepository(client *milvus.Client) *Repository {
	return &Repository{
		client:     client,
		collection: client.Config().Collection,
	}
}

// Upsert inserts or updates document chunks in the vector database
func (r *Repository) Upsert(ctx context.Context, chunks []*knowledge.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	// TODO: Implement actual Milvus upsert using the SDK
	// Reference: https://github.com/milvus-io/milvus-sdk-go/blob/main/docs/insert.md
	// Key steps:
	// 1. Prepare column data (ids, document_ids, chunk_indexes, project_ids, embeddings, contents)
	// 2. Call client.Upsert() with proper column format
	// 3. Call client.Flush() to persist

	return fmt.Errorf("vector Repository.Upsert: not yet implemented - Milvus SDK integration required")
}

// Search performs vector similarity search
func (r *Repository) Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*knowledge.DocumentChunk, error) {
	// TODO: Implement actual Milvus search using the SDK
	// Reference: https://github.com/milvus-io/milvus-sdk-go/blob/main/docs/search.md
	// Key steps:
	// 1. Build filter expression from filter map
	// 2. Call client.Search() with query vector and search parameters
	// 3. Convert results to DocumentChunk domain objects

	return nil, fmt.Errorf("vector Repository.Search: not yet implemented - Milvus SDK integration required")
}

// DeleteByDocumentID removes all chunks for a document
func (r *Repository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	// TODO: Implement actual Milvus delete using the SDK
	// Reference: https://github.com/milvus-io/milvus-sdk-go/blob/main/docs/delete.md
	// Key steps:
	// 1. Build delete expression for document_id
	// 2. Call client.Delete() with the expression
	// 3. Call client.Flush() to persist

	return fmt.Errorf("vector Repository.DeleteByDocumentID: not yet implemented - Milvus SDK integration required")
}

// CountByProjectID counts chunks for a project
func (r *Repository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	// TODO: Implement actual Milvus count query using the SDK
	// Reference: https://github.com/milvus-io/milvus-sdk-go/blob/main/docs/query.md
	// Key steps:
	// 1. Build query expression for project_id
	// 2. Call client.Query() with the expression
	// 3. Count returned results

	return 0, fmt.Errorf("vector Repository.CountByProjectID: not yet implemented - Milvus SDK integration required")
}

// Close closes the repository
func (r *Repository) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// Helper methods for future implementation

// buildFilterExpression creates a Milvus filter expression from the filter map
func (r *Repository) buildFilterExpression(filter map[string]any) string {
	if len(filter) == 0 {
		return ""
	}

	expr := ""
	if projectID, ok := filter["project_id"].(string); ok && projectID != "" {
		expr = fmt.Sprintf("project_id == '%s'", projectID)
	}
	if documentID, ok := filter["document_id"].(string); ok && documentID != "" {
		if expr != "" {
			expr += " && "
		}
		expr += fmt.Sprintf("document_id == '%s'", documentID)
	}

	return expr
}
