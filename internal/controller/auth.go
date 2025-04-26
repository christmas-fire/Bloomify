package controller

import (
	"errors"
	"net/http"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/gin-gonic/gin"
)

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} map[string]int "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Auth.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignIn godoc
// @Summary Sign in an existing user
// @Description Authenticate an existing user and return access and refresh JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param signInInput body controller.signInInput true "Sign in data"
// @Success 200 {object} service.Tokens "OK"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.Auth.GenerateToken(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, errors.New("invalid username or password")) {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Структура для запроса на обновление токена
type refreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// @Summary Refresh access and refresh tokens
// @Description Generate a new pair of access and refresh tokens using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshInput body controller.refreshInput true "Refresh token data"
// @Success 200 {object} service.Tokens "OK"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var input refreshInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid refresh token provided")
		return
	}

	tokens, err := h.services.Auth.RefreshToken(input.RefreshToken)
	if err != nil {
		if errors.Is(err, errors.New("refresh session not found")) || errors.Is(err, errors.New("refresh token expired")) {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}
