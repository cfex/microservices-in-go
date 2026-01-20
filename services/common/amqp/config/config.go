package config

import (
	"fmt"
	"net/url"

	"github.com/cfex/microservices-in-go/services/common/internal/env"
	"github.com/cfex/microservices-in-go/services/common/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	User     string
	Host     string
	Password string
	Port     string
	Vhost    string
}

func (c *Config) AmqpConnectinoString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, url.PathEscape(c.Vhost))
}

func Load() (*Config, error) {
	log := logger.GetLogger()

	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Err(err).Msg(".env not found, falling back to environment variables")
	}

	e := env.NewConfigLoader()

	config := &Config{
		User:     e.GetEnv("AMQP_USER"),
		Password: e.GetEnv("AMQP_PASSWORD"),
		Host:     e.GetEnv("AMQP_HOST"),
		Port:     e.GetEnv("AMQP_PORT"),
		Vhost:    e.GetEnv("AMQP_VHOST"),
	}

	if err := e.Error(); err != nil {
		return nil, err
	}

	return config, nil
}
