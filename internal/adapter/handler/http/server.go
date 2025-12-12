package http

import (
	"ishari-backend/pkg/config"
	"ishari-backend/pkg/errors"
	"ishari-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
	cfg config.ServerConfig
}

// NewServer creates a new HTTP server instance
func NewServer(cfg config.ServerConfig, l logger.Logger) *Server {
	app := fiber.New(fiber.Config{
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		ErrorHandler: errors.HandlerWithLogger(l),
	})

	return &Server{
		App: app,
		cfg: cfg,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	return s.App.Listen(addr)
}
