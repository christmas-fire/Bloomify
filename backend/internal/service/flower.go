package service

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type FlowerService struct {
	repo   repository.Flower
	logger *slog.Logger
}

func NewFlowerService(repo repository.Flower, logger *slog.Logger) *FlowerService {
	return &FlowerService{repo: repo, logger: logger}
}

func (s *FlowerService) CreateFlower(flower models.Flower) (int, error) {
	if err := validateFlower(flower); err != nil {
		return 0, err
	}

	return s.repo.CreateFlower(flower)
}

func (s *FlowerService) GetAll() ([]models.Flower, error) {
	return s.repo.GetAll()
}

func (s *FlowerService) GetById(flowerId int) (models.Flower, error) {
	return s.repo.GetById(flowerId)
}

func (s *FlowerService) Delete(flowerId int) error {
	return s.repo.Delete(flowerId)
}

func (s *FlowerService) GetFlowersByName(name string) ([]models.Flower, error) {
	return s.repo.GetFlowersByName(name)
}

func (s *FlowerService) GetFlowersByDescription(description string) ([]models.Flower, error) {
	return s.repo.GetFlowersByDescription(description)
}

func (s *FlowerService) GetFlowersByPrice(price string) ([]models.Flower, error) {
	floatPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, err
	}

	return s.repo.GetFlowersByPrice(floatPrice)
}

func (s *FlowerService) GetFlowersByStock(stock string) ([]models.Flower, error) {
	intStock, err := strconv.ParseInt(stock, 10, 64)
	if err != nil {
		return nil, err
	}

	return s.repo.GetFlowersByStock(intStock)
}

func (s *FlowerService) UpdateName(flowerId int, input models.UpdateNameInput) error {
	if err := validateName(input.NewName); err != nil {
		return err
	}
	return s.repo.UpdateName(flowerId, input)
}

func (s *FlowerService) UpdateDescription(flowerId int, input models.UpdateDescriptionInput) error {
	return s.repo.UpdateDescription(flowerId, input)
}

func (s *FlowerService) UpdatePrice(flowerId int, input models.UpdatePriceInput) error {
	if err := validatePrice(input.NewPrice); err != nil {
		return err
	}
	return s.repo.UpdatePrice(flowerId, input)
}

func (s *FlowerService) UpdateStock(flowerId int, input models.UpdateStockInput) error {
	if err := validateStock(input.NewStock); err != nil {
		return err
	}
	return s.repo.UpdateStock(flowerId, input)
}

func validateFlower(flower models.Flower) error {
	if flower.Name == "" {
		return fmt.Errorf("name can't be empty")
	}
	if flower.Price <= 0 {
		return fmt.Errorf("price can't be less or equal 0")
	}
	if flower.Stock < 0 {
		return fmt.Errorf("stock can't be less than 0")
	}
	return nil
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("name can't be empty")
	}

	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("price can't be less or equal 0")
	}

	return nil
}

func validateStock(stock int) error {
	if stock < 0 {
		return fmt.Errorf("stock can't be less than 0")
	}

	return nil
}
