package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshua468/p2p-payment-gateway/auth-service/internal/service"
)

type AuthHandler interface {
	RegisterRoutes(router *mux.Router)
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{service: service}
}

func (h *authHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/signup", h.Signup).Methods(http.MethodPost)
	router.HandleFunc("/login", h.Login).Methods(http.MethodPost)
}

func (h *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req service.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Signup(req); err != nil {
		http.Error(w, "Failed to signup", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		http.Error(w, "Failed to login", http.StatusUnauthorized)
		return
	}

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
