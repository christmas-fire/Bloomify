package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/christmas-fire/Bloomify/internal/auth"
	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository/postgres.go"
)

type UserHandler struct {
	repo postgres.UserRepository
}

func NewUserHandler(repo postgres.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "error decode json", http.StatusBadRequest)
			return
		}

		err := h.repo.Register(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)

		log.Printf("user '%s' has signed-up", u.Username)
	}
}

func (h *UserHandler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "error decode json", http.StatusBadRequest)
			return
		}

		if err := h.repo.Login(u); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateJWT(u.Username)
		if err != nil {
			http.Error(w, "error generate JWT", http.StatusInternalServerError)
			return
		}

		if err := h.repo.AddJWT(u, token); err != nil {
			http.Error(w, "error add JWT into database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(token)

		log.Printf("user: '%s' has signed-in", u.Username)

	}
}

func (h *UserHandler) GetAllTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.repo.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func (h *UserHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "error decode json", http.StatusBadRequest)
			return
		}

		if err := h.repo.DeleteUser(u); err != nil {
			http.Error(w, "error delete user: %w", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Printf("user: '%s' has deleted", u.Username)
	}
}
