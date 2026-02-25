package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/sqldb"
	"github.com/vukieuhaihoa/user-service/internal/api"
)

func TestHealthCheckEndpoint_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Health check returns OK status",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/health-check", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"OK","service_name":"bookmark-service","instance_id":"test_instance_id_1"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark-service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     redisPkg.InitMockRedis(t),
				SqlDB:           sqldb.InitMockDB(t),
				RandomCodeGen:   nil,
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    nil,
			})

			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.JSONEq(t, tc.expectedResponseBody, respRec.Body.String())
		})
	}
}
