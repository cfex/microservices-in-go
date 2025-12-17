package models

type Project struct {
	ID          string `json:"id"`
	UserId      string `json:"userId"`
	Title       string `json:"title"`
	Descritpion string `json:"description"`
	RepoUrl     string `json:"repo_url"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
