package services

import (
	"context"
	"errors"

	"github.com/cfex/microservices-in-go/projects-service/internal/models/dtos"
	"github.com/cfex/microservices-in-go/projects-service/internal/models/repositories"
)

type ProjectsService interface {
	GetById(ctx context.Context, projectID string) (*dtos.ProjectResponse, error)
}

type service struct {
	repo *repositories.Project
}

func NewProjectRepository(repo *repositories.Project) ProjectsService {
	return &service{repo: repo}
}

func (s *service) GetById(ctx context.Context, projectID string) (*dtos.ProjectResponse, error) {
	resp := &dtos.ProjectResponse{}

	proj, err := s.repo.GetById(ctx, projectID)
	if err != nil {
		return nil, errors.New("Not found")
	}

	return resp.FromEntity(proj), nil
}

func (s *service) GetAll(ctx context.Context) ([]*dtos.ProjectResponse, error) {

	res, err := s.repo.GetProjects(ctx)
	if err != nil {
		return nil, errors.New("internal server error")
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
