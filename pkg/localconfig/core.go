package localconfig

import (
	"os"
	"path/filepath"
	"strconv"

	"backend.atomicledger.com/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/samber/oops"
)

// GetConfig loads configuration from environment variables or .env file.
// It reads CONFIG_ENV_PATH and CONFIG_ENV_FILENAME environment variables to determine
// the location and name of the config file.
//
// Environment variables:
//   - CONFIG_ENV_PATH: Relative path to the directory containing the env file (default: ".")
//   - CONFIG_ENV_FILENAME: Name of the env file (default: ".env")
//
// If the config file is not found, it falls back to reading from system environment variables.
// This is useful for production environments where variables are already loaded.
func GetConfig(log logger.Logger) (LocalConfig, error) {
	opts := ConfigOptions{
		EnvPath:     getEnvOrDefault("CONFIG_ENV_PATH", "."),
		EnvFileName: getEnvOrDefault("CONFIG_ENV_FILENAME", ".env"),
	}
	return GetConfigWithOptions(opts, log)
}

// GetConfigWithOptions loads configuration using the provided options.
// This allows programmatic control over the config file location and name.
//
// If the config file is not found, it falls back to reading from system environment variables.
func GetConfigWithOptions(opts ConfigOptions, log logger.Logger) (LocalConfig, error) {
	// Construct the full path to the .env file
	envFilePath := filepath.Join(opts.EnvPath, opts.EnvFileName)

	// Try to load the .env file, but don't fail if it doesn't exist
	// This allows the app to work in production with only env vars
	// godotenv.Overload ensures new values override existing env vars
	if err := godotenv.Overload(envFilePath); err != nil {
		log.Warn("config file not found, using environment variables only",
			"path", envFilePath,
			"error", err.Error(),
		)
	} else {
		log.Info("config file loaded successfully",
			"path", envFilePath,
		)
	}

	// Manually parse environment variables and populate config struct
	servicePort, err := getEnvAsInt("SERVICE_PORT", log)
	if err != nil {
		return LocalConfig{}, oops.
			Code("invalid_config_value").
			With("key", "SERVICE_PORT").
			Wrapf(err, "failed to parse SERVICE_PORT as integer")
	}

	databasePort, err := getEnvAsInt("DATABASE_PORT", log)
	if err != nil {
		return LocalConfig{}, oops.
			Code("invalid_config_value").
			With("key", "DATABASE_PORT").
			Wrapf(err, "failed to parse DATABASE_PORT as integer")
	}

	config := LocalConfig{
		Service: Service{
			Port: servicePort,
			Name: getEnvAsString("SERVICE_NAME"),
		},
		Database: Database{
			Host:     getEnvAsString("DATABASE_HOST"),
			Port:     databasePort,
			Username: getEnvAsString("DATABASE_USERNAME"),
			Password: getEnvAsString("DATABASE_PASSWORD"),
			Name:     getEnvAsString("DATABASE_NAME"),
			SSLMode:  getEnvAsString("DATABASE_SSL_MODE"),
		},
	}

	log.Debug("configuration loaded",
		"service_name", config.Service.Name,
		"service_port", config.Service.Port,
		"database_host", config.Database.Host,
		"database_port", config.Database.Port,
	)

	return config, nil
}

// getEnvOrDefault returns the value of the environment variable or a default value if not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsString returns the environment variable value as a string.
func getEnvAsString(key string) string {
	return os.Getenv(key)
}

// getEnvAsInt returns the environment variable value as an int.
// Returns 0 if the environment variable is not set (default value).
// Returns error if the value exists but is not a valid integer.
func getEnvAsInt(key string, log logger.Logger) (int, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		log.Debug("environment variable not set, using default value 0",
			"key", key,
		)
		return 0, nil // Return 0 if not set (default value)
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Error("invalid integer value for environment variable",
			"key", key,
			"value", valueStr,
			"error", err,
		)
		return 0, oops.
			Code("invalid_env_int").
			With("key", key).
			With("value", valueStr).
			Wrapf(err, "environment variable %s has invalid integer value", key)
	}

	return value, nil
}
