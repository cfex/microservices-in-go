package transport

import (
	"context"

	"github.com/cfex/microservices-in-go/services/common/genproto"
	"github.com/cfex/microservices-in-go/services/users-service/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcSvc struct {
	genproto.UnimplementedUserServiceServer
	svc services.UserService
}

func NewGrpcSvc(svc services.UserService) *grpcSvc {
	return &grpcSvc{svc: svc}
}

func (s *grpcSvc) GetUserById(ctx context.Context, req *genproto.GetUserRequest) (*genproto.UserResponse, error) {

	dto, err := s.svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Message)
	}

	return &genproto.UserResponse{
		ID:       dto.ID,
		Username: dto.Username,
		Email:    dto.Email,
		Role:     dto.Role,
	}, nil

}
