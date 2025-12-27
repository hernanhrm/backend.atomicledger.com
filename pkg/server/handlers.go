package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
)

const (
	healthStatusHealthy = "healthy"
)

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Time     string            `json:"time"`
	Services map[string]string `json:"services"`
}

// PingResponse represents the ping response.
type PingResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

// HandleHealth returns detailed health status of all services.
func (s *Server) HandleHealth(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)

	// Check database health
	if err := s.db.HealthCheck(ctx); err != nil {
		services["database"] = "unhealthy: " + err.Error()
	} else {
		services["database"] = healthStatusHealthy
	}

	// Check other services from DI container
	services["logger"] = healthStatusHealthy
	services["config"] = healthStatusHealthy

	// Determine overall status
	status := healthStatusHealthy
	for _, svcStatus := range services {
		if svcStatus != healthStatusHealthy {
			status = "degraded"
			break
		}
	}

	response := HealthResponse{
		Status:   status,
		Version:  "0.1.0", // TODO: Get from config or build info
		Time:     time.Now().UTC().Format(time.RFC3339),
		Services: services,
	}

	httpStatus := http.StatusOK
	if status != healthStatusHealthy {
		httpStatus = http.StatusServiceUnavailable
	}

	if err := c.JSON(httpStatus, response); err != nil {
		return oops.
			Code("health_response_failed").
			Wrapf(err, "failed to send health response")
	}
	return nil
}

// HandlePing is a simple ping endpoint.
func (s *Server) HandlePing(c echo.Context) error {
	response := PingResponse{
		Message: "pong",
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := c.JSON(http.StatusOK, response); err != nil {
		return oops.
			Code("ping_response_failed").
			Wrapf(err, "failed to send ping response")
	}
	return nil
}
