// Package cache defines cache-specific errors
package cache

import "errors"

var (
	// ErrCacheNotFound is returned when a cache key does not exist
	ErrCacheNotFound = errors.New("cache not found")

	// ErrCacheWriteFailed is returned when writing to cache fails
	ErrCacheWriteFailed = errors.New("cache write failed")

	// ErrCacheExpired is returned when a cached value has expired
	ErrCacheExpired = errors.New("cache expired")

	// ErrCacheConnFailed is returned when cache connection fails
	ErrCacheConnFailed = errors.New("cache connection failed")

	// ErrCacheTimeout is returned when a cache operation times out
	ErrCacheTimeout = errors.New("cache timeout")
)
