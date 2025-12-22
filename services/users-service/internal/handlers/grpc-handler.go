package handlers

import (
	"context"

	"github.com/cfex/microservices-in-go/services/common/genproto"
	"github.com/cfex/microservices-in-go/services/users-service/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	svc transport.GrpcService
	genproto.UnimplementedUserServiceServer
}

func NewGRPCHandler(grpc *grpc.Server, svc transport.GrpcService) {
	grpcHandler := &grpcHandler{svc: svc}
	genproto.RegisterUserServiceServer(grpc, grpcHandler)
}

func (h *grpcHandler) GetUserById(ctx context.Context, req *genproto.GetUserRequest) (*genproto.UserResponse, error) {
	res, err := h.svc.GetUserById(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &genproto.UserResponse{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
		Role:     res.Role,
	}, nil
}
