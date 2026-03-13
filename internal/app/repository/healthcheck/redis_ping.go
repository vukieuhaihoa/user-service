package healthcheck

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// RedisPing checks the connectivity to the Redis server.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - error: An error object if the ping operation fails, otherwise nil
func (h *healthCheckStorage) RedisPing(ctx context.Context) error {
	s := newrelic.FromContext(ctx).StartSegment("Repo_RedisPing")
	defer s.End()

	return h.redisClient.Ping(ctx).Err()
}
