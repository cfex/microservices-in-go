package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cfex/microservices-in-go/users-service/internal/env"
	"github.com/cfex/microservices-in-go/users-service/internal/logger"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

type JwtConfig struct {
	SecretKey        string
	TTL              time.Duration
	RefreshSecretKey string
	RefreshTTL       time.Duration
}

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

type GithubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
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

	jwtTTL, err := time.ParseDuration(env.GetEnv("JWT_TTL"))
	if err != nil {
		jwtTTL = time.Minute
		log.Warn().Err(err).Msg("Invalid JWT_TTL")
	}

	p, err := strconv.Atoi(env.GetEnv("DB_PORT"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid DB_PORT")
	}

	return &Config{
		Server: ServerConfig{
			Port: env.GetEnv("PORT"),
			Env:  env.GetEnv("ENV"),
			Cors: cors.Config{
				AllowOrigins:     []string{env.GetEnv("FE_ORIGIN"), env.GetEnv("GITHUB_ORIGIN"), "*"}, // For testing purposes
				AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
			Domain: env.GetEnv("DOMAIN"),
		},
		Jwt: JwtConfig{SecretKey: env.GetEnv("SECRET_KEY"), TTL: jwtTTL},
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
		Github: GithubConfig{
			ClientID:     env.GetEnv("GITHUB_CLIENT_ID"),
			ClientSecret: env.GetEnv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  env.GetEnv("OAUTH_REDIRECT_URL"),
		},
	}, nil
}
