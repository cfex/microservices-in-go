package handlers

import (
	"net/http"

	"github.com/cfex/microservices-in-go/services/projects-service/cmd/config"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/logger"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/models/dtos"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/services"
	"github.com/gin-gonic/gin"
)

type handler struct {
	svc services.Service
	cfg *config.Config
}

func NewHandled(svc services.Service, cfg *config.Config) *handler {
	return &handler{svc: svc, cfg: cfg}
}

func (h *handler) CreateProject(c *gin.Context) {
	log := logger.GetLogger()
	var req dtos.ProjectCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})

		return
	}

	resp, err := h.svc.CreateProject(c.Request.Context(), &req)

	if err != nil {
		log.Error().Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
