// Package handler provides HTTP handlers for the API
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/liang21/aitestos/internal/service/identity"
)

// IdentityHandler handles identity-related HTTP requests
type IdentityHandler struct {
	authService identity.AuthService
}

// NewIdentityHandler creates a new IdentityHandler
func NewIdentityHandler(authService identity.AuthService) *IdentityHandler {
	return &IdentityHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *IdentityHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req identity.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		respondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// Login handles user login
func (h *IdentityHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req identity.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// RefreshToken handles token refresh
func (h *IdentityHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}
