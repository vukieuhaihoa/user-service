// Package healthcheck provides functionality to perform health checks on the Redis server and database.
// It defines the HealthCheck interface and its implementation.
// The health checks include pinging the Redis server and the database to ensure connectivity.
// This package is essential for monitoring the health of the application's dependencies.
// It uses the go-redis and gorm libraries for interacting with Redis and the database respectively.
package healthcheck

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repository defines the interface for performing health checks on the Redis server.
//
//go:generate mockery --name Repository --filename health_check_repo.go --output ./mocks
type Repository interface {
	// RedisPing checks the connectivity to the Redis server.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//
	// Returns:
	//   - error: An error object if the ping operation fails, otherwise nil
	RedisPing(ctx context.Context) error

	// DBPing checks the connectivity to the database.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//
	// Returns:
	//   - error: An error object if the ping operation fails, otherwise nil
	DBPing(ctx context.Context) error
}

type healthCheckStorage struct {
	redisClient *redis.Client
	db          *gorm.DB
}

// NewHealthCheck creates a new instance of Repository using the provided Redis client and Gorm DB.
//
// Parameters:
//   - redisClient: The Redis client used for health check operations
//   - db: The Gorm DB used for health check operations
//
// Returns:
//   - Repository: A new Repository instance
func NewHealthCheckRepository(redisClient *redis.Client, db *gorm.DB) Repository {
	return &healthCheckStorage{
		redisClient: redisClient,
		db:          db,
	}
}
