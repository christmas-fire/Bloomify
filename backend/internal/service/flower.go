package service

import (
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type FlowerService struct {
	repo   repository.Flower
	logger *slog.Logger
}

// Параметры для фильтрации списка цветов
type FlowerFilter struct {
	Name        string   // Поиск по имени
	Description string   // Поиск по описанию
	MaxPrice    *float64 // Максимальная цена
	MaxStock    *int     // Максимальное количество
}

func NewFlowerService(repo repository.Flower, logger *slog.Logger) *FlowerService {
	return &FlowerService{repo: repo, logger: logger}
}

func (s *FlowerService) CreateFlower(name, description string, price float64, stock int) (int, error) {
	return s.repo.CreateFlower(name, description, price, stock)
}

func (s *FlowerService) Get(filter FlowerFilter) ([]models.Flower, error) {
	repoFilter := repository.FlowerFilter{
		NameQuery:        filter.Name,
		DescriptionQuery: filter.Description,
		MaxPrice:         filter.MaxPrice,
		MaxStock:         filter.MaxStock,
	}
	return s.repo.Get(repoFilter)
}

func (s *FlowerService) GetById(flowerId int) (models.Flower, error) {
	return s.repo.GetById(flowerId)
}

func (s *FlowerService) Delete(flowerId int) error {
	return s.repo.Delete(flowerId)
}

func (s *FlowerService) UpdateName(flowerId int, newName string) error {
	return s.repo.UpdateName(flowerId, newName)
}

func (s *FlowerService) UpdateDescription(flowerId int, newDescription string) error {
	return s.repo.UpdateDescription(flowerId, newDescription)
}

func (s *FlowerService) UpdatePrice(flowerId int, newPrice float64) error {
	return s.repo.UpdatePrice(flowerId, newPrice)
}

func (s *FlowerService) UpdateStock(flowerId int, newStock int) error {
	return s.repo.UpdateStock(flowerId, newStock)
}
