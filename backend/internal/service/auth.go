package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
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
	repo repository.Repository
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	if err := validateUser(user); err != nil {
		return 0, err
	}
	user.Password = generatePasswordHash(user.Password)

	return s.repo.Auth.CreateUser(user)
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
	user, err := s.repo.Auth.GetUser(username, generatePasswordHash(password))
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
		logrus.Errorf("Error parsing token: %v", err)
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

func validateUser(u models.User) error {
	if len(u.Username) < 3 {
		return fmt.Errorf("username must have at least 3 characters")
	}

	if len(u.Password) < 8 {
		return fmt.Errorf("password must have at least 8 characters")
	}

	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
