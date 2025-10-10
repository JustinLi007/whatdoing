package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JustinLi007/whatdoing/services/gateway/internal/configs"
	"github.com/JustinLi007/whatdoing/services/gateway/internal/middleware"
	"github.com/JustinLi007/whatdoing/services/gateway/internal/service"
	"github.com/JustinLi007/whatdoing/services/gateway/internal/verifier"
)

type Server struct {
	Port       int
	Issuer     string
	Audience   string
	JwkUrl     string
	Middleware middleware.Middleware
	ServiceMap service.ServiceMap
	Verifier   verifier.Verifier
}

func NewServer(ctx context.Context) *http.Server {
	server := &Server{}

	configs := configs.NewConfigs()

	if err := configs.LoadEnv(); err != nil {
		log.Fatalf("error: %v", err)
	}

	server.Port = configs.ConfigServer.Port
	server.JwkUrl = configs.ConfigServer.JwkUrl
	server.Issuer = configs.ConfigServer.Issuer
	server.Audience = configs.ConfigServer.Audience

	// verifier
	verifier, err := verifier.NewVerifier(server.JwkUrl, server.Issuer, server.Audience)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// services
	serviceMap := service.NewServiceMap()

	// middleware
	middleware := middleware.NewMiddleware(verifier, serviceMap)

	// handlers

	server.Middleware = middleware
	server.ServiceMap = serviceMap
	server.Verifier = verifier

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", server.Port),
		Handler: server.RegisterServices(),
	}
}
