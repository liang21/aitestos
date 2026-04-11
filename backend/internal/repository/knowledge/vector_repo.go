// Package knowledge provides vector repository implementation
package knowledge

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	domainknowledge "github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/infrastructure/milvus"
)

// MilvusClient defines the interface for Milvus operations we need
type MilvusClient interface {
	Insert(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error)
	Search(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, vectors []entity.Vector, vectorField string, metricType entity.MetricType, topK int, sp entity.SearchParam) ([]client.SearchResult, error)
	Delete(ctx context.Context, collectionName string, partitionName string, expr string) error
	Query(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, opts ...client.SearchQueryOptionFunc) (client.ResultSet, error)
	Flush(ctx context.Context, collName string, async bool) error
}

// milvusClientAdapter wraps *milvus.Client to implement our MilvusClient interface
type milvusClientAdapter struct {
	client *milvus.Client
}

func (a *milvusClientAdapter) Insert(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error) {
	return a.client.Insert(ctx, collName, partitionName, columns...)
}

func (a *milvusClientAdapter) Search(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, vectors []entity.Vector, vectorField string, metricType entity.MetricType, topK int, sp entity.SearchParam) ([]client.SearchResult, error) {
	return a.client.Search(ctx, collectionName, partitionNames, expr, outputFields, vectors, vectorField, metricType, topK, sp)
}

func (a *milvusClientAdapter) Delete(ctx context.Context, collectionName string, partitionName string, expr string) error {
	return a.client.Delete(ctx, collectionName, partitionName, expr)
}

func (a *milvusClientAdapter) Query(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, opts ...client.SearchQueryOptionFunc) (client.ResultSet, error) {
	return a.client.Query(ctx, collectionName, partitionNames, expr, outputFields, opts...)
}

func (a *milvusClientAdapter) Flush(ctx context.Context, collName string, async bool) error {
	return a.client.Flush(ctx, collName, async)
}

// VectorRepository implements domainknowledge.VectorRepository using Milvus
type VectorRepository struct {
	client     MilvusClient
	collection string
}

// NewVectorRepository creates a new vector repository
func NewVectorRepository(client MilvusClient, collection string) *VectorRepository {
	return &VectorRepository{
		client:     client,
		collection: collection,
	}
}

// NewVectorRepositoryFromMilvusClient creates a VectorRepository from *milvus.Client
func NewVectorRepositoryFromMilvusClient(milvusClient *milvus.Client, collection string) *VectorRepository {
	return NewVectorRepository(&milvusClientAdapter{client: milvusClient}, collection)
}

// Upsert inserts or updates vectors for document chunks
func (r *VectorRepository) Upsert(ctx context.Context, chunks []*domainknowledge.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	// Prepare data for bulk insert
	ids := make([]string, 0, len(chunks))
	documentIDs := make([]string, 0, len(chunks))
	projectIDs := make([]string, 0, len(chunks))
	chunkIndexes := make([]int64, 0, len(chunks))
	embeddings := make([][]float32, 0, len(chunks))
	contents := make([]string, 0, len(chunks))

	for _, chunk := range chunks {
		ids = append(ids, chunk.ID().String())
		documentIDs = append(documentIDs, chunk.DocumentID().String())
		projectIDs = append(projectIDs, chunk.ProjectID().String())
		chunkIndexes = append(chunkIndexes, int64(chunk.ChunkIndex()))
		contents = append(contents, chunk.Content())

		// Convert []byte embedding to []float32
		embedding := bytesToFloat32(chunk.Embedding())
		if embedding == nil {
			return fmt.Errorf("invalid embedding for chunk %s", chunk.ID())
		}
		embeddings = append(embeddings, embedding)
	}

	// Prepare column-based data
	columns := []entity.Column{
		entity.NewColumnVarChar("id", ids),
		entity.NewColumnVarChar("document_id", documentIDs),
		entity.NewColumnInt64("chunk_index", chunkIndexes),
		entity.NewColumnVarChar("project_id", projectIDs),
		entity.NewColumnFloatVector("embedding", 1536, embeddings),
		entity.NewColumnVarChar("content", contents),
	}

	// Insert data (ignore returned ID column as we use our own IDs)
	_, err := r.client.Insert(ctx, r.collection, "", columns...)
	if err != nil {
		return fmt.Errorf("insert vectors: %w", err)
	}

	// Flush to ensure persistence
	if err := r.client.Flush(ctx, r.collection, false); err != nil {
		return fmt.Errorf("flush collection: %w", err)
	}

	return nil
}

