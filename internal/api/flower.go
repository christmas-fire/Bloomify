package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository/postgres.go"
)

type FlowerHandler struct {
	repo postgres.FlowerRepository
}

func NewFlowerHandler(repo postgres.FlowerRepository) *FlowerHandler {
	return &FlowerHandler{repo: repo}
}

func (h *FlowerHandler) AddFlower() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var f models.Flower

		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			http.Error(w, "error decode json", http.StatusBadRequest)
			return
		}

		if err := h.repo.AddFlower(f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(f)

		log.Printf("flower '%s' has added", f.Name)
	}
}

func (h *FlowerHandler) GetAllFlowers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		flowers, err := h.repo.GetAllFlowers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flowers)
	}
}

func (h *FlowerHandler) GetFlowerByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/flowers/")
		if path == "" {
			http.Error(w, "flower ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid flower ID", http.StatusBadRequest)
			return
		}

		flower, err := h.repo.GetFlowerByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flower)
	}
}
