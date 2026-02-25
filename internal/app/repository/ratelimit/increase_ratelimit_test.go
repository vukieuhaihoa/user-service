package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
)

func TestRepository_IncreaseRateLimit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		key string
		exp time.Duration

		expectedError error
		verifyFunc    func(ctx context.Context, redisClient *redis.Client)
	}{
		{
			name: "successful rate limit increase",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			key:           "test",
			exp:           time.Minute,
			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				count, err := redisClient.Get(ctx, "test").Int()
				assert.Nil(t, err)
				assert.Equal(t, 1, count)

				ttl, err := redisClient.TTL(ctx, "test").Result()
				assert.Nil(t, err)
				assert.True(t, ttl > 0)
			},
		},
		{
			name: "rate limit increase with existing count",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Set(ctx, "test", 5, time.Minute)
				return redisClient
			},

			key:           "test",
			exp:           time.Minute,
			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				count, err := redisClient.Get(ctx, "test").Int()
				assert.Nil(t, err)
				assert.Equal(t, 6, count)

				ttl, err := redisClient.TTL(ctx, "test").Result()
				assert.Nil(t, err)
				assert.True(t, ttl > 0)
			},
		},
		{
			name: "failed rate limit increase due to closed Redis client",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			key:           "test",
			exp:           time.Minute,
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			redisMockClient := tc.setupMock(ctx)

			rateLimitRepo := NewRedisRepo(redisMockClient)

			err := rateLimitRepo.IncreaseRateLimit(ctx, tc.key, tc.exp)
			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				tc.verifyFunc(ctx, redisMockClient)
			}
		})
	}
}
