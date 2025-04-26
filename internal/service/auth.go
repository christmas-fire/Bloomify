package service

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const (
	salt              = "romanbarma2005"
	signingKey        = "sonyboy"
	accessTokenTTL    = 15 * time.Minute
	refreshTokenTTL   = 720 * time.Hour
	refreshTokenBytes = 32
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

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) generateTokens(userId int) (Tokens, error) {
	var tokens Tokens
	var err error

	accesToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		customClaims{
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			},
			UserId: userId,
		},
	)
	tokens.AccessToken, err = accesToken.SignedString([]byte(signingKey))
	if err != nil {
		return Tokens{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenBytesArr := make([]byte, refreshTokenBytes)
	_, err = rand.Read(refreshTokenBytesArr)
	if err != nil {
		return Tokens{}, fmt.Errorf("failed to generate refresh token bytes: %w", err)
	}
	tokens.RefreshToken = base64.URLEncoding.EncodeToString(refreshTokenBytesArr)

	refreshTokenHash := sha256.Sum256([]byte(tokens.RefreshToken))

	session := models.RefreshSession{
		UserID:           userId,
		RefreshTokenHash: fmt.Sprintf("%x", refreshTokenHash[:]),
		ExpiresAt:        time.Now().Add(refreshTokenTTL),
	}

	if err := s.repo.Session.CreateSession(session); err != nil {
		return Tokens{}, fmt.Errorf("failed to create refresh session: %w", err)
	}

	return tokens, nil
}

func (s *AuthService) GenerateToken(username, password string) (Tokens, error) {
	user, err := s.repo.Auth.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return Tokens{}, errors.New("invalid username or password")
	}

	return s.generateTokens(user.Id)
}

func (s *AuthService) RefreshToken(refreshToken string) (Tokens, error) {
	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	hashString := fmt.Sprintf("%x", refreshTokenHash[:])

	session, err := s.repo.Session.GetSession(hashString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Tokens{}, errors.New("refresh session not found")
		}
		return Tokens{}, fmt.Errorf("failed to get session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		if delErr := s.repo.Session.DeleteSession(hashString); delErr != nil {
			logrus.Errorf("Failed to delete expired session: %v", delErr)
		}
		return Tokens{}, errors.New("refresh token expired")
	}

	if err := s.repo.Session.DeleteSession(hashString); err != nil {
		return Tokens{}, fmt.Errorf("failed to delete old session: %w", err)
	}

	return s.generateTokens(session.UserID)
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
