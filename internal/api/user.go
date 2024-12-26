package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/christmas-fire/Bloomify/internal/auth"
	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository/postgres.go"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	repo postgres.UserRepository
}

func NewUserHandler(repo postgres.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)

		log.Printf("user '%s' has signed-up", u.Username)
	}
}

func (h *UserHandler) SignIn(clint *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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

		if err := auth.AddTokenToRedis(context.Background(), clint, u.Username, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(token)

		log.Printf("user: '%s' has signed-in", u.Username)

	}
}

func (h *UserHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		users, err := h.repo.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}

func (h *UserHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/users/delete/")
		if path == "" {
			http.Error(w, "user ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		if err := h.repo.DeleteUserByID(id); err != nil {
			http.Error(w, "error delete user: %w", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Printf("user with id '%d' has deleted", id)
	}
}
