package http

import (
	"net/http"

	"github.com/go-clean/internal/ratelimit/application/command"
	"github.com/go-clean/platform/logger"
	"github.com/gofiber/fiber/v2"
)

// RateLimitHandler handles rate limit HTTP requests
type RateLimitHandler struct {
	logger         logger.Logger
	commandHandler *command.CheckRateLimitWithDetailCommandHandler
}

// NewRateLimitHandler creates a new rate limit handler
func NewRateLimitHandler(
	logger logger.Logger,
	commandHandler *command.CheckRateLimitWithDetailCommandHandler,
) *RateLimitHandler {
	return &RateLimitHandler{
		logger:         logger,
		commandHandler: commandHandler,
	}
}

// CheckRateLimit handles POST /rate-limit requests
// @Summary Check rate limit for a user
// @Description Checks if a user has exceeded their rate limit and returns detailed information
// @Tags Rate Limit
// @Accept json
// @Produce json
// @Param request body RateLimitRequest true "Rate limit check request"
// @Success 200 {object} RateLimitResponse "Rate limit check successful"
// @Success 429 {object} RateLimitResponse "Rate limit exceeded"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate-limit [post]
func (h *RateLimitHandler) CheckRateLimit(c *fiber.Ctx) error {
	h.logger.Info().Str("endpoint", "/rate-limit").Msg("Rate limit check endpoint called")
	ctx := c.Context()

	// Parse request body
	var req RateLimitRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if req.UserID == "" {
		h.logger.Error().Msg("Missing user_id in request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	if req.Limit <= 0 {
		h.logger.Error().Int("limit", req.Limit).Msg("Invalid limit in request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "limit must be greater than 0",
		})
	}

	// Create command
	cmd := command.CheckRateLimitWithDetailCommand{
		UserID: req.UserID,
		Limit:  req.Limit,
	}

	// Execute command
	result, err := h.commandHandler.Handle(ctx, cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Int("limit", req.Limit).Msg("Failed to check rate limit")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to check rate limit",
			"details": err.Error(),
		})
	}

	// Create response
	response := RateLimitResponse{
		Allowed:   result.Allowed,
		Remaining: result.Remaining,
		ResetTime: int64(result.ResetTime.Seconds()),
		UserID:    req.UserID,
		Limit:     req.Limit,
	}

	// Return appropriate HTTP status
	statusCode := http.StatusOK
	if !result.Allowed {
		statusCode = http.StatusTooManyRequests
		h.logger.Warn().Str("user_id", req.UserID).Int("limit", req.Limit).Int("remaining", result.Remaining).Msg("Rate limit exceeded")
	} else {
		h.logger.Info().Str("user_id", req.UserID).Int("limit", req.Limit).Int("remaining", result.Remaining).Msg("Rate limit check passed")
	}

	return c.Status(statusCode).JSON(response)
}

// RegisterRoutes registers rate limit related routes
func (h *RateLimitHandler) RegisterRoutes(router fiber.Router) {
	h.logger.Info().Msg("Registering rate limit routes")
	router.Post("/rate-limit", h.CheckRateLimit)
	h.logger.Debug().Str("route", "/rate-limit").Msg("Rate limit route registered")
}

// RateLimitRequest represents the request body for rate limit check
type RateLimitRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Limit  int    `json:"limit" validate:"required,min=1"`
}

// RateLimitResponse represents the response body for rate limit check
type RateLimitResponse struct {
	Allowed   bool   `json:"allowed"`
	Remaining int    `json:"remaining"`
	ResetTime int64  `json:"reset_time_seconds"`
	UserID    string `json:"user_id"`
	Limit     int    `json:"limit"`
}