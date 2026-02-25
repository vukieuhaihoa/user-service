package ratelimit

import (
	"context"
	"time"
)

// Repository defines the interface for managing rate limits.
// It provides methods to get the current rate limit counter and to increase the rate limit counter with an expiration time.
//
//go:generate mockery --name=Repository --output=./mocks --filename=repository.go
type Repository interface {
	// GetCurrentRateLimit retrieves the current value of the rate limit counter for the specified key.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - key: The specific key for which the rate limit counter should be retrieved.
	//
	// Returns:
	//   - int: The current value of the rate limit counter. If an error occurs, it returns -1.
	//   - error: An error if the operation fails, otherwise nil.
	GetCurrentRateLimit(ctx context.Context, key string) (int, error)

	// IncreaseRateLimit increments the rate limit counter for a given key and sets an expiration time.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - key: The specific key for which the rate limit counter should be incremented.
	//   - exp: The expiration duration for the rate limit counter.
	//
	// Returns:
	//   - error: An error if the operation fails, otherwise nil.
	IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) error
}
