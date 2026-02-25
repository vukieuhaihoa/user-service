package user

import (
	"context"

	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// UpdateUserByID updates a user's display name and email by their ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the user to be updated.
//   - displayName: The new display name for the user.
//   - email: The new email address for the user.
//
// Returns:
//   - *model.User: The updated user model.
//   - error: An error if the update fails, otherwise nil.
func (u *userService) UpdateUserByID(ctx context.Context, id, displayName, email string) error {
	updatedUser := &model.User{
		DisplayName: displayName,
		Email:       email,
	}

	return u.userRepo.UpdateUserByID(ctx, id, updatedUser)
}
