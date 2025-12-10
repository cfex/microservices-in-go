package api

import (
	"net/http"

	"github.com/cfex/microservices-in-go/users-service/internal/handlers"
	"github.com/cfex/microservices-in-go/users-service/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Users service is healthy"})
	})
}
func RegisterPublicUserRoutes(rg *gin.RouterGroup) {
	rg.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Public Users Endpoint"})
	})
}

func RegisterPublicAuthRoutes(rg *gin.RouterGroup, h *handlers.AuthHandler) {
	rg.POST("/login", h.Login)
	rg.POST("/register", h.RegisterUser)
	rg.POST("/logout", h.Logout)
}

func RegisterProtectedAuthRoutes(rg *gin.RouterGroup, h *handlers.AuthHandler) {
	rg.Use(middlewares.RequireAuthentication)
	rg.POST("/me", h.GetMe)
}

func RegisterProtectedUserRoutes(rg *gin.RouterGroup, h *handlers.UserHandler) {
	rg.Use(middlewares.RequireAuthentication)
	rg.GET("/:id", h.GetUserByID)
}
