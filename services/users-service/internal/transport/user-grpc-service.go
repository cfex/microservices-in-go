package transport

import (
	"context"

	"github.com/cfex/microservices-in-go/services/common/genproto"
	"github.com/cfex/microservices-in-go/services/users-service/internal/repositories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcService interface {
	GetUserById(context.Context, *genproto.GetUserRequest) (*genproto.UserResponse, error)
}

type grpcSvc struct {
	repo *repositories.User
	genproto.UnimplementedUserServiceServer
}

func NewUsrGrpcSvc(repo *repositories.User) GrpcService {
	return &grpcSvc{repo: repo}
}

func (s *grpcSvc) GetUserById(ctx context.Context, req *genproto.GetUserRequest) (*genproto.UserResponse, error) {
	dto, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &genproto.UserResponse{
		ID:       dto.ID,
		Username: dto.Username,
		Email:    dto.Email,
		Role:     dto.Role,
	}, nil
}
