package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	apperrors "github.com/cfex/microservices-in-go/services/common/errors"
	"github.com/cfex/microservices-in-go/services/users-service/internal/models/dtos"
	repositories "github.com/cfex/microservices-in-go/services/users-service/internal/repositories"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*dtos.AccountResponse, error)
	GetAllUsers(ctx context.Context) ([]*dtos.AccountResponse, error)
}

type userSvc struct {
	repo *repositories.User
}

func NewUserService(repo *repositories.User) UserService {
	return &userSvc{repo: repo}
}

func (s *userSvc) GetByID(ctx context.Context, id string) (*dtos.AccountResponse, error) {
	usr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	accountDto := &dtos.AccountResponse{}
	return accountDto.FromEntity(usr), nil
}

func (s *userSvc) GetAllUsers(ctx context.Context) ([]*dtos.AccountResponse, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	resp := make([]*dtos.AccountResponse, 0, len(users))

	for _, r := range users {
		re := &dtos.AccountResponse{
			ID:        r.ID,
			Username:  r.Username,
			Email:     r.Email,
			Role:      r.Role,
			CreatedAt: r.CreatedAt,
		}
		resp = append(resp, re)
	}

	return resp, nil
}
