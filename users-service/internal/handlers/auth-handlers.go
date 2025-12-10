package handlers

import (
	"fmt"
	"net/http"

	"github.com/cfex/microservices-in-go/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/users-service/internal/logger"
	"github.com/cfex/microservices-in-go/users-service/internal/models/dtos"
	"github.com/cfex/microservices-in-go/users-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.AuthService
	cfg     *config.ServerConfig
}

func NewAuthHandler(service services.AuthService, cfg *config.ServerConfig) *AuthHandler {
	return &AuthHandler{service: service, cfg: cfg}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	log := logger.GetLogger()
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})

		return
	}

	usr, err := h.service.Register(c.Request.Context(), &req)

	if err != nil {
		log.Error().Msg(err.Message)
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(201, usr)
}

func (h *AuthHandler) Login(c *gin.Context) {
	log := logger.GetLogger()
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	res, err := h.service.Login(c, &req)

	if err != nil {
		log.Error().Msg(err.Message)
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.SetCookie("access_token", res.AccessToken, 3600, "/", h.cfg.Domain, false, true)
	c.SetCookie("refresh_token", res.RefreshToken, 604800, "/", h.cfg.Domain, false, true)
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", h.cfg.Domain, false, true)
	c.SetCookie("refresh_token", "", -1, "/", h.cfg.Domain, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	log := logger.GetLogger()
	fmt.Print(c)
	userIdValue, exists := c.Get("user_id")
	if !exists {
		log.Error().Msg("No user id ")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No user id found"})
		return
	}

	userId, ok := userIdValue.(string)
	if !ok {
		log.Error().Msg("Invalid user id type")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
		return
	}

	res, err := h.service.Me(c, userId)
	if err != nil {
		log.Error().Msg("Unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, res)
}
