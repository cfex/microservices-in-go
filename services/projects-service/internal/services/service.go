package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/cfex/microservices-in-go/services/projects-service/internal/clients"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/models/dtos"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/models/repositories"
)

var (
	UserNotFound      = errors.New("user not found")
	ProjectNotFound   = errors.New("project not found")
	InternalServerErr = errors.New("internal server error")
)

type Service interface {
	GetById(ctx context.Context, projectID string) (*dtos.ProjectResponse, error)
	GetAll(ctx context.Context) ([]*dtos.ProjectResponse, error)
	CreateProject(ctx context.Context, req *dtos.ProjectCreateRequest) (*dtos.ProjectResponse, error)
}

type service struct {
	client *clients.UserClient
	repo   *repositories.Project
}

func NewProjectService(client *clients.UserClient, repo *repositories.Project) Service {
	return &service{client: client, repo: repo}
}

func (s *service) GetById(ctx context.Context, projectID string) (*dtos.ProjectResponse, error) {
	resp := &dtos.ProjectResponse{}

	proj, err := s.repo.GetById(ctx, projectID)
	if err != nil {
		return nil, errors.New("Not found")
	}

	return resp.FromEntity(proj), nil
}

func (s *service) CreateProject(ctx context.Context, req *dtos.ProjectCreateRequest) (*dtos.ProjectResponse, error) {
	userId, err := s.client.GetUserById(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	if userId == "" {
		return nil, UserNotFound
	}

	prj := req.ToEntity()
	err = s.repo.CreateProject(ctx, prj)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	response := &dtos.ProjectResponse{}
	return response.FromEntity(prj), nil
}

func (s *service) GetAll(ctx context.Context) ([]*dtos.ProjectResponse, error) {
	res, err := s.repo.GetProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerErr, err)
	}

	resp := make([]*dtos.ProjectResponse, 0, len(res))

	for _, r := range res {
		proj := &dtos.ProjectResponse{
			ID:          r.ID,
			UserId:      r.UserId,
			Title:       r.Title,
			Descritpion: r.Descritpion,
			RepoUrl:     r.RepoUrl,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
		resp = append(resp, proj)
	}

	return resp, nil
}
