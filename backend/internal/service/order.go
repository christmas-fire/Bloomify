package service

import (
	"log/slog"
	"strconv"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type OrderService struct {
	repo   repository.Order
	logger *slog.Logger
}

func NewOrderService(repo repository.Order, logger *slog.Logger) *OrderService {
	return &OrderService{repo: repo, logger: logger}
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

func (s *OrderService) GetAllOrderFlowers() ([]models.OrderFlowers, error) {
	return s.repo.GetAllOrderFlowers()
}

func (s *OrderService) GetOrderFlowersByOrderId(orderFlowersId int) ([]models.OrderFlowers, error) {
	return s.repo.GetOrderFlowersByOrderId(orderFlowersId)
}

func (s *OrderService) UpdateOrder(orderId int, input models.UpdateOrderInput) error {
	return s.repo.UpdateOrder(orderId, input)
}

func (s *OrderService) UpdateOrderFlowerId(orderId int, input models.UpdateOrderFlowerIdInput) error {
	return s.repo.UpdateOrderFlowerId(orderId, input)
}

func (s *OrderService) UpdateOrderQuantity(orderId int, input models.UpdateOrderQuantityInput) error {
	return s.repo.UpdateOrderQuantity(orderId, input)
}

func (s *OrderService) Delete(orderId int) error {
	return s.repo.Delete(orderId)
}

// RemoveFlowerFromOrder удаляет цветок из заказа пользователя
func (s *OrderService) RemoveFlowerFromOrder(userId int, flowerId int) error {
	s.logger.Info("Service: RemoveFlowerFromOrder started", "userId", userId, "flowerId", flowerId)
	err := s.repo.RemoveFlowerFromOrderByUser(userId, flowerId)
	if err != nil {
		s.logger.Error("Service: RemoveFlowerFromOrder - repository error", "error", err)
	}
	s.logger.Info("Service: RemoveFlowerFromOrder finished")
	return err
}

// IncrementFlowerQuantity увеличивает количество цветка в заказе пользователя на 1
func (s *OrderService) IncrementFlowerQuantity(userId int, flowerId int) error {
	s.logger.Info("Service: IncrementFlowerQuantity started", "userId", userId, "flowerId", flowerId)
	err := s.repo.IncrementFlowerQuantity(userId, flowerId)
	if err != nil {
		s.logger.Error("Service: IncrementFlowerQuantity - repository error", "error", err)
	}
	s.logger.Info("Service: IncrementFlowerQuantity finished")
	return err
}

// DecrementFlowerQuantity уменьшает количество цветка в заказе пользователя на 1.
// Если количество становится 0, цветок удаляется из заказа.
// Может вернуть специальную ошибку ErrFlowerRemoved (нужно определить).
func (s *OrderService) DecrementFlowerQuantity(userId int, flowerId int) error {
	s.logger.Info("Service: DecrementFlowerQuantity started", "userId", userId, "flowerId", flowerId)
	err := s.repo.DecrementFlowerQuantity(userId, flowerId)
	if err != nil {
		s.logger.Error("Service: DecrementFlowerQuantity - repository error", "error", err)
	}
	s.logger.Info("Service: DecrementFlowerQuantity finished")
	return err
}

func (s *OrderService) DeleteActiveOrder(userId int) error {
	return s.repo.DeleteActiveOrderByUserId(userId)
}
