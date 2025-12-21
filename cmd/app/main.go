package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Uranury/WorkoutTracker/internal/http_server"
	"github.com/Uranury/WorkoutTracker/internal/infra"
)

func main() {
	deps, cleanup, err := infra.New()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	server := http_server.NewHTTPServer(deps)
	srv := server.Start()

	deps.Logger.Info("Server started", "addr", deps.Config.ListenAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	deps.Logger.Info("Received shutdown signal", "signal", sig)

	// Attempt graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		deps.Logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	deps.Logger.Info("Server stopped gracefully")
}
