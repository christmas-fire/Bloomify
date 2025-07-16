package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/christmas-fire/Bloomify/internal/apperror"
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

func (s *AuthService) CreateUser(ctx context.Context, username, email, password string) (int, error) {
	passwordHash := generatePasswordHash(password)

	return s.repo.CreateUser(ctx, username, email, passwordHash)
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
		return "", &apperror.TokenError{
			Err:     err,
			Message: "failed to sign token",
		}
	}
	return signedToken, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	return s.generateToken(user.Id)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &apperror.TokenError{
				Err:     nil,
				Message: "invalid signing method",
			}
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, &apperror.TokenError{
			Err:     err,
			Message: "error parse token",
		}
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return 0, &apperror.TokenError{
			Err:     nil,
			Message: "invalid token claims",
		}
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
