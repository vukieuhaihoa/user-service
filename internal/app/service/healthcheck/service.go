// Package healthcheck provides functionality for performing health checks on the service.
// It defines the HealthCheck interface and its implementation.
// The health check service uses a repository to perform checks on dependencies like Redis and the database.
// This package is essential for monitoring the health status of the application.
package healthcheck

import (
	"context"

	"github.com/vukieuhaihoa/user-service/internal/app/repository/healthcheck"
)

const (
	StatusOK         = "OK"
	RedisPingTimeout = "redis: client is closed"
	DBPingConfused   = "database: database is closed"
)

// Service defines the interface for health check service.
// It provides a method to perform health checks and retrieve service status.
//
//go:generate mockery --name Service --filename health_check_service.go --output ./mocks
type Service interface {
	// Check performs a health check and returns the status, service name, and instance ID.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//
	// Returns:
	//   - string: The health status ("OK")
	//   - string: The name of the service
	//   - string: The unique instance ID of the service
	//   - error: An error if the health check fails, nil otherwise
	Check(ctx context.Context) (string, string, string, error)
}

// healthCheckService implements the Service interface and provides methods for performing health checks.
type healthCheckService struct {
	serviceName     string
	instanceID      string
	healthCheckRepo healthcheck.Repository
}

// NewHealthCheckService creates a new instance of HealthCheck service.
//
// Parameters:
//   - serviceName: The name of the service
//   - instanceID: The unique instance ID of the service
//   - healthCheckRepo: The repository used for performing health checks
//
// Returns:
//   - Service: The initialized HealthCheck service instance
func NewHealthCheckService(serviceName, instanceID string, healthCheckRepo healthcheck.Repository) Service {
	return &healthCheckService{
		serviceName:     serviceName,
		instanceID:      instanceID,
		healthCheckRepo: healthCheckRepo,
	}
}
