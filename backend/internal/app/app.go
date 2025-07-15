package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/christmas-fire/Bloomify/internal/controller"
	"github.com/christmas-fire/Bloomify/internal/database"
	"github.com/christmas-fire/Bloomify/internal/logger"
	"github.com/christmas-fire/Bloomify/internal/metrics"
	"github.com/christmas-fire/Bloomify/internal/repository"
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

// Структура приложения
type App struct {
	handler    *controller.Handler
	service    *service.Service
	repository *repository.Repository
	server     *http.Server
	db         *sqlx.DB
	logger     *slog.Logger
	metrics    *metrics.Metrics
}

// Создание нового приложения
func NewApp() (*App, error) {
	logger := logger.InitLogger()

	db, err := database.InitPostgres()
	if err != nil {
		return nil, err
	}

	if err := database.InitTables(db); err != nil {
		db.Close()
		return nil, err
	}

	validator := validator.New()
	metrics := metrics.NewMetrics()
	repository := repository.NewRepository(db, logger)
	service := service.NewService(repository, logger)
	handler := controller.NewHandler(service, validator, logger, metrics)

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
		logger:     logger,
		metrics:    metrics,
	}, nil
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info("Starting server", "port", os.Getenv("BACKEND_PORT"))
		a.logger.Info("Documentation on", "url", "http://localhost:"+os.Getenv("BACKEND_PORT")+"/swagger/index.html#/")

		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Failed to start server", "error", err)
		}
	}()

	<-quit
	a.logger.Info("Shutting down server...")

	// Даем 30 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Выключаем сервер
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", "error", err)
	}

	// Закрываем остальные ресурсы, например, соединение с БД
	if err := a.Close(); err != nil {
		a.logger.Error("Failed to close resources", "error", err)
	}

	a.logger.Info("Server exiting")
	return nil
}

func (a *App) Close() error {
	a.logger.Info("Closing database connection...")
	return a.db.Close()
}
