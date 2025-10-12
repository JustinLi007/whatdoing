package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/config"
	"github.com/JustinLi007/whatdoing/libs/go/util"
	"github.com/JustinLi007/whatdoing/services/anime/internal/database"
	"github.com/JustinLi007/whatdoing/services/anime/internal/handler"
	"github.com/JustinLi007/whatdoing/services/anime/migrations"
)

type Server struct {
	Port         int
	animeHandler handler.HandlerAnime
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

	animeService := database.NewServiceAnime(db)
	animeHandler := handler.NewHandlerAnime(animeService)
	server.animeHandler = animeHandler

	return &http.Server{
		Handler: server.RegisterRoutes(),
	}
}
