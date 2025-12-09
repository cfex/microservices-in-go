package api

import (
	"database/sql"

	"github.com/cfex/microservices-in-go/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/users-service/internal/handlers"
	repositories "github.com/cfex/microservices-in-go/users-service/internal/repositorties"
	"github.com/cfex/microservices-in-go/users-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type router struct {
	DB  *sql.DB
	cfg *config.ServerConfig
}

func CreateRouter(db *sql.DB, cfg *config.Config) *router {
	return &router{DB: db, cfg: &cfg.Server}
}

func (r *router) NewRouter() *gin.Engine {
	ro := gin.Default()
	ro.Use(cors.New(r.cfg.Cors))

	ur := repositories.NewUserRepository(r.DB)
	us := services.NewService(ur)
	h := handlers.NewHandler(us)

	prg := ro.Group("/api")
	rg := ro.Group("/api/users")

	RegisterPublicRoutes(prg)
	RegisterPublicUserRoutes(rg, h)

	return ro
}
