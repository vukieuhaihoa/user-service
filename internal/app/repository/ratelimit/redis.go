package ratelimit

import "github.com/redis/go-redis/v9"

// redisRepo is the concrete implementation of the Repository interface using Redis as the backend.
type redisRepo struct {
	c *redis.Client
}

// NewRedisRepo creates a new instance of the Redis-based rate limit repository.
//
// Parameters:
//   - c: A pointer to a redis.Client instance that will be used for rate limit operations.
//
// Returns:
//   - Repository: An implementation of the Repository interface that uses Redis for rate limiting.
func NewRedisRepo(c *redis.Client) Repository {
	return &redisRepo{c: c}
}
