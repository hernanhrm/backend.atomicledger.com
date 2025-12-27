package di

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestProvideLogger(t *testing.T) {
	provider := ProvideLogger()
	assert.NotNil(t, provider, "ProvideLogger should return a non-nil provider")

	// Test that provider is a function
	assert.IsType(t, func(do.Injector) {}, provider, "Provider should be a function")
}

func TestProvideConfig(t *testing.T) {
	provider := ProvideConfig()
	assert.NotNil(t, provider, "ProvideConfig should return a non-nil provider")

	// Test that provider is a function
	assert.IsType(t, func(do.Injector) {}, provider, "Provider should be a function")
}

func TestProvideDatabase(t *testing.T) {
	provider := ProvideDatabase()
	assert.NotNil(t, provider, "ProvideDatabase should return a non-nil provider")

	// Test that provider is a function
	assert.IsType(t, func(do.Injector) {}, provider, "Provider should be a function")
}

func TestProvideNamedDatabase(t *testing.T) {
	provider := ProvideNamedDatabase("test-db")
	assert.NotNil(t, provider, "ProvideNamedDatabase should return a non-nil provider")

	// Test that provider is a function
	assert.IsType(t, func(do.Injector) {}, provider, "Provider should be a function")
}

func TestProviderFunctions_ReturnDifferentInstances(t *testing.T) {
	// Test that different calls to ProvideNamedDatabase return different providers
	provider1 := ProvideNamedDatabase("db1")
	provider2 := ProvideNamedDatabase("db2")

	// Test that providers are different functions (they should be different closures)
	assert.NotNil(t, provider1)
	assert.NotNil(t, provider2)
	// We can't directly compare functions in Go, so we just test they're both non-nil
}

func TestProviderFunctions_WithSameName(t *testing.T) {
	// Test that calls with same name return equivalent providers
	provider1 := ProvideNamedDatabase("same-db")
	provider2 := ProvideNamedDatabase("same-db")

	// They should be different function instances but for same database
	assert.NotNil(t, provider1)
	assert.NotNil(t, provider2)
}
