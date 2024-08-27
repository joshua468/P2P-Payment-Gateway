package internal

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	JWTSecret string
	DB        *gorm.DB
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return err
	}

	user := User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user User
	if err := s.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	if !s.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := s.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
