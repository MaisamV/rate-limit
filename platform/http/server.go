package http

import (
	"time"

	"github.com/go-clean/platform/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Server represents the HTTP server configuration
type Server struct {
	app    *fiber.App
	port   string
	logger logger.Logger
}

// NewServer creates a new HTTP server with common middleware
func NewServer(port string, log logger.Logger) *Server {
	log.Info().Str("port", port).Msg("Initializing HTTP server")

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return errorHandler(c, err, log)
		},
	})

	// Add common middleware
	log.Debug().Msg("Configuring HTTP server middleware")
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "${time} ${status} - ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	log.Info().Msg("HTTP server initialized successfully")
	return &Server{
		app:    app,
		port:   port,
		logger: log,
	}
}

// GetApp returns the fiber app instance for route registration
func (s *Server) GetApp() *fiber.App {
	return s.app
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info().Str("port", s.port).Msg("Starting HTTP server")
	err := s.app.Listen(":" + s.port)
	if err != nil {
		s.logger.Error().Err(err).Str("port", s.port).Msg("Failed to start HTTP server")
	}
	return err
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.logger.Info().Msg("Shutting down HTTP server")
	err := s.app.Shutdown()
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to shutdown HTTP server gracefully")
	} else {
		s.logger.Info().Msg("HTTP server shutdown completed")
	}
	return err
}

// errorHandler handles fiber errors
func errorHandler(c *fiber.Ctx, err error, log logger.Logger) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Error().Err(err).Int("status_code", code).Str("method", c.Method()).Str("path", c.Path()).Msg("HTTP request error")

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}
