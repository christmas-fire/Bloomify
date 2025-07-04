package service

import (
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

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
	authService := NewAuthService(*repos, logger)

	return &Service{
		Auth:   authService,
		User:   NewUserService(repos.User, logger),
		Flower: NewFlowerService(repos.Flower, logger),
		Order:  NewOrderService(repos.Order, logger),
	}
}
