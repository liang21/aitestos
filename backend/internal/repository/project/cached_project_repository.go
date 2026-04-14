// Package project provides cached project repository implementation
package project

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/infrastructure/cache"
)

// CachedProjectRepository wraps a ProjectRepository with caching functionality
type CachedProjectRepository struct {
	repo  domainproject.ProjectRepository
	cache cache.Cache
}

// NewCachedProjectRepository creates a new cached project repository
func NewCachedProjectRepository(
	repo domainproject.ProjectRepository,
	cache cache.Cache,
) domainproject.ProjectRepository {
	return &CachedProjectRepository{
		repo:  repo,
		cache: cache,
	}
}

// Save persists a new project and invalidates relevant caches
func (c *CachedProjectRepository) Save(ctx context.Context, project *domainproject.Project) error {
	// Write to database first
	if err := c.repo.Save(ctx, project); err != nil {
		return err
	}

	// Invalidate related caches
	cacheKey := fmt.Sprintf(cache.KeyProjectDetail, project.ID())
	_ = c.cache.Delete(ctx, cacheKey)
	_ = c.cache.Delete(ctx, cache.KeyProjectList)

	return nil
}

// FindByID retrieves a project by ID with caching
func (c *CachedProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainproject.Project, error) {
	// Try cache first
	var cached domainproject.Project
	cacheKey := fmt.Sprintf(cache.KeyProjectDetail, id)

	err := c.cache.Get(ctx, cacheKey, &cached)
	if err == nil {
		return &cached, nil
	}

	// Cache miss, fetch from database
	p, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Async write to cache (don't block main flow)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = c.cache.Set(ctx, cacheKey, p, cache.TTLLong)
	}()

	return p, nil
}

// FindByName retrieves a project by name (bypasses cache for uniqueness check)
func (c *CachedProjectRepository) FindByName(ctx context.Context, name string) (*domainproject.Project, error) {
	// Uniqueness checks should always go to database
	return c.repo.FindByName(ctx, name)
}

// FindByPrefix retrieves a project by prefix (bypasses cache for uniqueness check)
func (c *CachedProjectRepository) FindByPrefix(ctx context.Context, prefix domainproject.ProjectPrefix) (*domainproject.Project, error) {
	// Uniqueness checks should always go to database
	return c.repo.FindByPrefix(ctx, prefix)
}

// FindAll retrieves all projects with caching
func (c *CachedProjectRepository) FindAll(ctx context.Context, opts domainproject.QueryOptions) ([]*domainproject.Project, error) {
	// For simplicity, cache only simple list queries without filters
	if opts.Keywords == "" && opts.Limit == 0 && opts.Offset == 0 {
		var cached []*domainproject.Project
		err := c.cache.Get(ctx, cache.KeyProjectList, &cached)
		if err == nil {
			return cached, nil
		}
	}

	// Fetch from database
	projects, err := c.repo.FindAll(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Cache only simple list queries
	if opts.Keywords == "" && opts.Limit == 0 && opts.Offset == 0 {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = c.cache.Set(ctx, cache.KeyProjectList, projects, cache.TTLLong)
		}()
	}

	return projects, nil
}

// Update updates an existing project and invalidates caches
func (c *CachedProjectRepository) Update(ctx context.Context, project *domainproject.Project) error {
	// Update database first
	if err := c.repo.Update(ctx, project); err != nil {
		return err
	}

	// Invalidate related caches
	cacheKey := fmt.Sprintf(cache.KeyProjectDetail, project.ID())
	_ = c.cache.Delete(ctx, cacheKey)
	_ = c.cache.Delete(ctx, cache.KeyProjectList)

	return nil
}

// Delete removes a project and invalidates caches
func (c *CachedProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Delete from database first
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate related caches
	cacheKey := fmt.Sprintf(cache.KeyProjectDetail, id)
	_ = c.cache.Delete(ctx, cacheKey)
	_ = c.cache.Delete(ctx, cache.KeyProjectList)

	return nil
}

// GetStatistics retrieves project statistics (bypasses cache for real-time data)
func (c *CachedProjectRepository) GetStatistics(ctx context.Context, id uuid.UUID) (*domainproject.ProjectStatistics, error) {
	// Statistics should always be fresh from database
	return c.repo.GetStatistics(ctx, id)
}

// SetStatistics stores statistics in cache (used for warmup)
func (c *CachedProjectRepository) SetStatistics(ctx context.Context, id uuid.UUID, stats *domainproject.ProjectStatistics) error {
	cacheKey := cache.KeyProjectStats
	key := fmt.Sprintf(cacheKey, id)

	return c.cache.Set(ctx, key, stats, cache.TTLMedium)
}
