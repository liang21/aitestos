// Package cache provides Redis-based cache implementation
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client  *redis.Client
	logger  *zerolog.Logger
	metrics *Metrics
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client, logger *zerolog.Logger) *RedisCache {
	return &RedisCache{
		client:  client,
		logger:  logger,
		metrics: NewMetrics(),
	}
}

// Get retrieves a cached value and unmarshals it into dest
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	start := time.Now()
	defer func() {
		c.metrics.RecordGet(time.Since(start))
	}()

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			c.metrics.RecordMiss()
			return ErrCacheNotFound
		}
		c.metrics.RecordError()
		return fmt.Errorf("cache get: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		c.metrics.RecordError()
		return fmt.Errorf("cache unmarshal: %w", err)
	}

	c.metrics.RecordHit()
	return nil
}

// Set stores a value in the cache with the specified TTL
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache set: %w", err)
	}

	return nil
}

// Delete removes one or more keys from the cache
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

// DeleteByPattern removes all keys matching the given pattern
func (c *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	keys := []string{}

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache scan: %w", err)
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists checks if a key exists in the cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Close closes the Redis connection
func (c *RedisCache) Close(ctx context.Context) error {
	return c.client.Close()
}

// GetMetrics returns the cache metrics
func (c *RedisCache) GetMetrics() *Metrics {
	return c.metrics
}
