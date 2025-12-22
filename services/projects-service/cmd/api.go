package api

import (
	"database/sql"

	"github.com/cfex/microservices-in-go/services/projects-service/cmd/config"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/clients"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/handlers"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/models/repositories"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type router struct {
	DB  *sql.DB
	cfg *config.Config
	uc  *clients.UserClient
}

func CreateRouter(db *sql.DB, cfg *config.Config, uc *clients.UserClient) *router {
	return &router{DB: db, cfg: cfg, uc: uc}
}

func (r *router) NewRouter() *gin.Engine {
	ro := gin.Default()
	ro.Use(cors.New(r.cfg.Server.Cors))

	prg := ro.Group("/api/projects")

	pr := repositories.NewProjectRepository(r.DB)
	svc := services.NewProjectService(r.uc, pr)
	hndl := handlers.NewHandled(svc, r.cfg)

	prg.GET("/", hndl.GetAllProjects).POST("/", hndl.CreateProject)

	return ro
}
