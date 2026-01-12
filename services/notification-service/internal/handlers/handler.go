package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/cfex/microservices-in-go/services/notification-service/config"
	"github.com/cfex/microservices-in-go/services/notification-service/internal/services"
)

type EmailHandler struct {
	srv services.EmailSenderService
	cfg *config.Config
}

func NewEmailHandler(srv services.EmailSenderService, cfg *config.Config) *EmailHandler {
	return &EmailHandler{srv: srv, cfg: cfg}
}

func (h *EmailHandler) SendEmail(to, subject string) {
	h.srv.SendEmail(to, subject, "")
}

func (h *EmailHandler) HandleMessage(body []byte) error {
	var event any

	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	fmt.Print(event)

	return nil
}
