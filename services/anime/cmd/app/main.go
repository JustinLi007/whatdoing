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

	"github.com/JustinLi007/whatdoing/libs/go/config"
	"github.com/JustinLi007/whatdoing/services/anime/internal/pubsub"
	"github.com/JustinLi007/whatdoing/services/anime/internal/server"
)

func main() {
	c := config.NewBuilder().
		Cli("mode").
		Cli("env").
		Env("DB_URL").
		Build()
	c.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)

	switch c.Get("mode") {
	case "service":
		server := server.NewServer(ctx, c)

		go gracefulShutdownServer(server, done, cancel)

		log.Printf("Listening on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %v", err)
		}
	case "pub":
		pub := pubsub.NewPublisher(c)

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
