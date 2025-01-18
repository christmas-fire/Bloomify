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
	UpdateUsername(userId int, input models.UpdateUsernameInput) error
	UpdatePassword(userId int, input models.UpdatePasswordInput) error
	Delete(userId int) error
}

type Flower interface {
	CreateFlower(flower models.Flower) (int, error)
	GetAll() ([]models.Flower, error)
	GetById(flowerId int) (models.Flower, error)
	GetFlowersByName(name string) ([]models.Flower, error)
	GetFlowersByDescription(description string) ([]models.Flower, error)
	GetFlowersByPrice(price string) ([]models.Flower, error)
	GetFlowersByStock(stock string) ([]models.Flower, error)
	UpdateName(flowerId int, input models.UpdateNameInput) error
	UpdateDescription(flowerId int, input models.UpdateDescriptionInput) error
	UpdatePrice(flowerId int, input models.UpdatePriceInput) error
	UpdateStock(flowerId int, input models.UpdateStockInput) error
	Delete(flowerId int) error
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
