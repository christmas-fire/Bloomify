package controller

import (
	"net/http"
	"strconv"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Summary Create a new order with a flower
// @Description Create a new order for a user with a specific flower and quantity
// @Tags orders
// @Accept json
// @Produce json
// @Param input body models.OrderFlowers true "Order Flowers Input"
// @Success 201 {object} map[string]int "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders [post]
func (h *Handler) createOrder(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input models.OrderFlowers
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order_id, err := h.services.Order.CreateOrder(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": order_id,
	})
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Get all orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} models.Order "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders [get]
func (h *Handler) getAllOrders(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	orders, err := h.services.Order.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderById godoc
// @Summary Get order by ID
// @Description Get an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders/{id} [get]
func (h *Handler) getOrderById(c *gin.Context) {
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

	order, err := h.services.Order.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrdersByUserId godoc
// @Summary Get orders by user ID
// @Description Get all orders for a specific user
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} models.Order "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders/user_id [get]
func (h *Handler) getOrdersByUserId(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	user_id := c.Query("user_id")
	if user_id == "" {
		newErrorResponse(c, http.StatusBadRequest, "not found query param")
		return
	}

	orders, err := h.services.Order.GetOrdersByUserId(user_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrder godoc
// @Summary Update an order
// @Description Update the details of a specific order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param input body models.UpdateOrderInput true "Update Order Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders/{id} [put]
func (h *Handler) updateOrder(c *gin.Context) {
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

	var input models.UpdateOrderInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Order.UpdateOrder(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// UpdateOrderFlowerId godoc
// @Summary Update order's flower ID
// @Description Update the flower ID of a specific order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param input body models.UpdateOrderFlowerIdInput true "Update Order Flower ID Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders/{id}/flower_id [patch]
func (h *Handler) updateOrderFlowerId(c *gin.Context) {
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

	var input models.UpdateOrderFlowerIdInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Order.UpdateOrderFlowerId(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// UpdateOrderQuantity godoc
// @Summary Update order's flower quantity
// @Description Update the quantity of a specific flower in an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param input body models.UpdateOrderQuantityInput true "Update Order Quantity Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders/{id}/quantity [patch]
func (h *Handler) updateOrderQuantity(c *gin.Context) {
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

	var input models.UpdateOrderQuantityInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Order.UpdateOrderQuantity(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete a specific order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/orders/{id} [delete]
func (h *Handler) deleteOrder(c *gin.Context) {
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

	if err := h.services.Order.Delete(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllOrderFlowers godoc
// @Summary Get all order flowers
// @Description Get all flowers in orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} models.OrderFlowers "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/order_flowers [get]
func (h *Handler) getAllOrderFlowers(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	orderFlowers, err := h.services.Order.GetAllOrderFlowers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orderFlowers)
}

// GetOrderFlowersByOrderId godoc
// @Summary Get order flowers by order ID
// @Description Get all flowers in a specific order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {array} models.OrderFlowers "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/order_flowers/{id} [get]
func (h *Handler) getOrderFlowersByOrderId(c *gin.Context) {
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

	orderFlowers, err := h.services.Order.GetOrderFlowersByOrderId(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orderFlowers)
}
