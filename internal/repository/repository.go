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
	Delete(userId int) error
	UpdateUsername(userId int, input models.UpdateUsernameInput) error
	UpdatePassword(userId int, input models.UpdatePasswordInput) error
}

type Flower interface {
	CreateFlower(flower models.Flower) (int, error) // Добавление цветка
	GetAll() ([]models.Flower, error)
	GetById(flowerId int) (models.Flower, error)
	Delete(flowerId int) error
	// GetAll() ([]models.Flower, error)                  // Получение всех цветов
	// GetFlowerByID(flowerId int) (models.Flower, error) // Получение цветка по ID
	// DeleteFlowerByID(flowerId int) error               // Удаление цветка по ID
	// GetFlowersByName(name string) ([]models.Flower, error)  // Получение цветов по названию
	// GetFlowersByPrice(price float64) ([]models.Flower, error) // Получение цветов по цене
	// GetFlowersByStock(stock int) ([]models.Flower, error) // Получение цветов по количеству в наличии
}

type Repository struct {
	Auth
	User
	Flower
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:   NewAuthPostgres(db),
		User:   NewUserPostgres(db),
		Flower: NewFlowerPostgres(db),
	}
}
