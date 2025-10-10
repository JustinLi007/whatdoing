package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JustinLi007/whatdoing/services/anime/internal/configs"
	"github.com/JustinLi007/whatdoing/services/anime/internal/database"
	"github.com/JustinLi007/whatdoing/services/anime/internal/handler"
	"github.com/JustinLi007/whatdoing/services/anime/migrations"
)

type Server struct {
	Port         int
	animeHandler handler.HandlerAnime
}

func NewServer(ctx context.Context) *http.Server {
	server := &Server{}

	configs := configs.NewConfigs()

	if err := configs.LoadEnv(); err != nil {
		log.Fatalf("error: %v", err)
	}

	server.Port = configs.ConfigServer.Port

	connStr := configs.ConfigDb.PostgresConnStr()
	if connStr == "" {
		log.Fatalf("error: %v", fmt.Errorf("invalid conn str"))
	}

	db, err := database.NewDb(connStr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := db.MigrateFS(migrations.Fs, "."); err != nil {
		log.Fatalf("error: %v", err)
	}

	animeService := database.NewServiceAnime(db)
	animeHandler := handler.NewHandlerAnime(animeService)
	server.animeHandler = animeHandler

	return &http.Server{
		// Addr:    fmt.Sprintf(":%d", server.Port),
		Handler: server.RegisterRoutes(),
	}
}
