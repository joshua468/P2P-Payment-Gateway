package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/joshua468/p2p-payment-gateway/auth-service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type AuthRepository interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(cfg *config.Config) AuthRepository {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return &authRepository{db: db}
}

func (r *authRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}
