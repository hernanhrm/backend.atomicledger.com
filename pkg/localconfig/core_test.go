package localconfig

import (
	"os"
	"path/filepath"
	"testing"

	"backend.atomicledger.com/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	// Create a temporary .env file for testing
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	envContent := `SERVICE_PORT=8080
SERVICE_NAME=test-api
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=testuser
DATABASE_PASSWORD=testpass
DATABASE_NAME=testdb
DATABASE_SSL_MODE=disable
`

	if err := os.WriteFile(envFile, []byte(envContent), 0o600); err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Set environment variables to point to temp directory
	if err := os.Setenv("CONFIG_ENV_PATH", tmpDir); err != nil {
		t.Fatalf("Failed to set CONFIG_ENV_PATH: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("CONFIG_ENV_PATH"); err != nil {
			t.Errorf("Failed to unset CONFIG_ENV_PATH: %v", err)
		}
	}()

	log := logger.NewNoop()
	config, err := GetConfig(log)
	assert.NoError(t, err, "GetConfig should not fail")

	// Verify Service config
	assert.Equal(t, 8080, config.Service.Port, "Service.Port should be 8080")
	assert.Equal(t, "test-api", config.Service.Name, "Service.Name should be 'test-api'")

	// Verify Database config
	assert.Equal(t, "localhost", config.Database.Host, "Database.Host should be 'localhost'")
	assert.Equal(t, 5432, config.Database.Port, "Database.Port should be 5432")
	assert.Equal(t, "testuser", config.Database.Username, "Database.Username should be 'testuser'")
	assert.Equal(t, "testdb", config.Database.Name, "Database.Name should be 'testdb'")
	assert.Equal(t, "disable", config.Database.SSLMode, "Database.SSLMode should be 'disable'")
}

func TestGetConfigWithOptions(t *testing.T) {
	// Create a temporary directory with a custom .env file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env.custom")

	envContent := `SERVICE_PORT=9090
SERVICE_NAME=custom-service
DATABASE_HOST=customhost
DATABASE_PORT=3306
DATABASE_USERNAME=customuser
DATABASE_PASSWORD=custompass
DATABASE_NAME=customdb
DATABASE_SSL_MODE=require
`

	if err := os.WriteFile(envFile, []byte(envContent), 0o600); err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Use GetConfigWithOptions to specify custom path and filename
	log := logger.NewNoop()
	config, err := GetConfigWithOptions(ConfigOptions{
		EnvPath:     tmpDir,
		EnvFileName: ".env.custom",
	}, log)
	assert.NoError(t, err, "GetConfigWithOptions should not fail")

	// Verify Service config
	assert.Equal(t, 9090, config.Service.Port, "Service.Port should be 9090")
	assert.Equal(t, "custom-service", config.Service.Name, "Service.Name should be 'custom-service'")

	// Verify Database config
	assert.Equal(t, "customhost", config.Database.Host, "Database.Host should be 'customhost'")
	assert.Equal(t, 3306, config.Database.Port, "Database.Port should be 3306")
	assert.Equal(t, "require", config.Database.SSLMode, "Database.SSLMode should be 'require'")
}

func TestGetConfigWithEnvVars(t *testing.T) {
	// Create a temporary directory with .env.dev file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env.dev")

	envContent := `SERVICE_PORT=7070
SERVICE_NAME=dev-service
DATABASE_HOST=devhost
DATABASE_PORT=5433
DATABASE_USERNAME=devuser
DATABASE_PASSWORD=devpass
DATABASE_NAME=devdb
DATABASE_SSL_MODE=disable
`

	if err := os.WriteFile(envFile, []byte(envContent), 0o600); err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Set environment variables
	if err := os.Setenv("CONFIG_ENV_PATH", tmpDir); err != nil {
		t.Fatalf("Failed to set CONFIG_ENV_PATH: %v", err)
	}
	if err := os.Setenv("CONFIG_ENV_FILENAME", ".env.dev"); err != nil {
		t.Fatalf("Failed to set CONFIG_ENV_FILENAME: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("CONFIG_ENV_PATH"); err != nil {
			t.Errorf("Failed to unset CONFIG_ENV_PATH: %v", err)
		}
	}()
	defer func() {
		if err := os.Unsetenv("CONFIG_ENV_FILENAME"); err != nil {
			t.Errorf("Failed to unset CONFIG_ENV_FILENAME: %v", err)
		}
	}()

	log := logger.NewNoop()
	config, err := GetConfig(log)
	assert.NoError(t, err, "GetConfig should not fail")

	// Verify Service config
	assert.Equal(t, 7070, config.Service.Port, "Service.Port should be 7070")
	assert.Equal(t, "dev-service", config.Service.Name, "Service.Name should be 'dev-service'")

	// Verify Database config
	assert.Equal(t, "devhost", config.Database.Host, "Database.Host should be 'devhost'")
}

func TestGetConfigFallbackToEnvVars(t *testing.T) {
	// Unset any existing CONFIG_ENV_* vars to ensure clean test
	if err := os.Unsetenv("CONFIG_ENV_PATH"); err != nil {
		t.Fatalf("Failed to unset CONFIG_ENV_PATH: %v", err)
	}
	if err := os.Unsetenv("CONFIG_ENV_FILENAME"); err != nil {
		t.Fatalf("Failed to unset CONFIG_ENV_FILENAME: %v", err)
	}

	// Set system environment variables (simulating production)
	// Viper converts env var names to lowercase, so we use uppercase names
	// that match the mapstructure tags (service_port, db_host, etc.)
	envVars := map[string]string{
		"SERVICE_PORT":      "3000",
		"SERVICE_NAME":      "prod-service",
		"DATABASE_HOST":     "prod.example.com",
		"DATABASE_PORT":     "5432",
		"DATABASE_USERNAME": "produser",
		"DATABASE_PASSWORD": "prodpass",
		"DATABASE_NAME":     "proddb",
		"DATABASE_SSL_MODE": "verify-full",
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set %s: %v", key, err)
		}
	}

	defer func() {
		for key := range envVars {
			if err := os.Unsetenv(key); err != nil {
				t.Errorf("Failed to unset %s: %v", key, err)
			}
		}
	}()

	// Point to a non-existent directory (simulating production without .env file)
	tmpDir := t.TempDir()
	log := logger.NewNoop()
	config, err := GetConfigWithOptions(ConfigOptions{
		EnvPath:     tmpDir,
		EnvFileName: ".env.nonexistent",
	}, log)
	assert.NoError(t, err, "GetConfigWithOptions should not fail")

	// Should read from system environment variables
	assert.Equal(t, 3000, config.Service.Port, "Service.Port should be 3000")
	assert.Equal(t, "prod-service", config.Service.Name, "Service.Name should be 'prod-service'")
	assert.Equal(t, "prod.example.com", config.Database.Host, "Database.Host should be 'prod.example.com'")
	assert.Equal(t, "verify-full", config.Database.SSLMode, "Database.SSLMode should be 'verify-full'")
}
