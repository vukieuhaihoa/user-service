package user

import (
	"context"

	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// GetUserByID retrieves a user from the database by their ID.
// It takes a context and an ID as input and returns the user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the user to be retrieved.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.GetUserByField(ctx, "id", id)
}
