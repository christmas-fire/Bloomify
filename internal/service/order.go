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
	// logrus.Infof("Service: RemoveFlowerFromOrder started for userId: %d, flowerId: %d", userId, flowerId)
	err := s.repo.RemoveFlowerFromOrderByUser(userId, flowerId)
	if err != nil {
		// logrus.Errorf("Service: RemoveFlowerFromOrder - repository error: %v", err)
	}
	// logrus.Infof("Service: RemoveFlowerFromOrder finished")
	return err
}

// IncrementFlowerQuantity увеличивает количество цветка в заказе пользователя на 1
func (s *OrderService) IncrementFlowerQuantity(userId int, flowerId int) error {
	// logrus.Infof("Service: IncrementFlowerQuantity started for userId: %d, flowerId: %d", userId, flowerId)
	err := s.repo.IncrementFlowerQuantity(userId, flowerId)
	if err != nil {
		// logrus.Errorf("Service: IncrementFlowerQuantity - repository error: %v", err)
	}
	// logrus.Infof("Service: IncrementFlowerQuantity finished")
	return err
}

// DecrementFlowerQuantity уменьшает количество цветка в заказе пользователя на 1.
// Если количество становится 0, цветок удаляется из заказа.
// Может вернуть специальную ошибку ErrFlowerRemoved (нужно определить).
func (s *OrderService) DecrementFlowerQuantity(userId int, flowerId int) error {
	// logrus.Infof("Service: DecrementFlowerQuantity started for userId: %d, flowerId: %d", userId, flowerId)
	err := s.repo.DecrementFlowerQuantity(userId, flowerId)
	if err != nil {
		// logrus.Errorf("Service: DecrementFlowerQuantity - repository error: %v", err)
	}
	// logrus.Infof("Service: DecrementFlowerQuantity finished")
	return err
}

func (s *OrderService) DeleteActiveOrder(userId int) error {
	return s.repo.DeleteActiveOrderByUserId(userId)
}
