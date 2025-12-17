package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cfex/microservices-in-go/services/users-service/cmd/config"
)

var ErrDatabaseRefuseConnection = errors.New("database refused connection")

func ConnectDB(cfg config.Config) *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		panic(ErrDatabaseRefuseConnection)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Database.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(fmt.Errorf("cannot connect to database: %w", err))
	}

	return db
}
