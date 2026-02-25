package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// healthCheckResponse represents the JSON structure returned by the health check endpoint.
type healthCheckResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// Check is a Gin framework handler that performs a health check.
// It calls the health check service to retrieve the current health status
// and returns it as a JSON response.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response
//
// Response:h
//   - 200 OK: Returns the health status as a JSON object
//
// @Summary Health Check
// @Description Performs a health check and returns the service status.
// @Tags health
// @Produce json
// @Success 200 {object} healthCheckResponse
// @Router /health-check [get]
func (h *healthCheckHandler) Check(c *gin.Context) {
	message, serviceName, instanceID, err := h.svc.Check(c)
	if err != nil {
		log.Error().Err(err).Msg("service return error when check health")
		c.JSON(http.StatusInternalServerError, healthCheckResponse{
			Message:     message,
			ServiceName: serviceName,
			InstanceID:  instanceID,
		})
		return
	}
	c.JSON(http.StatusOK, healthCheckResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