// Search performs vector similarity search
func (r *VectorRepository) Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*domainknowledge.DocumentChunk, error) {
	// Extract project_id from filter
	projectIDStr, ok := filter["project_id"].(string)
	if !ok {
		return nil, fmt.Errorf("project_id filter required")
	}

	// Validate project_id format
	if _, err := uuid.Parse(projectIDStr); err != nil {
		return nil, fmt.Errorf("invalid project_id format: %w", err)
	}

	// Build search parameters
	sp, err := entity.NewIndexHNSWSearchParam(64)
	if err != nil {
		return nil, fmt.Errorf("create search params: %w", err)
	}

	// Build vector search query
	vectors := []entity.Vector{entity.FloatVector(queryVector)}

	// Execute search with filter
	searchResult, err := r.client.Search(
		ctx,
		r.collection,
		[]string{}, // partitions
		fmt.Sprintf("project_id == '%s'", projectIDStr),
		[]string{"id", "document_id", "project_id", "chunk_index", "content"},
		vectors,
		"embedding",
		entity.L2,
		topK,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}

	// Parse results
	chunks := make([]*domainknowledge.DocumentChunk, 0, topK)
	for _, result := range searchResult {
		for i := 0; i < result.ResultCount; i++ {
			id := result.IDs.(*entity.ColumnVarChar).Data()[i]
			documentID := result.Fields.GetColumn("document_id").(*entity.ColumnVarChar).Data()[i]
			projectID := result.Fields.GetColumn("project_id").(*entity.ColumnVarChar).Data()[i]
			chunkIndex := int(result.Fields.GetColumn("chunk_index").(*entity.ColumnInt64).Data()[i])
			content := result.Fields.GetColumn("content").(*entity.ColumnVarChar).Data()[i]

			chunkUUID, err := uuid.Parse(id)
			if err != nil {
				continue
			}

			docUUID, err := uuid.Parse(documentID)
			if err != nil {
				continue
			}

			projUUID, err := uuid.Parse(projectID)
			if err != nil {
				continue
			}

			chunk := domainknowledge.ReconstructDocumentChunk(
				chunkUUID,
				docUUID,
				projUUID,
				chunkIndex,
				content,
				nil, // embedding not needed for search results
				time.Time{},
			)
			chunks = append(chunks, chunk)
		}
	}

	return chunks, nil
}

// DeleteByDocumentID removes all vectors for a document
func (r *VectorRepository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	expr := fmt.Sprintf("document_id == '%s'", documentID.String())

	if err := r.client.Delete(ctx, r.collection, "", expr); err != nil {
		return fmt.Errorf("delete vectors: %w", err)
	}

	// Flush to ensure deletion
	if err := r.client.Flush(ctx, r.collection, false); err != nil {
		return fmt.Errorf("flush after delete: %w", err)
	}

	return nil
}

// CountByProjectID counts vectors for a project
func (r *VectorRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	expr := fmt.Sprintf("project_id == '%s'", projectID.String())

	// Query with id field only since we just need count
	resultSet, err := r.client.Query(ctx, r.collection, []string{}, expr, []string{"id"})
	if err != nil {
		return 0, fmt.Errorf("count vectors: %w", err)
	}

	// Get count from result set length
	return int64(resultSet.Len()), nil
}

// bytesToFloat32 converts []byte to []float32
func bytesToFloat32(data []byte) []float32 {
	if len(data) == 0 {
		return nil
	}

	// Check if data length is valid (must be multiple of 4 for float32)
	if len(data)%4 != 0 {
		return nil
	}

	// Assuming embedding is stored as []float32 serialized to bytes
	vectors := make([]float32, 0, len(data)/4)
	for i := 0; i+4 <= len(data); i += 4 {
		bits := uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24
		vectors = append(vectors, math.Float32frombits(bits))
	}
	return vectors
}
