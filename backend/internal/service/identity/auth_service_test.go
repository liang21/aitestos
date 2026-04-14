// Package identity provides authentication services
package identity

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
	redispkg "github.com/liang21/aitestos/internal/infrastructure/redis"
)

// MockUserRepository implements identity.UserRepository for testing
type MockUserRepository struct {
	users      map[uuid.UUID]*identity.User
	emailIndex map[string]*identity.User
	nameIndex  map[string]*identity.User
	saveErr    error
	findErr    error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:      make(map[uuid.UUID]*identity.User),
		emailIndex: make(map[string]*identity.User),
		nameIndex:  make(map[string]*identity.User),
	}
}

func (m *MockUserRepository) Save(ctx context.Context, user *identity.User) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.users[user.ID()] = user
	m.emailIndex[user.Email()] = user
	m.nameIndex[user.Username()] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	user, ok := m.users[id]
	if !ok {
		return nil, identity.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*identity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	user, ok := m.emailIndex[email]
	if !ok {
		return nil, identity.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*identity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	user, ok := m.nameIndex[username]
	if !ok {
		return nil, identity.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *identity.User) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.users[user.ID()] = user
	m.emailIndex[user.Email()] = user
	m.nameIndex[user.Username()] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, opts identity.QueryOptions) ([]*identity.User, int64, error) {
	if m.findErr != nil {
		return nil, 0, m.findErr
	}
	users := make([]*identity.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, int64(len(users)), nil
}

// TestAuthService_Register tests user registration
func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret-key", redispkg.NewMockTokenStore())

	tests := []struct {
		name    string
		req     *RegisterRequest
		setup   func()
		wantErr error
	}{
		{
			name: "successful registration",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "normal",
			},
			setup:   func() {},
			wantErr: nil,
		},
		{
			name: "email already exists",
			req: &RegisterRequest{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
				Role:     "normal",
			},
			setup: func() {
				existingUser, _ := identity.NewUser("existing", "existing@example.com", "password123", identity.RoleNormal)
				mockRepo.users[existingUser.ID()] = existingUser
				mockRepo.emailIndex["existing@example.com"] = existingUser
			},
			wantErr: identity.ErrEmailDuplicate,
		},
		{
			name: "username already exists",
			req: &RegisterRequest{
				Username: "existing",
				Email:    "new@example.com",
				Password: "password123",
				Role:     "normal",
			},
			setup: func() {
				existingUser, _ := identity.NewUser("existing", "existing2@example.com", "password123", identity.RoleNormal)
				mockRepo.users[existingUser.ID()] = existingUser
				mockRepo.nameIndex["existing"] = existingUser
			},
			wantErr: identity.ErrUsernameDuplicate,
		},
		{
			name: "invalid email format",
			req: &RegisterRequest{
				Username: "testuser2",
				Email:    "invalid-email",
				Password: "password123",
				Role:     "normal",
			},
			setup:   func() {},
			wantErr: identity.ErrInvalidEmail,
		},
		{
			name: "password too short",
			req: &RegisterRequest{
				Username: "testuser3",
				Email:    "test3@example.com",
				Password: "short",
				Role:     "normal",
			},
			setup:   func() {},
			wantErr: identity.ErrPasswordTooShort,
		},
		{
			name: "invalid role",
			req: &RegisterRequest{
				Username: "testuser4",
				Email:    "test4@example.com",
				Password: "password123",
				Role:     "invalid_role",
			},
			setup:   func() {},
			wantErr: errors.New("invalid user role"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			user, err := service.Register(ctx, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Register() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("Register() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Register() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Error("Register() returned nil user")
				return
			}

			if user.Username() != tt.req.Username {
				t.Errorf("Register() username = %v, want %v", user.Username(), tt.req.Username)
			}
			if user.Email() != tt.req.Email {
				t.Errorf("Register() email = %v, want %v", user.Email(), tt.req.Email)
			}
			if !user.VerifyPassword(tt.req.Password) {
				t.Error("Register() password verification failed")
			}
		})
	}
}

// TestAuthService_Login tests user login
func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret-key", redispkg.NewMockTokenStore())

	// Create a test user
	testUser, _ := identity.NewUser("loginuser", "login@example.com", "correctpassword", identity.RoleNormal)
	mockRepo.users[testUser.ID()] = testUser
	mockRepo.emailIndex["login@example.com"] = testUser

	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr error
	}{
		{
			name: "successful login",
			req: &LoginRequest{
				Email:    "login@example.com",
				Password: "correctpassword",
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: &LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: identity.ErrUserNotFound,
		},
		{
			name: "wrong password",
			req: &LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			wantErr: identity.ErrPasswordMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Login(ctx, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Login() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("Login() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Login() unexpected error: %v", err)
				return
			}

			if resp == nil {
				t.Error("Login() returned nil response")
				return
			}

			if resp.AccessToken == "" {
				t.Error("Login() returned empty access token")
			}
			if resp.RefreshToken == "" {
				t.Error("Login() returned empty refresh token")
			}
			if resp.User == nil {
				t.Error("Login() returned nil user info")
			}
		})
	}
}

