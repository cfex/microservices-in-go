package handlers

import (
	"net/http"

	"github.com/cfex/microservices-in-go/users-service/internal/logger"
	"github.com/cfex/microservices-in-go/users-service/internal/models/dtos"
	"github.com/cfex/microservices-in-go/users-service/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(c *gin.Context) {
	log := logger.GetLogger()
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})

		return
	}

	usr, err := h.service.CreateUser(c.Request.Context(), &req)

	if err != nil {
		log.Error().Msg(err.Message)
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(201, usr)
}
