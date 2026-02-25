package healthcheck

import "context"

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
func (s *healthCheckService) Check(ctx context.Context) (string, string, string, error) {
	if err := s.healthCheckRepo.RedisPing(ctx); err != nil {
		return RedisPingTimeout, s.serviceName, s.instanceID, err
	}

	if err := s.healthCheckRepo.DBPing(ctx); err != nil {
		return DBPingConfused, s.serviceName, s.instanceID, err
	}

	return StatusOK, s.serviceName, s.instanceID, nil
}
