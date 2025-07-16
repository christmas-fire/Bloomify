package service

import (
	"context"
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

func (s *FlowerService) CreateFlower(ctx context.Context, name, description string, price float64, stock int) (int, error) {
	return s.repo.CreateFlower(ctx, name, description, price, stock)
}

func (s *FlowerService) Get(ctx context.Context, filter FlowerFilter) ([]models.Flower, error) {
	repoFilter := repository.FlowerFilter{
		NameQuery:        filter.Name,
		DescriptionQuery: filter.Description,
		MaxPrice:         filter.MaxPrice,
		MaxStock:         filter.MaxStock,
	}
	return s.repo.Get(ctx, repoFilter)
}

func (s *FlowerService) GetById(ctx context.Context, flowerId int) (models.Flower, error) {
	return s.repo.GetById(ctx, flowerId)
}

func (s *FlowerService) Delete(ctx context.Context, flowerId int) error {
	return s.repo.Delete(ctx, flowerId)
}

func (s *FlowerService) UpdateName(ctx context.Context, flowerId int, newName string) error {
	return s.repo.UpdateName(ctx, flowerId, newName)
}

func (s *FlowerService) UpdateDescription(ctx context.Context, flowerId int, newDescription string) error {
	return s.repo.UpdateDescription(ctx, flowerId, newDescription)
}

func (s *FlowerService) UpdatePrice(ctx context.Context, flowerId int, newPrice float64) error {
	return s.repo.UpdatePrice(ctx, flowerId, newPrice)
}

func (s *FlowerService) UpdateStock(ctx context.Context, flowerId int, newStock int) error {
	return s.repo.UpdateStock(ctx, flowerId, newStock)
}
