package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/christmas-fire/Bloomify/internal/repository"
	"github.com/dgrijalva/jwt-go"
)

var (
	salt           = os.Getenv("SALT")
	signingKey     = os.Getenv("SIGNING_KEY")
	accessTokenTTL = 24 * time.Hour
)

type customClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo   repository.Auth
	logger *slog.Logger
}

func NewAuthService(repo repository.Auth, logger *slog.Logger) *AuthService {
	return &AuthService{repo: repo, logger: logger}
}

func (s *AuthService) CreateUser(username, email, password string) (int, error) {
	passwordHash := generatePasswordHash(password)

	return s.repo.CreateUser(username, email, passwordHash)
}

func (s *AuthService) generateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		customClaims{
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			},
			UserId: userId,
		},
	)
	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return signedToken, nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	return s.generateToken(user.Id)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		s.logger.Error("Error parsing token", "error", err)
		return 0, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *customClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
