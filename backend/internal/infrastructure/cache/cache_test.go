// Package cache provides cache implementation tests
package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryCache(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()
	defer cache.Close(ctx)

	t.Run("Set and Get", func(t *testing.T) {
		key := "test:key:" + uuid.New().String()
		value := map[string]string{"name": "test", "value": "data"}

		err := cache.Set(ctx, key, value, TTLShort)
		require.NoError(t, err)

		var result map[string]string
		err = cache.Get(ctx, key, &result)
		require.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Get not found", func(t *testing.T) {
		key := "test:notfound:" + uuid.New().String()
		var result string

		err := cache.Get(ctx, key, &result)
		assert.ErrorIs(t, err, ErrCacheNotFound)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:delete:" + uuid.New().String()
		value := "delete-me"

		err := cache.Set(ctx, key, value, TTLShort)
		require.NoError(t, err)

		err = cache.Delete(ctx, key)
		require.NoError(t, err)

		var result string
		err = cache.Get(ctx, key, &result)
		assert.ErrorIs(t, err, ErrCacheNotFound)
	})

	t.Run("Exists", func(t *testing.T) {
		key := "test:exists:" + uuid.New().String()
		value := "exists-test"

		exists, err := cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.False(t, exists)

		err = cache.Set(ctx, key, value, TTLShort)
		require.NoError(t, err)

		exists, err = cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("TTL expiration", func(t *testing.T) {
		key := "test:ttl:" + uuid.New().String()
		value := "expires-soon"

		err := cache.Set(ctx, key, value, 10*time.Millisecond)
		require.NoError(t, err)

		var result string
		err = cache.Get(ctx, key, &result)
		require.NoError(t, err)

		time.Sleep(15 * time.Millisecond)

		err = cache.Get(ctx, key, &result)
		assert.ErrorIs(t, err, ErrCacheNotFound)
	})

	t.Run("DeleteByPattern", func(t *testing.T) {
		prefix := "test:pattern:" + uuid.New().String() + ":"

		keys := []string{}
		for i := 0; i < 3; i++ {
			key := prefix + uuid.New().String()
			keys = append(keys, key)
			err := cache.Set(ctx, key, i, TTLShort)
			require.NoError(t, err)
		}

		otherKey := "test:other:" + uuid.New().String()
		err := cache.Set(ctx, otherKey, "other", TTLShort)
		require.NoError(t, err)

		pattern := prefix + "*"
		err = cache.DeleteByPattern(ctx, pattern)
		require.NoError(t, err)

		for _, key := range keys {
			var result int
			err = cache.Get(ctx, key, &result)
			assert.ErrorIs(t, err, ErrCacheNotFound)
		}

		var result string
		err = cache.Get(ctx, otherKey, &result)
		require.NoError(t, err)
		assert.Equal(t, "other", result)
	})
}

func TestMetrics(t *testing.T) {
	m := NewMetrics()

	t.Run("Initial state", func(t *testing.T) {
		assert.Equal(t, int64(0), m.hits)
		assert.Equal(t, int64(0), m.misses)
		assert.Equal(t, float64(0), m.GetHitRate())
	})

	t.Run("Record hit", func(t *testing.T) {
		m := NewMetrics()
		m.RecordHit()
		m.RecordHit()

		hits, misses, _, _ := m.GetStats()
		assert.Equal(t, int64(2), hits)
		assert.Equal(t, int64(0), misses)
		assert.Equal(t, float64(1), m.GetHitRate())
	})

	t.Run("Record miss", func(t *testing.T) {
		m := NewMetrics()
		m.RecordHit()
		m.RecordMiss()
		m.RecordMiss()

		hits, misses, _, _ := m.GetStats()
		assert.Equal(t, int64(1), hits)
		assert.Equal(t, int64(2), misses)
		assert.InDelta(t, 0.333, m.GetHitRate(), 0.01)
	})

	t.Run("Reset", func(t *testing.T) {
		m := NewMetrics()
		m.RecordHit()
		m.RecordMiss()
		m.Reset()

		hits, misses, _, _ := m.GetStats()
		assert.Equal(t, int64(0), hits)
		assert.Equal(t, int64(0), misses)
	})
}

func TestTTLConstants(t *testing.T) {
	assert.Equal(t, 5*time.Minute, TTLShort)
	assert.Equal(t, 15*time.Minute, TTLMedium)
	assert.Equal(t, 1*time.Hour, TTLLong)
}

func TestKeyFormat(t *testing.T) {
	assert.Contains(t, KeyProjectStats, "%s")
	assert.Contains(t, KeyProjectDetail, "%s")
	assert.Contains(t, KeyProjectList, "list")
	assert.Contains(t, KeyModuleList, "modules")
}
