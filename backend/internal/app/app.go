package app

import (
	"net/http"
	"os"

	"github.com/christmas-fire/Bloomify/internal/controller"
	"github.com/christmas-fire/Bloomify/internal/database"
	"github.com/christmas-fire/Bloomify/internal/repository"
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/sirupsen/logrus"
)

// Структура приложения
type App struct {
	handler    *controller.Handler
	service    *service.Service
	repository *repository.Repository
	server     *http.Server
}

// Создание нового приложения
func NewApp() (*App, error) {
	db, err := database.InitPostgres()
	if err != nil {
		return nil, err
	}

	if err := database.InitTables(db); err != nil {
		return nil, err
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := controller.NewHandler(service)

	server := &http.Server{
		Addr:    ":" + os.Getenv("BACKEND_PORT"),
		Handler: handler.InitRoutes(),
	}

	return &App{
		handler:    handler,
		service:    service,
		repository: repository,
		server:     server,
	}, nil
}

// Запуск сервера
func (a *App) Run() error {
	logrus.Printf("Starting server on http://localhost:%s", os.Getenv("BACKEND_PORT"))
	logrus.Printf("Documentation on http://localhost:%s/swagger/index.html#/", os.Getenv("BACKEND_PORT"))
	return a.server.ListenAndServe()
}
