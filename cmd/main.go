package main

import (
	"github.com/christmas-fire/Bloomify/internal/app"
	"github.com/sirupsen/logrus"
)

// @title Bloomify API
// @version 3.0
// @description API server for flower's shop 'Bloomify'

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Инициализация приложения
	app, err := app.NewApp()
	if err != nil {
		logrus.Fatal(err)
	}

	// Запуск сервера
	if err := app.Run(); err != nil {
		logrus.Fatal(err)
	}
}
