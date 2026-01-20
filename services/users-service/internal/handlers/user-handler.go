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
		log.Error().Err(err).Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	log := logger.GetLogger()

	rep, err := h.service.GetAllUsers(c)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rep)
}
