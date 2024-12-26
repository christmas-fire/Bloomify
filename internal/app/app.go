package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/christmas-fire/Bloomify/internal/api"
	"github.com/christmas-fire/Bloomify/internal/database"
	"github.com/christmas-fire/Bloomify/internal/repository/postgres.go"
	"github.com/redis/go-redis/v9"
)

type App struct {
	db            *sql.DB
	client        *redis.Client
	userRepo      postgres.UserRepository
	userHandler   *api.UserHandler
	flowerRepo    postgres.FlowerRepository
	flowerHandler *api.FlowerHandler
	server        *http.Server
}

func NewApp() (*App, error) {
	db := database.InitPostgres()
	if err := database.InitTables(db); err != nil {
		return nil, err
	}

	client := database.InitRedis(context.Background())

	userRepo := postgres.NewUserRepository(db)
	flowerRepo := postgres.NewFlowerRepository(db)

	userHandler := api.NewUserHandler(*userRepo)
	flowerHandler := api.NewFlowerHandler(*flowerRepo)

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

func initRoutes(userHandler *api.UserHandler, flowerHandler *api.FlowerHandler, redisClient *redis.Client) http.Handler {
	r := http.DefaultServeMux

	r.Handle("/users/register", userHandler.SignUp())             // Регистрация
	r.HandleFunc("/users/login", userHandler.SignIn(redisClient)) // Логин
	r.HandleFunc("/users", userHandler.GetAllUsers())             // Получение всех пользователей
	r.HandleFunc("/users/delete/{id}", userHandler.DeleteUser())  // Удаление пользователя

	r.HandleFunc("/flowers/add", flowerHandler.AddFlower())      // Добавление нового цветка
	r.HandleFunc("/flowers", flowerHandler.GetAllFlowers())      // Получение всех цветов
	r.HandleFunc("/flowers/{id}", flowerHandler.GetFlowerByID()) // Получение цветка по ID

	return r
}

func (a *App) Run() error {
	log.Println("Starting server on :8080...")
	return a.server.ListenAndServe()
}

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
