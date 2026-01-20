package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cfex/microservices-in-go/services/common"
	"github.com/cfex/microservices-in-go/services/common/amqp"
	amqpConsts "github.com/cfex/microservices-in-go/services/common/amqp/consts"
	"github.com/cfex/microservices-in-go/services/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/services/users-service/internal/logger"
	"github.com/cfex/microservices-in-go/services/users-service/internal/models/dtos"
	"github.com/cfex/microservices-in-go/services/users-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service  services.AuthService
	cfg      *config.Config
	producer *amqp.Producer
}

func NewAuthHandler(service services.AuthService, cfg *config.Config, producer *amqp.Producer) *AuthHandler {
	return &AuthHandler{service: service, cfg: cfg, producer: producer}
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
		log.Error().Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	welcomeData := common.WelcomeData{
		To:       usr.Email,
		Subject:  "Welcome to Platform!",
		Username: usr.Username,
	}

	marshaled, marshalErr := json.Marshal(welcomeData)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msg("failed to marshal welcome data")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	event := common.Event{
		Type: "email.welcome",
		Data: marshaled,
	}

	h.producer.Publish(ctx, amqpConsts.EmailExchange, event.Type, event)

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
		log.Error().Err(err).Msg(err.Error())
		c.Error(err)
		return
	}

	c.SetCookie("access_token", res.AccessToken, 3600, "/", h.cfg.Server.Domain, false, true)
	c.SetCookie("refresh_token", res.RefreshToken, 604800, "/", h.cfg.Server.Domain, false, true)
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", h.cfg.Server.Domain, false, true)
	c.SetCookie("refresh_token", "", -1, "/", h.cfg.Server.Domain, false, true)
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
		log.Error().Msg("Invalid user id")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
		return
	}

	res, err := h.service.Me(c, userId)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, res)
}
