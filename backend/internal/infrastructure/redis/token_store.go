// Package redis provides Redis client implementations
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
)

// TokenStore implements identity.TokenStore using Redis
type TokenStore struct {
	client *redis.Client
}

// NewTokenStore creates a new Redis-based token store
func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{client: client}
}

// tokenInfo represents the data stored for each refresh token
type tokenInfo struct {
	UserID    string    `json:"user_id"`
	ExpiresAt int64     `json:"expires_at"` // Unix timestamp
}

// Store saves a refresh token with its expiration time
func (s *TokenStore) Store(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	key := s.tokenKey(token)
	info := tokenInfo{
		UserID:    userID.String(),
		ExpiresAt: expiresAt.Unix(),
	}

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal token info: %w", err)
	}

	ttl := time.Until(expiresAt)
	if ttl < 0 {
		ttl = time.Hour // Default TTL if token has expired
	}

	pipe := s.client.Pipeline()
	pipe.Set(ctx, key, data, ttl)
	pipe.SAdd(ctx, s.userTokensKey(userID), key)
	// Set TTL on the user tokens set to the max refresh token expiry (7 days)
	pipe.Expire(ctx, s.userTokensKey(userID), 7*24*time.Hour+time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("store refresh token: %w", err)
	}

	return nil
}

// Get retrieves the user ID and expiration for a token
func (s *TokenStore) Get(ctx context.Context, token string) (uuid.UUID, time.Time, bool, error) {
	key := s.tokenKey(token)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, time.Time{}, false, nil
		}
		return uuid.Nil, time.Time{}, false, fmt.Errorf("get token: %w", err)
	}

	var info tokenInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return uuid.Nil, time.Time{}, false, fmt.Errorf("unmarshal token info: %w", err)
	}

	userID, err := uuid.Parse(info.UserID)
	if err != nil {
		return uuid.Nil, time.Time{}, false, fmt.Errorf("parse user ID: %w", err)
	}

	expiresAt := time.Unix(info.ExpiresAt, 0)
	return userID, expiresAt, true, nil
}

// Delete removes a token from the store
func (s *TokenStore) Delete(ctx context.Context, token string) error {
	key := s.tokenKey(token)
	// Best-effort removal from user set — we read userID from the stored data first
	data, err := s.client.Get(ctx, key).Bytes()
	if err == nil {
		var info tokenInfo
		if json.Unmarshal(data, &info) == nil {
			if uid, parseErr := uuid.Parse(info.UserID); parseErr == nil {
				s.client.SRem(ctx, s.userTokensKey(uid), key)
			}
		}
	}
	return s.client.Del(ctx, key).Err()
}

// DeleteAllByUser removes all tokens for a specific user
func (s *TokenStore) DeleteAllByUser(ctx context.Context, userID uuid.UUID) error {
	setKey := s.userTokensKey(userID)
	keys, err := s.client.SMembers(ctx, setKey).Result()
	if err != nil {
		return fmt.Errorf("find user tokens: %w", err)
	}
	if len(keys) == 0 {
		return nil
	}
	pipe := s.client.Pipeline()
	pipe.Del(ctx, keys...)
	pipe.Del(ctx, setKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("delete user tokens: %w", err)
	}
	return nil
}

// tokenKey generates the Redis key for a token
func (s *TokenStore) tokenKey(token string) string {
	return "auth:refresh_token:" + token
}

// userTokensKey generates the Redis set key for a user's tokens
func (s *TokenStore) userTokensKey(userID uuid.UUID) string {
	return "auth:user_tokens:" + userID.String()
}

// NewClient creates a new Redis client from configuration
func NewClient(addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}

// MockTokenStore is an in-memory implementation for testing
type MockTokenStore struct {
	tokens map[string]*mockTokenInfo
	mu      sync.RWMutex
}

type mockTokenInfo struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

// NewMockTokenStore creates a new mock token store
func NewMockTokenStore() *MockTokenStore {
	return &MockTokenStore{
		tokens: make(map[string]*mockTokenInfo),
	}
}

func (m *MockTokenStore) Store(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[token] = &mockTokenInfo{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return nil
}

func (m *MockTokenStore) Get(ctx context.Context, token string) (uuid.UUID, time.Time, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	info, ok := m.tokens[token]
	if !ok {
		return uuid.Nil, time.Time{}, false, nil
	}
	return info.UserID, info.ExpiresAt, true, nil
}

func (m *MockTokenStore) Delete(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, token)
	return nil
}

func (m *MockTokenStore) DeleteAllByUser(ctx context.Context, userID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for token, info := range m.tokens {
		if info.UserID == userID {
			delete(m.tokens, token)
		}
	}
	return nil
}
