package query

import (
	"context"
	"time"

	"github.com/go-clean/internal/probes/domain"
	"github.com/go-clean/platform/logger"
)

// GetLivenessQuery represents a query to get liveness status
type GetLivenessQuery struct{}

// GetLivenessQueryHandler handles liveness check queries
type GetLivenessQueryHandler struct {
	logger    logger.Logger
	startTime time.Time
}

// NewGetLivenessQueryHandler creates a new liveness query handler
func NewGetLivenessQueryHandler(logger logger.Logger) *GetLivenessQueryHandler {
	return &GetLivenessQueryHandler{
		logger:    logger,
		startTime: time.Now().UTC(),
	}
}

// Handle executes the liveness check query
func (h *GetLivenessQueryHandler) Handle(ctx context.Context, query GetLivenessQuery) (*domain.LivenessResponse, error) {
	h.logger.Debug().Msg("Processing liveness check")
	
	// Create liveness response with uptime
	response := domain.NewLivenessResponse(h.startTime)
	
	h.logger.Debug().
		Int64("uptime_seconds", response.UptimeSeconds).
		Str("status", string(response.Status)).
		Msg("Liveness check completed")
	
	return response, nil
}

// LivenessService implements the liveness service
type LivenessService struct {
	logger       logger.Logger
	queryHandler *GetLivenessQueryHandler
}

// NewLivenessService creates a new liveness service
func NewLivenessService(logger logger.Logger, queryHandler *GetLivenessQueryHandler) *LivenessService {
	return &LivenessService{
		logger:       logger,
		queryHandler: queryHandler,
	}
}

// GetLivenessStatus returns the current liveness status
func (s *LivenessService) GetLivenessStatus(ctx context.Context) (*domain.LivenessResponse, error) {
	s.logger.Debug().Msg("Liveness status requested")
	return s.queryHandler.Handle(ctx, GetLivenessQuery{})
}