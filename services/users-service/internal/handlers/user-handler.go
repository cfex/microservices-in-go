package handlers

import (
	"net/http"

	"github.com/cfex/microservices-in-go/services/users-service/internal/logger"
	"github.com/cfex/microservices-in-go/services/users-service/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	log := logger.GetLogger()
	id := c.Param("id")

	usr, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Error().Msg(err.Message)
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, usr)
}
