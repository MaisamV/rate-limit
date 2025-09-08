package ports

// SwaggerProvider defines the interface for documentation operations
type SwaggerProvider interface {

	// GetOpenAPISpec loads the OpenAPI specification from file
	GetOpenAPISpec() ([]byte, error)

	// GetSwaggerHTML generates the Swagger UI HTML
	GetSwaggerHTML() ([]byte, error)
}
