package services

import (
	"context"

	"github.com/cfex/microservices-in-go/services/users-service/internal/models"
	"github.com/cfex/microservices-in-go/services/users-service/internal/models/dtos"
	repositories "github.com/cfex/microservices-in-go/services/users-service/internal/repositories"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*dtos.AccountResponse, *models.ErrorResponse)
}

type userSvc struct {
	repo *repositories.User
}

func NewUserService(repo *repositories.User) UserService {
	return &userSvc{repo: repo}
}

func (s *userSvc) GetByID(ctx context.Context, id string) (*dtos.AccountResponse, *models.ErrorResponse) {
	usr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve user",
			Err:     err,
		}
	}

	accountDto := &dtos.AccountResponse{}
	return accountDto.FromEntity(usr), nil
}
