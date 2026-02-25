package ratelimit

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// GetCurrentRateLimit retrieves the current value of the rate limit counter for the specified key.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - key: The specific key for which the rate limit counter should be retrieved.
//
// Returns:
//   - int: The current value of the rate limit counter. If an error occurs, it returns -1.
//   - error: An error if the operation fails, otherwise nil.
func (r *redisRepo) GetCurrentRateLimit(ctx context.Context, key string) (int, error) {
	// Get the current value of the rate limit counter for the given key
	current, err := r.c.Get(ctx, key).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// If the key does not exist, return 0 as the current rate limit
			return 0, nil
		}
		return -1, err
	}

	return current, nil
}
