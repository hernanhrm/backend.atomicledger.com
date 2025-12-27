package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"backend.atomicledger.com/pkg/database"
	"backend.atomicledger.com/pkg/localconfig"
	"backend.atomicledger.com/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

// Server represents the HTTP server with all its dependencies.
type Server struct {
	config   *localconfig.ConfigService
	db       *database.Database
	logger   logger.Logger
	injector do.Injector

	Echo *echo.Echo
}

// NewServer creates a new HTTP server instance with all dependencies injected.
func NewServer(
	config *localconfig.ConfigService,
	db *database.Database,
	log logger.Logger,
	injector do.Injector,
	routeSetupFunc func(*Server),
) (*Server, error) {
	s := &Server{
		config:   config,
		db:       db,
		logger:   log.With("component", "server"),
		injector: injector,
		Echo:     echo.New(),
	}

	// Configure Echo
	s.Echo.HideBanner = true
	s.Echo.HidePort = true

	// Setup middleware
	s.setupMiddleware()

	// Setup routes via callback
	if routeSetupFunc != nil {
		routeSetupFunc(s)
	}

	s.logger.Info("HTTP server initialized",
		"port", config.GetServicePort(),
		"service", config.GetServiceName(),
	)

	return s, nil
}

// setupMiddleware configures all HTTP middleware.
func (s *Server) setupMiddleware() {
	// Recover from panics
	s.Echo.Use(middleware.Recover())

	// Request logging
	s.Echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		LogLatency:  true,
		HandleError: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				s.logger.Info("request completed",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
				)
			} else {
				s.logger.Error("request failed",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"error", v.Error,
				)
			}
			return nil
		},
	}))

	// CORS (configure as needed)
	s.Echo.Use(middleware.CORS())
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.config.GetServicePort())

	s.logger.Info("starting HTTP server",
		"address", addr,
		"service", s.config.GetServiceName(),
	)

	// Start server in a goroutine
	go func() {
		if err := s.Echo.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server error", "error", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	s.logger.Info("shutting down HTTP server")

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Echo.Shutdown(shutdownCtx); err != nil {
		return oops.
			Code("server_shutdown_failed").
			Wrapf(err, "failed to shutdown HTTP server")
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")

	if err := s.Echo.Shutdown(ctx); err != nil {
		return oops.
			Code("server_shutdown_failed").
			Wrapf(err, "failed to shutdown HTTP server")
	}

	return nil
}

// HealthCheck performs a health check on the server.
func (s *Server) HealthCheck(_ context.Context) error {
	// Server is healthy if it's running
	return nil
}
