// @title Go Clean Architecture API
// @version 1.0
// @description A clean architecture implementation in Go with comprehensive API documentation
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize application with wire-generated dependency injection
	app, err := InitializeApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	app.Logger.Info().Msg("Starting application")
	app.Logger.Info().Msg("All modules initialized successfully")

	// Get the fiber app instance
	fiberApp := app.HTTPServer.GetApp()

	// Register routes
	app.Logger.Info().Msg("Registering routes")
	app.Probes.PingHandler.RegisterRoutes(fiberApp)
	app.Probes.HealthHandler.RegisterRoutes(fiberApp)
	app.RateLimit.RateLimitHandler.RegisterRoutes(fiberApp)
	app.Swagger.DocsHandler.RegisterRoutes(fiberApp, app.Config.Swagger.Enabled)
	app.Logger.Info().Msg("Routes registered successfully")

	// Start server
	app.Logger.Info().Str("port", app.Config.Server.Port).Msg("Starting HTTP server")
	go func() {
		if err := app.HTTPServer.Start(); err != nil {
			app.Logger.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	app.Logger.Info().Msg("Shutting down server...")

	// Gracefully shutdown the server
	if err := app.HTTPServer.Shutdown(); err != nil {
		app.Logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	app.Logger.Info().Msg("Server exited")
}
