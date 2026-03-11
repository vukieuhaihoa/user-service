package user

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// GetUserByID retrieves a user by their ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the user to be retrieved.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	s := newrelic.FromContext(ctx).StartSegment("Service_GetUserByID")
	defer s.End()

	return u.userRepo.GetUserByID(ctx, id)
}
