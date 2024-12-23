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

	r := http.DefaultServeMux

	r.HandleFunc("/register", userHandler.SignUp())
	r.HandleFunc("/login", userHandler.SignIn())
	r.HandleFunc("/users", userHandler.GetAllTasks())
	r.HandleFunc("/delete", userHandler.DeleteUser())

	log.Fatal(http.ListenAndServe(":8080", r))
}
