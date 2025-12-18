package http_server

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
)

type APIServer struct {
	DBConn      *sqlx.DB
	RedisClient *redis.Client
	HTTPClient  *http.Client
	Logger      *slog.Logger
}
