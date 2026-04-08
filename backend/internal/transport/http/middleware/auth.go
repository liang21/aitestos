// Package middleware provides HTTP middleware implementations
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// contextKey is the type for context keys
type contextKey string

// userContextKey is the context key for user ID
const userContextKey contextKey = "user_id"

// Auth creates an authentication middleware with optional role checking
func Auth(secret string, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Parse Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Parse and validate token
			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				respondWithError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Extract user ID from claims
			userIDStr, ok := (*claims)["user_id"].(string)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid user ID in token")
				return
			}

			// Check roles if specified
			if len(allowedRoles) > 0 {
				role, ok := (*claims)["role"].(string)
				if !ok {
					respondWithError(w, http.StatusForbidden, "missing role in token")
					return
				}

				roleAllowed := false
				for _, allowedRole := range allowedRoles {
					if role == allowedRole {
						roleAllowed = true
						break
					}
				}

				if !roleAllowed {
					respondWithError(w, http.StatusForbidden, "insufficient permissions")
					return
				}
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), userContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext extracts user ID from context
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userContextKey).(uuid.UUID)
	return userID, ok
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, httpStatus int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	response := map[string]interface{}{
		"error": message,
	}
	json.NewEncoder(w).Encode(response)
}
