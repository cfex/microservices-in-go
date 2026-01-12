package internal

type WelcomeEventEmail struct {
	Username string `json:"username" binding:"required"`
}

type ProjectCreatedEventEmail struct {
	ProjectName  string `json:"project_name" binding:"required"`
	ProjectOwner string `json:"project_owner" binding:"required"`
	ProjectID    string `json:"project_id"`
}
