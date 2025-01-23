package repository

import (
	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
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
	GetFlowersByPrice(price float64) ([]models.Flower, error)
	GetFlowersByStock(stock int64) ([]models.Flower, error)
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
	GetOrdersByUserId(userId int64) ([]models.Order, error)
	GetAllOrderFlowers() ([]models.OrderFlowers, error)
	GetOrderFlowersByOrderId(orderFlowersId int) ([]models.OrderFlowers, error)
	UpdateOrder(orderId int, input models.UpdateOrderInput) error
	UpdateOrderFlowerId(orderId int, input models.UpdateOrderFlowerIdInput) error
	UpdateOrderQuantity(orderId int, input models.UpdateOrderQuantityInput) error
	Delete(orderId int) error
}

type Repository struct {
	Auth
	User
	Flower
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:   NewAuthPostgres(db),
		User:   NewUserPostgres(db),
		Flower: NewFlowerPostgres(db),
		Order:  NewOrderPostgres(db),
	}
}
