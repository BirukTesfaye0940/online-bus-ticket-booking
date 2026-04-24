package handler

import (
	"encoding/json"
	"net/http"

	"github.com/biruk/bus-ticket/api-gateway/internal/middleware"
	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto"
)

// AuthHandler holds the gRPC client for the Auth Service.
type AuthHandler struct {
	authClient pb.AuthServiceClient
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(authClient pb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

// --- Request / Response DTOs ---

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// --- Handlers ---

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.authClient.Register(r.Context(), &pb.RegisterRequest{
		Email:       req.Email,
		Password:    req.Password,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Role:        pb.Role_ROLE_PASSENGER, // default role for self-registration
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.authClient.Login(r.Context(), &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetMe handles GET /api/v1/auth/me (protected)
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	resp, err := h.authClient.GetUser(r.Context(), &pb.GetUserRequest{UserId: userID})
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
