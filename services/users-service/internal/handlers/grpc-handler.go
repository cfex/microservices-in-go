package handlers

import (
	"context"

	"github.com/cfex/microservices-in-go/services/common/genproto"
	"github.com/cfex/microservices-in-go/services/users-service/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	usrSvc  services.UserService
	authSvc services.AuthService
}

func NewGRPCHandler(usrSvc services.UserService, authSvc services.AuthService) *grpcHandler {
	return &grpcHandler{usrSvc: usrSvc, authSvc: authSvc}
}

func (h *grpcHandler) GetUserById(ctx context.Context, req *genproto.GetUserRequest) (*genproto.UserResponse, error) {
	res, err := h.usrSvc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Message)
	}

	return &genproto.UserResponse{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
		Role:     res.Role,
	}, nil
}

func (h *grpcHandler) AuthenticateRequest(ctx context.Context, req *genproto.TokenRequest) (*genproto.TokenResponse, error) {

	return nil, nil
}
