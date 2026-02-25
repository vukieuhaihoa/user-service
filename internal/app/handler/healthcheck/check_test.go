package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	service "github.com/vukieuhaihoa/user-service/internal/app/service/healthcheck"
	"github.com/vukieuhaihoa/user-service/internal/app/service/healthcheck/mocks"
)

func TestHandler_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mocks.Service

		expectedError    error
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "successful health check",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("Check", ctx).Return(service.StatusOK, "TestService", "Instance123", nil)
				return svcMock
			},

			expectedError:    nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"OK","service_name":"TestService","instance_id":"Instance123"}`,
		},
		{
			name: "failed health check",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("Check", ctx).Return("Service Unhealthy", "TestService", "Instance123", assert.AnError)
				return svcMock
			},

			expectedError:    assert.AnError,
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Service Unhealthy","service_name":"TestService","instance_id":"Instance123"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(gc)
			testHandler := &healthCheckHandler{svc: mockSvc}

			testHandler.Check(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
