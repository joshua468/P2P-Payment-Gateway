package service

import (
	"errors"
	"time"

	"github.com/joshua468/p2p-payment-gateway/auth-service/config"
	"github.com/joshua468/p2p-payment-gateway/auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req LoginRequest) (string, error)
	Signup(req SignupRequest) error
}

type authService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (s *authService) Login(req LoginRequest) (string, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := generateJWT(user.Username, s.cfg.JWTSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Signup(req SignupRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &repository.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	return s.repo.CreateUser(user)
}

func generateJWT(username, secret string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
