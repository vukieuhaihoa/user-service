package healthcheck

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
)

func TestRepository_RedisPing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectedError error
	}{
		{
			name: "successful ping",

			setupMock: func() *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			expectedError: nil,
		},
		{
			name: "failed ping",

			setupMock: func() *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMockClient := tc.setupMock()

			healthCheckRepo := NewHealthCheckRepository(redisMockClient, nil)

			err := healthCheckRepo.RedisPing(ctx)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
