package config

import (
	"fmt"

	"github.com/cfex/microservices-in-go/services/notification-service/internal/env"
	"github.com/cfex/microservices-in-go/services/notification-service/logger"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type SmptConfig struct {
	Host   string `validate:"required,hostname"`
	Port   string `validate:"required,port"`
	Sender string `validate:"required,email"`
}

type ServerConfig struct {
	Port string `validate:"required,port"`
}

type Config struct {
	Server     *ServerConfig
	SmtpConfig *SmptConfig
}

func Load() (*Config, error) {
	log := logger.GetLogger()

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("cannot load .env file: %w", err)
	}

	e := env.NewConfigLoader()

	log.Info().Msg("Configuration loaded from environment variables")

	cfg := &Config{
		Server: &ServerConfig{
			Port: e.GetEnv("PORT"),
		},
		SmtpConfig: &SmptConfig{
			Host:   e.GetEnv("SMTP_HOST"),
			Port:   e.GetEnv("SMTP_PORT"),
			Sender: e.GetEnv("SENDER"),
		},
	}

	if err := e.Error(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	v := validator.New()
	return v.Struct(c)
}
