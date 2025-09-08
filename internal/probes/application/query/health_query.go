package query

import (
	"context"

	"github.com/go-clean/internal/probes/domain"
	"github.com/go-clean/internal/probes/ports"
	"github.com/go-clean/platform/logger"
)

// GetHealthQuery represents a query to get system health status
type GetHealthQuery struct{}

// GetHealthQueryHandler handles health check queries
type GetHealthQueryHandler struct {
	logger          logger.Logger
	databaseChecker ports.DatabaseChecker
	redisChecker    ports.RedisChecker
}

// NewGetHealthQueryHandler creates a new health query handler
func NewGetHealthQueryHandler(logger logger.Logger, databaseChecker ports.DatabaseChecker, redisChecker ports.RedisChecker) *GetHealthQueryHandler {
	return &GetHealthQueryHandler{
		logger:          logger,
		databaseChecker: databaseChecker,
		redisChecker:    redisChecker,
	}
}

// Handle executes the health check query
func (h *GetHealthQueryHandler) Handle(ctx context.Context, query GetHealthQuery) (*domain.HealthResponse, error) {
	h.logger.Info().Msg("Starting health check")
	response := domain.NewHealthResponse()

	// Check database connectivity
	if h.databaseChecker != nil {
		h.logger.Debug().Msg("Checking database connectivity")
		dbHealthy, dbResponseTime, err := h.databaseChecker.CheckDatabase(ctx)
		if err != nil {
			h.logger.Error().Err(err).Msg("Database health check failed")
			response.AddCheck("database", domain.CheckStatusDown, 0)
		} else {
			status := domain.CheckStatusUp
			if !dbHealthy {
				status = domain.CheckStatusDown
				h.logger.Warn().Msg("Database is not healthy")
			} else {
				h.logger.Info().Int64("response_time_ms", dbResponseTime.Milliseconds()).Msg("Database health check passed")
			}
			response.AddCheck("database", status, dbResponseTime.Milliseconds())
		}
	}

	// Check Redis connectivity
	if h.redisChecker != nil {
		h.logger.Debug().Msg("Checking Redis connectivity")
		redisHealthy, redisResponseTime, err := h.redisChecker.CheckRedis(ctx)
		if err != nil {
			h.logger.Error().Err(err).Msg("Redis health check failed")
			response.AddCheck("redis", domain.CheckStatusDown, 0)
		} else {
			status := domain.CheckStatusUp
			if !redisHealthy {
				status = domain.CheckStatusDown
				h.logger.Warn().Msg("Redis is not healthy")
			} else {
				h.logger.Info().Int64("response_time_ms", redisResponseTime.Milliseconds()).Msg("Redis health check passed")
			}
			response.AddCheck("redis", status, redisResponseTime.Milliseconds())
		}
	}

	// Determine overall status
	response.DetermineOverallStatus()
	h.logger.Info().Bool("is_healthy", response.IsHealthy()).Msg("Health check completed")

	return response, nil
}

// HealthService implements the HealthService port
type HealthService struct {
	logger       logger.Logger
	queryHandler *GetHealthQueryHandler
}

// NewHealthService creates a new health service
func NewHealthService(logger logger.Logger, queryHandler *GetHealthQueryHandler) *HealthService {
	return &HealthService{
		logger:       logger,
		queryHandler: queryHandler,
	}
}

// GetHealthStatus returns the current health status
func (s *HealthService) GetHealthStatus(ctx context.Context) (*domain.HealthResponse, error) {
	s.logger.Debug().Msg("Health status requested")
	return s.queryHandler.Handle(ctx, GetHealthQuery{})
}
