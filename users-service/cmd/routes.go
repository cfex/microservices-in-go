package api

import (
	"net/http"

	"github.com/cfex/microservices-in-go/users-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Users service is healthy"})
	})
}

func RegisterPublicUserRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	rg.GET("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get user by ID endpoint"})
	})
	rg.POST("/", h.CreateUser)
}
