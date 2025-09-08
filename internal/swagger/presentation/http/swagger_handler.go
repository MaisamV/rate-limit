package http

import (
	"github.com/go-clean/internal/swagger/application/query"
	"github.com/go-clean/platform/logger"
	"github.com/gofiber/fiber/v2"
)

// DocsHandler handles HTTP requests for Swagger documentation
type DocsHandler struct {
	logger              logger.Logger
	swaggerQueryHandler *query.SwaggerQueryHandler
}

// NewDocsHandler creates a new docs handler
func NewDocsHandler(logger logger.Logger, swaggerQueryHandler *query.SwaggerQueryHandler) *DocsHandler {
	return &DocsHandler{
		logger:              logger,
		swaggerQueryHandler: swaggerQueryHandler,
	}
}

// GetOpenAPISpec handles GET /api/docs/openapi.yaml
// @Summary Get OpenAPI Specification
// @Description Returns the OpenAPI specification in YAML format
// @Tags Documentation
// @Produce text/plain
// @Success 200 {string} string "OpenAPI specification in YAML format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/docs/openapi.yaml [get]
func (h *DocsHandler) GetOpenAPISpec(c *fiber.Ctx) error {
	h.logger.Info().Str("endpoint", "/openapi.yaml").Msg("OpenAPI specification requested")
	spec, err := h.swaggerQueryHandler.GetOpenAPISpec()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to serve OpenAPI specification")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load OpenAPI specification",
		})
	}

	c.Set("Content-Type", "text/yaml")
	h.logger.Debug().Int("size_bytes", len(spec)).Msg("OpenAPI specification served successfully")
	return c.Send(spec)
}

// GetSwaggerUI handles GET /api/docs
// @Summary Get Swagger UI
// @Description Returns the Swagger UI HTML page for API documentation
// @Tags Documentation
// @Produce text/html
// @Success 200 {string} string "Swagger UI HTML page"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/docs [get]
func (h *DocsHandler) GetSwaggerUI(c *fiber.Ctx) error {
	h.logger.Info().Str("endpoint", "/swagger").Msg("Swagger UI requested")
	html, err := h.swaggerQueryHandler.GetSwaggerHTML()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to serve Swagger UI")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate Swagger UI",
		})
	}

	c.Set("Content-Type", "text/html")
	h.logger.Debug().Int("size_bytes", len(html)).Msg("Swagger UI served successfully")
	return c.Send(html)
}

// RegisterRoutes registers the documentation routes
func (h *DocsHandler) RegisterRoutes(app *fiber.App, enabled bool) {
	if !enabled {
		h.logger.Info().Msg("Swagger documentation disabled, skipping route registration")
		return
	}
	h.logger.Info().Msg("Registering Swagger documentation routes")
	app.Get("/swagger", h.GetSwaggerUI)
	app.Get("/openapi.yaml", h.GetOpenAPISpec)
	h.logger.Info().Msg("Swagger documentation routes registered successfully")
}
