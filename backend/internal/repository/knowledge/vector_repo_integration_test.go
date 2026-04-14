//go:build !short
// +build !short

package knowledge_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liang21/aitestos/internal/config"
	domainknowledge "github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/infrastructure/milvus"
	"github.com/liang21/aitestos/internal/repository/knowledge"
)

// Helper function to create test chunk with embedding for integration tests
func createTestChunkWithEmbedding(t *testing.T, docID uuid.UUID, projectID uuid.UUID, index int, content string, embedding []float32) *domainknowledge.DocumentChunk {
	t.Helper()
	chunk, err := domainknowledge.NewDocumentChunk(docID, projectID, index, content)
	require.NoError(t, err)

	// Convert float32 to bytes
	embeddingBytes := make([]byte, len(embedding)*4)
	for i, v := range embedding {
		bits := uint32(v)
		embeddingBytes[i*4] = byte(bits)
		embeddingBytes[i*4+1] = byte(bits >> 8)
		embeddingBytes[i*4+2] = byte(bits >> 16)
		embeddingBytes[i*4+3] = byte(bits >> 24)
	}
	chunk.SetEmbedding(embeddingBytes)
	return chunk
}

// Helper function to generate random embedding
func generateRandomEmbedding(dim int) []float32 {
	embedding := make([]float32, dim)
	for i := range embedding {
		embedding[i] = float32(i) / float32(dim)
	}
	return embedding
}

