// Package healthcheck provides HTTP handlers for health check operations.
// It includes handlers for checking the health status of the application
// using the Gin web framework.
package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/user-service/internal/app/service/healthcheck"
)

// Handler defines the interface for health check HTTP handlers.
// It provides a method for handling health check requests using the Gin framework.
type Handler interface {
	// Check is a Gin framework handler that performs a health check.
	// It processes HTTP requests and returns the health status.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	Check(ctx *gin.Context)
}

// healthCheckHandler handles health check HTTP requests.
// It uses the HealthCheck service to retrieve health status information.
type healthCheckHandler struct {
	svc healthcheck.Service
}

// NewHealthCheckHandler creates a new instance of the HealthCheck handler.
// It accepts a health check service implementation and returns a handler
// that can process HTTP requests for health checks.
//
// Parameters:
//   - svc: The health check service used for retrieving health status
//
// Returns:
//   - Handler: A new health check handler instance
func NewHealthCheckHandler(svc healthcheck.Service) Handler {
	return &healthCheckHandler{svc: svc}
}
