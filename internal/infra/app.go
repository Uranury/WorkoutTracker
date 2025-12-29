package infra

import (
	"github.com/Uranury/WorkoutTracker/internal/auth"
	"github.com/Uranury/WorkoutTracker/internal/middleware"
	"github.com/Uranury/WorkoutTracker/internal/user"
	"github.com/Uranury/WorkoutTracker/pkg/config"
	"log/slog"
)

type App struct {
	deps *Deps

	// Shared services
	authService auth.Service

	// Module services (lazy-loaded or pre-initialized)
	userService user.Service
	// ...
	userHandler    *user.Handler
	authMiddleware *middleware.Auth
}

func NewApp(deps *Deps) *App {
	app := &App{
		deps: deps,
	}

	app.initAuth()
	app.authMiddleware = middleware.NewAuth(app.authService)

	// Initialize modules in dependency order
	app.initUser()
	// app.initWorkout()
	// ...

	return app
}

func (a *App) initAuth() {
	logger := a.deps.Logger.With("module", "auth")
	authRepo := auth.NewRepository(a.deps.DBConn)
	a.authService = auth.NewAuth(a.deps.Config.JWTKey, a.deps.DBConn, logger, authRepo)
}

func (a *App) initUser() {
	logger := a.deps.Logger.With("module", "user")
	userRepo := user.NewRepository(a.deps.DBConn, logger)
	a.userService = user.NewService(userRepo, logger)
	a.userHandler = user.NewHandler(a.userService, a.authService)
}

func (a *App) UserHandler() *user.Handler {
	return a.userHandler
}

func (a *App) initWorkout() {
	// workoutRepo := workout.NewRepository(a.deps.DBConn)
	// workoutService needs userService - it's already initialized
	// a.workoutService = workout.NewService(workoutRepo, a.userService)
}

/*
func (a *App) WorkoutHandler() *workout.Handler {
	return workout.NewHandler(a.workoutService, a.authService)
}
*/

func (a *App) AuthMiddleware() *middleware.Auth {
	return a.authMiddleware
}

func (a *App) Logger() *slog.Logger {
	return a.deps.Logger
}

func (a *App) Config() *config.Config {
	return a.deps.Config
}
