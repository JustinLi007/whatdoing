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

	"service-anime/internal/configs"
	"service-anime/internal/pubsub"
	"service-anime/internal/server"
)

func main() {
	c := configs.NewConfigs(configs.WithDbConfig())

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)

	switch c.ModeApp {
	case configs.APP_SERVICE:
		server := server.NewServer(ctx)

		go gracefulShutdownServer(server, done, cancel)

		log.Printf("Listening on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %v", err)
		}
	case configs.APP_PUB:
		pub := pubsub.NewPublisher()

		go gracefulShutdownPub(done, cancel)

		pub.Start(ctx)
	default:
		log.Panicf("error: unknown mode")
	}

	<-done
	log.Printf("Graceful shutdown complete.")
}

func gracefulShutdownServer(server *http.Server, done chan bool, ctxCancel context.CancelFunc) {
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	<-signalCtx.Done()

	ctxCancel()

	fmt.Println()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func gracefulShutdownPub(done chan bool, ctxCancel context.CancelFunc) {
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	<-signalCtx.Done()

	ctxCancel()

	fmt.Println()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	log.Println("Server exiting")
	done <- true
}