// TestAuthService_ValidateToken tests token validation
func TestAuthService_ValidateToken(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret-key", redispkg.NewMockTokenStore())

	// Create a test user
	testUser, _ := identity.NewUser("tokenuser", "token@example.com", "password123", identity.RoleNormal)
	mockRepo.users[testUser.ID()] = testUser
	mockRepo.emailIndex["token@example.com"] = testUser

	tests := []struct {
		name      string
		token     string
		setup     func() string // returns valid token if needed
		wantErr   error
		checkUser bool
	}{
		{
			name: "valid token",
			setup: func() string {
				// Login to get a valid token
				resp, _ := service.Login(ctx, &LoginRequest{
					Email:    "token@example.com",
					Password: "password123",
				})
				return resp.AccessToken
			},
			wantErr:   nil,
			checkUser: true,
		},
		{
			name:    "invalid token format",
			token:   "invalid-token",
			wantErr: identity.ErrTokenInvalid,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: identity.ErrTokenInvalid,
		},
		{
			name:    "expired token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2NzgtMTIzNC0xMjM0LTEyMzQtMTIzNDU2Nzg5YWJjIiwidXNlcm5hbWUiOiJ0ZXN0Iiwicm9sZSI6Im5vcm1hbCIsImV4cCI6MTAwMDAwMDAwMH0.invalid",
			wantErr: errors.New("invalid token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.token
			if tt.setup != nil {
				token = tt.setup()
			}

			claims, err := service.ValidateToken(ctx, token)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateToken() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("ValidateToken() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateToken() unexpected error: %v", err)
				return
			}

			if claims == nil {
				t.Error("ValidateToken() returned nil claims")
				return
			}

			if tt.checkUser {
				if claims.UserID != testUser.ID() {
					t.Errorf("ValidateToken() user ID = %v, want %v", claims.UserID, testUser.ID())
				}
				if claims.Username != testUser.Username() {
					t.Errorf("ValidateToken() username = %v, want %v", claims.Username, testUser.Username())
				}
			}
		})
	}
}

// TestAuthService_RefreshToken tests token refresh
func TestAuthService_RefreshToken(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret-key", redispkg.NewMockTokenStore())

	// Create a test user
	testUser, _ := identity.NewUser("refreshuser", "refresh@example.com", "password123", identity.RoleNormal)
	mockRepo.users[testUser.ID()] = testUser
	mockRepo.emailIndex["refresh@example.com"] = testUser

	tests := []struct {
		name    string
		setup   func() string
		wantErr error
	}{
		{
			name: "successful refresh",
			setup: func() string {
				resp, _ := service.Login(ctx, &LoginRequest{
					Email:    "refresh@example.com",
					Password: "password123",
				})
				return resp.RefreshToken
			},
			wantErr: nil,
		},
		{
			name:    "invalid refresh token",
			setup:   func() string { return "invalid-refresh-token" },
			wantErr: identity.ErrTokenInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshToken := tt.setup()

			resp, err := service.RefreshToken(ctx, refreshToken)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("RefreshToken() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("RefreshToken() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("RefreshToken() unexpected error: %v", err)
				return
			}

			if resp == nil {
				t.Error("RefreshToken() returned nil response")
				return
			}

			if resp.AccessToken == "" {
				t.Error("RefreshToken() returned empty access token")
			}
		})
	}
}
