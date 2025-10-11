package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/configs"
	"github.com/JustinLi007/whatdoing/libs/go/utils"
	"github.com/JustinLi007/whatdoing/progress/internal/database"
	"github.com/JustinLi007/whatdoing/progress/migrations"
)

type Server struct {
	Port int
}

func NewServer(ctx context.Context, c *configs.Config) *http.Server {
	server := &Server{}

	connStr := c.Get("DB_URL")
	if connStr == "" {
		log.Fatalf("error: %v", fmt.Errorf("invalid conn str"))
	}

	db, err := database.NewDb(connStr)
	utils.RequireNoError(err, "error: service failed to connect to db")

	err = db.MigrateFS(migrations.Fs, ".")
	utils.RequireNoError(err, "error: service failed to connect to db")

	return &http.Server{
		Handler: server.RegisterRoutes(),
	}
}
