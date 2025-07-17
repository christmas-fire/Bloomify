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

// Поля для обновления данных цветка
type UpdateFlowerInput struct {
	Name        *string  // Название
	Description *string  // Описание
	Price       *float64 // Цена
	Stock       *int     // Кол-во в наличии
}

func NewFlowerService(repo repository.Flower, logger *slog.Logger) *FlowerService {
	return &FlowerService{repo: repo, logger: logger}
}

func (s *FlowerService) CreateFlower(ctx context.Context, name, description string, price float64, stock int) (int, error) {
	return s.repo.CreateFlower(ctx, name, description, price, stock)
}

func (s *FlowerService) Get(ctx context.Context, filter FlowerFilter) ([]models.Flower, error) {
	repoFilter := repository.FlowerFilter{
		Name:        filter.Name,
		Description: filter.Description,
		MaxPrice:    filter.MaxPrice,
		MaxStock:    filter.MaxStock,
	}
	return s.repo.Get(ctx, repoFilter)
}

func (s *FlowerService) GetById(ctx context.Context, flowerId int) (models.Flower, error) {
	return s.repo.GetById(ctx, flowerId)
}

func (s *FlowerService) Delete(ctx context.Context, flowerId int) error {
	return s.repo.Delete(ctx, flowerId)
}

func (s *FlowerService) Update(ctx context.Context, flowerId int, input UpdateFlowerInput) error {
	repoInput := repository.UpdateFlowerInput{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
	}
	return s.repo.Update(ctx, flowerId, repoInput)
}
