package user

import (
	"context"

	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// GetUserByUsername retrieves a user from the database by their username.
// It takes a context and a username as input and returns the user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - username: The username of the user to be retrieved.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.GetUserByField(ctx, "username", username)
}
