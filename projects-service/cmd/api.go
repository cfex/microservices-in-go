package api

import (
	"database/sql"
	"net/http"

	"github.com/cfex/microservices-in-go/projects-service/cmd/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type router struct {
	DB  *sql.DB
	cfg *config.Config
}

func CreateRouter(db *sql.DB, cfg *config.Config) *router {
	return &router{DB: db, cfg: cfg}
}

func (r *router) NewRouter() *gin.Engine {
	ro := gin.Default()
	ro.Use(cors.New(r.cfg.Server.Cors))

	prg := ro.Group("/api")

	prg.GET("/projects", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"m": "hello"})
	})

	return ro
}
