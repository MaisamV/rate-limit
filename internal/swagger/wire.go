package swagger

import (
	swaggerQuery "github.com/go-clean/internal/swagger/application/query"
	"github.com/go-clean/internal/swagger/infrastructure"
	swaggerHttp "github.com/go-clean/internal/swagger/presentation/http"
	"github.com/go-clean/platform/logger"
	"github.com/google/wire"
)

// ProvideSwaggerConfig provides swagger configuration
func ProvideSwaggerConfig() infrastructure.SwaggerConfig {
	return infrastructure.SwaggerConfig{
		OpenApiFilePath: "./api/openapi.yaml",
		SwaggerFilePath: "./api/swagger.html",
	}
}

// ProvideSwaggerLoader provides a swagger loader
func ProvideSwaggerLoader(logger logger.Logger, config infrastructure.SwaggerConfig) (*infrastructure.SwaggerLoader, error) {
	loader := infrastructure.NewSwaggerLoader(logger, config)
	if err := loader.Init(); err != nil {
		return nil, err
	}
	return loader, nil
}

// ProvideSwaggerQueryHandler provides a swagger query handler
func ProvideSwaggerQueryHandler(logger logger.Logger, swaggerLoader *infrastructure.SwaggerLoader) *swaggerQuery.SwaggerQueryHandler {
	return swaggerQuery.NewSwaggerQueryHandler(logger, swaggerLoader)
}

// ProvideDocsHandler provides a docs HTTP handler
func ProvideDocsHandler(logger logger.Logger, swaggerQueryHandler *swaggerQuery.SwaggerQueryHandler) *swaggerHttp.DocsHandler {
	return swaggerHttp.NewDocsHandler(logger, swaggerQueryHandler)
}

// SwaggerSet is a wire provider set for all swagger dependencies
var SwaggerSet = wire.NewSet(
	ProvideSwaggerConfig,
	ProvideSwaggerLoader,
	ProvideSwaggerQueryHandler,
	ProvideDocsHandler,
)
