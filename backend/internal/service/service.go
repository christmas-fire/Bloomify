package service

import (
	"context"
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type Auth interface {
	CreateUser(ctx context.Context, username, email, password string) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type User interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetById(ctx context.Context, userId int) (models.User, error)
	UpdateUsername(ctx context.Context, userId int, oldUsername, newUsername string) error
	UpdatePassword(ctx context.Context, userId int, username, oldPassword, newPassword string) error
	Delete(ctx context.Context, userId int) error
}

type Flower interface {
	CreateFlower(ctx context.Context, name, description string, price float64, stock int) (int, error)
	Get(ctx context.Context, filter FlowerFilter) ([]models.Flower, error)
	GetById(ctx context.Context, flowerId int) (models.Flower, error)
	Update(ctx context.Context, flowerId int, input UpdateFlowerInput) error
	Delete(ctx context.Context, flowerId int) error
}

type Service struct {
	Auth
	User
	Flower
}

func NewService(repos *repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		Auth:   NewAuthService(repos.Auth, logger),
		User:   NewUserService(repos.User, logger),
		Flower: NewFlowerService(repos.Flower, logger),
	}
}
