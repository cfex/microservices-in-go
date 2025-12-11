package services

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/cfex/microservices-in-go/users-service/internal/jwt"
	"github.com/cfex/microservices-in-go/users-service/internal/models"
	"github.com/cfex/microservices-in-go/users-service/internal/models/dtos"
	repositories "github.com/cfex/microservices-in-go/users-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, *models.ErrorResponse)
	Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.AccountResponse, *models.ErrorResponse)
	Me(ctx context.Context, userID string) (*dtos.AccountResponse, *models.ErrorResponse)
}

type authSvc struct {
	repo *repositories.User
}

func NewAuthService(repo *repositories.User) AuthService {
	return &authSvc{repo: repo}
}

func (s *authSvc) Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.AccountResponse, *models.ErrorResponse) {
	response := &dtos.AccountResponse{}
	existingEmail, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	existingUsername, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if existingUsername != nil || existingEmail != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusConflict,
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

func (s *authSvc) Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, *models.ErrorResponse) {
	r := &dtos.AccountResponse{}

	foundUsr, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
				Err:     err,
			}
		}

		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Err:     err,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUsr.Password), []byte(req.Password))
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
			Err:     err,
		}
	}

	payload := jwt.Payload{UserID: foundUsr.ID, Role: foundUsr.Role}
	accessToken, refreshToken, err := jwt.Jwt.GenerateTokens(payload)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Token generation failed",
			Err:     err,
		}
	}

	return &dtos.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account:      r.FromEntity(foundUsr),
	}, nil
}

func (s *authSvc) Me(ctx context.Context, userID string) (*dtos.AccountResponse, *models.ErrorResponse) {

	usr, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Err:     err,
		}
	}

	r := dtos.AccountResponse{}
	return r.FromEntity(usr), nil
}
