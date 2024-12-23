package main

import (
	"log"
	"net/http"

	"github.com/christmas-fire/Bloomify/internal/api"
	"github.com/christmas-fire/Bloomify/internal/database"
	"github.com/christmas-fire/Bloomify/internal/repository/postgres.go"
)

func main() {
	db := database.InitDB()
	if err := database.InitTables(db); err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)
	userHandler := api.NewUserHandler(*userRepo)

	flowerRepo := postgres.NewFlowerRepository(db)
	flowerHandler := api.NewFlowerHandler(*flowerRepo)

	r := http.DefaultServeMux

	r.Handle("/users/register", userHandler.SignUp())       // Регистрация
	r.HandleFunc("/users/login", userHandler.SignIn())      // Логин
	r.HandleFunc("/users", userHandler.GetAllUsers())       // Получение всех пользователей
	r.HandleFunc("/users/delete", userHandler.DeleteUser()) // Удаление пользователя

	r.HandleFunc("/flowers/add", flowerHandler.AddFlower())      // Добавление нового цветка
	r.HandleFunc("/flowers", flowerHandler.GetAllFlowers())      // Получение всех цветов
	r.HandleFunc("/flowers/{id}", flowerHandler.GetFlowerByID()) // Получение цветка по ID

	log.Fatal(http.ListenAndServe(":8080", r))
}
