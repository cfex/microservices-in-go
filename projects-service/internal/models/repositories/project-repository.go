package repositories

import (
	"context"
	"database/sql"

	"github.com/cfex/microservices-in-go/projects-service/internal/database"
	"github.com/cfex/microservices-in-go/projects-service/internal/models"
)

type Project struct {
	database.BaseSqlRepository[models.Project]
}

func NewProjectRepository(db *sql.DB) *Project {
	return &Project{BaseSqlRepository: database.BaseSqlRepository[models.Project]{DB: db}}
}

func mapRow(row *sql.Row, u *models.Project) error {
	return row.Scan(&u.ID, &u.UserId, &u.Title, &u.Descritpion, &u.Status, &u.CreatedAt, &u.UpdatedAt)
}

func mapRows(rows *sql.Rows, u *models.Project) error {
	return rows.Scan(&u.ID, &u.UserId, &u.Title, &u.Descritpion, &u.Status, &u.CreatedAt, &u.UpdatedAt)
}

func (r *Project) GetById(ctx context.Context, id string) (*models.Project, error) {
	query := `SELECT id, user_id, title, description, status, created_at, updated_at FROM projects WHERE id = $1`

	return r.SelectSingle(mapRow, query, id)
}

func (r *Project) GetProjects(ctx context.Context) ([]*models.Project, error) {
	query := `SELECT id, user_id, title, description, status, created_at, updated_at FROM projects LIMIT 1000`

	return r.SelectMultiple(mapRows, query)
}
