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

// GetFlowersByName godoc
// @Summary Search flowers by name
// @Description Retrieve a list of flowers by their name
// @Tags flowers
// @Accept json
// @Produce json
// @Param name query string true "Flower Name"
// @Success 200 {array} models.Flower "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/flowers/search [get]
func (h *Handler) getFlowersByName(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	name := c.Query("name")
	if name == "" {
		newErrorResponse(c, http.StatusBadRequest, "not found query param")
		return
	}

	flowers, err := h.services.Flower.GetFlowersByName(name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, flowers)
}

// GetFlowersByDescription godoc
// @Summary Search flowers by description
// @Description Retrieve a list of flowers by their description
// @Tags flowers
// @Accept json
// @Produce json
// @Param description query string true "Flower Description"
// @Success 200 {array} models.Flower "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/flowers/search/description [get]
func (h *Handler) getFlowersByDescription(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	description := c.Query("description")
	if description == "" {
		newErrorResponse(c, http.StatusBadRequest, "not found query param")
		return
	}

	flowers, err := h.services.Flower.GetFlowersByDescription(description)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, flowers)
}

// GetFlowersByPrice godoc
// @Summary Search flowers by price
// @Description Retrieve a list of flowers with prices less than or equal to the specified value
// @Tags flowers
// @Accept json
// @Produce json
// @Param price query string true "Max Price"
// @Success 200 {array} models.Flower "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/flowers/price [get]
func (h *Handler) getFlowersByPrice(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	price := c.Query("price")
	if price == "" {
		newErrorResponse(c, http.StatusBadRequest, "not found query param")
		return
	}

	flowers, err := h.services.Flower.GetFlowersByPrice(price)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, flowers)
}

// GetFlowersByStock godoc
// @Summary Search flowers by stock
// @Description Retrieve a list of flowers with stock levels less than or equal to the specified value
// @Tags flowers
// @Accept json
// @Produce json
// @Param stock query string true "Max Stock"
// @Success 200 {array} models.Flower "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/flowers/stock [get]
func (h *Handler) getFlowersByStock(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	stock := c.Query("stock")
	if stock == "" {
		newErrorResponse(c, http.StatusBadRequest, "not found query param")
		return
	}

	flowers, err := h.services.Flower.GetFlowersByStock(stock)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, flowers)
}

// UpdateFlowerName godoc
// @Summary Update flower's name
// @Description Update the name of a specific flower by its ID
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Param input body models.UpdateNameInput true "Update Flower Name Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/flowers/{id}/name [patch]
func (h *Handler) updateFlowerName(c *gin.Context) {
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

	var input models.UpdateNameInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Flower.UpdateName(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// UpdateFlowerDescription godoc
// @Summary Update flower's description
// @Description Update the description of a specific flower by its ID
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Param input body models.UpdateDescriptionInput true "Update Flower Description Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/flowers/{id}/description [patch]
func (h *Handler) updateFlowerDescription(c *gin.Context) {
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

	var input models.UpdateDescriptionInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Flower.UpdateDescription(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// UpdateFlowerPrice godoc
// @Summary Update flower's price
// @Description Update the price of a specific flower by its ID
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Param input body models.UpdatePriceInput true "Update Flower Price Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/flowers/{id}/price [patch]
func (h *Handler) updateFlowerPrice(c *gin.Context) {
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

	var input models.UpdatePriceInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Flower.UpdatePrice(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// UpdateFlowerStock godoc
// @Summary Update flower's stock
// @Description Update the stock level of a specific flower by its ID
// @Tags flowers
// @Accept json
// @Produce json
// @Param id path int true "Flower ID"
// @Param input body models.UpdateStockInput true "Update Flower Stock Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/flowers/{id}/stock [patch]
func (h *Handler) updateFlowerStock(c *gin.Context) {
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

	var input models.UpdateStockInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Flower.UpdateStock(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
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
