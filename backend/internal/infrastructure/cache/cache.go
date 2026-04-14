// Package cache provides generic caching interfaces and implementations
package cache

import (
	"context"
	"errors"
	"time"
)

// Cache errors
var (
	ErrCacheNotFound    = errors.New("cache not found")
	ErrCacheWriteFailed = errors.New("cache write failed")
	ErrCacheExpired     = errors.New("cache expired")
	ErrCacheConnFailed  = errors.New("cache connection failed")
	ErrCacheTimeout     = errors.New("cache timeout")
)

// Cache defines a generic caching interface
type Cache interface {
	// Get retrieves a cached value and unmarshals it into dest
	// Returns ErrCacheNotFound if the key does not exist
	Get(ctx context.Context, key string, dest interface{}) error

	// Set stores a value in the cache with the specified TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes one or more keys from the cache
	Delete(ctx context.Context, keys ...string) error

	// DeleteByPattern removes all keys matching the given pattern (supports * wildcard)
	DeleteByPattern(ctx context.Context, pattern string) error

	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)

	// Close closes the cache connection and releases resources
	Close(ctx context.Context) error
}

// Cache key naming conventions
const (
	// Project-related keys
	KeyProjectStats  = "project:stats:%s"
	KeyProjectDetail = "project:detail:%s"
	KeyProjectList   = "project:list"
	KeyModuleList    = "project:%s:modules"

	// User-related keys
	KeyUserPermission = "user:perm:%s"

	// Generation-related keys
	KeyTaskStatus = "gen:task:status:%s"
)

// TTL constants
const (
	TTLShort  = 5 * time.Minute  // Hot data, high change frequency
	TTLMedium = 15 * time.Minute // Aggregate data, medium change frequency
	TTLLong   = 1 * time.Hour    // Metadata, low change frequency
)
