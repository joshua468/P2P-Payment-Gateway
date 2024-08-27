package internal

import (
	"encoding/json"
	"net/http"

	"github.com/joshua468/p2p-payment-gateway/auth-service/internal"
	"gorm.io/gorm"
)

type AuthHandler struct {
	AuthService *auth.AuthService
}

func NewAuthHandler(db *gorm.DB, jwtSecret string) *AuthHandler {
	authService := &auth.AuthService{
		DB:        db,
		JWTSecret: jwtSecret,
	}
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.AuthService.Register(r.Context(), req.Username, req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
