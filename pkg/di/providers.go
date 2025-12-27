package di

import (
	"context"

	"backend.atomicledger.com/pkg/database"
	"backend.atomicledger.com/pkg/localconfig"
	"backend.atomicledger.com/pkg/logger"
	"github.com/samber/do/v2"
)

// ProvideLogger returns a function that registers the logger service
// Logger is created based on APP_ENV environment variable.
func ProvideLogger() func(do.Injector) {
	return do.LazyNamed("logger", func(_ do.Injector) (logger.Logger, error) {
		if IsProduction() {
			return logger.NewProduction(), nil
		}
		return logger.NewDevelopment(), nil
	})
}

// ProvideConfig returns a function that registers the configuration service
// Config is loaded eagerly and will cause startup failure if invalid.
func ProvideConfig() func(do.Injector) {
	return do.LazyNamed("config", func(i do.Injector) (*localconfig.ConfigService, error) {
		log, err := do.InvokeNamed[logger.Logger](i, "logger")
		if err != nil {
			return nil, err
		}

		return localconfig.NewConfigService(log)
	})
}

// ProvideDatabase returns a function that registers the database service
// Database connection is lazy-loaded and implements graceful shutdown.
func ProvideDatabase() func(do.Injector) {
	return ProvideNamedDatabase("database")
}

// ProvideNamedDatabase returns a function that registers a named database service
// This allows multiple database connections (e.g., primary, analytics, cache).
func ProvideNamedDatabase(name string) func(do.Injector) {
	return do.LazyNamed(name, func(i do.Injector) (*database.Database, error) {
		log, err := do.InvokeNamed[logger.Logger](i, "logger")
		if err != nil {
			return nil, err
		}

		config, err := do.InvokeNamed[*localconfig.ConfigService](i, "config")
		if err != nil {
			return nil, err
		}

		ctx := context.Background()

		return database.NewConnection(ctx, config.GetConnectionString(), log)
	})
}
