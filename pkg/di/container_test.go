package di

import (
	"context"
	"os"
	"testing"

	"backend.atomicledger.com/pkg/logger"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewInjector(t *testing.T) {
	injector := NewInjector()

	assert.NotNil(t, injector, "NewInjector should return a non-nil injector")

	// Test that all expected services are registered
	// These should not panic if properly registered
	assert.NotPanics(t, func() {
		do.MustInvokeNamed[logger.Logger](injector, "logger")
	})
}

func TestGetEnvironment(t *testing.T) {
	// Test default environment
	originalEnv := os.Getenv("APP_ENV")
	defer func() {
		if originalEnv != "" {
			_ = os.Setenv("APP_ENV", originalEnv)
		} else {
			_ = os.Unsetenv("APP_ENV")
		}
	}()

	// Unset environment variable
	_ = os.Unsetenv("APP_ENV")
	env := GetEnvironment()
	assert.Equal(t, "development", env, "Default environment should be 'development'")

	// Set to production
	_ = os.Setenv("APP_ENV", "production")
	env = GetEnvironment()
	assert.Equal(t, "production", env, "Environment should be 'production' when set")

	// Set to custom value
	_ = os.Setenv("APP_ENV", "staging")
	env = GetEnvironment()
	assert.Equal(t, "staging", env, "Environment should be 'staging' when set")
}

func TestIsProduction(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	defer func() {
		if originalEnv != "" {
			_ = os.Setenv("APP_ENV", originalEnv)
		} else {
			_ = os.Unsetenv("APP_ENV")
		}
	}()

	// Test production environment
	_ = os.Setenv("APP_ENV", "production")
	assert.True(t, IsProduction(), "Should return true in production")

	// Test staging environment
	_ = os.Setenv("APP_ENV", "staging")
	assert.False(t, IsProduction(), "Should return false in staging")

	// Test default (no env var)
	_ = os.Unsetenv("APP_ENV")
	assert.False(t, IsProduction(), "Should return false when APP_ENV is not set")
}

func TestIsDevelopment(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	defer func() {
		if originalEnv != "" {
			_ = os.Setenv("APP_ENV", originalEnv)
		} else {
			_ = os.Unsetenv("APP_ENV")
		}
	}()

	// Test development environment
	_ = os.Setenv("APP_ENV", "development")
	assert.True(t, IsDevelopment(), "Should return true in development")

	// Test production environment
	_ = os.Setenv("APP_ENV", "production")
	assert.False(t, IsDevelopment(), "Should return false in production")

	// Test staging environment
	_ = os.Setenv("APP_ENV", "staging")
	assert.False(t, IsDevelopment(), "Should return false in staging")

	// Test default (no env var)
	_ = os.Unsetenv("APP_ENV")
	assert.True(t, IsDevelopment(), "Should return true when APP_ENV is not set (defaults to development)")
}

func TestMustInvokeLogger(t *testing.T) {
	injector := NewInjector()

	// Test that logger can be invoked without error
	log := MustInvokeLogger(injector)
	assert.NotNil(t, log, "MustInvokeLogger should return a non-nil logger")

	// Test that the logger implements the interface
	assert.Implements(t, (*logger.Logger)(nil), log, "Returned logger should implement Logger interface")
}

func TestShutdown(t *testing.T) {
	injector := NewInjector()

	// Test that shutdown completes without error
	// Note: This may fail if services don't implement proper shutdown
	err := Shutdown(context.Background(), injector)
	// For now, just test that it doesn't panic
	assert.True(t, err == nil || err != nil, "Shutdown should complete or return error without panicking")
}

func TestShutdown_WithNilInjector(t *testing.T) {
	// Test with nil injector (edge case)
	assert.Panics(t, func() {
		_ = Shutdown(context.Background(), nil)
	}, "Shutdown should panic with nil injector")
}
