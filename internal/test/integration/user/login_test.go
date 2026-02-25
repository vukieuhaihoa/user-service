package user

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	middleware "github.com/vukieuhaihoa/bookmark-libs/middlewares"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils/mocks"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
	"github.com/vukieuhaihoa/user-service/internal/api"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
)

func TestUserEndpoint_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		setupMockJWTGenerator func(t *testing.T) *mocks.JWTGenerator

		expectedStatusCode      int
		expectedMessageResponse string
	}{
		{
			name: "successful user login",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(`{"username":"testuser001","password":"my_SECURE_password123@"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				jwtGen := mocks.NewJWTGenerator(t)
				jwtGen.On("GenerateToken", mock.Anything).Return("mocked_jwt_token", nil)
				return jwtGen
			},

			expectedStatusCode:      http.StatusOK,
			expectedMessageResponse: `"message":"Logged in successfully!"`,
		},
		{
			name: "invalid user login payload",

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				return mocks.NewJWTGenerator(t)
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(`{"username":"","password":""}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"Invalid input fields"`,
		},
		{
			name: "user login failed - invalid credentials",

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				return mocks.NewJWTGenerator(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(`{"username":"testuser001","password":"wrong_password"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"invalid username or password"`,
		},
		{
			name: "user login failed - rate limit exceeded",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "192.0.2.1")
				redisClient.Set(ctx, key, middleware.IPRateLimitMaxCount, middleware.IPRateLimitInterval)
				return redisClient
			},

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				return mocks.NewJWTGenerator(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(`{"username":"testuser001","password":"my_SECURE_password123@"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusTooManyRequests,
			expectedMessageResponse: `"error":"Too many requests. Please try again later."`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// init mock db and migrate
			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtGen := tc.setupMockJWTGenerator(t)
			redisClient := redisPkg.InitMockRedis(t)

			if tc.setupMockRedis != nil {
				redisClient = tc.setupMockRedis(ctx, redisClient)
			}

			// Initialize API engine
			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     redisClient,
				SqlDB:           db,
				RandomCodeGen:   nil,
				PasswordHashing: utils.NewPasswordHashing(),
				JWTGenerator:    jwtGen,
				JWTValidator:    nil,
			})

			// Setup test HTTP request
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessageResponse)
		})
	}
}
