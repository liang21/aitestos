//go:build !integration
// +build !integration

package knowledge_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainknowledge "github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/repository/knowledge"
)

// MockMilvusClient is a mock implementation of MilvusClient interface
type MockMilvusClient struct {
	InsertFunc      func(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error)
	SearchFunc      func(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, vectors []entity.Vector, vectorField string, metricType entity.MetricType, topK int, sp entity.SearchParam) ([]client.SearchResult, error)
	DeleteFunc      func(ctx context.Context, collectionName string, partitionName string, expr string) error
	QueryFunc       func(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, opts ...client.SearchQueryOptionFunc) (client.ResultSet, error)
	FlushFunc       func(ctx context.Context, collName string, async bool) error
	InsertCallCount int
	SearchCallCount int
	DeleteCallCount int
	FlushCallCount  int
}

func (m *MockMilvusClient) Insert(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error) {
	m.InsertCallCount++
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, collName, partitionName, columns...)
	}
	return entity.NewColumnVarChar("id", []string{}), nil
}

func (m *MockMilvusClient) Search(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, vectors []entity.Vector, vectorField string, metricType entity.MetricType, topK int, sp entity.SearchParam) ([]client.SearchResult, error) {
	m.SearchCallCount++
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, collectionName, partitionNames, expr, outputFields, vectors, vectorField, metricType, topK, sp)
	}
	return []client.SearchResult{}, nil
}

func (m *MockMilvusClient) Delete(ctx context.Context, collectionName string, partitionName string, expr string) error {
	m.DeleteCallCount++
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, collectionName, partitionName, expr)
	}
	return nil
}

func (m *MockMilvusClient) Query(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, opts ...client.SearchQueryOptionFunc) (client.ResultSet, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, collectionName, partitionNames, expr, outputFields, opts...)
	}
	return client.ResultSet{}, nil
}

func (m *MockMilvusClient) Flush(ctx context.Context, collName string, async bool) error {
	m.FlushCallCount++
	if m.FlushFunc != nil {
		return m.FlushFunc(ctx, collName, async)
	}
	return nil
}

// Helper function to create test chunk with embedding for unit tests
func createTestChunkWithEmbeddingUnit(t *testing.T, docID uuid.UUID, projectID uuid.UUID, index int, content string, embedding []byte) *domainknowledge.DocumentChunk {
	t.Helper()
	chunk, err := domainknowledge.NewDocumentChunk(docID, projectID, index, content)
	require.NoError(t, err)
	chunk.SetEmbedding(embedding)
	return chunk
}

// Helper function to convert float32 to bytes
func float32ToBytes(vectors []float32) []byte {
	data := make([]byte, len(vectors)*4)
	for i, v := range vectors {
		bits := uint32(v)
		data[i*4] = byte(bits)
		data[i*4+1] = byte(bits >> 8)
		data[i*4+2] = byte(bits >> 16)
		data[i*4+3] = byte(bits >> 24)
	}
	return data
}

func TestVectorRepository_Upsert(t *testing.T) {
	tests := []struct {
		name    string
		chunks  []*domainknowledge.DocumentChunk
		setup   func(*MockMilvusClient)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful upsert with single chunk",
			chunks: func() []*domainknowledge.DocumentChunk {
				docID := uuid.New()
				projectID := uuid.New()
				embedding := float32ToBytes(make([]float32, 1536))
				return []*domainknowledge.DocumentChunk{
					createTestChunkWithEmbeddingUnit(t, docID, projectID, 0, "test content", embedding),
				}
			}(),
			setup: func(m *MockMilvusClient) {
				m.InsertFunc = func(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error) {
					return entity.NewColumnVarChar("id", []string{}), nil
				}
				m.FlushFunc = func(ctx context.Context, collName string, async bool) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "upsert with empty chunks",
			chunks: func() []*domainknowledge.DocumentChunk {
				return []*domainknowledge.DocumentChunk{}
			}(),
			setup: func(m *MockMilvusClient) {
				m.InsertFunc = func(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error) {
					return entity.NewColumnVarChar("id", []string{}), nil
				}
			},
			wantErr: false,
		},
		{
			name: "insert fails",
			chunks: func() []*domainknowledge.DocumentChunk {
				docID := uuid.New()
				projectID := uuid.New()
				embedding := float32ToBytes(make([]float32, 1536))
				return []*domainknowledge.DocumentChunk{
					createTestChunkWithEmbeddingUnit(t, docID, projectID, 0, "test content", embedding),
				}
			}(),
			setup: func(m *MockMilvusClient) {
				m.InsertFunc = func(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error) {
					return nil, errors.New("insert failed")
				}
			},
			wantErr: true,
			errMsg:  "insert vectors",
		},
		{
			name: "flush fails",
			chunks: func() []*domainknowledge.DocumentChunk {
				docID := uuid.New()
				projectID := uuid.New()
				embedding := float32ToBytes(make([]float32, 1536))
				return []*domainknowledge.DocumentChunk{
					createTestChunkWithEmbeddingUnit(t, docID, projectID, 0, "test content", embedding),
				}
			}(),
			setup: func(m *MockMilvusClient) {
				m.InsertFunc = func(ctx context.Context, collName string, partitionName string, columns ...entity.Column) (entity.Column, error) {
					return entity.NewColumnVarChar("id", []string{}), nil
				}
				m.FlushFunc = func(ctx context.Context, collName string, async bool) error {
					return errors.New("flush failed")
				}
			},
			wantErr: true,
			errMsg:  "flush collection",
		},
		{
			name: "invalid embedding size",
			chunks: func() []*domainknowledge.DocumentChunk {
				docID := uuid.New()
				projectID := uuid.New()
				embedding := []byte{1, 2, 3}
				return []*domainknowledge.DocumentChunk{
					createTestChunkWithEmbeddingUnit(t, docID, projectID, 0, "test content", embedding),
				}
			}(),
			setup: func(m *MockMilvusClient) {},
			wantErr: true,
			errMsg:  "invalid embedding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockMilvusClient{}
			tt.setup(mock)

			repo := knowledge.NewVectorRepository(mock, "test_collection")
			err := repo.Upsert(context.Background(), tt.chunks)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVectorRepository_Search(t *testing.T) {
	projectID := uuid.New()

	tests := []struct {
		name        string
		queryVector []float32
		topK        int
		filter      map[string]any
		setup       func(*MockMilvusClient)
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "successful search",
			queryVector: make([]float32, 1536),
			topK:        5,
			filter:      map[string]any{"project_id": projectID.String()},
			setup: func(m *MockMilvusClient) {
				m.SearchFunc = func(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, vectors []entity.Vector, vectorField string, metricType entity.MetricType, topK int, sp entity.SearchParam) ([]client.SearchResult, error) {
					id := uuid.New()
					docID := uuid.New()

					result := &client.SearchResult{
						IDs: entity.NewColumnVarChar("id", []string{id.String()}),
						Fields: []entity.Column{
							entity.NewColumnVarChar("document_id", []string{docID.String()}),
							entity.NewColumnVarChar("project_id", []string{projectID.String()}),
							entity.NewColumnInt64("chunk_index", []int64{0}),
							entity.NewColumnVarChar("content", []string{"test content"}),
						},
						ResultCount: 1,
					}
					return []client.SearchResult{*result}, nil
				}
			},
			wantErr: false,
		},
		{
			name:        "missing project_id filter",
			queryVector: make([]float32, 1536),
			topK:        5,
			filter:      map[string]any{},
			setup:       func(m *MockMilvusClient) {},
			wantErr:     true,
			errMsg:      "project_id filter required",
		},
		{
			name:        "search fails",
			queryVector: make([]float32, 1536),
			topK:        5,
			filter:      map[string]any{"project_id": projectID.String()},
			setup: func(m *MockMilvusClient) {
				m.SearchFunc = func(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, vectors []entity.Vector, vectorField string, metricType entity.MetricType, topK int, sp entity.SearchParam) ([]client.SearchResult, error) {
					return nil, errors.New("search failed")
				}
			},
			wantErr: true,
			errMsg:  "vector search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockMilvusClient{}
			tt.setup(mock)

			repo := knowledge.NewVectorRepository(mock, "test_collection")
			results, err := repo.Search(context.Background(), tt.queryVector, tt.topK, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, results)
			}
		})
	}
}

