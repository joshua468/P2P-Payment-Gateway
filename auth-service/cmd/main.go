package main

import (
	"log"
	"net/http"

	"github.com/joshua468/p2p-payment-gateway/auth-service/config"
	"github.com/joshua468/p2p-payment-gateway/auth-service/internal/handler"
	"github.com/joshua468/p2p-payment-gateway/auth-service/internal/middleware"
	"github.com/joshua468/p2p-payment-gateway/auth-service/internal/repository"
	"github.com/joshua468/p2p-payment-gateway/auth-service/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize repository, service, and handler
	authRepo := repository.NewAuthRepository(cfg)
	authService := service.NewAuthService(authRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	// Setup router and middleware
	router := mux.NewRouter()
	router.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // Apply middleware globally

	authHandler.RegisterRoutes(router)

	log.Printf("Auth service running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
