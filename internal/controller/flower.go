package controller

import (
	"net/http"
	"strconv"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateFlower godoc
// @Summary Create a new flower
// @Description Add a new flower to the database
// @Tags flowers
// @Accept json
// @Produce json
// @Param flower body models.Flower true "Flower data"
// @Success 201 {object} map[string]int "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/flowers [post]
func (h *Handler) createFlower(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input models.Flower
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := h.services.Flower.CreateFlower(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// GetAllFlowers godoc
// @Summary Get all flowers
// @Description Retrieve a list of all flowers
// @Tags flowers
// @Accept json
// @Produce json
// @Success 200 {array} models.Flower "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/flowers [get]
func (h *Handler) getAllFlowers(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	flowers, err := h.services.Flower.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, flowers)
}

// GetFlowerById godoc
// @Summary Get flower by ID
// @Description Retrieve a flower by its ID
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Success 200 {object} models.Flower "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/flowers/{id} [get]
func (h *Handler) getFlowerById(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := h.services.Flower.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteFlower godoc
// @Summary Delete a flower by ID
// @Description Delete a flower by its ID
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/flowers/{id} [delete]
func (h *Handler) deleteFlower(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.Flower.Delete(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// // Получение цветов по названию
// func (h *FlowerHandler) GetFlowersByName() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		name := r.URL.Query().Get("name")
// 		if name == "" {
// 			http.Error(w, "name parameter is required", http.StatusBadRequest)
// 			return
// 		}

// 		flowers, err := h.repo.GetFlowersByName(models.Flower{Name: name})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(flowers)
// 	}
// }

// // Получение цветов по цене
// func (h *FlowerHandler) GetFlowersByPrice() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		priceStr := r.URL.Query().Get("price")
// 		if priceStr == "" {
// 			http.Error(w, "price parameter is required", http.StatusBadRequest)
// 			return
// 		}

// 		price, err := strconv.ParseFloat(priceStr, 64)
// 		if err != nil {
// 			http.Error(w, "invalid price value", http.StatusBadRequest)
// 			return
// 		}

// 		flowers, err := h.repo.GetFlowersByPrice(models.Flower{Price: price})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(flowers)
// 	}
// }

// // Получение цветов по количеству в наличии
// func (h *FlowerHandler) GetFlowersByStock() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		stockStr := r.URL.Query().Get("stock")
// 		if stockStr == "" {
// 			http.Error(w, "stock parameter is required", http.StatusBadRequest)
// 			return
// 		}

// 		stock, err := strconv.Atoi(stockStr)
// 		if err != nil {
// 			http.Error(w, "invalid stock value", http.StatusBadRequest)
// 			return
// 		}

// 		flowers, err := h.repo.GetFlowersByStock(models.Flower{Stock: stock})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(flowers)
// 	}
// }
