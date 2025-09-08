package probes

import (
	healthQuery "github.com/go-clean/internal/probes/application/query"
	pingQuery "github.com/go-clean/internal/probes/application/query"
	healthInfra "github.com/go-clean/internal/probes/infrastructure"
	healthHttp "github.com/go-clean/internal/probes/presentation/http"
	pingHttp "github.com/go-clean/internal/probes/presentation/http"
	"github.com/go-clean/platform/logger"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// ProvidePingQueryHandler provides a ping query handler
func ProvidePingQueryHandler(logger logger.Logger) *pingQuery.PingQueryHandler {
	return pingQuery.NewPingQueryHandler(logger)
}

// ProvidePingHandler provides a ping HTTP handler
func ProvidePingHandler(logger logger.Logger, pingQueryHandler *pingQuery.PingQueryHandler) *pingHttp.PingHandler {
	return pingHttp.NewPingHandler(logger, pingQueryHandler)
}

// ProvideDatabaseChecker provides a database checker
func ProvideDatabaseChecker(logger logger.Logger, db *pgxpool.Pool) *healthInfra.DatabaseChecker {
	return healthInfra.NewDatabaseChecker(logger, db)
}

// ProvideRedisChecker provides a Redis checker
func ProvideRedisChecker(logger logger.Logger, redisClient *redis.Client) *healthInfra.RedisChecker {
	return healthInfra.NewRedisChecker(logger, redisClient)
}

// ProvideHealthQueryHandler provides a health query handler
func ProvideHealthQueryHandler(logger logger.Logger, databaseChecker *healthInfra.DatabaseChecker, redisChecker *healthInfra.RedisChecker) *healthQuery.GetHealthQueryHandler {
	return healthQuery.NewGetHealthQueryHandler(logger, databaseChecker, redisChecker)
}

// ProvideHealthService provides a health service
func ProvideHealthService(logger logger.Logger, healthQueryHandler *healthQuery.GetHealthQueryHandler) *healthQuery.HealthService {
	return healthQuery.NewHealthService(logger, healthQueryHandler)
}

// ProvideLivenessQueryHandler provides a liveness query handler
func ProvideLivenessQueryHandler(logger logger.Logger) *healthQuery.GetLivenessQueryHandler {
	return healthQuery.NewGetLivenessQueryHandler(logger)
}

// ProvideLivenessService provides a liveness service
func ProvideLivenessService(logger logger.Logger, livenessQueryHandler *healthQuery.GetLivenessQueryHandler) *healthQuery.LivenessService {
	return healthQuery.NewLivenessService(logger, livenessQueryHandler)
}

// ProvideHealthHandler provides a health HTTP handler
func ProvideHealthHandler(logger logger.Logger, healthService *healthQuery.HealthService, livenessService *healthQuery.LivenessService) *healthHttp.HealthHandler {
	return healthHttp.NewHealthHandler(logger, healthService, livenessService)
}

// ProbesSet is a wire provider set for all probes dependencies
var ProbesSet = wire.NewSet(
	ProvidePingQueryHandler,
	ProvidePingHandler,
	ProvideDatabaseChecker,
	ProvideRedisChecker,
	ProvideHealthQueryHandler,
	ProvideHealthService,
	ProvideLivenessQueryHandler,
	ProvideLivenessService,
	ProvideHealthHandler,
)
