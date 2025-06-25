package main

import (
	"github.com/christmas-fire/Bloomify/internal/app"
	"github.com/sirupsen/logrus"

	_ "github.com/christmas-fire/Bloomify/docs"
)

// @title Bloomify API
// @version 1.0
// @description API server for flower's shop 'Bloomify'

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
