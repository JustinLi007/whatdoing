package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/JustinLi007/whatdoing/libs/go/configs"
	"github.com/JustinLi007/whatdoing/services/gateway/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)

	c := configs.NewBuilder().
		Env("SERVER_PORT").
		Env("JWK_URL").
		Env("JWT_ISSUER").
		Env("JWT_AUDIENCE").
		Build()
	c.Parse()

	server := server.NewServer(ctx, c)

	go gracefullShutdown(server, done, cancel)

	log.Printf("listening on %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error: %v", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefullShutdown(server *http.Server, done chan bool, ctxCancel context.CancelFunc) {
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-signalCtx.Done()

	ctxCancel()

	fmt.Println()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}
