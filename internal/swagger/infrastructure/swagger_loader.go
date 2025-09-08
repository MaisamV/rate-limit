package infrastructure

import (
	"fmt"
	"os"

	"github.com/go-clean/platform/logger"
)

// SwaggerLoader implements the ports.SwaggerProvider interface
type SwaggerLoader struct {
	logger      logger.Logger
	config      SwaggerConfig
	openapi     []byte
	swaggerHtml []byte
}

// SwaggerConfig holds configuration for swagger loader
type SwaggerConfig struct {
	OpenApiFilePath string
	SwaggerFilePath string
}

// NewSwaggerLoader creates a new instance of DocsAdapter
func NewSwaggerLoader(logger logger.Logger, config SwaggerConfig) *SwaggerLoader {
	return &SwaggerLoader{
		logger: logger,
		config: config,
	}
}

func (a *SwaggerLoader) Init() error {
	a.logger.Info().Str("openapi_path", a.config.OpenApiFilePath).Msg("Loading OpenAPI specification")
	data, err := os.ReadFile(a.config.OpenApiFilePath)
	if err != nil {
		a.logger.Error().Err(err).Str("path", a.config.OpenApiFilePath).Msg("Failed to read OpenAPI spec")
		return fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}
	a.openapi = data
	a.logger.Info().Int("size_bytes", len(data)).Msg("OpenAPI specification loaded successfully")

	a.logger.Info().Str("swagger_path", a.config.SwaggerFilePath).Msg("Loading Swagger HTML")
	html, err := os.ReadFile(a.config.SwaggerFilePath)
	if err != nil {
		a.logger.Error().Err(err).Str("path", a.config.SwaggerFilePath).Msg("Failed to read Swagger HTML")
		return fmt.Errorf("failed to read Swagger HTML: %w", err)
	}
	a.swaggerHtml = html
	a.logger.Info().Int("size_bytes", len(html)).Msg("Swagger HTML loaded successfully")
	return nil
}

// GetOpenAPISpec loads the OpenAPI specification from file
func (a *SwaggerLoader) GetOpenAPISpec() ([]byte, error) {
	a.logger.Debug().Msg("Serving OpenAPI specification")
	return a.openapi, nil
}

// GetSwaggerHTML generates the Swagger UI HTML
func (a *SwaggerLoader) GetSwaggerHTML() ([]byte, error) {
	a.logger.Debug().Msg("Serving Swagger UI HTML")
	return a.swaggerHtml, nil
}
