// Package middleware provides HTTP middleware implementations
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	secret := "test-secret-key-at-least-32-chars"
	validToken := generateTestToken(t, secret, uuid.New(), "testuser", "normal")

	tests := []struct {
		name       string
		token      string
		wantStatus int
		wantUserID bool
	}{
		{
			name:       "no token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
			wantUserID: false,
		},
		{
			name:       "malformed authorization header",
			token:      "InvalidFormat",
			wantStatus: http.StatusUnauthorized,
			wantUserID: false,
		},
		{
			name:       "invalid token",
			token:      "invalid.jwt.token",
			wantStatus: http.StatusUnauthorized,
			wantUserID: false,
		},
		{
			name:       "valid token",
			token:      "Bearer " + validToken,
			wantStatus: http.StatusOK,
			wantUserID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var ctxUserID uuid.UUID
			handler := Auth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if userID, ok := UserIDFromContext(r.Context()); ok {
					ctxUserID = userID
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}))

			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantUserID {
				assert.NotEqual(t, uuid.Nil, ctxUserID)
			}
		})
	}
}

func TestAuthMiddlewareWithRoles(t *testing.T) {
	t.Parallel()

	secret := "test-secret-key-at-least-32-chars"
	adminToken := generateTestToken(t, secret, uuid.New(), "admin", "admin")
	normalToken := generateTestToken(t, secret, uuid.New(), "normaluser", "normal")

	tests := []struct {
		name       string
		token      string
		roles      []string
		wantStatus int
	}{
		{
			name:       "admin accessing admin route",
			token:      "Bearer " + adminToken,
			roles:      []string{"admin", "super_admin"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "normal user accessing admin route",
			token:      "Bearer " + normalToken,
			roles:      []string{"admin", "super_admin"},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := Auth(secret, tt.roles...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/admin-only", nil)
			req.Header.Set("Authorization", tt.token)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestUserIDFromContext(t *testing.T) {
	t.Parallel()

	t.Run("no user in context", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userID, ok := UserIDFromContext(ctx)
		assert.False(t, ok)
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("user in context", func(t *testing.T) {
		t.Parallel()
		expectedID := uuid.New()
		ctx := context.WithValue(context.Background(), userContextKey, expectedID)
		userID, ok := UserIDFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, expectedID, userID)
	})
}

// generateTestToken creates a valid JWT token for testing
func generateTestToken(t *testing.T, secret string, userID uuid.UUID, username, role string) string {
	t.Helper()

	now := time.Now()
	claims := &jwt.MapClaims{
		"user_id":  userID.String(),
		"username": username,
		"role":     role,
		"exp":      jwt.NewNumericDate(now.Add(24 * time.Hour)),
		"iat":      jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	return tokenString
}
