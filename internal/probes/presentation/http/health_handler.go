package http

import (
	"net/http"

	"github.com/go-clean/internal/probes/application/query"
	"github.com/go-clean/platform/logger"
	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check HTTP requests
type HealthHandler struct {
	logger          logger.Logger
	healthService   *query.HealthService
	livenessService *query.LivenessService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(logger logger.Logger, healthService *query.HealthService, livenessService *query.LivenessService) *HealthHandler {
	return &HealthHandler{
		logger:          logger,
		healthService:   healthService,
		livenessService: livenessService,
	}
}

// GetHealth handles GET /health requests
// @Summary Get system health status
// @Description Returns the health status of the system including database and Redis connectivity
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} domain.HealthResponse "System is healthy"
// @Success 503 {object} domain.HealthResponse "System is unhealthy"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /health [get]
func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
	h.logger.Info().Str("endpoint", "/health").Msg("Health check endpoint called")
	ctx := c.Context()

	// Get health status from service
	healthResponse, err := h.healthService.GetHealthStatus(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get health status from service")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to check system health",
			"details": err.Error(),
		})
	}

	// Return appropriate HTTP status based on health
	statusCode := http.StatusOK
	if !healthResponse.IsHealthy() {
		statusCode = http.StatusServiceUnavailable
		h.logger.Warn().Int("status_code", statusCode).Bool("is_healthy", false).Msg("System is unhealthy")
	} else {
		h.logger.Info().Int("status_code", statusCode).Bool("is_healthy", true).Msg("System is healthy")
	}

	return c.Status(statusCode).JSON(healthResponse)
}

// GetLiveness handles GET /liveness requests
// @Summary Get liveness status
// @Description Returns the liveness status of the service for Kubernetes liveness probes
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} domain.LivenessResponse "Service is alive"
// @Success 503 {object} domain.LivenessResponse "Service is dead"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /liveness [get]
func (h *HealthHandler) GetLiveness(c *fiber.Ctx) error {
	h.logger.Info().Str("endpoint", "/liveness").Msg("Liveness check endpoint called")
	ctx := c.Context()

	// Get liveness status from service
	livenessResponse, err := h.livenessService.GetLivenessStatus(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get liveness status from service")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to check service liveness",
			"details": err.Error(),
		})
	}

	// Return appropriate HTTP status based on liveness
	statusCode := http.StatusOK
	if !livenessResponse.IsAlive() {
		statusCode = http.StatusServiceUnavailable
		h.logger.Warn().Int("status_code", statusCode).Bool("is_alive", false).Msg("Service is not alive")
	} else {
		h.logger.Info().Int("status_code", statusCode).Bool("is_alive", true).Int64("uptime_seconds", livenessResponse.UptimeSeconds).Msg("Service is alive")
	}

	return c.Status(statusCode).JSON(livenessResponse)
}

// RegisterRoutes registers health-related routes
func (h *HealthHandler) RegisterRoutes(router fiber.Router) {
	h.logger.Info().Msg("Registering health routes")
	router.Get("/health", h.GetHealth)
	router.Get("/liveness", h.GetLiveness)
	h.logger.Debug().Str("route", "/health").Msg("Health route registered")
	h.logger.Debug().Str("route", "/liveness").Msg("Liveness route registered")
}
