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
	middleware "github.com/vukieuhaihoa/bookmark-libs/middlewares"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
	"github.com/vukieuhaihoa/user-service/internal/api"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
)

func TestUserEndpoint_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode      int
		expectedMessageResponse string
	}{
		{
			name: "successful user registration",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(`{"username":"testuser","password":"my_SECURE_password123@","display_name":"Test User","email":"testuser@example.com"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusCreated,
			expectedMessageResponse: `message":"Register an user successfully!"`,
		},
		{
			name: "invalid user registration payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(`{"username":"","password":"weak","display_name":"Test User","email":"invalid-email"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"Invalid input fields"`,
		},
		{
			name: "register failed - rate limit exceeded",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "192.0.2.1")
				redisClient.Set(ctx, key, middleware.IPRateLimitMaxCount, middleware.IPRateLimitInterval)
				return redisClient
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(`{"username":"testuser","password":"my_SECURE_password123@","display_name":"Test User","email":"testuser@example.com"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusTooManyRequests,
			expectedMessageResponse: `"error":"Too many requests. Please try again later."`,
		},
		{
			name: "register failed - username already exists",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(`{"username":"testuser001","password":"my_SECURE_password123@","display_name":"Test User","email":"testuser@example.com"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"username or email already exists"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Initialize mock Redis client
			redisClient := redisPkg.InitMockRedis(t)

			if tc.setupMockRedis != nil {
				redisClient = tc.setupMockRedis(ctx, redisClient)
			}

			// init mock db and migrate
			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})

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
				JWTGenerator:    nil,
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
