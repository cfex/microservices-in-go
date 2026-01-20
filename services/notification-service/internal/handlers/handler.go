package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/cfex/microservices-in-go/services/common"
	"github.com/cfex/microservices-in-go/services/notification-service/config"
	"github.com/cfex/microservices-in-go/services/notification-service/internal/services"
	"github.com/rs/zerolog"
)

type EmailHandler struct {
	srv services.EmailSenderService
	cfg *config.Config
	log *zerolog.Logger
}

func NewEmailHandler(srv services.EmailSenderService, cfg *config.Config, log *zerolog.Logger) *EmailHandler {
	return &EmailHandler{srv: srv, cfg: cfg, log: log}
}

func (h *EmailHandler) SendEmail(to, subject, body string) {
	h.srv.SendEmail(to, subject, body)
}

func (h *EmailHandler) HandleMessage(body []byte) error {
	var event *common.Event

	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	switch event.Type {
	case "email.welcome":
		var data common.WelcomeData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return fmt.Errorf("failed to unmarshal welcome data: %w", err)
		}

		if err := h.srv.SendEmail(data.To, data.Subject, data.Username); err != nil {
			return fmt.Errorf("send email failed: %w", err)
		}

		return nil
	default:
		h.log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}
