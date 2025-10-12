package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/config"
	"github.com/JustinLi007/whatdoing/libs/go/util"
	"github.com/JustinLi007/whatdoing/services/progress/internal/database"
	"github.com/JustinLi007/whatdoing/services/progress/migrations"
)

type Server struct {
	Port int
}

func NewServer(ctx context.Context, c *config.Config) *http.Server {
	server := &Server{}

	connStr := c.Get("DB_URL")
	if connStr == "" {
		log.Fatalf("error: %v", fmt.Errorf("invalid conn str"))
	}

	db, err := database.NewDb(connStr)
	util.RequireNoError(err, "error: service failed to connect to db")

	err = db.MigrateFS(migrations.Fs, ".")
	util.RequireNoError(err, "error: service failed to connect to db")

	return &http.Server{
		Handler: server.RegisterRoutes(),
	}
}
