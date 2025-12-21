package main

import (
	"context"
	"github.com/Uranury/WorkoutTracker/internal/http_server"
	"github.com/Uranury/WorkoutTracker/internal/infra"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	deps, cleanup, err := infra.NewDeps()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	app := infra.NewApp(deps)
	server := http_server.NewHTTPServer(app)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}
