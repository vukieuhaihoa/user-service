package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/validators"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	svcMocks "github.com/vukieuhaihoa/user-service/internal/app/service/user/mocks"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
)

func init() {
	// Register custom validators before running tests
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password_strength", validators.PasswordStrength)
	}
}

func TestHandler_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *createUserRequest

		setupRequest func(ctx *gin.Context, inputRequest *createUserRequest)

		setupMockSvc func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful create user",

			inputRequest: &createUserRequest{
				Username:    "testuser",
				Password:    "my_SECURE_password123@",
				DisplayName: "Test User",
				Email:       "testuser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("CreateUser", ctx, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
					Return(&model.User{
						Username:    inputRequest.Username,
						DisplayName: inputRequest.DisplayName,
						Email:       inputRequest.Email,
						Base: model.Base{
							ID:        "de305d54-75b4-431b-adb2-eb6b9e546099",
							CreatedAt: fixture.TestTime,
							UpdatedAt: fixture.TestTime,
						},
					}, nil)
				return mockUserSvc
			},

			expectedCode:     http.StatusCreated,
			expectedResponse: `{"data":{"id":"de305d54-75b4-431b-adb2-eb6b9e546099","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","username":"testuser","email":"testuser@example.com","display_name":"Test User"},"message":"Register an user successfully!"}`,
		},
		{
			name: "invalid request body",

			inputRequest: &createUserRequest{
				Username:    "",
				Password:    "short",
				DisplayName: "Test User",
				Email:       "invalid-email-format",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Username is invalid (required)","Password is invalid (min)","Email is invalid (email)"]}`,
		},
		{
			name: "invalid request body - weak password",

			inputRequest: &createUserRequest{
				Username:    "testuser",
				Password:    "shortshort",
				DisplayName: "Test User",
				Email:       "testuser@gmail.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.Service {
				return svcMocks.NewService(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Password is invalid (password_strength)"]}`,
		},
		{
			name: "duplicate username or email",

			inputRequest: &createUserRequest{
				Username:    "existinguser",
				Password:    "my_SECURE_password123@",
				DisplayName: "Existing User",
				Email:       "existinguser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("CreateUser", mock.Anything, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
					Return(nil, dbutils.ErrDuplicationType)
				return mockUserSvc
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"username or email already exists"}`,
		},
		{
			name: "service layer error",

			inputRequest: &createUserRequest{
				Username:    "testuser",
				Password:    "my_SECURE_password123@",
				DisplayName: "Test User",
				Email:       "testuser@gmail.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.Service {
				mockUserSvc := svcMocks.NewService(t)
				mockUserSvc.On("CreateUser", mock.Anything, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
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

			tc.setupRequest(ctx, tc.inputRequest)
			mockUserSvc := tc.setupMockSvc(ctx, tc.inputRequest)

			userHandler := NewUserHandler(mockUserSvc)
			userHandler.CreateUser(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
