package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	svcMocks "github.com/vukieuhaihoa/user-service/internal/app/service/user/mocks"
)

func TestUser_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *loginRequest
		setupRequest func(ctx *gin.Context, inputRequest *loginRequest)

		setupMockSvc func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful login",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "my_SECURE_password123@",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("Login", mock.Anything, inputRequest.Username, inputRequest.Password).
					Return("mocked-jwt-token", nil)
				return mockUserSvc
			},
			expectedCode:     http.StatusOK,
			expectedResponse: `{"data":"mocked-jwt-token","message":"Logged in successfully!"}`,
		},
		{
			name: "invalid request body",
			inputRequest: &loginRequest{
				Username: "",
				Password: "",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Username is invalid (required)","Password is invalid (required)"]}`,
		},
		{
			name: "password too short",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "short",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Password is invalid (gte)"]}`,
		},
		{
			name: "invalid credentials",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "wrong_password",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("Login", mock.Anything, inputRequest.Username, inputRequest.Password).
					Return("", dbutils.ErrRecordNotFoundType)
				return mockUserSvc
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"invalid username or password"}`,
		},
		{
			name: "service layer error",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "my_SECURE_password123@",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("Login", mock.Anything, inputRequest.Username, inputRequest.Password).
					Return("", assert.AnError)
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
			userHandler.Login(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
