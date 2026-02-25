package user

import (
	"context"

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
	return u.userRepo.GetUserByID(ctx, id)
}
