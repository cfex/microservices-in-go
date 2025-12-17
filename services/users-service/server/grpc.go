package server

import (
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer()
	log.Print("Starting grpc server")

	return gRPCServer.Serve(lis)
}
