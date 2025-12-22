package transport

import (
	"context"

	"github.com/cfex/microservices-in-go/services/common/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	svc GrpcService
	pb.UnimplementedUserServiceServer
}

func NewGRPCHandler(grpc *grpc.Server, svc GrpcService) {
	grpcHandler := &grpcHandler{svc: svc}
	pb.RegisterUserServiceServer(grpc, grpcHandler)
}

func (h *grpcHandler) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	res, err := h.svc.GetUserById(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
		Role:     res.Role,
	}, nil
}
