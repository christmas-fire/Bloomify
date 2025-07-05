package controller

import (
	"net/http"
	"strconv"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/gin-gonic/gin"
)

// DTO для отправки клиенту данных пользователя
type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// DTO для смены имени пользователя
type UpdateUsernameRequest struct {
	OldUsername string `json:"oldUsername" binding:"required" validate:"required,min=3,max=50"` // Старое имя пользователя
	NewUsername string `json:"newUsername" binding:"required" validate:"required,min=3,max=50"` // Новое имя пользователя
}

// DTO для смены пароля
type UpdatePasswordRequest struct {
	Username    string `json:"username" binding:"required" validate:"required,min=3,max=50"` // Имя пользователя
	OldPassword string `json:"oldPassword" binding:"required" validate:"required,min=8"`     // Старый пароль пользователя
	NewPassword string `json:"newPassword" binding:"required" validate:"required,min=8"`     // Новый пароль пользователя
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} UserResponse "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/users [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusUnauthorized, err.Error())
		return
	}

	users, err := h.services.User.GetAll()
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		return
	}

	var res []UserResponse
	for _, user := range users {
		res = append(res, toUserResponse(user))
	}

	c.JSON(http.StatusOK, res)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *Handler) getUserById(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := h.services.User.GetById(id)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

// UpdateUserUsername godoc
// @Summary Update user's username
// @Description Update the username of a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body UpdateUsernameRequest true "Update Username Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/users/{id}/username [patch]
func (h *Handler) updateUserUsername(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, "invalid id param")
		return
	}

	var req UpdateUsernameRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.UpdateUsername(id, req.OldUsername, req.NewUsername); err != nil {
		newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateUserPassword godoc
// @Summary Update user's password
// @Description Update the password of a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body UpdatePasswordRequest true "Update Password Input"
// @Success 200 {object} statusResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/users/{id}/password [patch]
func (h *Handler) updateUserPassword(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, "invalid id param")
		return
	}

	var req UpdatePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.UpdatePassword(id, req.Username, req.OldPassword, req.NewPassword); err != nil {
		newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	_, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.User.Delete(id); err != nil {
		newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetMe godoc
// @Summary Get current user info
// @Description Get info about the currently authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} UserResponse "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/users/me [get]
func (h *Handler) getMe(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.services.User.GetById(userId)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func toUserResponse(user models.User) UserResponse {
	return UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}
