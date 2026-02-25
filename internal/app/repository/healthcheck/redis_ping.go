package healthcheck

import "context"

// RedisPing checks the connectivity to the Redis server.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - error: An error object if the ping operation fails, otherwise nil
func (h *healthCheckStorage) RedisPing(ctx context.Context) error {
	return h.redisClient.Ping(ctx).Err()
}
