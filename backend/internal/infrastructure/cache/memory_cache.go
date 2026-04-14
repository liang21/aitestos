// Package cache provides in-memory cache implementation for fallback/testing
package cache

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sync"
	"time"
)

// memoryItem represents a cached item with expiration
type memoryItem struct {
	data      []byte
	expiresAt time.Time
}

// MemoryCache implements Cache interface using in-memory storage
type MemoryCache struct {
	items   map[string]*memoryItem
	mu      sync.RWMutex
	closed  bool
	stopCh  chan struct{}
}

// NewMemoryCache creates a new in-memory cache instance
func NewMemoryCache() *MemoryCache {
	mc := &MemoryCache{
		items:  make(map[string]*memoryItem),
		stopCh: make(chan struct{}),
	}
	go mc.cleanupExpired()
	return mc
}

// Get retrieves a cached value and unmarshals it into dest
func (m *MemoryCache) Get(ctx context.Context, key string, dest interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed || m.items == nil {
		return ErrCacheNotFound
	}

	item, ok := m.items[key]
	if !ok {
		return ErrCacheNotFound
	}

	if time.Now().After(item.expiresAt) {
		return ErrCacheNotFound
	}

	return json.Unmarshal(item.data, dest)
}

// Set stores a value in the cache with the specified TTL
func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed || m.items == nil {
		return nil // Silently fail if closed
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.items[key] = &memoryItem{
		data:      data,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

// Delete removes one or more keys from the cache
func (m *MemoryCache) Delete(ctx context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed || m.items == nil {
		return nil
	}

	for _, key := range keys {
		delete(m.items, key)
	}

	return nil
}

// DeleteByPattern removes all keys matching the given pattern
func (m *MemoryCache) DeleteByPattern(ctx context.Context, pattern string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed || m.items == nil {
		return nil
	}

	for key := range m.items {
		if matchPattern(key, pattern) {
			delete(m.items, key)
		}
	}

	return nil
}

// Exists checks if a key exists in the cache
func (m *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed || m.items == nil {
		return false, nil
	}

	item, ok := m.items[key]
	if !ok {
		return false, nil
	}

	return time.Now().Before(item.expiresAt), nil
}

// Close clears all cached items and stops cleanup goroutine
func (m *MemoryCache) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	m.closed = true
	m.items = nil
	close(m.stopCh)
	return nil
}

// cleanupExpired periodically removes expired items
func (m *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			if m.closed {
				m.mu.Unlock()
				return
			}
			now := time.Now()
			for key, item := range m.items {
				if now.After(item.expiresAt) {
					delete(m.items, key)
				}
			}
			m.mu.Unlock()
		case <-m.stopCh:
			return
		}
	}
}

// matchPattern checks if a key matches a glob pattern
func matchPattern(key, pattern string) bool {
	matched, err := filepath.Match(pattern, key)
	if err != nil {
		return false
	}
	return matched
}