func TestVectorRepository_DeleteByDocumentID(t *testing.T) {
	docID := uuid.New()

	tests := []struct {
		name    string
		docID   uuid.UUID
		setup   func(*MockMilvusClient)
		wantErr bool
		errMsg  string
	}{
		{
			name:  "successful delete",
			docID: docID,
			setup: func(m *MockMilvusClient) {
				m.DeleteFunc = func(ctx context.Context, collectionName string, partitionName string, expr string) error {
					return nil
				}
				m.FlushFunc = func(ctx context.Context, collName string, async bool) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:  "delete fails",
			docID: docID,
			setup: func(m *MockMilvusClient) {
				m.DeleteFunc = func(ctx context.Context, collectionName string, partitionName string, expr string) error {
					return errors.New("delete failed")
				}
			},
			wantErr: true,
			errMsg:  "delete vectors",
		},
		{
			name:  "flush fails after delete",
			docID: docID,
			setup: func(m *MockMilvusClient) {
				m.DeleteFunc = func(ctx context.Context, collectionName string, partitionName string, expr string) error {
					return nil
				}
				m.FlushFunc = func(ctx context.Context, collName string, async bool) error {
					return errors.New("flush failed")
				}
			},
			wantErr: true,
			errMsg:  "flush after delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockMilvusClient{}
			tt.setup(mock)

			repo := knowledge.NewVectorRepository(mock, "test_collection")
			err := repo.DeleteByDocumentID(context.Background(), tt.docID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVectorRepository_CountByProjectID(t *testing.T) {
	projectID := uuid.New()

	tests := []struct {
		name        string
		projectID   uuid.UUID
		setup       func(*MockMilvusClient)
		wantCount   int64
		wantErr     bool
		errMsg      string
	}{
		{
			name:      "successful count",
			projectID: projectID,
			setup: func(m *MockMilvusClient) {
				m.QueryFunc = func(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, opts ...client.SearchQueryOptionFunc) (client.ResultSet, error) {
					col := entity.NewColumnVarChar("id", []string{uuid.New().String()})
					return client.ResultSet{col}, nil
				}
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "count with zero results",
			projectID: projectID,
			setup: func(m *MockMilvusClient) {
				m.QueryFunc = func(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, opts ...client.SearchQueryOptionFunc) (client.ResultSet, error) {
					col := entity.NewColumnVarChar("id", []string{})
					return client.ResultSet{col}, nil
				}
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "count fails",
			projectID: projectID,
			setup: func(m *MockMilvusClient) {
				m.QueryFunc = func(ctx context.Context, collectionName string, partitionNames []string, expr string, outputFields []string, opts ...client.SearchQueryOptionFunc) (client.ResultSet, error) {
					return client.ResultSet{}, errors.New("count failed")
				}
			},
			wantErr: true,
			errMsg:  "count vectors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockMilvusClient{}
			tt.setup(mock)

			repo := knowledge.NewVectorRepository(mock, "test_collection")
			count, err := repo.CountByProjectID(context.Background(), tt.projectID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
			}
		})
	}
}
