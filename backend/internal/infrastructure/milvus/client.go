// Package milvus provides Milvus vector database client wrapper
package milvus

import (
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/liang21/aitestos/internal/config"
)

// Client wraps Milvus SDK client with connection management
type Client struct {
	client.Client
	config *config.MilvusConfig
}

// NewClient creates a new Milvus client with connection pooling
func NewClient(cfg *config.MilvusConfig) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build address in format host:port
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	milvusClient, err := client.NewGrpcClient(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("create milvus client: %w", err)
	}

	return &Client{
		Client: milvusClient,
		config: cfg,
	}, nil
}

// EnsureCollection creates collection if not exists with proper schema and indexes
func (c *Client) EnsureCollection(ctx context.Context) error {
	has, err := c.HasCollection(ctx, c.config.Collection)
	if err != nil {
		return fmt.Errorf("check collection exists: %w", err)
	}

	if has {
		return nil
	}

	// Define schema
	schema := &entity.Schema{
		CollectionName: c.config.Collection,
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "36"},
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:       "document_id",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "36"},
			},
			{
				Name:     "chunk_index",
				DataType: entity.FieldTypeInt64,
			},
			{
				Name:       "project_id",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "36"},
			},
			{
				Name:       "embedding",
				DataType:   entity.FieldTypeFloatVector,
				TypeParams: map[string]string{"dim": "1536"},
			},
			{
				Name:       "content",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "65535"},
			},
		},
	}

	// Create collection with shard count
	shardNum := int32(2)
	if err := c.CreateCollection(ctx, schema, shardNum); err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	// Create index on embedding field using HNSW
	idx, err := entity.NewIndexHNSW(entity.L2, 16, 256)
	if err != nil {
		return fmt.Errorf("create index config: %w", err)
	}

	if err := c.CreateIndex(ctx, c.config.Collection, "embedding", idx, false); err != nil {
		return fmt.Errorf("create vector index: %w", err)
	}

	// Create scalar indexes for filtering
	scalarIdx, err := entity.NewIndexAUTOINDEX(entity.L2)
	if err != nil {
		return fmt.Errorf("create scalar index config: %w", err)
	}

	// Index on document_id
	if err := c.CreateIndex(ctx, c.config.Collection, "document_id", scalarIdx, false); err != nil {
		return fmt.Errorf("create document_id index: %w", err)
	}

	// Index on project_id
	if err := c.CreateIndex(ctx, c.config.Collection, "project_id", scalarIdx, false); err != nil {
		return fmt.Errorf("create project_id index: %w", err)
	}

	// Load collection into memory for search
	if err := c.LoadCollection(ctx, c.config.Collection, false); err != nil {
		return fmt.Errorf("load collection: %w", err)
	}

	return nil
}

// Close gracefully closes the connection
func (c *Client) Close() error {
	return c.Client.Close()
}
