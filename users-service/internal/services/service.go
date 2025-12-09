package services

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/cfex/microservices-in-go/users-service/internal/models"
	"github.com/cfex/microservices-in-go/users-service/internal/models/dtos"
	repositories "github.com/cfex/microservices-in-go/users-service/internal/repositorties"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, req *dtos.CreateUserRequest) (*dtos.UserCreateResponse, *models.ErrorResponse)
}

type svc struct {
	repo *repositories.User
}

func NewService(repo *repositories.User) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateUser(ctx context.Context, req *dtos.CreateUserRequest) (*dtos.UserCreateResponse, *models.ErrorResponse) {
	response := &dtos.UserCreateResponse{}
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if existingUser != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User already exists",
			Err:     err,
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Err:     err,
		}
	}

	usr := req.ToEntity(string(hashedPassword))
	if usr == nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "User conversion failed",
		}
	}

	err = s.repo.Create(ctx, usr)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Err:     err,
		}
	}

	return response.FromEntity(usr), nil
}
