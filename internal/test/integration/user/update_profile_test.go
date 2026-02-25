package user

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	middleware "github.com/vukieuhaihoa/bookmark-libs/middlewares"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils/mocks"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/user-service/internal/api"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
)

func TestUserEndpoint_UpdateProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful update user profile",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `"message":"Edit current user successfully!"`,
		},
		{
			name: "get user profile failed - invalid token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"`))
				req.Header.Set("Authorization", "Bearer invalid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "invalid_jwt_token").Return(nil, assert.AnError)
				return jwtValidator
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `"message":"Invalid token"`,
		},
		{
			name: "token does not contain user ID",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"}`))
				req.Header.Set("Authorization", "Bearer token_without_user_id")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "token_without_user_id").Return(jwt.MapClaims{}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `"message":"Unauthorized"`,
		},
		{
			name: "user not found",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"}`))
				req.Header.Set("Authorization", "Bearer user_not_found_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "user_not_found_token").Return(jwt.MapClaims{"sub": "non_existent_user_id"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `"message":"Unauthorized"`,
		},
		{
			name: "update user profile failed - invalid payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"","email":"invalid-email"}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `"message":"Invalid input fields"`,
		},
		{
			name: "update user profile failed - email already exists",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"alice@example.com"}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `"message":"email already exists"`,
		},
		{
			name: "rate limit exceeded",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "4d9326d6-980c-4c62-9709-dbc70a82cbfe")
				redisClient.Set(ctx, key, middleware.UserIDRateLimitMaxCount, middleware.UserIDRateLimitInterval)
				return redisClient
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusTooManyRequests,
			expectedResponse: `"error":"Too many requests. Please try again later."`,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtValidator := tc.setupMockJWTValidator(t)
			redisClient := redisPkg.InitMockRedis(t)

			if tc.setupMockRedis != nil {
				redisClient = tc.setupMockRedis(context.Background(), redisClient)
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
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    jwtValidator,
			})

			// Setup test HTTP request
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedResponse)
		})
	}
}
