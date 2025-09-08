package http

import (
	"github.com/go-clean/internal/probes/application/query"
	"github.com/go-clean/platform/logger"
	"github.com/gofiber/fiber/v2"
)

// PingHandler handles HTTP requests for ping endpoints
type PingHandler struct {
	logger           logger.Logger
	pingQueryHandler *query.PingQueryHandler
}

// NewPingHandler creates a new ping HTTP handler
func NewPingHandler(logger logger.Logger, pingQueryHandler *query.PingQueryHandler) *PingHandler {
	return &PingHandler{
		logger:           logger,
		pingQueryHandler: pingQueryHandler,
	}
}

// Ping handles GET /ping requests
// @Summary Ping endpoint
// @Description Returns a simple PONG response to verify service is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} domain.PingResponse
// @Router /ping [get]
func (h *PingHandler) Ping(c *fiber.Ctx) error {
	h.logger.Info().Str("endpoint", "/ping").Msg("Ping endpoint called")
	ctx := c.Context()

	response, err := h.pingQueryHandler.Handle(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to handle ping request")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	h.logger.Debug().Msg("Ping request handled successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}

// RegisterRoutes registers ping routes with the fiber app
func (h *PingHandler) RegisterRoutes(app *fiber.App) {
	h.logger.Info().Msg("Registering ping routes")
	app.Get("/ping", h.Ping)
	h.logger.Debug().Str("route", "/ping").Msg("Ping route registered")
}