func TestVectorRepository_Integration_Upsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This test requires a running Milvus instance
	// Use docker-compose to start Milvus before running tests

	ctx := context.Background()

	// Create Milvus client configuration
	cfg := &config.MilvusConfig{
		Host:       "localhost",
		Port:       19530,
		Database:   "aitestos_test",
		Collection: "test_chunks",
	}

	// Create client
	client, err := milvus.NewClient(cfg)
	if err != nil {
		t.Skipf("cannot connect to Milvus: %v (start with docker-compose up milvus)", err)
	}
	defer func() {
		// Clean up collection
		_ = client.DropCollection(ctx, cfg.Collection)
		_ = client.Close()
	}()

	// Ensure collection exists
	require.NoError(t, client.EnsureCollection(ctx))

	// Create repository
	repo := knowledge.NewVectorRepositoryFromMilvusClient(client, cfg.Collection)

	t.Run("upsert single chunk", func(t *testing.T) {
		docID := uuid.New()
		projectID := uuid.New()
		embedding := generateRandomEmbedding(1536)

		chunk := createTestChunkWithEmbedding(t, docID, projectID, 0, "test content", embedding)

		err := repo.Upsert(ctx, []*domainknowledge.DocumentChunk{chunk})
		require.NoError(t, err)

		// Wait for index to be updated
		time.Sleep(1 * time.Second)

		// Verify count
		count, err := repo.CountByProjectID(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("upsert multiple chunks", func(t *testing.T) {
		docID := uuid.New()
		projectID := uuid.New()

		chunks := make([]*domainknowledge.DocumentChunk, 5)
		for i := 0; i < 5; i++ {
			embedding := generateRandomEmbedding(1536)
			chunks[i] = createTestChunkWithEmbedding(t, docID, projectID, i, fmt.Sprintf("content %d", i), embedding)
		}

		err := repo.Upsert(ctx, chunks)
		require.NoError(t, err)

		// Wait for index to be updated
		time.Sleep(1 * time.Second)

		// Verify count
		count, err := repo.CountByProjectID(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestVectorRepository_Integration_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	cfg := &config.MilvusConfig{
		Host:       "localhost",
		Port:       19530,
		Database:   "aitestos_test",
		Collection: "test_chunks_search",
	}

	client, err := milvus.NewClient(cfg)
	if err != nil {
		t.Skipf("cannot connect to Milvus: %v (start with docker-compose up milvus)", err)
	}
	defer func() {
		_ = client.DropCollection(ctx, cfg.Collection)
		_ = client.Close()
	}()

	require.NoError(t, client.EnsureCollection(ctx))
	repo := knowledge.NewVectorRepositoryFromMilvusClient(client, cfg.Collection)

	// Insert test data
	projectID := uuid.New()
	docID := uuid.New()

	queryEmbedding := generateRandomEmbedding(1536)

	// Insert chunks with different embeddings
	chunks := make([]*domainknowledge.DocumentChunk, 3)
	for i := 0; i < 3; i++ {
		embedding := make([]float32, 1536)
		for j := range embedding {
			// Create slightly different embeddings
			embedding[j] = queryEmbedding[j] + float32(i)*0.01
		}
		chunks[i] = createTestChunkWithEmbedding(t, docID, projectID, i, fmt.Sprintf("content %d", i), embedding)
	}

	err = repo.Upsert(ctx, chunks)
	require.NoError(t, err)

	// Wait for index to be updated
	time.Sleep(2 * time.Second)

	t.Run("search with project filter", func(t *testing.T) {
		results, err := repo.Search(ctx, queryEmbedding, 3, map[string]any{
			"project_id": projectID.String(),
		})
		require.NoError(t, err)
		assert.NotEmpty(t, results)
		assert.LessOrEqual(t, len(results), 3)
	})

	t.Run("search with topK limit", func(t *testing.T) {
		results, err := repo.Search(ctx, queryEmbedding, 2, map[string]any{
			"project_id": projectID.String(),
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(results), 2)
	})

	t.Run("search with no matching project", func(t *testing.T) {
		otherProjectID := uuid.New()
		results, err := repo.Search(ctx, queryEmbedding, 3, map[string]any{
			"project_id": otherProjectID.String(),
		})
		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestVectorRepository_Integration_DeleteByDocumentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	cfg := &config.MilvusConfig{
		Host:       "localhost",
		Port:       19530,
		Database:   "aitestos_test",
		Collection: "test_chunks_delete",
	}

	client, err := milvus.NewClient(cfg)
	if err != nil {
		t.Skipf("cannot connect to Milvus: %v (start with docker-compose up milvus)", err)
	}
	defer func() {
		_ = client.DropCollection(ctx, cfg.Collection)
		_ = client.Close()
	}()

	require.NoError(t, client.EnsureCollection(ctx))
	repo := knowledge.NewVectorRepositoryFromMilvusClient(client, cfg.Collection)

	// Insert test data
	projectID := uuid.New()
	docID := uuid.New()

	chunks := make([]*domainknowledge.DocumentChunk, 3)
	for i := 0; i < 3; i++ {
		embedding := generateRandomEmbedding(1536)
		chunks[i] = createTestChunkWithEmbedding(t, docID, projectID, i, fmt.Sprintf("content %d", i), embedding)
	}

	err = repo.Upsert(ctx, chunks)
	require.NoError(t, err)

	// Wait for index to be updated
	time.Sleep(1 * time.Second)

	// Verify initial count
	count, err := repo.CountByProjectID(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// Delete by document ID
	err = repo.DeleteByDocumentID(ctx, docID)
	require.NoError(t, err)

	// Wait for deletion to complete
	time.Sleep(1 * time.Second)

	// Verify deletion
	count, err = repo.CountByProjectID(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestVectorRepository_Integration_CountByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	cfg := &config.MilvusConfig{
		Host:       "localhost",
		Port:       19530,
		Database:   "aitestos_test",
		Collection: "test_chunks_count",
	}

	client, err := milvus.NewClient(cfg)
	if err != nil {
		t.Skipf("cannot connect to Milvus: %v (start with docker-compose up milvus)", err)
	}
	defer func() {
		_ = client.DropCollection(ctx, cfg.Collection)
		_ = client.Close()
	}()

	require.NoError(t, client.EnsureCollection(ctx))
	repo := knowledge.NewVectorRepositoryFromMilvusClient(client, cfg.Collection)

	t.Run("count with no chunks", func(t *testing.T) {
		projectID := uuid.New()
		count, err := repo.CountByProjectID(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("count with chunks", func(t *testing.T) {
		projectID := uuid.New()
		docID := uuid.New()

		// Insert chunks
		chunks := make([]*domainknowledge.DocumentChunk, 10)
		for i := 0; i < 10; i++ {
			embedding := generateRandomEmbedding(1536)
			chunks[i] = createTestChunkWithEmbedding(t, docID, projectID, i, fmt.Sprintf("content %d", i), embedding)
		}

		err := repo.Upsert(ctx, chunks)
		require.NoError(t, err)

		// Wait for index to be updated
		time.Sleep(1 * time.Second)

		count, err := repo.CountByProjectID(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, int64(10), count)
	})

	t.Run("count with multiple documents", func(t *testing.T) {
		projectID := uuid.New()

		// Insert chunks from multiple documents
		for docIdx := 0; docIdx < 3; docIdx++ {
			docID := uuid.New()
			chunks := make([]*domainknowledge.DocumentChunk, 5)
			for i := 0; i < 5; i++ {
				embedding := generateRandomEmbedding(1536)
				chunks[i] = createTestChunkWithEmbedding(t, docID, projectID, i, fmt.Sprintf("doc %d content %d", docIdx, i), embedding)
			}
			err := repo.Upsert(ctx, chunks)
			require.NoError(t, err)
		}

		// Wait for index to be updated
		time.Sleep(2 * time.Second)

		count, err := repo.CountByProjectID(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, int64(15), count) // 3 docs * 5 chunks
	})
}
