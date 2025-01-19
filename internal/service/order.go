package service

import (
	"strconv"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(userId int, order_flowers models.OrderFlowers) (int, error) {
	return s.repo.CreateOrder(userId, order_flowers)
}

func (s *OrderService) GetAll() ([]models.Order, error) {
	return s.repo.GetAll()
}

func (s *OrderService) GetById(orderId int) (models.Order, error) {
	return s.repo.GetById(orderId)
}

func (s *OrderService) GetOrdersByUserId(userId string) ([]models.Order, error) {
	intUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, err
	}

	return s.repo.GetOrdersByUserId(intUserId)
}
