package main

import (
	"log"

	"github.com/christmas-fire/Bloomify/internal/app"
)

func main() {
	// Инициализация приложения
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	// Запуск сервера
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
