package app

import (
	"fmt"
	"net/http"

	"github.com/christmas-fire/Bloomify/internal/controller"
	"github.com/christmas-fire/Bloomify/internal/repository"
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	if err := initConfig(); err != nil {
		return nil, fmt.Errorf("error read config: %s", err.Error())
	}

	db, err := repository.InitPostgreSQL(repository.PostgreSQLConfig{
		User:     viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		Database: viper.GetString("postgres.database"),
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
		Sslmode:  viper.GetString("postgres.sslmode"),
	})

	if err != nil {
		return nil, err
	}

	if err := repository.InitTables(db); err != nil {
		return nil, err
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := controller.NewHandler(service)

	server := &http.Server{
		Addr:    viper.GetString("port"),
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
	logrus.Println("Starting server on :8080...")
	return a.server.ListenAndServe()
}

// Чтение файла конфигурации
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
