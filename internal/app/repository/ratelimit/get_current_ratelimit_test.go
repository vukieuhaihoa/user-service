package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
)

func TestRepository_GetCurrentRateLimit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectedCount int
		expectedError error
	}{
		{
			name: "successful rate limit retrieval",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Set(ctx, "testget", 5, time.Hour)
				return redisClient
			},

			expectedCount: 5,
			expectedError: nil,
		},
		{
			name: "rate limit not set for the key",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			expectedCount: 0,
			expectedError: nil,
		},
		{
			name: "failed rate limit retrieval due to closed Redis client",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedCount: -1,
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMockClient := tc.setupMock(ctx)

			rateLimitRepo := NewRedisRepo(redisMockClient)

			cnt, err := rateLimitRepo.GetCurrentRateLimit(ctx, "testget")
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedCount, cnt)
		})

	}
}
