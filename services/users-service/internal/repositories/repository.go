package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cfex/microservices-in-go/services/users-service/internal/database"
	"github.com/cfex/microservices-in-go/services/users-service/internal/models"
)

var ErrUserNotCreated = errors.New("user not created")
var ErrUserNotFound = errors.New("user not found")

type User struct {
	database.BaseSqlRepository[models.User]
}

func NewUserRepository(db *sql.DB) *User {
	return &User{
		BaseSqlRepository: database.BaseSqlRepository[models.User]{DB: db},
	}
}

func mapRowToUser(row *sql.Row, u *models.User) error {
	return row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
}

func (r *User) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, email, password, role, created_at FROM users WHERE username = $1`

	return r.SelectSingle(mapRowToUser, query, username)
}

func (r *User) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password, role, created_at FROM users WHERE email = $1`

	return r.SelectSingle(mapRowToUser, query, email)
}

func (r *User) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, username, email, password, role, created_at FROM users WHERE id = $1`

	return r.SelectSingle(mapRowToUser, query, id)
}

func (r *User) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, email, password, role, created_at) 
        VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) 
    `
	id, err := r.Insert(query, user.Username, user.Email, user.Password, user.Role)
	user.ID = id

	return err
}
