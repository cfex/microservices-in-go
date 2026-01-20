package config

import (
	"strconv"
	"time"

	"github.com/cfex/microservices-in-go/services/notification-service/logger"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/env"
	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port        string      `validate:"required,port"`
	Env         string      `validate:"required,oneof=dev staging prod"`
	Cors        cors.Config `validate:"required"`
	Domain      string      `validate:"required,hostname"`
	UsrGrpcPort string      `validate:"required,port"`
}

type DatabaseConfig struct {
	Host            string        `validate:"required,hostname"`
	Port            int           `validate:"required,port"`
	User            string        `validate:"required"`
	Password        string        `validate:"required"`
	Name            string        `validate:"required"`
	DSN             string        `validate:"required"`
	MaxOpenConns    int           `validate:"required,numeric,min=1"`
	MaxIdleConns    int           `validate:"required,numeric,min=1"`
	ConnMaxLifetime time.Duration `validate:"required"`
	SSLMode         bool          `validate:"required"`
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

func Load() (*Config, error) {
	log := logger.GetLogger()

	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Err(err).Msg(".env not found, falling back to environment variables")
	}

	e := env.NewConfigLoader()

	p, err := strconv.Atoi(e.GetEnv("DB_PORT"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid DB_PORT")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:            e.GetEnv("DB_HOST"),
			Port:            p,
			User:            e.GetEnv("DB_USER"),
			Password:        e.GetEnv("DB_PASSWORD"),
			Name:            e.GetEnv("DB_NAME"),
			DSN:             e.GetEnv("DSN"),
			MaxOpenConns:    20,
			MaxIdleConns:    10,
			ConnMaxLifetime: 5 * time.Minute,
			SSLMode:         false,
		},
		Server: ServerConfig{
			Port: e.GetEnv("PORT"),
			Env:  e.GetEnv("ENV"),
			Cors: cors.Config{
				AllowOrigins:     []string{"*"}, // for now
				AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
			Domain:      e.GetEnv("DOMAIN"),
			UsrGrpcPort: e.GetEnv("USER_GRPC_ADDR"),
		},
	}

	if err := e.Error(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	v := validator.New()
	return v.Struct(c)
}
