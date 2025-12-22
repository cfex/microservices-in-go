package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cfex/microservices-in-go/projects-service/internal/env"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Port   string
	Env    string
	Cors   cors.Config
	Domain string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	SSLMode         bool
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("cannot load .env file: %w", err)
	}

	p, err := strconv.Atoi(env.GetEnv("DB_PORT"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid DB_PORT")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:            env.GetEnv("DB_HOST"),
			Port:            p,
			User:            env.GetEnv("DB_USER"),
			Password:        env.GetEnv("DB_PASSWORD"),
			Name:            env.GetEnv("DB_NAME"),
			MaxOpenConns:    20,
			MaxIdleConns:    10,
			ConnMaxLifetime: 5 * time.Minute,
			SSLMode:         false,
		},
		Server: ServerConfig{
			Port: env.GetEnv("PORT"),
			Env:  env.GetEnv("ENV"),
			Cors: cors.Config{
				AllowOrigins:     []string{"*"}, // for now
				AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
			Domain: env.GetEnv("DOMAIN"),
		},
	}, nil
}
