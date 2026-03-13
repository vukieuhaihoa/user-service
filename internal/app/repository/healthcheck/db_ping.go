package healthcheck

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// DBPing checks the connectivity to the database.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - error: An error object if the ping operation fails, otherwise nil
func (h *healthCheckStorage) DBPing(ctx context.Context) error {
	s := newrelic.FromContext(ctx).StartSegment("Repo_DBPing")
	defer s.End()

	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
