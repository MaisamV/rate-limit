package command

import (
	"context"
	"fmt"
	"time"
	
	"github.com/go-clean/internal/ratelimit/ports"
	"github.com/go-clean/platform/logger"
)

// CheckRateLimitWithDetailCommand represents a command to check rate limit with detailed response
type CheckRateLimitWithDetailCommand struct {
	UserID string
	Limit  int
}

// CheckRateLimitWithDetailResponse represents the detailed response from rate limit check
type CheckRateLimitWithDetailResponse struct {
	Remaining int
	ResetTime time.Duration
	Allowed   bool
}

// CheckRateLimitWithDetailCommandHandler handles rate limit checking commands with detailed response
type CheckRateLimitWithDetailCommandHandler struct {
	logger     logger.Logger
	repository ports.RateLimitRepository
}

// NewCheckRateLimitWithDetailCommandHandler creates a new CheckRateLimitWithDetailCommandHandler
func NewCheckRateLimitWithDetailCommandHandler(
	logger logger.Logger,
	repository ports.RateLimitRepository,
) *CheckRateLimitWithDetailCommandHandler {
	return &CheckRateLimitWithDetailCommandHandler{
		logger:     logger,
		repository: repository,
	}
}

// Handle processes the CheckRateLimitWithDetailCommand
func (h *CheckRateLimitWithDetailCommandHandler) Handle(ctx context.Context, cmd CheckRateLimitWithDetailCommand) (*CheckRateLimitWithDetailResponse, error) {
	h.logger.Info().Str("user_id", cmd.UserID).Int("limit", cmd.Limit).Msg("Processing rate limit check with detail")
	
	if cmd.UserID == "" {
		h.logger.Error().Str("user_id", cmd.UserID).Msg("Invalid user ID provided")
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	
	if cmd.Limit <= 0 {
		h.logger.Error().Int("limit", cmd.Limit).Msg("Invalid limit provided")
		return nil, fmt.Errorf("limit must be greater than 0")
	}
	
	remaining, resetTime, err := h.repository.RateLimitWithDetail(cmd.UserID, cmd.Limit)
	if err != nil {
		h.logger.Error().Str("user_id", cmd.UserID).Err(err).Msg("Failed to check rate limit with detail")
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}
	
	allowed := remaining > 0
	
	response := &CheckRateLimitWithDetailResponse{
		Remaining: remaining,
		ResetTime: resetTime,
		Allowed:   allowed,
	}
	
	h.logger.Info().Str("user_id", cmd.UserID).Int("limit", cmd.Limit).Int("remaining", remaining).Dur("reset_time", resetTime).Bool("allowed", allowed).Msg("Rate limit check with detail completed")
	
	return response, nil
}