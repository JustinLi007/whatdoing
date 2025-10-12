package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/JustinLi007/whatdoing/libs/go/config"
	"github.com/JustinLi007/whatdoing/libs/go/util"
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

func NewServer(ctx context.Context, c *config.Config) *http.Server {
	server := &Server{}

	port, err := strconv.Atoi(c.Get("SERVER_PORT"))
	util.RequireNoError(err, "error: failed to parse port")
	server.Port = port
	server.JwkUrl = c.Get("JWK_URL")
	server.Issuer = c.Get("JWT_ISSUER")
	server.Audience = c.Get("JWT_AUDIENCE")

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
