package platform

import (
	"github.com/go-clean/platform/config"
	"github.com/go-clean/platform/database"
	"github.com/go-clean/platform/http"
	"github.com/go-clean/platform/logger"
	platformRedis "github.com/go-clean/platform/redis"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// ProvideLogger provides a logger instance
func ProvideLogger() logger.Logger {
	return logger.New()
}

// ProvideConfig provides a configuration instance
func ProvideConfig(log logger.Logger) (*config.Config, error) {
	return config.Load(log)
}

// ProvideDatabase provides a database connection pool
func ProvideDatabase(cfg *config.Config, log logger.Logger) (*pgxpool.Pool, error) {
	return database.NewConnection(cfg.Database, log)
}

// ProvideRedis provides a Redis client
func ProvideRedis(cfg *config.Config, log logger.Logger) (*redis.Client, error) {
	return platformRedis.NewClient(cfg.Redis, log)
}

// ProvideHTTPServer provides an HTTP server instance
func ProvideHTTPServer(cfg *config.Config, log logger.Logger) *http.Server {
	return http.NewServer(cfg.Server.Port, log)
}

// PlatformSet is a wire provider set for all platform dependencies
var PlatformSet = wire.NewSet(
	ProvideLogger,
	ProvideConfig,
	ProvideDatabase,
	ProvideRedis,
	ProvideHTTPServer,
)
