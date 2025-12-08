package dtos

import "github.com/cfex/microservices-in-go/users-service/internal/models"

type UserReponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type UserCreateReponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type FinddAllUserReponse struct {
	Users []*UserReponse `json:"users"`
}

type CreateUserRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (cr *CreateUserRequest) ToEntity(password string) *models.User {
	return &models.User{
		Username: cr.Username,
		Email:    cr.Email,
		Role:     "user",
		Password: password,
	}
}

func (cr *UserCreateReponse) ToEntity() *models.User {
	return &models.User{
		ID:       cr.ID,
		Username: cr.Username,
		Email:    cr.Email,
		Role:     cr.Role,
	}
}

func (ur *UserCreateReponse) FromEntity(e *models.User) *UserCreateReponse {
	return &UserCreateReponse{
		ID:       e.ID,
		Username: e.Username,
		Email:    e.Email,
		Role:     e.Role,
	}
}
