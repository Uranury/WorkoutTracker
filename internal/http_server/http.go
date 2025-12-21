package http_server

import (
	"context"
	"github.com/Uranury/WorkoutTracker/internal/infra"
	"github.com/Uranury/WorkoutTracker/internal/middleware"
	"github.com/Uranury/WorkoutTracker/internal/services"
	"github.com/Uranury/WorkoutTracker/internal/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HTTPServer struct {
	router         *gin.Engine
	deps           *infra.Deps
	userHandler    *user.Handler
	authMiddleware *middleware.Auth
}

func NewHTTPServer(deps *infra.Deps) *HTTPServer {
	authService := services.NewAuth(deps.Config.JWTKey)
	authMiddleware := middleware.NewAuth(authService)

	userRepo := user.NewRepository(deps.DBConn)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, authService)

	server := &HTTPServer{
		router:         gin.Default(),
		deps:           deps,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}

	server.setupRoutes()
	return server
}

func (h *HTTPServer) Start() *http.Server {
	srv := &http.Server{
		Addr:           h.deps.Config.ListenAddr,
		Handler:        h.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.deps.Logger.Error("Server failed", "error", err)
		}
	}()

	return srv
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
	srv := &http.Server{
		Addr:    h.deps.Config.ListenAddr,
		Handler: h.router,
	}
	return srv.Shutdown(ctx)
}
