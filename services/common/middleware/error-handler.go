package middleware

import (
	"errors"
	"net/http"

	apperrors "github.com/cfex/microservices-in-go/services/common/errors"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			statusCode := http.StatusInternalServerError
			message := "Internal server error"

			switch {
			case errors.Is(err, apperrors.ErrNotFound),
				errors.Is(err, apperrors.ErrUserNotFound),
				errors.Is(err, apperrors.ErrProjectNotFound):
				statusCode = http.StatusNotFound
				message = err.Error()

			case errors.Is(err, apperrors.ErrUnauthorized),
				errors.Is(err, apperrors.ErrInvalidPassword):
				statusCode = http.StatusUnauthorized
				message = err.Error()

			case errors.Is(err, apperrors.ErrValidation):
				statusCode = http.StatusBadRequest
				message = err.Error()

			case errors.Is(err, apperrors.ErrConflict),
				errors.Is(err, apperrors.ErrUserExists),
				errors.Is(err, apperrors.ErrProjectExists):
				statusCode = http.StatusConflict
				message = err.Error()

			default:
				message = "Internal server error"
			}

			c.JSON(statusCode, ErrorResponse{Error: message})
		}
	}
}
