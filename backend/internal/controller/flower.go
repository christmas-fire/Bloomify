package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/christmas-fire/Bloomify/internal/apperror"
	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/gin-gonic/gin"
)

// DTO для отправки клиенту данных цветка
type FlowerResponse struct {
	Id          int     `json:"id"`          // ID цветка
	Name        string  `json:"name"`        // Название
	Description string  `json:"description"` // Описание
	Price       float64 `json:"price"`       // Цена
	Stock       int     `json:"stock"`       // Кол-во в наличии
}

// DTO для создания цветка
type CreateFlowerRequest struct {
	Name        string  `json:"name" binding:"required" validate:"required,min=3,max=50"` // Название
	Description string  `json:"description" validate:"max=1024"`                          // Описание (опционально)
	Price       float64 `json:"price" binding:"required" validate:"required,gt=0"`        // Цена
	Stock       int     `json:"stock" binding:"required" validate:"required,gte=0"`       // Кол-во в наличии
}

// DTO для обновления данных цветка
type UpdateFlowerRequest struct {
	Name        *string  `json:"name" validate:"omitempty,min=3,max=50"`    // Новое название
	Description *string  `json:"description" validate:"omitempty,max=1024"` // Новое описание
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`           // Новая цена
	Stock       *int     `json:"stock" validate:"omitempty,gte=0"`          // Новое кол-во в наличии
}

func toFlowerResponse(flower models.Flower) FlowerResponse {
	return FlowerResponse{
		Id:          flower.Id,
		Name:        flower.Name,
		Description: flower.Description,
		Price:       flower.Price,
		Stock:       flower.Stock,
	}
}

// CreateFlower godoc
// @Summary Create a new flower
// @Description Add a new flower to the database
// @Tags flowers
// @Accept json
// @Produce json
// @Param flower body CreateFlowerRequest true "Flower data"
// @Success 201 {object} map[string]int "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/flowers [post]
func (h *Handler) createFlower(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	var req CreateFlowerRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
		return
	}

	parentCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(parentCtx, defaultTimeout)
	defer cancel()

	id, err := h.services.Flower.CreateFlower(ctx, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// getFlowers godoc
// @Summary      Get a list of flowers with optional filters
// @Description  Retrieve a list of all flowers. Can be filtered by `name`, `description`, `max_price`, and `max_stock`. All filters are combined with AND logic.
// @Tags         flowers
// @Accept       json
// @Produce      json
// @Param        name        query string  false "Filter by flower name (case-insensitive, partial match)"
// @Param        description query string  false "Filter by description (case-insensitive, partial match)"
// @Param        max_price   query number  false "Filter by maximum price (inclusive)"
// @Param        max_stock   query integer false "Filter by maximum stock (inclusive)"
// @Success      200 {array} FlowerResponse "OK"
// @Failure      400 {object} map[string]string "Bad Request (e.g., invalid number format for price/stock)"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /api/v1/flowers [get]
func (h *Handler) getFlowers(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	var filter service.FlowerFilter

	filter.Name = c.Query("name")
	filter.Description = c.Query("description")

	if priceStr := c.Query("max_price"); priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
			return
		}
		filter.MaxPrice = &price
	}

	if stockStr := c.Query("max_stock"); stockStr != "" {
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
			return
		}
		filter.MaxStock = &stock
	}

	parentCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(parentCtx, defaultTimeout)
	defer cancel()

	flowers, err := h.services.Flower.Get(ctx, filter)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	var res []FlowerResponse
	for _, flower := range flowers {
		res = append(res, toFlowerResponse(flower))
	}

	if res == nil {
		res = []FlowerResponse{}
	}

	c.JSON(http.StatusOK, res)
}

// GetFlowerById godoc
// @Summary Get flower by ID
// @Description Retrieve a flower by its ID
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Success 200 {object} FlowerResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/flowers/{id} [get]
func (h *Handler) getFlowerById(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
		return
	}

	parentCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(parentCtx, defaultTimeout)
	defer cancel()

	flower, err := h.services.Flower.GetById(ctx, id)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, toFlowerResponse(flower))
}

// updateFlower godoc
// @Summary Partially update a flower
// @Description Update one or more fields of a flower by its ID. Only include the fields you want to change.
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Param input body UpdateFlowerRequest true "Fields to update"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/flowers/{id} [patch]
func (h *Handler) updateFlower(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
		return
	}

	var req UpdateFlowerRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
		return
	}

	// Проверяем, что хотя бы одно поле было передано для обновления
	if req.Name == nil && req.Description == nil && req.Price == nil && req.Stock == nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: fmt.Errorf("update body is empty")})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
		return
	}

	parentCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(parentCtx, defaultTimeout)
	defer cancel()

	input := service.UpdateFlowerInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := h.services.Flower.Update(ctx, id, input); err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	c.Status(http.StatusNoContent)
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
// @Security BearerAuth
// @Router /api/flowers/{id} [delete]
func (h *Handler) deleteFlower(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, h.logger, &apperror.BadRequestError{Err: err})
		return
	}

	parentCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(parentCtx, defaultTimeout)
	defer cancel()

	if err := h.services.Flower.Delete(ctx, id); err != nil {
		newErrorResponse(c, h.logger, err)
		return
	}

	c.Status(http.StatusNoContent)
}
