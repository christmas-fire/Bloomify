package service

import (
	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Auth interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type User interface {
	GetAll() ([]models.User, error)
	GetById(userId int) (models.User, error)
	Delete(userId int) error
	UpdateUsername(userId int, input models.UpdateUsernameInput) error
	UpdatePassword(userId int, input models.UpdatePasswordInput) error
}

type Flower interface {
	CreateFlower(flower models.Flower) (int, error)
	GetAll() ([]models.Flower, error)
	GetById(flowerId int) (models.Flower, error)
	Delete(flowerId int) error
	GetFlowersByName(name string) ([]models.Flower, error)
	// Update(flowerId int, input models.Flower) error
}

type Service struct {
	Auth
	User
	Flower
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth:   NewAuthService(repos.Auth),
		User:   NewUserService(repos.User),
		Flower: NewFlowerService(repos.Flower),
	}
}
