package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"interviews-ai/internal/common/config"
	"interviews-ai/internal/common/routes"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv, err := routes.NewServer(cfg)

	if err != nil {
		log.Fatal("failed to create server: %v", err)
	}

	httpServer := &http.Server{
		Addr:    cfg.ServerAddress(),
		Handler: srv,
	}

	go func() {
		log.Printf("Server is listening on %s\n", cfg.ServerAddress())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe() error: %V", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
