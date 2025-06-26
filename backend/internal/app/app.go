package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/christmas-fire/Bloomify/internal/controller"
	"github.com/christmas-fire/Bloomify/internal/database"
	"github.com/christmas-fire/Bloomify/internal/repository"
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Структура приложения
type App struct {
	handler    *controller.Handler
	service    *service.Service
	repository *repository.Repository
	server     *http.Server
	db         *sqlx.DB
}

// Создание нового приложения
func NewApp() (*App, error) {
	db, err := database.InitPostgres()
	if err != nil {
		return nil, err
	}

	if err := database.InitTables(db); err != nil {
		db.Close()
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
		db:         db,
	}, nil
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logrus.Printf("Starting server on http://localhost:%s", os.Getenv("BACKEND_PORT"))
		logrus.Printf("Documentation on http://localhost:%s/swagger/index.html#/", os.Getenv("BACKEND_PORT"))

		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	logrus.Println("Shutting down server...")

	// Даем 30 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Выключаем сервер
	if err := a.server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	// Закрываем остальные ресурсы, например, соединение с БД
	if err := a.Close(); err != nil {
		logrus.Errorf("Failed to close resources: %v", err)
	}

	logrus.Println("Server exiting")
	return nil
}

func (a *App) Close() error {
	logrus.Println("Closing database connection...")
	return a.db.Close()
}
