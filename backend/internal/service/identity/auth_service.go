// Package identity provides authentication services
package identity

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
)

// emailRegex validates email format
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// accessTokenExpiry is the duration before access tokens expire
const accessTokenExpiry = 15 * time.Minute

// refreshTokenExpiry is the duration before refresh tokens expire
const refreshTokenExpiry = 7 * 24 * time.Hour

// RegisterRequest contains user registration data
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=super_admin admin normal"`
}

// LoginRequest contains login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse contains authentication tokens
type LoginResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	User         *identity.UserJSON `json:"user"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// AuthService provides authentication operations
type AuthService interface {
	// Register creates a new user account
	Register(ctx context.Context, req *RegisterRequest) (*identity.User, error)

	// Login authenticates user and returns tokens
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

	// ValidateToken validates JWT token and returns claims
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)

	// RefreshToken refreshes access token using refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
}

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepo   identity.UserRepository
	jwtSecret  []byte
	tokenStore identity.TokenStore
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo identity.UserRepository, jwtSecret string, tokenStore identity.TokenStore) AuthService {
	return &AuthServiceImpl{
		userRepo:   userRepo,
		jwtSecret:  []byte(jwtSecret),
		tokenStore: tokenStore,
	}
}

// Register creates a new user account
func (s *AuthServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*identity.User, error) {
	// Validate email format
	if !emailRegex.MatchString(req.Email) {
		return nil, identity.ErrInvalidEmail
	}

	// Validate password length
	if len(req.Password) < 8 {
		return nil, identity.ErrPasswordTooShort
	}

	// Validate role
	role, err := identity.ParseUserRole(req.Role)
	if err != nil {
		return nil, identity.ErrInvalidRole
	}

	// Check if email already exists
	_, err = s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, identity.ErrEmailDuplicate
	}
	if !errors.Is(err, identity.ErrUserNotFound) {
		return nil, fmt.Errorf("check email existence: %w", err)
	}

	// Check if username already exists
	_, err = s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, identity.ErrUsernameDuplicate
	}
	if !errors.Is(err, identity.ErrUserNotFound) {
		return nil, fmt.Errorf("check username existence: %w", err)
	}

	// Create new user
	user, err := identity.NewUser(req.Username, req.Email, req.Password, role)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Save user
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	return user, nil
}

// Login authenticates user and returns tokens
func (s *AuthServiceImpl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, identity.ErrUserNotFound) {
			return nil, identity.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	// Verify password
	if !user.VerifyPassword(req.Password) {
		return nil, identity.ErrPasswordMismatch
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToJSON(),
	}, nil
}

// ValidateToken validates JWT token and returns claims
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	if token == "" {
		return nil, identity.ErrTokenInvalid
	}

	claims := &TokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, identity.ErrTokenInvalid
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, identity.ErrTokenExpired
		}
		return nil, identity.ErrTokenInvalid
	}

	if !parsedToken.Valid {
		return nil, identity.ErrTokenInvalid
	}

	return claims, nil
}

// RefreshToken refreshes access token using refresh token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// Validate refresh token format
	if refreshToken == "" {
		return nil, identity.ErrTokenInvalid
	}

	// Parse and validate refresh token
	claims := &TokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, identity.ErrTokenInvalid
		}
		return s.jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, identity.ErrTokenInvalid
	}

	// Verify token exists in store (not revoked)
	userID, expiresAt, found, err := s.tokenStore.Get(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("get token from store: %w", err)
	}
	if !found {
		return nil, identity.ErrTokenRevoked
	}
	if time.Now().After(expiresAt) {
		if delErr := s.tokenStore.Delete(ctx, refreshToken); delErr != nil {
			return nil, fmt.Errorf("delete expired token: %w", delErr)
		}
		return nil, identity.ErrTokenExpired
	}

	// Invalidate old refresh token
	if err := s.tokenStore.Delete(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("delete old token: %w", err)
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user.ToJSON(),
	}, nil
}

// generateAccessToken generates a new access token for a user
func (s *AuthServiceImpl) generateAccessToken(user *identity.User) (string, error) {
	now := time.Now()
	claims := &TokenClaims{
		UserID:   user.ID(),
		Username: user.Username(),
		Role:     string(user.Role()),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// generateRefreshToken generates a new refresh token for a user
func (s *AuthServiceImpl) generateRefreshToken(ctx context.Context, user *identity.User) (string, error) {
	now := time.Now()
	claims := &TokenClaims{
		UserID:   user.ID(),
		Username: user.Username(),
		Role:     string(user.Role()),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	// Store refresh token info
	if err := s.tokenStore.Store(ctx, signedToken, user.ID(), now.Add(refreshTokenExpiry)); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
	}

	return signedToken, nil
}
