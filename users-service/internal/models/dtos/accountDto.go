package dtos

import "github.com/cfex/microservices-in-go/users-service/internal/models"

type AccountResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type RegisterRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	Account      *AccountResponse `json:"account"`
}

func (cr *RegisterRequest) ToEntity(password string) *models.User {
	return &models.User{
		Username: cr.Username,
		Email:    cr.Email,
		Role:     "user",
		Password: password,
	}
}

func (ur *AccountResponse) FromEntity(e *models.User) *AccountResponse {
	return &AccountResponse{
		ID:        e.ID,
		Username:  e.Username,
		Email:     e.Email,
		Role:      e.Role,
		CreatedAt: e.CreatedAt,
	}
}
