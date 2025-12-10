package api

import (
	"database/sql"

	"github.com/cfex/microservices-in-go/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/users-service/internal/handlers"
	repositories "github.com/cfex/microservices-in-go/users-service/internal/repositories"
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
	as := services.NewAuthService(ur)
	us := services.NewUserService(ur)
	ah := handlers.NewAuthHandler(as, r.cfg)
	uh := handlers.NewUserHandler(us)

	prg := ro.Group("/api")
	rg := ro.Group("/api/users")
	ag := ro.Group("/api/auth")

	RegisterPublicRoutes(prg)
	RegisterPublicUserRoutes(rg)
	RegisterPublicAuthRoutes(ag, ah)
	RegisterProtectedAuthRoutes(ag, ah)
	RegisterProtectedUserRoutes(rg, uh)

	return ro
}
