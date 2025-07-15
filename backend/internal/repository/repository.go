package repository

import (
	"context"
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(ctx context.Context, username, email, password string) (int, error)
	GetUser(ctx context.Context, username, password string) (models.User, error)
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

type Repository struct {
	Auth
	User
	Flower
}

func NewRepository(db *sqlx.DB, logger *slog.Logger) *Repository {
	return &Repository{
		Auth:   NewAuthPostgres(db, logger),
		User:   NewUserPostgres(db, logger),
		Flower: NewFlowerPostgres(db, logger),
	}
}
