package localconfig

import (
	"fmt"

	"backend.atomicledger.com/pkg/logger"
)

// ConfigService wraps LocalConfig and provides convenient access methods.
type ConfigService struct {
	config LocalConfig
	logger logger.Logger
}

// NewConfigService creates a new config service by loading configuration.
func NewConfigService(log logger.Logger) (*ConfigService, error) {
	config, err := GetConfig(log)
	if err != nil {
		return nil, err
	}

	return &ConfigService{
		config: config,
		logger: log,
	}, nil
}

// Get returns the entire configuration.
func (s *ConfigService) Get() LocalConfig {
	return s.config
}

// GetConnectionString builds and returns the database connection string.
func (s *ConfigService) GetConnectionString() string {
	db := s.config.Database
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host,
		db.Port,
		db.Username,
		db.Password,
		db.Name,
		db.SSLMode,
	)
}

// GetServicePort returns the HTTP server port.
func (s *ConfigService) GetServicePort() int {
	return s.config.Service.Port
}

// GetServiceName returns the service name.
func (s *ConfigService) GetServiceName() string {
	return s.config.Service.Name
}
