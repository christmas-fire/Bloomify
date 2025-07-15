package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DTO для регистрации
type SignUpRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=50"` // Имя пользователя
	Email    string `json:"email" binding:"required" validate:"required,email"`           // Email пользователя
	Password string `json:"password" binding:"required" validate:"required,min=8"`        // Пароль пользователя
}

// DTO для логина
type SignInRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=50"` // Имя пользователя
	Password string `json:"password" binding:"required" validate:"required,min=8"`        // Пароль пользователя
}

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param SignUpRequest body SignUpRequest true "User data"
// @Success 201 {object} map[string]int "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Auth.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// SignIn godoc
// @Summary Sign in an existing user
// @Description Authenticate an existing user and return access and refresh JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param SignInRequest body SignInRequest true "Sign in data"
// @Success 200 {object} map[string]int "accessToken"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var req SignInRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		newErrorResponse(c, h.logger, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Auth.GenerateToken(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, errors.New("invalid username or password")) {
			newErrorResponse(c, h.logger, http.StatusUnauthorized, err.Error())
		} else {
			newErrorResponse(c, h.logger, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"accessToken": token,
	})
}
