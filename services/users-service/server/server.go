package server

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/cfex/microservices-in-go/services/users-service/cmd/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type server struct {
	log    zerolog.Logger
	router *gin.Engine
	cfg    *config.Config
}

func NewServer(log zerolog.Logger, router *gin.Engine, config *config.Config) *server {
	return &server{
		log:    log,
		router: router,
		cfg:    config,
	}
}

func (s *server) Serve() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := http.Server{
		Addr:              fmt.Sprintf(":%s", s.cfg.Server.Port),
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       time.Minute,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	go func() {
		s.log.Info().Str("port", s.cfg.Server.Port).Msg("Server starting...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	<-ctx.Done()

	stop()
	log.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
