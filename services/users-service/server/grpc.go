package server

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/cfex/microservices-in-go/services/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/services/users-service/internal/logger"
	"github.com/cfex/microservices-in-go/services/users-service/internal/repositories"
	"github.com/cfex/microservices-in-go/services/users-service/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	db     *sql.DB
	cfg    *config.Config
	server *grpc.Server
}

func NewGRPCServer(db *sql.DB, cfg *config.Config) *GRPCServer {
	return &GRPCServer{db: db, cfg: cfg}
}

func (s *GRPCServer) Run() error {
	log := logger.GetLogger()

	lis, err := net.Listen("tcp", s.cfg.Server.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.Server.GRPCPort, err)
	}

	kasp := keepalive.ServerParameters{
		Time:    20 * time.Second,
		Timeout: 5 * time.Second,
	}

	kep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}

	s.server = grpc.NewServer(
		grpc.KeepaliveParams(kasp),
		grpc.KeepaliveEnforcementPolicy(kep),
		grpc.ConnectionTimeout(10*time.Second),
	)

	ur := repositories.NewUserRepository(s.db)
	svc := transport.NewUsrGrpcSvc(ur)
	transport.NewGRPCHandler(s.server, svc)

	reflection.Register(s.server)

	log.Info().Msgf("gRPC server listening on %s", s.cfg.Server.GRPCPort)

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve grpc: %w", err)
	}

	return nil
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		log := logger.GetLogger()
		log.Info().Msg("Stopping gRPC server...")
		s.server.GracefulStop()
	}
}
