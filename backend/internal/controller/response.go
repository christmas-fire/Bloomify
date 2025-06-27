package controller

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"error"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, logger *slog.Logger, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
