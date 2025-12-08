package models

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}
