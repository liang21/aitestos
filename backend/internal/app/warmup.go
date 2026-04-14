// Package app provides application initialization and warmup utilities
package app

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/infrastructure/cache"
)

// WarmupProjectStats preloads statistics for recent projects into cache
func WarmupProjectStats(ctx context.Context, projectRepo project.ProjectRepository, cacheClient cache.Cache, limit int) error {
	// Fetch recent projects
	opts := project.QueryOptions{
		Offset:  0,
		Limit:   limit,
		OrderBy: "created_at DESC",
	}

	projects, err := projectRepo.FindAll(ctx, opts)
	if err != nil {
		return fmt.Errorf("fetch projects for warmup: %w", err)
	}

	log.Info().Int("count", len(projects)).Msg("Starting project statistics warmup")

	warmed := 0
	for _, proj := range projects {
		stats, err := projectRepo.GetStatistics(ctx, proj.ID())
		if err != nil {
			log.Warn().Err(err).Str("project_id", proj.ID().String()).Msg("Failed to fetch stats for warmup")
			continue
		}

		// Store in cache
		cacheKey := cache.KeyProjectStats
		key := fmt.Sprintf(cacheKey, proj.ID())
		if err := cacheClient.Set(ctx, key, stats, cache.TTLMedium); err != nil {
			log.Warn().Err(err).Str("project_id", proj.ID().String()).Msg("Failed to cache stats")
			continue
		}

		warmed++
	}

	log.Info().Int("warmed", warmed).Int("total", len(projects)).Msg("Project statistics warmup completed")
	return nil
}
