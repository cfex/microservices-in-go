package services

import (
	"fmt"
	"net/smtp"

	"github.com/cfex/microservices-in-go/services/notification-service/config"
)

type EmailSenderService struct {
	cfg *config.Config
}

func NewSender(cfg *config.Config) *EmailSenderService {
	return &EmailSenderService{cfg: cfg}
}

func (s *EmailSenderService) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.cfg.SmtpConfig.Host, s.cfg.SmtpConfig.Port)

	message := []byte(fmt.Sprintf(
		"From: <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		s.cfg.SmtpConfig.Sender,
		to,
		subject,
		body,
	))

	err := smtp.SendMail(addr, nil, s.cfg.SmtpConfig.Sender, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
