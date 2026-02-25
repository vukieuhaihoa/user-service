package healthcheck

import (
	"context"
	"database/sql"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/user-service/internal/app/repository/healthcheck/mocks"
)

func TestService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputServiceName string
		inputInstanceID  string

		setupMockRepo func(ctx context.Context) *mocks.Repository

		expectedError       error
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
	}{
		{
			name: "Health check returns OK status",

			inputServiceName: "TestService",
			inputInstanceID:  "Instance123",

			setupMockRepo: func(ctx context.Context) *mocks.Repository {
				repoMock := mocks.NewRepository(t)
				repoMock.On("RedisPing", ctx).Return(nil)
				repoMock.On("DBPing", ctx).Return(nil)
				return repoMock
			},

			expectedError:       nil,
			expectedMessage:     StatusOK,
			expectedServiceName: "TestService",
			expectedInstanceID:  "Instance123",
		},
		{
			name: "Health Check return timeout error",

			inputServiceName: "TestService",
			inputInstanceID:  "Instance123",

			setupMockRepo: func(ctx context.Context) *mocks.Repository {
				repoMock := mocks.NewRepository(t)
				repoMock.On("RedisPing", ctx).Return(redis.ErrPoolTimeout)
				return repoMock
			},

			expectedError:       redis.ErrPoolTimeout,
			expectedMessage:     RedisPingTimeout,
			expectedServiceName: "TestService",
			expectedInstanceID:  "Instance123",
		},
		{
			name: "Health Check return database error",

			inputServiceName: "TestService",
			inputInstanceID:  "Instance123",

			setupMockRepo: func(ctx context.Context) *mocks.Repository {
				repoMock := mocks.NewRepository(t)
				repoMock.On("RedisPing", ctx).Return(nil)
				repoMock.On("DBPing", ctx).Return(sql.ErrConnDone)
				return repoMock
			},

			expectedError:       sql.ErrConnDone,
			expectedMessage:     DBPingConfused,
			expectedServiceName: "TestService",
			expectedInstanceID:  "Instance123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := tc.setupMockRepo(t.Context())

			testSvc := NewHealthCheckService(tc.inputServiceName, tc.inputInstanceID, mockRepo)

			status, serviceName, instanceID, err := testSvc.Check(t.Context())

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedMessage, status)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
		})
	}
}
