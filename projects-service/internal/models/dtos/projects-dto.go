package dtos

import (
	"time"

	"github.com/cfex/microservices-in-go/projects-service/internal/models"
)

type ProjectResponse struct {
	ID          string `json:"id"`
	UserId      string `json:"userId"`
	Title       string `json:"title"`
	Descritpion string `json:"description"`
	RepoUrl     string `json:"repo_url"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ProjectCreateRequest struct {
	Title       string `json:"title" binding:"required,min=10,max=244"`
	Description string `json:"description" binding:"required,min=1,max=500"`
	RepoUrl     string `json:"repo_url" binding:"required"`
	Status      string `json:"status" bindng:"required"`
}

func (pr *ProjectResponse) ToEntity() *models.Project {
	return &models.Project{
		ID:          pr.ID,
		UserId:      pr.UserId,
		Title:       pr.Title,
		Descritpion: pr.Descritpion,
		RepoUrl:     pr.RepoUrl,
		Status:      pr.Status,
		CreatedAt:   time.Now().String(),
	}
}

func (pr *ProjectResponse) FromEntity(p *models.Project) *ProjectResponse {
	return &ProjectResponse{
		ID:          p.ID,
		UserId:      p.UserId,
		Title:       p.Title,
		Descritpion: p.Descritpion,
		RepoUrl:     p.RepoUrl,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
