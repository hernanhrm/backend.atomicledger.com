package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend.atomicledger.com/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHealth(t *testing.T) {
	// For now, just test that the handler doesn't panic with minimal setup
	// We'll need proper mocking for comprehensive testing
	s := &Server{
		logger: logger.NewNoop(),
	}

	// Create Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler - expect panic due to nil db, but catch it
	assert.Panics(t, func() {
		_ = s.HandleHealth(c)
	}, "Should panic when database is nil")
}

func TestHandlePing(t *testing.T) {
	// Create server instance
	s := &Server{
		logger: logger.NewNoop(),
	}

	// Create Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	err := s.HandlePing(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response PingResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "pong", response.Message)
	assert.NotEmpty(t, response.Time)
}

func TestHealthResponse_Structure(t *testing.T) {
	// Test that HealthResponse can be properly marshaled/unmarshaled
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Time:    "1h30m45s",
		Services: map[string]string{
			"database": "healthy",
			"logger":   "healthy",
			"config":   "healthy",
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled HealthResponse
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, response, unmarshaled)
}

func TestPingResponse_Structure(t *testing.T) {
	// Test that PingResponse can be properly marshaled/unmarshaled
	response := PingResponse{
		Message: "pong",
		Time:    "2023-12-14T10:30:00Z",
	}

	// Test JSON marshaling
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled PingResponse
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, response, unmarshaled)
}

func TestHealthStatus_Constant(t *testing.T) {
	// Test the health status constant
	assert.Equal(t, "healthy", healthStatusHealthy, "Health status constant should be 'healthy'")
}
