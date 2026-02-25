package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// IncreaseRateLimit increments the rate limit counter for a given key and sets an expiration time.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - key: The specific key for which the rate limit counter should be incremented.
//   - exp: The expiration duration for the rate limit counter.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func (r *redisRepo) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) error {
	// Execute INCR and ExpireNX atomically in a Redis transaction to avoid race conditions.
	_, err := r.c.TxPipelined(ctx, func(p redis.Pipeliner) error {

		// Increment the rate limit counter for the given key
		p.Incr(ctx, key)

		// Set the expiration time for the rate limit counter
		p.ExpireNX(ctx, key, exp)

		return nil
	})

	return err
}
