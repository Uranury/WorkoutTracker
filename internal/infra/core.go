package infra

import (
	"fmt"
	"github.com/Uranury/WorkoutTracker/pkg/config"
	db2 "github.com/Uranury/WorkoutTracker/pkg/db"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"time"
)

type Deps struct {
	DBConn      *sqlx.DB
	RedisClient *redis.Client
	HTTPClient  *http.Client
	Logger      *slog.Logger
	Config      *config.Config
}

func New() (*Deps, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}

	logger := slog.Default()

	if err := db2.RunMigrations(cfg.Driver, cfg.DSN(), cfg.MigrationsPath, logger); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	dbConn, err := db2.InitDB(cfg.Driver, cfg.DSN(), logger)
	if err != nil {
		return nil, nil, err
	}

	//rdb := redis.NewClient(&redis.Options{
	//	Addr: cfg.RedisAddr,
	//})

	//if err := rdb.Ping(context.Background()).Err(); err != nil {
	//	logger.Warn("Could not connect to redis", "error", err)
	//}

	httpClient := &http.Client{
		Timeout: time.Second * 20,
	}

	deps := &Deps{
		DBConn:     dbConn,
		HTTPClient: httpClient,
		Logger:     logger,
		Config:     cfg,
	}

	cleanup := func() {
		if err := dbConn.Close(); err != nil {
			logger.Error("Failed to close database connection", "error", err)
		}
		logger.Info("Infrastructure cleaned up")
	}

	return deps, cleanup, nil
}
