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
	Password string `json:"password" binding:"required,min=8,max=72"`
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" binding:"required,email,max=254"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
		Password: password,
		Role:     "user",
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
