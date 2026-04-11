// Package handler provides HTTP handlers for the API
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/liang21/aitestos/internal/domain/identity"
	authservice "github.com/liang21/aitestos/internal/service/identity"
)

// MockAuthService implements authservice.AuthService for testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *authservice.RegisterRequest) (*identity.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*identity.User), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *authservice.LoginRequest) (*authservice.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authservice.LoginResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*authservice.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authservice.TokenClaims), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*authservice.LoginResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authservice.LoginResponse), args.Error(1)
}

func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful registration", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		mockSvc.On("Register", mock.Anything, mock.Anything).Return(&identity.User{}, nil).Run(func(args mock.Arguments) {
			req := args.Get(1).(*authservice.RegisterRequest)
			assert.Equal(t, "testuser", req.Username)
			assert.Equal(t, "test@example.com", req.Email)
		})

		handler := NewIdentityHandler(mockSvc)
		require.NotNil(t, handler)

		body := map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
			"role":     "normal",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		handler := NewIdentityHandler(mockSvc)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		handler := NewIdentityHandler(mockSvc)

		body := map[string]string{
			"username": "",
			"email":    "",
			"password": "",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful login", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		mockSvc.On("Login", mock.Anything, mock.Anything).Return(&authservice.LoginResponse{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			User:         &identity.UserJSON{ID: uuid.New()},
		}, nil)

		handler := NewIdentityHandler(mockSvc)

		body := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "access_token")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		mockSvc.On("Login", mock.Anything, mock.Anything).Return(nil, identity.ErrPasswordMismatch)

		handler := NewIdentityHandler(mockSvc)

		body := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		mockSvc.On("Login", mock.Anything, mock.Anything).Return(nil, identity.ErrUserNotFound)

		handler := NewIdentityHandler(mockSvc)

		body := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRefreshTokenHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful refresh", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		mockSvc.On("RefreshToken", mock.Anything, "valid_refresh_token").Return(&authservice.LoginResponse{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		}, nil)

		handler := NewIdentityHandler(mockSvc)

		body := map[string]string{
			"refresh_token": "valid_refresh_token",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.RefreshToken(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockAuthService)
		mockSvc.On("RefreshToken", mock.Anything, "invalid_token").Return(nil, errors.New("invalid refresh token"))

		handler := NewIdentityHandler(mockSvc)

		body := map[string]string{
			"refresh_token": "invalid_token",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.RefreshToken(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
