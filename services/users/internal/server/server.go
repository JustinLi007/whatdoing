package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JustinLi007/whatdoing/services/users/internal/configs"
	"github.com/JustinLi007/whatdoing/services/users/internal/database"
	"github.com/JustinLi007/whatdoing/services/users/internal/handlers"
	"github.com/JustinLi007/whatdoing/services/users/internal/middleware"
	"github.com/JustinLi007/whatdoing/services/users/internal/signer"
	"github.com/JustinLi007/whatdoing/services/users/migrations"
)

type Server struct {
	Port          int
	Iss           string
	Aud           string
	Middleware    middleware.Middleware
	HandlerSigner handlers.HandlerSigner
	HandlerUsers  handlers.HandlerUsers
}

func NewServer(ctx context.Context) *http.Server {
	server := &Server{}

	configs := configs.NewConfigs()

	if err := configs.LoadEnv(); err != nil {
		log.Fatalf("error: %v", err)
	}

	server.Port = configs.ConfigServer.Port
	server.Iss = configs.ConfigServer.Iss
	server.Aud = configs.ConfigServer.Aud

	// database
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

	// middleware
	middleware := middleware.NewMiddleware()

	// signer
	signer, err := signer.NewSigner(server.Iss, server.Aud)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// services
	usersService := database.NewServiceUsers(db)

	// handlers
	signerHandler := handlers.NewHandlerSigner(signer)
	usersHandler := handlers.NewHandlerUsers(signer, usersService)

	server.Middleware = middleware
	server.HandlerSigner = signerHandler
	server.HandlerUsers = usersHandler

	return &http.Server{
		// Addr:    fmt.Sprintf(":%d", server.Port),
		Handler: server.RegisterRoutes(),
	}
}
