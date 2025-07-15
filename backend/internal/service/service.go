package service

import (
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type Auth interface {
	CreateUser(username, email, password string) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type User interface {
	GetAll() ([]models.User, error)
	GetById(userId int) (models.User, error)
	UpdateUsername(userId int, oldUsername, newUsername string) error
	UpdatePassword(userId int, username, oldPassword, newPassword string) error
	Delete(userId int) error
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

type Order interface {
	CreateOrder(userId int, order_flowers models.OrderFlowers) (int, error)
	GetAll() ([]models.Order, error)
	GetById(orderId int) (models.Order, error)
	GetOrdersByUserId(userId string) ([]models.Order, error)
	GetAllOrderFlowers() ([]models.OrderFlowers, error)
	GetOrderFlowersByOrderId(orderFlowersId int) ([]models.OrderFlowers, error)
	UpdateOrder(orderId int, input models.UpdateOrderInput) error
	UpdateOrderFlowerId(orderId int, input models.UpdateOrderFlowerIdInput) error
	UpdateOrderQuantity(orderId int, input models.UpdateOrderQuantityInput) error
	Delete(orderId int) error
	RemoveFlowerFromOrder(userId int, flowerId int) error
	IncrementFlowerQuantity(userId int, flowerId int) error
	DecrementFlowerQuantity(userId int, flowerId int) error
	DeleteActiveOrder(userId int) error
}

type Service struct {
	Auth
	User
	Flower
	Order
}

func NewService(repos *repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		Auth:   NewAuthService(repos.Auth, logger),
		User:   NewUserService(repos.User, logger),
		Flower: NewFlowerService(repos.Flower, logger),
		Order:  NewOrderService(repos.Order, logger),
	}
}
