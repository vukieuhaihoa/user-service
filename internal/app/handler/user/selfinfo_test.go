package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	svcMocks "github.com/vukieuhaihoa/user-service/internal/app/service/user/mocks"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
)

func TestHandler_GetProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)

		setupMockSvc func(ctx *gin.Context) *svcMocks.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful get profile",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": "de305d54-75b4-431b-adb2-eb6b9e546099",
				})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("GetUserByID", mock.Anything, "de305d54-75b4-431b-adb2-eb6b9e546099").
					Return(&model.User{
						Base: model.Base{
							ID:        "de305d54-75b4-431b-adb2-eb6b9e546099",
							CreatedAt: fixture.TestTime,
							UpdatedAt: fixture.TestTime,
						},
						Username:    "testuser",
						Email:       "testuser@example.com",
						DisplayName: "Test User",
					}, nil)
				return mockUserSvc
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"data":{"id":"de305d54-75b4-431b-adb2-eb6b9e546099","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","username":"testuser","email":"testuser@example.com","display_name":"Test User"},"message":"User profile retrieved successfully!"}`,
		},
		{
			name: "unauthenticated request",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// No userID set in context to simulate unauthenticated request
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid userID in context",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Set invalid userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": 12345, // should be a string
				})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "user not found",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": "nonexistent-user-id",
				})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("GetUserByID", mock.Anything, "nonexistent-user-id").
					Return(nil, dbutils.ErrRecordNotFoundType)
				return mockUserSvc
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "service layer error",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": "de305d54-75b4-431b-adb2-eb6b9e546099",
				})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("GetUserByID", mock.Anything, "de305d54-75b4-431b-adb2-eb6b9e546099").
					Return(nil, assert.AnError)
				return mockUserSvc
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			mockUserSvc := tc.setupMockSvc(ctx)

			userHandler := NewUserHandler(mockUserSvc)
			userHandler.GetProfile(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestHandler_UpdateProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *updateProfileRequest

		setupRequest func(ctx *gin.Context, inputRequest *updateProfileRequest)

		setupMockSvc func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful update profile",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": "de305d54-75b4-431b-adb2-eb6b9e546099",
				})
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", inputRequest.DisplayName, inputRequest.Email).
					Return(nil)
				return mockUserSvc
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"message":"Edit current user successfully!"}`,
		},
		{
			name: "unauthenticated request",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// No userID set in context to simulate unauthenticated request
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid userID in context",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Set invalid userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": 12345, // should be a string
				})
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "user not found",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": "nonexistent-user-id",
				})
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("UpdateUserByID", ctx, "nonexistent-user-id", inputRequest.DisplayName, inputRequest.Email).
					Return(dbutils.ErrRecordNotFoundType)
				return mockUserSvc
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "service layer error",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("claims", jwt.MapClaims{
					"sub": "de305d54-75b4-431b-adb2-eb6b9e546099",
				})
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", inputRequest.DisplayName, inputRequest.Email).
					Return(assert.AnError)
				return mockUserSvc
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx, tc.inputRequest)
			mockUserSvc := tc.setupMockSvc(ctx, tc.inputRequest)

			userHandler := NewUserHandler(mockUserSvc)
			userHandler.UpdateProfile(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
