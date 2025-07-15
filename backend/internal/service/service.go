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
	CreateFlower(name, description string, price float64, stock int) (int, error)
	Get(filter FlowerFilter) ([]models.Flower, error)
	GetById(flowerId int) (models.Flower, error)
	UpdateName(flowerId int, newName string) error
	UpdateDescription(flowerId int, newDescription string) error
	UpdatePrice(flowerId int, newPrice float64) error
	UpdateStock(flowerId int, newStock int) error
	Delete(flowerId int) error
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
