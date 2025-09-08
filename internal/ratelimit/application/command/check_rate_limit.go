package command

import (
	"context"
	"fmt"
	
	"github.com/go-clean/internal/ratelimit/ports"
	"github.com/go-clean/platform/logger"
)

// CheckRateLimitCommand represents a command to check rate limit for a user
type CheckRateLimitCommand struct {
	UserID string
	Limit  int
}

// CheckRateLimitCommandHandler handles rate limit checking commands
type CheckRateLimitCommandHandler struct {
	logger     logger.Logger
	repository ports.RateLimitRepository
}

// NewCheckRateLimitCommandHandler creates a new CheckRateLimitCommandHandler
func NewCheckRateLimitCommandHandler(
	logger logger.Logger,
	repository ports.RateLimitRepository,
) *CheckRateLimitCommandHandler {
	return &CheckRateLimitCommandHandler{
		logger:     logger,
		repository: repository,
	}
}

// Handle processes the CheckRateLimitCommand
func (h *CheckRateLimitCommandHandler) Handle(ctx context.Context, cmd CheckRateLimitCommand) (bool, error) {
	h.logger.Info().Str("user_id", cmd.UserID).Int("limit", cmd.Limit).Msg("Processing rate limit check")
	
	if cmd.UserID == "" {
		h.logger.Error().Str("user_id", cmd.UserID).Msg("Invalid user ID provided")
		return false, fmt.Errorf("user ID cannot be empty")
	}
	
	if cmd.Limit <= 0 {
		h.logger.Error().Int("limit", cmd.Limit).Msg("Invalid limit provided")
		return false, fmt.Errorf("limit must be greater than 0")
	}
	
	allowed := h.repository.RateLimit(cmd.UserID, cmd.Limit)
	
	h.logger.Info().Str("user_id", cmd.UserID).Int("limit", cmd.Limit).Bool("allowed", allowed).Msg("Rate limit check completed")
	
	return allowed, nil
}