package healthcheck

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Check performs a health check and returns the status, service name, and instance ID.
// It always returns "OK" as the status for a healthy service.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - string: The health status ("OK")
//   - string: The name of the service
//   - string: The unique instance ID of the service
//   - error: An error if the health check fails, nil otherwise
func (h *healthCheckService) Check(ctx context.Context) (string, string, string, error) {
	s := newrelic.FromContext(ctx).StartSegment("Service_Check")
	defer s.End()

	if err := h.healthCheckRepo.RedisPing(ctx); err != nil {
		return RedisPingTimeout, h.serviceName, h.instanceID, err
	}

	if err := h.healthCheckRepo.DBPing(ctx); err != nil {
		return DBPingConfused, h.serviceName, h.instanceID, err
	}

	return StatusOK, h.serviceName, h.instanceID, nil
}
