// Package di provides dependency injection container and service management.
package di

import (
	"context"
	"os"
	"time"

	"backend.atomicledger.com/pkg/logger"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

// NewInjector creates and configures the root DI container with all application services.
func NewInjector() do.Injector {
	injector := do.New()

	// Register all services using the provider functions
	// Each provider function returns a func(do.Injector) that registers the service
	ProvideLogger()(injector)
	ProvideConfig()(injector)
	ProvideDatabase()(injector)

	return injector
}

// Shutdown gracefully shuts down all services in the injector
// Services are shut down in reverse dependency order.
func Shutdown(_ context.Context, injector do.Injector) error {
	log := MustInvokeLogger(injector)
	log.Info("initiating graceful shutdown")

	startTime := time.Now()

	// Shutdown the injector (this calls Shutdown on all services that implement it)
	if err := injector.Shutdown(); err != nil {
		log.Error("error during shutdown",
			"error", err,
			"duration", time.Since(startTime),
		)
		return oops.
			Code("shutdown_failed").
			With("duration", time.Since(startTime)).
			Wrapf(err, "failed to shutdown services")
	}

	log.Info("graceful shutdown completed",
		"duration", time.Since(startTime),
	)

	return nil
}

// GetEnvironment returns the current application environment (production, development, etc.)
func GetEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "development"
	}
	return env
}

// IsProduction returns true if running in production environment.
func IsProduction() bool {
	return GetEnvironment() == "production"
}

// IsDevelopment returns true if running in development environment.
func IsDevelopment() bool {
	return GetEnvironment() == "development"
}

// MustInvokeLogger is a helper to get the logger from the injector without error handling
// This is safe to use because logger is always provided and never fails.
func MustInvokeLogger(injector do.Injector) logger.Logger {
	return do.MustInvokeNamed[logger.Logger](injector, "logger")
}
