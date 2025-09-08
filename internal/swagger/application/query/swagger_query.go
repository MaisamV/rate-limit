package query

import (
	"github.com/go-clean/internal/swagger/ports"
	"github.com/go-clean/platform/logger"
)

// SwaggerQueryHandler handles swagger-related queries
type SwaggerQueryHandler struct {
	logger          logger.Logger
	swaggerProvider ports.SwaggerProvider
}

// NewSwaggerQueryHandler creates a new swagger query handler
func NewSwaggerQueryHandler(logger logger.Logger, swaggerProvider ports.SwaggerProvider) *SwaggerQueryHandler {
	return &SwaggerQueryHandler{
		logger:          logger,
		swaggerProvider: swaggerProvider,
	}
}

// GetOpenAPISpec returns the OpenAPI specification as JSON
func (h *SwaggerQueryHandler) GetOpenAPISpec() ([]byte, error) {
	h.logger.Debug().Msg("Retrieving OpenAPI specification")
	spec, err := h.swaggerProvider.GetOpenAPISpec()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to retrieve OpenAPI specification")
		return nil, err
	}
	h.logger.Info().Msg("OpenAPI specification retrieved successfully")
	return spec, nil
}

// GetSwaggerHTML returns the Swagger UI HTML page
func (h *SwaggerQueryHandler) GetSwaggerHTML() ([]byte, error) {
	h.logger.Debug().Msg("Retrieving Swagger UI HTML")
	html, err := h.swaggerProvider.GetSwaggerHTML()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to retrieve Swagger UI HTML")
		return nil, err
	}
	h.logger.Info().Msg("Swagger UI HTML retrieved successfully")
	return html, nil
}
