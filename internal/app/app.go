package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/christmas-fire/Bloomify/internal/database"
	delivery "github.com/christmas-fire/Bloomify/internal/delivery/http/v1"
	"github.com/christmas-fire/Bloomify/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
)

// Структура приложения
type App struct {
	db            *sql.DB
	client        *redis.Client
	userRepo      postgres.UserRepository
	userHandler   *delivery.UserHandler
	flowerRepo    postgres.FlowerRepository
	flowerHandler *delivery.FlowerHandler
	server        *http.Server
}

// Создание нового приложения
func NewApp() (*App, error) {
	db := database.InitPostgres()
	if err := database.InitTables(db); err != nil {
		return nil, err
	}

	client := database.InitRedis(context.Background())

	userRepo := postgres.NewUserRepository(db)
	flowerRepo := postgres.NewFlowerRepository(db)

	userHandler := delivery.NewUserHandler(*userRepo)
	flowerHandler := delivery.NewFlowerHandler(*flowerRepo)

	server := &http.Server{
		Addr:    ":8080",
		Handler: initRoutes(userHandler, flowerHandler, client),
	}

	return &App{
		db:            db,
		client:        client,
		userRepo:      *userRepo,
		userHandler:   userHandler,
		flowerRepo:    *flowerRepo,
		flowerHandler: flowerHandler,
		server:        server,
	}, nil
}

// Инициализация маршрутов
func initRoutes(userHandler *delivery.UserHandler, flowerHandler *delivery.FlowerHandler, redisClient *redis.Client) http.Handler {
	r := http.DefaultServeMux

	r.Handle("/api/v1/users/sign-up", userHandler.SignUp())                // Регистрация
	r.HandleFunc("/api/v1/users/sign-in", userHandler.SignIn(redisClient)) // Логин
	r.HandleFunc("/api/v1/users", userHandler.GetAllUsers())               // Получение всех пользователей
	r.HandleFunc("/api/v1/users/delete/{id}", userHandler.DeleteUser())    // Удаление пользователя по ID

	r.HandleFunc("/api/v1/flowers/add", flowerHandler.AddFlower())                  // Добавление нового цветка
	r.HandleFunc("/api/v1/flowers", flowerHandler.GetAllFlowers())                  // Получение всех цветов
	r.HandleFunc("/api/v1/flowers/{id}", flowerHandler.GetFlowerByID())             // Получение цветка по ID
	r.HandleFunc("/api/v1/flowers/search/name", flowerHandler.GetFlowersByName())   // Поиск цветов по имени
	r.HandleFunc("/api/v1/flowers/search/price", flowerHandler.GetFlowersByPrice()) // Поиск цветов по цене
	r.HandleFunc("/api/v1/flowers/search/stock", flowerHandler.GetFlowersByStock()) // Поиск цветов по наличию
	r.HandleFunc("/api/v1/flowers/delete/{id}", flowerHandler.DeleteFlowerByID())   // Удаление цветка по ID

	return r
}

// Запуск сервера
func (a *App) Run() error {
	log.Println("Starting server on :8080...")
	return a.server.ListenAndServe()
}

// Завершение работы приложения
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.db.Close(); err != nil {
		return err
	}
	if err := a.client.Close(); err != nil {
		return err
	}

	return nil
}
