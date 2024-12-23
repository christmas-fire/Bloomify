package repository

import (
	"github.com/christmas-fire/Bloomify/internal/models"
)

type UserRepository interface {
	Register(u models.User) error
	Login(u models.User) error
	GetAllUsers() ([]models.User, error)
	DeleteUser(u models.User) error
	AddJWT(u models.User, token string) error
}

type FlowerRepository interface {
	AddFlower(f models.Flower) error
	GetAllFlowers() ([]models.Flower, error)
	GetFlowerByID(f models.Flower) (models.Flower, error)
	GetFlowersByName(f models.Flower) ([]models.Flower, error)
	GetFlowersByPrice(f models.Flower) ([]models.Flower, error)
	GetFlowersByStock(f models.Flower) ([]models.Flower, error)
}
