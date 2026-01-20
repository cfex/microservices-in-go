package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cfex/microservices-in-go/services/users-service/internal/jwt"
	"github.com/cfex/microservices-in-go/services/users-service/internal/models/dtos"
	repositories "github.com/cfex/microservices-in-go/services/users-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error)
	Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.AccountResponse, error)
	Me(ctx context.Context, userID string) (*dtos.AccountResponse, error)
}

type authSvc struct {
	repo *repositories.User
}

func NewAuthService(repo *repositories.User) AuthService {
	return &authSvc{repo: repo}
}

func (s *authSvc) Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.AccountResponse, error) {
	response := &dtos.AccountResponse{}
	existingEmail, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	existingUsername, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	if existingUsername != nil || existingEmail != nil {
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	usr := req.ToEntity(string(hashedPassword))
	err = s.repo.Create(ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return response.FromEntity(usr), nil
}

func (s *authSvc) Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	r := &dtos.AccountResponse{}

	foundUsr, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid credentials: %w", err)
		}

		return nil, fmt.Errorf("internal server error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUsr.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	payload := jwt.Payload{UserID: foundUsr.ID, Role: foundUsr.Role}
	accessToken, refreshToken, err := jwt.Jwt.GenerateTokens(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &dtos.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account:      r.FromEntity(foundUsr),
	}, nil
}

func (s *authSvc) Me(ctx context.Context, userID string) (*dtos.AccountResponse, error) {

	usr, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	r := dtos.AccountResponse{}
	return r.FromEntity(usr), nil
}
