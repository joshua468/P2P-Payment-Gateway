package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joshua468/p2p-payment-gateway/auth-service/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	authHandler := handler.NewAuthHandler(db, jwtSecret)

	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
