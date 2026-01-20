package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cfex/microservices-in-go/services/users-service/internal/env"
	"github.com/cfex/microservices-in-go/services/users-service/internal/logger"
	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type JwtConfig struct {
	SecretKey        string        `validate:"required,min=32"`
	TTL              time.Duration `validate:"required"`
	RefreshSecretKey string        `validate:"required,min=32"`
	RefreshTTL       time.Duration `validate:"required"`
}

type ServerConfig struct {
	Port     string      `validate:"required,port"`
	Env      string      `validate:"required,oneof=dev staging prod"`
	Cors     cors.Config `validate:"required"`
	Domain   string      `validate:"required,hostname"`
	GRPCPort string      `validate:"required,port"`
}

type DatabaseConfig struct {
	Host            string        `validate:"required,hostname"`
	Port            int           `validate:"required,port"`
	User            string        `validate:"required"`
	Password        string        `validate:"required"`
	Name            string        `validate:"required,alphanum"`
	DSN             string        `validate:"required"`
	MaxOpenConns    int           `validate:"required,numeric,min=1"`
	MaxIdleConns    int           `validate:"required,numeric,min=1"`
	ConnMaxLifetime time.Duration `validate:"required"`
	SSLMode         bool          `validate:"required"`
}

type GithubConfig struct {
	ClientID     string `validate:"required,alphanum"`
	ClientSecret string `validate:"required,alphanum"`
	RedirectURL  string `validate:"required,url"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Jwt      JwtConfig
	Github   GithubConfig
}

func Load() (*Config, error) {
	log := logger.GetLogger()

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("cannot load .env file: %w", err)
	}

	e := env.NewConfigLoader()

	jwtTTL, err := time.ParseDuration(e.GetEnv("JWT_TTL"))
	if err != nil {
		jwtTTL = time.Minute
		log.Warn().Err(err).Msg("Invalid JWT_TTL")
	}

	p, err := strconv.Atoi(e.GetEnv("DB_PORT"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid DB_PORT")
	}

	config := &Config{
		Server: ServerConfig{
			Port: e.GetEnv("PORT"),
			Env:  e.GetEnv("ENV"),
			Cors: cors.Config{
				AllowOrigins:     []string{e.GetEnv("FE_ORIGIN"), e.GetEnv("GITHUB_ORIGIN")},
				AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
			Domain:   e.GetEnv("DOMAIN"),
			GRPCPort: e.GetEnv("GRPC_PORT"),
		},
		Jwt: JwtConfig{SecretKey: e.GetEnv("SECRET_KEY"), TTL: jwtTTL},
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
		Github: GithubConfig{
			ClientID:     e.GetEnv("GITHUB_CLIENT_ID"),
			ClientSecret: e.GetEnv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  e.GetEnv("OAUTH_REDIRECT_URL"),
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
