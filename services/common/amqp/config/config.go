package config

import (
	"fmt"
	"net/url"

	"github.com/cfex/microservices-in-go/services/common/internal/env"
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

func Load() *Config {
	confing := &Config{
		User:     env.GetEnv("AMQP_USER"),
		Password: env.GetEnv("AMQP_PASSWORD"),
		Host:     env.GetEnv("AMQP_HOST"),
		Port:     env.GetEnv("AMQP_PORT"),
		Vhost:    env.GetEnv("AMQP_VHOST"),
	}

	return confing
}
