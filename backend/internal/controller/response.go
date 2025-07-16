package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/christmas-fire/Bloomify/internal/apperror"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"error"`
}

func newErrorResponse(c *gin.Context, logger *slog.Logger, err error) {
	var statusCode int
	var message string

	var alreadyExistsErr *apperror.AlreadyExistsError
	var noChangesErr *apperror.NoChangesError
	var timeoutErr *apperror.TimeoutError
	var notFoundErr *apperror.NotFoundError
	var internalServerErr *apperror.InternalServerError
	var invalidCredentialsErr *apperror.InvalidCredentialsError
	var badRequestErr *apperror.BadRequestError
	var tokenError *apperror.TokenError

	switch {
	case errors.As(err, &alreadyExistsErr):
		statusCode = http.StatusBadRequest
		message = err.Error()

	case errors.As(err, &noChangesErr):
		statusCode = http.StatusBadRequest
		message = err.Error()

	case errors.As(err, &timeoutErr):
		statusCode = http.StatusGatewayTimeout
		message = err.Error()

	case errors.As(err, &notFoundErr):
		statusCode = http.StatusNotFound
		message = err.Error()

	case errors.As(err, &internalServerErr):
		statusCode = http.StatusInternalServerError
		message = err.Error()

	case errors.As(err, &badRequestErr):
		statusCode = http.StatusBadRequest
		message = err.Error()

	case errors.As(err, &invalidCredentialsErr):
		statusCode = http.StatusUnauthorized
		message = err.Error()

	case errors.As(err, &tokenError):
		statusCode = http.StatusUnauthorized
		message = err.Error()

	default:
		statusCode = http.StatusInternalServerError
		message = "internal server error"
	}

	logAttrs := []slog.Attr{
		slog.String("message", err.Error()),
	}

	if cause := errors.Unwrap(err); cause != nil {
		logAttrs = append(logAttrs, slog.String("error", cause.Error()))
	}

	logger.LogAttrs(c.Request.Context(), slog.LevelError, "request failed", logAttrs...)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
