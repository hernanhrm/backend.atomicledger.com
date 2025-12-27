// Package localconfig provides configuration loading from environment variables and .env files.
package localconfig

// ConfigOptions defines options for loading configuration.
type ConfigOptions struct {
	EnvPath     string // Relative path to directory containing the env file (default: ".")
	EnvFileName string // Name of the env file (default: ".env")
}

// LocalConfig holds the complete application configuration.
type LocalConfig struct {
	Service  Service
	Database Database
}

// Service holds service-specific configuration.
type Service struct {
	Port int
	Name string
}

// Database holds database connection configuration.
type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	SSLMode  string
}
