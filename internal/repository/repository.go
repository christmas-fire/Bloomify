package repository

import (
	"github.com/christmas-fire/Bloomify/internal/models"
)

// Репозиторий для работы с пользователями
type UserRepository interface {
	Register(u models.User) error // Регистрация пользователя
	Login(u models.User) error // Вход пользователя
	GetAllUsers() ([]models.User, error) // Получение всех пользователей
	DeleteUserByID(id int) error // Удаление пользователя по ID
	AddJWT(u models.User, token string) error // Добавление JWT токена пользователю
}

// Репозиторий для работы с цветами
type FlowerRepository interface {
	AddFlower(f models.Flower) error // Добавление цветка
	GetAllFlowers() ([]models.Flower, error) // Получение всех цветов
	GetFlowerByID(id int) (*models.Flower, error) // Получение цветка по ID
	GetFlowersByName(f models.Flower) ([]models.Flower, error) // Получение цветов по названию
	GetFlowersByPrice(f models.Flower) ([]models.Flower, error) // Получение цветов по цене
	GetFlowersByStock(f models.Flower) ([]models.Flower, error) // Получение цветов по количеству в наличии
	DeleteFlowerByID(id int) error // Удаление цветка по ID
}
