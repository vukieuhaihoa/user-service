package infrastructure

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/nrtrace"
)

// CreateNewRelicClient initializes and returns a New Relic application client.
// It reads the configuration from environment variables and handles any errors that occur during initialization.
// Returns:
//   - *newrelic.Application: A pointer to the initialized New Relic application client
func CreateNewRelicClient() *newrelic.Application {
	nrClient, err := nrtrace.NewClient("")
	common.HandlerError(err)

	return nrClient
}
