package healthcheck

import "context"

// DBPing checks the connectivity to the database.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - error: An error object if the ping operation fails, otherwise nil
func (h *healthCheckStorage) DBPing(ctx context.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
