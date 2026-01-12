package config

import "github.com/cfex/microservices-in-go/services/notification-service/internal/env"

type SmptConfig struct {
	Host   string
	Port   string
	Sender string
}

type ServerConfig struct {
	Port string
}

type Config struct {
	Server     *ServerConfig
	SmtpConfig *SmptConfig
}

func InitConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Port: env.GetEnv("PORT"),
		},
		SmtpConfig: &SmptConfig{
			Host:   env.GetEnv("SMTP_HOST"),
			Port:   env.GetEnv("SMTP_PORT"),
			Sender: env.GetEnv("SENDER"),
		},
	}
}
