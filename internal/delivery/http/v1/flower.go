package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository/postgres.go"
)

// Обработчик цветов
type FlowerHandler struct {
	repo postgres.FlowerRepository
}

// Создание нового обработчика цветов
func NewFlowerHandler(repo postgres.FlowerRepository) *FlowerHandler {
	return &FlowerHandler{repo: repo}
}

// Добавление цветка
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

// Получение всех цветов
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

// Получение цветка по ID
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

// Получение цветов по названию
func (h *FlowerHandler) GetFlowersByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "name parameter is required", http.StatusBadRequest)
			return
		}

		flowers, err := h.repo.GetFlowersByName(models.Flower{Name: name})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flowers)
	}
}

// Получение цветов по цене
func (h *FlowerHandler) GetFlowersByPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		priceStr := r.URL.Query().Get("price")
		if priceStr == "" {
			http.Error(w, "price parameter is required", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			http.Error(w, "invalid price value", http.StatusBadRequest)
			return
		}

		flowers, err := h.repo.GetFlowersByPrice(models.Flower{Price: price})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flowers)
	}
}

// Получение цветов по количеству в наличии
func (h *FlowerHandler) GetFlowersByStock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stockStr := r.URL.Query().Get("stock")
		if stockStr == "" {
			http.Error(w, "stock parameter is required", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			http.Error(w, "invalid stock value", http.StatusBadRequest)
			return
		}

		flowers, err := h.repo.GetFlowersByStock(models.Flower{Stock: stock})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flowers)
	}
}

// Удаление цветка по ID
func (h *FlowerHandler) DeleteFlowerByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем ID из URL path
		path := strings.TrimPrefix(r.URL.Path, "/flowers/delete/")
		if path == "" {
			http.Error(w, "flower ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid flower ID", http.StatusBadRequest)
			return
		}

		if err := h.repo.DeleteFlowerByID(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Printf("flower with ID '%d' has been deleted", id)
	}
}
