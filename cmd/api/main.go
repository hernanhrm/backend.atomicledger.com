// Package main provides the entry point for the backend API application.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend.atomicledger.com/pkg/database"
	"backend.atomicledger.com/pkg/di"
	"backend.atomicledger.com/pkg/localconfig"
	"backend.atomicledger.com/pkg/server"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

func main() {
	if err := run(); err != nil {
		// Extract oops error details if available
		if oopsErr, ok := oops.AsOops(err); ok {
			fmt.Fprintf(os.Stderr, "Fatal error: %s\n", oopsErr.Error())
			fmt.Fprintf(os.Stderr, "Error code: %s\n", oopsErr.Code())
			fmt.Fprintf(os.Stderr, "Details: %+v\n", oopsErr.Context())
		} else {
			fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		}
		os.Exit(1)
	}
}

func run() error {
	// Create DI injector with services (excluding server)
	injector := do.New()

	// Register providers manually
	di.ProvideLogger()(injector)
	di.ProvideConfig()(injector)
	di.ProvideDatabase()(injector)

	// Get dependencies
	logSvc := di.MustInvokeLogger(injector)

	configSvc, err := do.InvokeNamed[*localconfig.ConfigService](injector, "config")
	if err != nil {
		return oops.Wrapf(err, "failed to get config")
	}

	dbSvc, err := do.InvokeNamed[*database.Database](injector, "database")
	if err != nil {
		return oops.Wrapf(err, "failed to get database")
	}

	logSvc.Info("application starting",
		"environment", di.GetEnvironment(),
		"version", "0.1.0",
	)

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Define route setup function
	setupRoutes := func(s *server.Server) {
		// Health check endpoint
		s.Echo.GET("/health", s.HandleHealth)

		// API routes group
		api := s.Echo.Group("/api")
		api.GET("/ping", s.HandlePing)

		// Add more routes here as needed
	}

	// Create server with route setup
	srv, err := server.NewServer(configSvc, dbSvc, logSvc, injector, setupRoutes)
	if err != nil {
		return oops.
			Code("server_initialization_failed").
			Wrapf(err, "failed to initialize HTTP server")
	}

	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	logSvc.Info("application ready")

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		logSvc.Info("shutdown signal received")
	case err := <-errChan:
		return oops.
			Code("server_runtime_error").
			Wrapf(err, "server encountered a runtime error")
	}

	// Graceful shutdown with timeout
	logSvc.Info("initiating graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logSvc.Error("error during shutdown", "error", err)
		return oops.
			Code("shutdown_error").
			Wrapf(err, "failed to shutdown gracefully")
	}

	logSvc.Info("application stopped gracefully")
	return nil
}
