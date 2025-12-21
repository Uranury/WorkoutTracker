package http_server

import (
	"context"
	"errors"
	"github.com/Uranury/WorkoutTracker/internal/infra"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HTTPServer struct {
	router *gin.Engine
	server *http.Server
	app    *infra.App
}

func NewHTTPServer(app *infra.App) *HTTPServer {
	server := &HTTPServer{
		router: gin.Default(),
		app:    app,
	}

	server.setupRoutes()
	return server
}

func (h *HTTPServer) Start() error {
	h.server = &http.Server{
		Addr:           h.app.Config().ListenAddr,
		Handler:        h.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	h.app.Logger().Info("Starting HTTP server", "address", h.server.Addr)

	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.app.Logger().Error("Server error", "error", err)
		}
	}()

	return nil
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
	if h.server == nil {
		return nil
	}

	h.app.Logger().Info("Shutting down HTTP server")

	if err := h.server.Shutdown(ctx); err != nil {
		h.app.Logger().Error("Failed to gracefully shutdown", "error", err)
		return err
	}

	h.app.Logger().Info("HTTP server shut down gracefully")
	return nil
}
