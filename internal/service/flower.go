package service

import (
	"fmt"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type FlowerService struct {
	repo repository.Flower
}

func NewFlowerService(repo repository.Flower) *FlowerService {
	return &FlowerService{repo: repo}
}

func (s *FlowerService) CreateFlower(flower models.Flower) (int, error) {
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
