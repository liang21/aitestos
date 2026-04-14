// Package cache provides caching metrics tracking
package cache

import (
	"sync"
	"time"
)

// Metrics tracks cache performance statistics
type Metrics struct {
	mu         sync.RWMutex
	hits       int64
	misses     int64
	errors     int64
	getLatency time.Duration
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{}
}

// RecordHit records a cache hit
func (m *Metrics) RecordHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits++
}

// RecordMiss records a cache miss
func (m *Metrics) RecordMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.misses++
}

// RecordError records a cache operation error
func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors++
}

// RecordGet records the latency of a cache get operation
func (m *Metrics) RecordGet(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Simple moving average
	m.getLatency = (m.getLatency + d) / 2
}

// GetHitRate returns the cache hit rate (0.0 to 1.0)
func (m *Metrics) GetHitRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	total := m.hits + m.misses
	if total == 0 {
		return 0
	}
	return float64(m.hits) / float64(total)
}

// GetStats returns all metrics
func (m *Metrics) GetStats() (hits, misses, errors int64, latency time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hits, m.misses, m.errors, m.getLatency
}

// Reset clears all metrics
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits = 0
	m.misses = 0
	m.errors = 0
	m.getLatency = 0
}
